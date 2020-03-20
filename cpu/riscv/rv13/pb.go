//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Program Buffer Command Operations

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// newProgramBuffer returns a n word program buffer filled with EBREAKs.
func (dbg *Debug) newProgramBuffer(n uint) []uint32 {
	pb := make([]uint32, n)
	for i := range pb {
		pb[i] = rv.InsEBREAK()
	}
	return pb
}

// pbOps converts program buffer words into dmi operations.
func pbOps(pb []uint32, n int) []dmiOp {
	ops := make([]dmiOp, len(pb), len(pb)+n)
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
	ops := pbOps(pb, 4)
	// postexec
	ops = append(ops, dmiWr(command, cmdRegister(0, 0, cmdPostExec)))
	// transfer GPR s0 to data0
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), sizeMap[size], cmdRead)))
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
	ops := pbOps(pb, 5)
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
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), sizeMap[size], cmdWrite|cmdPostExec)))
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

// pbRdMemRV32 performs n X 8/16/32-bit memory reads using RV32 instructions.
// This sequence checks for command errors at the end of the sequence.
// Slower devices may return busy errors.
func (dbg *Debug) pbRdMemRV32(addr, n uint, pb []uint32) ([]uint32, error) {
	// build the operations buffer
	ops := pbOps(pb, int(n)+10)
	// setup the address in dataX
	mxlen := dbg.GetCurrentHart().MXLEN
	switch mxlen {
	case 32:
		// setup the 32-bit address in data0
		ops = append(ops, dmiWr(data0, uint32(addr)))
		// transfer the address to s0 and postexec to read the first value
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, cmdWrite|cmdPostExec)))
	case 64:
		// setup the 64-bit address in data0/1
		ops = append(ops, dmiWr(data0, uint32(addr)))
		ops = append(ops, dmiWr(data1, uint32(addr>>32)))
		// transfer the address to s0 and postexec to read the first value
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size64, cmdWrite|cmdPostExec)))
	default:
		return nil, fmt.Errorf("memory reads from a %d-bit address are not supported", mxlen)
	}
	// the value read from memory is in s1
	if n == 1 {
		// transfer s1 to data0
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size32, cmdRead)))
	} else {
		// transfer s1 to data0 and then postexec to get the next value in s1
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size32, cmdRead|cmdPostExec)))
		// turn on autoexec for data0
		ops = append(ops, dmiWr(abstractauto, 1<<0))
		// do n-1 data reads
		for i := 0; i < int(n)-1; i++ {
			ops = append(ops, dmiRd(data0))
		}
		// turn off autoexec
		ops = append(ops, dmiWr(abstractauto, 0))
	}
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

// pbRdMemSingleRV32 performs a single 8/16/32-bit memory reads using RV32 instructions.
// This sequence checks for a busy condition on the initial read and waits for completion.
func (dbg *Debug) pbRdMemSingleRV32(addr uint, pb []uint32) ([]uint32, error) {
	// build the operations buffer
	ops := pbOps(pb, 5)
	// setup the address in dataX
	mxlen := dbg.GetCurrentHart().MXLEN
	switch mxlen {
	case 32:
		// setup the 32-bit address in data0
		ops = append(ops, dmiWr(data0, uint32(addr)))
		// transfer the address to s0 and postexec to read the first value
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, cmdWrite|cmdPostExec)))
	case 64:
		// setup the 64-bit address in data0/1
		ops = append(ops, dmiWr(data0, uint32(addr)))
		ops = append(ops, dmiWr(data1, uint32(addr>>32)))
		// transfer the address to s0 and postexec to read the value.
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size64, cmdWrite|cmdPostExec)))
	default:
		return nil, fmt.Errorf("memory reads from a %d-bit address are not supported", mxlen)
	}
	// read the command status
	ops = append(ops, dmiRd(abstractcs))
	// done
	ops = append(ops, dmiEnd())
	// run the operations
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return nil, err
	}
	// wait for command completion
	cs := cmdStatus(data[0])
	err = dbg.cmdWait(cs, cmdTimeout)
	if err != nil {
		return nil, err
	}
	// the value read from memory is in s1
	ops = []dmiOp{
		// transfer s1 to data0
		dmiWr(command, cmdRegister(regGPR(rv.RegS1), size32, cmdRead)),
		// read the data0 value
		dmiRd(data0),
		// read the command status
		dmiRd(abstractcs),
		// done
		dmiEnd(),
	}
	// run the operations
	data, err = dbg.dmiOps(ops)
	if err != nil {
		return nil, err
	}
	// check the command status
	cs = cmdStatus(data[1])
	err = dbg.checkError(cs)
	if err != nil {
		return nil, err
	}
	// return the data
	return data[0:1], nil
}

//-----------------------------------------------------------------------------

