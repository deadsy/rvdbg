//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Program Buffer Command Operations

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"errors"
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// newProgramBuffer returns a n word program buffer filled with EBREAKs.
func (dbg *Debug) newProgramBuffer(n uint) []uint32 {
	if n > dbg.progbufsize {
		return nil
	}
	pb := make([]uint32, n)
	for i := range pb {
		pb[i] = rv.InsEBREAK()
	}
	return pb
}

//-----------------------------------------------------------------------------
// program buffer read operations

// pbRead reads a size-bit value using an program buffer operation.
func (dbg *Debug) pbRead(size uint, pb []uint32) (uint64, error) {
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
	ops[n+1] = dmiWr(command, cmdRegister(regGPR(rv.RegS0), sizeMap[size], false, false, true, false))
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
	switch size {
	case 32:
		val, err := dbg.rdData32()
		return uint64(val), err
	case 64:
		return dbg.rdData64()
	}
	return 0, fmt.Errorf("%-bit read size not supported", size)
}

//-----------------------------------------------------------------------------
// program buffer write operations

// pbWr32 writes a 32-bit value using an program buffer operation.
func (dbg *Debug) pbWrite(reg, size uint, val uint32, pb []uint32) error {
	return nil
}

//-----------------------------------------------------------------------------

func pbRdCSR(dbg *Debug, reg, size uint) (uint64, error) {
	pb := dbg.newProgramBuffer(2)
	pb[0] = rv.InsCSRR(rv.RegS0, reg)
	return dbg.pbRead(size, pb)
}

func pbWrCSR(dbg *Debug, reg, size uint, val uint64) error {
	//pb := dbg.newProgramBuffer(2)
	//pb[0] = rv.InsCSRR(rv.RegS0, reg)
	//return dbg.pbRead(size, pb)
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
