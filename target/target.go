//-----------------------------------------------------------------------------
/*

Common target functions.

*/
//-----------------------------------------------------------------------------

package target

import (
	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/itf"
)

//-----------------------------------------------------------------------------

// Info provides general target information.
type Info struct {
	Name  string   // short name for target (command line)
	Descr string   // description of target
	Itf   itf.Type // default debugger type
	Mode  itf.Mode // debugger interface mode
}

//-----------------------------------------------------------------------------
// cli related leaf functions

// CmdHelp provides general CLI help.
var CmdHelp = cli.Leaf{
	Descr: "general help",
	F: func(c *cli.CLI, args []string) {
		c.GeneralHelp()
	},
}

// CmdHistory lists the CLI command history.
var CmdHistory = cli.Leaf{
	Descr: "command history",
	F: func(c *cli.CLI, args []string) {
		c.SetLine(c.DisplayHistory(args))
	},
}

// CmdExit exits the CLI.
var CmdExit = cli.Leaf{
	Descr: "exit application",
	F: func(c *cli.CLI, args []string) {
		c.Exit()
	},
}

//-----------------------------------------------------------------------------
