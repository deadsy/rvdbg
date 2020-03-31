//-----------------------------------------------------------------------------
/*

SiPeed MaixGo Board using Kendryte K210 RISC-V

Version 0.11 of the debugger spec is implemented.

*/
//-----------------------------------------------------------------------------

package maixgo

import (
	"os"

	"kendryte/k210"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/cpu/riscv/rv11"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "maixgo",
	Descr:    "SiPeed MaixGo (Kendryte K210, Dual Core RISC-V RV64)",
	DbgType:  itf.TypeDapLink,
	DbgMode:  itf.ModeJtag,
	DbgSpeed: 4000,
	Volts:    3300,
}

//-----------------------------------------------------------------------------

// menuRoot is the root menu.
var menuRoot = cli.Menu{
	{"cpu", riscv.Menu, "cpu functions"},
	{"dbg", rv11.Menu, "debugger functions"},
	{"exit", target.CmdExit},
	{"help", target.CmdHelp},
	{"history", target.CmdHistory, cli.HistoryHelp},
	{"jtag", jtag.Menu, "jtag functions"},
	{"map", soc.CmdMap},
	{"regs", soc.CmdRegs, soc.RegsHelp},
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

// New returns a new maixgo target.
func New(jtagDriver jtag.Driver) (target.Target, error) {

	// make the jtag chain
	jtagChain, err := jtag.NewChain(jtagDriver, k210.Chain)
	if err != nil {
		return nil, err
	}

	// make the jtag device for the cpu core
	jtagDevice, err := jtagChain.GetDevice(k210.CoreIndex)
	if err != nil {
		return nil, err
	}

	rvDebug, err := riscv.NewDebug(jtagDevice)
	if err != nil {
		return nil, err
	}

	// create the SoC device
	socDevice := k210.NewSoC().Setup()

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
