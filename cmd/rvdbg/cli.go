//-----------------------------------------------------------------------------
/*

RISC-V Debugger CLI

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------
// cli related leaf functions

var cmdHelp = cli.Leaf{
	Descr: "general help",
	F: func(c *cli.CLI, args []string) {
		c.GeneralHelp()
	},
}

var cmdHistory = cli.Leaf{
	Descr: "command history",
	F: func(c *cli.CLI, args []string) {
		c.SetLine(c.DisplayHistory(args))
	},
}

var cmdExit = cli.Leaf{
	Descr: "exit application",
	F: func(c *cli.CLI, args []string) {
		c.Exit()
	},
}

//-----------------------------------------------------------------------------

var cmdJtagChain = cli.Leaf{
	Descr: "display jtag chain state",
	F: func(c *cli.CLI, args []string) {
		jtagChain := c.User.(*debugApp).jtagChain
		c.User.Put(fmt.Sprintf("%s\n", jtagChain))
	},
}

var cmdJtagDriver = cli.Leaf{
	Descr: "display jtag driver state",
	F: func(c *cli.CLI, args []string) {
		jtagDriver := c.User.(*debugApp).jtagDriver
		c.User.Put(fmt.Sprintf("%s\n", jtagDriver))
	},
}

// jtagMenu submenu items
var jtagMenu = cli.Menu{
	{"chain", cmdJtagChain},
	{"driver", cmdJtagDriver},
}

//-----------------------------------------------------------------------------

// root menu
var menuRoot = cli.Menu{
	{"exit", cmdExit},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"jtag", jtagMenu, "jtag functions"},
}

//-----------------------------------------------------------------------------
