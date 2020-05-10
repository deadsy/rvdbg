//-----------------------------------------------------------------------------
/*

I2C Access

*/
//-----------------------------------------------------------------------------

package i2c

import "github.com/deadsy/go-cli"

//-----------------------------------------------------------------------------

// Driver is the I2C driver api.
type Driver interface {
	Init() error
	Read(adr byte, buf []byte) (int, error)
	Write(adr byte, buf []byte) (int, error)
}

// target provides a method for getting the I2C driver.
type target interface {
	GetI2CDriver() Driver
}

//-----------------------------------------------------------------------------

var helpRead = []cli.Help{
	{"<bus> <adr> [n]", "bus number (hex)"},
	{"", "device address (hex)"},
	{"", "n bytes to read (hex), default is 1"},
}

var helpWrite = []cli.Help{
	{"<bus> <adr> <bytes>", "bus number (hex)"},
	{"", "device address (hex)"},
	{"", "bytes to write (hex)"},
}

var helpBus = []cli.Help{
	{"<bus>", "bus number (hex), default is 0"},
}

//-----------------------------------------------------------------------------

var cmdInit = cli.Leaf{
	Descr: "initialize a bus",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdRead = cli.Leaf{
	Descr: "read bytes from a bus:device",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdWrite = cli.Leaf{
	Descr: "write bytes to a bus:device",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdScan = cli.Leaf{
	Descr: "scan a bus for devices",
	F: func(c *cli.CLI, args []string) {
	},
}

//-----------------------------------------------------------------------------

// Menu GPIO submenu items
var Menu = cli.Menu{
	{"init", cmdInit, helpBus},
	{"rd", cmdRead, helpRead},
	{"scan", cmdScan, helpBus},
	{"wr", cmdWrite, helpWrite},
}

//-----------------------------------------------------------------------------
