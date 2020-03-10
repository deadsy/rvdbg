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
// slice conversion

// convert32to16 converts a 32-bit slice to a 16-bit slice.
func convert32to16(x []uint32) []uint16 {
	y := make([]uint16, len(x))
	for i := range x {
		y[i] = uint16(x[i])
	}
	return y
}

// convert16to32 converts a 16-bit slice to a 32-bit slice.
func convert16to32(x []uint16) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

// convert32to8 converts a 32-bit slice to an 8-bit slice.
func convert32to8(x []uint32) []uint8 {
	y := make([]uint8, len(x))
	for i := range x {
		y[i] = uint8(x[i])
	}
	return y
}

// convert8to32 converts an 8-bit slice to a 32-bit slice.
func convert8to32(x []uint8) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

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

// pbOps converts program buffer words into dmi operations.
func pbOps(pb []uint32) []dmiOp {
	ops := make([]dmiOp, len(pb))
	for i, v := range pb {
		ops[i] = dmiWr(progbuf(i), v)
	}
	return ops
}

//-----------------------------------------------------------------------------
// program buffer read operations

// pbRead reads a size-bit value using an program buffer operation.
func (dbg *Debug) pbRead(size uint, pb []uint32) (uint64, error) {
	// build the operations buffer
	ops := pbOps(pb)
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
	ops := pbOps(pb)
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
// read memory 8/16/32-bits

// pbRdMem_RV32 performs 8/16/32-bit memory reads using RV32 instructions.
func (dbg *Debug) pbRdMem_RV32(addr, n uint, pb []uint32) ([]uint32, error) {
	// build the operations buffer
	ops := pbOps(pb)
	// setup the address in dataX
	switch dbg.GetCurrentHart().MXLEN {
	case 32:
		// setup the 32-bit address in data0
		ops = append(ops, dmiWr(data0, uint32(addr)))
		// transfer the address to s0 and postexec to read the first value
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, false, true, true, true)))
	case 64:
		// setup the 64-bit address in data0/1
		ops = append(ops, dmiWr(data0, uint32(addr)))
		ops = append(ops, dmiWr(data1, uint32(addr>>32)))
		// transfer the address to s0 and postexec to read the first value
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size64, false, true, true, true)))
	default:
		return nil, fmt.Errorf("memory reads from a %d-bit address are not supported")
	}
	// the value read from memory is in s1
	// transfer s1 to data0 and then postexec to get the next value in s1
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size32, false, true, true, false)))
	// turn on autoexec for data0
	ops = append(ops, dmiWr(abstractauto, 1<<0))
	// do n-1 data reads
	for i := 0; i < int(n)-1; i++ {
		ops = append(ops, dmiRd(data0))
	}
	// turn off autoexec
	ops = append(ops, dmiWr(abstractauto, 0))
	// read the final data0 value
	ops = append(ops, dmiRd(data0))
	// read the command status
	ops = append(ops, dmiRd(abstractcs))
	// done
	ops = append(ops, dmiEnd())
	// run the operations
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return nil, err
	}
	// check the command status
	cs := cmdStatus(data[len(data)-1])
	err = dbg.checkError(cs)
	if err != nil {
		return nil, err
	}
	// return the data
	return data[:len(data)-1], nil
}

// pbRdMem8 reads n x 8-bit values from memory using program buffer operations.
func pbRdMem8(dbg *Debug, addr, n uint) ([]uint8, error) {
	// 8-bit reads
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsLB(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 1)
	// read the memory
	data, err := dbg.pbRdMem_RV32(addr, n, pb)
	if err != nil {
		return nil, err
	}
	return convert32to8(data), nil
}

// pbRdMem16 reads n x 16-bit values from memory using program buffer operations.
func pbRdMem16(dbg *Debug, addr, n uint) ([]uint16, error) {
	// 16-bit reads
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsLH(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 2)
	// read the memory
	data, err := dbg.pbRdMem_RV32(addr, n, pb)
	if err != nil {
		return nil, err
	}
	return convert32to16(data), nil
}

// pbRdMem32 reads n x 32-bit values from memory using program buffer operations.
func pbRdMem32(dbg *Debug, addr, n uint) ([]uint32, error) {
	// 32-bit reads
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsLW(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 4)
	// read the memory
	return dbg.pbRdMem_RV32(addr, n, pb)
}

//-----------------------------------------------------------------------------
// read memory 64-bits

// pbRdMem_RV64 performs 64-bit memory reads using RV64 instructions.
func (dbg *Debug) pbRdMem_RV64(addr, n uint, pb []uint32) ([]uint64, error) {
	return nil, errors.New("TODO")
}

// pbRdMem64 reads n x 64-bit values from memory using program buffer operations.
func pbRdMem64(dbg *Debug, addr, n uint) ([]uint64, error) {
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsLD(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 8)
	// read the memory
	return dbg.pbRdMem_RV64(addr, n, pb)
}

//-----------------------------------------------------------------------------
// write memory 8/16/32-bits

// pbWrMem_RV32 performs 8/16/32-bit memory writes using RV32 instructions.
func (dbg *Debug) pbWrMem_RV32(addr uint, val, pb []uint32) error {
	return errors.New("TODO")
}

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

//-----------------------------------------------------------------------------
// write memory 64-bits

// pbWrMem_RV64 performs 64-bit memory writes using RV64 instructions.
func (dbg *Debug) pbWrMem_RV64(addr uint, val []uint64, pb []uint32) error {
	return errors.New("TODO")
}

// pbWrMem64 writes n x 64-bit values to memory using program buffer operations.
func pbWrMem64(dbg *Debug, addr uint, val []uint64) error {
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
