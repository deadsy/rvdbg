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
	abits       uint // address bits in dtmcs
	amask       uint // mask for address bits
	idle        uint // idle value in dtmcs
	irlen       int  // IR length
	drDmiLength int  // DR length for dmi
	dmiZeros    *bitstr.BitString
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
	dbg.amask = (1 << dbg.abits) - 1
	dbg.dmiZeros = bitstr.Zeros(dbg.drDmiLength)

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
	tdo, err := dbg.dev.RdWrDR(bitstr.Zeros(drDtmcsLength))
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
	return dbg.dev.WrDR(bitstr.FromUint(val, drDtmcsLength))
}

//-----------------------------------------------------------------------------
// debug module registers

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

// dmi operations
const opIgnore = 0
const opRd = 1
const opWr = 2

// dmi errors
const opOk = 0
const opFail = 2
const opBusy = 3
const opMask = (1 << 2) - 1

type dmiOp uint

func dmiRd(addr uint) dmiOp {
	return dmiOp((addr << 34) | opRd)
}

func dmiWr(addr uint, data uint32) dmiOp {
	return dmiOp((addr << 34) | (uint(data) << 2) | opWr)
}

func dmiEnd() dmiOp {
	return dmiOp(opIgnore)
}

func (x dmiOp) isRead() bool {
	return (x & opMask) == opRd
}

func (dbg *Debug) dmiOps(ops []dmiOp) ([]uint32, error) {
	data := []uint32{}

	// select dmi
	err := dbg.wrIR(irDmi)
	if err != nil {
		return nil, err
	}

	read := false
	for _, dmi := range ops {
		// run the operation
		tdo, err := dbg.dev.RdWrDR(bitstr.FromUint(uint(dmi), dbg.drDmiLength))
		if err != nil {
			return nil, err
		}
		x := tdo.Split([]int{dbg.drDmiLength})[0]
		// check the result
		result := x & opMask
		if result != opOk {
			// clear error condition
			dbg.wrDtmcs(dmireset)
			// re-select dmi
			dbg.wrIR(irDmi)

			// TODO auto-adjust timing
			return nil, fmt.Errorf("dmi operation error %d", result)
		}
		// get the read data
		if read {
			data = append(data, uint32((x>>2)&mask32))
			read = false
		}
		// setup the next read
		read = dmi.isRead()
	}
	return data, nil
}

//-----------------------------------------------------------------------------

func (dbg *Debug) rdDebugModule(addr uint) (uint32, error) {
	ops := []dmiOp{
		dmiRd(addr),
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (dbg *Debug) wrDebugModule(addr uint, data uint32) error {
	ops := []dmiOp{
		dmiWr(addr, data),
		dmiEnd(),
	}
	_, err := dbg.dmiOps(ops)
	return err
}

//-----------------------------------------------------------------------------

func (dbg *Debug) Test() string {
	s := []string{}
	for i := 0x04; i <= 0x40; i++ {
		x, err := dbg.rdDebugModule(uint(i))
		if err != nil {
			s = append(s, fmt.Sprintf("%02x: %s", i, err))
		} else {
			s = append(s, fmt.Sprintf("%02x: %08x", i, x))
		}
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
