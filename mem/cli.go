//-----------------------------------------------------------------------------
/*

Memory Menu Items

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
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

var cmdPic = cli.Leaf{
	Descr: "display a pictorial summary of memory",
	F: func(c *cli.CLI, args []string) {
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
