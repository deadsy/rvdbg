//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 RV64 Memory Operations

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/util"
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
// rv64 single 8/16/32/64-bit reads.

func (dbg *Debug) rdMemSingle8(addr uint) (uint8, error) {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLB(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wr32(2, rv.InsSB(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wrResume(3)
	dbg.cache.read(6)
	// run the code
	err := dbg.cache.flush(true)
	if err != nil {
		return 0, err
	}
	return uint8(dbg.cache.rd32(6)), nil
}

func (dbg *Debug) rdMemSingle16(addr uint) (uint16, error) {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLH(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wr32(2, rv.InsSH(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wrResume(3)
	dbg.cache.read(6)
	// run the code
	err := dbg.cache.flush(true)
	if err != nil {
		return 0, err
	}
	return uint16(dbg.cache.rd32(6)), nil
}

func (dbg *Debug) rdMemSingle32(addr uint) (uint32, error) {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wr32(2, rv.InsSW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wrResume(3)
	dbg.cache.read(6)
	// run the code
	err := dbg.cache.flush(true)
	if err != nil {
		return 0, err
	}
	return dbg.cache.rd32(6), nil
}

func (dbg *Debug) rdMemSingle64(addr uint) (uint64, error) {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLD(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wr32(2, rv.InsSD(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wrResume(3)
	dbg.cache.read(6)
	dbg.cache.read(7)
	// run the code
	err := dbg.cache.flush(true)
	if err != nil {
		return 0, err
	}
	return dbg.cache.rd64(6), nil
}

//-----------------------------------------------------------------------------
// rv64 reads

func (dbg *Debug) rdMem8(addr, n uint) ([]uint, error) {
	if n == 1 {
		x, err := dbg.rdMemSingle8(addr)
		return []uint{uint(x)}, err
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return nil, err
	}
	// setup the program, do the first read.
	dbg.cache.wr32(0, rv.InsLB(rv.RegS1, 0, rv.RegT0))
	dbg.cache.wr32(1, rv.InsSB(rv.RegS1, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 1))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(true)
	if err != nil {
		return nil, err
	}
	// perform the reads
	data, err := dbg.rdOps32(4, n, util.Mask8)
	if err != nil {
		return nil, err
	}
	// restore t0
	return data, dbg.WrGPR(rv.RegT0, 0, saved)
}

func (dbg *Debug) rdMem16(addr, n uint) ([]uint, error) {
	if n == 1 {
		x, err := dbg.rdMemSingle16(addr)
		return []uint{uint(x)}, err
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return nil, err
	}
	// setup the program
	dbg.cache.wr32(0, rv.InsLH(rv.RegS1, 0, rv.RegT0))
	dbg.cache.wr32(1, rv.InsSH(rv.RegS1, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 2))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(true)
	if err != nil {
		return nil, err
	}
	// perform the reads, do the first read.
	data, err := dbg.rdOps32(4, n, util.Mask16)
	if err != nil {
		return nil, err
	}
	// restore t0
	return data, dbg.WrGPR(rv.RegT0, 0, saved)
}

func (dbg *Debug) rdMem32(addr, n uint) ([]uint, error) {
	if n == 1 {
		x, err := dbg.rdMemSingle32(addr)
		return []uint{uint(x)}, err
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return nil, err
	}
	// setup the program, do the first read.
	dbg.cache.wr32(0, rv.InsLW(rv.RegS1, 0, rv.RegT0))
	dbg.cache.wr32(1, rv.InsSW(rv.RegS1, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 4))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(true)
	if err != nil {
		return nil, err
	}
	// perform the reads
	data, err := dbg.rdOps32(4, n, util.Mask32)
	if err != nil {
		return nil, err
	}
	// restore t0
	return data, dbg.WrGPR(rv.RegT0, 0, saved)
}

func (dbg *Debug) rdMem64(addr, n uint) ([]uint, error) {
	if n == 1 {
		x, err := dbg.rdMemSingle64(addr)
		return []uint{uint(x)}, err
	}
	// save t0 and load the address in t0
	saved, err := dbg.rwReg(rv.RegT0, uint64(addr))
	if err != nil {
		return nil, err
	}
	// setup the program, do the first read.
	dbg.cache.wr32(0, rv.InsLD(rv.RegS1, 0, rv.RegT0))
	dbg.cache.wr32(1, rv.InsSD(rv.RegS1, ramAddr(4), rv.RegZero))
	dbg.cache.wr32(2, rv.InsADDI(rv.RegT0, rv.RegT0, 8))
	dbg.cache.wrResume(3)
	err = dbg.cache.flush(true)
	if err != nil {
		return nil, err
	}
	// perform the reads
	data, err := dbg.rdOps64(4, n)
	if err != nil {
		return nil, err
	}
	// restore t0
	return data, dbg.WrGPR(rv.RegT0, 0, saved)
}

func rdMem(dbg *Debug, width, addr, n uint) ([]uint, error) {
	switch width {
	case 8:
		return dbg.rdMem8(addr, n)
	case 16:
		return dbg.rdMem16(addr, n)
	case 32:
		return dbg.rdMem32(addr, n)
	case 64:
		return dbg.rdMem64(addr, n)
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
}

//-----------------------------------------------------------------------------
// rv64 single 8/16/32/64-bit writes

func (dbg *Debug) wrMemSingle8(addr uint, val uint8) error {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSB(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr32(6, uint32(val))
	// run the code
	return dbg.cache.flush(true)
}

func (dbg *Debug) wrMemSingle16(addr uint, val uint16) error {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSH(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr32(6, uint32(val))
	// run the code
	return dbg.cache.flush(true)
}

func (dbg *Debug) wrMemSingle32(addr uint, val uint32) error {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLW(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSW(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr32(6, val)
	// run the code
	return dbg.cache.flush(true)
}

func (dbg *Debug) wrMemSingle64(addr uint, val uint64) error {
	dbg.cache.setAddr(addr)
	dbg.cache.wr32(1, rv.InsLD(rv.RegS1, ramAddr(6), rv.RegZero))
	dbg.cache.wr32(2, rv.InsSD(rv.RegS1, 0, rv.RegS0))
	dbg.cache.wrResume(3)
	dbg.cache.wr64(6, val)
	// run the code
	return dbg.cache.flush(true)
}

//-----------------------------------------------------------------------------
// rv64 writes

func (dbg *Debug) wrMem8(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.wrMemSingle8(addr, uint8(val[0]))
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

func (dbg *Debug) wrMem16(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.wrMemSingle16(addr, uint16(val[0]))
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

func (dbg *Debug) wrMem32(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.wrMemSingle32(addr, uint32(val[0]))
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

func (dbg *Debug) wrMem64(addr uint, val []uint) error {
	// do a single write
	if len(val) == 1 {
		return dbg.wrMemSingle64(addr, uint64(val[0]))
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

func wrMem(dbg *Debug, width, addr uint, val []uint) error {
	switch width {
	case 8:
		return dbg.wrMem8(addr, val)
	case 16:
		return dbg.wrMem16(addr, val)
	case 32:
		return dbg.wrMem32(addr, val)
	case 64:
		return dbg.wrMem64(addr, val)
	}
	return fmt.Errorf("%d-bit memory writes are not supported", width)
}

//-----------------------------------------------------------------------------
// implement the mem.Driver interface methods

// GetAddressSize returns the current hart's address size in bits.
func (dbg *Debug) GetAddressSize() uint {
	return dbg.hart[dbg.hartid].info.MXLEN
}

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
