//-----------------------------------------------------------------------------
/*

Segger J-Link Driver

*/
//-----------------------------------------------------------------------------

package jlink

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/deadsy/libjaylink"
	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------

const colorGreen = "\033[0;32m"
const colorNone = "\033[0m"

func logCallback(domain, msg string, user interface{}) {
	s := []string{colorGreen, domain, msg, colorNone}
	log.Debug.Printf("%s\n", strings.Join(s, ""))
}

//-----------------------------------------------------------------------------

// Jlink stores the J-Link library context.
type Jlink struct {
	ctx *libjaylink.Context
	dev []libjaylink.Device
}

// Init initializes the J-Link library.
func Init() (*Jlink, error) {
	// initialise the library
	ctx, err := libjaylink.Init()
	if err != nil {
		return nil, err
	}
	// setup the logging callback
	err = ctx.LogSetCallback(logCallback, nil)
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	err = ctx.LogSetLevel(libjaylink.LOG_LEVEL_DEBUG)
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	// discover devices
	err = ctx.DiscoveryScan(libjaylink.HIF_USB)
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	dev, err := ctx.GetDevices()
	if err != nil {
		ctx.Exit()
		return nil, err
	}
	// return the library context
	j := &Jlink{
		ctx: ctx,
		dev: dev,
	}
	return j, nil
}

// Shutdown closes the J-Link library.
func (j *Jlink) Shutdown() {
	j.ctx.FreeDevices(j.dev, true)
	j.ctx.Exit()
}

// NumDevices returns the number of devices discovered.
func (j *Jlink) NumDevices() int {
	return len(j.dev)
}

// DeviceByIndex returns a J-Link device by index number.
func (j *Jlink) DeviceByIndex(idx int) (*libjaylink.Device, error) {
	if idx < 0 || idx >= len(j.dev) {
		return nil, fmt.Errorf("device index %d out of range", idx)
	}
	return &j.dev[idx], nil
}

//-----------------------------------------------------------------------------

// pre-canned TAP state machine transitions
var xToIdle = bitstr.FromString("011111")     // any state -> run-test/idle
var idleToIRshift = bitstr.FromString("0011") // run-test/idle -> shift-ir
var idleToDRshift = bitstr.FromString("001")  // run-test/idle -> shift-dr
var xShiftToIdle = bitstr.FromString("011")   // shift-x -> run-test/idle

// Jtag is a driver for J-link JTAG operations.
type Jtag struct {
	dev     *libjaylink.Device
	hdl     *libjaylink.DeviceHandle
	version libjaylink.JtagVersion
}

// NewJtag returns a new J-Link JTAG driver.
func NewJtag(dev *libjaylink.Device, speed uint16) (*Jtag, error) {
	// get the device handle
	hdl, err := dev.Open()
	if err != nil {
		return nil, err
	}
	// get the device capabilities
	caps, err := hdl.GetAllCaps()
	if err != nil {
		hdl.Close()
		return nil, err
	}
	// get the JTAG command version
	version, err := hdl.GetJtagCommandVersion()
	if err != nil {
		hdl.Close()
		return nil, err
	}
	// check and select the target interface
	if !caps.HasCap(libjaylink.DEV_CAP_SELECT_TIF) {
		return nil, errors.New("jtag interface can't be selected")
	}
	itf, err := hdl.GetAvailableInterfaces()
	if err != nil {
		hdl.Close()
		return nil, err
	}
	if itf&(1<<libjaylink.TIF_JTAG) == 0 {
		hdl.Close()
		return nil, errors.New("jtag interface not available")
	}
	_, err = hdl.SelectInterface(libjaylink.TIF_JTAG)
	if err != nil {
		hdl.Close()
		return nil, err
	}
	// check the JTAG state
	state, err := hdl.GetHardwareStatus()
	if err != nil {
		hdl.Close()
		return nil, err
	}
	if state.TargetVoltage < 1500 {
		hdl.Close()
		return nil, fmt.Errorf("Target voltage is too low (%dmV). Is the target connected and powered?", state.TargetVoltage)
	}
	if state.Tres {
		hdl.Close()
		return nil, errors.New("Target ~SRST line asserted. Target is held in reset.")
	}
	// check the desired interface speed
	if caps.HasCap(libjaylink.DEV_CAP_GET_SPEEDS) {
		maxSpeed, err := hdl.GetMaxSpeed()
		if err != nil {
			hdl.Close()
			return nil, err
		}
		if speed > maxSpeed {
			hdl.Close()
			return nil, fmt.Errorf("JTAG speed is too high: %dkHz > %dkHz (max)", speed, maxSpeed)
		}
	}
	// set the interface speed
	err = hdl.SetSpeed(speed)
	if err != nil {
		hdl.Close()
		return nil, err
	}
	j := &Jtag{
		dev:     dev,
		hdl:     hdl,
		version: version,
	}
	return j, nil
}

