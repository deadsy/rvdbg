//-----------------------------------------------------------------------------
/*

GPIO Access

*/
//-----------------------------------------------------------------------------

package gpio

import (
	"fmt"

	"github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// Driver is the GPIO driver api.
type Driver interface {
	Status() string                        // return a status string for GPIOs
	Pin(name string) (string, uint, error) // convert a pin name to a port/bit tuple
	Set(port string, bit uint) error       // set an output bit
	Clr(port string, bit uint) error       // clear an output bit
}

// target provides a method for getting the GPIO driver.
type target interface {
	GetGpioDriver() Driver
}

//-----------------------------------------------------------------------------

var helpGpio = []cli.Help{
	{"<name>", "gpio name (string), see \"gpio status\" command"},
}

//-----------------------------------------------------------------------------

func gpioArg(drv Driver, args []string) (string, uint, error) {
	err := cli.CheckArgc(args, []int{1})
	if err != nil {
		return "", 0, err
	}
	return drv.Pin(args[0])
}

//-----------------------------------------------------------------------------

var cmdClr = cli.Leaf{
	Descr: "clear gpio (0)",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetGpioDriver()
		port, bit, err := gpioArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		err = drv.Clr(port, bit)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
		}
	},
}

var cmdSet = cli.Leaf{
	Descr: "set gpio (1)",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetGpioDriver()
		port, bit, err := gpioArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		err = drv.Set(port, bit)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
		}
	},
}

var cmdStatus = cli.Leaf{
	Descr: "display gpio status",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetGpioDriver()
		c.User.Put(fmt.Sprintf("%s\n", drv.Status()))
	},
}

//-----------------------------------------------------------------------------

// Menu GPIO submenu items
var Menu = cli.Menu{
	{"clr", cmdClr, helpGpio},
	{"set", cmdSet, helpGpio},
	{"status", cmdStatus},
}

//-----------------------------------------------------------------------------
