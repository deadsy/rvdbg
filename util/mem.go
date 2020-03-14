//-----------------------------------------------------------------------------
/*

Memory related utilities.

*/
//-----------------------------------------------------------------------------

package util

import "fmt"

//-----------------------------------------------------------------------------

const KiB = 1 << 10
const MiB = 1 << 20
const GiB = 1 << 30

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
