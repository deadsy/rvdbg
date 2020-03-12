//-----------------------------------------------------------------------------
/*

Memory Display

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

const bytesPerLine = 16 // must be a power of 2

//-----------------------------------------------------------------------------

func Display(dbg rv.Debug, addr, n, width uint) string {

	s := []string{}

	hi := dbg.GetCurrentHart()

	fmtLine := fmt.Sprintf("%%0%dx %%s %%s", [2]int{16, 8}[util.BoolToInt(hi.MXLEN == 32)])
	fmtData := fmt.Sprintf("%%0%dx", width>>2)

	// round down address to width alignment
	addr &= ^uint((width >> 3) - 1)
	// round up size to an integral multiple of bytesPerLine bytes
	n = (n + bytesPerLine - 1) & ^uint(bytesPerLine-1)

	// read and print the data
	for i := 0; i < int(n/bytesPerLine); i++ {

		// read bytesPerLine bytes
		buf, err := dbg.RdBuf(addr, bytesPerLine/(width>>3), width)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}

		// create the data string
		xStr := make([]string, len(buf))
		for j := range xStr {
			xStr[j] = fmt.Sprintf(fmtData, buf[j])
		}
		dataStr := strings.Join(xStr[:], " ")

		// create the ascii string
		data, err := dbg.RdMem8(addr, bytesPerLine)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		var ascii [bytesPerLine]string
		for j := range data {
			if data[j] >= 32 && data[j] <= 126 {
				ascii[j] = fmt.Sprintf("%c", data[j])
			} else {
				ascii[j] = "."
			}
		}
		asciiStr := strings.Join(ascii[:], "")

		s = append(s, fmt.Sprintf(fmtLine, addr, dataStr, asciiStr))
		addr += bytesPerLine
		addr &= (1 << hi.MXLEN) - 1
	}

	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
