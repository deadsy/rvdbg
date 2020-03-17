//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Abstract Command Operations

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"
)

//-----------------------------------------------------------------------------
// abstract command read operations

// acRd32 reads a 32-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd32(reg uint) (uint32, error) {
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(reg, size32, cmdRead)),
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

// acRd64 reads a 64-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd64(reg uint) (uint64, error) {
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(reg, size64, cmdRead)),
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

// acRd128 reads a 128-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd128(reg uint) (uint64, uint64, error) {
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(reg, size128, cmdRead)),
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
// abstract command write operations

// acWr32 writes a 32-bit GPR/FPR/CSR using an abstract register write command.
func (dbg *Debug) acWr32(reg uint, val uint32) error {
	ops := []dmiOp{
		// write val
		dmiWr(data0, val),
		// write the register
		dmiWr(command, cmdRegister(reg, size32, cmdWrite)),
		// readback the command status
		dmiRd(abstractcs),
		// done
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return err
	}
	return dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
}

// acWr64 writes a 64-bit GPR/FPR/CSR using an abstract register write command.
func (dbg *Debug) acWr64(reg uint, val uint64) error {
	ops := []dmiOp{
		// write val
		dmiWr(data0, uint32(val)),
		dmiWr(data0+1, uint32(val>>32)),
		// write the register
		dmiWr(command, cmdRegister(reg, size64, cmdWrite)),
		// readback the command status
		dmiRd(abstractcs),
		// done
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return err
	}
	return dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
}

//-----------------------------------------------------------------------------
// general purpose registers

func acRdGPR(dbg *Debug, reg, size uint) (uint64, error) {
	switch size {
	case 32:
		val, err := dbg.acRd32(regGPR(reg))
		return uint64(val), err
	case 64:
		return dbg.acRd64(regGPR(reg))
	}
	return 0, fmt.Errorf("%d-bit gpr read not supported", size)
}

func acWrGPR(dbg *Debug, reg, size uint, val uint64) error {
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

func acRdFPR(dbg *Debug, reg, size uint) (uint64, error) {
	switch size {
	case 32:
		val, err := dbg.acRd32(regFPR(reg))
		return uint64(val), err
	case 64:
		return dbg.acRd64(regFPR(reg))
	}
	return 0, fmt.Errorf("%d-bit fpr read not supported", size)
}

func acWrFPR(dbg *Debug, reg, size uint, val uint64) error {
	switch size {
	case 32:
		return dbg.acWr32(regFPR(reg), uint32(val))
	case 64:
		return dbg.acWr64(regFPR(reg), val)
	}
	return fmt.Errorf("%d-bit fpr write not supported", size)
}

//-----------------------------------------------------------------------------
// control and status registers

func acRdCSR(dbg *Debug, reg, size uint) (uint64, error) {
	switch size {
	case 32:
		val, err := dbg.acRd32(regCSR(reg))
		return uint64(val), err
	case 64:
		return dbg.acRd64(regCSR(reg))
	}
	return 0, fmt.Errorf("%d-bit csr read not supported", size)
}

func acWrCSR(dbg *Debug, reg, size uint, val uint64) error {
	switch size {
	case 32:
		return dbg.acWr32(regCSR(reg), uint32(val))
	case 64:
		return dbg.acWr64(regCSR(reg), val)
	}
	return fmt.Errorf("%d-bit csr write not supported", size)
}

//-----------------------------------------------------------------------------
