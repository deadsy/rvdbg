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
	GetCpu() *Cpu
}

//-----------------------------------------------------------------------------

var cmdRiscvTest = cli.Leaf{
	Descr: "test routine",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		c.User.Put(fmt.Sprintf("%s\n", cpu.dbg.Test()))
	},
}

// Menu submenu items
var Menu = cli.Menu{
	{"test", cmdRiscvTest},
}

//-----------------------------------------------------------------------------
