//-----------------------------------------------------------------------------
/*

k210
See:

*/
//-----------------------------------------------------------------------------

package k210

import "github.com/deadsy/rvdbg/jtag"

//-----------------------------------------------------------------------------
/*

K210 JTAG Layout

There is 1 device on the JTAG chain.

chain: irlen 5 devices 1
device 0: k210-rv64 irlen 5 idcode 0x04e4796b mfg 0x4b5 (Kendryte) part 0x4e47 ver 0x0

*/

// CoreIndex is the index of the RV64 core on the JTAG chain.
const CoreIndex = 0

var Chain = []jtag.DeviceInfo{
	// irlen, idcode, name
	jtag.DeviceInfo{5, jtag.IDCode(0x04e4796b), "k210-rv64"},
}

//-----------------------------------------------------------------------------
