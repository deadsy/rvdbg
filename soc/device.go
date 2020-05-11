//-----------------------------------------------------------------------------
/*

SoC Device

*/
//-----------------------------------------------------------------------------

package soc

import (
	"errors"
	"fmt"
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
func (dev *Device) GetPeripheral(name string) (*Peripheral, error) {
	if dev == nil {
		return nil, errors.New("nil device")
	}
	for i := range dev.Peripherals {
		p := &dev.Peripherals[i]
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("peripheral \"%s\" not found", name)
}

// AddePeripheral adds periphals to the device.
func (dev *Device) AddPeripheral(p []Peripheral) {
	dev.Peripherals = append(dev.Peripherals, p...)
}

// GetPeripheralRegister looks up a register within a peripheral.
func (dev *Device) GetPeripheralRegister(pname, rname string) (*Register, error) {
	p, err := dev.GetPeripheral(pname)
	if err != nil {
		return nil, err
	}
	r, err := p.GetRegister(rname)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// WrPeripheralRegister writes a register within a peripheral.
func (dev *Device) WrPeripheralRegister(drv Driver, pname, rname string, val uint) error {
	r, err := dev.GetPeripheralRegister(pname, rname)
	if err != nil {
		return err
	}
	return r.Wr(drv, 0, val)
}

// RdPeripheralRegister reads a register within a peripheral.
func (dev *Device) RdPeripheralRegister(drv Driver, pname, rname string) (uint, error) {
	r, err := dev.GetPeripheralRegister(pname, rname)
	if err != nil {
		return 0, err
	}
	return r.Rd(drv, 0)
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
