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
	// build the operations buffer
	ops := []dmiOp{}
	// write the program buffer
	for i, v := range pb {
		ops = append(ops, dmiWr(progbuf(i), v))
	}
	// postexec
	ops = append(ops, dmiWr(command, cmdRegister(0, 0, false, true, false, false)))
	// transfer GPR s0 to data0
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), sizeMap[size], false, false, true, false)))
	// read the command status
	ops = append(ops, dmiRd(abstractcs))
	// done
	ops = append(ops, dmiEnd())
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

// pbRdCSR reads a CSR using program buffer operations.
func pbRdCSR(dbg *Debug, reg, size uint) (uint64, error) {
	pb := dbg.newProgramBuffer(2)
	pb[0] = rv.InsCSRR(rv.RegS0, reg)
	return dbg.pbRead(size, pb)
}

//-----------------------------------------------------------------------------
// program buffer write operations

// pbWrite writes a size-bit value using an program buffer operation.
func (dbg *Debug) pbWrite(size uint, val uint64, pb []uint32) error {
	// build the operations buffer
	ops := []dmiOp{}
	// write the program buffer
	for i, v := range pb {
		ops = append(ops, dmiWr(progbuf(i), v))
	}
	// setup dataX with the value to write
	switch size {
	case 32:
		ops = append(ops, dmiWr(data0, uint32(val)))
	case 64:
		ops = append(ops, dmiWr(data0, uint32(val)))
		ops = append(ops, dmiWr(data1, uint32(val>>32)))
	default:
		return fmt.Errorf("%-bit write size not supported", size)
	}
	// transfer dataX to GPR s0 and then postexec
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), sizeMap[size], false, true, true, true)))
	// read the command status
	ops = append(ops, dmiRd(abstractcs))
	// done
	ops = append(ops, dmiEnd())
	// run the operations
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return err
	}
	// wait for command completion
	return dbg.cmdWait(cmdStatus(data[0]), cmdTimeout)
}

// pbWrCSR writes a CSR using program buffer operations.
func pbWrCSR(dbg *Debug, reg, size uint, val uint64) error {
	pb := dbg.newProgramBuffer(2)
	pb[0] = rv.InsCSRW(reg, rv.RegS0)
	return dbg.pbWrite(size, val, pb)
}

//-----------------------------------------------------------------------------
// read memory

// pbRdMem8 reads n x 8-bit values from memory using program buffer operations.
func pbRdMem8(dbg *Debug, addr, n uint) ([]uint8, error) {
	return nil, errors.New("TODO")
}

// pbRdMem16 reads n x 16-bit values from memory using program buffer operations.
func pbRdMem16(dbg *Debug, addr, n uint) ([]uint16, error) {
	return nil, errors.New("TODO")
}

// pbRdMem32 reads n x 32-bit values from memory using program buffer operations.
func pbRdMem32(dbg *Debug, addr, n uint) ([]uint32, error) {
	return nil, errors.New("TODO")
}

// pbRdMem64 reads n x 64-bit values from memory using program buffer operations.
func pbRdMem64(dbg *Debug, addr, n uint) ([]uint64, error) {
	return nil, errors.New("TODO")
}

//-----------------------------------------------------------------------------
// write memory

// pbWrMem8 writes n x 8-bit values to memory using program buffer operations.
func pbWrMem8(dbg *Debug, addr uint, val []uint8) error {
	return errors.New("TODO")
}

// pbWrMem16 writes n x 16-bit values to memory using program buffer operations.
func pbWrMem16(dbg *Debug, addr uint, val []uint16) error {
	return errors.New("TODO")
}

// pbWrMem32 writes n x 32-bit values to memory using program buffer operations.
func pbWrMem32(dbg *Debug, addr uint, val []uint32) error {
	return errors.New("TODO")
}

// pbWrMem64 writes n x 64-bit values to memory using program buffer operations.
func pbWrMem64(dbg *Debug, addr uint, val []uint64) error {
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
