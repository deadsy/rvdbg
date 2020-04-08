//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 RV64 (64-bit address) Memory Operations

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"errors"
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------
// rv64 reads

func (dbg *Debug) rv64RdMem8(addr, n uint) ([]uint, error) {
	data := make([]uint, n)
	for i := range data {
		dbg.cache.rv64Addr(addr)
		dbg.cache.wr32(1, rv.InsLB(rv.RegS1, 0, rv.RegS0))
		dbg.cache.wr32(2, rv.InsSB(rv.RegS1, ramAddr(6), rv.RegZero))
		dbg.cache.wrResume(3)
		dbg.cache.read(6)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return nil, err
		}
		data[i] = uint(uint8(dbg.cache.rd32(6)))
		addr += 1
	}
	return data, nil
}

func (dbg *Debug) rv64RdMem16(addr, n uint) ([]uint, error) {
	data := make([]uint, n)
	for i := range data {
		dbg.cache.rv64Addr(addr)
		dbg.cache.wr32(1, rv.InsLH(rv.RegS1, 0, rv.RegS0))
		dbg.cache.wr32(2, rv.InsSH(rv.RegS1, ramAddr(6), rv.RegZero))
		dbg.cache.wrResume(3)
		dbg.cache.read(6)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return nil, err
		}
		data[i] = uint(uint16(dbg.cache.rd32(6)))
		addr += 2
	}
	return data, nil
}

func (dbg *Debug) rv64RdMem32(addr, n uint) ([]uint, error) {
	data := make([]uint, n)
	for i := range data {
		dbg.cache.rv64Addr(addr)
		dbg.cache.wr32(1, rv.InsLW(rv.RegS1, 0, rv.RegS0))
		dbg.cache.wr32(2, rv.InsSW(rv.RegS1, ramAddr(6), rv.RegZero))
		dbg.cache.wrResume(3)
		dbg.cache.read(6)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return nil, err
		}
		data[i] = uint(dbg.cache.rd32(6))
		addr += 4
	}
	return data, nil
}

func (dbg *Debug) rv64RdMem64(addr, n uint) ([]uint, error) {
	data := make([]uint, n)
	for i := range data {
		dbg.cache.rv64Addr(addr)
		dbg.cache.wr32(1, rv.InsLD(rv.RegS1, 0, rv.RegS0))
		dbg.cache.wr32(2, rv.InsSD(rv.RegS1, ramAddr(6), rv.RegZero))
		dbg.cache.wrResume(3)
		dbg.cache.read(6)
		dbg.cache.read(7)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return nil, err
		}
		data[i] = uint(dbg.cache.rd64(6))
		addr += 8
	}
	return data, nil
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
	for _, v := range val {
		dbg.cache.rv64Addr(addr)
		dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
		dbg.cache.wr32(2, rv.InsSW(rv.RegS1, 0, rv.RegS0))
		dbg.cache.wrResume(3)
		dbg.cache.wr32(6, uint32(v))
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return err
		}
		addr += 4
	}
	return nil
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
