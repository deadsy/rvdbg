//-----------------------------------------------------------------------------
/*

ADIv5 JTAG Debug Port

*/
//-----------------------------------------------------------------------------

package arm

import (
	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
)

//-----------------------------------------------------------------------------
// JTAG-DP Registers

// length of instruction register
const irlen = 4

// addresses of data registers
const ir_ABORT = 0x8
const ir_DPACC = 0xa
const ir_APACC = 0xb
const ir_IDCODE = 0xe
const ir_BYPASS = 0xf

// lengths of data registers
const dr_ABORT_LEN = 35
const dr_DPACC_LEN = 35
const dr_APACC_LEN = 35
const dr_IDCODE_LEN = 32
const dr_BYPASS_LEN = 1

var irName = map[int]string{
	ir_ABORT:  "abort",
	ir_DPACC:  "dpacc",
	ir_APACC:  "apacc",
	ir_IDCODE: "idcode",
	ir_BYPASS: "bypass",
}

//-----------------------------------------------------------------------------

// ACK[2:0]
const ack_OK_FAULT = 2
const ack_WAIT = 1

// RnW[0]
const dp_WR = 0
const dp_RD = 1

//-----------------------------------------------------------------------------

type Error struct {
	name string
}

func (e *Error) Error() string {
	return e.name
}

func dpError(name string) error {
	return &Error{name}
}

//-----------------------------------------------------------------------------

type JtagDP struct {
	dev *jtag.Device
	ir  uint
}

func NewJtagDP() *JtagDP {
	dp := JtagDP{}
	return &dp
}

// wr_IR writes the instruction register.
func (dp *JtagDP) wr_IR(ir uint) error {
	if dp.ir == ir {
		// no changes
		return nil
	}
	err := dp.dev.WrIR(bitstr.FromUint(ir, irlen))
	if err != nil {
		return err
	}
	dp.ir = ir
	return nil
}

// rd_IDCODE reads the IDCODE.
func (dp *JtagDP) rd_IDCODE() (uint32, error) {
	err := dp.wr_IR(ir_IDCODE)
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

// wr_ABORT writes the ABORT register.
func (dp *JtagDP) wr_ABORT(val uint) error {
	err := dp.wr_IR(ir_ABORT)
	if err != nil {
		return err
	}
	return dp.dev.WrDR(bitstr.FromUint(val, dr_ABORT_LEN))
}

// rw_DPACC writes to and reads back the DPACC register.
func (dp *JtagDP) rw_DPACC(rnw, addr, val uint) (uint, error) {
	err := dp.wr_IR(ir_DPACC)
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
		return 0, dpError("JTAG-DP ack timeout")
	}
	if ack != ack_OK_FAULT {
		return 0, dpError("JTAG-DP invalid ack")
	}
	return val, nil
}

// rw_APACC writes to and reads back from the selected APACC register.
func (dp *JtagDP) rw_APACC(rnw, addr, val uint) (uint, error) {
	err := dp.wr_IR(ir_APACC)
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
		return 0, dpError("JTAG-DP ack timeout")
	}
	if ack != ack_OK_FAULT {
		return 0, dpError("JTAG-DP invalid ack")
	}
	return val, nil
}

// clr_Errors clears and returns the error bits from the control/status register.
func (dp *JtagDP) clr_Errors() error {
	/*
	   self.rw_dpacc(_DP_RD, _DPACC_CTRL_STAT, 0)
	   x = self.rw_dpacc(_DP_WR, _DPACC_CTRL_STAT, CS_PWR_REQ | CS_ORUNDETECT | CS_ERR)
	   return x & CS_ERR
	*/
	return nil
}

// rd_DPACC reads a DPACC register.
func (dp *JtagDP) rd_DPACC(addr uint) error {
	/*
	  def rd_dpacc(self, adr):
	    self.rw_dpacc(_DP_RD, adr, 0)
	    return self.rw_dpacc(_DP_RD, adr, 0)
	*/
	return nil
}

// wr_DPACC writes a DPACC register.
func (dp *JtagDP) wr_DPACC(addr, val uint) error {
	/*
	  def wr_dpacc(self, adr, val):
	    self.rw_dpacc(_DP_WR, adr, val)
	*/
	return nil
}

// wr_DPACC_Select writes the DPACC select register.
func (dp *JtagDP) wr_DPACC_Select() error {
	/*
	  def wr_dpacc_select(self, ap, reg, dp=0):
	    x = ((ap & 0xff) << 24) | (reg & 0xf0) | (dp & 0xf)
	    self.wr_dpacc(_DPACC_SELECT, x)
	*/
	return nil
}

// rd_RDBUFF returns the RDBUFF value.
func (dp *JtagDP) rd_RDBUFF() error {
	/*
	   def rd_rdbuff(self):
	     return self.rw_dpacc(_DP_RD, _DPACC_RDBUFF, 0)
	*/
	return nil
}

// rd_APACC selects the AP and reads an APACC register.
func (dp *JtagDP) rd_APACC(ap, addr uint) error {
	/*
	  def rd_apacc(self, ap, adr):
	    self.wr_dpacc_select(ap, adr)
	    self.rw_apacc(_DP_RD, adr, 0)
	    return self.rw_apacc(_DP_RD, adr, 0)
	*/
	return nil
}

// wr_APACC selects the AP and writes an APACC register.
func (dp *JtagDP) wr_APACC(ap, addr, val uint) error {
	/*
	  def wr_apacc(self, ap, adr, val):
	    self.wr_dpacc_select(ap, adr)
	    self.rw_apacc(_DP_WR, adr, val)
	*/
	return nil
}

//-----------------------------------------------------------------------------
