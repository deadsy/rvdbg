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
	GetAddressSize() uint                         // get address size in bits
	LookupSymbol(name string) (uint, uint, error) // lookup the address of a symbol
	RdMem(width, addr, n uint) ([]uint, error)    // read width-bit memory buffer
	WrMem(width, addr uint, val []uint) error     // write width-bit memory buffer
}

// target provides a method for getting the memory driver.
type target interface {
	GetMemoryDriver() Driver
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
		buf, err := tgt.RdMem(width, addr, bytesPerLine/(width>>3))
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
		data := util.ConvertXY(width, 8, buf)
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
