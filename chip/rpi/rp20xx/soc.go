//-----------------------------------------------------------------------------
/*

RP20xx SoC

Raspberry Pi Dual Core Cortex M0+ (various memory configurations)

*/
//-----------------------------------------------------------------------------

package rp20xx

import (
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

type Variant int

// rp20xx device variants.
const (
	RP2040 Variant = iota
)

var romSize = map[Variant]uint{
	RP2040: 0 * util.KiB,
}

var sramSize = map[Variant]uint{
	RP2040: 256 * util.KiB,
}

// NewSoC returns the SoC device for a rp20xx chip.
func NewSoC(variant Variant) *soc.Device {
	dev := baseSoC()
	// setup the sram /flash
	p := []soc.Peripheral{
		{"SRAM", 0x20000000, sramSize[variant], "Static RAM", nil},
	}
	if romSize[variant] != 0 {
		p = append(p, soc.Peripheral{"ROM", 0x00000000, romSize[variant], "on-chip ROM", nil})
	}
	dev.AddPeripheral(p)
	return dev
}

//-----------------------------------------------------------------------------
