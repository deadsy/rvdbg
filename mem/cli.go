//-----------------------------------------------------------------------------
/*

Memory Menu Items

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv"
)

//-----------------------------------------------------------------------------

// targetMem is the interface for a target using these memory functions.
type targetMem interface {
	GetCPU() *riscv.CPU
}

//-----------------------------------------------------------------------------

var cmdDisplay8 = cli.Leaf{
	Descr: "memory display 8-bit",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(targetMem).GetCPU()
		_ = cpu
		c.User.Put("TODO")
	},
}

var cmdDisplay16 = cli.Leaf{
	Descr: "memory display 16-bit",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(targetMem).GetCPU()
		_ = cpu
		c.User.Put("TODO")
	},
}

var cmdDisplay32 = cli.Leaf{
	Descr: "memory display 32-bit",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(targetMem).GetCPU().Dbg
		c.User.Put(fmt.Sprintf("%s\n", Display(dbg, 0x20000000, 0x100, 32)))
	},
}

var cmdDisplay64 = cli.Leaf{
	Descr: "memory display 64-bit",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(targetMem).GetCPU()
		_ = cpu
		c.User.Put("TODO")
	},
}

// DisplayMenu submenu items
var DisplayMenu = cli.Menu{
	{"b", cmdDisplay8},
	{"h", cmdDisplay16},
	{"w", cmdDisplay32},
	{"d", cmdDisplay64},
}

//-----------------------------------------------------------------------------
