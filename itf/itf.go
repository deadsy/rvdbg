//-----------------------------------------------------------------------------
/*

Debugger Interface

*/
//-----------------------------------------------------------------------------

package itf

import (
	"errors"
	"fmt"
	"sort"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/itf/daplink"
	"github.com/deadsy/rvdbg/itf/jlink"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

// Type is the debugger interface type.
type Type int

const (
	TypeNone    Type = iota // user must specify the debugger interface to use
	TypeDapLink             // ARM DAPLink
	TypeJlink               // Segger J-Link
	TypeStLink              // ST-LinkV2
)

func (t Type) String() string {
	for k, v := range interfaceDb {
		if t == v.Type {
			return k
		}
	}
	return fmt.Sprintf("unknown (%d)", int(t))
}

//-----------------------------------------------------------------------------

type Info struct {
	Name  string // short name for interface (command line)
	Descr string // description of interface
	Type  Type   // enumerated type for interface
}

var interfaceDb = map[string]*Info{}

// List the supported debugger interface types.
func List() string {
	// sort the interface names
	names := make([]string, 0, len(interfaceDb))
	for k, _ := range interfaceDb {
		names = append(names, k)
	}
	sort.Strings(names)
	// create the ordered interface list
	s := make([][]string, 0, len(names))
	for _, k := range names {
		s = append(s, []string{"", k, interfaceDb[k].Descr})
	}
	return cli.TableString(s, []int{8, 12, 0}, 1)
}

// Add an interface to the database.
func add(info *Info) {
	interfaceDb[info.Name] = info
}

// Lookup an interface by name.
func Lookup(name string) *Info {
	return interfaceDb[name]
}

//-----------------------------------------------------------------------------

// Mode is the debugger interface mode.
type Mode int

const (
	ModeJtag Mode = iota // JTAG
	ModeSwd              // ARM Serial Wire Debug
)

func (m Mode) String() string {
	x := map[Mode]string{
		ModeJtag: "jtag",
		ModeSwd:  "swd",
	}
	if s, ok := x[m]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", int(m))
}

//-----------------------------------------------------------------------------

func init() {
	add(&Info{"daplink", "ARM DAPLink", TypeDapLink})
	add(&Info{"jlink", "Segger J-Link", TypeJlink})
	add(&Info{"stlink", "ST-LinkV2", TypeStLink})
}

//-----------------------------------------------------------------------------

func NewJtagDriver(typ Type, speed int) (jtag.Driver, error) {

	var jtagDriver jtag.Driver

	switch typ {
	case TypeJlink:
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
		jtagDriver, err = jlink.NewJtag(dev, speed)
		if err != nil {
			jlinkLibrary.Shutdown()
			return nil, err
		}

	case TypeDapLink:
		dapLibrary, err := daplink.Init()
		if err != nil {
			return nil, err
		}
		if dapLibrary.NumDevices() == 0 {
			dapLibrary.Shutdown()
			return nil, errors.New("no DAPLink devices found")
		}
		devInfo, err := dapLibrary.DeviceByIndex(0)
		if err != nil {
			dapLibrary.Shutdown()
			return nil, err
		}
		jtagDriver, err = daplink.NewJtag(devInfo, speed)
		if err != nil {
			dapLibrary.Shutdown()
			return nil, err
		}

	default:
		return nil, fmt.Errorf("%s does not support JTAG operations", typ)
	}

	return jtagDriver, nil
}

//-----------------------------------------------------------------------------
