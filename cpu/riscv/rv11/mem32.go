//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 RV32 (32-bit address) Memory Operations

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"errors"
	"fmt"
)

//-----------------------------------------------------------------------------
// rv32 reads

func (dbg *Debug) rv32RdMem8(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func (dbg *Debug) rv32RdMem16(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func (dbg *Debug) rv32RdMem32(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func rv32RdMem(dbg *Debug, width, addr, n uint) ([]uint, error) {
	switch width {
	case 8:
		return dbg.rv32RdMem8(addr, n)
	case 16:
		return dbg.rv32RdMem16(addr, n)
	case 32:
		return dbg.rv32RdMem32(addr, n)
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
}

//-----------------------------------------------------------------------------
// rv32 writes

func (dbg *Debug) rv32WrMem8(addr uint, val []uint) error {
	return errors.New("TODO")
}

func (dbg *Debug) rv32WrMem16(addr uint, val []uint) error {
	return errors.New("TODO")
}

func (dbg *Debug) rv32WrMem32(addr uint, val []uint) error {
	return errors.New("TODO")
}

func rv32WrMem(dbg *Debug, width, addr uint, val []uint) error {
	switch width {
	case 8:
		return dbg.rv32WrMem8(addr, val)
	case 16:
		return dbg.rv32WrMem16(addr, val)
	case 32:
		return dbg.rv32WrMem32(addr, val)
	}
	return fmt.Errorf("%d-bit memory writes are not supported", width)
}

//-----------------------------------------------------------------------------
