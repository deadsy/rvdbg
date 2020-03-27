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
)

//-----------------------------------------------------------------------------

// Swd is a driver for CMSIS-DAP SWD operations.
type Swd struct {
	dev *device
}

func (swd *Swd) String() string {
	return fmt.Sprintf("%s", swd.dev)
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

	swd := &Swd{
		dev: dev,
	}

	// make sure the CMSIS-DAP can do SWD
	if !swd.dev.hasCap(capSwd) {
		swd.Close()
		return nil, errors.New("swd not supported")
	}

	// connect in SWD mode
	err = swd.dev.cmdConnect(modeSwd)
	if err != nil {
		swd.Close()
		return nil, err
	}

	// set the clock speed
	err = swd.dev.cmdSwjClock(speed)
	if err != nil {
		swd.Close()
		return nil, err
	}

	return swd, nil
}

// Close closes a CMSIS-DAP SWD driver.
func (swd *Swd) Close() error {
	swd.dev.cmdDisconnect()
	swd.dev.close()
	return nil
}

// SystemReset pulses the system reset line.
func (swd *Swd) SystemReset(delay time.Duration) error {
	err := swd.dev.setPins(pinSRST)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return swd.dev.clrPins(pinSRST)
}

//-----------------------------------------------------------------------------
