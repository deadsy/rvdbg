//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Memory Operations
Implements the mem.Driver interface methods.

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"

	"github.com/deadsy/rvdbg/util"
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
	switch width {
	case 8:
		x, err := hi.rdMem8(dbg, addr, n)
		return util.Cast8toUint(x), err
	case 16:
		x, err := hi.rdMem16(dbg, addr, n)
		return util.Cast16toUint(x), err
	case 32:
		x, err := hi.rdMem32(dbg, addr, n)
		return util.Cast32toUint(x), err
	case 64:
		x, err := hi.rdMem64(dbg, addr, n)
		return util.Cast64toUint(x), err
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
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
	switch width {
	case 8:
		return hi.wrMem8(dbg, addr, util.CastUintto8(val))
	case 16:
		return hi.wrMem16(dbg, addr, util.CastUintto16(val))
	case 32:
		return hi.wrMem32(dbg, addr, util.CastUintto32(val))
	case 64:
		return hi.wrMem64(dbg, addr, util.CastUintto64(val))
	}
	return fmt.Errorf("%d-bit memory writes are not supported", width)
}

//-----------------------------------------------------------------------------
