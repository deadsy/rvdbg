//-----------------------------------------------------------------------------
/*

GigaDevice gd32vf103 JTAG Layout

See: https://www.gigadevice.com/products/microcontrollers/gd32/risc-v/

*/
//-----------------------------------------------------------------------------

package gd32vf103

import (
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------

// CoreIndex is the index of the RISC-V core within the JTAG chain.
const CoreIndex = 0

// Chain is the the JTAG chain description.
var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	jtag.DeviceInfo{5, jtag.IDCode(0x1000563d), "gd32v.rv32"},
	jtag.DeviceInfo{5, jtag.IDCode(0x790007a3), "gd32v.dev1"},
}

//-----------------------------------------------------------------------------
