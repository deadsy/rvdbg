//-----------------------------------------------------------------------------
/*

CMSIS-DAP SWD Driver

*/
//-----------------------------------------------------------------------------

package daplink

import (
	"errors"
	"fmt"
	"time"

	"github.com/deadsy/hidapi"
	"github.com/deadsy/rvdbg/swd"
)

//-----------------------------------------------------------------------------

// Swd is a driver for CMSIS-DAP SWD operations.
type Swd struct {
	dev *device
}

func (drv *Swd) String() string {
	return fmt.Sprintf("%s", drv.dev)
}

// NewSwd returns a new CMSIS-DAP SWD driver.
func NewSwd(devInfo *hidapi.DeviceInfo, speed int) (*Swd, error) {

	// get the hid device
	hid, err := hidapi.Open(devInfo.VendorID, devInfo.ProductID, devInfo.SerialNumber)
	if err != nil {
		return nil, err
	}

	dev, err := newDevice(hid)
	if err != nil {
		hid.Close()
		return nil, err
	}

	drv := &Swd{
		dev: dev,
	}

	// make sure the CMSIS-DAP can do SWD
	if !drv.dev.hasCap(capSwd) {
		drv.Close()
		return nil, errors.New("swd not supported")
	}

	// connect in SWD mode
	err = drv.dev.cmdConnect(modeSwd)
	if err != nil {
		drv.Close()
		return nil, err
	}

	// set the clock speed
	err = drv.dev.cmdSwjClock(speed)
	if err != nil {
		drv.Close()
		return nil, err
	}

	return drv, nil
}

// Close closes a CMSIS-DAP SWD driver.
func (drv *Swd) Close() error {
	drv.dev.cmdDisconnect()
	drv.dev.close()
	return nil
}

// GetState returns the SWD hardware state.
func (drv *Swd) GetState() (*swd.State, error) {
	pins, err := drv.dev.getPins()
	if err != nil {
		return nil, err
	}
	return &swd.State{
		TargetVoltage: -1, // not supported
		Srst:          pins&pinSRST != 0,
	}, nil
}

// SystemReset pulses the system reset line.
func (drv *Swd) SystemReset(delay time.Duration) error {
	err := drv.dev.setPins(pinSRST)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return drv.dev.clrPins(pinSRST)
}

//-----------------------------------------------------------------------------
