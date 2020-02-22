//-----------------------------------------------------------------------------
/*

RISC-V Debugger Functions

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv11"
	"github.com/deadsy/rvdbg/cpu/riscv/rv13"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

const irLength = 5

const irIDCode = 0x01 // ID code
const irDtm = 0x10    // dtm register (for version)

const drIDCodeLength = 32
const drDtmLength = 32

//-----------------------------------------------------------------------------

// Debug is the RISC-V debug interface.
type Debug interface {
	Test() string
}

// CPU is a RISC-V cpu.
type CPU struct {
	dbg Debug
}

//-----------------------------------------------------------------------------

// NewCPU returns a new RISC-V cpu.
func NewCPU(dev *jtag.Device) (*CPU, error) {

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
	version &= 15

	cpu := &CPU{}

	switch version {
	case 0:
		cpu.dbg, err = rv11.New(dev)
		if err != nil {
			return nil, err
		}
	case 1:
		cpu.dbg, err = rv13.New(dev)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown dtm version %d", version)
	}

	return cpu, nil
}

//-----------------------------------------------------------------------------
