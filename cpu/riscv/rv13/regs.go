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

// rdReg32 reads a 32-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) rdReg32(reg uint) (uint32, error) {
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(reg, size32, false, false, true, false)),
		// readback the command status
		dmiRd(abstractcs),
		// done
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	err = dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
	if err != nil {
		return 0, err
	}
	return dbg.rdData32()
}

// rdReg64 reads a 64-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) rdReg64(reg uint) (uint64, error) {
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(reg, size64, false, false, true, false)),
		// readback the command status
		dmiRd(abstractcs),
		// done
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	err = dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
	if err != nil {
		return 0, err
	}
	return dbg.rdData64()
}

// rdReg128 reads a 128-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) rdReg128(reg uint) (uint64, uint64, error) {
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(reg, size128, false, false, true, false)),
		// readback the command status
		dmiRd(abstractcs),
		// done
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, 0, err
	}
	err = dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
	if err != nil {
		return 0, 0, err
	}
	return dbg.rdData128()
}

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
		x, err = dbg.rdReg32(regCSR(reg))
		val = uint(x)
	case 64:
		var x uint64
		x, err = dbg.rdReg64(regCSR(reg))
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

// RdGPR reads a general purpose register.
func (dbg *Debug) RdGPR(reg uint) (uint64, error) {
	hi := dbg.GetCurrentHart()
	if reg >= uint(hi.Nregs) {
		return 0, fmt.Errorf("gpr %d is invalid", reg)
	}
	var err error
	var val uint64
	size := hi.MXLEN
	switch size {
	case 32:
		var x uint32
		x, err = dbg.rdReg32(regGPR(reg))
		val = uint64(x)
	case 64:
		val, err = dbg.rdReg64(regGPR(reg))
	default:
		return 0, fmt.Errorf("%d-bit gpr read not supported", size)
	}
	return val, err
}

// WrGPR writes a general purpose register.
func (dbg *Debug) WrGPR(reg uint, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
// floating point registers

// RdFPR reads a floating point register.
func (dbg *Debug) RdFPR(reg uint) (uint64, error) {
	if reg >= 32 {
		return 0, fmt.Errorf("fpr %d is invalid", reg)
	}
	var err error
	var val uint64
	size := dbg.GetCurrentHart().FLEN
	switch size {
	case 32:
		var x uint32
		x, err = dbg.rdReg32(regGPR(reg))
		val = uint64(x)
	case 64:
		val, err = dbg.rdReg64(regGPR(reg))
	default:
		return 0, fmt.Errorf("%d-bit gpr read not supported", size)
	}
	return val, err
}

// WrFPR writes a floating point register.
func (dbg *Debug) WrFPR(reg uint, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
