//-----------------------------------------------------------------------------
/*

SparkFun RED-V RedBoard - SiFive RISC-V FE310 SoC

See: https://www.sparkfun.com/products/15594

*/
//-----------------------------------------------------------------------------

package redv

import (
	"errors"
	"fmt"
	"os"
	"sifive/fe310"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "redv",
	Descr:    "SparkFun RED-V RedBoard (SiFive FE310 RISC-V RV32)",
	DbgType:  itf.TypeJlink,
	DbgMode:  itf.ModeJtag,
	DbgSpeed: 4000,
	Volts:    3300,
}

//-----------------------------------------------------------------------------

// menuRoot is the root menu.
var menuRoot = cli.Menu{
	{"cpu", riscv.Menu, "cpu functions"},
	{"exit", target.CmdExit},
	{"gpr", riscv.CmdGpr},
	{"halt", riscv.CmdHalt},
	{"hart", riscv.CmdHart, riscv.HartHelp},
	{"help", target.CmdHelp},
	{"history", target.CmdHistory, cli.HistoryHelp},
	{"jtag", jtag.Menu, "jtag functions"},
	{"md", mem.DisplayMenu, "memory display functions"},
	{"resume", riscv.CmdResume},
}

//-----------------------------------------------------------------------------

// Target is the application structure for the target.
type Target struct {
	jtagDriver jtag.Driver
	jtagChain  *jtag.Chain
	jtagDevice *jtag.Device
	riscvCPU   *riscv.CPU
}

// New returns a new redv target.
func New(jtagDriver jtag.Driver) (target.Target, error) {

	// get the JTAG state
	state, err := jtagDriver.GetState()
	if err != nil {
		return nil, err
	}

	// check the voltage
	if float32(state.TargetVoltage) < 0.9*float32(Info.Volts) {
		return nil, fmt.Errorf("target voltage is too low (%dmV), is the target connected and powered?", state.TargetVoltage)
	}

	// check the ~SRST state
	if !state.Srst {
		return nil, errors.New("target ~SRST line asserted, target is held in reset")
	}

	// make the jtag chain
	jtagChain, err := jtag.NewChain(jtagDriver, fe310.Chain)
	if err != nil {
		return nil, err
	}

	// make the jtag device for the cpu core
	jtagDevice, err := jtagChain.GetDevice(fe310.CoreIndex)
	if err != nil {
		return nil, err
	}

	riscvCPU, err := riscv.NewCPU(jtagDevice)
	if err != nil {
		return nil, err
	}

	return &Target{
		jtagDriver: jtagDriver,
		jtagChain:  jtagChain,
		jtagDevice: jtagDevice,
		riscvCPU:   riscvCPU,
	}, nil

}

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return fmt.Sprintf("redv.%s> ", t.riscvCPU.PromptState())
}

// GetMenuRoot returns the target root menu.
func (t *Target) GetMenuRoot() []cli.MenuItem {
	return menuRoot
}

// Shutdown shuts down the target application.
func (t *Target) Shutdown() {
}

// Put outputs a string to the user application.
func (t *Target) Put(s string) {
	os.Stdout.WriteString(s)
}

// GetMemoryDriver returns a memory driver for this target.
func (t *Target) GetMemoryDriver() mem.Driver {
	return t.riscvCPU.Dbg
}

// GetRiscvDebug returns a RISC-V debug driver for this target.
func (t *Target) GetRiscvDebug() riscv.Driver {
	return t.riscvCPU.Dbg
}

//-----------------------------------------------------------------------------

// GetJtagDevice returns the JTAG device.
func (t *Target) GetJtagDevice() *jtag.Device {
	return t.jtagDevice
}

// GetJtagChain returns the JTAG chain.
func (t *Target) GetJtagChain() *jtag.Chain {
	return t.jtagChain
}

// GetJtagDriver returns the JTAG driver.
func (t *Target) GetJtagDriver() jtag.Driver {
	return t.jtagDriver
}

// GetCPU returns the RISC-V CPU.
func (t *Target) GetCPU() *riscv.CPU {
	return t.riscvCPU
}

//-----------------------------------------------------------------------------
