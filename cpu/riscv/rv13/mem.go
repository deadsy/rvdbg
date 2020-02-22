//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Memory Operations

*/
//-----------------------------------------------------------------------------

package rv13

import "github.com/deadsy/rvdbg/cpu/riscv/rv"

//-----------------------------------------------------------------------------

// Rd32 reads a 32-bit value from a 32-bit address.
func (dbg *Debug) Rd32(addr uint32) (uint32, error) {

	ops := []dmiOp{
		// setup the program buffer
		dmiWr(progbuf0, rv.InsLW(rv.RegS0, rv.RegS0, 0)), // lw s0, 0(s0)
		dmiWr(progbuf1, rv.InsEBREAK()),                  // ebreak
		// s0 = addr, execute program buffer
		dmiWr(data0, addr),
		dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, false, true, true, true)),
		// read s0
		dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, false, false, true, false)),
		dmiRd(data0),
		// done
		dmiEnd(),
	}

	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

// Wr32 writes a 32-bit value to a 32-bit address.
func (dbg *Debug) Wr32(addr, val uint32) error {

	ops := []dmiOp{
		// setup the program buffer
		dmiWr(progbuf0, rv.InsSW(rv.RegS1, rv.RegS0, 0)), // sw s1, 0(s0)
		dmiWr(progbuf1, rv.InsEBREAK()),                  // ebreak
		// s0 = addr
		dmiWr(data0, addr),
		dmiWr(command, cmdRegister(regGPR(rv.RegS0), size32, false, false, true, true)),
		// s1 = val, execute program buffer
		dmiWr(data0, val),
		dmiWr(command, cmdRegister(regGPR(rv.RegS1), size32, false, true, true, true)),
		// done
		dmiEnd(),
	}

	_, err := dbg.dmiOps(ops)
	return err
}

//-----------------------------------------------------------------------------
