//-----------------------------------------------------------------------------
/*

Peripherals

*/
//-----------------------------------------------------------------------------

package soc

import (
	"strings"

	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// Peripheral is functionally grouped set of registers.
type Peripheral struct {
	Name      string
	Addr      uint
	Size      uint
	Descr     string
	Registers []Register
}

// PeripheralSet is a set of peripherals.
type PeripheralSet []Peripheral

//-----------------------------------------------------------------------------
// Sort peripherals set by address.

func (a PeripheralSet) Len() int      { return len(a) }
func (a PeripheralSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a PeripheralSet) Less(i, j int) bool {
	// Base addresses for peripherals are not always unique. e.g. Nordic chips.
	// So: tie break with the name to give a well-defined sort order.
	if a[i].Addr == a[j].Addr {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	}
	return a[i].Addr < a[j].Addr
}

//-----------------------------------------------------------------------------

// GetRegister returns the named register if it exists.
func (p *Peripheral) GetRegister(name string) *Register {
	for i := range p.Registers {
		r := &p.Registers[i]
		if r.Name == name {
			return r
		}
	}
	return nil
}

// Display returns a string for the decoded registers of the peripheral.
func (p *Peripheral) Display(drv Driver, r *Register, fields bool) string {
	s := [][]string{}
	if r != nil {
		// decode a single register
		s = append(s, r.Display(drv, fields)...)
	} else {
		// decode all registers
		for i := range p.Registers {
			s = append(s, p.Registers[i].Display(drv, fields)...)
		}
	}
	return cli.TableString(s, []int{0, 0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------
