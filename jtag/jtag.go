//-----------------------------------------------------------------------------
/*

JTAG Chain Management

*/
//-----------------------------------------------------------------------------

package jtag

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/deadsy/rvdbg/bitstr"
)

//-----------------------------------------------------------------------------

const maxDevices = 6
const flushSize = maxDevices * 32
const idcodeLength = 32

//-----------------------------------------------------------------------------

// DeviceInfo describes how the deice will appear on the JTAG chain.
type DeviceInfo struct {
	IRLength int    // length in bits of instruction register
	ID       IDCode // expected id code for the device
	Name     string // name of the device
}

type ChainInfo []DeviceInfo

// irLengthBefore returns the total IR length before device at idx position.
func (ci ChainInfo) irLengthBefore(idx int) int {
	irlen := 0
	for i, d := range ci {
		if i < idx {
			irlen += d.IRLength
		}
	}
	return irlen
}

// irLengthTotal returns the total IR length in the chain information.
func (ci ChainInfo) irLengthTotal() int {
	return ci.irLengthBefore(len(ci))
}

//-----------------------------------------------------------------------------

// Driver is the interface for a JTAG driver.
type Driver interface {
	TestReset(delay time.Duration) error                                   // pulse TRST
	SystemReset(delay time.Duration) error                                 // pulse SRST
	TapReset() error                                                       // reset TAP state machine
	ScanIR(tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) // scan IR
	ScanDR(tdi *bitstr.BitString, needTdo bool) (*bitstr.BitString, error) // scan DR
}

//-----------------------------------------------------------------------------

// Chain stores the state for JTAG chain.
type Chain struct {
	drv   Driver    // jtag driver
	info  ChainInfo // device chain information
	dev   []*Device // devices on the chain
	n     int       // number of devices on the chain
	irlen int       // total IR length
}

// NewChain returns the interface object for a JTAG chain.
func NewChain(drv Driver, info []DeviceInfo) (*Chain, error) {
	ch := &Chain{
		drv:  drv,
		info: info,
	}
	err := ch.scan()
	return ch, err
}

func (ch *Chain) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("jtag chain irlen %d devices %d", ch.irlen, len(ch.dev)))
	for i := range ch.dev {
		s = append(s, ch.dev[i].String())
	}
	return strings.Join(s, "\n")
}

func (ch *Chain) scan() error {
	// reset the TAP state machine for all devices
	err := ch.drv.TapReset()
	if err != nil {
		return err
	}
	// how many devices are on the chain?
	ch.n, err = ch.numDevices()
	if err != nil {
		return err
	}
	// sanity check the number of devices
	if len(ch.info) != ch.n {
		return fmt.Errorf("expecting %d devices, found %d", len(ch.info), ch.n)
	}
	// get the total IR length
	ch.irlen, err = ch.irLength()
	if err != nil {
		return err
	}
	// sanity check the total IR length
	irlen := ch.info.irLengthTotal()
	if irlen != ch.irlen {
		return fmt.Errorf("expecting irlen %d bits, found %d bits", irlen, ch.irlen)
	}
	// sanity check the device id codes
	code, err := ch.readIDCodes()
	if err != nil {
		return err
	}
	for i, d := range ch.info {
		if d.ID != code[i] {
			return fmt.Errorf("expecting idcode 0x%08x at position %d, found 0x%08x", d.ID, i, code[i])
		}
	}
	// build the devices
	ch.dev = make([]*Device, ch.n)
	for i := range ch.dev {
		ch.dev[i] = ch.NewDevice(i)
	}
	// test the IR capture value for all devices
	for _, d := range ch.dev {
		good, err := d.testIRCapture()
		if err != nil {
			return err
		}
		if !good {
			return fmt.Errorf("failed ir capture for idcode 0x%08x at position %d", d.idcode, d.idx)
		}
	}
	return nil
}

