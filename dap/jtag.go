//-----------------------------------------------------------------------------
/*

CMSIS-DAP JTAG Driver

*/
//-----------------------------------------------------------------------------

package dap

import (
	"errors"
	"fmt"
	"time"

	"github.com/deadsy/hidapi"
	"github.com/deadsy/rvdbg/bitstr"
)

//-----------------------------------------------------------------------------

// Jtag is a driver for CMSIS-DAP JTAG operations.
type Jtag struct {
	dev *device
}

func (j *Jtag) String() string {
	return fmt.Sprintf("%s", j.dev)
}

// NewJtag returns a new CMSIS-DAP JTAG driver.
func NewJtag(devInfo *hidapi.DeviceInfo, speed int) (*Jtag, error) {

	// get the hid device
	hid, err := hidapi.Open(devInfo.VendorID, devInfo.ProductID, devInfo.SerialNumber)
	if err != nil {
		return nil, err
	}
	j := &Jtag{
		dev: newDevice(hid),
	}

	// make sure the CMSIS-DAP can do JTAG
	if !j.dev.hasCap(capJtag) {
		j.Close()
		return nil, errors.New("jtag not supported")
	}

	// connect in JTAG mode
	err = j.dev.connect(modeJtag)
	if err != nil {
		j.Close()
		return nil, err
	}

	// set the clock speed
	err = j.dev.setClock(speed)
	if err != nil {
		j.Close()
		return nil, err
	}

	return j, nil
}

// Close closes a CMSIS-DAP JTAG driver.
func (j *Jtag) Close() error {
	j.dev.disconnect()
	j.dev.close()
	return nil
}

// jtagIO performs jtag IO operations.
func (j *Jtag) jtagIO(tms, tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	return nil, errors.New("TODO")
}

// TestReset pulses the test reset line.
func (j *Jtag) TestReset(delay time.Duration) error {
	err := j.dev.setPins(0, pinTRST)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return j.dev.setPins(pinTRST, pinTRST)
}

// SystemReset pulses the system reset line.
func (j *Jtag) SystemReset(delay time.Duration) error {
	err := j.dev.setPins(0, pinSRST)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return j.dev.setPins(pinSRST, pinSRST)
}

// TapReset resets the TAP state machine.
func (j *Jtag) TapReset() error {
	return errors.New("TODO")
}

// ScanIR scans bits through the JTAG IR chain
func (j *Jtag) ScanIR(tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	return nil, errors.New("TODO")
}

// ScanDR scans bits through the JTAG DR chain
func (j *Jtag) ScanDR(tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	return nil, errors.New("TODO")
}

//-----------------------------------------------------------------------------
