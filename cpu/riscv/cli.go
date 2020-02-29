//-----------------------------------------------------------------------------
/*

RISC-V CPU Menu Items

*/
//-----------------------------------------------------------------------------

package riscv

import (
	"fmt"
	"math"
	"strings"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
)

//-----------------------------------------------------------------------------

// target is the interface for a target using a RISC-V CPU.
type target interface {
	GetCpu() *CPU
}

//-----------------------------------------------------------------------------
// display general purpose register set

var abiXName = [32]string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}

func gprString(reg []uint64, pc uint64, xlen int) string {
	fmtx := "%08x"
	if xlen == 64 {
		fmtx = "%016x"
	}
	s := make([]string, len(reg)+1)
	for i := 0; i < len(reg); i++ {
		regStr := fmt.Sprintf("x%d", i)
		valStr := "0"
		if reg[i] != 0 {
			valStr = fmt.Sprintf(fmtx, reg[i])
		}
		s[i] = fmt.Sprintf("%-4s %-4s %s", regStr, abiXName[i], valStr)
	}
	s[len(reg)] = fmt.Sprintf("%-9s "+fmtx, "pc", pc)
	return strings.Join(s, "\n")
}

// CmdGpr displays the general purpose registers.
var CmdGpr = cli.Leaf{
	Descr: "display general purpose registers",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(target).GetCpu().dbg
		hi := dbg.GetCurrentHart()
		err := dbg.HaltHart()
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to halt hart%d: %v\n", hi.ID, err))
			return
		}
		reg := make([]uint64, hi.Nregs)
		// read the GPRs
		for i := range reg {
			var err error
			reg[i], err = dbg.RdGPR(uint(i))
			if err != nil {
				c.User.Put(fmt.Sprintf("unable to read gpr%d: %v\n", i, err))
				return
			}
		}
		// read the PC
		pc, err := dbg.RdCSR(rv.DPC)
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to read pc: %v\n", err))
			return
		}
		c.User.Put(fmt.Sprintf("%s\n", gprString(reg, uint64(pc), hi.MXLEN)))
	},
}

//-----------------------------------------------------------------------------
// display floating point register set

var abiFName = [32]string{
	"ft0", "ft1", "ft2", "ft3", "ft4", "ft5", "ft6", "ft7",
	"fs0", "fs1", "fa0", "fa1", "fa2", "fa3", "fa4", "fa5",
	"fa6", "fa7", "fs2", "fs3", "fs4", "fs5", "fs6", "fs7",
	"fs8", "fs9", "fs10", "fs11", "ft8", "ft9", "ft10", "ft11",
}

func fprString(reg []uint64, flen int) string {
	s := make([]string, len(reg))
	for i := 0; i < len(reg); i++ {
		regStr := fmt.Sprintf("f%d", i)
		valStr := "0"
		if reg[i] != 0 {
			valStr = fmt.Sprintf("%016x", reg[i])
		}
		f32 := math.Float32frombits(uint32(reg[i]))
		s[i] = fmt.Sprintf("%-4s %-4s %-16s %f", regStr, abiFName[i], valStr, f32)
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------

// CmdHalt halts the current hart.
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

// CmdResume resumes the current hart.
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

var cmdDebugInfo = cli.Leaf{
	Descr: "debug module information",
	F: func(c *cli.CLI, args []string) {
		dbg := c.User.(target).GetCpu().dbg
		di, err := dbg.GetInfo()
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to get debug info: %v\n", err))
			return
		}
		c.User.Put(fmt.Sprintf("%s\n", di))
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
	{"debug", cmdDebugInfo},
	{"test", cmdRiscvTest},
}

//-----------------------------------------------------------------------------
