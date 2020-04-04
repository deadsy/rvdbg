//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11
CSR Access

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"errors"
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// rdCSR reads a CSR using debug ram operations.
func rdCSR(dbg *Debug, reg, size uint) (uint64, error) {

	if size == 32 {
		dbg.cache.wr(0, rv.InsCSRR(rv.RegS0, reg))
		dbg.cache.wr(1, rv.InsSW(rv.RegS0, ramAddr(0), rv.RegZero))
		dbg.cache.wrResume(2)
		dbg.cache.read(0)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return 0, err
		}
		return uint64(dbg.cache.rd(0)), nil
	}

	if size == 64 {
		dbg.cache.wr(0, rv.InsCSRR(rv.RegS0, reg))
		dbg.cache.wr(1, rv.InsSD(rv.RegS0, ramAddr(0), rv.RegZero))
		dbg.cache.wrResume(2)
		dbg.cache.read(0)
		dbg.cache.read(1)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return 0, err
		}
		lx := uint64(dbg.cache.rd(0))
		ux := uint64(dbg.cache.rd(1))
		return (ux << 32) | lx, nil
	}

	return 0, fmt.Errorf("%d-bit csr reads are not supported", size)
}

// wrCSR writes a CSR using debug ram operations.
func wrCSR(dbg *Debug, reg, size uint, val uint64) error {
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
