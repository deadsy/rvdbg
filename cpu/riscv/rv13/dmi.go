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
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------
// debug module registers

const data0 = 0x04 // Abstract Data 0-11
const data1 = 0x05
const data2 = 0x06
const data3 = 0x07
const data4 = 0x08
const data5 = 0x09
const data6 = 0x0a
const data7 = 0x0b
const data8 = 0x0c
const data9 = 0x0d
const data10 = 0x0e
const data11 = 0x0f

const progbuf0 = 0x20 // Program Buffer 0-15
const progbuf1 = 0x21
const progbuf2 = 0x22
const progbuf3 = 0x23
const progbuf4 = 0x24
const progbuf5 = 0x25
const progbuf6 = 0x26
const progbuf7 = 0x27
const progbuf8 = 0x28
const progbuf9 = 0x29
const progbuf10 = 0x2a
const progbuf11 = 0x2b
const progbuf12 = 0x2c
const progbuf13 = 0x2d
const progbuf14 = 0x2e
const progbuf15 = 0x2f

func progbuf(n int) uint {
	return progbuf0 + uint(n)
}

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

func newDMI() *soc.Device {
	return &soc.Device{
		Name: "DMI",
		Peripherals: []soc.Peripheral{
			{
				Name:  "DMI",
				Descr: "DMI Registers",
				Registers: []soc.Register{
					{Offset: data0, Name: "data0", Descr: "abstract data 0"},
					{Offset: data1, Name: "data1", Descr: "abstract data 1"},
					{Offset: data2, Name: "data2", Descr: "abstract data 2"},
					{Offset: data3, Name: "data3", Descr: "abstract data 3"},
					{Offset: data4, Name: "data4", Descr: "abstract data 4"},
					{Offset: data5, Name: "data5", Descr: "abstract data 5"},
					{Offset: data6, Name: "data6", Descr: "abstract data 6"},
					{Offset: data7, Name: "data7", Descr: "abstract data 7"},
					{Offset: data8, Name: "data8", Descr: "abstract data 8"},
					{Offset: data9, Name: "data9", Descr: "abstract data 9"},
					{Offset: data10, Name: "data10", Descr: "abstract data 10"},
					{Offset: data11, Name: "data11", Descr: "abstract data 11"},
					{Offset: progbuf0, Name: "progbuf0", Descr: "program buffer 0"},
					{Offset: progbuf1, Name: "progbuf1", Descr: "program buffer 1"},
					{Offset: progbuf2, Name: "progbuf2", Descr: "program buffer 2"},
					{Offset: progbuf3, Name: "progbuf3", Descr: "program buffer 3"},
					{Offset: progbuf4, Name: "progbuf4", Descr: "program buffer 4"},
					{Offset: progbuf5, Name: "progbuf5", Descr: "program buffer 5"},
					{Offset: progbuf6, Name: "progbuf6", Descr: "program buffer 6"},
					{Offset: progbuf7, Name: "progbuf7", Descr: "program buffer 7"},
					{Offset: progbuf8, Name: "progbuf8", Descr: "program buffer 8"},
					{Offset: progbuf9, Name: "progbuf9", Descr: "program buffer 9"},
					{Offset: progbuf10, Name: "progbuf10", Descr: "program buffer 10"},
					{Offset: progbuf11, Name: "progbuf11", Descr: "program buffer 11"},
					{Offset: progbuf12, Name: "progbuf12", Descr: "program buffer 12"},
					{Offset: progbuf13, Name: "progbuf13", Descr: "program buffer 13"},
					{Offset: progbuf14, Name: "progbuf14", Descr: "program buffer 14"},
					{Offset: progbuf15, Name: "progbuf15", Descr: "program buffer 15"},
					{Offset: dmcontrol,
						Name:  "dmcontrol",
						Descr: "debug module control",
						Fields: []soc.Field{
							{Name: "haltreq", Msb: 31, Lsb: 31},
							{Name: "resumereq", Msb: 30, Lsb: 30},
							{Name: "hartreset", Msb: 29, Lsb: 29},
							{Name: "ackhavereset", Msb: 28, Lsb: 28},
							{Name: "hasel", Msb: 26, Lsb: 26},
							{Name: "hartsello", Msb: 25, Lsb: 16},
							{Name: "hartselhi", Msb: 15, Lsb: 6},
							{Name: "setresethaltreq", Msb: 3, Lsb: 3},
							{Name: "clrresethaltreq", Msb: 2, Lsb: 2},
							{Name: "ndmreset", Msb: 1, Lsb: 1},
							{Name: "dmactive", Msb: 0, Lsb: 0},
						},
					},
					{Offset: dmstatus,
						Name:  "dmstatus",
						Descr: "debug module status",
						Fields: []soc.Field{
							{Name: "impebreak", Msb: 22, Lsb: 22},
							{Name: "allhavereset", Msb: 19, Lsb: 19},
							{Name: "anyhavereset", Msb: 18, Lsb: 18},
							{Name: "allresumeack", Msb: 17, Lsb: 17},
							{Name: "anyresumeack", Msb: 16, Lsb: 16},
							{Name: "allnonexistent", Msb: 15, Lsb: 15},
							{Name: "anynonexistent", Msb: 14, Lsb: 14},
							{Name: "allunavail", Msb: 13, Lsb: 13},
							{Name: "anyunavail", Msb: 12, Lsb: 12},
							{Name: "allrunning", Msb: 11, Lsb: 11},
							{Name: "anyrunning", Msb: 10, Lsb: 10},
							{Name: "allhalted", Msb: 9, Lsb: 9},
							{Name: "anyhalted", Msb: 8, Lsb: 8},
							{Name: "authenticated", Msb: 7, Lsb: 7},
							{Name: "authbusy", Msb: 6, Lsb: 6},
							{Name: "hasresethaltreq", Msb: 5, Lsb: 5},
							{Name: "confstrptrvalid", Msb: 4, Lsb: 4},
							{Name: "version", Msb: 3, Lsb: 0},
						},
					},
					{Offset: hartinfo,
						Name:  "hartinfo",
						Descr: "hart info",
						Fields: []soc.Field{
							{Name: "nscratch", Msb: 23, Lsb: 20},
							{Name: "dataaccess", Msb: 16, Lsb: 16},
							{Name: "datasize", Msb: 15, Lsb: 12},
							{Name: "dataaddr", Msb: 11, Lsb: 0},
						},
					},
					{Offset: hawindowsel, Name: "hawindowsel", Descr: "hart array window select"},
					{Offset: hawindow, Name: "hawindow", Descr: "hart array window"},
					{Offset: abstractcs,
						Name:  "abstractcs",
						Descr: "abstract control and status",
						Fields: []soc.Field{
							{Name: "progbufsize", Msb: 28, Lsb: 24},
							{Name: "busy", Msb: 12, Lsb: 12},
							{Name: "cmderr", Msb: 10, Lsb: 8},
							{Name: "datacount", Msb: 3, Lsb: 0},
						},
					},
					{Offset: command, Name: "command", Descr: "abstract command"},
					{Offset: abstractauto, Name: "abstractauto", Descr: "abstract command autoexec"},
					{Offset: confstrptr0, Name: "confstrptr0", Descr: "configuration string pointer 0"},
					{Offset: confstrptr1, Name: "confstrptr1", Descr: "configuration string pointer 1"},
					{Offset: confstrptr2, Name: "confstrptr2", Descr: "configuration string pointer 2"},
					{Offset: confstrptr3, Name: "confstrptr3", Descr: "configuration string pointer 3"},
					{Offset: nextdm, Name: "nextdm", Descr: "next debug module"},
					{Offset: authdata, Name: "authdata", Descr: "authentication data"},
					{Offset: haltsum0, Name: "haltsum0", Descr: "halt summary 0"},
					{Offset: haltsum1, Name: "haltsum1", Descr: "halt summary 1"},
					{Offset: haltsum2, Name: "haltsum2", Descr: "halt summary 2"},
					{Offset: haltsum3, Name: "haltsum3", Descr: "halt summary 3"},
					{Offset: sbcs, Name: "sbcs", Descr: "system bus access control and status"},
					{Offset: sbaddress0, Name: "sbaddress0", Descr: "system bus address 31:0"},
					{Offset: sbaddress1, Name: "sbaddress1", Descr: "system bus address 63:32"},
					{Offset: sbaddress2, Name: "sbaddress2", Descr: "system bus address 95:64"},
					{Offset: sbaddress3, Name: "sbaddress3", Descr: "system bus address 127:96"},
					{Offset: sbdata0, Name: "sbdata0", Descr: "system bus data 31:0"},
					{Offset: sbdata1, Name: "sbdata1", Descr: "system bus data 63:32"},
					{Offset: sbdata2, Name: "sbdata2", Descr: "system bus data 95:64"},
					{Offset: sbdata3, Name: "sbdata3", Descr: "system bus data 127:96"},
				},
			},
		},
	}
}

