//-----------------------------------------------------------------------------
/*

CMSIS-DAP Driver

This package implements CMSIS-DAP JTAG/SWD drivers using the hidapi library.

*/
//-----------------------------------------------------------------------------

package dap

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/deadsy/hidapi"
)

//-----------------------------------------------------------------------------
// cmsis-dap protocol constants

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
// CMSIS-DAP Device

const dapReport = 0
const usbTimeout = 500 // milliseconds

type device struct {
	hid     *hidapi.Device
	caps    capabilities
	pktSize int
}

func newDevice(hid *hidapi.Device) *device {
	dev := &device{
		hid:     hid,
		pktSize: 64,
	}
	// get the max packet size
	maxPktSize, err := dev.getMaxPacketSize()
	if err == nil && int(maxPktSize) < dev.pktSize {
		dev.pktSize = int(maxPktSize)
	}
	// get the capabilities
	caps, _ := dev.getCapabilities()
	dev.caps = caps
	return dev
}

func (dev *device) String() string {
	return fmt.Sprintf("%s", dev.hid)
}

func (dev *device) txrx(buf []byte) ([]byte, error) {
	//fmt.Printf("tx (%d) %v\n", len(buf), buf)
	err := dev.hid.Write(buf)
	if err != nil {
		return nil, err
	}
	rx, err := dev.hid.ReadTimeout(dapReport, dev.pktSize, usbTimeout)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("rx (%d) %v\n", len(rx), rx)
	return rx, nil
}

func (dev *device) close() error {
	return dev.hid.Close()
}

//-----------------------------------------------------------------------------
// General Commands

// General Commands
const cmdInfo = 0x00
const cmdHostStatus = 0x01
const cmdConnect = 0x02
const cmdDisconnect = 0x03
const cmdWriteAbort = 0x08
const cmdDelay = 0x09
const cmdResetTarget = 0x0A

func (dev *device) info(id byte) ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdInfo, id})
}

func (dev *device) hostStatus(typ, status byte) ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdHostStatus, typ, status})
}

func (dev *device) connect(port byte) ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdConnect, port})
}

func (dev *device) disconnect() ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdDisconnect})
}

func (dev *device) writeAbort(index byte, abort uint32) ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdWriteAbort, index, byte(abort), byte(abort >> 8), byte(abort >> 16), byte(abort >> 24)})
}

func (dev *device) delay(delay uint16) ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdDelay, byte(delay), byte(delay >> 8)})
}

func (dev *device) resetTarget() ([]byte, error) {
	return dev.txrx([]byte{dapReport, cmdResetTarget})
}

//-----------------------------------------------------------------------------
// DAP Information Commands

// cmdInfo
const infoVendorID = 0x01        // Get the Vendor ID (string)
const infoProductID = 0x02       // Get the Product ID (string)
const infoSerialNumber = 0x03    // Get the Serial Number (string)
const infoFirmwareVersion = 0x04 // Get the CMSIS-DAP Firmware Version (string)
const infoVendorName = 0x05      // Get the Target Device Vendor (string)
const infoDeviceName = 0x06      // Get the Target Device Name (string)
const infoCapabilities = 0xF0    // Get information about the Capabilities (BYTE) of the Debug Unit
const infoTestDomainTimer = 0xF1 // Get the Test Domain Timer parameter information
const infoSwoTraceSize = 0xFD    // Get the SWO Trace Buffer Size (WORD)
const infoMaxPacketCount = 0xFE  // Get the maximum Packet Count (BYTE)
const infoMaxPacketSize = 0xFF   // Get the maximum Packet Size (SHORT)

// getString gets a string type information item.
func (dev *device) getString(id byte) (string, error) {
	rx, err := dev.info(id)
	if err != nil {
		return "", err
	}
	if len(rx) < 2 || rx[0] != 0 {
		return "", errors.New("bad info response")
	}
	if rx[1] == 0 {
		return "", errors.New("no information")
	}
	n := int(rx[1])
	if len(rx) < n+2 {
		return "", errors.New("response too short")
	}
	if rx[n+1] != 0 {
		return "", errors.New("string not null terminated")
	}
	return string(rx[2 : 2+n-1]), nil
}

// getByte gets a byte type information item.
func (dev *device) getByte(id byte) (byte, error) {
	rx, err := dev.info(id)
	if err != nil {
		return 0, err
	}
	if len(rx) < 3 || rx[0] != 0 || rx[1] != 1 {
		return 0, errors.New("bad info response")
	}
	return rx[2], nil
}

