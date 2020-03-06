//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Program Buffer Command Operations

*/
//-----------------------------------------------------------------------------

package rv13

import "github.com/deadsy/rvdbg/cpu/riscv/rv"

//-----------------------------------------------------------------------------
// program buffer read operations

// pbRd32 reads a 32-bit value using an program buffer operation.
func (dbg *Debug) pbRd32(reg uint, pb []uint32) (uint32, error) {

	n := len(pb)

	// build the operations buffer
	ops := make([]dmiOp, n+4)
	// write the program buffer
	for i, v := range pb {
		ops[i] = dmiWr(progbuf(i), v)
	}
	// postexec
	ops[n+0] = dmiWr(command, cmdRegister(0, 0, false, true, false, false))
	// transfer GPR s0 to data0
	ops[n+1] = dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, false, false, true, false))
	// read the command status
	ops[n+2] = dmiRd(abstractcs)
	// done
	ops[n+3] = dmiEnd()

	// run the operations
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}

	// wait for command completion
	err = dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
	if err != nil {
		return 0, err
	}

	// read the data
	return dbg.rdData32()
}

// pbRd64 reads a 64-bit value using an program buffer operation.
func (dbg *Debug) pbRd64(reg uint, pb []uint32) (uint32, error) {
	return 0, nil
}

//-----------------------------------------------------------------------------
// program buffer write operations

// pbWr32 writes a 32-bit value using an program buffer operation.
func (dbg *Debug) pbWr32(reg uint, val uint32, pn []uint32) error {
	return nil
}

// pbWr64 writes a 64-bit value using an program buffer operation.
func (dbg *Debug) pbWr64(reg uint, val uint64, pn []uint32) error {
	return nil
}

//-----------------------------------------------------------------------------
