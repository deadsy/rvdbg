//-----------------------------------------------------------------------------
/*

JTAG Device Functions

*/
//-----------------------------------------------------------------------------

package jtag

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvdbg/bitstr"
)

//-----------------------------------------------------------------------------

// Device stores the state for a single device on a JTAG chain.
type Device struct {
	idx         int    // index of device on the JTAG chain
	chain       *Chain // pointer back to full JTAG chain
	drv         Driver // jtag driver
	name        string // device name
	idcode      IDCode // ID code for the device
	irlen       int    // IR length for this device
	irlenBefore int    // IR bits before this device
	irlenAfter  int    // IR bits after this device
	devsBefore  int    // number of devices before this one in the chain
	devsAfter   int    // number of devices after this one in the chain
}

// NewDevice returns the interface object for a single device on a JTAG chain.
func (ch *Chain) NewDevice(idx int) *Device {
	return &Device{
		idx:         idx,
		chain:       ch,
		drv:         ch.drv,
		name:        ch.info[idx].Name,
		idcode:      ch.info[idx].ID,
		irlen:       ch.info[idx].IRLength,
		irlenBefore: ch.info.irLengthBefore(idx),
		irlenAfter:  ch.info.irLengthAfter(idx),
		devsBefore:  idx,
		devsAfter:   len(ch.info) - idx - 1,
	}
}

func (dev *Device) String() string {
	return fmt.Sprintf("idx %d %s irlen %d %s", dev.idx, dev.name, dev.irlen, dev.idcode)
}

// wrIR writes to IR for a device
func (dev *Device) wrIR(wr *bitstr.BitString) error {
	// place other devices into bypass mode (IR = all 1's)
	tdi := bitstr.Ones(dev.irlenBefore).Tail(wr).Tail1(dev.irlenAfter)
	_, err := dev.drv.ScanIR(tdi, false)
	return err
}

// rwIR reads and writes IR for a device.
func (dev *Device) rwIR(wr *bitstr.BitString) (*bitstr.BitString, error) {
	tdi := bitstr.Ones(dev.irlenBefore).Tail(wr).Tail1(dev.irlenAfter)
	tdo, err := dev.drv.ScanIR(tdi, true)
	if err != nil {
		return nil, err
	}
	// strip the IR bits from the other devices
	tdo.DropHead(dev.irlenBefore).DropTail(dev.irlenAfter)
	return tdo, nil
}

// wrDR writes to DR for a device
func (dev *Device) wrDR(wr *bitstr.BitString) error {
	// other devices are assumed to be in bypass mode (DR length = 1)
	tdi := bitstr.Ones(dev.devsBefore).Tail(wr).Tail1(dev.devsAfter)
	_, err := dev.drv.ScanDR(tdi, false)
	return err
}

// rwIR reads and writes DR for a device.
func (dev *Device) rwDR(wr *bitstr.BitString) (*bitstr.BitString, error) {
	tdi := bitstr.Ones(dev.devsBefore).Tail(wr).Tail1(dev.devsAfter)
	tdo, err := dev.drv.ScanDR(tdi, true)
	if err != nil {
		return nil, err
	}
	// strip the DR bits from the bypassed devices
	tdo.DropHead(dev.devsBefore).DropTail(dev.devsAfter)
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

// Survey returns a string with all IR values and corresponding DR lengths.
func (dev *Device) Survey() string {
	s := []string{}
	for ir := 0; ir < (1 << dev.irlen); ir++ {
		err := dev.wrIR(bitstr.FromUint(uint(ir), dev.irlen))
		if err != nil {
			s = append(s, fmt.Sprintf("ir %d can't write ir", ir))
			continue
		}
		n, err := dev.chain.drLength()
		if err != nil {
			s = append(s, fmt.Sprintf("ir %d drlen unknown", ir))
			continue
		}
		n = n - dev.devsAfter - dev.devsBefore
		s = append(s, fmt.Sprintf("ir %d drlen %d", ir, n))
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
