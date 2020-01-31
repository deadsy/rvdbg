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
	dev *hidapi.Device
}

func (j *Jtag) String() string {
	return fmt.Sprintf("%s", j.dev)
}

// NewJtag returns a new CMSIS-DAP JTAG driver.
func NewJtag(devInfo *hidapi.DeviceInfo, speed, volts uint16) (*Jtag, error) {

	dev, err := hidapi.Open(devInfo.VendorID, devInfo.ProductID, devInfo.SerialNumber)
	if err != nil {
		return nil, err
	}

	j := &Jtag{
		dev: dev,
	}
	return j, nil
}

// Close closes a CMSIS-DAP JTAG driver.
func (j *Jtag) Close() error {
	return j.dev.Close()
}

// jtagIO performs jtag IO operations.
func (j *Jtag) jtagIO(tms, tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) {
	return nil, errors.New("TODO")
}

// TestReset pulses the test reset line.
func (j *Jtag) TestReset(delay time.Duration) error {
	return errors.New("TODO")
}

// SystemReset pulses the system reset line.
func (j *Jtag) SystemReset(delay time.Duration) error {
	return errors.New("TODO")
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
