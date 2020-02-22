//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 Functions

*/
//-----------------------------------------------------------------------------

package rv11

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------

// Debug is a RISC-V 0.11 debugger.
type Debug struct {
}

// New returns a RISC-V 0.11 debugger.
func New(dev *jtag.Device) (*Debug, error) {
	return nil, nil
}

//-----------------------------------------------------------------------------

// Test is a test routine.
func (dbg *Debug) Test() string {
	return "here"
}

//-----------------------------------------------------------------------------
