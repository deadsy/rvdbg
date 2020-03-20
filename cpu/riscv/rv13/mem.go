//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Memory Operations
Implements the mem.Driver interface methods.

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"
)

//-----------------------------------------------------------------------------

// GetAddressSize returns the current hart's address size in bits.
func (dbg *Debug) GetAddressSize() uint {
	return dbg.hart[dbg.hartid].info.MXLEN
}

//-----------------------------------------------------------------------------
// read memory

// RdMem reads n x width-bit values from memory.
func (dbg *Debug) RdMem(width, addr, n uint) ([]uint, error) {
	if n == 0 {
		return nil, nil
	}
	hi := dbg.hart[dbg.hartid]
	if width == 64 && hi.info.MXLEN < 64 {
		return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
	}
	return hi.rdMem(dbg, width, addr, n)
}

//-----------------------------------------------------------------------------
// write memory

// WrMem writes n x width-bit values to memory.
func (dbg *Debug) WrMem(width, addr uint, val []uint) error {
	if len(val) == 0 {
		return nil
	}
	hi := dbg.hart[dbg.hartid]
	if width == 64 && hi.info.MXLEN < 64 {
		return fmt.Errorf("%d-bit memory writes are not supported", width)
	}
	return hi.wrMem(dbg, width, addr, val)
}

//-----------------------------------------------------------------------------
