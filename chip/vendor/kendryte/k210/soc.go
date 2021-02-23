//-----------------------------------------------------------------------------
/*

Kendryte K210 SoC

*/
//-----------------------------------------------------------------------------

package k210

import (
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// NewSoC returns the SoC device for a K210 chip.
func NewSoC() *soc.Device {
	dev := baseSoC()
	p := []soc.Peripheral{
		{"mem0", 0x40000000, 4 * util.MiB, "CPU SRAM non-cached", nil},
		{"mem1", 0x40400000, 2 * util.MiB, "CPU SRAM non-cached", nil},
		{"mem0c", 0x80000000, 4 * util.MiB, "CPU SRAM cached", nil},
		{"mem1c", 0x80400000, 2 * util.MiB, "CPU SRAM cached", nil},
		{"aimem", 0x40600000, 2 * util.MiB, "AI SRAM non-cached", nil},
		{"aimemc", 0x80600000, 2 * util.MiB, "AI SRAM cached", nil},
		{"rom", 0x88000000, 128 * util.KiB, "Boot ROM", nil},
	}
	dev.AddPeripheral(p)
	return dev
}

//-----------------------------------------------------------------------------
