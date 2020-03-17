//-----------------------------------------------------------------------------
/*

Memory Menu Items

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"math"

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

// regionArg converts memory region arguments to an (address, n) tuple.
func regionArg(defAddr, maxAddr uint, args []string) (uint, uint, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return 0, 0, err
	}
	// address
	addr := defAddr
	if len(args) >= 1 {
		addr, err = cli.UintArg(args[0], [2]uint{0, maxAddr}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	// number of bytes
	n := uint(0x100) // default size
	if len(args) >= 2 {
		n, err = cli.UintArg(args[1], [2]uint{1, 0x100000000}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return addr, n, nil
}

//-----------------------------------------------------------------------------
// memory display

func display(c *cli.CLI, args []string, width uint) {
	tgt := c.User.(target).GetMemoryDriver()
	maxAddr := uint((1 << tgt.GetAddressSize()) - 1)
	addr, n, err := regionArg(0, maxAddr, args)
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
		return
	}
	c.User.Put(fmt.Sprintf("%s\n", displayMem(tgt, addr, n, width)))
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

var cmdRead8 = cli.Leaf{
	Descr: "memory read 8-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdRead16 = cli.Leaf{
	Descr: "memory read 16-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdRead32 = cli.Leaf{
	Descr: "memory read 32-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdRead64 = cli.Leaf{
	Descr: "memory read 64-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

//-----------------------------------------------------------------------------
// memory write

var helpMemWrite = []cli.Help{
	{"<addr> <val>", "write value to memory address"},
	{"  addr", "address (hex)"},
	{"  val", "value (hex)"},
}

var cmdWrite8 = cli.Leaf{
	Descr: "memory write 8-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdWrite16 = cli.Leaf{
	Descr: "memory write 16-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdWrite32 = cli.Leaf{
	Descr: "memory write 32-bit",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdWrite64 = cli.Leaf{
	Descr: "memory write 64-bit",
	F: func(c *cli.CLI, args []string) {
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
		tgt := c.User.(target).GetMemoryDriver()
		// get the arguments
		maxAddr := uint((1 << tgt.GetAddressSize()) - 1)
		addr, n, err := regionArg(0, maxAddr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// round down address to 32-bit byte boundary
		addr &= ^uint(3)
		// round up n to an integral multiple of 4 bytes
		n = (n + 3) & ^uint(3)
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
		data32, err := tgt.RdMem(32, addr, n>>2)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		data8 := util.Convert32to8Little(util.ConvertUintto32(data32))
		// display the summary
		c.User.Put("'.' all ones, '-' all zeroes, '$' various\n")
		c.User.Put(fmt.Sprintf("%d (0x%x) bytes per symbol\n", bytesPerSymbol, bytesPerSymbol))
		c.User.Put(fmt.Sprintf("%d (0x%x) bytes per row\n", bytesPerRow, bytesPerRow))
		c.User.Put(fmt.Sprintf("%d cols x %d rows\n", cols, rows))
		// display the matrix
		var ofs int
		for y := 0; y < rows; y++ {
			s := []rune{}
			addrStr := fmt.Sprintf("0x%08x: ", addr+uint(ofs))
			for x := 0; x < cols; x++ {
				s = append(s, analyze(data8, ofs, bytesPerSymbol))
				ofs += bytesPerSymbol
			}
			c.User.Put(fmt.Sprintf("%s%s\n", addrStr, string(s)))
		}
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
	{"wb", cmdWrite8, helpMemWrite},
	{"wh", cmdWrite16, helpMemWrite},
	{"ww", cmdWrite32, helpMemWrite},
	{"wd", cmdWrite64, helpMemWrite},
	{">file", cmdToFile, helpMemToFile},
	{"pic", cmdPic, helpMemRegion},
}

//-----------------------------------------------------------------------------
