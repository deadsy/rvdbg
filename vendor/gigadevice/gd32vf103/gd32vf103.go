//-----------------------------------------------------------------------------
/*

gd32vf103
See: https://www.gigadevice.com/products/microcontrollers/gd32/risc-v/

*/
//-----------------------------------------------------------------------------

package gd32vf103

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------
/*

GD32VF103 JTAG Layout

*/

var ChainInfo = []jtag.DeviceInfo{
	// irlen, idcode, name
	jtag.DeviceInfo{4, jtag.IDCode(0x4ba00477), "gd32v.dev0"},
	jtag.DeviceInfo{5, jtag.IDCode(0x790007a3), "gd32v.dev1"},
}

//-----------------------------------------------------------------------------
