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

// fmtFunc is formatting function for a uint.
type fmtFunc func(x uint) string

// Enum maps a field value to a display string.
type Enum map[uint]string

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
	return a[i].Msb > a[j].Msb
}

//-----------------------------------------------------------------------------

// Display returns display strings for a bit field.
func (f *Field) Display(val uint) []string {

	x := util.Bits(val, f.Msb, f.Lsb)

	// has the value changed?
	changed := ""

	// name string
	var nameStr string
	if f.Msb == f.Lsb {
		nameStr = fmt.Sprintf("  %s[%d]", f.Name, f.Lsb)
	} else {
		nameStr = fmt.Sprintf("  %s[%d:%d]", f.Name, f.Msb, f.Lsb)
	}

	// value string

	var valName string

	if f.Fmt != nil {
		valName = f.Fmt(x)
	} else if f.Enums != nil {

	}

	var valStr string
	if x < 10 {
		valStr = fmt.Sprintf(": %d %s%s", x, valName, changed)
	} else {
		valStr = fmt.Sprintf(": 0x%x %s%s", x, valName, changed)
	}

	return []string{nameStr, valStr, "", f.Descr}
}

//-----------------------------------------------------------------------------

// DisplayH returns the horizontal display string for the bit fields of a uint value.
func DisplayH(fs []Field, val uint) string {
	s := []string{}
	for _, f := range fs {
		x := util.Bits(val, f.Msb, f.Lsb)
		s = append(s, fmt.Sprintf("%s %s", f.Name, f.Fmt(x)))
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
