//-----------------------------------------------------------------------------
/*

Flash Memory Commands

*/
//-----------------------------------------------------------------------------

package flash

import (
	"fmt"

	"github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/mem"
)

//-----------------------------------------------------------------------------

// Driver is the Flash driver api.
type Driver interface {
	GetSectors() []*mem.Region // return the set of flash sectors
}

// target provides a method for getting the Flash driver.
type target interface {
	GetFlashDriver() Driver
}

//-----------------------------------------------------------------------------

var helpFlashErase = []cli.Help{
	{"*", "erase all"},
	{"<addr/name> [len]", "erase memory region"},
	{"  addr", "address (hex), default is 0"},
	{"  name", "region name (string), see \"map\" command"},
	{"  len", "length (hex), defaults to region size"},
}

var helpFlashWrite = []cli.Help{
	{"<filename> <addr/name> [len]", "write a file to flash"},
	{"  filename", "name of file (string)"},
	{"  addr", "address (hex), default is 0"},
	{"  name", "region name (string), see \"map\" command"},
	{"  len", "length (hex), defaults to file size"},
}

var helpFlashProgram = []cli.Help{
	{"<filename>", "write firmware file to flash"},
	{"  filename", "name of file (string)"},
}

//-----------------------------------------------------------------------------

var cmdErase = cli.Leaf{
	Descr: "erase flash",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		_ = drv
	},
}

var cmdInfo = cli.Leaf{
	Descr: "display flash info",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		s := [][]string{}
		for _, r := range drv.GetSectors() {
			s = append(s, r.ColString())
		}
		c.User.Put(fmt.Sprintf("%s\n", cli.TableString(s, []int{0, 0, 0, 0}, 1)))
	},
}

var cmdProgram = cli.Leaf{
	Descr: "write firmware file to flash",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		_ = drv
	},
}

var cmdWrite = cli.Leaf{
	Descr: "write to flash",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		_ = drv
	},
}

//-----------------------------------------------------------------------------

// Menu GPIO submenu items
var Menu = cli.Menu{
	{"erase", cmdErase, helpFlashErase},
	{"info", cmdInfo},
	{"program", cmdProgram, helpFlashProgram},
	{"write", cmdWrite, helpFlashWrite},
}

//-----------------------------------------------------------------------------
