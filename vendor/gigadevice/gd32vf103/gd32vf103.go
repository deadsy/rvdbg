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

device 0 (RISC-V?):

ir 1 drlen 32
ir 16 drlen 32
ir 17 drlen 41
(other ir values have drlen == 1)

device 1 (?)

*/

const CoreIndex = 0

var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	jtag.DeviceInfo{5, jtag.IDCode(0x1000563d), "gd32v.rv32"},
	jtag.DeviceInfo{5, jtag.IDCode(0x790007a3), "gd32v.dev1"},
}

//-----------------------------------------------------------------------------
