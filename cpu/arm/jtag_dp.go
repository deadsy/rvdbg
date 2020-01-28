//-----------------------------------------------------------------------------
/*

ADIv5 JTAG Debug Port

*/
//-----------------------------------------------------------------------------

package arm

import (
	"errors"

	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------
// JTAG-DP Registers

// length of instruction register
const irLength = 4

// addresses of data registers
const irABORT = 0x8
const irDPACC = 0xa
const irAPACC = 0xb
const irIDCODE = 0xe
const irBYPASS = 0xf

// lengths of data registers
const dr_ABORT_LEN = 35
const dr_DPACC_LEN = 35
const dr_APACC_LEN = 35
const dr_IDCODE_LEN = 32
const dr_BYPASS_LEN = 1

var irName = map[int]string{
	irABORT:  "abort",
	irDPACC:  "dpacc",
	irAPACC:  "apacc",
	irIDCODE: "idcode",
	irBYPASS: "bypass",
}

//-----------------------------------------------------------------------------

// ACK[2:0]
const ack_OK_FAULT = 2
const ack_WAIT = 1

// RnW[0]
const dp_WR = 0
const dp_RD = 1

//-----------------------------------------------------------------------------
// Debug Port Register Access (DPACC)

const dpacc_CTRL_STAT = 0x4 // read/write
const dpacc_SELECT = 0x8    // read/write
const dpacc_RDBUFF = 0xc    // read only

var dpaccName = map[int]string{
	dpacc_CTRL_STAT: "ctrl/stat",
	dpacc_SELECT:    "select",
	dpacc_RDBUFF:    "rdbuff",
}

//-----------------------------------------------------------------------------
// ABORT Register

const abort_ORUNERRCLR = (1 << 4) // Clear the STICKYORUN flag, SW-DP only
const abort_WDERRCLR = (1 << 3)   // Clear the WDATAERR flag, SW-DP only
const abort_STKERRCLR = (1 << 2)  // Clear the STICKYERR flag, SW-DP only
const abort_STKCMPCLR = (1 << 1)  // Clear the STICKYCMP flag, SW-DP only
const abort_DAPABORT = (1 << 0)   // Generate a DAP abort

//-----------------------------------------------------------------------------
// DPACC CTRL/STAT register

const cs_ORUNDETECT = (1 << 0)
const cs_STICKYORUN = (1 << 1)

// 3:2 - transaction mode (e.g. pushed compare)
const cs_STICKYCMP = (1 << 4)
const cs_STICKYERR = (1 << 5)
const cs_READOK = (1 << 6)   // SWD-only
const cs_WDATAERR = (1 << 7) // SWD-only
// 11:8 - mask lanes for pushed compare or verify ops
// 21:12 - transaction counter
const cs_DBGRSTREQ = (1 << 26)
const cs_DBGRSTACK = (1 << 27)
const cs_DBGPWRUPREQ = (1 << 28) // debug power up request
const cs_DBGPWRUPACK = (1 << 29) // debug power up acknowledge (read only)
const cs_SYSPWRUPREQ = (1 << 30) // system power up request
const cs_SYSPWRUPACK = (1 << 31) // system power up acknowledge (read only)

const cs_PWR_REQ = cs_DBGPWRUPREQ | cs_SYSPWRUPREQ
const cs_PWR_ACK = cs_DBGPWRUPACK | cs_SYSPWRUPACK
const cs_ERR = cs_STICKYORUN | cs_STICKYCMP | cs_STICKYERR

//-----------------------------------------------------------------------------

// JtagDP is a JTAG-DP access object.
type JtagDP struct {
	dev *jtag.Device
	ir  uint
}

// NewJtagDP returns a new JTAG-DP access object.
func NewJtagDP(dev *jtag.Device) *JtagDP {
	return &JtagDP{
		dev: dev,
	}
}

// WrIR writes the instruction register.
func (dp *JtagDP) WrIR(ir uint) error {
	if dp.ir == ir {
		// no changes
		return nil
	}
	err := dp.dev.WrIR(bitstr.FromUint(ir, irLength))
	if err != nil {
		return err
	}
	dp.ir = ir
	return nil
}

// RdIDCODE reads the IDCODE.
func (dp *JtagDP) RdIDCODE() (uint32, error) {
	err := dp.WrIR(irIDCODE)
	if err != nil {
		return 0, err
	}
	x, err := dp.dev.RdWrDR(bitstr.Zeroes(dr_IDCODE_LEN))
	if err != nil {
		return 0, err
	}
	idcode := uint32(x.Split([]int{32})[0])
	return idcode, nil
}

// WrABORT writes the ABORT register.
func (dp *JtagDP) WrABORT(val uint) error {
	err := dp.WrIR(irABORT)
	if err != nil {
		return err
	}
	return dp.dev.WrDR(bitstr.FromUint(val, dr_ABORT_LEN))
}

// RdWrDPACC reads and writes from a DPACC register.
func (dp *JtagDP) RdWrDPACC(rnw, addr, val uint) (uint, error) {
	err := dp.WrIR(irDPACC)
	if err != nil {
		return 0, err
	}
	val = (val << 3) | ((addr >> 1) & 0x06) | rnw
	rd, err := dp.dev.RdWrDR(bitstr.FromUint(val, dr_DPACC_LEN))
	if err != nil {
		return 0, err
	}
	x := rd.Split([]int{3, 32})
	ack := x[0]
	val = x[1]
	if ack == ack_WAIT {
		return 0, errors.New("JTAG-DP ack timeout")
	}
	if ack != ack_OK_FAULT {
		return 0, errors.New("JTAG-DP invalid ack")
	}
	return val, nil
}

// RdWrAPACC reads and writes from a APACC register.
func (dp *JtagDP) RdWrAPACC(rnw, addr, val uint) (uint, error) {
	err := dp.WrIR(irAPACC)
	if err != nil {
		return 0, err
	}
	val = (val << 3) | ((addr >> 1) & 0x06) | rnw
	rd, err := dp.dev.RdWrDR(bitstr.FromUint(val, dr_APACC_LEN))
	if err != nil {
		return 0, err
	}
	x := rd.Split([]int{3, 32})
	ack := x[0]
	val = x[1]
	if ack == ack_WAIT {
		return 0, errors.New("JTAG-DP ack timeout")
	}
	if ack != ack_OK_FAULT {
		return 0, errors.New("JTAG-DP invalid ack")
	}
	return val, nil
}

// ClrErrors clears and returns the error bits from the control/status register.
func (dp *JtagDP) ClrErrors() (uint, error) {
	_, err := dp.RdWrDPACC(dp_RD, dpacc_CTRL_STAT, 0)
	if err != nil {
		return 0, err
	}
	val, err := dp.RdWrDPACC(dp_WR, dpacc_CTRL_STAT, cs_PWR_REQ|cs_ORUNDETECT|cs_ERR)
	if err != nil {
		return 0, err
	}
	return val & cs_ERR, nil
}

// RdDPACC reads a DPACC register.
func (dp *JtagDP) RdDPACC(addr uint) (uint, error) {
	_, err := dp.RdWrDPACC(dp_RD, addr, 0)
	if err != nil {
		return 0, err
	}
	return dp.RdWrDPACC(dp_RD, addr, 0)
}

// WrDPACC writes a DPACC register.
func (dp *JtagDP) WrDPACC(addr, val uint) error {
	_, err := dp.RdWrDPACC(dp_WR, addr, val)
	return err
}

// WrDPACC_Select writes the DPACC select register.
func (dp *JtagDP) WrDPACC_Select(ap, reg, xdp uint) error {
	val := ((ap & 0xff) << 24) | (reg & 0xf0) | (xdp & 0xf)
	return dp.WrDPACC(dpacc_SELECT, val)
}

// RdRDBUFF returns the RDBUFF value.
func (dp *JtagDP) RdRDBUFF() (uint, error) {
	return dp.RdWrDPACC(dp_RD, dpacc_RDBUFF, 0)
}

// RdAPACC selects the AP and reads an APACC register.
func (dp *JtagDP) RdAPACC(ap, addr uint) (uint, error) {
	err := dp.WrDPACC_Select(ap, addr, 0)
	if err != nil {
		return 0, err
	}
	_, err = dp.RdWrAPACC(dp_RD, addr, 0)
	if err != nil {
		return 0, err
	}
	return dp.RdWrAPACC(dp_RD, addr, 0)
}

// WrAPACC selects the AP and writes an APACC register.
func (dp *JtagDP) WrAPACC(ap, addr, val uint) error {
	err := dp.WrDPACC_Select(ap, addr, 0)
	if err != nil {
		return err
	}
	_, err = dp.RdWrAPACC(dp_WR, addr, val)
	return err
}

//-----------------------------------------------------------------------------
