//-----------------------------------------------------------------------------
/*

File Reader/Writer

*/
//-----------------------------------------------------------------------------

package util

import (
	"bufio"
	"io"
	"os"

	"github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// Reader is a data source interface.
type Reader interface {
	Read(buf []uint) (int, error) // read a []uint buffer from the object
	NumReads(n int) int           // how many reads of n-uint buffers will it take?
	Close() error                 // close the read object
}

// Writer is a data sink interface.
type Writer interface {
	Write(buf []uint) (int, error) // write a []uint buffer to the write object
	Close() error                  // close the write object
}

//-----------------------------------------------------------------------------
// copy from reader to writer

// CopyState records the state of read from, write to copying.
type CopyState struct {
	rd       Reader    // read from
	wr       Writer    // write to
	size     int       // buffer size
	progress *Progress // progress indicator
	idx      int       // progress index
	err      error     // stored error
}

// NewCopyState returns a new copy state object.
func NewCopyState(rd Reader, wr Writer, size int) *CopyState {
	return &CopyState{
		rd:   rd,
		wr:   wr,
		size: size,
	}
}

// CopyLoop is the looping function for read from, write to copying
func (cs *CopyState) CopyLoop() bool {
	buf := make([]uint, cs.size)
	n, err := cs.rd.Read(buf)
	if err != nil && err != io.EOF {
		cs.err = err
		return true
	}
	done := err == io.EOF
	_, err = cs.wr.Write(buf[0:n])
	if err != nil {
		cs.err = err
		return true
	}
	if cs.progress != nil {
		cs.idx++
		cs.progress.Update(cs.idx)
	}
	return done
}

// Start starts the progress indicator for a copy.
func (cs *CopyState) Start(ui cli.USER) {
	cs.progress = NewProgress(ui, cs.rd.NumReads(cs.size))
	cs.progress.Update(0)
}

// Stop stops the progress indicator for a copy.
func (cs *CopyState) Stop() {
	if cs.progress != nil {
		cs.progress.Erase()
	}
}

// GetError returns the recorded error for a copy.
func (cs *CopyState) GetError() error {
	return cs.err
}

//-----------------------------------------------------------------------------
// file writer

type fileWriter struct {
	f     *os.File
	w     *bufio.Writer
	width uint // data has width-bit values
}

// NewFileWriter returns a Writer interface to a file.
func NewFileWriter(name string, width uint) (Writer, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &fileWriter{
		f:     f,
		w:     bufio.NewWriter(f),
		width: width,
	}, nil
}

func (fw *fileWriter) Write(buf []uint) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	_, err := fw.w.Write(ConvertToUint8(fw.width, buf))
	if err != nil {
		return 0, err
	}
	return len(buf), nil
}

func (fw *fileWriter) Close() error {
	fw.w.Flush()
	return fw.f.Close()
}

//-----------------------------------------------------------------------------
// file reader

type fileReader struct {
	name  string // filename
	f     *os.File
	size  uint // size of file in bytes
	width uint // data has width-bit values
	shift int  // shift for width-bits
}

// NewFileReader returns a Reader interface to a file.
func NewFileReader(name string, width uint) (Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	shift := WidthToShift(width)
	return &fileReader{
		name:  name,
		f:     f,
		size:  (uint(info.Size()) >> shift) << shift,
		width: width,
		shift: shift,
	}, nil
}

// TotalReads returns the total number of calls to Read() required.
func (fr *fileReader) NumReads(n int) int {
	bytesPerRead := n << fr.shift
	return (int(fr.size) + bytesPerRead - 1) / bytesPerRead
}

func (fr *fileReader) Read(buf []uint) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	fbuf := make([]byte, len(buf)<<fr.shift)
	n, err := fr.f.Read(fbuf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	if n == 0 {
		return 0, err
	}
	// resize buffers to an integral number of width-bit values
	buf = buf[0 : n>>fr.shift]
	fbuf = fbuf[0 : (n>>fr.shift)<<fr.shift]
	// convert the file buffer
	ConvertFromUint8(fr.width, fbuf, buf)
	return len(buf), err
}

func (fr *fileReader) Close() error {
	return fr.f.Close()
}

func (fr *fileReader) String() string {
	return fr.name
}

//-----------------------------------------------------------------------------
