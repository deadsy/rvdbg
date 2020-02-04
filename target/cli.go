//-----------------------------------------------------------------------------
/*

Common target CLI functions.

*/
//-----------------------------------------------------------------------------

package target

import (
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
