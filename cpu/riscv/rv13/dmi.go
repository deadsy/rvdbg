//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13
Debug Module Interface

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"

	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/util"
)

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

var dmstatusFields = util.FieldSet{
	{"impebreak", 22, 22, util.FmtDec},
	{"allhavereset", 19, 19, util.FmtDec},
	{"anyhavereset", 18, 18, util.FmtDec},
	{"allresumeack", 17, 17, util.FmtDec},
	{"anyresumeack", 16, 16, util.FmtDec},
	{"allnonexistent", 15, 15, util.FmtDec},
	{"anynonexistent", 14, 14, util.FmtDec},
	{"allunavail", 13, 13, util.FmtDec},
	{"anyunavail", 12, 12, util.FmtDec},
	{"allrunning", 11, 11, util.FmtDec},
	{"anyrunning", 10, 10, util.FmtDec},
	{"allhalted", 9, 9, util.FmtDec},
	{"anyhalted", 8, 8, util.FmtDec},
	{"authenticated", 7, 7, util.FmtDec},
	{"authbusy", 6, 6, util.FmtDec},
	{"hasresethaltreq", 5, 5, util.FmtDec},
	{"confstrptrvalid", 4, 4, util.FmtDec},
	{"version", 3, 0, util.FmtDec},
}

var hartinfoFields = util.FieldSet{
	{"nscratch", 23, 20, util.FmtDec},
	{"dataaccess", 16, 16, util.FmtDec},
	{"datasize", 15, 12, util.FmtDec},
	{"dataaddr", 11, 0, util.FmtHex},
}

var abstractcsFields = util.FieldSet{
	{"progbufsize", 28, 24, util.FmtDec},
	{"busy", 12, 12, util.FmtDec},
	{"cmderr", 10, 8, util.FmtDec},
	{"datacount", 3, 0, util.FmtDec},
}

//-----------------------------------------------------------------------------

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
	for i := 0; i < len(ops); i++ {
		dmi := ops[i]
		// run the operation
		tdo, err := dbg.dev.RdWrDR(bitstr.FromUint(uint(dmi), dbg.drDmiLength), dbg.idle)
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
			if result == opBusy {
				// auto-adjust timing
				dbg.idle++
				if dbg.idle > jtag.MaxIdle {
					return nil, fmt.Errorf("dmi operation error %d", result)
				}
				// redo the operation
				i--
			} else {
				return nil, fmt.Errorf("dmi operation error %d", result)
			}
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

// wrDmi reads a debug module interface register.
func (dbg *Debug) rdDmi(addr uint) (uint32, error) {
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

// wrDmi writes a debug module interface register.
func (dbg *Debug) wrDmi(addr uint, data uint32) error {
	ops := []dmiOp{
		dmiWr(addr, data),
		dmiEnd(),
	}
	_, err := dbg.dmiOps(ops)
	return err
}

//-----------------------------------------------------------------------------
