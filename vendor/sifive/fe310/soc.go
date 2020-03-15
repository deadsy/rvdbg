//-----------------------------------------------------------------------------
/*

SiFive FE310 SoC

*/
//-----------------------------------------------------------------------------

package fe310

import (
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// NewSoC returns the SoC device for a FE310 chip.
func NewSoC() *soc.Device {
	dev := baseSoC()

	// add peripherals
	p := []soc.Peripheral{
		{"debug", 0, 0x1000, "Debug ROM", nil},
		{"mode", 0x00001000, 4 * util.KiB, "Mode Select ROM", nil},
		{"error", 0x00003000, 4 * util.KiB, "Error Device ROM", nil},
		{"mask", 0x00010000, 8 * util.KiB, "Mask ROM", nil},
		{"flash", 0x20000000, 512 * util.MiB, "QSPI0 External Flash", nil},
		{"sram", 0x80000000, 16 * util.KiB, "E31 DTIM", nil},
	}
	// TODO some peripherals need to add/removed as a function of the device variant.
	dev.AddPeripheral(p)
	return dev.Sort()
}

//-----------------------------------------------------------------------------
