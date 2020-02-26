//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Register Operations

*/
//-----------------------------------------------------------------------------

package rv13

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

//-----------------------------------------------------------------------------
// general purpose registers

// RdGPR reads a general purpose register.
func (dbg *Debug) RdGPR(reg uint) (uint64, error) {
	// TODO 64-bit regs
	ops := []dmiOp{
		// read the register
		dmiWr(command, cmdRegister(regGPR(reg), size32, false, false, true, false)),
		dmiRd(data0),
		// done
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	return uint64(data[0]), nil
}

// WrGPR writes a general purpose register.
func (dbg *Debug) WrGPR(reg uint, val uint64) error {
	// TODO 64-bit regs
	ops := []dmiOp{
		// write the register
		dmiWr(data0, uint32(val)),
		dmiWr(command, cmdRegister(regGPR(reg), size32, false, false, true, true)),
		// done
		dmiEnd(),
	}
	_, err := dbg.dmiOps(ops)
	return err
}

//-----------------------------------------------------------------------------
// floating point registers

// RdFPR reads a floating point register.
func (dbg *Debug) RdFPR(reg uint) (uint64, error) {
	return 0, nil
}

// WrFPR writes a floating point register.
func (dbg *Debug) WrFPR(reg uint, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
// control and status registers

// RdCSR reads a control and status register.
func (dbg *Debug) RdCSR(csr uint) (uint32, error) {
	return 0, nil
}

// WrCSR writes a control and status register.
func (dbg *Debug) WrCSR(csr, val uint) error {
	return nil
}

//-----------------------------------------------------------------------------