//-----------------------------------------------------------------------------
// DM control

const haltreq = (1 << 31)
const resumereq = (1 << 30)
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
	x, err := dbg.rdDmi(dmcontrol)
	if err != nil {
		return err
	}
	x = setHartSelect(x, id)
	return dbg.wrDmi(dmcontrol, x)
}

//-----------------------------------------------------------------------------
// DM status

const anyhavereset = (1 << 18)
const allresumeack = (1 << 17)
const anyresumeack = (1 << 16)
const anynonexistent = (1 << 14)
const anyunavail = (1 << 12)
const allrunning = (1 << 11)
const allhalted = (1 << 9)

// checkStatus checks the dmstatus register for a flag.
func (dbg *Debug) checkStatus(flag uint32) (bool, error) {
	x, err := dbg.rdDmi(dmstatus)
	if err != nil {
		return false, err
	}
	return x&flag != 0, nil
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
			data = append(data, uint32((x>>2)&util.Mask32))
		}
		// setup the next read
		read = dmi.isRead()
	}
	return data, nil
}

//-----------------------------------------------------------------------------
// abstract commands

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

const errClear = (7 << 8 /*cmderr*/)

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
	return dbg.wrDmi(abstractcs, errClear)
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

var sizeMap = map[uint]uint{
	8:   size8,
	16:  size16,
	32:  size32,
	64:  size64,
	128: size128,
}

