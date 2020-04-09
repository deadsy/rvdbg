//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 RV64 (64-bit address) Memory Operations

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// rwReg reads a value from and writes a value to a register.
func (dbg *Debug) rwReg(reg uint, val uint64) (uint64, error) {
	old, err := dbg.RdGPR(reg, 0)
	if err != nil {
		return 0, err
	}
	err = dbg.WrGPR(reg, 0, val)
	return old, err
}

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
// rv64 single 8/16/32/64-bit writes

func (dbg *Debug) rv64WrMemSingle8(addr uint, val uint8) error {
	dbg.cache.rv64Addr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSB(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr32(6, uint32(val))
	// run the code
	return dbg.cache.flush(true)
}

func (dbg *Debug) rv64WrMemSingle16(addr uint, val uint16) error {
	dbg.cache.rv64Addr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSH(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr32(6, uint32(val))
	// run the code
	return dbg.cache.flush(true)
}

func (dbg *Debug) rv64WrMemSingle32(addr uint, val uint32) error {
	dbg.cache.rv64Addr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSW(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr32(6, val)
	// run the code
	return dbg.cache.flush(true)
}

func (dbg *Debug) rv64WrMemSingle64(addr uint, val uint64) error {
	dbg.cache.rv64Addr(addr)
	dbg.cache.wr32(1, rv.InsLD(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSD(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr64(6, val)
	// run the code
	return dbg.cache.flush(true)
}

//-----------------------------------------------------------------------------
// rv64 writes

func (dbg *Debug) rv64WrMem8(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.rv64WrMemSingle8(addr, uint8(val[0]))
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return err
	}
	// setup the program
	dbg.cache.wr32(0, rv.InsLW(rv.RegS0, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(1, rv.InsSB(rv.RegS0, 0, rv.RegT0))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 1))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(false)
	if err != nil {
		return err
	}
	// perform the writes
	err = dbg.wrOps32(4, val)
	if err != nil {
		return err
	}
	// restore t0
	return dbg.WrGPR(rv.RegT0, 0, saved)
}

func (dbg *Debug) rv64WrMem16(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.rv64WrMemSingle16(addr, uint16(val[0]))
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return err
	}
	// setup the program
	dbg.cache.wr32(0, rv.InsLW(rv.RegS0, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(1, rv.InsSH(rv.RegS0, 0, rv.RegT0))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 2))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(false)
	if err != nil {
		return err
	}
	// perform the writes
	err = dbg.wrOps32(4, val)
	if err != nil {
		return err
	}
	// restore t0
	return dbg.WrGPR(rv.RegT0, 0, saved)
}

func (dbg *Debug) rv64WrMem32(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.rv64WrMemSingle32(addr, uint32(val[0]))
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return err
	}
	// setup the program
	dbg.cache.wr32(0, rv.InsLW(rv.RegS0, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(1, rv.InsSW(rv.RegS0, 0, rv.RegT0))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 4))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(false)
	if err != nil {
		return err
	}
	// perform the writes
	err = dbg.wrOps32(4, val)
	if err != nil {
		return err
	}
	// restore t0
	return dbg.WrGPR(rv.RegT0, 0, uint64(saved))
}

func (dbg *Debug) rv64WrMem64(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.rv64WrMemSingle64(addr, uint64(val[0]))
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return err
	}
	// setup the program
	dbg.cache.wr32(0, rv.InsLD(rv.RegS0, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(1, rv.InsSD(rv.RegS0, 0, rv.RegT0))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 8))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(false)
	if err != nil {
		return err
	}
	// perform the writes
	err = dbg.wrOps64(4, val)
	if err != nil {
		return err
	}
	// restore t0
	return dbg.WrGPR(rv.RegT0, 0, uint64(saved))
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
