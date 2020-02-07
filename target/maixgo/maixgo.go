//-----------------------------------------------------------------------------
/*

SiPeed MaixGo Board using Kendryte K210 RISC-V

*/
//-----------------------------------------------------------------------------

package maixgo

import (
	"os"

	"kendryte/k210"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:  "maixgo",
	Descr: "SiPeed MaixGo (Kendryte K210, Dual Core RISC-V RV64)",
	Itf:   itf.TypeDap,
	Mode:  itf.ModeJtag,
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

// NewMaixGo returns a new target.
func NewTarget(jtagDriver jtag.Driver) (*Target, error) {

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

	return &Target{
		jtagDriver: jtagDriver,
		jtagChain:  jtagChain,
		jtagDevice: jtagDevice,
	}, nil

}

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return "maixgo> "
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