type cmdFlag uint

const (
	cmdPostInc  = cmdFlag(1 << 19)               // post increment register number
	cmdPostExec = cmdFlag(1 << 18)               // post execute program buffer
	cmdRead     = cmdFlag(1 << 17)               // dataX = register
	cmdWrite    = cmdFlag((1 << 17) | (1 << 16)) // register = dataX
)

func cmdRegister(reg, size uint, flags cmdFlag) uint32 {
	return uint32((0 << 24) | (size << 20) | uint(flags) | (reg << 0))
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

// checkError checks a command status and clears and reports any error.
func (dbg *Debug) checkError(cs cmdStatus) error {
	ce := cs.getError()
	// are we done with no errors?
	if cs.isDone() && ce == errOk {
		return nil
	}
	// clear the error
	err := dbg.cmdErrorClr()
	if err != nil {
		return err
	}
	return fmt.Errorf("error: %s(%d)", ce, ce)
}

const cmdTimeout = 10 * time.Millisecond

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
				return fmt.Errorf("error: %s(%d)", ce, ce)
			}
			return nil
		}
		// wait a while
		time.Sleep(1 * time.Millisecond)
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

// rdDmi reads a dmi register.
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

// rmwDmi read/modify/write a dmi register.
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

// rdData32 reads a 32-bit data value.
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

// rdData64 reads a 64-bit data value.
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

// rdData128 reads a 128-bit data value.
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
// decode/display dmi registers

type dmiDriver struct {
	dbg *Debug
}

func (drv *dmiDriver) GetAddressSize() uint {
	return drv.dbg.abits
}

func (drv *dmiDriver) GetRegisterSize(r *soc.Register) uint {
	return 32
}

func (drv *dmiDriver) Rd(width, addr uint) (uint, error) {
	x, err := drv.dbg.rdDmi(addr)
	if err != nil {
		return 0, err
	}
	return uint(x), nil
}

func (dbg *Debug) dmiDump() (string, error) {
	p := dbg.dmiDevice.GetPeripheral("DMI")
	drv := &dmiDriver{dbg}
	return p.Display(drv, nil, false), nil
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
			return fmt.Errorf("w/r mismatch at 0x%x", addr+uint(i))
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
