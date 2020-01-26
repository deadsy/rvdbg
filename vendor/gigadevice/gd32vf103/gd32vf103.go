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
	//jtag.DeviceInfo{5, jtag.IDCode(0x0d31017f), "bcm47622.dev3"},  // some broadcom device
}

//-----------------------------------------------------------------------------
