//-----------------------------------------------------------------------------
/*

RISC-V Control and Status Registers

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"

	"github.com/deadsy/rvdbg/decode"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

const (
	MSTATUS   = 0x300
	MISA      = 0x301
	DCSR      = 0x7b0
	DPC       = 0x7b1
	DCRATCH0  = 0x7b2
	DCRATCH1  = 0x7b3
	MVENDORID = 0xf11
	MARCHID   = 0xf12
	MIMPID    = 0xf13
	MHARTID   = 0xf14
)

const modeMask = (3 << 8)
const modeUser = (0 << 8)
const modeSupervisor = (1 << 8)
const modeHypervisor = (2 << 8)
const modeMachine = (3 << 8)

//-----------------------------------------------------------------------------

/*

var csrRegs = decode.RegisterSet{
	{"mstatus", MSTATUS, 0, nil, ""},
	{"misa", MISA, 0, nil, ""},
	{"dcsr", DCSR, 0, nil, ""},
	{"dpc", DPC, 0, nil, ""},
	{"dcratch0", DCRATCH0, 0, nil, ""},
	{"dcratch1", DCRATCH1, 0, nil, ""},
	{"mvendorid", MVENDORID, 0, nil, ""},
	{"marchid", MARCHID, 0, nil, ""},
	{"mimpid", MIMPID, 0, nil, ""},
	{"mhartid", MHARTID, 0, nil, ""},
}

*/

//-----------------------------------------------------------------------------
// MISA

func fmtMXL(x uint) string {
	return []string{"?", "32", "64", "128"}[x]
}

func fmtExtensions(x uint) string {
	s := []rune{}
	for i := 0; i < 26; i++ {
		if x&1 != 0 {
			s = append(s, 'a'+rune(i))
		}
		x >>= 1
	}
	if len(s) != 0 {
		return fmt.Sprintf("\"%s\"", string(s))
	}
	return "none"
}

func DisplayMISA(misa, mxlen uint) string {
	fs := decode.FieldSet{
		{"mxl", mxlen - 1, mxlen - 2, fmtMXL},
		{"extensions", 25, 0, fmtExtensions},
	}
	return fs.Display(misa)
}

// GetMxlMISA returns the bit length in the MISA.mxl field.
func GetMxlMISA(misa, mxlen uint) int {
	mxl := util.Bits(misa, mxlen-1, mxlen-2)
	return []int{0, 32, 64, 128}[mxl]
}

// CheckExtMISA returns is the extension is present in MISA.
func CheckExtMISA(misa uint, ext rune) bool {
	n := int(ext) - int('a')
	if n < 0 || n >= 26 {
		return false
	}
	return (misa & (1 << n)) != 0
}

//-----------------------------------------------------------------------------

// GetCSRLength returns the bit length of a CSR register.
func GetCSRLength(reg uint, hi *HartInfo) int {
	// exceptions
	switch reg {
	case MVENDORID, DCSR:
		return 32
	case DPC:
		return hi.DXLEN
	}
	// normal
	switch reg & modeMask {
	case modeUser:
		return hi.UXLEN
	case modeSupervisor:
		return hi.SXLEN
	case modeHypervisor:
		return hi.HXLEN
	case modeMachine:
		return hi.MXLEN
	}
	// default
	return hi.MXLEN
}

//-----------------------------------------------------------------------------
