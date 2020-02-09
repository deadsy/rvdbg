//-----------------------------------------------------------------------------
/*

Debugger Interface

*/
//-----------------------------------------------------------------------------

package itf

import (
	"errors"
	"fmt"

	"github.com/deadsy/rvdbg/itf/dap"
	"github.com/deadsy/rvdbg/itf/jlink"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

// Type is the debugger interface type.
type Type int

const (
	TypeCmsisDap Type = iota // ARM CMSIS-DAP
	TypeJlink                // Segger J-Link
	TypeStLink               // ST-LinkV2
)

func (t Type) String() string {
	x := map[Type]string{
		TypeCmsisDap: "cmsis-dap",
		TypeJlink:    "jlink",
		TypeStLink:   "stlink",
	}
	if s, ok := x[t]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", int(t))
}

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

	case TypeCmsisDap:
		dapLibrary, err := dap.Init()
		if err != nil {
			return nil, err
		}
		if dapLibrary.NumDevices() == 0 {
			dapLibrary.Shutdown()
			return nil, errors.New("no CMSIS-DAP devices found")
		}
		devInfo, err := dapLibrary.DeviceByIndex(0)
		if err != nil {
			dapLibrary.Shutdown()
			return nil, err
		}
		jtagDriver, err = dap.NewJtag(devInfo, speed)
		if err != nil {
			dapLibrary.Shutdown()
			return nil, err
		}

	}

	return jtagDriver, nil
}

//-----------------------------------------------------------------------------
