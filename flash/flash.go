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
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the Flash driver api.
type Driver interface {
	GetDefaultRegion() *mem.Region        // get a default region
	GetAddressSize() uint                 // get address size in bits
	LookupSymbol(name string) *mem.Region // lookup a symbol
	GetSectors() []*mem.Region            // return the set of flash sectors
	Erase(r *mem.Region) error            // erase a flash sector
	EraseAll() error                      // erase all of the flash
}

// target provides a method for getting the Flash driver.
type target interface {
	GetFlashDriver() Driver
}

//-----------------------------------------------------------------------------

var helpFlashErase = []cli.Help{
	{"all", "erase all"},
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
		// check for erase all
		if len(args) == 1 && args[0] == "all" {
			c.User.Put("erase all: ")
			err := drv.EraseAll()
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
			} else {
				c.User.Put("done\n")
			}
			return
		}
		// get the memory region
		r, err := mem.RegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// build a list of the flash sectors to be erased
		eraseList := []*mem.Region{}
		for _, s := range drv.GetSectors() {
			if s.Overlaps(r) {
				eraseList = append(eraseList, s)
			}
		}
		if len(eraseList) == 0 {
			c.User.Put("nothing to erase\n")
			return
		}
		// do the erase
		c.User.Put("erasing: ")
		progress := util.NewProgress(c.User, len(eraseList))
		nErased := 0
		nErrors := 0
		for _, s := range eraseList {
			err := drv.Erase(s)
			if err != nil {
				nErrors += 1
			}
			nErased += 1
			progress.Update(nErased)
		}
		progress.Erase()
		c.User.Put(fmt.Sprintf("done (%d errors)\n", nErrors))
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
