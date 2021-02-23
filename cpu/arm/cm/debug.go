//-----------------------------------------------------------------------------
/*

ARM Cortex-M Debugger API

*/
//-----------------------------------------------------------------------------

package cm

import (
	"errors"

	"github.com/deadsy/rvdbg/swd"
)

//-----------------------------------------------------------------------------

// Debug is the RISC-V debug interface.
type Debug interface {
	GetPrompt(name string) string // get the target prompt
	// registers
	RdReg(reg uint) (uint32, error)   // read general purpose register
	WrReg(reg uint, val uint32) error // write general purpose register
	// memory
	GetAddressSize() uint                      // get address size in bits
	RdMem(width, addr, n uint) ([]uint, error) // read width-bit memory buffer
	WrMem(width, addr uint, val []uint) error  // write width-bit memory buffer
}

// NewDebug returns a new ARM Cortex-M debugger interface.
func NewDebug(dev *swd.Device) (Debug, error) {
	return nil, errors.New("TODO")
}

//-----------------------------------------------------------------------------
