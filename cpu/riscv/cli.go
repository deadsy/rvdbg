//-----------------------------------------------------------------------------
/*

RISC-V CPU Menu Items

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// target is the interface for a target using a RISC-V CPU.
type target interface {
	GetCpu() *CPU
}

//-----------------------------------------------------------------------------

var CmdHalt = cli.Leaf{
	Descr: "halt the current hart",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		hi := cpu.dbg.GetCurrentHart()
		if hi.State == rv.Halted {
			c.User.Put(fmt.Sprintf("hart%d already halted\n", hi.ID))
			return
		}
		err := cpu.dbg.HaltHart()
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to halt hart%d: %v\n", hi.ID, err))
			return
		}
	},
}

var CmdResume = cli.Leaf{
	Descr: "resume the current hart",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		hi := cpu.dbg.GetCurrentHart()
		if hi.State == rv.Running {
			c.User.Put(fmt.Sprintf("hart%d already running\n", hi.ID))
			return
		}
		err := cpu.dbg.ResumeHart()
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to resume hart%d: %v\n", hi.ID, err))
			return
		}
	},
}

//-----------------------------------------------------------------------------

// HartHelp is help for the hart command.
var HartHelp = []cli.Help{
	{"<cr>", "display info for current hart"},
	{"<id>", "select hart<id> as the current hart"},
}

var CmdHart = cli.Leaf{
	Descr: "hart info/select",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		hi := cpu.dbg.GetCurrentHart()

		if len(args) == 0 {
			c.User.Put(fmt.Sprintf("%s\n", hi))
			return
		}

	},
}

//-----------------------------------------------------------------------------

var cmdRiscvDebug = cli.Leaf{
	Descr: "debug module status",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		c.User.Put(fmt.Sprintf("%s\n", cpu.dbg))
	},
}

var cmdRiscvTest = cli.Leaf{
	Descr: "test routine",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(target).GetCpu()
		c.User.Put(fmt.Sprintf("%s\n", cpu.dbg.Test()))
	},
}

// Menu submenu items
var Menu = cli.Menu{
	{"debug", cmdRiscvDebug},
	{"test", cmdRiscvTest},
}

//-----------------------------------------------------------------------------
