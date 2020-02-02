//-----------------------------------------------------------------------------
/*

CMSIS-DAP SWD Driver

*/
//-----------------------------------------------------------------------------

package dap

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
	swd := &Swd{
		dev: newDevice(hid),
	}

	// make sure the CMSIS-DAP can do SWD
	if !swd.dev.hasCap(capSwd) {
		swd.Close()
		return nil, errors.New("swd not supported")
	}

	// connect in SWD mode
	err = swd.dev.connect(modeSwd)
	if err != nil {
		swd.Close()
		return nil, err
	}

	// set the clock speed
	err = swd.dev.setClock(speed)
	if err != nil {
		swd.Close()
		return nil, err
	}

	return swd, nil
}

// Close closes a CMSIS-DAP SWD driver.
func (swd *Swd) Close() error {
	swd.dev.disconnect()
	swd.dev.close()
	return nil
}

// SystemReset pulses the system reset line.
func (swd *Swd) SystemReset(delay time.Duration) error {
	err := swd.dev.setPins(0, pinSRST)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return swd.dev.setPins(pinSRST, pinSRST)
}

//-----------------------------------------------------------------------------
