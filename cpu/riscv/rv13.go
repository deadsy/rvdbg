//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Functions

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

const irDtmcs = 0x10 // debug transport module control and status
const irDmi = 0x11   // debug module interface access

const drDtmcsLength = 32

//-----------------------------------------------------------------------------

type rv13 struct {
	dev         *jtag.Device
	ir          uint // cache of ir value
	abits       uint // address bits in dtmcs
	idle        uint // idle value in dtmcs
	drDmiLength int
}

func newRv13(dev *jtag.Device) (*rv13, error) {

	dbg := &rv13{
		dev: dev,
	}

	dtmcs, err := dbg.rdDtmcs()
	if err != nil {
		return nil, err
	}

	dbg.abits = (dtmcs >> 4) & 0x3f
	dbg.idle = (dtmcs >> 12) & 7
	dbg.drDmiLength = 33 + int(dbg.abits) + 1

	// check the DMI register
	_, err = checkDR(dev, irDmi, dbg.drDmiLength)
	if err != nil {
		return nil, err
	}

	return dbg, nil
}

//-----------------------------------------------------------------------------

// wrIR writes the instruction register.
func (dbg *rv13) wrIR(ir uint) error {
	if ir == dbg.ir {
		return nil
	}
	err := dbg.dev.WrIR(bitstr.FromUint(ir, irLength))
	if err != nil {
		return err
	}
	dbg.ir = ir
	return nil
}

// rdDtmcs reads the DTMCS register.
func (dbg *rv13) rdDtmcs() (uint, error) {
	err := dbg.wrIR(irDtmcs)
	if err != nil {
		return 0, err
	}
	tdo, err := dbg.dev.RdWrDR(bitstr.Zeroes(drDtmcsLength))
	if err != nil {
		return 0, err
	}
	return tdo.Split([]int{drDtmcsLength})[0], nil
}

//-----------------------------------------------------------------------------
