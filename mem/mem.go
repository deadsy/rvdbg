//-----------------------------------------------------------------------------
/*

Memory Display

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the memory driver api.
type Driver interface {
	GetAddressSize() uint                      // get address size in bits
	GetDefaultRegion() *Region                 // get a default region
	LookupSymbol(name string) *Region          // lookup a symbol
	RdMem(width, addr, n uint) ([]uint, error) // read width-bit memory buffer
	WrMem(width, addr uint, val []uint) error  // write width-bit memory buffer
}

// target provides a method for getting the memory driver.
type target interface {
	GetMemoryDriver() Driver
}

//-----------------------------------------------------------------------------
// memory reader, implements io.Reader for memory with 32-bit reads.

type memReader struct {
	drv  Driver // memory driver
	addr uint   // address to read from
	size uint   // size of memory region
	n    uint   // bytes remaining to read
	err  error  // error state
}

func newMemReader(drv Driver, addr, n uint) *memReader {
	// round down address to 32-bit byte boundary
	addr &= ^uint(3)
	// round up n to an integral multiple of 4 bytes
	n = (n + 3) & ^uint(3)
	return &memReader{
		drv:  drv,
		addr: addr,
		size: n,
		n:    n,
	}
}

// totalReads returns the total number of calls to Read() required.
func (mr *memReader) totalReads(size uint) int {
	return int((mr.size + size - 1) / size)
}

func (mr *memReader) Read(p []byte) (n int, err error) {
	if len(p)&3 != 0 {
		return 0, errors.New("length of read buffer must be a multiple of 4 bytes")
	}
	if mr.err != nil {
		return 0, mr.err
	}
	// read from memory
	nread := min(mr.n, uint(len(p)))
	buf, err := mr.drv.RdMem(32, mr.addr, nread>>2)
	if err != nil {
		mr.err = err
		return 0, err
	}
	mr.addr += nread
	mr.n -= nread
	// copy the buffer
	i := 0
	for _, x := range buf {
		p[(4*i)+0] = byte(x >> 0)
		p[(4*i)+1] = byte(x >> 8)
		p[(4*i)+2] = byte(x >> 16)
		p[(4*i)+3] = byte(x >> 24)
		i += 4
	}
	// return
	if mr.n == 0 {
		mr.err = io.EOF
		return int(nread), io.EOF
	}
	return int(nread), nil
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
