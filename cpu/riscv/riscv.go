//-----------------------------------------------------------------------------
/*

RISC-V Debugger Functions

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/cpu/riscv/rv11"
	"github.com/deadsy/rvdbg/cpu/riscv/rv13"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

const irLength = 5

// ID Code
const irIDCode = 0x01
const drIDCodeLength = 32

// DTM register (for version)
const irDtm = 0x10
const drDtmLength = 32

//-----------------------------------------------------------------------------

// NewDebug returns a new RISC-V debugger interface.
func NewDebug(dev *jtag.Device) (rv.Debug, error) {

	// check the IR length
	if dev.GetIRLength() != irLength {
		return nil, fmt.Errorf("device ir length is %d, expected %d", dev.GetIRLength(), irLength)
	}

	// check the ID code
	idcode, err := dev.CheckDR(irIDCode, drIDCodeLength)
	if err != nil {
		return nil, err
	}
	if dev.GetIDCode() != jtag.IDCode(idcode) {
		return nil, fmt.Errorf("device idcode is 0x%08x, expected 0x%08x", uint32(idcode), uint32(dev.GetIDCode()))
	}

	// check the DTM register
	version, err := dev.CheckDR(irDtm, drDtmLength)
	if err != nil {
		return nil, err
	}

	// create a debug object based on version
	version &= 15
	switch version {
	case 0:
		return rv11.New(dev)
	case 1:
		return rv13.New(dev)
	}

	return nil, fmt.Errorf("unknown dtm version %d", version)
}

//-----------------------------------------------------------------------------
