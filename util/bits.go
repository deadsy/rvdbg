//-----------------------------------------------------------------------------
/*

Utilities to Read/Write Bit Fields

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

const Mask8 = (1 << 8) - 1
const Mask16 = (1 << 16) - 1
const Mask32 = (1 << 32) - 1
const Mask64 = (1 << 64) - 1

//-----------------------------------------------------------------------------

// BitMask returns a bit mask from the msb to lsb bits.
func BitMask(msb, lsb uint) uint {
	n := msb - lsb + 1
	return ((1 << n) - 1) << lsb
}

// Bits reads a bit field from a value.
func Bits(x, msb, lsb uint) uint {
	return (x & BitMask(msb, lsb)) >> lsb
}

// Bit reads a bit from a value.
func Bit(x, n uint) uint {
	return (x & BitMask(n, n)) >> n
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

// BoolToUint converts a boolean to an usnigned int (1 or 0).
func BoolToUint(x bool) uint {
	if x {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------------