// Close closes a J-Link JTAG driver.
func (j *Jtag) Close() error {
	return j.hdl.Close()
}

func (j *Jtag) String() string {
	s := []string{}
	hw, err := j.hdl.GetHardwareVersion()
	if err == nil {
		s = append(s, fmt.Sprintf("hardware %s", hw))
	}
	ver, err := j.hdl.GetFirmwareVersion()
	if err == nil {
		s = append(s, fmt.Sprintf("firmware %s", ver))
	}
	sn, err := j.dev.GetSerialNumber()
	if err == nil {
		s = append(s, fmt.Sprintf("serial %d", sn))
	}
	return strings.Join(s, "\n")
}

// jtagIO performs jtag IO operations.
func (j *Jtag) jtagIO(tms, tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	tdo, err := j.hdl.JtagIO(tms.GetBytes(), tdi.GetBytes(), j.version)
	if needTdo {
		return bitstr.FromBytes(tdo, tdi.Len()), err
	}
	return nil, err
}

// TestReset pulses the test reset line.
func (j *Jtag) TestReset(delay time.Duration) error {
	err := j.hdl.JtagClearTrst()
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return j.hdl.JtagSetTrst()
}

// SystemReset pulses the system reset line.
func (j *Jtag) SystemReset(delay time.Duration) error {
	err := j.hdl.ClearReset()
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return j.hdl.SetReset()
}

// TapReset resets the TAP state machine.
func (j *Jtag) TapReset() error {
	tdi := bitstr.Zeroes(xToIdle.Len())
	_, err := j.jtagIO(xToIdle, tdi, false)
	return err
}

// ScanIR scans bits through the JTAG IR chain
func (j *Jtag) ScanIR(tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	tms := bitstr.Null().Tail(idleToIRshift).Tail0(tdi.Len() - 1).Tail(xShiftToIdle)
	tdi = bitstr.Zeroes(idleToIRshift.Len()).Tail(tdi).Tail0(xShiftToIdle.Len() - 1)
	//log.Debug.Printf("tms %s\n", tms.LenBits())
	//log.Debug.Printf("tdi %s\n", tdi.LenBits())
	tdo, err := j.jtagIO(tms, tdi, needTdo)
	if err != nil {
		return nil, err
	}
	if needTdo {
		tdo.DropHead(idleToIRshift.Len()).DropTail(xShiftToIdle.Len() - 1)
		//log.Debug.Printf("tdo %s\n", tdo.LenBits())
		return tdo, nil
	}
	return nil, nil
}

// ScanDR scans bits through the JTAG DR chain
func (j *Jtag) ScanDR(tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	tms := bitstr.Null().Tail(idleToDRshift).Tail0(tdi.Len() - 1).Tail(xShiftToIdle)
	tdi = bitstr.Zeroes(idleToDRshift.Len()).Tail(tdi).Tail0(xShiftToIdle.Len() - 1)
	//log.Debug.Printf("tms %s\n", tms.LenBits())
	//log.Debug.Printf("tdi %s\n", tdi.LenBits())
	tdo, err := j.jtagIO(tms, tdi, needTdo)
	if err != nil {
		return nil, err
	}
	if needTdo {
		tdo.DropHead(idleToDRshift.Len()).DropTail(xShiftToIdle.Len() - 1)
		//log.Debug.Printf("tdo %s\n", tdo.LenBits())
		return tdo, nil
	}
	return nil, nil
}

//-----------------------------------------------------------------------------
