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

// target provides a method for getting the JTAG device.
type target interface {
	GetJtagDevice() *Device
}

//-----------------------------------------------------------------------------

var cmdJtagChain = cli.Leaf{
	Descr: "display jtag chain state",
	F: func(c *cli.CLI, args []string) {
		chain := c.User.(target).GetJtagDevice().chain
		c.User.Put(fmt.Sprintf("%s\n", chain))
	},
}

var cmdJtagDriver = cli.Leaf{
	Descr: "display jtag driver state",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetJtagDevice().drv
		c.User.Put(fmt.Sprintf("%s\n", drv))
	},
}

var cmdJtagSurvey = cli.Leaf{
	Descr: "display jtag device survey",
	F: func(c *cli.CLI, args []string) {
		dev := c.User.(target).GetJtagDevice()
		c.User.Put(fmt.Sprintf("%s\n", dev.Survey()))
	},
}

// Menu submenu items
var Menu = cli.Menu{
	{"chain", cmdJtagChain},
	{"driver", cmdJtagDriver},
	//{"survey", cmdJtagSurvey},
}

//-----------------------------------------------------------------------------
