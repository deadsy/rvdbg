//-----------------------------------------------------------------------------
/*

Memory Display

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the memory driver api.
type Driver interface {
	GetAddressSize() uint                   // get address size in bits
	RdMem8(addr, n uint) ([]uint8, error)   // read 8-bit memory buffer
	RdMem16(addr, n uint) ([]uint16, error) // read 16-bit memory buffer
	RdMem32(addr, n uint) ([]uint32, error) // read 32-bit memory buffer
	RdMem64(addr, n uint) ([]uint64, error) // read 64-bit memory buffer
	WrMem8(addr uint, val []uint8) error    // write 8-bit memory buffer
	WrMem16(addr uint, val []uint16) error  // write 16-bit memory buffer
	WrMem32(addr uint, val []uint32) error  // write 32-bit memory buffer
	WrMem64(addr uint, val []uint64) error  // write 64-bit memory buffer
}

// target provides a method for getting the memory driver.
type target interface {
	GetMemoryDriver() Driver
}

//-----------------------------------------------------------------------------

// rdBuf reads a n x width-bit values from memory.
func rdBuf(tgt Driver, addr, n, width uint) ([]uint, error) {
	switch width {
	case 8:
		x, err := tgt.RdMem8(addr, n)
		return util.Convert8toUint(x), err
	case 16:
		x, err := tgt.RdMem16(addr, n)
		return util.Convert16toUint(x), err
	case 32:
		x, err := tgt.RdMem32(addr, n)
		return util.Convert32toUint(x), err
	case 64:
		x, err := tgt.RdMem64(addr, n)
		return util.Convert64toUint(x), err
	}
	return nil, fmt.Errorf("%d-bit memory reads are not supported", width)
}

//-----------------------------------------------------------------------------

const bytesPerLine = 16 // must be a power of 2

func displayMem(tgt Driver, addr, n, width uint) string {

	s := []string{}

	addrLength := tgt.GetAddressSize()
	addrMask := uint((1 << addrLength) - 1)

	fmtLine := fmt.Sprintf("%%0%dx  %%s  %%s", [2]int{16, 8}[util.BoolToInt(addrLength == 32)])
	fmtData := fmt.Sprintf("%%0%dx", width>>2)

	// round down address to width alignment
	addr &= ^uint((width >> 3) - 1)
	// round up size to an integral multiple of bytesPerLine bytes
	n = (n + bytesPerLine - 1) & ^uint(bytesPerLine-1)

	// read and print the data
	for i := 0; i < int(n/bytesPerLine); i++ {

		// read bytesPerLine bytes
		buf, err := rdBuf(tgt, addr, bytesPerLine/(width>>3), width)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}

		// create the data string
		xStr := make([]string, len(buf))
		for j := range xStr {
			xStr[j] = fmt.Sprintf(fmtData, buf[j])
		}
		dataStr := strings.Join(xStr, " ")

		// create the ascii string
		data, err := tgt.RdMem8(addr, bytesPerLine)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		var ascii [bytesPerLine]rune
		for j, val := range data {
			if val >= 32 && val <= 126 {
				ascii[j] = rune(val)
			} else {
				ascii[j] = '.'
			}
		}
		asciiStr := string(ascii[:])

		s = append(s, fmt.Sprintf(fmtLine, addr, dataStr, asciiStr))
		addr += bytesPerLine
		addr &= addrMask
	}

	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
