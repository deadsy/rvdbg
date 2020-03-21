//-----------------------------------------------------------------------------
/*

SoC Device

*/
//-----------------------------------------------------------------------------

package soc

import (
	"sort"
)

//-----------------------------------------------------------------------------

// CPU provides high-level CPU information.
type CPU struct {
}

// Device is the top-level device description.
type Device struct {
	Vendor      string
	Name        string
	Descr       string
	Version     string
	CPU         *CPU
	Interrupts  []Interrupt
	Peripherals []Peripheral
}

//-----------------------------------------------------------------------------

// GetPeripheral returns the named peripheral if it exists.
func (dev *Device) GetPeripheral(name string) *Peripheral {
	if dev == nil {
		return nil
	}
	for i := range dev.Peripherals {
		p := &dev.Peripherals[i]
		if p.Name == name {
			return p
		}
	}
	return nil
}

// AddePeripheral adds periphals to the device.
func (dev *Device) AddPeripheral(p []Peripheral) {
	dev.Peripherals = append(dev.Peripherals, p...)
}

//-----------------------------------------------------------------------------

// Setup performs port-creation setup work on the device structure.
func (dev *Device) Setup() *Device {

	// sort interrupts
	sort.Sort(InterruptSet(dev.Interrupts))
	// sort peripherals
	sort.Sort(PeripheralSet(dev.Peripherals))
	// sort registers
	for _, p := range dev.Peripherals {
		sort.Sort(RegisterSet(p.Registers))
		// sort fields
		for _, r := range p.Registers {
			sort.Sort(FieldSet(r.Fields))
		}
	}

	// setup register parents
	for i := range dev.Peripherals {
		p := &dev.Peripherals[i]
		for j := range p.Registers {
			r := &p.Registers[j]
			r.parent = p
		}
	}

	return dev
}

//-----------------------------------------------------------------------------
