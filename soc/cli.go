//-----------------------------------------------------------------------------
/*

SoC Device CLI

*/
//-----------------------------------------------------------------------------

package soc

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the SoC driver api.
type Driver interface {
}

// target provides a method for getting the SoC device and driver.
type target interface {
	GetSoC() (*Device, Driver)
}

//-----------------------------------------------------------------------------

var CmdMap = cli.Leaf{
	Descr: "display memory map",
	F: func(c *cli.CLI, args []string) {
		dev, drv := c.User.(target).GetSoC()
		_ = drv // TODO address format based on MXLEN
		s := make([][]string, len(dev.Peripherals))
		for i, p := range dev.Peripherals {
			var region string
			if p.Size == 0 {
				region = fmt.Sprintf(": %08x", p.Addr)
			} else {
				region = fmt.Sprintf(": %08x %08x %s", p.Addr, p.Addr+p.Size-1, util.MemSize(p.Size))
			}
			s[i] = []string{p.Name, region, p.Descr}
		}
		c.User.Put(fmt.Sprintf("%s\n", cli.TableString(s, []int{0, 0, 0}, 1)))
	},
}

//-----------------------------------------------------------------------------
