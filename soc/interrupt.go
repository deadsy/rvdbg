//-----------------------------------------------------------------------------
/*

SoC Interrupts

*/
//-----------------------------------------------------------------------------

package soc

import "strings"

//-----------------------------------------------------------------------------

// Interrupt describes an SoC interrupt.
type Interrupt struct {
	Name  string // name
	IRQ   uint   // interrupt request number
	Descr string // description
}

// InterruptSet is a set of interrupts.
type InterruptSet []Interrupt

//-----------------------------------------------------------------------------
// Sort interrupts by IRQ.

func (a InterruptSet) Len() int      { return len(a) }
func (a InterruptSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a InterruptSet) Less(i, j int) bool {
	// Tie break with the name to give a well-defined sort order.
	if a[i].IRQ == a[j].IRQ {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	}
	return a[i].IRQ < a[j].IRQ
}

//-----------------------------------------------------------------------------
