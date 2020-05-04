//-----------------------------------------------------------------------------
/*

Memory related utilities.

*/
//-----------------------------------------------------------------------------

package util

import "fmt"

//-----------------------------------------------------------------------------

type Reader interface {
	Read(buf []uint) (int, error)
}

type Writer interface {
	Write(buf []uint) (int, error)
}

//-----------------------------------------------------------------------------

// memory size constants
const (
	KiB = 1 << 10
	MiB = 1 << 20
	GiB = 1 << 30
)

// MemSize returns a scaled string for the memory size.
func MemSize(x uint) string {
	if (x >= GiB) && (x&(GiB-1) == 0) {
		return fmt.Sprintf("%dGiB", x/GiB)
	}
	if (x >= MiB) && (x&(MiB-1) == 0) {
		return fmt.Sprintf("%dMiB", x/MiB)
	}
	if (x >= KiB) && (x&(KiB-1) == 0) {
		return fmt.Sprintf("%dKiB", x/KiB)
	}
	return fmt.Sprintf("%dB", x)
}

//-----------------------------------------------------------------------------

// UintFormat returns a format string for the bit size.
func UintFormat(size uint) string {
	// round size up to a multiple of 4 bits
	size = (size + 3) & ^uint(3)
	switch size {
	case 8:
		return "%02x"
	case 12:
		return "%03x"
	case 16:
		return "%04x"
	case 32:
		return "%08x"
	case 64:
		return "%016x"
	}
	return "%x"
}

//-----------------------------------------------------------------------------
