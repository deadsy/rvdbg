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
func (dbg *Debug) RdCSR(reg uint) (uint, error) {
	if reg > 0xfff {
		return 0, fmt.Errorf("csr 0x%x is invalid", reg)
	}
	var err error
	var val uint
	size := rv.GetCSRLength(reg, dbg.GetCurrentHart())
	switch size {
	case 32:
		var x uint32
		x, err = dbg.acRd32(regCSR(reg))
		val = uint(x)
	case 64:
		var x uint64
		x, err = dbg.acRd64(regCSR(reg))
		val = uint(x)
	default:
		return 0, fmt.Errorf("%d-bit csr read not supported", size)
	}
	return val, err
}

// WrCSR writes a control and status register.
func (dbg *Debug) WrCSR(reg, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
// general purpose registers

// rdGPR reads a sized general purpose register.
func (dbg *Debug) rdGPR(reg uint, size int) (uint64, error) {
	var err error
	var val uint64
	switch size {
	case 32:
		var x uint32
		x, err = dbg.acRd32(regGPR(reg))
		val = uint64(x)
	case 64:
		val, err = dbg.acRd64(regGPR(reg))
	default:
		return 0, fmt.Errorf("%d-bit gpr read not supported", size)
	}
	return val, err
}

// RdGPR reads a general purpose register.
func (dbg *Debug) RdGPR(reg uint) (uint64, error) {
	hi := dbg.GetCurrentHart()
	if reg >= uint(hi.Nregs) {
		return 0, fmt.Errorf("gpr%d is invalid", reg)
	}
	return dbg.rdGPR(reg, hi.MXLEN)
}

// WrGPR writes a general purpose register.
func (dbg *Debug) WrGPR(reg uint, val uint64) error {
	hi := dbg.GetCurrentHart()
	if reg >= uint(hi.Nregs) {
		return fmt.Errorf("gpr%d is invalid", reg)
	}
	size := hi.MXLEN
	switch size {
	case 32:
		return dbg.acWr32(regGPR(reg), uint32(val))
	case 64:
		return dbg.acWr64(regGPR(reg), val)
	}
	return fmt.Errorf("%d-bit gpr write not supported", size)
}

//-----------------------------------------------------------------------------
// floating point registers

// rdFPR reads a sized floating point register.
func (dbg *Debug) rdFPR(reg uint, size int) (uint64, error) {
	var err error
	var val uint64
	switch size {
	case 32:
		var x uint32
		x, err = dbg.acRd32(regFPR(reg))
		val = uint64(x)
	case 64:
		val, err = dbg.acRd64(regFPR(reg))
	default:
		return 0, fmt.Errorf("%d-bit fpr read not supported", size)
	}
	return val, err
}

// RdFPR reads a floating point register.
func (dbg *Debug) RdFPR(reg uint) (uint64, error) {
	if reg >= 32 {
		return 0, fmt.Errorf("fpr%d is invalid", reg)
	}
	return dbg.rdFPR(reg, dbg.GetCurrentHart().FLEN)
}

// WrFPR writes a floating point register.
func (dbg *Debug) WrFPR(reg uint, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