// pbRdMem8 reads n x 8-bit values from memory using program buffer operations.
func (dbg *Debug) pbRdMem8(addr, n uint) ([]uint, error) {
	var err error
	var data []uint32
	if n == 1 {
		pb := dbg.newProgramBuffer(2)
		pb[0] = rv.InsLB(rv.RegS1, 0, rv.RegS0)
		data, err = dbg.pbRdMemSingleRV32(addr, pb)
	} else {
		pb := dbg.newProgramBuffer(3)
		pb[0] = rv.InsLB(rv.RegS1, 0, rv.RegS0)
		pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 1)
		data, err = dbg.pbRdMemRV32(addr, n, pb)
	}
	if err != nil {
		return nil, err
	}
	return util.Cast32toUint(data, util.Mask8), nil
}

// pbRdMem16 reads n x 16-bit values from memory using program buffer operations.
func (dbg *Debug) pbRdMem16(addr, n uint) ([]uint, error) {
	var err error
	var data []uint32
	if n == 1 {
		pb := dbg.newProgramBuffer(2)
		pb[0] = rv.InsLH(rv.RegS1, 0, rv.RegS0)
		data, err = dbg.pbRdMemSingleRV32(addr, pb)
	} else {
		pb := dbg.newProgramBuffer(3)
		pb[0] = rv.InsLH(rv.RegS1, 0, rv.RegS0)
		pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 2)
		data, err = dbg.pbRdMemRV32(addr, n, pb)
	}
	if err != nil {
		return nil, err
	}
	return util.Cast32toUint(data, util.Mask16), nil
}

// pbRdMem32 reads n x 32-bit values from memory using program buffer operations.
func (dbg *Debug) pbRdMem32(addr, n uint) ([]uint, error) {
	var err error
	var data []uint32
	if n == 1 {
		pb := dbg.newProgramBuffer(2)
		pb[0] = rv.InsLW(rv.RegS1, 0, rv.RegS0)
		data, err = dbg.pbRdMemSingleRV32(addr, pb)
	} else {
		pb := dbg.newProgramBuffer(3)
		pb[0] = rv.InsLW(rv.RegS1, 0, rv.RegS0)
		pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 4)
		data, err = dbg.pbRdMemRV32(addr, n, pb)
	}
	if err != nil {
		return nil, err
	}
	return util.Cast32toUint(data, util.Mask32), nil
}

//-----------------------------------------------------------------------------
// read memory 64-bits

// pbRdMem_RV64 performs n X 64-bit memory reads using RV64 instructions.
func (dbg *Debug) pbRdMemRV64(addr, n uint, pb []uint32) ([]uint64, error) {
	// build the operations buffer
	ops := pbOps(pb, (int(n)<<1)+10)
	// setup the address in dataX
	mxlen := dbg.GetCurrentHart().MXLEN
	switch mxlen {
	case 64:
		// setup the 64-bit address in data0/1
		ops = append(ops, dmiWr(data0, uint32(addr)))
		ops = append(ops, dmiWr(data1, uint32(addr>>32)))
		// transfer the address to s0 and postexec to read the first value
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size64, cmdWrite|cmdPostExec)))
	default:
		return nil, fmt.Errorf("memory reads from a %d-bit address are not supported", mxlen)
	}
	// the value read from memory is in s1
	if n == 1 {
		// transfer s1 to data0/1
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size64, cmdRead)))
	} else {
		// transfer s1 to data0/1 and then postexec to get the next value in s1
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size64, cmdRead|cmdPostExec)))
		// turn on autoexec for data1
		ops = append(ops, dmiWr(abstractauto, 1<<1))
		// do n-1 data reads
		for i := 0; i < int(n)-1; i++ {
			ops = append(ops, dmiRd(data0))
			ops = append(ops, dmiRd(data1))
		}
		// turn off autoexec
		ops = append(ops, dmiWr(abstractauto, 0))
	}
	// read the final data0/1 value
	ops = append(ops, dmiRd(data0))
	ops = append(ops, dmiRd(data1))
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
	return util.Convert32to64(data[:len(data)-1]), nil
}

// pbRdMem64 reads n x 64-bit values from memory using program buffer operations.
func (dbg *Debug) pbRdMem64(addr, n uint) ([]uint, error) {
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsLD(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 8)
	// read the memory
	data, err := dbg.pbRdMemRV64(addr, n, pb)
	if err != nil {
		return nil, err
	}
	return util.Cast64toUint(data, util.Mask64), nil
}

//-----------------------------------------------------------------------------

