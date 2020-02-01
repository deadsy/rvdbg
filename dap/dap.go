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
	"strings"

	"github.com/deadsy/hidapi"
)

//-----------------------------------------------------------------------------
// cmsis-dap protocol constants

// General Commands
const cmdInfo = 0x00
const cmdLED = 0x01
const cmdConnect = 0x02
const cmdDisconnect = 0x03
const cmdWriteAbort = 0x08
const cmdDelay = 0x09
const cmdResetTarget = 0x0A

// cmdInfo
const InfoVendorID = 0x01        // Get the Vendor ID (string)
const InfoProductID = 0x02       // Get the Product ID (string)
const InfoSerialNumber = 0x03    // Get the Serial Number (string)
const InfoFirmwareVersion = 0x04 // Get the CMSIS-DAP Firmware Version (string)
const InfoVendorName = 0x05      // Get the Target Device Vendor (string)
const InfoDeviceName = 0x06      // Get the Target Device Name (string)
const InfoCapabilities = 0xF0    // Get information about the Capabilities (BYTE) of the Debug Unit
const InfoTestDomainTimer = 0xF1 // Get the Test Domain Timer parameter information
const InfoSwoTraceSize = 0xFD    // Get the SWO Trace Buffer Size (WORD)
const InfoMaxPacketCount = 0xFE  // Get the maximum Packet Count (BYTE)
const InfoMaxPacketSize = 0xFF   // Get the maximum Packet Size (SHORT)

// DAP Status Code
const dapOk = 0
const dapError = 0xFF

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
	for _, devInfo := range hidDevice {
		dev, err := hidapi.Open(devInfo.VendorID, devInfo.ProductID, "")
		if err != nil {
			continue
		}
		product, err := dev.GetProductString()
		if err != nil {
			continue
		}
		if strings.Contains(product, "CMSIS-DAP") {
			dapDevice = append(dapDevice, devInfo)
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

// Device is a DAP device.
type Device struct {
	dev *hidapi.Device
}

func (dev *Device) String() string {
	return fmt.Sprintf("%s", dev.dev)
}

// info issues a DAP_info command
func (dev *Device) info(id byte) ([]byte, error) {

	buf := []byte{cmdInfo, id}
	err := dev.dev.Write(0, buf)
	if err != nil {
		return nil, err
	}
	resp, err := dev.dev.ReadTimeout(0, 32, 1000)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//-----------------------------------------------------------------------------
