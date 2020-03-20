//-----------------------------------------------------------------------------
/*

SoC Driver

Implements the soc.Driver interface for the CPUs SoC device.

*/
//-----------------------------------------------------------------------------

package redv

import (
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/soc"
)

//-----------------------------------------------------------------------------

type socDriver struct {
	dbg rv.Debug
}

func newSocDriver(dbg rv.Debug) *socDriver {
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

//-----------------------------------------------------------------------------
