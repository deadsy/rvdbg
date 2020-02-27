//-----------------------------------------------------------------------------
/*

RISC-V Debugger API

*/
//-----------------------------------------------------------------------------

package rv

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
	ID    int       // hart identifier
	State HartState // hart state
	Mxlen int       // machine XLEN
	Sxlen int       // supervisor XLEN (0 == no S-mode)
	Uxlen int       // user XLEN (0 == no U-mode)
}

// Debug is the RISC-V debug interface.
type Debug interface {
	GetHartCount() int
	GetHartInfo(id int) (*HartInfo, error)
	SetCurrentHart(id int) error

	Test() string
}

//-----------------------------------------------------------------------------
