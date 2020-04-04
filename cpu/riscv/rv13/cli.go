//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13

CLI Functions

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// target provides a method for getting the CPU debugger driver.
type target interface {
	GetRiscvDebug() rv.Debug
}

//-----------------------------------------------------------------------------

var cmdCache = cli.Leaf{
	Descr: "display debug ram cache state",
	F: func(c *cli.CLI, args []string) {
		//dbg := c.User.(target).GetRiscvDebug().(*Debug)
		//c.User.Put(dbg.cache.String())
		c.User.Put("TODO\n")
	},
}

var cmdDmi = cli.Leaf{
	Descr: "display dmi registers",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(target).GetRiscvDebug().(*Debug)
		dump, err := dbg.dmiDump()
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to get dmi registers: %v\n", err))
		}
		c.User.Put(fmt.Sprintf("%s\n", dump))
	},
}

var cmdInfo = cli.Leaf{
	Descr: "display debug information",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(target).GetRiscvDebug().(*Debug)
		c.User.Put(fmt.Sprintf("%s\n", dbg.String()))
	},
}

// Menu debug submenu items
var Menu = cli.Menu{
	{"cache", cmdCache},
	{"dmi", cmdDmi},
	{"info", cmdInfo},
}

//-----------------------------------------------------------------------------
