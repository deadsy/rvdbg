//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 RV64 (64-bit address) Memory Operations

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"errors"
	"fmt"
)

//-----------------------------------------------------------------------------
// rv64 reads

func (dbg *Debug) rv64RdMem8(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func (dbg *Debug) rv64RdMem16(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func (dbg *Debug) rv64RdMem32(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func (dbg *Debug) rv64RdMem64(addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func rv64RdMem(dbg *Debug, width, addr, n uint) ([]uint, error) {
	switch width {
	case 8:
		return dbg.rv64RdMem8(addr, n)
	case 16:
		return dbg.rv64RdMem16(addr, n)
	case 32:
		return dbg.rv64RdMem32(addr, n)
	case 64:
		return dbg.rv64RdMem64(addr, n)
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
}

//-----------------------------------------------------------------------------
// rv64 writes

func (dbg *Debug) rv64WrMem8(addr uint, val []uint) error {
	return errors.New("TODO")
}

func (dbg *Debug) rv64WrMem16(addr uint, val []uint) error {
	return errors.New("TODO")
}

func (dbg *Debug) rv64WrMem32(addr uint, val []uint) error {
	return errors.New("TODO")
}

func (dbg *Debug) rv64WrMem64(addr uint, val []uint) error {
	return errors.New("TODO")
}

func rv64WrMem(dbg *Debug, width, addr uint, val []uint) error {
	switch width {
	case 8:
		return dbg.rv64WrMem8(addr, val)
	case 16:
		return dbg.rv64WrMem16(addr, val)
	case 32:
		return dbg.rv64WrMem32(addr, val)
	case 64:
		return dbg.rv64WrMem64(addr, val)

	}
	return fmt.Errorf("%d-bit memory writes are not supported", width)
}

//-----------------------------------------------------------------------------
