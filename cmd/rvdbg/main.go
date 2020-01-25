//-----------------------------------------------------------------------------
/*

RISC-V Debugger

*/
//-----------------------------------------------------------------------------

package main

import (
	"errors"
	"fmt"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/jlink"
)

//-----------------------------------------------------------------------------

const historyPath = ".rvdbg_history"

//-----------------------------------------------------------------------------

// debugApp is state associated with the RISC-V debugger application.
type debugApp struct {
	jlinkLibrary *jlink.Jlink
	jtagDriver   *jlink.Jtag
	prompt       string
}

// newDebugApp returns a new RISC-V debugger application.
func newDebugApp() (*debugApp, error) {

	jlinkLibrary, err := jlink.Init()
	if err != nil {
		return nil, err
	}

	if jlinkLibrary.NumDevices() == 0 {
		jlinkLibrary.Shutdown()
		return nil, errors.New("no J-Link devices found")
	}

	dev, err := jlinkLibrary.DeviceByIndex(0)
	if err != nil {
		jlinkLibrary.Shutdown()
		return nil, err
	}

	jtagDriver, err := jlink.NewJtag(dev)
	if err != nil {
		jlinkLibrary.Shutdown()
		return nil, err
	}

	return &debugApp{
		jlinkLibrary: jlinkLibrary,
		jtagDriver:   jtagDriver,
		prompt:       "rvdbg> ",
	}, nil
}

func (app *debugApp) Shutdown() {
	app.jtagDriver.Close()
	app.jlinkLibrary.Shutdown()
}

// Put outputs a string to the user application.
func (app *debugApp) Put(s string) {
	os.Stdout.WriteString(s)
}

//-----------------------------------------------------------------------------

func main() {

	// create the application
	app, err := newDebugApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(menuRoot)
	c.SetPrompt(app.prompt)

	// run the cli
	for c.Running() {
		c.Run()
	}

	// exit
	c.HistorySave(historyPath)
	app.Shutdown()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
