//-----------------------------------------------------------------------------
/*

Flash Memory Commands

*/
//-----------------------------------------------------------------------------

package flash

import (
	"fmt"
	"time"

	"github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Driver is the Flash driver api.
type Driver interface {
	GetDefaultRegion() *mem.Region         // get a default region
	GetAddressSize() uint                  // get address size in bits
	LookupSymbol(name string) *mem.Region  // lookup a symbol
	GetSectors() []*mem.Region             // return the set of flash regions
	Erase(r *mem.Region) error             // erase a flash region
	EraseAll() error                       // erase all of the flash
	Write(r *mem.Region, buf []byte) error // write a flash region
}

// target provides a method for getting the Flash driver.
type target interface {
	GetFlashDriver() Driver
}

//-----------------------------------------------------------------------------

type flashWriter struct {
	drv    Driver      // flash driver
	region *mem.Region // flash mem ory region to write
}

func newFlashWriter(drv Driver, region *mem.Region) *flashWriter {
	return &flashWriter{
		drv:    drv,
		region: region,
	}
}

func (fw *flashWriter) Write(buf []uint) (int, error) {
	time.Sleep(50 * time.Millisecond)
	return len(buf), nil
}

func (fw *flashWriter) Close() error {
	return nil
}

func (fw *flashWriter) String() string {
	return fmt.Sprintf("%s", fw.region)
}

//-----------------------------------------------------------------------------

var helpFlashErase = []cli.Help{
	{"all", "erase all"},
	{"<addr/name> [len]", "erase memory region"},
	{"  addr", "address (hex), default is 0"},
	{"  name", "region name (string), see \"map\" command"},
	{"  len", "length (hex), defaults to region size"},
}

var helpFlashWrite = []cli.Help{
	{"<filename> <addr/name> [len]", "write a file to flash"},
	{"  filename", "name of file (string)"},
	{"  addr", "address (hex), default is 0"},
	{"  name", "region name (string), see \"map\" command"},
	{"  len", "length (hex), defaults to file size"},
}

var helpFlashProgram = []cli.Help{
	{"<filename>", "write firmware file to flash"},
	{"  filename", "name of file (string)"},
}

//-----------------------------------------------------------------------------

// erase State stores the erase loop state
type eraseState struct {
	drv      Driver         // flash driver
	progress *util.Progress // progress indicator
	sectors  []*mem.Region  // sectors to be erased
	errors   []error        // errors for erase
	idx      int            // index into erase list
}

// eraseLoop is the looping function for sector erasing
func eraseLoop(es *eraseState) bool {
	err := es.drv.Erase(es.sectors[es.idx])
	if err != nil {
		es.errors = append(es.errors, err)
	}
	es.idx++
	es.progress.Update(es.idx)
	return es.idx == len(es.sectors)
}

var cmdErase = cli.Leaf{
	Descr: "erase flash",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		// check for erase all
		if len(args) == 1 && args[0] == "all" {
			c.User.Put("erase all: ")
			err := drv.EraseAll()
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
			} else {
				c.User.Put("done\n")
			}
			return
		}
		// get the memory region
		r, err := mem.RegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// build a list of the flash sectors to be erased
		eraseList := []*mem.Region{}
		for _, s := range drv.GetSectors() {
			if s.Overlaps(r) {
				eraseList = append(eraseList, s)
			}
		}
		if len(eraseList) == 0 {
			c.User.Put("nothing to erase\n")
			return
		}
		// do the erase
		c.User.Put("erasing (ctrl-d to abort): ")
		es := &eraseState{
			drv:      drv,
			progress: util.NewProgress(c.User, len(eraseList)),
			sectors:  eraseList,
			errors:   []error{},
		}
		es.progress.Update(0)
		done := c.Loop(func() bool { return eraseLoop(es) }, cli.KeycodeCtrlD)
		es.progress.Erase()
		status := []string{"abort", "done"}[util.BoolToInt(done)]
		c.User.Put(fmt.Sprintf("%s (%d errors)\n", status, len(es.errors)))
	},
}

var cmdInfo = cli.Leaf{
	Descr: "display flash info",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		s := [][]string{}
		for _, r := range drv.GetSectors() {
			s = append(s, r.ColString())
		}
		c.User.Put(fmt.Sprintf("%s\n", cli.TableString(s, []int{0, 0, 0, 0}, 1)))
	},
}

var cmdProgram = cli.Leaf{
	Descr: "write firmware file to flash",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()
		_ = drv
	},
}

var cmdWrite = cli.Leaf{
	Descr: "write to flash",
	F: func(c *cli.CLI, args []string) {
		drv := c.User.(target).GetFlashDriver()

		// process the arguments
		err := cli.CheckArgc(args, []int{2, 3})
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		name, region, err := mem.FileRegionArg(drv, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		if region.Size == 0 {
			c.User.Put("target region has 0 size\n")
			return
		}
		// work with 32-bit alignment
		region.Align32()

		// file reader
		rd, err := util.NewFileReader(name, 32)
		if err != nil {
			c.User.Put(fmt.Sprintf("unable to open %s (%s)\n", name, err))
			return
		}

		// flash writer
		wr := newFlashWriter(drv, region)

		// copy from file to flash
		cs := util.NewCopyState(rd, wr, 256)
		c.User.Put(fmt.Sprintf("read from %s, write to flash (ctrl-d to abort): ", rd))
		cs.Start(c.User)
		done := c.Loop(func() bool { return cs.CopyLoop() }, cli.KeycodeCtrlD)
		cs.Stop()

		// report result
		if !done {
			c.User.Put("abort\n")
			return
		}
		err = cs.GetError()
		if err != nil {
			c.User.Put(fmt.Sprintf("error (%s)\n", err))
			return
		}
		c.User.Put("done\n")

	},
}

//-----------------------------------------------------------------------------

// Menu GPIO submenu items
var Menu = cli.Menu{
	{"erase", cmdErase, helpFlashErase},
	{"info", cmdInfo},
	{"program", cmdProgram, helpFlashProgram},
	{"write", cmdWrite, helpFlashWrite},
}

//-----------------------------------------------------------------------------
