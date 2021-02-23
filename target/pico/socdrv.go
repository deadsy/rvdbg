//-----------------------------------------------------------------------------
/*

SoC Driver

Implements the soc.Driver interface for the CPUs SoC device.

*/
//-----------------------------------------------------------------------------

package pico

import (
	"github.com/deadsy/rvdbg/cpu/arm/cm"
	"github.com/deadsy/rvdbg/soc"
)

//-----------------------------------------------------------------------------

type socDriver struct {
	dbg cm.Debug
}

func newSocDriver(dbg cm.Debug) *socDriver {
	return &socDriver{
		dbg: dbg,
	}
}

func (drv *socDriver) GetAddressSize() uint {
	return drv.dbg.GetAddressSize()
}

func (drv *socDriver) GetRegisterSize(r *soc.Register) uint {
	return 32
}

func (drv *socDriver) Rd(width, addr uint) (uint, error) {
	x, err := drv.dbg.RdMem(width, addr, 1)
	if err != nil {
		return 0, err
	}
	return x[0], nil
}

func (drv *socDriver) Wr(width, addr, val uint) error {
	return drv.dbg.WrMem(width, addr, []uint{val})
}

//-----------------------------------------------------------------------------
