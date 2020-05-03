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

func ConvertToUint8(width uint, buf []uint) []uint8 {
	switch width {
	case 64:
		out := make([]uint8, len(buf)*8)
		for i := range buf {
			out[(8*i)+0] = uint8(buf[i] >> 0)
			out[(8*i)+1] = uint8(buf[i] >> 8)
			out[(8*i)+2] = uint8(buf[i] >> 16)
			out[(8*i)+3] = uint8(buf[i] >> 24)
			out[(8*i)+4] = uint8(buf[i] >> 32)
			out[(8*i)+5] = uint8(buf[i] >> 40)
			out[(8*i)+6] = uint8(buf[i] >> 48)
			out[(8*i)+7] = uint8(buf[i] >> 56)
		}
		return out
	case 32:
		out := make([]uint8, len(buf)*4)
		for i := range buf {
			out[(4*i)+0] = uint8(buf[i] >> 0)
			out[(4*i)+1] = uint8(buf[i] >> 8)
			out[(4*i)+2] = uint8(buf[i] >> 16)
			out[(4*i)+3] = uint8(buf[i] >> 24)
		}
		return out
	case 16:
		out := make([]uint8, len(buf)*2)
		for i := range buf {
			out[(2*i)+0] = uint8(buf[i] >> 0)
			out[(2*i)+1] = uint8(buf[i] >> 8)
		}
		return out
	case 8:
		out := make([]uint8, len(buf))
		for i := range buf {
			out[i] = uint8(buf[i])
		}
		return out
	}
	panic(fmt.Sprintf("%d-bit to uint8 conversion not supported", width))
	return nil
}

//-----------------------------------------------------------------------------