// readIDcodes returns a slice of idcodes for the JTAG chain.
func (ch *Chain) readIDCodes() ([]IDCode, error) {
	// a TAP reset leaves the idcodes in the DR chain
	ch.drv.TapReset()
	tdi := bitstr.Ones(ch.n * idcodeLength)
	tdo, err := ch.drv.ScanDR(tdi, true)
	if err != nil {
		return nil, err
	}
	splits := make([]int, ch.n)
	for i := range splits {
		splits[i] = 32
	}
	x, err := tdo.Split(splits), nil
	if err != nil {
		return nil, err
	}
	// re-type uint to IDCode
	code := make([]IDCode, ch.n)
	for i := range code {
		code[i] = IDCode(x[i])
	}
	return code, nil
}

type scanFunc func(tdi *bitstr.BitString) (*bitstr.BitString, error)

// chainLength returns the length of the JTAG chain.
func (ch *Chain) chainLength(scan scanFunc) (int, error) {
	// build a 000...001000...000 flush buffer for tdi
	tdi := bitstr.Zeroes(flushSize).Tail1(1).Tail0(flushSize)
	tdo, err := scan(tdi)
	if err != nil {
		return 0, err
	}
	// the first bits are junk
	tdo.DropHead(flushSize)
	// work out how many bits tdo is behind tdi
	s := tdo.String()
	s = strings.TrimLeft(s, "0")
	if strings.Count(s, "1") != 1 {
		return 0, errors.New("unexpected result from jtag chain, there should be a single 1")
	}
	return len(s) - 1, nil
}

// drLength returns the DR chain length.
// The DR chain length is a function of current IR chain state.
func (ch *Chain) drLength() (int, error) {
	scan := func(tdi *bitstr.BitString) (*bitstr.BitString, error) {
		return ch.drv.ScanDR(tdi, true)
	}
	return ch.chainLength(scan)
}

// irLength returns the IR chain length.
func (ch *Chain) irLength() (int, error) {
	scan := func(tdi *bitstr.BitString) (*bitstr.BitString, error) {
		return ch.drv.ScanIR(tdi, true)
	}
	return ch.chainLength(scan)
}

// numDevices returns the number of JTAG devices in the chain.
func (ch *Chain) numDevices() (int, error) {
	// put every device into bypass mode (IR = all 1's)
	_, err := ch.drv.ScanIR(bitstr.Ones(flushSize), false)
	if err != nil {
		return 0, err
	}
	// now each DR is a single bit
	// the DR chain length is the number of devices
	return ch.drLength()
}

//-----------------------------------------------------------------------------

// Device stores the state for a single device on a JTAG chain.
type Device struct {
	drv         Driver // jtag driver
	name        string // device name
	idx         int    // index of device on the JTAG chain
	idcode      IDCode // ID code for the device
	irlen       int    // IR length for this device
	irlenAfter  int    // IR bits after this device
	irlenBefore int    // IR bits before this device
}

// NewDevice returns the interface object for a single device on a JTAG chain.
func (ch *Chain) NewDevice(idx int) *Device {
	return &Device{}
}

func (dev *Device) String() string {
	return fmt.Sprintf("idx %d %s irlen %d %s", dev.idx, dev.name, dev.irlen, dev.idcode.Decode())
	return ""
}

// rwIR reads and writes IR for a device.
func (dev *Device) rwIR(wr *bitstr.BitString) (*bitstr.BitString, error) {
	tdi := bitstr.Ones(dev.irlenBefore).Tail(wr).Tail1(dev.irlenAfter)
	tdo, err := dev.drv.ScanIR(tdi, true)
	if err != nil {
		return nil, err
	}
	// strip the ir bits from the other devices
	tdo.DropHead(dev.irlenBefore)
	tdo.DropTail(dev.irlenAfter)
	return tdo, nil
}

// testIRCapture tests the IR capture result.
func (dev *Device) testIRCapture() (bool, error) {
	// write all-1s to the IR
	rd, err := dev.rwIR(bitstr.Ones(dev.irlen))
	if err != nil {
		return false, err
	}
	val := rd.Split([]int{dev.irlen})[0]
	// the lowest 2 bits should be "01"
	return val&3 == 1, nil
}

//-----------------------------------------------------------------------------
