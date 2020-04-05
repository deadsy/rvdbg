//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11
CSR Access

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// rdCSR reads a CSR using debug ram operations.
func rdCSR(dbg *Debug, reg, size uint) (uint64, error) {

	if size == 32 {
		dbg.cache.wr32(0, rv.InsCSRR(rv.RegS0, reg))
		dbg.cache.wr32(1, rv.InsSW(rv.RegS0, ramAddr(3), rv.RegZero))
		dbg.cache.wrResume(2)
		dbg.cache.wr32(3, 0xdeadbeef)
		dbg.cache.read(3)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return 0, err
		}
		return uint64(dbg.cache.rd32(3)), nil
	}

	if size == 64 {
		dbg.cache.wr32(0, rv.InsCSRR(rv.RegS0, reg))
		dbg.cache.wr32(1, rv.InsSD(rv.RegS0, ramAddr(4), rv.RegZero))
		dbg.cache.wrResume(2)
		dbg.cache.wr32(4, 0xcafebabe)
		dbg.cache.wr32(5, 0xdeadbeef)
		dbg.cache.read(4)
		dbg.cache.read(5)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return 0, err
		}
		return dbg.cache.rd64(4), nil
	}

	return 0, fmt.Errorf("%d-bit csr reads are not supported", size)
}

// wrCSR writes a CSR using debug ram operations.
func wrCSR(dbg *Debug, reg, size uint, val uint64) error {

	if size == 32 {
		dbg.cache.wr32(0, rv.InsLW(rv.RegS0, ramAddr(3), rv.RegZero))
		dbg.cache.wr32(1, rv.InsCSRW(reg, rv.RegS0))
		dbg.cache.wrResume(2)
		dbg.cache.wr32(3, uint32(val))
		// run the code
		return dbg.cache.flush(true)
	}

	if size == 64 {
		dbg.cache.wr32(0, rv.InsLD(rv.RegS0, ramAddr(4), rv.RegZero))
		dbg.cache.wr32(1, rv.InsCSRW(reg, rv.RegS0))
		dbg.cache.wrResume(2)
		dbg.cache.wr64(4, val)
		// run the code
		return dbg.cache.flush(true)
	}

	return fmt.Errorf("%d-bit csr writes are not supported", size)
}

//-----------------------------------------------------------------------------
