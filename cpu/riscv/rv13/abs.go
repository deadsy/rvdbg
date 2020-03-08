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
		// clear any command error
		//dmiWr(abstractcs, errClear),
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

// acRd64 reads a 64-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd64(reg uint) (uint64, error) {
	ops := []dmiOp{
		// clear any command error
		//dmiWr(abstractcs, errClear),
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

// acRd128 reads a 128-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd128(reg uint) (uint64, uint64, error) {
	ops := []dmiOp{
		// clear any command error
		//dmiWr(abstractcs, errClear),
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
// abstract command write operations

// acWr32 writes a 32-bit GPR/FPR/CSR using an abstract register write command.
func (dbg *Debug) acWr32(reg uint, val uint32) error {
	ops := []dmiOp{
		// write val
		dmiWr(data0, val),
		// clear any command error
		//dmiWr(abstractcs, errClear),
		// write the register
		dmiWr(command, cmdRegister(reg, size32, false, false, true, true)),
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
		// clear any command error
		//dmiWr(abstractcs, errClear),
		// write the register
		dmiWr(command, cmdRegister(reg, size64, false, false, true, true)),
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

//-----------------------------------------------------------------------------
