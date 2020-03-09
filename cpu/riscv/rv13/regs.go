//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Register Operations

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------
// control and status registers

// RdCSR reads a control and status register for the current hart.
func (dbg *Debug) RdCSR(reg, size uint) (uint64, error) {
	hi := dbg.hart[dbg.hartid]
	if reg > 0xfff {
		return 0, fmt.Errorf("csr 0x%x is invalid", reg)
	}
	if size == 0 {
		size = rv.GetCSRLength(reg, &hi.info)
	}
	return hi.rdCSR(dbg, reg, size)
}

//-----------------------------------------------------------------------------
// general purpose registers

// RdGPR reads a general purpose register.
func (dbg *Debug) RdGPR(reg, size uint) (uint64, error) {
	hi := dbg.hart[dbg.hartid]
	if reg >= uint(hi.info.Nregs) {
		return 0, fmt.Errorf("gpr%d is invalid", reg)
	}
	if size == 0 {
		size = hi.info.MXLEN
	}
	return hi.rdGPR(dbg, reg, size)
}

// WrGPR writes a general purpose register.
func (dbg *Debug) WrGPR(reg, size uint, val uint64) error {
	hi := dbg.hart[dbg.hartid]
	if reg >= uint(hi.info.Nregs) {
		return fmt.Errorf("gpr%d is invalid", reg)
	}
	if size == 0 {
		size = hi.info.MXLEN
	}
	return hi.wrGPR(dbg, reg, size, val)
}

//-----------------------------------------------------------------------------
// floating point registers

// RdFPR reads a floating point register.
func (dbg *Debug) RdFPR(reg, size uint) (uint64, error) {
	hi := dbg.hart[dbg.hartid]
	if reg >= 32 {
		return 0, fmt.Errorf("fpr%d is invalid", reg)
	}
	if size == 0 {
		size = hi.info.FLEN
	}
	return hi.rdFPR(dbg, reg, size)
}

// WrFPR writes a floating point register.
func (dbg *Debug) WrFPR(reg, size uint, val uint64) error {
	hi := dbg.hart[dbg.hartid]
	if reg >= 32 {
		return fmt.Errorf("fpr%d is invalid", reg)
	}
	if size == 0 {
		size = hi.info.FLEN
	}
	return hi.wrFPR(dbg, reg, size, val)
}

//-----------------------------------------------------------------------------
