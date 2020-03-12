//-----------------------------------------------------------------------------
/*

RISC-V Debugger

*/
//-----------------------------------------------------------------------------

package main

import (
	"flag"
	"fmt"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/target"
	"github.com/deadsy/rvdbg/target/gd32v"
	"github.com/deadsy/rvdbg/target/maixgo"
	"github.com/deadsy/rvdbg/target/redv"
	"github.com/deadsy/rvdbg/target/wap"
	"github.com/deadsy/rvdbg/util/log"
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

	var app target.Target
	switch info.Name {
	case "wap":
		app, err = wap.New(jtagDriver)
	case "maixgo":
		app, err = maixgo.New(jtagDriver)
	case "gd32v":
		app, err = gd32v.New(jtagDriver)
	case "redv":
		app, err = redv.New(jtagDriver)
	}

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(app.GetMenuRoot())

	// run the cli
	for c.Running() {
		// update the prompt to indicate state
		c.SetPrompt(app.GetPrompt())
		// run the cli
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
	target.Add(&redv.Info)
}

//-----------------------------------------------------------------------------

func main() {

	addTargets()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\ntargets:\n%s\n", target.List())
		fmt.Fprintf(os.Stderr, "\ndebug interfaces:\n%s\n", itf.List())
	}

	targetName := flag.String("t", "", "target name")
	interfaceName := flag.String("i", "", "debug interface name")

	flag.Parse()

	if *targetName == "" {
		fmt.Fprintf(os.Stderr, "use -t to specify a target name\n")
		fmt.Fprintf(os.Stderr, "\ntargets:\n%s\n", target.List())
		os.Exit(1)
	}

	infoPtr := target.Lookup(*targetName)
	if infoPtr == nil {
		fmt.Fprintf(os.Stderr, "target \"%s\" not found\n", *targetName)
		fmt.Fprintf(os.Stderr, "\ntargets:\n%s\n", target.List())
		os.Exit(1)
	}

	// work out the debugger interface type
	info := *infoPtr
	if *interfaceName == "" {
		log.Info.Printf(fmt.Sprintf("using default debug interface: %s", info.DbgType))
	} else {
		x := itf.Lookup(*interfaceName)
		if x == nil {
			fmt.Fprintf(os.Stderr, "debug interface \"%s\" not found\n", *interfaceName)
			fmt.Fprintf(os.Stderr, "debug interfaces:\n%s\n", itf.List())
			os.Exit(1)
		}
		log.Info.Printf(fmt.Sprintf("using debug interface: %s", x.Type))
		info.DbgType = x.Type
	}

	err := run(&info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------
