//-----------------------------------------------------------------------------
/*

Memory Display

*/
//-----------------------------------------------------------------------------

package mem

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"math"
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
	shift := util.WidthToShift(width)
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
func (mr *memReader) totalReads(size uint) int {
	bytesPerRead := size << mr.shift
	return int((mr.n + bytesPerRead - 1) / bytesPerRead)
}

func (mr *memReader) Read(buf []uint) (int, error) {
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

func (mr *memReader) Close() error {
	return nil
}

//-----------------------------------------------------------------------------
// memory picture

// analyze the buffer and return a character to represent it
func analyze(data []uint8, ofs, n int) rune {
	// are we off the end of the buffer?
	if ofs >= len(data) {
		return ' '
	}
	// trim the length we will check
	if ofs+n > len(data) {
		n = len(data) - ofs
	}
	var c rune
	b0 := data[ofs]
	if b0 == 0 {
		c = '-'
	} else if b0 == 0xff {
		c = '.'
	} else {
		return '$'
	}
	for i := 0; i < n; i++ {
		if data[ofs+i] != b0 {
			return '$'
		}
	}
	return c
}

type memPicture struct {
	ui                          cli.USER // access to user interface
	addr                        uint     // address of memory buffer
	fmtAddr                     string   // format of address string
	cols, rows                  int      // number of cols/rows
	bytesPerSymbol, bytesPerRow int      // bytes per symbol/row
	rowBuffer                   []uint8  // running buffer for row data
}

func (mp *memPicture) rowString() string {
	s := []rune{}
	addrStr := fmt.Sprintf(mp.fmtAddr, mp.addr)
	for ofs := 0; ofs < mp.bytesPerRow; ofs += mp.bytesPerSymbol {
		s = append(s, analyze(mp.rowBuffer, ofs, mp.bytesPerSymbol))
	}
	return fmt.Sprintf("%s %s", addrStr, string(s))
}

// headerString returns the header string for memory picture.
func (mp *memPicture) headerString() string {
	s := []string{}
	s = append(s, "'.' all ones, '-' all zeroes, '$' various")
	s = append(s, fmt.Sprintf("%d (0x%x) bytes per symbol", mp.bytesPerSymbol, mp.bytesPerSymbol))
	s = append(s, fmt.Sprintf("%d (0x%x) bytes per row", mp.bytesPerRow, mp.bytesPerRow))
	s = append(s, fmt.Sprintf("%d cols x %d rows", mp.cols, mp.rows))
	return strings.Join(s, "\n")
}

const colsMax = 70

func newMemPicture(ui cli.USER, addr, n, addrWidth uint) *memPicture {
	// work out how many rows, columns and bytes per symbol we should display
	cols := colsMax + 1
	bytesPerSymbol := 1
	// we try to display a matrix that is roughly square
	for cols > colsMax {
		bytesPerSymbol *= 2
		cols = int(math.Sqrt(float64(n) / float64(bytesPerSymbol)))
	}
	rows := int(math.Ceil(float64(n) / (float64(cols) * float64(bytesPerSymbol))))
	// bytes per row
	bytesPerRow := cols * bytesPerSymbol
	return &memPicture{
		ui:             ui,
		addr:           addr,
		fmtAddr:        util.UintFormat(addrWidth),
		cols:           cols,
		rows:           rows,
		bytesPerSymbol: bytesPerSymbol,
		bytesPerRow:    bytesPerRow,
	}
}

func (mp *memPicture) Write(buf []uint) (int, error) {
	// assume the []uint buf has 32-bit values
	mp.rowBuffer = append(mp.rowBuffer, util.ConvertToUint8(32, buf)...)
	for len(mp.rowBuffer) > mp.bytesPerRow {
		mp.ui.Put(fmt.Sprintf("%s\r\n", mp.rowString()))
		mp.rowBuffer = mp.rowBuffer[mp.bytesPerRow:]
		mp.addr += uint(mp.bytesPerRow)
	}
	return len(buf), nil
}

func (mp *memPicture) Close() error {
	if len(mp.rowBuffer) != 0 {
		mp.ui.Put(fmt.Sprintf("%s\n", mp.rowString()))
	}
	return nil
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
	shift := util.WidthToShift(width)
	// round down address to be width-bit aligned.
	align := uint((1 << shift) - 1)
	addr &= ^align
	return &memDisplay{
		ui:       ui,
		fmtLine:  fmt.Sprintf("%s  %%s  %%s\r\n", util.UintFormat(addrWidth)),
		fmtData:  util.UintFormat(width),
		addrMask: (1 << addrWidth) - 1,
		addr:     addr,
		width:    width,
		shift:    shift,
	}
}

func (md *memDisplay) Write(buf []uint) (int, error) {
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
		data := util.ConvertToUint8(md.width, lineBuf)
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

func (md *memDisplay) Close() error {
	return nil
}

//-----------------------------------------------------------------------------
// MD5 writer

type md5Writer struct {
	h     hash.Hash
	width uint // data has width-bit values
}

func newMd5Writer(width uint) *md5Writer {
	return &md5Writer{
		h:     md5.New(),
		width: width,
	}
}

func (mw *md5Writer) String() string {
	return hex.EncodeToString(mw.h.Sum(nil))
}

func (mw *md5Writer) Write(buf []uint) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	_, err := mw.h.Write(util.ConvertToUint8(mw.width, buf))
	if err != nil {
		return 0, err
	}
	return len(buf), nil
}

func (mw *md5Writer) Close() error {
	return nil
}

//-----------------------------------------------------------------------------
