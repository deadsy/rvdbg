//-----------------------------------------------------------------------------
/*

RISC-V Debugger Functions

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"fmt"

	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

const irLength = 5

const irIDCode = 0x01 // ID code

const drIDCodeLength = 32

//-----------------------------------------------------------------------------

// checkDR verifies the DR length for a given IR and returns the DR value.
func checkDR(dev *jtag.Device, ir, drlen int) (uint, error) {
	// write IR
	err := dev.WrIR(bitstr.FromUint(uint(ir), irLength))
	if err != nil {
		return 0, nil
	}
	// check the DR length
	n, err := dev.GetDRLength()
	if err != nil {
		return 0, nil
	}
	if n != drlen {
		return 0, fmt.Errorf("ir %d dr length is %d, expected %d", ir, n, drlen)
	}
	// get the value
	tdo, err := dev.RdWrDR(bitstr.Zeros(drlen))
	if err != nil {
		return 0, err
	}
	return uint(tdo.Split([]int{drlen})[0]), nil
}

//-----------------------------------------------------------------------------

type Debug interface {
}

//-----------------------------------------------------------------------------

func NewDebug(dev *jtag.Device) (Debug, error) {

	// check the IR length
	if dev.GetIRLength() != irLength {
		return nil, fmt.Errorf("device ir length is %d, expected %d", dev.GetIRLength(), irLength)
	}

	// check the ID code
	idcode, err := checkDR(dev, irIDCode, drIDCodeLength)
	if err != nil {
		return nil, err
	}
	if dev.GetIDCode() != jtag.IDCode(idcode) {
		return nil, fmt.Errorf("device idcode is 0x%08x, expected 0x%08x", uint32(idcode), uint32(dev.GetIDCode()))
	}

	// check the DTM register
	version, err := checkDR(dev, irDtmcs, drDtmcsLength)
	if err != nil {
		return nil, err
	}
	version &= 15

	// return the version specific debugger
	if version == 0 {
		return newRv11(dev)
	}
	if version == 1 {
		return newRv13(dev)
	}
	return nil, fmt.Errorf("unknown dtm version %d", version)
}

//-----------------------------------------------------------------------------
