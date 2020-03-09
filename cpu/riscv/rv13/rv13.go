//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Functions

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"errors"
	"fmt"
	"strings"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/decode"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/util"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------

const irDtmcs = 0x10 // debug transport module control and status
const irDmi = 0x11   // debug module interface access

const drDtmcsLength = 32

//-----------------------------------------------------------------------------

// Debug is a RISC-V 0.13 debugger.
// It implements the rv.Debug interface.
type Debug struct {
	dev             *jtag.Device
	hart            []*hartInfo // implemented harts
	hartid          int         // currently selected hart
	ir              uint        // cache of ir value
	irlen           int         // IR length
	drDmiLength     int         // DR length for dmi
	abits           uint        // address bits in dtmcs
	idle            uint        // idle value in dtmcs
	progbufsize     uint        // number of progbuf words implemented
	datacount       uint        // number of data words implemented
	autoexecprogbuf bool        // can we autoexec on progbufX access?
	autoexecdata    bool        // can we autoexec on dataX access?
	sbasize         uint        // width of system bus address (0 = no access)
	hartsellen      uint        // hart select length 0..20
	impebreak       uint        // implicit ebreak in progbuf
}

func (dbg *Debug) String() string {
	s := [][]string{}
	s = append(s, []string{"version", "0.13"})
	s = append(s, []string{"idle cycles", fmt.Sprintf("%d", dbg.idle)})
	s = append(s, []string{"sbasize", fmt.Sprintf("%d bits", dbg.sbasize)})
	s = append(s, []string{"progbufsize", fmt.Sprintf("%d words", dbg.progbufsize)})
	s = append(s, []string{"datacount", fmt.Sprintf("%d words", dbg.datacount)})
	s = append(s, []string{"autoexecprogbuf", fmt.Sprintf("%t", dbg.autoexecprogbuf)})
	s = append(s, []string{"autoexecdata", fmt.Sprintf("%t", dbg.autoexecdata)})
	return cli.TableString(s, []int{0, 0}, 1)
}

