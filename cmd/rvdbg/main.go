//-----------------------------------------------------------------------------
/*

RISC-V Debugger

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/target"
	"github.com/deadsy/rvdbg/target/gd32v"
	"github.com/deadsy/rvdbg/target/maixgo"
	"github.com/deadsy/rvdbg/target/wap"
)

//-----------------------------------------------------------------------------

const historyPath = ".rvdbg_history"
const MHz = 1000

//-----------------------------------------------------------------------------

func run(info *target.Info) error {

	jtagDriver, err := itf.NewJtagDriver(info.DbgType, info.DbgSpeed)
	if err != nil {
		return err
	}
	defer jtagDriver.Close()

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

func addTargets() {
	target.Add(&gd32v.Info)
	target.Add(&wap.Info)
	target.Add(&maixgo.Info)
}

//-----------------------------------------------------------------------------

//const targetName = "wap"
const targetName = "gd32v"

//const targetName = "maixgo"

func main() {

	addTargets()

	info := target.Lookup(targetName)
	if info == nil {
		fmt.Fprintf(os.Stderr, "target %s not found\n", targetName)
		os.Exit(1)
	}

	err := run(info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------
