//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11
floating pont register access

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// rdFPR reads a FPR using debug ram operations.
func rdFPR(dbg *Debug, reg, size uint) (uint64, error) {

	if size == 32 {
		dbg.cache.wr(0, rv.InsFSW(reg, ramAddr(0), rv.RegZero))
		dbg.cache.wrResume(1)
		dbg.cache.read(0)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return 0, err
		}
		return uint64(dbg.cache.rd(0)), nil
	}

	if size == 64 {
		dbg.cache.wr(0, rv.InsFSD(reg, ramAddr(2), rv.RegZero))
		dbg.cache.wrResume(1)
		dbg.cache.read(2)
		dbg.cache.read(3)
		// run the code
		err := dbg.cache.flush(true)
		if err != nil {
			return 0, err
		}
		lx := uint64(dbg.cache.rd(2))
		ux := uint64(dbg.cache.rd(3))
		return (ux << 32) | lx, nil
	}

	return 0, fmt.Errorf("%d-bit fpr reads are not supported", size)
}

// wrFPR writes a FPR using debug ram operations.
func wrFPR(dbg *Debug, reg, size uint, val uint64) error {

	if size == 32 {
		dbg.cache.wr(0, rv.InsFLW(reg, ramAddr(2), rv.RegZero))
		dbg.cache.wrResume(1)
		dbg.cache.wr(2, uint32(val))
		// run the code
		return dbg.cache.flush(true)
	}

	if size == 64 {
		dbg.cache.wr(0, rv.InsFLD(reg, ramAddr(2), rv.RegZero))
		dbg.cache.wrResume(1)
		dbg.cache.wr(2, uint32(val))
		dbg.cache.wr(3, uint32(val>>32))
		// run the code
		return dbg.cache.flush(true)
	}

	return fmt.Errorf("%d-bit fpr writes are not supported", size)
}

//-----------------------------------------------------------------------------
