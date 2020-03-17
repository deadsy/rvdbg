//-----------------------------------------------------------------------------
/*

Utilities to Convert Slice Types

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

// Convert32to16 converts a 32-bit slice to a 16-bit slice.
func Convert32to16(x []uint32) []uint16 {
	y := make([]uint16, len(x))
	for i := range x {
		y[i] = uint16(x[i])
	}
	return y
}

// Convert16to32 converts a 16-bit slice to a 32-bit slice.
func Convert16to32(x []uint16) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

// Convert32to8 converts a 32-bit slice to an 8-bit slice.
func Convert32to8(x []uint32) []uint8 {
	y := make([]uint8, len(x))
	for i := range x {
		y[i] = uint8(x[i])
	}
	return y
}

// Convert32to8Little converts a 32-bit slice to a little endian 8-bit slice.
func Convert32to8Little(x []uint32) []uint8 {
	y := make([]uint8, len(x)*4)
	for i := range x {
		j := i * 4
		y[j+0] = uint8(x[i])
		y[j+1] = uint8(x[i] >> 8)
		y[j+2] = uint8(x[i] >> 16)
		y[j+3] = uint8(x[i] >> 24)
	}
	return y
}

// Convert8to32 converts an 8-bit slice to a 32-bit slice.
func Convert8to32(x []uint8) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

// Convert32to64 converts an 32-bit slice to a 64-bit slice.
func Convert32to64(x []uint32) []uint64 {
	if len(x)&1 != 0 {
		panic("len(x) must be a multiple of 2")
	}
	y := make([]uint64, len(x)>>1)
	i := 0
	for j := range y {
		y[j] = uint64(x[i+0]) | uint64(x[i+1]<<32)
		i += 2
	}
	return y
}

// Convert8toUint converts an 8-bit slice to a uint slice.
func Convert8toUint(x []uint8) []uint {
	y := make([]uint, len(x))
	for i := range x {
		y[i] = uint(x[i])
	}
	return y
}

// Convert16toUint converts a 16-bit slice to a uint slice.
func Convert16toUint(x []uint16) []uint {
	y := make([]uint, len(x))
	for i := range x {
		y[i] = uint(x[i])
	}
	return y
}

// Convert32toUint converts a 32-bit slice to a uint slice.
func Convert32toUint(x []uint32) []uint {
	y := make([]uint, len(x))
	for i := range x {
		y[i] = uint(x[i])
	}
	return y
}

// Convert64toUint converts a 64-bit slice to a uint slice.
func Convert64toUint(x []uint64) []uint {
	y := make([]uint, len(x))
	for i := range x {
		y[i] = uint(x[i])
	}
	return y
}

// ConvertUintto8 converts an uint slice to an 8-bit slice.
func ConvertUintto8(x []uint) []uint8 {
	y := make([]uint8, len(x))
	for i := range x {
		y[i] = uint8(x[i])
	}
	return y
}

// ConvertUintto16 converts an uint slice to a 16-bit slice.
func ConvertUintto16(x []uint) []uint16 {
	y := make([]uint16, len(x))
	for i := range x {
		y[i] = uint16(x[i])
	}
	return y
}

// ConvertUintto32 converts an uint slice to a 32-bit slice.
func ConvertUintto32(x []uint) []uint32 {
	y := make([]uint32, len(x))
	for i := range x {
		y[i] = uint32(x[i])
	}
	return y
}

// ConvertUintto64 converts an uint slice to a 64-bit slice.
func ConvertUintto64(x []uint) []uint64 {
	y := make([]uint64, len(x))
	for i := range x {
		y[i] = uint64(x[i])
	}
	return y
}

//-----------------------------------------------------------------------------
