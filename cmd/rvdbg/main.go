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
	"github.com/deadsy/rvdbg/itf/dap"
	"github.com/deadsy/rvdbg/itf/jlink"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/target/gd32v"
)

//-----------------------------------------------------------------------------

const historyPath = ".rvdbg_history"
const MHz = 1000

//-----------------------------------------------------------------------------

func run(jtagMode string) error {

	var jtagDriver jtag.Driver
	var err error

	switch jtagMode {
	case "J-Link":
		jlinkLibrary, err := jlink.Init()
		if err != nil {
			return err
		}
		defer jlinkLibrary.Shutdown()
		if jlinkLibrary.NumDevices() == 0 {
			return errors.New("no J-Link devices found")
		}
		dev, err := jlinkLibrary.DeviceByIndex(0)
		if err != nil {
			return err
		}
		jtagDriver, err = jlink.NewJtag(dev, 4*MHz)
		if err != nil {
			return err
		}
		defer jtagDriver.Close()

	case "CMSIS-DAP":
		dapLibrary, err := dap.Init()
		if err != nil {
			return err
		}
		defer dapLibrary.Shutdown()
		if dapLibrary.NumDevices() == 0 {
			return errors.New("no CMSIS-DAP devices found")
		}
		devInfo, err := dapLibrary.DeviceByIndex(0)
		if err != nil {
			return err
		}
		jtagDriver, err = dap.NewJtag(devInfo, 4*MHz)
		if err != nil {
			return err
		}
		defer jtagDriver.Close()

	}

	//app, err := wap.NewTarget(jtagDriver)
	//app, err := maixgo.NewTarget(jtagDriver)
	app, err := gd32v.NewTarget(jtagDriver)
	if err != nil {
		return err
	}

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(app.GetMenuRoot())
	c.SetPrompt(app.GetPrompt())

	// run the cli
	for c.Running() {
		c.Run()
	}

	// exit
	c.HistorySave(historyPath)
	app.Shutdown()
	return nil
}

//-----------------------------------------------------------------------------

const jtagMode = "J-Link"

//const jtagMode = "CMSIS-DAP"

func main() {
	err := run(jtagMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------
