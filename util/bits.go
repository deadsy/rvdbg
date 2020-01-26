//-----------------------------------------------------------------------------
/*

Utilities to Read/Write Bit Fields

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

// BitMask returns a bit mask from the msb to lsb bits.
func BitMask(msb, lsb uint) uint {
	n := msb - lsb + 1
	return ((1 << n) - 1) << lsb
}

// GetBits reads a bit field from a value.
func GetBits(x, msb, lsb uint) uint {
	return (x & BitMask(msb, lsb)) >> lsb
}

// MaskBits masks a bit field within a value.
func MaskBits(x, msb, lsb uint) uint {
	return x & BitMask(msb, lsb)
}

// SetBits writes a bit field within a value.
func SetBits(x, val, msb, lsb uint) uint {
	mask := BitMask(msb, lsb)
	val = (val << lsb) & mask
	return (x & ^mask) | val
}

//-----------------------------------------------------------------------------

// BoolToInt converts a boolean to an int (1 or 0).
func BoolToInt(x bool) int {
	if x {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------------
