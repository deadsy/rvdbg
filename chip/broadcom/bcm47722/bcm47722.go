//-----------------------------------------------------------------------------
/*

BCM47722: SoC with ARM, 4 core ARMv8 CPU
See: https://www.broadcom.com/products/wireless/wireless-lan-infrastructure/bcm47722

*/
//-----------------------------------------------------------------------------

package bcm47722

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------

// CoreIndex is the index of the ARM core on the JTAG chain.
const CoreIndex = 3

// note: this chain layout is wrong
var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	{32, jtag.IDCode(0x476220a0), "bcm47722.dev0"}, // some broadcom device
	{32, jtag.IDCode(0x006dc17f), "bcm47722.dev1"}, // some broadcom device
	{32, jtag.IDCode(0x006dc17f), "bcm47722.dev2"}, // some broadcom device
	{4, jtag.IDCode(0x5ba00477), "bcm47722.arm0"},  // ARM core
	{5, jtag.IDCode(0x0d31017f), "bcm47722.dev3"},  // some broadcom device
}

//-----------------------------------------------------------------------------
