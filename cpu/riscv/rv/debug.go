//-----------------------------------------------------------------------------
/*

RISC-V Debugger API

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

type HartState int

const (
	Unknown HartState = iota // unknown
	Running                  // hart is running
	Halted                   // hart is halted
)

type HartInfo struct {
	ID    int       // hart identifier
	Xlen  int       // general purpose register bit length
	State HartState // hart state
}

// Debug is the RISC-V debug interface.
type Debug interface {
	GetHartCount() int
	GetHartInfo(id int) (*HartInfo, error)
	SetCurrentHart(id int) error

	Test() string
}

//-----------------------------------------------------------------------------
