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
)

//-----------------------------------------------------------------------------

// Reader is a data source interface.
type Reader interface {
	Read(buf []uint) (int, error)
	Close() error
}

// Writer is a data sink interface.
type Writer interface {
	Write(buf []uint) (int, error)
	Close() error
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
	f     *os.File
	size  int64 // size of file in bytes
	width uint  // data has width-bit values
	shift int   // shift for width-bits
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
		f:     f,
		size:  (info.Size() >> shift) << shift,
		width: width,
		shift: shift,
	}, nil
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

//-----------------------------------------------------------------------------
