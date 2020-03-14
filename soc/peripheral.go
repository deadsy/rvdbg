//-----------------------------------------------------------------------------
/*

Peripherals

*/
//-----------------------------------------------------------------------------

package soc

import "strings"

//-----------------------------------------------------------------------------

// Peripheral is functionally grouped set of registers.
type Peripheral struct {
	Name      string
	Addr      uint
	Size      uint
	Descr     string
	Registers []Register
}

//-----------------------------------------------------------------------------
// Sort by Address for a Peripheral Set

// PeripheralByAddr is used to sort peripherals by address.
type PeripheralByAddr []Peripheral

func (a PeripheralByAddr) Len() int      { return len(a) }
func (a PeripheralByAddr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a PeripheralByAddr) Less(i, j int) bool {
	// Base addresses for peripherals are not always unique. e.g. Nordic chips.
	// So: tie break with the name to give a well-defined sort order.
	if a[i].Addr == a[j].Addr {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	}
	return a[i].Addr < a[j].Addr
}

//-----------------------------------------------------------------------------
