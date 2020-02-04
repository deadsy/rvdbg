//-----------------------------------------------------------------------------
/*

WAP is a development platform using a BCM47622

*/
//-----------------------------------------------------------------------------

package target

import (
	"broadcom/bcm47622"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

// wapRoot is the root menu.
var wapRoot = cli.Menu{
	{"exit", cmdExit},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"jtag", jtag.Menu, "jtag functions"},
}

//-----------------------------------------------------------------------------

// Wap is the application structure for the target.
type Wap struct {
	jtagDriver jtag.Driver
	jtagChain  *jtag.Chain
	jtagDevice *jtag.Device
}

// NewWap returns a new target.
func NewWap(jtagDriver jtag.Driver) (*Wap, error) {

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

	return &Wap{
		jtagDriver: jtagDriver,
		jtagChain:  jtagChain,
		jtagDevice: jtagDevice,
	}, nil

}

// GetPrompt returns the target prompt string.
func (t *Wap) GetPrompt() string {
	return "wap> "
}

// GetMenuRoot returns the target root menu.
func (t *Wap) GetMenuRoot() []cli.MenuItem {
	return wapRoot
}

// GetJtagDevice returns the JTAG device.
func (t *Wap) GetJtagDevice() *jtag.Device {
	return t.jtagDevice
}

// GetJtagChain returns the JTAG chain.
func (t *Wap) GetJtagChain() *jtag.Chain {
	return t.jtagChain
}

// GetJtagDriver returns the JTAG driver.
func (t *Wap) GetJtagDriver() jtag.Driver {
	return t.jtagDriver
}

// Shutdown shuts down the target application.
func (t *Wap) Shutdown() {
}

// Put outputs a string to the user application.
func (t *Wap) Put(s string) {
	os.Stdout.WriteString(s)
}

//-----------------------------------------------------------------------------
