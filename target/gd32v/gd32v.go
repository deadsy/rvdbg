//-----------------------------------------------------------------------------
/*

GD32V is a development platform using a GD32VF103VBT6 RISC-V RV32.

See: https://www.seeedstudio.com/SeeedStudio-GD32-RISC-V-Dev-Board-p-4302.html

*/
//-----------------------------------------------------------------------------

package gd32v

import (
	"errors"
	"fmt"
	"gigadevice/gd32vf103"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "gd32v",
	Descr:    "GD32V Board (GigaDevice GD32VF103VBT6 RISC-V RV32)",
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
	jtagDriver jtag.Driver
	jtagChain  *jtag.Chain
	jtagDevice *jtag.Device
	riscvCPU   *riscv.CPU
	socDevice  *soc.Device
}

// New returns a new gd32v target.
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
	jtagChain, err := jtag.NewChain(jtagDriver, gd32vf103.Chain)
	if err != nil {
		return nil, err
	}

	// make the jtag device for the cpu core
	jtagDevice, err := jtagChain.GetDevice(gd32vf103.CoreIndex)
	if err != nil {
		return nil, err
	}

	riscvCPU, err := riscv.NewCPU(jtagDevice)
	if err != nil {
		return nil, err
	}

	// create the SoC device
	socDevice := gd32vf103.NewSoC(gd32vf103.VB)
	socDevice = socDevice.Setup()

	return &Target{
		jtagDriver: jtagDriver,
		jtagChain:  jtagChain,
		jtagDevice: jtagDevice,
		riscvCPU:   riscvCPU,
		socDevice:  socDevice,
	}, nil

}

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return fmt.Sprintf("gd32v.%s> ", t.riscvCPU.PromptState())
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

// GetSoC returns the SoC device and driver.
func (t *Target) GetSoC() (*soc.Device, soc.Driver) {
	return t.socDevice, t.riscvCPU.Dbg
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
