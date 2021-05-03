//-----------------------------------------------------------------------------
/*

ARM Cortex-M Debugger API

*/
//-----------------------------------------------------------------------------

package cm

import (
	"errors"
	"fmt"

	"github.com/deadsy/rvdbg/swd"
	"github.com/deadsy/rvdbg/util/log"
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

type CmDebug struct {
}

// NewDebug returns a new ARM Cortex-M debugger interface.
func NewDebug(dev *swd.Device) (Debug, error) {

	log.Info.Printf("cortex-m debug module")

	dbg := &CmDebug{}

	return dbg, nil
}

//-----------------------------------------------------------------------------

func (dbg *CmDebug) GetPrompt(name string) string {
	// TODO
	//hi := dbg.GetCurrentHart()
	//state := []rune{'h', 'r'}[util.BoolToInt(hi.State == rv.Running)]
	//return fmt.Sprintf("%s.%d%c> ", name, hi.ID, state)
	return fmt.Sprintf("%s.%d%c> ", name, 0, 'r')
}
func (dbg *CmDebug) RdReg(reg uint) (uint32, error) {
	// TODO
	return 0, nil
}

func (dbg *CmDebug) WrReg(reg uint, val uint32) error {
	// TODO
	return nil
}

func (dbg *CmDebug) GetAddressSize() uint {
	return 32
}

func (dbg *CmDebug) RdMem(width, addr, n uint) ([]uint, error) {
	return nil, errors.New("TODO")
}

func (dbg *CmDebug) WrMem(width, addr uint, val []uint) error {
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
