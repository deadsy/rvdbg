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
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

const irDtmcs = 0x10 // debug transport module control and status
const irDmi = 0x11   // debug module interface access

const drDtmcsLength = 32

//-----------------------------------------------------------------------------

// Debug is a RISC-V 0.13 debugger.
type Debug struct {
	dev             *jtag.Device
	ir              uint // cache of ir value
	irlen           int  // IR length
	drDmiLength     int  // DR length for dmi
	abits           uint // address bits in dtmcs
	idle            uint // idle value in dtmcs
	progbufsize     uint // number of progbuf words implemented
	datacount       uint // number of data words implemented
	autoexecprogbuf bool // can we autoexec on progbufX access?
	autoexecdata    bool // can we autoexec on dataX access?
	sbasize         uint // width of system bus address (0 = no access)
}

// New returns a RISC-V 0.13 debugger.
func New(dev *jtag.Device) (*Debug, error) {

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

	// reset dmi
	err = dbg.wrDtmcs(dmihardreset | dmireset)
	if err != nil {
		return nil, err
	}

	// make the dmi active
	err = dbg.wrDmi(dmcontrol, 1)
	if err != nil {
		return nil, err
	}

	// work out how many program and data words we have
	x, err := dbg.rdDmi(abstractcs)
	if err != nil {
		return nil, err
	}
	dbg.progbufsize = util.Bits(uint(x), 28, 24)
	dbg.datacount = util.Bits(uint(x), 3, 0)

	// test program buffers
	err = dbg.testBuffers(progbuf0, dbg.progbufsize)
	if err != nil {
		return nil, err
	}

	// test data buffers
	err = dbg.testBuffers(data0, dbg.datacount)
	if err != nil {
		return nil, err
	}

	// work out if we can autoexec on progbuf/data access
	err = dbg.wrDmi(abstractauto, 0xffffffff)
	if err != nil {
		return nil, err
	}
	x, err = dbg.rdDmi(abstractauto)
	if err != nil {
		return nil, err
	}
	if util.Bits(uint(x), 31, 16) == ((1 << dbg.progbufsize) - 1) {
		dbg.autoexecprogbuf = true
	}
	if util.Bits(uint(x), 11, 0) == ((1 << dbg.datacount) - 1) {
		dbg.autoexecdata = true
	}

	// work out the system bus address size
	x, err = dbg.rdDmi(sbcs)
	if err != nil {
		return nil, err
	}
	dbg.sbasize = util.Bits(uint(x), 11, 5)

	return dbg, nil
}

func (dbg *Debug) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("version 0.13"))
	s = append(s, fmt.Sprintf("idle cycles %d", dbg.idle))
	s = append(s, fmt.Sprintf("sbasize %d", dbg.sbasize))
	s = append(s, fmt.Sprintf("progbufsize %d", dbg.progbufsize))
	s = append(s, fmt.Sprintf("datacount %d", dbg.datacount))
	s = append(s, fmt.Sprintf("autoexecprogbuf %t", dbg.autoexecprogbuf))
	s = append(s, fmt.Sprintf("autoexecdata %t", dbg.autoexecdata))
	return strings.Join(s, "\n")
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

// Test is a test routine.
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
