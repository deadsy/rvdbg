//-----------------------------------------------------------------------------
/*

RISC-V CPU Menu Items

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// target is the interface for a target using a RISC-V CPU.
type target interface {
	GetCpu() *CPU
}

//-----------------------------------------------------------------------------

var cmdRiscvDebug = cli.Leaf{
	Descr: "debug module status",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		c.User.Put(fmt.Sprintf("%s\n", cpu.dbg))
	},
}

//-----------------------------------------------------------------------------

var cmdRiscvTest = cli.Leaf{
	Descr: "test routine",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		c.User.Put(fmt.Sprintf("%s\n", cpu.dbg.Test()))
	},
}

//-----------------------------------------------------------------------------

// Menu submenu items
var Menu = cli.Menu{
	{"debug", cmdRiscvDebug},
	{"test", cmdRiscvTest},
}

//-----------------------------------------------------------------------------
