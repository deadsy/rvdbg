//-----------------------------------------------------------------------------
/*

CMSIS-DAP Driver

This package implements CMSIS-DAP JTAG/SWD drivers using the hidapi library.

*/
//-----------------------------------------------------------------------------

package dap

import (
	"errors"
	"fmt"

	"github.com/deadsy/hidapi"
)

//-----------------------------------------------------------------------------

// Dap stores the DAP library context.
type Dap struct {
	device []*hidapi.DeviceInfo // CMSIS-DAP devices found
}

// Init initializes the DAP library.
func Init() (*Dap, error) {

	err := hidapi.Init()
	if err != nil {
		return nil, err
	}

	// get all HID devices
	hidDevice := hidapi.Enumerate(0, 0)
	if len(hidDevice) == 0 {
		hidapi.Exit()
		return nil, errors.New("no HID devices found")
	}

	// filter in the CMSIS-DAP devices
	dapDevice := []*hidapi.DeviceInfo{}
	for _, dInfo := range hidDevice {
		dev, err := hidapi.Open(dInfo.VendorID, dInfo.ProductID, "")
		if err != nil {
			continue
		}
		product, err := dev.GetProductString()
		if err != nil {
			continue
		}
		if product == "CMSIS-DAP" {
			dapDevice = append(dapDevice, dInfo)
		}
		dev.Close()
	}

	dap := &Dap{
		device: dapDevice,
	}

	return dap, nil
}

// Shutdown closes the DAP library.
func (dap *Dap) Shutdown() {
	hidapi.Exit()
}

// NumDevices returns the number of devices discovered.
func (dap *Dap) NumDevices() int {
	return len(dap.device)
}

// DeviceByIndex returns DAP device information by index number.
func (dap *Dap) DeviceByIndex(idx int) (*hidapi.DeviceInfo, error) {
	if idx < 0 || idx >= len(dap.device) {
		return nil, fmt.Errorf("device index %d out of range", idx)
	}
	return dap.device[idx], nil
}

//-----------------------------------------------------------------------------
