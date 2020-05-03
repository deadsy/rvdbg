//-----------------------------------------------------------------------------
/*

Memory Display

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"io"
	"strings"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the memory driver api.
type Driver interface {
	GetAddressSize() uint                      // get address width in bits
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
// memory reader

type memReader struct {
	drv   Driver // memory driver
	addr  uint   // address to read from
	n     uint   // number of bytes to read
	width uint   // data has width-bit values
	shift int    // shift for width-bits
	err   error  // error state
}

func newMemReader(drv Driver, addr, n, width uint) *memReader {
	shift := map[uint]int{8: 0, 16: 1, 32: 2, 64: 3}[width]
	// round down address, round up n to be width-bit aligned.
	align := uint((1 << shift) - 1)
	addr &= ^align
	n = (n + align) & ^align
	return &memReader{
		drv:   drv,
		addr:  addr,
		n:     n,
		width: width,
		shift: shift,
	}
}

// totalReads returns the total number of calls to Read() required.
func (mr *memReader) totalReads(n uint) int {
	bytesPerRead := n << mr.shift
	return int((mr.n + bytesPerRead - 1) / bytesPerRead)
}

func (mr *memReader) Read(buf []uint) (n int, err error) {
	if len(buf) == 0 || mr.err != nil {
		return 0, mr.err
	}
	// read from memory
	nread := min(mr.n, uint(len(buf)<<mr.shift))
	mbuf, err := mr.drv.RdMem(mr.width, mr.addr, nread>>mr.shift)
	if err != nil {
		mr.err = err
		return 0, err
	}
	// copy the buffer
	for i := range mbuf {
		buf[i] = mbuf[i]
	}
	mr.addr += nread
	mr.n -= nread
	// return
	if mr.n == 0 {
		mr.err = io.EOF
	}
	return int(nread) >> mr.shift, mr.err
}

//-----------------------------------------------------------------------------
// memory display

const bytesPerLine = 16 // must be a power of 2

type memDisplay struct {
	ui       cli.USER // access to user interface
	fmtLine  string   // format for the line
	fmtData  string   // format for data item
	addrMask uint     // address mask
	addr     uint     // address of memory buffer
	width    uint     // data has width-bit values
	shift    int      // shift for width-bits
}

func newMemDisplay(ui cli.USER, addr, addrWidth, width uint) *memDisplay {
	shift := map[uint]int{8: 0, 16: 1, 32: 2, 64: 3}[width]
	// round down address to be width-bit aligned.
	align := uint((1 << shift) - 1)
	addr &= ^align
	return &memDisplay{
		ui:       ui,
		fmtLine:  fmt.Sprintf("%%0%dx  %%s  %%s\n", [2]int{16, 8}[util.BoolToInt(addrWidth == 32)]),
		fmtData:  fmt.Sprintf("%%0%dx", width>>2),
		addrMask: (1 << addrWidth) - 1,
		addr:     addr,
		width:    width,
		shift:    shift,
	}
}

func (md *memDisplay) Write(buf []uint) (n int, err error) {
	if (len(buf)<<md.shift)&(bytesPerLine-1) != 0 {
		return 0, fmt.Errorf("write buffer must be a multiple of %d bytes", bytesPerLine)
	}
	if len(buf) == 0 {
		return 0, nil
	}
	// print each line
	perLine := bytesPerLine >> md.shift
	for i := 0; i < len(buf); i += perLine {
		lineBuf := buf[i : i+perLine]
		// create the data string
		xStr := make([]string, perLine)
		for j := range xStr {
			xStr[j] = fmt.Sprintf(md.fmtData, lineBuf[j])
		}
		dataStr := strings.Join(xStr, " ")
		// create the ascii string
		data := util.ConvertXY(md.width, 8, lineBuf)
		var ascii [bytesPerLine]rune
		for j, val := range data {
			if val >= 32 && val <= 126 {
				ascii[j] = rune(val)
			} else {
				ascii[j] = '.'
			}
		}
		asciiStr := string(ascii[:])
		// output
		md.ui.Put(fmt.Sprintf(md.fmtLine, md.addr, dataStr, asciiStr))
		// next address
		md.addr += bytesPerLine
		md.addr &= md.addrMask
	}
	return len(buf), nil
}

//-----------------------------------------------------------------------------
