//-----------------------------------------------------------------------------
/*

Segger J-Link SWD Driver

*/
//-----------------------------------------------------------------------------

package jlink

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/deadsy/jaylink"
)

//-----------------------------------------------------------------------------

// Swd is a driver for J-link SWD operations.
type Swd struct {
	dev *jaylink.Device
	hdl *jaylink.DeviceHandle
}

// NewSwd returns a new J-Link SWD driver.
func NewSwd(dev *jaylink.Device, speed uint16) (*Swd, error) {
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
	// check and select the target interface
	if !caps.HasCap(jaylink.DEV_CAP_SELECT_TIF) {
		return nil, errors.New("swd interface can't be selected")
	}
	itf, err := hdl.GetAvailableInterfaces()
	if err != nil {
		hdl.Close()
		return nil, err
	}
	if itf&(1<<jaylink.TIF_SWD) == 0 {
		hdl.Close()
		return nil, errors.New("swd interface not available")
	}
	_, err = hdl.SelectInterface(jaylink.TIF_SWD)
	if err != nil {
		hdl.Close()
		return nil, err
	}
	// check the hardware state
	state, err := hdl.GetHardwareStatus()
	if err != nil {
		hdl.Close()
		return nil, err
	}
	if state.TargetVoltage < 1500 {
		hdl.Close()
		return nil, fmt.Errorf("target voltage is too low (%dmV), is the target connected and powered?", state.TargetVoltage)
	}
	if state.Tres {
		hdl.Close()
		return nil, errors.New("target ~SRST line asserted, target is held in reset")
	}
	// check the desired interface speed
	if caps.HasCap(jaylink.DEV_CAP_GET_SPEEDS) {
		maxSpeed, err := hdl.GetMaxSpeed()
		if err != nil {
			hdl.Close()
			return nil, err
		}
		if speed > maxSpeed {
			hdl.Close()
			return nil, fmt.Errorf("SWD speed is too high: %dkHz > %dkHz (max)", speed, maxSpeed)
		}
	}
	// set the interface speed
	err = hdl.SetSpeed(speed)
	if err != nil {
		hdl.Close()
		return nil, err
	}
	swd := &Swd{
		dev: dev,
		hdl: hdl,
	}
	return swd, nil
}

// Close closes a J-Link SWD driver.
func (swd *Swd) Close() error {
	return swd.hdl.Close()
}

func (swd *Swd) String() string {
	s := []string{}
	hw, err := swd.hdl.GetHardwareVersion()
	if err == nil {
		s = append(s, fmt.Sprintf("hardware %s", hw))
	}
	ver, err := swd.hdl.GetFirmwareVersion()
	if err == nil {
		s = append(s, fmt.Sprintf("firmware %s", ver))
	}
	sn, err := swd.dev.GetSerialNumber()
	if err == nil {
		s = append(s, fmt.Sprintf("serial %d", sn))
	}
	return strings.Join(s, "\n")
}

// SystemReset pulses the system reset line.
func (swd *Swd) SystemReset(delay time.Duration) error {
	err := swd.hdl.ClearReset()
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return swd.hdl.SetReset()
}

//-----------------------------------------------------------------------------
