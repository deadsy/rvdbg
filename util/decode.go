//-----------------------------------------------------------------------------
/*

Utilities to decode and display bit fields.

*/
//-----------------------------------------------------------------------------

package util

import (
	"fmt"
	"strings"
)

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

// fmtFunc is formatting function for a uint.
type fmtFunc func(x uint) string

// Field is a bit field within a uint value.
type Field struct {
	Name     string
	Msb, Lsb uint
	Fmt      fmtFunc
}

// Display returns a display string for a bit field.
func (f *Field) Display(x uint) string {
	val := GetBits(x, f.Msb, f.Lsb)
	return fmt.Sprintf("%s %s", f.Name, f.Fmt(val))
}

//-----------------------------------------------------------------------------

// FieldSet is a set of field definitions.
type FieldSet []Field

// Display returns a display string for the bit fields of a uint value.
func (fs FieldSet) Display(x uint) string {
	s := make([]string, len(fs))
	for i := range fs {
		s[i] = (&fs[i]).Display(x)
	}
	return strings.Join(s, " ")
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
