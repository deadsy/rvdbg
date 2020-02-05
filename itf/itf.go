//-----------------------------------------------------------------------------
/*

Debugger Interface

*/
//-----------------------------------------------------------------------------

package itf

import "fmt"

//-----------------------------------------------------------------------------

// Type is the debugger interface type.
type Type int

const (
	TypeDap    Type = iota // ARM CMSIS-DAP
	TypeJlink              // Segger J-Link
	TypeStLink             // ST-LinkV2
)

func (t Type) String() string {
	x := map[Type]string{
		TypeDap:    "dap",
		TypeJlink:  "jlink",
		TypeStLink: "stlink",
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
