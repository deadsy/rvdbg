//-----------------------------------------------------------------------------
/*

RP20xx Flash Driver

This code implements the flash.Driver interface.

*/
//-----------------------------------------------------------------------------

package rp20xx

import (
	"errors"

	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// FlashDriver is a flash driver for the gd32vf103.
type FlashDriver struct {
	drv soc.Driver
	dev *soc.Device
}

// NewFlashDriver returns a new gd32vf103 flash driver.
func NewFlashDriver(drv soc.Driver, dev *soc.Device) (*FlashDriver, error) {
	return &FlashDriver{
		drv: drv,
		dev: dev,
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
	return nil
}

// Erase erases a flash region.
func (drv *FlashDriver) Erase(r *mem.Region) error {
	return errors.New("TODO")
}

// EraseAll erases all of the device flash.
func (drv *FlashDriver) EraseAll() error {
	return errors.New("TODO")
}

// Write a flash region.
func (drv *FlashDriver) Write(r *mem.Region, buf []byte) error {
	// check arguments
	if len(buf)&3 != 0 {
		return errors.New("write buffer must be a multiple of 4 bytes")
	}
	if r.Size&3 != 0 {
		return errors.New("flash region must be a multiple of 4 bytes")
	}
	if int(r.Size) < len(buf) {
		return errors.New("flash region size is smaller than the write buffer")
	}
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
