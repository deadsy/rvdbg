//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13
Debug Module Interface

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------
// debug module registers

const data0 = 0x04    // Abstract Data 0-11
const progbuf0 = 0x20 // Program Buffer 0-15

const dmcontrol = 0x10 // Debug Module Control
const dmstatus = 0x11  // Debug Module Status

const hartinfo = 0x12    // Hart Info
const hawindowsel = 0x14 // Hart Array Window Select
const hawindow = 0x15    // Hart Array Window

const abstractcs = 0x16   // Abstract Control and Status
const command = 0x17      // Abstract Command
const abstractauto = 0x18 // Abstract Command Autoexec

const confstrptr0 = 0x19 // Configuration String Pointer 0
const confstrptr1 = 0x1a // Configuration String Pointer 1
const confstrptr2 = 0x1b // Configuration String Pointer 2
const confstrptr3 = 0x1c // Configuration String Pointer 3

const nextdm = 0x1d   // Next Debug Module
const authdata = 0x30 // Authentication Data

const haltsum0 = 0x40 // Halt Summary 0
const haltsum1 = 0x13 // Halt Summary 1
const haltsum2 = 0x34 // Halt Summary 2
const haltsum3 = 0x35 // Halt Summary 3

const sbcs = 0x38 // System Bus Access Control and Status

const sbaddress0 = 0x39 // System Bus Address 31:0
const sbaddress1 = 0x3a // System Bus Address 63:32
const sbaddress2 = 0x3b // System Bus Address 95:64
const sbaddress3 = 0x37 // System Bus Address 127:96

const sbdata0 = 0x3c // System Bus Data 31:0
const sbdata1 = 0x3d // System Bus Data 63:32
const sbdata2 = 0x3e // System Bus Data 95:64
const sbdata3 = 0x3f // System Bus Data 127:96

//-----------------------------------------------------------------------------
// DM control

var dmcontrolFields = util.FieldSet{
	{"haltreq", 31, 31, util.FmtDec},
	{"resumereq", 30, 30, util.FmtDec},
	{"hartreset", 29, 29, util.FmtDec},
	{"ackhavereset", 28, 28, util.FmtDec},
	{"hasel", 26, 26, util.FmtDec},
	{"hartsello", 25, 16, util.FmtDec},
	{"hartselhi", 15, 6, util.FmtDec},
	{"setresethaltreq", 3, 3, util.FmtDec},
	{"clrresethaltreq", 2, 2, util.FmtDec},
	{"ndmreset", 1, 1, util.FmtDec},
	{"dmactive", 0, 0, util.FmtDec},
}

const haltreq = (1 << 31)
const ackhavereset = (1 << 28)
const hartsello = ((1 << 10) - 1) << 16
const hartselhi = ((1 << 10) - 1) << 6
const ndmreset = (1 << 1)
const dmactive = (1 << 0)

func (dbg *Debug) ndmResetPulse() error {
	// write 1
	err := dbg.setDmi(dmcontrol, ndmreset)
	if err != nil {
		return err
	}
	// write 0
	return dbg.clrDmi(dmcontrol, ndmreset)
}

func (dbg *Debug) dmActivePulse() error {
	// write 0
	err := dbg.clrDmi(dmcontrol, dmactive)
	if err != nil {
		return err
	}
	// write 1
	return dbg.setDmi(dmcontrol, dmactive)
}

//-----------------------------------------------------------------------------
// hart selection

// setHartSelect sets the hart select value in a dmcontrol value.
func setHartSelect(x uint32, id int) uint32 {
	x &= ^uint32(hartselhi | hartsello)
	hi := ((id >> 10) << 6) & hartselhi
	lo := (id << 16) & hartsello
	return x | uint32(hi) | uint32(lo)
}

// getHartSelect gets the hart select value from a dmcontrol value.
func getHartSelect(x uint32) int {
	return int((util.Bits(uint(x), 15, 6) << 10) | util.Bits(uint(x), 25, 16))
}

// selectHart sets the dmcontrol hartsel value.
func (dbg *Debug) selectHart(id int) error {
	if dbg.hartid == id {
		// already selected
		return nil
	}
	x, err := dbg.rdDmi(dmcontrol)
	if err != nil {
		return err
	}
	x = setHartSelect(x, id)
	dbg.hartid = id
	return dbg.wrDmi(dmcontrol, x)
}

//-----------------------------------------------------------------------------
// DM status

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

const anyhavereset = (1 << 18)
const anynonexistent = (1 << 14)
const anyunavail = (1 << 12)
const allhalted = (1 << 9)

//-----------------------------------------------------------------------------

