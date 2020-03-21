//-----------------------------------------------------------------------------
/*

SparkFun RED-V RedBoard - SiFive RISC-V FE310-G002 SoC

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
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "redv",
	Descr:    "SparkFun RED-V RedBoard (SiFive FE310-G002 RISC-V RV32)",
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
	{"map", soc.CmdMap},
	{"mem", mem.Menu, "memory functions"},
	{"regs", soc.CmdRegs, soc.RegsHelp},
	{"resume", riscv.CmdResume},
}

//-----------------------------------------------------------------------------

// Target is the application structure for the target.
type Target struct {
	jtagDevice *jtag.Device
	rvDebug    rv.Debug
	socDevice  *soc.Device
	memDriver  *memDriver
	csrDriver  *csrDriver
	socDriver  *socDriver
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

	// create the CPU debug interface
	rvDebug, err := riscv.NewDebug(jtagDevice)
	if err != nil {
		return nil, err
	}

	// create the SoC device
	socDevice := fe310.NewSoC(fe310.G002).Setup()

	return &Target{
		jtagDevice: jtagDevice,
		rvDebug:    rvDebug,
		socDevice:  socDevice,
		memDriver:  newMemDriver(rvDebug, socDevice),
		socDriver:  newSocDriver(rvDebug),
		csrDriver:  newCsrDriver(rvDebug),
	}, nil

}

//-----------------------------------------------------------------------------

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return t.rvDebug.GetPrompt(Info.Name)
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

//-----------------------------------------------------------------------------

// GetMemoryDriver returns a memory driver for this target.
func (t *Target) GetMemoryDriver() mem.Driver {
	return t.memDriver
}

// GetRiscvDebug returns a RISC-V debug driver for this target.
func (t *Target) GetRiscvDebug() rv.Debug {
	return t.rvDebug
}

// GetSoC returns the SoC device and driver.
func (t *Target) GetSoC() (*soc.Device, soc.Driver) {
	return t.socDevice, t.socDriver
}

// GetCSR returns the CSR device and driver.
func (t *Target) GetCSR() (*soc.Device, soc.Driver) {
	return t.rvDebug.GetCurrentHart().CSR, t.csrDriver
}

// GetJtagDevice returns the JTAG device.
func (t *Target) GetJtagDevice() *jtag.Device {
	return t.jtagDevice
}

//-----------------------------------------------------------------------------
