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
// general purpose registers

// rdGPR reads a general purpose register.
func (hi *hartInfo) rdGPR(reg uint) (uint64, error) {
	return 0, nil
}

// wrGPR writes a general purpose register.
func (hi *hartInfo) wrGPR(reg uint, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
// floating point registers

// rdFPR reads a floating point register.
func (hi *hartInfo) rdFPR(reg uint) (uint64, error) {
	return 0, nil
}

// wrFPR writes a floating point register.
func (hi *hartInfo) wrFPR(reg uint, val uint64) error {
	return nil
}

//-----------------------------------------------------------------------------
// control and status registers

// rdCSR reads a control and status register.
func (hi *hartInfo) rdCSR(csr uint) (uint, error) {
	return 0, nil
}

// wrCSR writes a control and status register.
func (hi *hartInfo) wrCSR(csr, val uint) error {
	return nil
}

//-----------------------------------------------------------------------------
