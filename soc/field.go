//-----------------------------------------------------------------------------
/*

Bit fields within Registers.

*/
//-----------------------------------------------------------------------------

package soc

//-----------------------------------------------------------------------------

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Enum provides descriptive strings for the enumeration of a bit field.
type Enum map[uint]string

//-----------------------------------------------------------------------------

// fmtFunc is formatting function for a uint.
type fmtFunc func(x uint) string

// Field is a bit field within a register value.
type Field struct {
	Name  string  // name
	Msb   uint    // most significant bit
	Lsb   uint    // least significant bit
	Descr string  // description
	Fmt   fmtFunc // formatting function
	Enums Enum    // enumeration values
}

// FieldSet is a set of fields.
type FieldSet []Field

//-----------------------------------------------------------------------------
// Sort fields by Msb.

func (a FieldSet) Len() int      { return len(a) }
func (a FieldSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a FieldSet) Less(i, j int) bool {
	// MSBs for registers may not be unique.
	// Tie break with the name to give a well-defined sort order.
	if a[i].Msb == a[j].Msb {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	}
	return a[i].Msb < a[j].Msb
}

//-----------------------------------------------------------------------------

// Display returns a display string for a bit field.
func (f *Field) Display(x uint) string {
	val := util.Bits(x, f.Msb, f.Lsb)
	return fmt.Sprintf("%s %s", f.Name, f.Fmt(val))
}

//-----------------------------------------------------------------------------

// Display returns a display string for the bit fields of a uint value.
func (fs FieldSet) Display(x uint) string {
	s := make([]string, len(fs))
	for i := range fs {
		s[i] = (&fs[i]).Display(x)
	}
	return strings.Join(s, " ")
}

//-----------------------------------------------------------------------------
// standard formatting functions

// FmtDec formats a uint as a decimal string.
func FmtDec(x uint) string {
	return fmt.Sprintf("%d", x)
}

// FmtHex formats a uint as a hexadecimal string.
func FmtHex(x uint) string {
	return fmt.Sprintf("%x", x)
}

// FmtHex8 formats a uint as a 2-nybble hexadecimal string.
func FmtHex8(x uint) string {
	return fmt.Sprintf("%02x", x)
}

// FmtHex16 formats a uint as a 4-nybble hexadecimal string.
func FmtHex16(x uint) string {
	return fmt.Sprintf("%04x", x)
}

//-----------------------------------------------------------------------------

// DisplayEnum returns the display string for an enumeration.
func DisplayEnum(x uint, m map[uint]string, unknown string) string {
	s, ok := m[x]
	if !ok {
		s = unknown
	}
	return fmt.Sprintf("%s(%d)", s, x)
}

//-----------------------------------------------------------------------------