// New returns a RISC-V 0.13 debugger.
func New(dev *jtag.Device) (rv.Debug, error) {

	dbg := &Debug{
		dev:   dev,
		irlen: dev.GetIRLength(),
	}

	// get dtmcs
	dtmcs, err := dbg.rdDtmcs()
	if err != nil {
		return nil, err
	}
	// check the version
	if util.Bits(dtmcs, 3, 0) != 1 {
		return nil, errors.New("unknown dtmcs version")
	}
	// get the dmi address bits
	dbg.abits = util.Bits(dtmcs, 9, 4)
	// get the cycles to wait in run-time/idle.
	dbg.idle = util.Bits(dtmcs, 14, 12)

	// check dmi for the correct length
	dbg.drDmiLength = 33 + int(dbg.abits) + 1
	_, err = dev.CheckDR(irDmi, dbg.drDmiLength)
	if err != nil {
		return nil, err
	}

	// reset the debug module
	err = dbg.wrDtmcs(dmihardreset | dmireset)
	if err != nil {
		return nil, err
	}

	// make the dmi active
	err = dbg.wrDmi(dmcontrol, 0)
	if err != nil {
		return nil, err
	}
	err = dbg.wrDmi(dmcontrol, dmactive)
	if err != nil {
		return nil, err
	}

	// write all-ones to hartsel
	err = dbg.selectHart((1 << 20) - 1)
	if err != nil {
		return nil, err
	}

	// read back dmcontrol
	x, err := dbg.rdDmi(dmcontrol)
	if err != nil {
		return nil, err
	}
	// check dmi is active
	if x&dmactive == 0 {
		return nil, errors.New("dmi is not active")
	}
	// work out hartsellen
	hartsel := getHartSelect(x)
	for hartsel&1 != 0 {
		dbg.hartsellen++
		hartsel >>= 1
	}
	log.Info.Printf(fmt.Sprintf("hartsellen %d", dbg.hartsellen))

	// check dmstatus fields
	x, err = dbg.rdDmi(dmstatus)
	if err != nil {
		return nil, err
	}
	// check version
	if util.Bits(uint(x), 3, 0) != 2 {
		return nil, errors.New("unknown dmstatus version")
	}
	// check authentication
	if util.Bit(uint(x), 7) != 1 {
		return nil, errors.New("debugger is not authenticated")
	}
	// implicit ebreak after progbuf
	dbg.impebreak = util.Bit(uint(x), 22)

	// work out the system bus address size
	x, err = dbg.rdDmi(sbcs)
	if err != nil {
		return nil, err
	}
	dbg.sbasize = util.Bits(uint(x), 11, 5)
	log.Info.Printf(fmt.Sprintf("sbasize %d", dbg.sbasize))

	// work out how many program and data words we have
	x, err = dbg.rdDmi(abstractcs)
	if err != nil {
		return nil, err
	}
	dbg.progbufsize = util.Bits(uint(x), 28, 24)
	dbg.datacount = util.Bits(uint(x), 3, 0)

	// check progbuf/impebreak consistency
	if dbg.progbufsize == 1 && dbg.impebreak != 1 {
		return nil, fmt.Errorf("progbufsize == 1 and impebreak != 1")
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
	// turn off autoexec
	err = dbg.wrDmi(abstractauto, 0)
	if err != nil {
		return nil, err
	}

	log.Info.Printf(fmt.Sprintf("progbufsize %d impebreak %d autoexecprogbuf %t", dbg.progbufsize, dbg.impebreak, dbg.autoexecprogbuf))
	log.Info.Printf(fmt.Sprintf("datacount %d autoexecdata %t", dbg.datacount, dbg.autoexecdata))

	// clear any pending command errors
	err = dbg.cmdErrorClr()
	if err != nil {
		return nil, err
	}

	// 1st pass: enumerate the harts
	maxHarts := 1 << dbg.hartsellen
	for id := 0; id < maxHarts; id++ {
		// select the hart
		err := dbg.selectHart(id)
		if err != nil {
			return nil, err
		}
		// get the hart status
		x, err = dbg.rdDmi(dmstatus)
		if err != nil {
			return nil, err
		}
		if x&anynonexistent != 0 {
			// this hart does not exist - we're done
			break
		}
		// add a hart to the list
		dbg.hart = append(dbg.hart, dbg.newHart(id))
	}

	if len(dbg.hart) == 0 {
		return nil, errors.New("no harts found")
	}

	// 2nd pass: examine each hart
	log.Info.Printf("%d hart(s) found", len(dbg.hart))
	for i := range dbg.hart {
		err := dbg.hart[i].examine()
		if err != nil {
			return nil, err
		}
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

var dtmcsFields = decode.FieldSet{
	{"dmihardreset", 17, 17, decode.FmtDec},
	{"dmireset", 16, 16, decode.FmtDec},
	{"idle", 14, 12, decode.FmtDec},
	{"dmistat", 11, 10, decode.FmtDec},
	{"abits", 9, 4, decode.FmtDec},
	{"version", 3, 0, decode.FmtDec},
}

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
// hart control

// GetHartCount returns the number of harts.
func (dbg *Debug) GetHartCount() int {
	return len(dbg.hart)
}

// GetHartInfo returns the hart info. id < 0 gives the current hart info.
func (dbg *Debug) GetHartInfo(id int) (*rv.HartInfo, error) {
	if id < 0 || id >= len(dbg.hart) {
		return nil, errors.New("hart id is out of range")
	}
	return &dbg.hart[id].info, nil
}

func (dbg *Debug) GetCurrentHart() *rv.HartInfo {
	return &dbg.hart[dbg.hartid].info
}

// SetCurrentHart sets the current hart.
func (dbg *Debug) SetCurrentHart(id int) (*rv.HartInfo, error) {
	if id < 0 || id >= len(dbg.hart) {
		return nil, errors.New("hart id is out of range")
	}
	err := dbg.selectHart(id)
	dbg.hartid = id
	return &dbg.hart[dbg.hartid].info, err
}

// HaltHart halts the current hart.
func (dbg *Debug) HaltHart() error {
	_, err := dbg.halt()
	halted, _ := dbg.isHalted()
	if halted {
		dbg.hart[dbg.hartid].info.State = rv.Halted
	}
	return err
}

// ResumeHart resumes the current hart.
func (dbg *Debug) ResumeHart() error {
	_, err := dbg.resume()
	running, _ := dbg.isRunning()
	if running {
		dbg.hart[dbg.hartid].info.State = rv.Running
	}
	return err
}

//-----------------------------------------------------------------------------

// GetInfo returns a string of debugger information.
func (dbg *Debug) GetInfo() string {
	s := []string{}
	s = append(s, dbg.String())
	s = append(s, "")
	dump, err := dbg.dmiDump()
	if err != nil {
		s = append(s, fmt.Sprintf("unable to get dmi registers: %v", err))
	} else {
		s = append(s, dump)
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------

// Test is a test routine.
func (dbg *Debug) Test1() string {
	s := []string{}

	// test progbuf buffers
	err := dbg.testBuffers(progbuf0, dbg.progbufsize)
	if err != nil {
		s = append(s, fmt.Sprintf("progbuf: %v", err))
	} else {
		s = append(s, "progbuf: passed")
	}

	// test data buffers
	err = dbg.testBuffers(data0, dbg.datacount)
	if err != nil {
		s = append(s, fmt.Sprintf("data: %v", err))
	} else {
		s = append(s, "data: passed")
	}

	return strings.Join(s, "\n")
}

const testReg = 25

// Test2 is a test routine.
func (dbg *Debug) Test2() string {
	s := []string{}

	err := dbg.WrGPR(testReg, 0, 0xdeadbeef)
	s = append(s, fmt.Sprintf("wr %v", err))

	val, err := dbg.RdGPR(testReg, 0)
	s = append(s, fmt.Sprintf("rd %x %v", val, err))

	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
