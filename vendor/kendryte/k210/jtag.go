//-----------------------------------------------------------------------------
/*

k210 JTAG Layout

*/
//-----------------------------------------------------------------------------

package k210

import (
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

// CoreIndex is the index of the RV64 core on the JTAG chain.
const CoreIndex = 0

// Chain is the the JTAG chain description.
var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	{5, jtag.IDCode(0x04e4796b), "k210-rv64"},
}

//-----------------------------------------------------------------------------
