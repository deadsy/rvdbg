//-----------------------------------------------------------------------------
/*

Memory Menu Items

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"time"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

var helpMemRegion = []cli.Help{
	{"<addr/name> [len]", "memory region"},
	{"  addr", "address (hex), default is 0"},
	{"  name", "region name (string), see \"map\" command"},
	{"  len", "length (hex), defaults to region size or 0x100"},
}

//-----------------------------------------------------------------------------

type copyState struct {
	rd       util.Reader    // read from
	wr       util.Writer    // write to
	size     int            // buffer size
	progress *util.Progress // progress indicator
	idx      int            // progress index
	err      error          // stored error
}

// copyLoop is the looping function for read from, write to copying
func copyLoop(cs *copyState) bool {
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

//-----------------------------------------------------------------------------
// memory display

func display(c *cli.CLI, args []string, width uint) {
	drv := c.User.(target).GetMemoryDriver()
	r, err := RegionArg(drv, args)
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
		return
	}
	// read from memory, write to the display
	cs := &copyState{
		rd:   newMemReader(drv, r.addr, r.size, width),
		wr:   newMemDisplay(c.User, r.addr, drv.GetAddressSize(), width),
		size: 16,
	}
	c.Loop(func() bool { return copyLoop(cs) }, cli.KeycodeCtrlD)
	if cs.err != nil {
		c.User.Put(fmt.Sprintf("%s\n", cs.err))
	}
}

var cmdDisplay8 = cli.Leaf{
	Descr: "memory display 8-bit",
	F: func(c *cli.CLI, args []string) {
		display(c, args, 8)
	},
}

var cmdDisplay16 = cli.Leaf{
	Descr: "memory display 16-bit",
	F: func(c *cli.CLI, args []string) {
		display(c, args, 16)
	},
}

var cmdDisplay32 = cli.Leaf{
	Descr: "memory display 32-bit",
	F: func(c *cli.CLI, args []string) {
		display(c, args, 32)
	},
}

var cmdDisplay64 = cli.Leaf{
	Descr: "memory display 64-bit",
	F: func(c *cli.CLI, args []string) {
		display(c, args, 64)
	},
}

//-----------------------------------------------------------------------------
// memory read

var helpMemRead = []cli.Help{
	{"<addr>", "read value from memory address"},
	{"  addr", "address (hex)"},
}

func cmdRead(c *cli.CLI, args []string, width uint) {
	drv := c.User.(target).GetMemoryDriver()
	_ = drv
}

var cmdRead8 = cli.Leaf{
	Descr: "memory read 8-bit",
	F: func(c *cli.CLI, args []string) {
		cmdRead(c, args, 8)
	},
}

var cmdRead16 = cli.Leaf{
	Descr: "memory read 16-bit",
	F: func(c *cli.CLI, args []string) {
		cmdRead(c, args, 16)
	},
}

var cmdRead32 = cli.Leaf{
	Descr: "memory read 32-bit",
	F: func(c *cli.CLI, args []string) {
		cmdRead(c, args, 32)
	},
}

var cmdRead64 = cli.Leaf{
	Descr: "memory read 64-bit",
	F: func(c *cli.CLI, args []string) {
		cmdRead(c, args, 64)
	},
}

//-----------------------------------------------------------------------------
// memory write

var helpMemWrite = []cli.Help{
	{"<addr> <val>", "write value to memory address"},
	{"  addr", "address (hex)"},
	{"  val", "value (hex)"},
}

func cmdWrite(c *cli.CLI, args []string, width uint) {
	drv := c.User.(target).GetMemoryDriver()
	_ = drv
}

var cmdWrite8 = cli.Leaf{
	Descr: "memory write 8-bit",
	F: func(c *cli.CLI, args []string) {
		cmdWrite(c, args, 8)
	},
}

var cmdWrite16 = cli.Leaf{
	Descr: "memory write 16-bit",
	F: func(c *cli.CLI, args []string) {
		cmdWrite(c, args, 16)
	},
}

var cmdWrite32 = cli.Leaf{
	Descr: "memory write 32-bit",
	F: func(c *cli.CLI, args []string) {
		cmdWrite(c, args, 32)
	},
}

var cmdWrite64 = cli.Leaf{
	Descr: "memory write 64-bit",
	F: func(c *cli.CLI, args []string) {
		cmdWrite(c, args, 64)
	},
}

//-----------------------------------------------------------------------------
// memory to file

var helpMemToFile = []cli.Help{
	{"<filename> <addr/name> [len]", "read from memory, write to file"},
	{"  filename", "filename (string)"},
	{"  addr", "address (hex), default is 0"},
	{"  name", "region name (string), see \"map\" command"},
	{"  len", "length (hex), defaults to region size or 0x100"},
}

// fileRegionArg converts filename and memory region arguments to a (name, addr, n) tuple.
func fileRegionArg(drv Driver, args []string) (string, uint, uint, error) {
	err := cli.CheckArgc(args, []int{1, 2, 3})
	if err != nil {
		return "", 0, 0, err
	}
	// args[0] is the filename
	name := args[0]
	// the remaining arguments define the memory region
	r, err := RegionArg(drv, args[1:])
	if err != nil {
		return "", 0, 0, err
	}
	return name, r.addr, r.size, nil
}

var cmdToFile = cli.Leaf{
	Descr: "read from memory, write to file",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetMemoryDriver()
		// process the arguments
		err := cli.CheckArgc(args, []int{2, 3})
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		name, addr, n, err := fileRegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// round down address to 32-bit byte boundary
		addr &= ^uint(3)
		// round up n to an integral multiple of 4 bytes
		n = (n + 3) & ^uint(3)
		if n == 0 {
			c.User.Put("nothing to read\n")
			return
		}

		// read from memory, write to file
		const readSize = 1024
		const width = 32
		rd := newMemReader(drv, addr, n, width)
		wr, err := newFileWriter(name, width)
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to open %s (%s)\n", name, err))
			return
		}
		cs := &copyState{
			rd:       rd,
			wr:       wr,
			size:     readSize,
			progress: util.NewProgress(c.User, rd.totalReads(readSize)),
		}
		c.User.Put(fmt.Sprintf("writing %s (ctrl-d to abort): ", name))
		cs.progress.Update(0)
		done := c.Loop(func() bool { return copyLoop(cs) }, cli.KeycodeCtrlD)
		cs.progress.Erase()
		// flush and close the output file
		wr.Close()

		// report result
		if !done {
			c.User.Put("abort\n")
			return
		}
		if cs.err != nil {
			c.User.Put(fmt.Sprintf("error (%s)\n", cs.err))
			return
		}
		c.User.Put("done\n")
	},
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

var cmdPic = cli.Leaf{
	Descr: "display a pictorial summary of memory",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetMemoryDriver()
		// get the arguments
		r, err := RegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// round down address to 32-bit byte boundary
		addr := r.addr & ^uint(3)
		// round up n to an integral multiple of 4 bytes
		n := (r.size + 3) & ^uint(3)
		// work out how many rows, columns and bytes per symbol we should display
		colsMax := 70
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
		// read the memory
		if n > 16*util.KiB {
			c.User.Put("reading memory ...\n")
		}
		data32, err := drv.RdMem(32, addr, n>>2)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		data8 := util.ConvertToUint8(32, data32)
		// display the summary
		c.User.Put("'.' all ones, '-' all zeroes, '$' various\n")
		c.User.Put(fmt.Sprintf("%d (0x%x) bytes per symbol\n", bytesPerSymbol, bytesPerSymbol))
		c.User.Put(fmt.Sprintf("%d (0x%x) bytes per row\n", bytesPerRow, bytesPerRow))
		c.User.Put(fmt.Sprintf("%d cols x %d rows\n", cols, rows))
		// display the matrix
		addrFmt := fmt.Sprintf("0x%s: ", util.UintFormat(drv.GetAddressSize()))
		var ofs int
		for y := 0; y < rows; y++ {
			s := []rune{}
			addrStr := fmt.Sprintf(addrFmt, addr+uint(ofs))
			for x := 0; x < cols; x++ {
				s = append(s, analyze(data8, ofs, bytesPerSymbol))
				ofs += bytesPerSymbol
			}
			c.User.Put(fmt.Sprintf("%s%s\n", addrStr, string(s)))
		}
	},
}

//-----------------------------------------------------------------------------
// memory region checksum

var cmdCheckSum = cli.Leaf{
	Descr: "calcuate md5 checksum of memory region",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetMemoryDriver()

		// get the arguments
		r, err := RegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// round down address to 32-bit byte boundary
		addr := r.addr & ^uint(3)
		// round up n to an integral multiple of 4 bytes
		n := (r.size + 3) & ^uint(3)

		// read from memory, write to checksum
		const readSize = 1024
		const width = 32
		rd := newMemReader(drv, addr, n, width)
		wr := newMd5Writer(width)
		cs := &copyState{
			rd:       rd,
			wr:       wr,
			size:     readSize,
			progress: util.NewProgress(c.User, rd.totalReads(readSize)),
		}
		c.User.Put(fmt.Sprintf("reading memory (ctrl-d to abort): "))
		cs.progress.Update(0)
		done := c.Loop(func() bool { return copyLoop(cs) }, cli.KeycodeCtrlD)
		cs.progress.Erase()

		// report result
		if !done {
			c.User.Put("abort\n")
			return
		}
		if cs.err != nil {
			c.User.Put(fmt.Sprintf("error (%s)\n", cs.err))
			return
		}
		c.User.Put("done\n")
		c.User.Put(fmt.Sprintf("%s\n", wr))
	},
}

//-----------------------------------------------------------------------------
// memory test

// randBuf returns a random buffer of masked values.
func randBuf(n, mask uint) []uint {
	buf := make([]uint, n)
	for i := range buf {
		buf[i] = uint(rand.Uint64()) & mask
	}
	return buf
}

// cmpBuf returns true if the buffers are the same.
func cmpBuf(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func cmdTest(c *cli.CLI, args []string, width uint) {
	drv := c.User.(target).GetMemoryDriver()
	// get the arguments
	r, err := RegionArg(drv, args)
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
		return
	}
	// round down address to 32-bit byte boundary
	addr := r.addr & ^uint(3)
	// round up n to an integral multiple of 4 bytes
	n := (r.size + 3) & ^uint(3)
	// build a random buffer of width-bit words
	nx := n / (width >> 3)
	mask := uint((1 << width) - 1)
	wrbuf := randBuf(nx, mask)
	// TODO halt the cpu
	// write memory
	start := time.Now()
	err = drv.WrMem(width, addr, wrbuf)
	if err != nil {
		c.User.Put(fmt.Sprintf("write error: %s\n", err))
		return
	}
	delta := time.Now().Sub(start)
	c.User.Put(fmt.Sprintf("write %.2f KiB/sec\n", float64(n)/(1024.0*delta.Seconds())))
	// read memory
	start = time.Now()
	rdbuf, err := drv.RdMem(width, addr, nx)
	if err != nil {
		c.User.Put(fmt.Sprintf("read error: %s\n", err))
		return
	}
	delta = time.Now().Sub(start)
	c.User.Put(fmt.Sprintf("read %.2f KiB/sec\n", float64(n)/(1024.0*delta.Seconds())))
	c.User.Put(fmt.Sprintf("read %s write\n", []string{"!=", "=="}[util.BoolToInt(cmpBuf(rdbuf, wrbuf))]))
}

var cmdTest8 = cli.Leaf{
	Descr: "memory test 8-bit write/read",
	F: func(c *cli.CLI, args []string) {
		cmdTest(c, args, 8)
	},
}

var cmdTest16 = cli.Leaf{
	Descr: "memory test 16-bit write/read",
	F: func(c *cli.CLI, args []string) {
		cmdTest(c, args, 16)
	},
}

var cmdTest32 = cli.Leaf{
	Descr: "memory test 32-bit write/read",
	F: func(c *cli.CLI, args []string) {
		cmdTest(c, args, 32)
	},
}

var cmdTest64 = cli.Leaf{
	Descr: "memory test 64-bit write/read",
	F: func(c *cli.CLI, args []string) {
		cmdTest(c, args, 64)
	},
}

//-----------------------------------------------------------------------------

// Menu memory submenu items
var Menu = cli.Menu{
	{"db", cmdDisplay8, helpMemRegion},
	{"dh", cmdDisplay16, helpMemRegion},
	{"dw", cmdDisplay32, helpMemRegion},
	{"dd", cmdDisplay64, helpMemRegion},
	{"rb", cmdRead8, helpMemRead},
	{"rh", cmdRead16, helpMemRead},
	{"rw", cmdRead32, helpMemRead},
	{"rd", cmdRead64, helpMemRead},
	{"tb", cmdTest8, helpMemRegion},
	{"th", cmdTest16, helpMemRegion},
	{"tw", cmdTest32, helpMemRegion},
	{"td", cmdTest64, helpMemRegion},
	{"wb", cmdWrite8, helpMemWrite},
	{"wh", cmdWrite16, helpMemWrite},
	{"ww", cmdWrite32, helpMemWrite},
	{"wd", cmdWrite64, helpMemWrite},
	{">file", cmdToFile, helpMemToFile},
	{"md5", cmdCheckSum, helpMemRegion},
	{"pic", cmdPic, helpMemRegion},
}

//-----------------------------------------------------------------------------
