//-----------------------------------------------------------------------------
/*

Utilities to Read/Write Bit Fields

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

// common bit mask values
const (
	Mask8   = (1 << 8) - 1
	Mask16  = (1 << 16) - 1
	Mask32  = (1 << 32) - 1
	Mask34  = (1 << 34) - 1
	Mask64  = (1 << 64) - 1
	Upper32 = Mask32 << 32
)

//-----------------------------------------------------------------------------

// Mask returns a bit mask from the msb to lsb bits.
func Mask(msb, lsb uint) uint {
	n := msb - lsb + 1
	return ((1 << n) - 1) << lsb
}

// Bits reads a bit field from a value.
func Bits(x, msb, lsb uint) uint {
	return (x & Mask(msb, lsb)) >> lsb
}

// Bit reads a single bit from a value.
func Bit(x, n uint) uint {
	return (x >> n) & 1
}

//-----------------------------------------------------------------------------

// BoolToInt converts a boolean to an int (1 or 0).
func BoolToInt(x bool) int {
	if x {
		return 1
	}
	return 0
}

// BoolToUint converts a boolean to an unsigned int (1 or 0).
func BoolToUint(x bool) uint {
	if x {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------------
