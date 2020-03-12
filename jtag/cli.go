//-----------------------------------------------------------------------------
/*

JTAG Menu Items

*/
//-----------------------------------------------------------------------------

package jtag

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// targetJtag is the interface for a target using JTAG.
type targetJtag interface {
	GetJtagDevice() *Device
	GetJtagChain() *Chain
	GetJtagDriver() Driver
}

//-----------------------------------------------------------------------------

var cmdJtagChain = cli.Leaf{
	Descr: "display jtag chain state",
	F: func(c *cli.CLI, args []string) {
		jtagChain := c.User.(targetJtag).GetJtagChain()
		c.User.Put(fmt.Sprintf("%s\n", jtagChain))
	},
}

var cmdJtagDriver = cli.Leaf{
	Descr: "display jtag driver state",
	F: func(c *cli.CLI, args []string) {
		jtagDriver := c.User.(targetJtag).GetJtagDriver()
		c.User.Put(fmt.Sprintf("%s\n", jtagDriver))
	},
}

var cmdJtagSurvey = cli.Leaf{
	Descr: "display jtag device survey",
	F: func(c *cli.CLI, args []string) {
		jtagDevice := c.User.(targetJtag).GetJtagDevice()
		c.User.Put(fmt.Sprintf("%s\n", jtagDevice.Survey()))
	},
}

// Menu submenu items
var Menu = cli.Menu{
	{"chain", cmdJtagChain},
	{"driver", cmdJtagDriver},
	//{"survey", cmdJtagSurvey},
}

//-----------------------------------------------------------------------------
