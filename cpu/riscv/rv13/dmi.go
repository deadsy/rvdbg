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

const data0 = 0x04 + 0      // Abstract Data 0
const data1 = 0x04 + 1      // Abstract Data 1
const data2 = 0x04 + 2      // Abstract Data 2
const data3 = 0x04 + 3      // Abstract Data 3
const data4 = 0x04 + 4      // Abstract Data 4
const data5 = 0x04 + 5      // Abstract Data 5
const data6 = 0x04 + 6      // Abstract Data 6
const data7 = 0x04 + 7      // Abstract Data 7
const data8 = 0x04 + 8      // Abstract Data 8
const data9 = 0x04 + 9      // Abstract Data 9
const data10 = 0x04 + 10    // Abstract Data 10
const data11 = 0x04 + 11    // Abstract Data 11
const dmcontrol = 0x10      // Debug Module Control
const dmstatus = 0x11       // Debug Module Status
const hartinfo = 0x12       // Hart Info
const haltsum1 = 0x13       // Halt Summary 1
const hawindowsel = 0x14    // Hart Array Window Select
const hawindow = 0x15       // Hart Array Window
const abstractcs = 0x16     // Abstract Control and Status
const command = 0x17        // Abstract Command
const abstractauto = 0x18   // Abstract Command Autoexec
const confstrptr0 = 0x19    // Configuration String Pointer 0
const confstrptr1 = 0x1a    // Configuration String Pointer 1
const confstrptr2 = 0x1b    // Configuration String Pointer 2
const confstrptr3 = 0x1c    // Configuration String Pointer 3
const nextdm = 0x1d         // Next Debug Module
const progbuf0 = 0x20 + 0   // Program Buffer 0
const progbuf1 = 0x20 + 1   // Program Buffer 1
const progbuf2 = 0x20 + 2   // Program Buffer 2
const progbuf3 = 0x20 + 3   // Program Buffer 3
const progbuf4 = 0x20 + 4   // Program Buffer 4
const progbuf5 = 0x20 + 5   // Program Buffer 5
const progbuf6 = 0x20 + 6   // Program Buffer 6
const progbuf7 = 0x20 + 7   // Program Buffer 7
const progbuf8 = 0x20 + 8   // Program Buffer 8
const progbuf9 = 0x20 + 9   // Program Buffer 9
const progbuf10 = 0x20 + 10 // Program Buffer 10
const progbuf11 = 0x20 + 11 // Program Buffer 11
const progbuf12 = 0x20 + 12 // Program Buffer 12
const progbuf13 = 0x20 + 13 // Program Buffer 13
const progbuf14 = 0x20 + 14 // Program Buffer 14
const progbuf15 = 0x20 + 15 // Program Buffer 15
const authdata = 0x30       // Authentication Data
const haltsum2 = 0x34       // Halt Summary 2
const haltsum3 = 0x35       // Halt Summary 3
const sbaddress3 = 0x37     // System Bus Address 127:96
const sbcs = 0x38           // System Bus Access Control and Status
const sbaddress0 = 0x39     // System Bus Address 31:0
const sbaddress1 = 0x3a     // System Bus Address 63:32
const sbaddress2 = 0x3b     // System Bus Address 95:64
const sbdata0 = 0x3c        // System Bus Data 31:0
const sbdata1 = 0x3d        // System Bus Data 63:32
const sbdata2 = 0x3e        // System Bus Data 95:64
const sbdata3 = 0x3f        // System Bus Data 127:96
const haltsum0 = 0x40       // Halt Summary 0

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

const mask32 = (1 << 32) - 1

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
// abstract commands

// regCSR returns the abstract register number for a control and status register.
func regCSR(i uint) uint {
	return 0 + (i & 0xfff)
}

// regGPR returns the abstract register number for a general purpose register.
func regGPR(i uint) uint {
	return 0x1000 + (i & 0x1f)
}

// regFPR returns the abstract register number for a floating point register.
func regFPR(i uint) uint {
	return 0x1020 + (i & 0x1f)
}

const (
	size8   = 0 // lower 8 bits
	size16  = 1 // lower 16 bits
	size32  = 2 // lower 32 bits
	size64  = 3 // lower 64 bits
	size128 = 4 // lower 128 bits
)

func cmdRegister(reg, size uint, postinc, postexec, transfer, write bool) uint32 {
	return uint32((0 << 24) |
		(size << 20) |
		(util.BoolToUint(postinc) << 19) |
		(util.BoolToUint(postexec) << 18) |
		(util.BoolToUint(transfer) << 17) |
		(util.BoolToUint(write) << 16) |
		(reg << 0))
}

func cmdQuick() uint32 {
	return uint32((1 << 24))
}

func cmdMemory(size uint, virtual, postinc, write bool) uint32 {
	return uint32((2 << 24) |
		(util.BoolToUint(virtual) << 23) |
		(size << 20) |
		(util.BoolToUint(postinc) << 19) |
		(util.BoolToUint(write) << 16))
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
