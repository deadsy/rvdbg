//-----------------------------------------------------------------------------
/*

SiFive FE310 SoC

*/
//-----------------------------------------------------------------------------

package fe310

import (
	"fmt"

	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

func newI2C(idx, addr uint) soc.Peripheral {
	return soc.Peripheral{
		Name:  fmt.Sprintf("I2C%d", idx),
		Addr:  addr,
		Size:  4 * util.KiB,
		Descr: "Inter-Integrated Circuit",
		Registers: []soc.Register{
			{
				Name:   "PRERlo",
				Offset: 0x00,
				Size:   32,
				Descr:  "Clock Prescale register lo-byte",
			},
			{
				Name:   "PRERhi",
				Offset: 0x04,
				Size:   32,
				Descr:  "Clock Prescale register hi-byte",
			},
			{
				Name:   "CTR",
				Offset: 0x08,
				Size:   32,
				Descr:  "Control register",
			},
			{
				Name:   "TXR",
				Offset: 0x0c,
				Size:   32,
				Descr:  "Transmit register",
			},
			{
				Name:   "RXR",
				Offset: 0x0c,
				Size:   32,
				Descr:  "Receive register",
			},
			{
				Name:   "CR",
				Offset: 0x10,
				Size:   32,
				Descr:  "Command register",
			},
			{
				Name:   "SR",
				Offset: 0x10,
				Size:   32,
				Descr:  "Status register",
			},
		},
	}
}

//-----------------------------------------------------------------------------

type Variant int

const (
	G000 Variant = iota
	G002
)

// NewSoC returns the SoC device for a FE310 chip.
func NewSoC(variant Variant) *soc.Device {
	dev := baseSoC()

	switch variant {
	case G000:
		p := []soc.Peripheral{
			{"debug", 0x100, 0xeff, "Debug ROM", nil},
			{"mask", 0x00001000, 4 * util.KiB, "Mask ROM", nil},
			{"flash", 0x20000000, 512 * util.MiB, "QSPI0 External Flash", nil},
			{"DTIM", 0x80000000, 16 * util.KiB, "E31 DTIM", nil},
		}
		dev.AddPeripheral(p)
	case G002:
		p := []soc.Peripheral{
			{"debug", 0, 0x1000, "Debug ROM", nil},
			{"mode", 0x00001000, 4 * util.KiB, "Mode Select ROM", nil},
			{"error", 0x00003000, 4 * util.KiB, "Error Device ROM", nil},
			{"mask", 0x00010000, 8 * util.KiB, "Mask ROM", nil},
			{"flash", 0x20000000, 512 * util.MiB, "QSPI0 External Flash", nil},
			{"DTIM", 0x80000000, 16 * util.KiB, "E31 DTIM", nil},
			{"ITIM", 0x08000000, 8 * util.KiB, "E31 ITIM", nil},
			newI2C(0, 0x10016000),
		}
		dev.AddPeripheral(p)
	}

	return dev
}

//-----------------------------------------------------------------------------
