//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11

CLI Functions

*/
//-----------------------------------------------------------------------------

package rv11

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
		dbg := c.User.(target).GetRiscvDebug().(*Debug)
		c.User.Put(dbg.cache.String())
	},
}

var cmdDbus = cli.Leaf{
	Descr: "display dbus registers",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(target).GetRiscvDebug().(*Debug)
		dump, err := dbg.dbusDump()
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to get dbus registers: %v", err))
		}
		c.User.Put(dump)
	},
}

var cmdInfo = cli.Leaf{
	Descr: "display debug information",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(target).GetRiscvDebug().(*Debug)
		c.User.Put(dbg.String())
	},
}

// Menu debug submenu items
var Menu = cli.Menu{
	{"cache", cmdCache},
	{"dbus", cmdDbus},
	{"info", cmdInfo},
}

//-----------------------------------------------------------------------------
