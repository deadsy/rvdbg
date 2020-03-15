//-----------------------------------------------------------------------------
/*

Peripheral Registers

*/
//-----------------------------------------------------------------------------

package soc

import "strings"

//-----------------------------------------------------------------------------

// Register is peripheral register.
type Register struct {
	Name   string
	Offset uint
	Size   uint
	Descr  string
	Fields []Field
}

// RegisterSet is a set of registers.
type RegisterSet []Register

//-----------------------------------------------------------------------------
// Sort registers by offset.

func (a RegisterSet) Len() int      { return len(a) }
func (a RegisterSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a RegisterSet) Less(i, j int) bool {
	// Offsets for registers may not be unique.
	// Tie break with the name to give a well-defined sort order.
	if a[i].Offset == a[j].Offset {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	}
	return a[i].Offset < a[j].Offset
}

//-----------------------------------------------------------------------------
