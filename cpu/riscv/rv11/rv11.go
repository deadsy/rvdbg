//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 Functions

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------

// Debug is a RISC-V 0.11 debugger.
// It implements the rv.Debug interface.
type Debug struct {
}

// New returns a RISC-V 0.11 debugger.
func New(dev *jtag.Device) (rv.Debug, error) {
	log.Info.Printf("0.11 debug module")
	return nil, nil
}

//-----------------------------------------------------------------------------

// Test is a test routine.
func (dbg *Debug) Test() string {
	return "here"
}

//-----------------------------------------------------------------------------
