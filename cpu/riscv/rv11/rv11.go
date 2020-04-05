//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 Functions

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"errors"
	"fmt"
	"strings"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------

const irDtmcontrol = 0x10 // debug transport module control
const drDtmcontrolLength = 32

const irDbus = 0x11 // debug bus access

//-----------------------------------------------------------------------------

// Debug is a RISC-V 0.11 debugger. It implements the rv.Debug interface.
type Debug struct {
	dev          *jtag.Device
	dbusDevice   *soc.Device // dbus device for decode/display
	cache        *ramCache   // cache of debug ram words
	hart         []*hartInfo // implemented harts
	hartid       int         // currently selected hart
	ir           uint        // cache of ir value
	irlen        int         // IR length
	drDbusLength int         // DR length for dbus
	abits        uint        // address bits in dtmcontrol
	idle         uint        // idle value in dtmcontrol
	dramsize     uint        // number of debug ram words implemented
	dbusops      uint        // running count of total dbus operations
}

func (dbg *Debug) String() string {
	s := [][]string{}
	s = append(s, []string{"version", "0.11"})
	s = append(s, []string{"idle cycles", fmt.Sprintf("%d", dbg.idle)})
	s = append(s, []string{"dramsize", fmt.Sprintf("%d words", dbg.dramsize)})
	s = append(s, []string{"dbusops", fmt.Sprintf("%d", dbg.dbusops)})
	return cli.TableString(s, []int{0, 0}, 1)
}

// New returns a RISC-V 0.11 debugger.
func New(dev *jtag.Device) (*Debug, error) {
	log.Info.Printf("0.11 debug module")
	dbg := &Debug{
		dev:        dev,
		irlen:      dev.GetIRLength(),
		dbusDevice: newDBUS().Setup(),
	}

	// get dtmcontrol
	dtmcontrol, err := dbg.rdDtmcontrol()
	if err != nil {
		return nil, err
	}
	if dtmcontrol == 0 {
		return nil, errors.New("bad value for dtmcontrol (0)")
	}
	// check the version
	if util.Bits(dtmcontrol, 3, 0) != 0 {
		return nil, errors.New("unknown dtmcontrol version")
	}
	// get the dbus address bits
	dbg.abits = (util.Bits(dtmcontrol, 14, 13) << 4) | util.Bits(dtmcontrol, 7, 4)
	// get the cycles to wait in run-time/idle.
	dbg.idle = util.Bits(dtmcontrol, 12, 10)

	// check dbus for the correct length
	dbg.drDbusLength = 35 + int(dbg.abits) + 1
	_, err = dev.CheckDR(irDbus, dbg.drDbusLength)
	if err != nil {
		return nil, err
	}

	// reset the debug module
	err = dbg.wrDtmcontrol(dbusreset)
	if err != nil {
		return nil, err
	}

	// get dminfo
	x, err := dbg.rdDbus(dminfo)
	if err != nil {
		return nil, err
	}
	// check dminfo.version
	version := (util.Bits(x, 7, 6) << 2) | (util.Bits(x, 1, 0) << 0)
	if version != 1 {
		return nil, fmt.Errorf("dminfo.version expected 1, actual %d", version)
	}
	// get number of words of debug ram
	dbg.dramsize = util.Bits(x, 15, 10) + 1
	// check dminfo.authtype
	authtype := util.Bits(x, 3, 2)
	if authtype != 0 {
		return nil, fmt.Errorf("dminfo.authtype %d not supported", authtype)
	}

	// create the debug ram cache
	dbg.cache, err = dbg.newCache(debugRamStart, dbg.dramsize)
	if err != nil {
		return nil, err
	}

	// 1st pass: enumerate the harts
	maxHarts := 32
	for id := 0; id < maxHarts; id++ {
		// select the hart
		err := dbg.selectHart(id)
		if err != nil {
			return nil, err
		}
		// try to get MXLEN
		_, err = dbg.getMXLEN()
		if err != nil {
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

	// select the 0th hart
	_, err = dbg.SetCurrentHart(0)
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
// dtmcontrol

const dbusreset = (1 << 16)

// rdDtmcontrol reads the dtmcontrol register.
func (dbg *Debug) rdDtmcontrol() (uint, error) {
	err := dbg.wrIR(irDtmcontrol)
	if err != nil {
		return 0, err
	}
	tdo, err := dbg.dev.RdWrDR(bitstr.Zeros(drDtmcontrolLength), 0)
	if err != nil {
		return 0, err
	}
	return tdo.Split([]int{drDtmcontrolLength})[0], nil
}

// wrDtmcontrol writes the dtmcontrol register.
func (dbg *Debug) wrDtmcontrol(val uint) error {
	err := dbg.wrIR(irDtmcontrol)
	if err != nil {
		return err
	}
	return dbg.dev.WrDR(bitstr.FromUint(val, drDtmcontrolLength), 0)
}

//-----------------------------------------------------------------------------
// hart control

// GetHartCount returns the number of harts.
func (dbg *Debug) GetHartCount() int {
	return len(dbg.hart)
}

// GetHartInfo returns hart info.
func (dbg *Debug) GetHartInfo(id int) (*rv.HartInfo, error) {
	if id < 0 || id >= len(dbg.hart) {
		return nil, errors.New("hart id is out of range")
	}
	return &dbg.hart[id].info, nil
}

// GetCurrentHart returns the current hart info.
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

// GetPrompt returns a target prompt string.
func (dbg *Debug) GetPrompt(name string) string {
	hi := dbg.GetCurrentHart()
	state := []rune{'h', 'r'}[util.BoolToInt(hi.State == rv.Running)]
	return fmt.Sprintf("%s.%d%c> ", name, hi.ID, state)
}

//-----------------------------------------------------------------------------

// Test1 is a test routine.
func (dbg *Debug) Test1() string {
	s := []string{}
	// test debug ram buffers
	err := dbg.testBuffers(ram0, dbg.dramsize)
	if err != nil {
		s = append(s, fmt.Sprintf("debug ram: %v", err))
	} else {
		s = append(s, "debug ram: passed")
	}
	return strings.Join(s, "\n")
}

// Test2 is a test routine.
func (dbg *Debug) Test2() string {

	//dbg.WrGPR(30, 0, 0x123456789abcdef)
	dbg.WrFPR(16, 0, 0xdeadbeefcafebabe)
	//dbg.RdFPR(25, 0)

	return ""
}

//-----------------------------------------------------------------------------
