//-----------------------------------------------------------------------------
/*

BCM47722: SoC with ARM, 4 core ARMv8 CPU
See: https://www.broadcom.com/products/wireless/wireless-lan-infrastructure/bcm47722

*/
//-----------------------------------------------------------------------------

package bcm47722

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------
/*

BCM47722 JTAG Layout

chain: irlen 106 devices 10
device 0: bcm47722.dev0 irlen 16 idcode 0x11a6d17f mfg 0x0bf (Broadcom) part 0x1a6d ver 0x1
device 1: bcm47722.dev1 irlen 32 idcode 0x006f517f mfg 0x0bf (Broadcom) part 0x06f5 ver 0x0
device 2: bcm47722.arm0 irlen 4 idcode 0x0ba02477 mfg 0x23b (ARM Ltd.) part 0xba02 ver 0x0
device 3: bcm47722.arm1 irlen 4 idcode 0x0ba02477 mfg 0x23b (ARM Ltd.) part 0xba02 ver 0x0
device 4: bcm47722.arm2 irlen 4 idcode 0x5ba00477 mfg 0x23b (ARM Ltd.) part 0xba00 ver 0x5
device 5: bcm47722.arm3 irlen 4 idcode 0x5ba00477 mfg 0x23b (ARM Ltd.) part 0xba00 ver 0x5
device 6: bcm47722.dev2 irlen 32 idcode 0x006f517f mfg 0x0bf (Broadcom) part 0x06f5 ver 0x0
device 7: bcm47722.arm4 irlen 4 idcode 0x0ba02477 mfg 0x23b (ARM Ltd.) part 0xba02 ver 0x0
device 8: bcm47722.arm5 irlen 4 idcode 0x5ba00477 mfg 0x23b (ARM Ltd.) part 0xba00 ver 0x5
device 9: bcm47722.dev3 irlen 2 idcode 0x03cb017f mfg 0x0bf (Broadcom) part 0x3cb0 ver 0x0

device 2,3,4,5,7,8 (ARM Cores):
ir 8 drlen 35 # abort
ir 10 drlen 35 # dpacc
ir 11 drlen 35 # apacc
ir 14 drlen 32 # idcode
ir 15 drlen 1 # bypass
(other ir values have drlen == 1)

device 9 (bcm mystery device):
ir 0 drlen 32
ir 1 drlen 64
ir 2 drlen 32
ir 3 drlen 1

*/

// CoreIndex is the index of the ARM core on the JTAG chain.
const CoreIndex = 4

var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	{16, jtag.IDCode(0x11a6d17f), "bcm47722.dev0"}, // some broadcom device
	{32, jtag.IDCode(0x006f517f), "bcm47722.dev1"}, // some broadcom device
	{4, jtag.IDCode(0x0ba02477), "bcm47722.arm0"},  // ARM core
	{4, jtag.IDCode(0x0ba02477), "bcm47722.arm1"},  // ARM core
	{4, jtag.IDCode(0x5ba00477), "bcm47722.arm2"},  // ARM core (main)
	{4, jtag.IDCode(0x5ba00477), "bcm47722.arm3"},  // ARM core
	{32, jtag.IDCode(0x006f517f), "bcm47722.dev2"}, // some broadcom device
	{4, jtag.IDCode(0x0ba02477), "bcm47722.arm4"},  // ARM core
	{4, jtag.IDCode(0x5ba00477), "bcm47722.arm5"},  // ARM core
	{2, jtag.IDCode(0x03cb017f), "bcm47722.dev3"},  // some broadcom device
}

//-----------------------------------------------------------------------------
