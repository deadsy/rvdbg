//-----------------------------------------------------------------------------
/*

CSR Driver

Implements the soc.Driver interface for the CPUs control and status registers.

*/
//-----------------------------------------------------------------------------

package maixgo

import (
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/soc"
)

//-----------------------------------------------------------------------------

type csrDriver struct {
	dbg rv.Debug
}

func newCsrDriver(dbg rv.Debug) *csrDriver {
	return &csrDriver{
		dbg: dbg,
	}
}

func (drv *csrDriver) GetAddressSize() uint {
	// 12-bits for the CSR register number.
	return 12
}

func (drv *csrDriver) GetRegisterSize(r *soc.Register) uint {
	return rv.GetCSRSize(r.Offset, drv.dbg.GetCurrentHart())
}

func (drv *csrDriver) Rd(width, addr uint) (uint, error) {
	val, err := drv.dbg.RdCSR(addr, width)
	return uint(val), err
}

//-----------------------------------------------------------------------------