func pbRdMem(dbg *Debug, width, addr, n uint) ([]uint, error) {
	switch width {
	case 8:
		return dbg.pbRdMem8(addr, n)
	case 16:
		return dbg.pbRdMem16(addr, n)
	case 32:
		return dbg.pbRdMem32(addr, n)
	case 64:
		return dbg.pbRdMem64(addr, n)
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
}

//-----------------------------------------------------------------------------
// write memory 8/16/32-bits

// pbWrMemRV32 performs 8/16/32-bit memory writes using RV32 instructions.
func (dbg *Debug) pbWrMemRV32(addr uint, val, pb []uint32) error {
	// build the operations buffer
	ops := pbOps(pb, len(val)+10)
	// setup the address in dataX
	mxlen := dbg.GetCurrentHart().MXLEN
	switch mxlen {
	case 32:
		// setup the 32-bit address in data0
		ops = append(ops, dmiWr(data0, uint32(addr)))
		// transfer data0 to s0
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, cmdWrite)))
	case 64:
		// setup the 64-bit address in data0/1
		ops = append(ops, dmiWr(data0, uint32(addr)))
		ops = append(ops, dmiWr(data1, uint32(addr>>32)))
		// transfer data0/1 to s0
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size64, cmdWrite)))
	default:
		return fmt.Errorf("memory writes to a %d-bit address are not supported", mxlen)
	}
	// setup val[0] in data0
	ops = append(ops, dmiWr(data0, val[0]))
	// transfer data0 to s1 and then postexec to write the value to memory.
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size32, cmdWrite|cmdPostExec)))
	if len(val) > 1 {
		// turn on autoexec for data0
		ops = append(ops, dmiWr(abstractauto, 1<<0))
		// write the rest of the buffer
		for i := 1; i < len(val); i++ {
			ops = append(ops, dmiWr(data0, val[i]))
		}
		// turn off autoexec
		ops = append(ops, dmiWr(abstractauto, 0))
	}
	// read the command status
	ops = append(ops, dmiRd(abstractcs))
	// done
	ops = append(ops, dmiEnd())
	// run the operations
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return err
	}
	// check the command status
	return dbg.checkError(cmdStatus(data[0]))
}

// pbWrMem8 writes n x 8-bit values to memory using program buffer operations.
func (dbg *Debug) pbWrMem8(addr uint, val []uint) error {
	// 8-bit writes
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsSB(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 1)
	return dbg.pbWrMemRV32(addr, util.CastUintto32(val), pb)
}

// pbWrMem16 writes n x 16-bit values to memory using program buffer operations.
func (dbg *Debug) pbWrMem16(addr uint, val []uint) error {
	// 16-bit writes
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsSH(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 2)
	return dbg.pbWrMemRV32(addr, util.CastUintto32(val), pb)
}

// pbWrMem32 writes n x 32-bit values to memory using program buffer operations.
func (dbg *Debug) pbWrMem32(addr uint, val []uint) error {
	// 32-bit writes
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsSW(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 4)
	return dbg.pbWrMemRV32(addr, util.CastUintto32(val), pb)
}

//-----------------------------------------------------------------------------
// write memory 64-bits

// pbWrMemRV64 performs 64-bit memory writes using RV64 instructions.
func (dbg *Debug) pbWrMemRV64(addr uint, val []uint64, pb []uint32) error {
	// build the operations buffer
	ops := pbOps(pb, (len(val)<<1)+10)
	// setup the address in dataX
	mxlen := dbg.GetCurrentHart().MXLEN
	switch mxlen {
	case 64:
		// setup the 64-bit address in data0/1
		ops = append(ops, dmiWr(data0, uint32(addr)))
		ops = append(ops, dmiWr(data1, uint32(addr>>32)))
		// transfer data0/1 to s0
		ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS0), size64, cmdWrite)))
	default:
		return fmt.Errorf("memory writes to a %d-bit address are not supported", mxlen)
	}
	// setup val[0] in data0/1
	ops = append(ops, dmiWr(data0, uint32(val[0])))
	ops = append(ops, dmiWr(data1, uint32(val[0]>>32)))
	// transfer data0/1 to s1 and then postexec to write the value to memory.
	ops = append(ops, dmiWr(command, cmdRegister(regGPR(rv.RegS1), size64, cmdWrite|cmdPostExec)))
	if len(val) > 1 {
		// turn on autoexec for data1
		ops = append(ops, dmiWr(abstractauto, 1<<1))
		// write the rest of the buffer
		for i := 1; i < len(val); i++ {
			ops = append(ops, dmiWr(data0, uint32(val[i])))
			ops = append(ops, dmiWr(data1, uint32(val[i]>>32)))
		}
		// turn off autoexec
		ops = append(ops, dmiWr(abstractauto, 0))
	}
	// read the command status
	ops = append(ops, dmiRd(abstractcs))
	// done
	ops = append(ops, dmiEnd())
	// run the operations
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return err
	}
	// check the command status
	return dbg.checkError(cmdStatus(data[0]))
}

// pbWrMem64 writes n x 64-bit values to memory using program buffer operations.
func (dbg *Debug) pbWrMem64(addr uint, val []uint) error {
	// 64-bit writes
	pb := dbg.newProgramBuffer(3)
	pb[0] = rv.InsSD(rv.RegS1, 0, rv.RegS0)
	pb[1] = rv.InsADDI(rv.RegS0, rv.RegS0, 8)
	return dbg.pbWrMemRV64(addr, util.CastUintto64(val), pb)
}

//-----------------------------------------------------------------------------

func pbWrMem(dbg *Debug, width, addr uint, val []uint) error {
	switch width {
	case 8:
		return dbg.pbWrMem8(addr, val)
	case 16:
		return dbg.pbWrMem16(addr, val)
	case 32:
		return dbg.pbWrMem32(addr, val)
	case 64:
		return dbg.pbWrMem64(addr, val)

	}
	return fmt.Errorf("%d-bit memory writes are not supported", width)
}

//-----------------------------------------------------------------------------