// getShort gets a short type information item.
func (dev *device) getShort(id byte) (uint16, error) {
	rx, err := dev.info(id)
	if err != nil {
		return 0, err
	}
	if len(rx) < 4 || rx[0] != 0 || rx[1] != 2 {
		return 0, errors.New("bad info response")
	}
	return binary.LittleEndian.Uint16(rx[2:4]), nil
}

// getWord gets a word type information item.
func (dev *device) getWord(id byte) (uint32, error) {
	rx, err := dev.info(id)
	if err != nil {
		return 0, err
	}
	if len(rx) < 6 || rx[0] != 0 || rx[1] != 4 {
		return 0, errors.New("bad info response")
	}
	return binary.LittleEndian.Uint32(rx[2:6]), nil
}

// GetVendorID gets the Vendor ID
func (dev *device) getVendorID() (string, error) {
	return dev.getString(infoVendorID)
}

// GetProductID gets the Product ID
func (dev *device) getProductID() (string, error) {
	return dev.getString(infoProductID)
}

// GetSerialNumber gets the Serial Number
func (dev *device) getSerialNumber() (string, error) {
	return dev.getString(infoSerialNumber)
}

// GetFirmwareVersion gets the CMSIS-DAP Firmware Version
func (dev *device) getFirmwareVersion() (string, error) {
	return dev.getString(infoFirmwareVersion)
}

// GetVendorName gets the Target Device Vendor
func (dev *device) getVendorName() (string, error) {
	return dev.getString(infoVendorName)
}

// GetDeviceName gets the Target Device Name
func (dev *device) getDeviceName() (string, error) {
	return dev.getString(infoDeviceName)
}

// GetTestDomainTimer gets the Test Domain Timer parameter information
func (dev *device) getTestDomainTimer() (uint32, error) {
	return dev.getWord(infoTestDomainTimer)
}

// GetSwoTraceSize gets the SWO Trace Buffer Size
func (dev *device) getSwoTraceSize() (uint32, error) {
	return dev.getWord(infoSwoTraceSize)
}

// GetMaxPacketCount gets the maximum Packet Count
func (dev *device) getMaxPacketCount() (byte, error) {
	return dev.getByte(infoMaxPacketCount)
}

// GetMaxPacketSize gets the maximum Packet Size
func (dev *device) getMaxPacketSize() (uint16, error) {
	return dev.getShort(infoMaxPacketSize)
}

//-----------------------------------------------------------------------------
// Device Capabilities

type capabilities uint

const (
	capSwd              capabilities = (1 << 0) // SWD Serial Wire Debug
	capJtag             capabilities = (1 << 1) // JTAG communication
	capSwoUart          capabilities = (1 << 2) // UART Serial Wire Output
	capSwoManchester    capabilities = (1 << 3) // Manchester Serial Wire Output
	capAtomicCommand    capabilities = (1 << 4) // Atomic Commands
	capTestDomainTimer  capabilities = (1 << 5) // Test Domain Timer
	capSwoStreamTracing capabilities = (1 << 6) // SWO Streaming Trace
)

func (caps capabilities) String() string {
	s := []string{}
	if caps&capSwd != 0 {
		s = append(s, "Swd")
	}
	if caps&capJtag != 0 {
		s = append(s, "Jtag")
	}
	if caps&capSwoUart != 0 {
		s = append(s, "SwoUart")
	}
	if caps&capSwoManchester != 0 {
		s = append(s, "SwoManchester")
	}
	if caps&capAtomicCommand != 0 {
		s = append(s, "AtomicCommand")
	}
	if caps&capTestDomainTimer != 0 {
		s = append(s, "TestDomainTimer")
	}
	if caps&capSwoStreamTracing != 0 {
		s = append(s, "SwoStreamTracing")
	}
	return strings.Join(s, ",")
}

// getCapabilities gets information about the Capabilities of the Debug Unit
func (dev *device) getCapabilities() (capabilities, error) {
	rx, err := dev.info(infoCapabilities)
	if err != nil {
		return 0, err
	}
	if len(rx) < 3 || rx[0] != 0 {
		return 0, errors.New("bad info response")
	}
	var info0, info1 uint
	if rx[1] == 1 {
		info0 = uint(rx[2])
	}
	if len(rx) >= 4 && rx[1] == 2 {
		info1 = uint(rx[3])
	}
	return capabilities((info1 << 8) + info0), nil
}

// hasCap returns true if the device has the capability.
func (dev *device) hasCap(x capabilities) bool {
	return dev.caps&x != 0
}

//-----------------------------------------------------------------------------
