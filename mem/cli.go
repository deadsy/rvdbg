//-----------------------------------------------------------------------------
/*

Memory Menu Items

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvemu/util"
)

//-----------------------------------------------------------------------------

// MemArg converts memory arguments to an (address, n) tuple.
func MemArg(defAddr, maxAddr uint, args []string) (uint, uint, error) {
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
	n := uint(0x80) // default size
	if len(args) >= 2 {
		n, err = cli.UintArg(args[1], [2]uint{1, 0x100000000}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return addr, n, nil
}

//-----------------------------------------------------------------------------

func display(c *cli.CLI, args []string, width uint) {
	tgt := c.User.(target).GetMemoryDriver()
	maxAddr := uint((1 << tgt.GetAddressSize()) - 1)
	addr, n, err := util.MemArg(0, maxAddr, args)
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\n", err))
		return
	}
	c.User.Put(fmt.Sprintf("%s\n", displayMem(tgt, addr, n, width)))
}

//-----------------------------------------------------------------------------

var helpMemDisplay = []cli.Help{
	{"<adr> [len]", "address (hex) - default is 0"},
	{"", "length (hex) - default is 0x80"},
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

// DisplayMenu submenu items
var DisplayMenu = cli.Menu{
	{"b", cmdDisplay8, helpMemDisplay},
	{"h", cmdDisplay16, helpMemDisplay},
	{"w", cmdDisplay32, helpMemDisplay},
	{"d", cmdDisplay64, helpMemDisplay},
}

//-----------------------------------------------------------------------------
