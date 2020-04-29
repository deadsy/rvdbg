//-----------------------------------------------------------------------------
/*

Common target functions.

*/
//-----------------------------------------------------------------------------

package target

import (
	"sort"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/itf"
)

//-----------------------------------------------------------------------------

type Target interface {
	GetPrompt() string
	GetMenuRoot() []cli.MenuItem
	Shutdown()
	Put(s string)
}

// Info provides general target information.
type Info struct {
	Name     string   // short name for target (command line)
	Descr    string   // description of target
	DbgType  itf.Type // default debugger type
	DbgMode  itf.Mode // debugger interface mode
	DbgSpeed int      // debugger clock speed
	Volts    int      // target voltage
}

//-----------------------------------------------------------------------------

var targetDb = map[string]*Info{}

// Add a target to the database.
func Add(info *Info) {
	targetDb[info.Name] = info
}

// Lookup target information by name.
func Lookup(name string) *Info {
	return targetDb[name]
}

// List all the supported targets
func List() string {
	// sort the target names
	names := make([]string, 0, len(targetDb))
	for k := range targetDb {
		names = append(names, k)
	}
	sort.Strings(names)
	// create the ordered target list
	s := make([][]string, 0, len(names))
	for _, k := range names {
		s = append(s, []string{"", k, targetDb[k].Descr})
	}
	return cli.TableString(s, []int{8, 12, 0}, 1)
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
