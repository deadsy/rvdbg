//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Functions

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

const irDtmcs = 0x10 // debug transport module control and status
const irDmi = 0x11   // debug module interface access

const drDtmcsLength = 32

const mask32 = (1 << 32) - 1

//-----------------------------------------------------------------------------

type Debug struct {
	dev         *jtag.Device
	ir          uint // cache of ir value
	irlen       int  // IR length
	drDmiLength int  // DR length for dmi
	abits       uint // address bits in dtmcs
	idle        uint // idle value in dtmcs
}

func NewDebug(dev *jtag.Device) (*Debug, error) {

	dbg := &Debug{
		dev:   dev,
		irlen: dev.GetIRLength(),
	}

	// get parameters from dtmcs
	dtmcs, err := dbg.rdDtmcs()
	if err != nil {
		return nil, err
	}

	dbg.abits = (dtmcs >> 4) & 0x3f
	dbg.idle = (dtmcs >> 12) & 7
	dbg.drDmiLength = 33 + int(dbg.abits) + 1

	// check dmi for the correct length
	_, err = dev.CheckDR(irDmi, dbg.drDmiLength)
	if err != nil {
		return nil, err
	}

	err = dbg.wrDtmcs(dmihardreset | dmireset)
	if err != nil {
		return nil, err
	}

	return dbg, nil
}

//-----------------------------------------------------------------------------

// wrIR writes the instruction register.
func (dbg *Debug) wrIR(ir uint) error {
	if ir == dbg.ir {
		return nil
	}
	err := dbg.dev.WrIR(bitstr.FromUint(ir, dbg.irlen))
	if err != nil {
		return err
	}
	dbg.ir = ir
	return nil
}

//-----------------------------------------------------------------------------
// dtmcs

const dmireset = (1 << 16)
const dmihardreset = (1 << 17)

// rdDtmcs reads the DTMCS register.
func (dbg *Debug) rdDtmcs() (uint, error) {
	err := dbg.wrIR(irDtmcs)
	if err != nil {
		return 0, err
	}
	tdo, err := dbg.dev.RdWrDR(bitstr.Zeros(drDtmcsLength), 0)
	if err != nil {
		return 0, err
	}
	return tdo.Split([]int{drDtmcsLength})[0], nil
}

// wrDtmcs writes the DTMCS register.
func (dbg *Debug) wrDtmcs(val uint) error {
	err := dbg.wrIR(irDtmcs)
	if err != nil {
		return err
	}
	return dbg.dev.WrDR(bitstr.FromUint(val, drDtmcsLength), 0)
}

//-----------------------------------------------------------------------------

func (dbg *Debug) Test() string {
	s := []string{}

	for i := 0x04; i <= 0x40; i++ {
		x, err := dbg.rdDmi(uint(i))
		if err != nil {
			s = append(s, fmt.Sprintf("%02x: %s", i, err))
		} else {
			s = append(s, fmt.Sprintf("%02x: %08x", i, x))
		}
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