var hartinfoFields = util.FieldSet{
	{"nscratch", 23, 20, util.FmtDec},
	{"dataaccess", 16, 16, util.FmtDec},
	{"datasize", 15, 12, util.FmtDec},
	{"dataaddr", 11, 0, util.FmtHex},
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

var abstractcsFields = util.FieldSet{
	{"progbufsize", 28, 24, util.FmtDec},
	{"busy", 12, 12, util.FmtDec},
	{"cmderr", 10, 8, util.FmtDec},
	{"datacount", 3, 0, util.FmtDec},
}

// command error
type cmdErr uint

// command error values
const (
	errOk           cmdErr = 0
	errBusy         cmdErr = 1
	errNotSupported cmdErr = 2
	errException    cmdErr = 3
	errHaltResume   cmdErr = 4
	errBusError     cmdErr = 5
	errReserved     cmdErr = 6
	errOther        cmdErr = 7
)

func (ce cmdErr) String() string {
	return [8]string{
		"ok",
		"busy",
		"not supported",
		"exception",
		"halt/resume",
		"bus error",
		"reserved",
		"other",
	}[ce]
}

// getError returns the error field of the command status.
func (cs cmdStatus) getError() cmdErr {
	return cmdErr(util.Bits(uint(cs), 10, 8))
}

// cmdErrorClr resets a command error.
func (dbg *Debug) cmdErrorClr() error {
	// write all-ones to the cmderr field.
	return dbg.wrDmi(abstractcs, 7<<8 /*cmderr*/)
}

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

type cmdStatus uint32

// isDone returns if a command has completed.
func (cs cmdStatus) isDone() bool {
	return cs&(1<<12 /*busy*/) == 0
}

const cmdTimeout = 500 * time.Millisecond

// cmdWait waits for command completion.
func (dbg *Debug) cmdWait(cs cmdStatus, to time.Duration) error {
	// wait for the command to complete
	t := time.Now().Add(to)
	for t.After(time.Now()) {
		// is the command complete?
		if cs.isDone() {
			// check for command error
			ce := cs.getError()
			if ce != errOk {
				// clear the error
				err := dbg.cmdErrorClr()
				if err != nil {
					return err
				}
				return fmt.Errorf("command error %s (%d)", ce, ce)
			}
			return nil
		}
		// wait a while
		time.Sleep(5 * time.Millisecond)
		// read the command status
		x, err := dbg.rdDmi(abstractcs)
		if err != nil {
			return err
		}
		cs = cmdStatus(x)
	}
	// timeout
	// reset the hart
	err := dbg.ndmResetPulse()
	if err != nil {
		return err

	}
	// reset the debug module
	err = dbg.dmActivePulse()
	if err != nil {
		return err
	}
	return errors.New("command timeout")
}

//-----------------------------------------------------------------------------
// debug module interface

// wrDmi reads a dmi register.
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

// wrDmi writes a dmi register.
func (dbg *Debug) wrDmi(addr uint, data uint32) error {
	ops := []dmiOp{
		dmiWr(addr, data),
		dmiEnd(),
	}
	_, err := dbg.dmiOps(ops)
	return err
}

// rmwDmi read/write/modify a dmi register.
func (dbg *Debug) rmwDmi(addr uint, mask, bits uint32) error {
	// read
	x, err := dbg.rdDmi(addr)
	if err != nil {
		return err
	}
	// modify
	x &= ^mask
	x |= bits
	// write
	return dbg.wrDmi(addr, x)
}

// setDmi sets bits in a dmi register.
func (dbg *Debug) setDmi(addr uint, bits uint32) error {
	return dbg.rmwDmi(addr, bits, bits)
}

// clrDmi clears bits in a dmi register.
func (dbg *Debug) clrDmi(addr uint, bits uint32) error {
	return dbg.rmwDmi(addr, bits, 0)
}

//-----------------------------------------------------------------------------
// read/write data value buffers

func (dbg *Debug) rdData32() (uint32, error) {
	if dbg.datacount < 1 {
		return 0, errors.New("need datacount >= 1 for 32-bit reads")
	}
	ops := []dmiOp{
		dmiRd(data0),
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (dbg *Debug) rdData64() (uint64, error) {
	if dbg.datacount < 2 {
		return 0, errors.New("need datacount >= 2 for 64-bit reads")
	}
	ops := []dmiOp{
		dmiRd(data0),
		dmiRd(data0 + 1),
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, err
	}
	return (uint64(data[1]) << 32) | uint64(data[0]), nil
}

func (dbg *Debug) rdData128() (uint64, uint64, error) {
	if dbg.datacount < 4 {
		return 0, 0, errors.New("need datacount >= 4 for 128-bit reads")
	}
	ops := []dmiOp{
		dmiRd(data0),
		dmiRd(data0 + 1),
		dmiRd(data0 + 2),
		dmiRd(data0 + 3),
		dmiEnd(),
	}
	data, err := dbg.dmiOps(ops)
	if err != nil {
		return 0, 0, err
	}
	lo := (uint64(data[1]) << 32) | uint64(data[0])
	hi := (uint64(data[3]) << 32) | uint64(data[2])
	return lo, hi, nil
}

//-----------------------------------------------------------------------------

// testBuffers tests dmi r/w buffers.
func (dbg *Debug) testBuffers(addr, n uint) error {

	// random write values
	wr := make([]uint32, n)
	for i := range wr {
		wr[i] = rand.Uint32()
	}

	// write to dmi registers
	for i := range wr {
		err := dbg.wrDmi(addr+uint(i), wr[i])
		if err != nil {
			return err
		}
	}

	// read back from dmi registers
	for i := range wr {
		x, err := dbg.rdDmi(addr + uint(i))
		if err != nil {
			return err
		}
		if x != wr[i] {
			return fmt.Errorf("dmi buffer r/w mismatch at addr 0x%x", addr+uint(i))
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
