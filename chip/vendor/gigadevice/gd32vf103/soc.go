//-----------------------------------------------------------------------------
/*

GigaDevice gd32vf103 SoC

See: https://www.gigadevice.com/products/microcontrollers/gd32/risc-v/

*/
//-----------------------------------------------------------------------------

package gd32vf103

import (
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// Variant is the gd32vf103 device variant.
type Variant int

// gd32vf103 device variants.
const (
	RB Variant = iota
	R8
	R6
	R4
	VB
	V8
	TB
	T8
	T6
	T4
	CB
	C8
	C6
	C4
)

var flashSize = map[Variant]uint{
	RB: 128 * util.KiB,
	R8: 64 * util.KiB,
	R6: 32 * util.KiB,
	R4: 16 * util.KiB,
	VB: 128 * util.KiB,
	V8: 64 * util.KiB,
	TB: 128 * util.KiB,
	T8: 64 * util.KiB,
	T6: 32 * util.KiB,
	T4: 16 * util.KiB,
	CB: 128 * util.KiB,
	C8: 64 * util.KiB,
	C6: 32 * util.KiB,
	C4: 16 * util.KiB,
}

var sramSize = map[Variant]uint{
	RB: 32 * util.KiB,
	R8: 20 * util.KiB,
	R6: 10 * util.KiB,
	R4: 6 * util.KiB,
	VB: 32 * util.KiB,
	V8: 20 * util.KiB,
	TB: 32 * util.KiB,
	T8: 20 * util.KiB,
	T6: 10 * util.KiB,
	T4: 6 * util.KiB,
	CB: 32 * util.KiB,
	C8: 20 * util.KiB,
	C6: 10 * util.KiB,
	C4: 6 * util.KiB,
}

// NewSoC returns the SoC device for a GD32VF103 chip.
func NewSoC(variant Variant) *soc.Device {
	dev := baseSoC()
	// setup the sram /flash
	p := []soc.Peripheral{
		{"sram", 0x20000000, sramSize[variant], "Static RAM", nil},
		{"flash", 0x08000000, flashSize[variant], "Main Flash Block", nil},
		{"boot", 0x1fffb000, 18 * util.KiB, "Boot Loader Block", nil},
		{"option", 0x1ffff800, 16, "SoC Option Bytes", nil},
	}
	// TODO some peripherals need to add/removed as a function of the device variant.
	dev.AddPeripheral(p)
	return dev
}

//-----------------------------------------------------------------------------
