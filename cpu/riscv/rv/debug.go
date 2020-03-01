//-----------------------------------------------------------------------------
/*

RISC-V Debugger API

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
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

var stateName = map[HartState]string{
	Running: "running",
	Halted:  "halted",
}

func (s HartState) String() string {
	if name, ok := stateName[s]; ok {
		return name
	}
	return "unknown"
}

// HartInfo stores hart information.
type HartInfo struct {
	ID      int       // hart identifier
	State   HartState // hart state
	Nregs   int       // number of GPRs (normally 32, 16 for rv32e)
	MXLEN   int       // machine XLEN
	SXLEN   int       // supervisor XLEN (0 == no S-mode)
	UXLEN   int       // user XLEN (0 == no U-mode)
	HXLEN   int       // hypervisor XLEN (0 == no H-mode)
	DXLEN   int       // debug XLEN
	FLEN    int       // foating point register width
	MISA    uint      // MISA value
	MHARTID uint      // MHARTID value
}

func xlenString(n int, msg string) string {
	if n != 0 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("no %s", msg)
}

func (hi *HartInfo) String() string {
	s := make([][]string, 0)
	s = append(s, []string{fmt.Sprintf("hart%d", hi.ID), fmt.Sprintf("%s", hi.State)})
	s = append(s, []string{"mhartid", fmt.Sprintf("%d", hi.MHARTID)})
	s = append(s, []string{"misa", fmt.Sprintf("%s", DisplayMISA(hi.MISA, uint(hi.MXLEN)))})
	s = append(s, []string{"nregs", fmt.Sprintf("%d", hi.Nregs)})
	s = append(s, []string{"mxlen", fmt.Sprintf("%d", hi.MXLEN)})
	s = append(s, []string{"sxlen", xlenString(hi.SXLEN, "s-mode")})
	s = append(s, []string{"uxlen", xlenString(hi.UXLEN, "u-mode")})
	s = append(s, []string{"hxlen", xlenString(hi.HXLEN, "h-mode")})
	s = append(s, []string{"flen", xlenString(hi.FLEN, "floating point")})
	s = append(s, []string{"dxlen", fmt.Sprintf("%d", hi.DXLEN)})
	return cli.TableString(s, []int{0, 0}, 1)
}

// Debug is the RISC-V debug interface.
type Debug interface {
	GetInfo() string
	GetHartCount() int
	GetHartInfo(id int) (*HartInfo, error)
	GetCurrentHart() *HartInfo
	SetCurrentHart(id int) (*HartInfo, error)
	HaltHart() error
	ResumeHart() error
	RdGPR(reg uint) (uint64, error)
	RdCSR(reg uint) (uint, error)
	RdFPR(reg uint) (uint64, error)

	Test() string
}

//-----------------------------------------------------------------------------
