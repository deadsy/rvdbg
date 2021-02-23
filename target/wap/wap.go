//-----------------------------------------------------------------------------
/*

WAP is a development platform using a BCM47622

*/
//-----------------------------------------------------------------------------

package wap

import (
	"errors"
	"fmt"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/chip/broadcom/bcm47622"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "wap",
	Descr:    "WAP Board (Broadcom BCM47622, Quad Core ARM 32-bit Cortex-A7)",
	DbgType:  itf.TypeJlink,
	DbgMode:  itf.ModeJtag,
	DbgSpeed: 4000,
	Volts:    3300,
}

//-----------------------------------------------------------------------------

// menuRoot is the root menu.
var menuRoot = cli.Menu{
	{"exit", target.CmdExit},
	{"help", target.CmdHelp},
	{"history", target.CmdHistory, cli.HistoryHelp},
	{"jtag", jtag.Menu, "jtag functions"},
}

//-----------------------------------------------------------------------------

// Target is the application structure for the target.
type Target struct {
	jtagDriver jtag.Driver
	jtagChain  *jtag.Chain
	jtagDevice *jtag.Device
}

// New returns a new wap target.
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
	jtagChain, err := jtag.NewChain(jtagDriver, bcm47622.Chain1)
	if err != nil {
		return nil, err
	}

	// make the jtag device for the cpu core
	jtagDevice, err := jtagChain.GetDevice(bcm47622.CoreIndex)
	if err != nil {
		return nil, err
	}

	return &Target{
		jtagDriver: jtagDriver,
		jtagChain:  jtagChain,
		jtagDevice: jtagDevice,
	}, nil

}

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return "wap> "
}

// GetMenuRoot returns the target root menu.
func (t *Target) GetMenuRoot() []cli.MenuItem {
	return menuRoot
}

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

// Shutdown shuts down the target application.
func (t *Target) Shutdown() {
}

// Put outputs a string to the user application.
func (t *Target) Put(s string) {
	os.Stdout.WriteString(s)
}

//-----------------------------------------------------------------------------
