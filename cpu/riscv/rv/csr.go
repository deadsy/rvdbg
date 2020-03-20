//-----------------------------------------------------------------------------
/*

RISC-V Control and Status Registers

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"

	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

func NewCsr() *soc.Device {
	return &soc.Device{
		Name: "CSR",
		Peripherals: []soc.Peripheral{
			{
				Name:  "CSR",
				Descr: "CSR Registers",
				Registers: []soc.Register{
					// User CSRs 0x000 - 0x0ff (read/write)
					{Offset: 0x000, Name: "ustatus"},
					{Offset: 0x001, Name: "fflags"},
					{Offset: 0x002, Name: "frm"},
					{Offset: 0x003, Name: "fcsr"},
					{Offset: 0x004, Name: "uie"},
					{Offset: 0x005, Name: "utvec"},
					{Offset: 0x040, Name: "uscratch"},
					{Offset: 0x041, Name: "uepc"},
					{Offset: 0x042, Name: "ucause"},
					{Offset: 0x043, Name: "utval"},
					{Offset: 0x044, Name: "uip"},
					// User CSRs 0xc00 - 0xc7f (read only)
					{Offset: 0xc00, Name: "cycle"},
					{Offset: 0xc01, Name: "time"},
					{Offset: 0xc02, Name: "instret"},
					{Offset: 0xc03, Name: "hpmcounter3"},
					{Offset: 0xc04, Name: "hpmcounter4"},
					{Offset: 0xc05, Name: "hpmcounter5"},
					{Offset: 0xc06, Name: "hpmcounter6"},
					{Offset: 0xc07, Name: "hpmcounter7"},
					{Offset: 0xc08, Name: "hpmcounter8"},
					{Offset: 0xc09, Name: "hpmcounter9"},
					{Offset: 0xc0a, Name: "hpmcounter10"},
					{Offset: 0xc0b, Name: "hpmcounter11"},
					{Offset: 0xc0c, Name: "hpmcounter12"},
					{Offset: 0xc0d, Name: "hpmcounter13"},
					{Offset: 0xc0e, Name: "hpmcounter14"},
					{Offset: 0xc0f, Name: "hpmcounter15"},
					{Offset: 0xc10, Name: "hpmcounter16"},
					{Offset: 0xc11, Name: "hpmcounter17"},
					{Offset: 0xc12, Name: "hpmcounter18"},
					{Offset: 0xc13, Name: "hpmcounter19"},
					{Offset: 0xc14, Name: "hpmcounter20"},
					{Offset: 0xc15, Name: "hpmcounter21"},
					{Offset: 0xc16, Name: "hpmcounter22"},
					{Offset: 0xc17, Name: "hpmcounter23"},
					{Offset: 0xc18, Name: "hpmcounter24"},
					{Offset: 0xc19, Name: "hpmcounter25"},
					{Offset: 0xc1a, Name: "hpmcounter26"},
					{Offset: 0xc1b, Name: "hpmcounter27"},
					{Offset: 0xc1c, Name: "hpmcounter28"},
					{Offset: 0xc1d, Name: "hpmcounter29"},
					{Offset: 0xc1e, Name: "hpmcounter30"},
					{Offset: 0xc1f, Name: "hpmcounter31"},
					// User CSRs 0xc80 - 0xcbf (read only)
					{Offset: 0xc80, Name: "cycleh"},
					{Offset: 0xc81, Name: "timeh"},
					{Offset: 0xc82, Name: "instreth"},
					{Offset: 0xc83, Name: "hpmcounter3h"},
					{Offset: 0xc84, Name: "hpmcounter4h"},
					{Offset: 0xc85, Name: "hpmcounter5h"},
					{Offset: 0xc86, Name: "hpmcounter6h"},
					{Offset: 0xc87, Name: "hpmcounter7h"},
					{Offset: 0xc88, Name: "hpmcounter8h"},
					{Offset: 0xc89, Name: "hpmcounter9h"},
					{Offset: 0xc8a, Name: "hpmcounter10h"},
					{Offset: 0xc8b, Name: "hpmcounter11h"},
					{Offset: 0xc8c, Name: "hpmcounter12h"},
					{Offset: 0xc8d, Name: "hpmcounter13h"},
					{Offset: 0xc8e, Name: "hpmcounter14h"},
					{Offset: 0xc8f, Name: "hpmcounter15h"},
					{Offset: 0xc90, Name: "hpmcounter16h"},
					{Offset: 0xc91, Name: "hpmcounter17h"},
					{Offset: 0xc92, Name: "hpmcounter18h"},
					{Offset: 0xc93, Name: "hpmcounter19h"},
					{Offset: 0xc94, Name: "hpmcounter20h"},
					{Offset: 0xc95, Name: "hpmcounter21h"},
					{Offset: 0xc96, Name: "hpmcounter22h"},
					{Offset: 0xc97, Name: "hpmcounter23h"},
					{Offset: 0xc98, Name: "hpmcounter24h"},
					{Offset: 0xc99, Name: "hpmcounter25h"},
					{Offset: 0xc9a, Name: "hpmcounter26h"},
					{Offset: 0xc9b, Name: "hpmcounter27h"},
					{Offset: 0xc9c, Name: "hpmcounter28h"},
					{Offset: 0xc9d, Name: "hpmcounter29h"},
					{Offset: 0xc9e, Name: "hpmcounter30h"},
					{Offset: 0xc9f, Name: "hpmcounter31h"},
					// Supervisor CSRs 0x100 - 0x1ff (read/write)
					{Offset: 0x100, Name: "sstatus"},
					{Offset: 0x102, Name: "sedeleg"},
					{Offset: 0x103, Name: "sideleg"},
					{Offset: 0x104, Name: "sie"},
					{Offset: 0x105, Name: "stvec"},
					{Offset: 0x106, Name: "scounteren"},
					{Offset: 0x140, Name: "sscratch"},
					{Offset: 0x141, Name: "sepc"},
					{Offset: 0x142, Name: "scause"},
					{Offset: 0x143, Name: "stval"},
					{Offset: 0x144, Name: "sip"},
					{Offset: 0x180, Name: "satp"},
					// Machine CSRs 0xf00 - 0xf7f (read only)
					{Offset: 0xf11, Name: "mvendorid"},
					{Offset: 0xf12, Name: "marchid"},
					{Offset: 0xf13, Name: "mimpid"},
					{Offset: 0xf14, Name: "mhartid"},
					// Machine CSRs 0x300 - 0x3ff (read/write)
					{Offset: 0x300, Name: "mstatus"},
					{Offset: 0x301, Name: "misa"},
					{Offset: 0x302, Name: "medeleg"},
					{Offset: 0x303, Name: "mideleg"},
					{Offset: 0x304, Name: "mie"},
					{Offset: 0x305, Name: "mtvec"},
					{Offset: 0x306, Name: "mcounteren"},
					{Offset: 0x320, Name: "mucounteren"},
					{Offset: 0x321, Name: "mscounteren"},
					{Offset: 0x322, Name: "mhcounteren"},
					{Offset: 0x323, Name: "mhpmevent3"},
					{Offset: 0x324, Name: "mhpmevent4"},
					{Offset: 0x325, Name: "mhpmevent5"},
					{Offset: 0x326, Name: "mhpmevent6"},
					{Offset: 0x327, Name: "mhpmevent7"},
					{Offset: 0x328, Name: "mhpmevent8"},
					{Offset: 0x329, Name: "mhpmevent9"},
					{Offset: 0x32a, Name: "mhpmevent10"},
					{Offset: 0x32b, Name: "mhpmevent11"},
					{Offset: 0x32c, Name: "mhpmevent12"},
					{Offset: 0x32d, Name: "mhpmevent13"},
					{Offset: 0x32e, Name: "mhpmevent14"},
					{Offset: 0x32f, Name: "mhpmevent15"},
					{Offset: 0x330, Name: "mhpmevent16"},
					{Offset: 0x331, Name: "mhpmevent17"},
					{Offset: 0x332, Name: "mhpmevent18"},
					{Offset: 0x333, Name: "mhpmevent19"},
					{Offset: 0x334, Name: "mhpmevent20"},
					{Offset: 0x335, Name: "mhpmevent21"},
					{Offset: 0x336, Name: "mhpmevent22"},
					{Offset: 0x337, Name: "mhpmevent23"},
					{Offset: 0x338, Name: "mhpmevent24"},
					{Offset: 0x339, Name: "mhpmevent25"},
					{Offset: 0x33a, Name: "mhpmevent26"},
					{Offset: 0x33b, Name: "mhpmevent27"},
					{Offset: 0x33c, Name: "mhpmevent28"},
					{Offset: 0x33d, Name: "mhpmevent29"},
					{Offset: 0x33e, Name: "mhpmevent30"},
					{Offset: 0x33f, Name: "mhpmevent31"},
					{Offset: 0x340, Name: "mscratch"},
					{Offset: 0x341, Name: "mepc"},
					{Offset: 0x342, Name: "mcause"},
					{Offset: 0x343, Name: "mtval"},
					{Offset: 0x344, Name: "mip"},
					{Offset: 0x380, Name: "mbase"},
					{Offset: 0x381, Name: "mbound"},
					{Offset: 0x382, Name: "mibase"},
					{Offset: 0x383, Name: "mibound"},
					{Offset: 0x384, Name: "mdbase"},
					{Offset: 0x385, Name: "mdbound"},
					{Offset: 0x3a0, Name: "pmpcfg0"},
					{Offset: 0x3a1, Name: "pmpcfg1"},
					{Offset: 0x3a2, Name: "pmpcfg2"},
					{Offset: 0x3a3, Name: "pmpcfg3"},
					{Offset: 0x3b0, Name: "pmpaddr0"},
					{Offset: 0x3b1, Name: "pmpaddr1"},
					{Offset: 0x3b2, Name: "pmpaddr2"},
					{Offset: 0x3b3, Name: "pmpaddr3"},
					{Offset: 0x3b4, Name: "pmpaddr4"},
					{Offset: 0x3b5, Name: "pmpaddr5"},
					{Offset: 0x3b6, Name: "pmpaddr6"},
					{Offset: 0x3b7, Name: "pmpaddr7"},
					{Offset: 0x3b8, Name: "pmpaddr8"},
					{Offset: 0x3b9, Name: "pmpaddr9"},
					{Offset: 0x3ba, Name: "pmpaddr10"},
					{Offset: 0x3bb, Name: "pmpaddr11"},
					{Offset: 0x3bc, Name: "pmpaddr12"},
					{Offset: 0x3bd, Name: "pmpaddr13"},
					{Offset: 0x3be, Name: "pmpaddr14"},
					{Offset: 0x3bf, Name: "pmpaddr15"},
					// Machine CSRs 0xb00 - 0xb7f (read/write)
					{Offset: 0xb00, Name: "mcycle"},
					{Offset: 0xb02, Name: "minstret"},
					{Offset: 0xb03, Name: "mhpmcounter3"},
					{Offset: 0xb04, Name: "mhpmcounter4"},
					{Offset: 0xb05, Name: "mhpmcounter5"},
					{Offset: 0xb06, Name: "mhpmcounter6"},
					{Offset: 0xb07, Name: "mhpmcounter7"},
					{Offset: 0xb08, Name: "mhpmcounter8"},
					{Offset: 0xb09, Name: "mhpmcounter9"},
					{Offset: 0xb0a, Name: "mhpmcounter10"},
					{Offset: 0xb0b, Name: "mhpmcounter11"},
					{Offset: 0xb0c, Name: "mhpmcounter12"},
					{Offset: 0xb0d, Name: "mhpmcounter13"},
					{Offset: 0xb0e, Name: "mhpmcounter14"},
					{Offset: 0xb0f, Name: "mhpmcounter15"},
					{Offset: 0xb10, Name: "mhpmcounter16"},
					{Offset: 0xb11, Name: "mhpmcounter17"},
					{Offset: 0xb12, Name: "mhpmcounter18"},
					{Offset: 0xb13, Name: "mhpmcounter19"},
					{Offset: 0xb14, Name: "mhpmcounter20"},
					{Offset: 0xb15, Name: "mhpmcounter21"},
					{Offset: 0xb16, Name: "mhpmcounter22"},
					{Offset: 0xb17, Name: "mhpmcounter23"},
					{Offset: 0xb18, Name: "mhpmcounter24"},
					{Offset: 0xb19, Name: "mhpmcounter25"},
					{Offset: 0xb1a, Name: "mhpmcounter26"},
					{Offset: 0xb1b, Name: "mhpmcounter27"},
					{Offset: 0xb1c, Name: "mhpmcounter28"},
					{Offset: 0xb1d, Name: "mhpmcounter29"},
					{Offset: 0xb1e, Name: "mhpmcounter30"},
					{Offset: 0xb1f, Name: "mhpmcounter31"},
					// Machine CSRs 0xb80 - 0xbbf (read/write)
					{Offset: 0xb80, Name: "mcycleh"},
					{Offset: 0xb82, Name: "minstreth"},
					{Offset: 0xb83, Name: "mhpmcounter3h"},
					{Offset: 0xb84, Name: "mhpmcounter4h"},
					{Offset: 0xb85, Name: "mhpmcounter5h"},
					{Offset: 0xb86, Name: "mhpmcounter6h"},
					{Offset: 0xb87, Name: "mhpmcounter7h"},
					{Offset: 0xb88, Name: "mhpmcounter8h"},
					{Offset: 0xb89, Name: "mhpmcounter9h"},
					{Offset: 0xb8a, Name: "mhpmcounter10h"},
					{Offset: 0xb8b, Name: "mhpmcounter11h"},
					{Offset: 0xb8c, Name: "mhpmcounter12h"},
					{Offset: 0xb8d, Name: "mhpmcounter13h"},
					{Offset: 0xb8e, Name: "mhpmcounter14h"},
					{Offset: 0xb8f, Name: "mhpmcounter15h"},
					{Offset: 0xb90, Name: "mhpmcounter16h"},
					{Offset: 0xb91, Name: "mhpmcounter17h"},
					{Offset: 0xb92, Name: "mhpmcounter18h"},
					{Offset: 0xb93, Name: "mhpmcounter19h"},
					{Offset: 0xb94, Name: "mhpmcounter20h"},
					{Offset: 0xb95, Name: "mhpmcounter21h"},
					{Offset: 0xb96, Name: "mhpmcounter22h"},
					{Offset: 0xb97, Name: "mhpmcounter23h"},
					{Offset: 0xb98, Name: "mhpmcounter24h"},
					{Offset: 0xb99, Name: "mhpmcounter25h"},
					{Offset: 0xb9a, Name: "mhpmcounter26h"},
					{Offset: 0xb9b, Name: "mhpmcounter27h"},
					{Offset: 0xb9c, Name: "mhpmcounter28h"},
					{Offset: 0xb9d, Name: "mhpmcounter29h"},
					{Offset: 0xb9e, Name: "mhpmcounter30h"},
					{Offset: 0xb9f, Name: "mhpmcounter31h"},
					// Machine Debug CSRs 0x7a0 - 0x7af (read/write)
					{Offset: 0x7a0, Name: "tselect"},
					{Offset: 0x7a1, Name: "tdata1"},
					{Offset: 0x7a2, Name: "tdata2"},
					{Offset: 0x7a3, Name: "tdata3"},
					// Machine Debug Mode Only CSRs 0x7b0 - 0x7bf (read/write)
					{Offset: 0x7b0, Name: "dcsr"},
					{Offset: 0x7b1, Name: "dpc"},
					{Offset: 0x7b2, Name: "dscratch"},
					// Hypervisor CSRs 0x200 - 0x2ff (read/write)
					{Offset: 0x200, Name: "hstatus"},
					{Offset: 0x202, Name: "hedeleg"},
					{Offset: 0x203, Name: "hideleg"},
					{Offset: 0x204, Name: "hie"},
					{Offset: 0x205, Name: "htvec"},
					{Offset: 0x240, Name: "hscratch"},
					{Offset: 0x241, Name: "hepc"},
					{Offset: 0x242, Name: "hcause"},
					{Offset: 0x243, Name: "hbadaddr"},
					{Offset: 0x244, Name: "hip"},
				},
			},
		},
	}
}

//-----------------------------------------------------------------------------

// CSR register addresses.
const (
	MSTATUS   = 0x300
	MISA      = 0x301
	DCSR      = 0x7b0
	DPC       = 0x7b1
	DSCRATCH0 = 0x7b2
	DSCRATCH1 = 0x7b3
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

// DisplayMISA returns a string decoding a MISA value.
func DisplayMISA(misa, mxlen uint) string {
	fs := []soc.Field{
		{Name: "mxl", Msb: mxlen - 1, Lsb: mxlen - 2, Fmt: fmtMXL},
		{Name: "extensions", Msb: 25, Lsb: 0, Fmt: fmtExtensions},
	}
	return soc.DisplayH(fs, misa)
}

// GetMxlMISA returns the bit length in the MISA.mxl field.
func GetMxlMISA(misa, mxlen uint) uint {
	mxl := util.Bits(misa, mxlen-1, mxlen-2)
	return []uint{0, 32, 64, 128}[mxl]
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

// GetCSRSize returns the CSR register bit size.
func GetCSRSize(reg uint, hi *HartInfo) uint {
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
