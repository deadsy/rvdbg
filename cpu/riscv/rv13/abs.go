//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Abstract Command Operations

*/
//-----------------------------------------------------------------------------

package rv13

//-----------------------------------------------------------------------------
// abstract command read operations

// acRd32 reads a 32-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd32(reg uint) (uint32, error) {
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

// acRd64 reads a 64-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd64(reg uint) (uint64, error) {
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

// acRd128 reads a 128-bit GPR/FPR/CSR using an abstract register read command.
func (dbg *Debug) acRd128(reg uint) (uint64, uint64, error) {
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
// abstract command write operations

// acWr32 writes a 32-bit GPR/FPR/CSR using an abstract register write command.
func (dbg *Debug) acWr32(reg uint, val uint32) error {
	ops := []dmiOp{
		// write val
		dmiWr(data0, val),
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
