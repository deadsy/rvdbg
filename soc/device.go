//-----------------------------------------------------------------------------
/*

SoC Device

*/
//-----------------------------------------------------------------------------

package soc

import "sort"

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

// AddPeripheral adds periphals to the device.
func (dev *Device) AddPeripheral(p []Peripheral) {
	dev.Peripherals = append(dev.Peripherals, p...)
}

// RenamePeripheral renames a periphal in the device.
func (dev *Device) RenamePeripheral(oldname, newname string) {
	for i := range dev.Peripherals {
		p := &dev.Peripherals[i]
		if p.Name == oldname {
			p.Name = newname
		}
	}
}

//-----------------------------------------------------------------------------

// Sort the device peripherals, registers and fields.
func (dev *Device) Sort() *Device {
	// interrupts
	sort.Sort(InterruptSet(dev.Interrupts))
	// peripherals
	sort.Sort(PeripheralSet(dev.Peripherals))
	// registers
	for _, p := range dev.Peripherals {
		sort.Sort(RegisterSet(p.Registers))
		// fields
		for _, r := range p.Registers {
			sort.Sort(FieldSet(r.Fields))
		}
	}
	return dev
}

//-----------------------------------------------------------------------------
