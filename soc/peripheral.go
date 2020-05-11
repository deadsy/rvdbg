//-----------------------------------------------------------------------------
/*

Peripherals

*/
//-----------------------------------------------------------------------------

package soc

import (
	"errors"
	"fmt"
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
func (p *Peripheral) GetRegister(name string) (*Register, error) {
	if p == nil {
		return nil, errors.New("nil peripheral")
	}
	for i := range p.Registers {
		r := &p.Registers[i]
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("register \"%s\" not found", name)
}

// RemoveRegister ignores a register within the peripheral.
func (p *Peripheral) RemoveRegister(name string) {
	if p == nil {
		return
	}
	for i := range p.Registers {
		r := &p.Registers[i]
		if r.Name == name {
			r.ignore = true
			return
		}
	}
}

// Wr writes a register within a peripheral.
func (p *Peripheral) Wr(drv Driver, rname string, val uint) error {
	r, err := p.GetRegister(rname)
	if err != nil {
		return err
	}
	return r.Wr(drv, 0, val)
}

// Rd reads a register within a peripheral.
func (p *Peripheral) Rd(drv Driver, rname string) (uint, error) {
	r, err := p.GetRegister(rname)
	if err != nil {
		return 0, err
	}
	return r.Rd(drv, 0)
}

// Rmw read/modify/writes a register within a peripheral.
func (p *Peripheral) Rmw(drv Driver, rname string, setbits, clrbits uint) error {
	r, err := p.GetRegister(rname)
	if err != nil {
		return err
	}
	val, err := r.Rd(drv, 0)
	if err != nil {
		return err
	}
	val |= setbits
	val &= ^clrbits
	return r.Wr(drv, 0, val)
}

// Set bits for a register within a peripheral.
func (p *Peripheral) Set(drv Driver, rname string, bits uint) error {
	return p.Rmw(drv, rname, bits, 0)
}

// Clr bits for a register within a peripheral.
func (p *Peripheral) Clr(drv Driver, rname string, bits uint) error {
	return p.Rmw(drv, rname, 0, bits)
}

// Display returns a string for the decoded registers of the peripheral.
func (p *Peripheral) Display(drv Driver, r *Register, fields bool) string {
	s := [][]string{}
	if r != nil {
		// decode a single register
		if !r.ignore && r.regSize(drv) != 0 {
			s = append(s, r.Display(drv, fields)...)
		}
	} else {
		// decode all registers
		for i := range p.Registers {
			r := &p.Registers[i]
			if !r.ignore && r.regSize(drv) != 0 {
				s = append(s, r.Display(drv, fields)...)
			}
		}
	}
	return cli.TableString(s, []int{0, 0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------
