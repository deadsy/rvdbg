//-----------------------------------------------------------------------------
/*

Memory Menu Items

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
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
// memory display

func display(c *cli.CLI, args []string, width uint) {
	drv := c.User.(target).GetMemoryDriver()
	r, err := RegionArg(drv, args)
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
		return
	}
	// read from memory, write to the display
	cs := util.NewCopyState(
		newMemReader(drv, r.Addr, r.Size, width),
		newMemDisplay(c.User, r.Addr, drv.GetAddressSize(), width),
		16)
	c.Loop(func() bool { return cs.CopyLoop() }, cli.KeycodeCtrlD)
	err = cs.GetError()
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
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
		name, region, err := FileRegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		if region.Size == 0 {
			c.User.Put("nothing to read\n")
			return
		}
		// work with 32-bit alignment
		region.Align32()

		// read from memory, write to file
		const width = 32
		rd := newMemReader(drv, region.Addr, region.Size, width)
		wr, err := util.NewFileWriter(name, width)
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to open %s (%s)\n", name, err))
			return
		}
		cs := util.NewCopyState(rd, wr, 1024)
		c.User.Put(fmt.Sprintf("writing %s (ctrl-d to abort): ", name))
		cs.Start(c.User)
		done := c.Loop(func() bool { return cs.CopyLoop() }, cli.KeycodeCtrlD)
		cs.Stop()
		// flush and close the output file
		wr.Close()

		// report result
		if !done {
			c.User.Put("abort\n")
			return
		}

		err = cs.GetError()
		if err != nil {
			c.User.Put(fmt.Sprintf("error (%s)\n", err))
			return
		}
		c.User.Put("done\n")
	},
}

//-----------------------------------------------------------------------------
// memory picture

var cmdPic = cli.Leaf{
	Descr: "display a pictorial summary of memory",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetMemoryDriver()

		// get the arguments
		region, err := RegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		if region.Size == 0 {
			c.User.Put("nothing to read\n")
			return
		}
		// work with 32-bit alignment
		region.Align32()

		// read from memory, write to memory picture
		rd := newMemReader(drv, region.Addr, region.Size, 32)
		wr := newMemPicture(c.User, region.Addr, region.Size, 32)
		cs := util.NewCopyState(rd, wr, 1024)
		c.User.Put(fmt.Sprintf("%s\n", wr.headerString()))
		done := c.Loop(func() bool { return cs.CopyLoop() }, cli.KeycodeCtrlD)
		wr.Close()
		// report result
		if !done {
			c.User.Put("abort\n")
			return
		}
		err = cs.GetError()
		if err != nil {
			c.User.Put(fmt.Sprintf("error (%s)\n", err))
			return
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
		region, err := RegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// work with 32-bit alignment
		region.Align32()

		// read from memory, write to checksum
		const width = 32
		rd := newMemReader(drv, region.Addr, region.Size, width)
		wr := newMd5Writer(width)
		cs := util.NewCopyState(rd, wr, 1024)
		c.User.Put(fmt.Sprintf("reading memory (ctrl-d to abort): "))
		cs.Start(c.User)
		done := c.Loop(func() bool { return cs.CopyLoop() }, cli.KeycodeCtrlD)
		cs.Stop()

		// report result
		if !done {
			c.User.Put("abort\n")
			return
		}
		err = cs.GetError()
		if err != nil {
			c.User.Put(fmt.Sprintf("error (%s)\n", err))
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
	region, err := RegionArg(drv, args)
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
		return
	}
	// work with 32-bit alignment
	region.Align32()

	// build a random buffer of width-bit words
	nx := region.Size / (width >> 3)
	mask := uint((1 << width) - 1)
	wrbuf := randBuf(nx, mask)

	// TODO halt the cpu
	// write memory
	start := time.Now()
	err = drv.WrMem(width, region.Addr, wrbuf)
	if err != nil {
		c.User.Put(fmt.Sprintf("write error: %s\n", err))
		return
	}
	delta := time.Now().Sub(start)
	c.User.Put(fmt.Sprintf("write %.2f KiB/sec\n", float64(region.Size)/(1024.0*delta.Seconds())))
	// read memory
	start = time.Now()
	rdbuf, err := drv.RdMem(width, region.Addr, nx)
	if err != nil {
		c.User.Put(fmt.Sprintf("read error: %s\n", err))
		return
	}
	delta = time.Now().Sub(start)
	c.User.Put(fmt.Sprintf("read %.2f KiB/sec\n", float64(region.Size)/(1024.0*delta.Seconds())))
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
	{"d8", cmdDisplay8, helpMemRegion},
	{"d16", cmdDisplay16, helpMemRegion},
	{"d32", cmdDisplay32, helpMemRegion},
	{"d64", cmdDisplay64, helpMemRegion},
	{"r8", cmdRead8, helpMemRead},
	{"r16", cmdRead16, helpMemRead},
	{"r32", cmdRead32, helpMemRead},
	{"r64", cmdRead64, helpMemRead},
	{"t8", cmdTest8, helpMemRegion},
	{"t16", cmdTest16, helpMemRegion},
	{"t32", cmdTest32, helpMemRegion},
	{"t64", cmdTest64, helpMemRegion},
	{"w8", cmdWrite8, helpMemWrite},
	{"w16", cmdWrite16, helpMemWrite},
	{"w32", cmdWrite32, helpMemWrite},
	{"w64", cmdWrite64, helpMemWrite},
	{">file", cmdToFile, helpMemToFile},
	{"md5", cmdCheckSum, helpMemRegion},
	{"pic", cmdPic, helpMemRegion},
}

//-----------------------------------------------------------------------------
