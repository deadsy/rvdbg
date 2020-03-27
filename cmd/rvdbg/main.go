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

	// create the debug interface
	jtagDriver, err := itf.NewJtagDriver(info.DbgType, info.DbgSpeed)
	if err != nil {
		return err
	}
	defer jtagDriver.Close()

	// create the target
	var tgt target.Target
	switch info.Name {
	case "wap":
		tgt, err = wap.New(jtagDriver)
	case "maixgo":
		tgt, err = maixgo.New(jtagDriver)
	case "gd32v":
		tgt, err = gd32v.New(jtagDriver)
	case "redv":
		tgt, err = redv.New(jtagDriver)
	}
	if err != nil {
		return err
	}

	// create the cli
	c := cli.NewCLI(tgt)
	c.HistoryLoad(historyPath)
	c.SetRoot(tgt.GetMenuRoot())

	// run the cli
	for c.Running() {
		// update the prompt to indicate state
		c.SetPrompt(tgt.GetPrompt())
		// run the cli
		c.Run()
	}

	// exit
	c.HistorySave(historyPath)
	tgt.Shutdown()
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
		fmt.Fprintf(os.Stderr, "\ndebug interfaces:\n%s\n", itf.List())
		fmt.Fprintf(os.Stderr, "\ntargets:\n%s\n", target.List())
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
		if info.DbgType == itf.TypeNone {
			fmt.Fprintf(os.Stderr, "use -i to specify an interface name\n")
			fmt.Fprintf(os.Stderr, "\ndebug interfaces:\n%s\n", itf.List())
			os.Exit(1)
		}
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
