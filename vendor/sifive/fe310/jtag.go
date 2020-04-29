//-----------------------------------------------------------------------------
/*

SiFive FE310 JTAG Layout

*/
//-----------------------------------------------------------------------------

package fe310

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------

// CoreIndex is the index of the RISC-V core within the JTAG chain.
const CoreIndex = 0

// Chain is the the JTAG chain description.
var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	{5, jtag.IDCode(0x20000913), "fe310.rv32"},
}

//-----------------------------------------------------------------------------
