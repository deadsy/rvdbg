//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13

Hart Functions

*/
//-----------------------------------------------------------------------------

package rv13

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/util"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------
// halt a hart

// isHalted returns true if the currently selected hart is halted.
func (dbg *Debug) isHalted() (bool, error) {
	return dbg.checkStatus(allhalted)
}

const haltTimeout = 5 * time.Millisecond

// halt the current hart, return true if it was already halted.
func (dbg *Debug) halt() (bool, error) {
	// get the current state
	halted, err := dbg.isHalted()
	if err != nil {
		return false, err
	}
	// is the current hart already halted?
	if halted {
		return true, nil
	}
	// request the halt
	err = dbg.setDmi(dmcontrol, haltreq)
	if err != nil {
		return false, err
	}
	// wait for the hart to halt
	t := time.Now().Add(haltTimeout)
	for t.After(time.Now()) {
		halted, err = dbg.isHalted()
		if err != nil {
			return false, err
		}
		if halted {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	// did we timeout?
	if !halted {
		return false, fmt.Errorf("unable to halt hart%d", dbg.hartid)
	}
	// clear the halt request
	err = dbg.clrDmi(dmcontrol, haltreq)
	if err != nil {
		return false, err
	}
	return false, nil
}

//-----------------------------------------------------------------------------
// resume a hart

// isRunning returns true if the currently selected hart is running.
func (dbg *Debug) isRunning() (bool, error) {
	return dbg.checkStatus(allrunning)
}

const resumeTimeout = 5 * time.Millisecond

// resume the current hart, return true if it was already running.
func (dbg *Debug) resume() (bool, error) {
	// get the current state
	running, err := dbg.isRunning()
	if err != nil {
		return false, err
	}
	// is the current hart already running?
	if running {
		return true, nil
	}
	// request the resume
	err = dbg.setDmi(dmcontrol, resumereq)
	if err != nil {
		return false, err
	}
	// wait for the hart to resume
	t := time.Now().Add(resumeTimeout)
	ack := false
	for t.After(time.Now()) {
		ack, err = dbg.checkStatus(allresumeack)
		if err != nil {
			return false, err
		}
		if ack {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	// did we timeout?
	if !ack {
		return false, fmt.Errorf("unable to resume hart%d", dbg.hartid)
	}
	// clear the resume request
	err = dbg.clrDmi(dmcontrol, resumereq)
	if err != nil {
		return false, err
	}
	return false, nil
}

//-----------------------------------------------------------------------------

// getMXLEN returns the GPR length for the current hart.
func (dbg *Debug) getMXLEN() (int, error) {
	// try a 128-bit register read
	_, _, err := dbg.rdReg128(regGPR(rv.RegS0))
	if err == nil {
		return 128, nil
	}
	// try a 64-bit register read
	_, err = dbg.rdReg64(regGPR(rv.RegS0))
	if err == nil {
		return 64, nil
	}
	// try a 32-bit register read
	_, err = dbg.rdReg32(regGPR(rv.RegS0))
	if err == nil {
		return 32, nil
	}
	return 0, errors.New("unable to determine MXLEN")
}

// getFLEN returns the FPR length for the current hart.
func (dbg *Debug) getFLEN() (int, error) {
	// try a 128-bit register read
	_, _, err := dbg.rdReg128(regFPR(0))
	if err == nil {
		return 128, nil
	}
	// try a 64-bit register read
	_, err = dbg.rdReg64(regFPR(0))
	if err == nil {
		return 64, nil
	}
	// try a 32-bit register read
	_, err = dbg.rdReg32(regFPR(0))
	if err == nil {
		return 32, nil
	}
	return 0, errors.New("unable to determine FLEN")
}

// getDXLEN returns the debug register length for the current hart.
func (dbg *Debug) getDXLEN() (int, error) {
	// try a 128-bit register read
	_, _, err := dbg.rdReg128(regCSR(rv.DPC))
	if err == nil {
		return 128, nil
	}
	// try a 64-bit register read
	_, err = dbg.rdReg64(regCSR(rv.DPC))
	if err == nil {
		return 64, nil
	}
	// try a 32-bit register read
	_, err = dbg.rdReg32(regCSR(rv.DPC))
	if err == nil {
		return 32, nil
	}
	return 0, errors.New("unable to determine DXLEN")
}

//-----------------------------------------------------------------------------

// hartInfo stores per hart information.
type hartInfo struct {
	dbg        *Debug      // pointer back to parent debugger
	info       rv.HartInfo // public information
	nscratch   uint        // number of dscratch registers
	datasize   uint        // number of data registers in csr/memory
	dataaccess uint        // data registers in csr(0)/memory(1)
	dataaddr   uint        // csr/memory address
}

func (hi *hartInfo) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("%s", &hi.info))
	s = append(s, fmt.Sprintf("nscratch %d words", hi.nscratch))
	s = append(s, fmt.Sprintf("datasize %d %s", hi.datasize, []string{"csr", "words"}[hi.dataaccess]))
	s = append(s, fmt.Sprintf("dataaccess %s(%d)", []string{"csr", "memory"}[hi.dataaccess], hi.dataaccess))
	s = append(s, fmt.Sprintf("dataaddr 0x%x", hi.dataaddr))
	return strings.Join(s, "\n")
}

func (hi *hartInfo) examine() error {

	dbg := hi.dbg

	// select the hart
	_, err := dbg.SetCurrentHart(hi.info.ID)
	if err != nil {
		return err
	}

	// get the hart status
	x, err := dbg.rdDmi(dmstatus)
	if err != nil {
		return err
	}

	if x&anyhavereset != 0 {
		err := dbg.setDmi(dmcontrol, ackhavereset)
		if err != nil {
			return err
		}
	}

	// halt the hart
	wasHalted, err := dbg.halt()
	if err != nil {
		return err
	}
	hi.info.State = rv.Halted

	// get the MXLEN value
	hi.info.MXLEN, err = dbg.getMXLEN()
	if err != nil {
		return err
	}

	// read the MISA value
	hi.info.MISA, err = dbg.RdCSR(rv.MISA)
	if err != nil {
		return err
	}

	// does MISA.mxl match our MXLEN?
	if rv.GetMxlMISA(hi.info.MISA, uint(hi.info.MXLEN)) != hi.info.MXLEN {
		return errors.New("MXLEN != misa.mxl")
	}

	// are we rv32e?
	if rv.CheckExtMISA(hi.info.MISA, 'e') {
		hi.info.Nregs = 16
	}

	// do we have supervisor mode?
	if rv.CheckExtMISA(hi.info.MISA, 's') {
		if hi.info.MXLEN == 32 {
			hi.info.SXLEN = 32
		} else {
			log.Debug.Printf("TODO")
		}
	}

	// do we have user mode?
	if rv.CheckExtMISA(hi.info.MISA, 'u') {
		if hi.info.MXLEN == 32 {
			hi.info.UXLEN = 32
		} else {
			log.Debug.Printf("TODO")
		}
	}

	// do we have hypervisor mode?
	if rv.CheckExtMISA(hi.info.MISA, 'h') {
		if hi.info.MXLEN == 32 {
			hi.info.HXLEN = 32
		} else {
			log.Debug.Printf("TODO")
		}
	}

	// get the DXLEN value
	hi.info.DXLEN, err = dbg.getDXLEN()
	if err != nil {
		return err
	}

	// get the FLEN value
	hi.info.FLEN, err = dbg.getFLEN()
	if err != nil {
		// ignore errors - we probably don't have floating point support.
		log.Info.Printf(fmt.Sprintf("hart%d: %v", hi.info.ID, err))
	}

	// check 32-bit float support
	if rv.CheckExtMISA(hi.info.MISA, 'f') && hi.info.FLEN < 32 {
		log.Error.Printf("hart%d: misa has 32-bit floating point but FLEN < 32", hi.info.ID)
	}

	// check 64-bit float support
	if rv.CheckExtMISA(hi.info.MISA, 'd') && hi.info.FLEN < 64 {
		log.Error.Printf("hart%d: misa has 64-bit floating point but FLEN < 64", hi.info.ID)
	}

	// check 128-bit float support
	if rv.CheckExtMISA(hi.info.MISA, 'q') && hi.info.FLEN < 128 {
		log.Error.Printf("hart%d: misa has 128-bit floating point but FLEN < 128", hi.info.ID)
	}

	// get the hart id per the CSR
	hi.info.MHARTID, err = dbg.RdCSR(rv.MHARTID)
	if err != nil {
		return err
	}

	// get hartinfo parameters
	x, err = dbg.rdDmi(hartinfo)
	if err != nil {
		return err
	}
	hi.nscratch = util.Bits(uint(x), 23, 20)
	hi.datasize = util.Bits(uint(x), 15, 12)
	hi.dataaccess = util.Bit(uint(x), 16)
	hi.dataaddr = util.Bits(uint(x), 11, 0)

	if !wasHalted {
		// resume the hart
		_, err := dbg.resume()
		if err != nil {
			return err
		}
		hi.info.State = rv.Running
	}

	acs, _ := hi.dbg.rdDmi(abstractcs)
	log.Info.Printf(fmt.Sprintf("acs %08x", acs))

	return nil
}

// newHart creates a hart info structure.
func (dbg *Debug) newHart(id int) *hartInfo {
	hi := &hartInfo{
		dbg: dbg,
	}
	hi.info.ID = id
	hi.info.Nregs = 32
	return hi
}

//-----------------------------------------------------------------------------
