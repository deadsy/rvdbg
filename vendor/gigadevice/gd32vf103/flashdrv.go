//-----------------------------------------------------------------------------
/*

GigaDevice gd32vf103 Flash Driver

This code implements the flash.Driver interface.

*/
//-----------------------------------------------------------------------------

package gd32vf103

import (
	"errors"
	"fmt"
	"time"

	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

type flashMeta struct {
	name string
}

func (m *flashMeta) String() string {
	return m.name
}

//-----------------------------------------------------------------------------

// flashSectors returns a set of flash sectors for the device.
func flashSectors(dev *soc.Device) []*mem.Region {
	r := []*mem.Region{}
	// main flash
	p, _ := dev.GetPeripheral("flash")
	sectorSize := uint(1 * util.KiB)
	i := 0
	for addr := uint(p.Addr); addr < p.Addr+p.Size; addr += sectorSize {
		r = append(r, mem.NewRegion(p.Name, addr, sectorSize, &flashMeta{fmt.Sprintf("page %d", i)}))
		i++
	}
	// boot
	p, _ = dev.GetPeripheral("boot")
	r = append(r, mem.NewRegion(p.Name, p.Addr, p.Size, &flashMeta{"boot loader area"}))
	// option
	p, _ = dev.GetPeripheral("option")
	r = append(r, mem.NewRegion(p.Name, p.Addr, p.Size, &flashMeta{"option bytes"}))
	return r
}

//-----------------------------------------------------------------------------

// control register bits
const (
	ctlENDIE = (1 << 12) // End of operation interrupt enable bit
	ctlERRIE = (1 << 10) // Error interrupt enable bit
	ctlOBWEN = (1 << 9)  // Option byte erase/program enable bit
	ctlLK    = (1 << 7)  // FMC_CTL0 lock bit
	ctlSTART = (1 << 6)  // Send erase command to FMC bit
	ctlOBER  = (1 << 5)  // Option bytes erase command bit
	ctlOBPG  = (1 << 4)  // Option bytes program command bit
	ctlMER   = (1 << 2)  // Main flash mass erase for bank0 command bit
	ctlPER   = (1 << 1)  // Main flash page erase for bank0 command bit
	ctlPG    = (1 << 0)  // Main flash program for bank0 command bit
)

// status register bits
const (
	statENDF  = (1 << 5) // End of operation flag bit
	statWPERR = (1 << 4) // Erase/Program protection error flag bit
	statPGERR = (1 << 2) // Program error flag bit
	statBUSY  = (1 << 0) // The flash is busy bit
)

func (drv *FlashDriver) wrCtl(val uint) error {
	return drv.fmc.Wr(drv.drv, "CTL0", val)
}

func (drv *FlashDriver) rdCtl() (uint, error) {
	return drv.fmc.Rd(drv.drv, "CTL0")
}

func (drv *FlashDriver) setCtl(bits uint) error {
	return drv.fmc.Set(drv.drv, "CTL0", bits)
}

func (drv *FlashDriver) clrCtl(bits uint) error {
	return drv.fmc.Clr(drv.drv, "CTL0", bits)
}

func (drv *FlashDriver) wrKey(val uint) error {
	return drv.fmc.Wr(drv.drv, "KEY0", val)
}

func (drv *FlashDriver) wrStat(val uint) error {
	return drv.fmc.Wr(drv.drv, "STAT0", val)
}

func (drv *FlashDriver) rdStat() (uint, error) {
	return drv.fmc.Rd(drv.drv, "STAT0")
}

func (drv *FlashDriver) wrAddr(val uint) error {
	return drv.fmc.Wr(drv.drv, "ADDR0", val)
}

/*

WS     : 40022000[31:0] = 0x00000030 wait state counter register
KEY0   : 40022004[31:0] = 0          Unlock key register 0
OBKEY  : 40022008[31:0] = 0          Option byte unlock key register
STAT0  : 4002200c[31:0] = 0          Status register 0
CTL0   : 40022010[31:0] = 0x00000080 Control register 0
ADDR0  : 40022014[31:0] = 0          Address register 0
OBSTAT : 4002201c[31:0] = 0x03fffffc Option byte status register
WP     : 40022020[31:0] = 0xffffffff Erase/Program Protection register
PID    : 40022100[31:0] = 0x4a425633 Product ID register

*/

//-----------------------------------------------------------------------------

// unlock the flash
func (drv *FlashDriver) unlock() error {
	ctl, err := drv.rdCtl()
	if err != nil {
		return err
	}
	if ctl&ctlLK == 0 {
		// already unlocked
		return nil
	}
	// write the unlock sequence
	err = drv.wrKey(0x45670123)
	if err != nil {
		return err
	}
	err = drv.wrKey(0xCDEF89AB)
	if err != nil {
		return err
	}
	// clear any set CR bits
	err = drv.wrCtl(0)
	if err != nil {
		return err
	}
	return nil
}

// lock the flash
func (drv *FlashDriver) lock() error {
	return drv.setCtl(ctlLK)
}

// wait for flash operation completion
func (drv *FlashDriver) wait() error {
	const pollMax = 5
	const pollTime = 100 * time.Millisecond
	var stat uint
	var err error
	i := 0
	for i = 0; i < pollMax; i++ {
		stat, err = drv.rdStat()
		if err != nil {
			return err
		}
		if stat&statBUSY == 0 {
			break
		}
		time.Sleep(pollTime)
	}
	// clear status bits
	err = drv.wrStat(statENDF | statWPERR | statPGERR)
	if err != nil {
		return err
	}
	// check for errors
	if i >= pollMax {
		return errors.New("timeout")
	}
	return checkErrors(stat)
}

func checkErrors(stat uint) error {
	if stat&statWPERR != 0 {
		return errors.New("write protect error")
	}
	if stat&statPGERR != 0 {
		return errors.New("programming error")
	}
	return nil
}

//-----------------------------------------------------------------------------

// FlashDriver is a flash driver for the gd32vf103.
type FlashDriver struct {
	drv     soc.Driver
	dev     *soc.Device
	fmc     *soc.Peripheral
	sectors []*mem.Region
}

// NewFlashDriver returns a new gd32vf103 flash driver.
func NewFlashDriver(drv soc.Driver, dev *soc.Device) (*FlashDriver, error) {
	fmc, err := dev.GetPeripheral("FMC")
	if err != nil {
		return nil, err
	}
	return &FlashDriver{
		drv:     drv,
		dev:     dev,
		fmc:     fmc,
		sectors: flashSectors(dev),
	}, nil
}

// GetAddressSize returns the address size in bits.
func (drv *FlashDriver) GetAddressSize() uint {
	return 32
}

// GetDefaultRegion returns a default memory region.
func (drv *FlashDriver) GetDefaultRegion() *mem.Region {
	return mem.NewRegion("", 0, 1*util.KiB, nil)
}

// LookupSymbol returns an address and size for a symbol.
func (drv *FlashDriver) LookupSymbol(name string) *mem.Region {
	p, err := drv.dev.GetPeripheral(name)
	if err != nil {
		return nil
	}
	return mem.NewRegion(name, p.Addr, p.Size, nil)
}

// GetSectors returns the flash memory regions.
func (drv *FlashDriver) GetSectors() []*mem.Region {
	return drv.sectors
}

// Erase erases a flash region.
func (drv *FlashDriver) Erase(r *mem.Region) error {
	// make sure the flash is not busy
	err := drv.wait()
	if err != nil {
		return err
	}
	// unlock the flash
	err = drv.unlock()
	if err != nil {
		return err
	}
	// set the page erase bit
	err = drv.setCtl(ctlPER)
	if err != nil {
		return err
	}
	// set the page address
	err = drv.wrAddr(r.GetAddr())
	if err != nil {
		return err
	}
	// set the start bit
	err = drv.setCtl(ctlSTART)
	if err != nil {
		return err
	}
	// wait for completion
	err = drv.wait()
	if err != nil {
		return err
	}
	// clear the page erase bit
	err = drv.clrCtl(ctlPER)
	if err != nil {
		return err
	}
	// lock the flash
	return drv.lock()
}

// EraseAll erases all of the device flash.
func (drv *FlashDriver) EraseAll() error {
	// make sure the flash is not busy
	err := drv.wait()
	if err != nil {
		return err
	}
	// unlock the flash
	err = drv.unlock()
	if err != nil {
		return err
	}
	// set the mass erase bit
	err = drv.setCtl(ctlMER)
	if err != nil {
		return err
	}
	// set the start bit
	err = drv.setCtl(ctlSTART)
	if err != nil {
		return err
	}
	// wait for completion
	err = drv.wait()
	if err != nil {
		return err
	}
	// clear the mass erase bit
	err = drv.clrCtl(ctlMER)
	if err != nil {
		return err
	}
	// lock the flash
	return drv.lock()
}

// Write a flash region.
func (drv *FlashDriver) Write(r *mem.Region, buf []byte) error {
	// check arguments
	if len(buf)&3 != 0 {
		return errors.New("write buffer must be a multiple of 4 bytes")
	}
	if r.GetSize()&3 != 0 {
		return errors.New("flash region must be a multiple of 4 bytes")
	}
	if int(r.GetSize()) < len(buf) {
		return errors.New("flash region size is smaller than the write buffer")
	}

	// make sure the flash is not busy
	err := drv.wait()
	if err != nil {
		return err
	}

	// unlock the flash
	err = drv.unlock()
	if err != nil {
		return err
	}

	// convert the write buffer to 32-bit values
	buf32 := make([]uint, len(buf)>>2)
	util.ConvertFromUint8(32, buf, buf32)

	// write to flash
	addr := r.GetAddr()
	for i := 0; i < len(buf32); i++ {
		// set the program bit
		err = drv.setCtl(ctlPG)
		if err != nil {
			return err
		}
		// write to flash
		err = drv.drv.Wr(32, addr, buf32[i])
		if err != nil {
			return err
		}
		addr += 4
		// wait for completion
		err := drv.wait()
		if err != nil {
			return err
		}
	}

	// lock the flash
	return drv.lock()
}

//-----------------------------------------------------------------------------
