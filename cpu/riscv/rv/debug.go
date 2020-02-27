//-----------------------------------------------------------------------------
/*

RISC-V Debugger API

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

// HartState is the running state of a hart.
type HartState int

// HartState values.
const (
	Unknown HartState = iota // unknown
	Running                  // hart is running
	Halted                   // hart is halted
)

// HartInfo stores hart information.
type HartInfo struct {
	ID      int       // hart identifier
	State   HartState // hart state
	MXLEN   int       // machine XLEN
	SXLEN   int       // supervisor XLEN (0 == no S-mode)
	UXLEN   int       // user XLEN (0 == no U-mode)
	HXLEN   int       // hypervisor XLEN (0 == no H-mode)
	FLEN    int       // foating point register width
	MISA    uint      // MISA value
	MHARTID uint      // MHARTID value
}

func (hi *HartInfo) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("hartid %d", hi.ID))
	s = append(s, fmt.Sprintf("mxlen %d", hi.MXLEN))
	s = append(s, fmt.Sprintf("sxlen %d", hi.SXLEN))
	s = append(s, fmt.Sprintf("uxlen %d", hi.UXLEN))
	s = append(s, fmt.Sprintf("hxlen %d", hi.HXLEN))
	s = append(s, fmt.Sprintf("flen %d", hi.FLEN))
	s = append(s, fmt.Sprintf("misa %s", DisplayMISA(hi.MISA, uint(hi.MXLEN))))
	s = append(s, fmt.Sprintf("mhartid 0x%x", hi.MHARTID))
	return strings.Join(s, "\n")
}

// Debug is the RISC-V debug interface.
type Debug interface {
	GetHartCount() int
	GetHartInfo(id int) (*HartInfo, error)
	GetCurrentHart() *HartInfo
	SetCurrentHart(id int) error

	Test() string
}

//-----------------------------------------------------------------------------
