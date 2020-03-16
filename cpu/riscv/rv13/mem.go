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

// RdMem8 reads n x 8-bit values from memory.
func (dbg *Debug) RdMem8(addr, n uint) ([]uint8, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem8(dbg, addr, n)
}

// RdMem16 reads n x 16-bit values from memory.
func (dbg *Debug) RdMem16(addr, n uint) ([]uint16, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem16(dbg, addr, n)
}

// RdMem32 reads n x 32-bit values from memory.
func (dbg *Debug) RdMem32(addr, n uint) ([]uint32, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem32(dbg, addr, n)
}

// RdMem64 reads n x 64-bit values from memory.
func (dbg *Debug) RdMem64(addr, n uint) ([]uint64, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem64(dbg, addr, n)
}

// RdMem reads n x width-bit values from memory.
func (dbg *Debug) RdMem(width, addr, n uint) ([]uint, error) {
	switch width {
	case 8:
		x, err := dbg.RdMem8(addr, n)
		return util.Convert8toUint(x), err
	case 16:
		x, err := dbg.RdMem16(addr, n)
		return util.Convert16toUint(x), err
	case 32:
		x, err := dbg.RdMem32(addr, n)
		return util.Convert32toUint(x), err
	case 64:
		x, err := dbg.RdMem64(addr, n)
		return util.Convert64toUint(x), err
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
}

//-----------------------------------------------------------------------------
// write memory

// WrMem8 writes n x 8-bit values to memory.
func (dbg *Debug) WrMem8(addr uint, val []uint8) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem8(dbg, addr, val)
}

// WrMem16 writes n x 16-bit values to memory.
func (dbg *Debug) WrMem16(addr uint, val []uint16) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem16(dbg, addr, val)
}

// WrMem32 writes n x 32-bit values to memory.
func (dbg *Debug) WrMem32(addr uint, val []uint32) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem32(dbg, addr, val)
}

// WrMem64 writes n x 64-bit values to memory.
func (dbg *Debug) WrMem64(addr uint, val []uint64) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem64(dbg, addr, val)
}

// WrMem writes n x width-bit values to memory.
func (dbg *Debug) WrMem(width, addr uint, val []uint) error {
	switch width {
	case 8:
		return dbg.WrMem8(addr, util.ConvertUintto8(val))
	case 16:
		return dbg.WrMem16(addr, util.ConvertUintto16(val))
	case 32:
		return dbg.WrMem32(addr, util.ConvertUintto32(val))
	case 64:
		return dbg.WrMem64(addr, util.ConvertUintto64(val))
	}
	return fmt.Errorf("%d-bit memory writes are not supported", width)
}

//-----------------------------------------------------------------------------
