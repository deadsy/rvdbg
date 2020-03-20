//-----------------------------------------------------------------------------
/*

Utilities to Convert Slice Types

*/
//-----------------------------------------------------------------------------

package util

import "fmt"

//-----------------------------------------------------------------------------
// 1-1 conversion of []uintX to []uintY

// Cast32to16 one-to-one converts a 32-bit slice to a 16-bit slice.
func Cast32to16(x []uint32) []uint16 {
	y := make([]uint16, len(x))
	for i := range x {
		y[i] = uint16(x[i])
	}
	return y
}

// Cast16to32 one-to-one converts a 16-bit slice to a 32-bit slice.
func Cast16to32(x []uint16) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

// Cast32to8 one-to-one converts a 32-bit slice to an 8-bit slice.
func Cast32to8(x []uint32) []uint8 {
	y := make([]uint8, len(x))
	for i := range x {
		y[i] = uint8(x[i])
	}
	return y
}

// Cast8to32 one-to-one converts an 8-bit slice to a 32-bit slice.
func Cast8to32(x []uint8) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

//-----------------------------------------------------------------------------
// 1-1 conversion of []uintX to []uint

// Cast32toUint one-to-one converts a 32-bit slice to a uint slice.
func Cast32toUint(x []uint32, mask uint32) []uint {
	y := make([]uint, len(x))
	for i := range x {
		y[i] = uint(x[i] & mask)
	}
	return y
}

// Cast64toUint one-to-one converts a 64-bit slice to a uint slice.
func Cast64toUint(x []uint64, mask uint64) []uint {
	y := make([]uint, len(x))
	for i := range x {
		y[i] = uint(x[i] & mask)
	}
	return y
}

//-----------------------------------------------------------------------------
// 1-1 conversion of []uint to []uintX

// CastUintto8 one-to-one converts an uint slice to an 8-bit slice.
func CastUintto8(x []uint) []uint8 {
	y := make([]uint8, len(x))
	for i := range x {
		y[i] = uint8(x[i])
	}
	return y
}

// CastUintto16 one-to-one converts an uint slice to a 16-bit slice.
func CastUintto16(x []uint) []uint16 {
	y := make([]uint16, len(x))
	for i := range x {
		y[i] = uint16(x[i])
	}
	return y
}

// CastUintto32 one-to-one converts an uint slice to a 32-bit slice.
func CastUintto32(x []uint) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

// CastUintto64 one-to-one converts an uint slice to a 64-bit slice.
func CastUintto64(x []uint) []uint64 {
	y := make([]uint64, len(x))
	for i := range x {
		y[i] = uint64(x[i])
	}
	return y
}

//-----------------------------------------------------------------------------

// Convert32to64 converts []uint32 to []uint64 (little endian).
func Convert32to64(x []uint32) []uint64 {
	if len(x)&1 != 0 {
		panic("len(x) must be a multiple of 2")
	}
	y := make([]uint64, len(x)>>1)
	i := 0
	for j := range y {
		y[j] = uint64(x[i+0]<<0) | uint64(x[i+1]<<32)
		i += 2
	}
	return y
}

//-----------------------------------------------------------------------------

// ConvertXY converts x-bit []uint to y-bit []uint (little endian).
func ConvertXY(x, y uint, in []uint) []uint {
	if x == y {
		// no conversion
		return in
	}
	if y == 8 {
		switch x {
		case 64:
			out := make([]uint, len(in)*8)
			for i := range in {
				out[(8*i)+0] = (in[i] >> 0) & 0xff
				out[(8*i)+1] = (in[i] >> 8) & 0xff
				out[(8*i)+2] = (in[i] >> 16) & 0xff
				out[(8*i)+3] = (in[i] >> 24) & 0xff
				out[(8*i)+4] = (in[i] >> 32) & 0xff
				out[(8*i)+5] = (in[i] >> 40) & 0xff
				out[(8*i)+6] = (in[i] >> 48) & 0xff
				out[(8*i)+7] = (in[i] >> 56) & 0xff
			}
			return out
		case 32:
			out := make([]uint, len(in)*4)
			for i := range in {
				out[(4*i)+0] = (in[i] >> 0) & 0xff
				out[(4*i)+1] = (in[i] >> 8) & 0xff
				out[(4*i)+2] = (in[i] >> 16) & 0xff
				out[(4*i)+3] = (in[i] >> 24) & 0xff
			}
			return out
		case 16:
			out := make([]uint, len(in)*2)
			for i := range in {
				out[(2*i)+0] = (in[i] >> 0) & 0xff
				out[(2*i)+1] = (in[i] >> 8) & 0xff
			}
			return out
		}
	}
	panic(fmt.Sprintf("%d to %d conversion not supported", x, y))
	return nil
}

//-----------------------------------------------------------------------------
