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

// SortedPeripherals returns a sorted list of device peripherals.
func (dev *Device) SortedPeripherals() []Peripheral {
	// Build a list of peripherals in base address order.
	ps := dev.Peripherals
	sort.Sort(PeripheralByAddr(ps))
	return ps
}

//-----------------------------------------------------------------------------
