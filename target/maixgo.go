//-----------------------------------------------------------------------------
/*

SiPeed MaixGo Board using Kendryte K210 RISC-V

*/
//-----------------------------------------------------------------------------

package target

import (
	"os"

	"kendryte/k210"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

// maixGoRoot is the root menu.
var maixGoRoot = cli.Menu{
	{"exit", cmdExit},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"jtag", jtag.Menu, "jtag functions"},
}

//-----------------------------------------------------------------------------

// MaixGo is the application structure for the target.
type MaixGo struct {
	jtagDriver jtag.Driver
	jtagChain  *jtag.Chain
	jtagDevice *jtag.Device
}

// NewMaixGo returns a new target.
func NewMaixGo(jtagDriver jtag.Driver) (*MaixGo, error) {

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

	return &MaixGo{
		jtagDriver: jtagDriver,
		jtagChain:  jtagChain,
		jtagDevice: jtagDevice,
	}, nil

}

// GetPrompt returns the target prompt string.
func (t *MaixGo) GetPrompt() string {
	return "maixgo> "
}

// GetMenuRoot returns the target root menu.
func (t *MaixGo) GetMenuRoot() []cli.MenuItem {
	return maixGoRoot
}

// GetJtagDevice returns the JTAG device.
func (t *MaixGo) GetJtagDevice() *jtag.Device {
	return t.jtagDevice
}

// GetJtagChain returns the JTAG chain.
func (t *MaixGo) GetJtagChain() *jtag.Chain {
	return t.jtagChain
}

// GetJtagDriver returns the JTAG driver.
func (t *MaixGo) GetJtagDriver() jtag.Driver {
	return t.jtagDriver
}

// Shutdown shuts down the target application.
func (t *MaixGo) Shutdown() {
}

// Put outputs a string to the user application.
func (t *MaixGo) Put(s string) {
	os.Stdout.WriteString(s)
}

//-----------------------------------------------------------------------------
