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
// debug module debug bus registers

const data0 = 0x04        // Abstract Data 0
const data11 = 0x0f       // Abstract Data 11
const dmcontrol = 0x10    // Debug Module Control
const dmstatus = 0x11     // Debug Module Status
const hartinfo = 0x12     // Hart Info
const haltsum1 = 0x13     // Halt Summary 1
const hawindowsel = 0x14  // Hart Array Window Select
const hawindow = 0x15     // Hart Array Window
const abstractcs = 0x16   // Abstract Control and Status
const command = 0x17      // Abstract Command
const abstractauto = 0x18 // Abstract Command Autoexec
const confstrptr0 = 0x19  // Configuration String Pointer 0
const confstrptr1 = 0x1a  // Configuration String Pointer 1
const confstrptr2 = 0x1b  // Configuration String Pointer 2
const confstrptr3 = 0x1c  // Configuration String Pointer 3
const nextdm = 0x1d       // Next Debug Module
const progbuf0 = 0x20     // Program Buffer 0
const progbuf15 = 0x2f    // Program Buffer 15
const authdata = 0x30     // Authentication Data
const haltsum2 = 0x34     // Halt Summary 2
const haltsum3 = 0x35     // Halt Summary 3
const sbaddress3 = 0x37   // System Bus Address 127:96
const sbcs = 0x38         // System Bus Access Control and Status
const sbaddress0 = 0x39   // System Bus Address 31:0
const sbaddress1 = 0x3a   // System Bus Address 63:32
const sbaddress2 = 0x3b   // System Bus Address 95:64
const sbdata0 = 0x3c      // System Bus Data 31:0
const sbdata1 = 0x3d      // System Bus Data 63:32
const sbdata2 = 0x3e      // System Bus Data 95:64
const sbdata3 = 0x3f      // System Bus Data 127:96
const haltsum0 = 0x40     // Halt Summary 0

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
