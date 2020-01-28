//-----------------------------------------------------------------------------
/*

BCM47622: SoC with ARM Cortex-A7, 4 core 32-bit ARM.
See: https://www.broadcom.com/products/wireless/wireless-lan-infrastructure/bcm47622

*/
//-----------------------------------------------------------------------------

package bcm47622

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------
/*

BCM47622 JTAG Layout

There are 5 devices on the JTAG chain:

idx 0 bcm47622.dev0 irlen 32 idcode 0x476220a0 mfg 0x050 (Klic) part 0x7622 ver 0x4 leading bit != 1
idx 1 bcm47622.dev1 irlen 32 idcode 0x006dc17f mfg 0x0bf (Broadcom) part 0x06dc ver 0x0
idx 2 bcm47622.dev2 irlen 32 idcode 0x006dc17f mfg 0x0bf (Broadcom) part 0x06dc ver 0x0
idx 3 bcm47622.arm0 irlen 4 idcode 0x5ba00477 mfg 0x23b (ARM Ltd.) part 0xba00 ver 0x5
idx 4 bcm47622.dev3 irlen 5 idcode 0x0d31017f mfg 0x0bf (Broadcom) part 0xd310 ver 0x0

device 3 (ARM Core):

ir 8 drlen 35 # abort
ir 10 drlen 35 # dpacc
ir 11 drlen 35 # apacc
ir 14 drlen 32 # idcode
ir 15 drlen 1 # bypass
(other ir values have drlen == 1)

device 3 has APs:
ap 0: idr 0x44770002 rev 4 jedec 4:3b (ARM) class 8 (MEM-AP) ap 0:2 (APB)
ap 1: idr 0x34770004 rev 3 jedec 4:3b (ARM) class 8 (MEM-AP) ap 0:4 (AXI)

This APB MEM-AP has the following components:

Core 0:
INFO:adiv5:mem-ap 3:0 pidr 00000004004bb906 CTI (Cross Trigger)
INFO:adiv5:mem-ap 3:0 pidr 00000004005bbc07 Cortex-A7 Debug Unit
INFO:adiv5:mem-ap 3:0 pidr 00000004005bb9a7 Cortex-A7 PMU (Performance Monitor Unit)

Core 1:
INFO:adiv5:mem-ap 3:0 pidr 00000004005bbc07 Cortex-A7 Debug Unit
INFO:adiv5:mem-ap 3:0 pidr 00000004005bb9a7 Cortex-A7 PMU (Performance Monitor Unit)

Core 2:
INFO:adiv5:mem-ap 3:0 pidr 00000004005bbc07 Cortex-A7 Debug Unit
INFO:adiv5:mem-ap 3:0 pidr 00000004005bb9a7 Cortex-A7 PMU (Performance Monitor Unit)

Core 3:
INFO:adiv5:mem-ap 3:0 pidr 00000004005bbc07 Cortex-A7 Debug Unit
INFO:adiv5:mem-ap 3:0 pidr 00000004005bb9a7 Cortex-A7 PMU (Performance Monitor Unit)


device 4 (mystery):

This looks the same as the mystery device 3 on the bcm49408.

ir 0 drlen unknown
ir 1 drlen unknown
ir 2 drlen unknown
ir 3 drlen 4
ir 4 drlen unknown
ir 5 drlen unknown
ir 6 drlen 32
ir 7 drlen 1
ir 8 drlen 40
ir 9 drlen 32
ir 10 drlen 4
ir 11 drlen unknown
ir 12 drlen unknown
ir 13 drlen unknown
ir 14 drlen 32
ir 15 drlen 1
ir 16 drlen unknown
ir 17 drlen unknown
ir 18 drlen unknown
ir 19 drlen unknown
ir 20 drlen unknown
ir 21 drlen unknown
ir 22 drlen 32
ir 23 drlen 1
ir 24 drlen unknown
ir 25 drlen unknown
ir 26 drlen unknown
ir 27 drlen unknown
ir 28 drlen unknown
ir 29 drlen unknown
ir 30 drlen 32
ir 31 drlen 1

*/

// ChainInfo0 is the JTAG chain layout for early BCM47622 devices.
var ChainInfo0 = []jtag.DeviceInfo{
	// irlen, idcode, name
	jtag.DeviceInfo{32, jtag.IDCode(0x476220a0), "bcm47622.dev0"}, // some broadcom device
	jtag.DeviceInfo{32, jtag.IDCode(0x006dc17f), "bcm47622.dev1"}, // some broadcom device
	jtag.DeviceInfo{32, jtag.IDCode(0x006dc17f), "bcm47622.dev2"}, // some broadcom device
	jtag.DeviceInfo{4, jtag.IDCode(0x5ba00477), "bcm47622.arm0"},  // ARM core
	jtag.DeviceInfo{5, jtag.IDCode(0x0d31017f), "bcm47622.dev3"},  // some broadcom device
}

// ChainInfo1 is the JTAG chain layout for later BCM47622 devices.
var ChainInfo1 = []jtag.DeviceInfo{
	// irlen, idcode, name
	jtag.DeviceInfo{32, jtag.IDCode(0x11f0617f), "bcm47622.dev0"}, // some broadcom device
	jtag.DeviceInfo{32, jtag.IDCode(0x206dc17f), "bcm47622.dev1"}, // some broadcom device
	jtag.DeviceInfo{32, jtag.IDCode(0x206dc17f), "bcm47622.dev2"}, // some broadcom device
	jtag.DeviceInfo{4, jtag.IDCode(0x5ba00477), "bcm47622.arm0"},  // ARM core
	jtag.DeviceInfo{5, jtag.IDCode(0x0d31017f), "bcm47622.dev3"},  // some broadcom device
}

//-----------------------------------------------------------------------------
