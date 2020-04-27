//-----------------------------------------------------------------------------
/*

GigaDevice gd32vf103 Flash Driver

This code implements the flash.Driver interface.

*/
//-----------------------------------------------------------------------------

package gd32vf103

import (
	"fmt"

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
	p := dev.GetPeripheral("flash")
	sectorSize := uint(1 * util.KiB)
	i := 0
	for addr := uint(p.Addr); addr < p.Addr+p.Size; addr += sectorSize {
		r = append(r, mem.NewRegion(p.Name, addr, sectorSize, &flashMeta{fmt.Sprintf("page %d", i)}))
		i++
	}
	// boot
	p = dev.GetPeripheral("boot")
	r = append(r, mem.NewRegion(p.Name, p.Addr, p.Size, &flashMeta{"boot loader area"}))
	// option
	p = dev.GetPeripheral("option")
	r = append(r, mem.NewRegion(p.Name, p.Addr, p.Size, &flashMeta{"option bytes"}))
	return r
}

//-----------------------------------------------------------------------------

// FlashDriver is a flash driver for the gd32vf103.
type FlashDriver struct {
	drv     soc.Driver
	dev     *soc.Device
	sectors []*mem.Region
}

// NewFlashDriver returns a new gd32vf103 flash driver.
func NewFlashDriver(drv soc.Driver, dev *soc.Device) *FlashDriver {
	return &FlashDriver{
		drv:     drv,
		dev:     dev,
		sectors: flashSectors(dev),
	}
}

// GetSectors returns the flash sector memory regions for the gd32vf103.
func (drv *FlashDriver) GetSectors() []*mem.Region {
	return drv.sectors
}

//-----------------------------------------------------------------------------
