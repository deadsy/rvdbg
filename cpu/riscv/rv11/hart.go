//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11

Hart Functions

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deadsy/rvda"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/util"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------
// halt a hart

// isHalted returns true if the currently selected hart is halted.
func (dbg *Debug) isHalted() (bool, error) {
	return false, errors.New("TODO")
}

// halt the current hart, return true if it was already halted.
func (dbg *Debug) halt() (bool, error) {
	dbg.cache.wr32(0, rv.InsCSRSI(rv.DCSR, 1<<3 /*halt*/))
	dbg.cache.wr32(1, rv.InsCSRR(rv.RegS0, rv.MHARTID))
	dbg.cache.wr32(2, rv.InsSW(rv.RegS0, debugSetHaltNotification, rv.RegZero))
	dbg.cache.wrResume(3)
	// run the code
	err := dbg.cache.flush(true)
	return true, err
}

//-----------------------------------------------------------------------------
// resume a hart

// isRunning returns true if the currently selected hart is running.
func (dbg *Debug) isRunning() (bool, error) {
	return false, errors.New("TODO")
}

// resume the current hart, return true if it was already running.
func (dbg *Debug) resume() (bool, error) {

	return false, nil

	//return false, errors.New("TODO")
}

//-----------------------------------------------------------------------------
// access probing- setup pointers to access functions

// probeGPR works out how we can access GPRs
func (hi *hartInfo) probeGPR() error {
	hi.rdGPR = rdGPR
	hi.wrGPR = wrGPR
	return nil
}

// probeFPR works out how we can access FPRs
func (hi *hartInfo) probeFPR() error {
	hi.rdFPR = rdFPR
	hi.wrFPR = wrFPR
	return nil
}

// probeCSR works out how we can access CSRs
func (hi *hartInfo) probeCSR() error {
	hi.rdCSR = rdCSR
	hi.wrCSR = wrCSR
	return nil
}

// probeMemory works out how we can access memory.
func (hi *hartInfo) probeMemory() error {
	if hi.info.MXLEN == 64 {
		hi.rdMem = rv64RdMem
		hi.wrMem = rv64WrMem
		return nil
	}
	if hi.info.MXLEN == 32 {
		hi.rdMem = rv32RdMem
		hi.wrMem = rv32WrMem
		return nil
	}
	return errors.New("unable to support memory access")
}

func (dbg *Debug) probeAccess() error {
	hi := dbg.hart[dbg.hartid]
	// GPRs
	err := hi.probeGPR()
	if err != nil {
		return err
	}
	// FPRs
	err = hi.probeFPR()
	if err != nil {
		return err
	}
	// CSRs
	err = hi.probeCSR()
	if err != nil {
		return err
	}
	// Memory
	err = hi.probeMemory()
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

// getMXLEN returns the GPR length for the current hart.
func (dbg *Debug) getMXLEN() (uint, error) {

	// here's the method:
	// s1 = -1,   s1 = 0xffffffff ffffffff ffffffff ffffffff
	// s1 >>= 31, s1 = 0x00000001 ffffffff ffffffff ffffffff -> x0 result
	// s1 >>= 31, s1 = 0x00000000 00000003 ffffffff ffffffff -> x1 result
	// looking at the least significant 32-bits:
	// rv32:  x0 = 0x00000001, x1 = 0x00000000
	// rv64:  x0 = 0xffffffff, x1 = 0x00000003
	// rv128: x0 = 0xffffffff, x1 = 0xffffffff

	dbg.cache.wr32(0, rv.InsXORI(rv.RegS1, rv.RegZero, ^uint(0)))
	dbg.cache.wr32(1, rv.InsSRLI(rv.RegS1, rv.RegS1, 31))
	dbg.cache.wr32(2, rv.InsSW(rv.RegS1, ramAddr(0), rv.RegZero))
	dbg.cache.wr32(3, rv.InsSRLI(rv.RegS1, rv.RegS1, 31))
	dbg.cache.wr32(4, rv.InsSW(rv.RegS1, ramAddr(1), rv.RegZero))
	dbg.cache.wrResume(5)
	dbg.cache.read(0)
	dbg.cache.read(1)

	// run the code
	err := dbg.cache.flush(true)
	if err != nil {
		return 0, err
	}

	// get the results
	x0 := dbg.cache.rd32(0)
	x1 := dbg.cache.rd32(1)
	if x0 == 1 && x1 == 0 {
		return 32, nil
	}
	if x0 == 0xffffffff && x1 == 3 {
		return 64, nil
	}
	if x0 == 0xffffffff && x1 == 0xffffffff {
		return 128, nil
	}
	return 0, errors.New("unable to determine MXLEN")
}

// getSXLEN returns the supervisor mode register length for the current hart.
func (dbg *Debug) getSXLEN(hi *hartInfo) (uint, error) {
	if hi.info.MXLEN == 32 {
		return 32, nil
	}
	// read the MSTATUS value
	mstatus, err := dbg.RdCSR(rv.MSTATUS, 0)
	if err != nil {
		return 0, err
	}
	sxl := util.Bits(uint(mstatus), 35, 34)
	if sxl == 0 {
		log.Info.Printf("hart%d: misa indicates s-mode, but mstatus.sxl = 0", hi.info.ID)
		return hi.info.MXLEN, nil
	}
	return []uint{0, 32, 64, 128}[sxl], nil
}

// getUXLEN returns the user mode register length for the current hart.
func (dbg *Debug) getUXLEN(hi *hartInfo) (uint, error) {
	if hi.info.MXLEN == 32 {
		return 32, nil
	}
	// read the MSTATUS value
	mstatus, err := dbg.RdCSR(rv.MSTATUS, 0)
	if err != nil {
		return 0, err
	}
	uxl := util.Bits(uint(mstatus), 33, 32)
	if uxl == 0 {
		log.Info.Printf("hart%d: misa indicates u-mode, but mstatus.uxl = 0", hi.info.ID)
		return hi.info.MXLEN, nil
	}
	return []uint{0, 32, 64, 128}[uxl], nil
}

// getHXLEN returns the hypervisor mode register length for the current hart.
func (dbg *Debug) getHXLEN(hi *hartInfo) (uint, error) {
	if hi.info.MXLEN == 32 {
		return 32, nil
	}
	return 0, errors.New("TODO")
}

// getFLEN returns the FPR length for the current hart.
func (dbg *Debug) getFLEN() (uint, error) {
	// try a 64-bit register read
	_, err := dbg.RdFPR(0, 64)
	if err == nil {
		return 64, nil
	}
	// try a 32-bit register read
	_, err = dbg.RdFPR(0, 32)
	if err == nil {
		return 32, nil
	}
	return 0, errors.New("unable to determine FLEN")
}

// getDXLEN returns the debug register length for the current hart.
func (dbg *Debug) getDXLEN(hi *hartInfo) (uint, error) {
	// try a 64-bit register read
	_, err := dbg.RdCSR(rv.DPC, 64)
	if err == nil {
		return 64, nil
	}
	// try a 32-bit register read
	_, err = dbg.RdCSR(rv.DPC, 32)
	if err == nil {
		return 32, nil
	}
	return 0, errors.New("unable to determine DXLEN")
}

//-----------------------------------------------------------------------------

type rdRegFunc func(dbg *Debug, reg, size uint) (uint64, error)
type wrRegFunc func(dbg *Debug, reg, size uint, val uint64) error
type rdMemFunc func(dbg *Debug, width, addr, n uint) ([]uint, error)
type wrMemFunc func(dbg *Debug, width, addr uint, val []uint) error

// hartInfo stores generic/rv13 hart information.
type hartInfo struct {
	dbg   *Debug      // pointer back to parent debugger
	info  rv.HartInfo // generic information
	rdGPR rdRegFunc   // read GPR function
	rdFPR rdRegFunc   // read FPR function
	rdCSR rdRegFunc   // read CSR function
	wrGPR wrRegFunc   // write GPR function
	wrFPR wrRegFunc   // write FPR function
	wrCSR wrRegFunc   // write CSR function
	rdMem rdMemFunc   // read memory buffer
	wrMem wrMemFunc   // write memory buffer

	/*
		nscratch   uint        // number of dscratch registers
		datasize   uint        // number of data registers in csr/memory
		dataaccess uint        // data registers in csr(0)/memory(1)
		dataaddr   uint        // csr/memory address
	*/
}

func (hi *hartInfo) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("%s", &hi.info))
	//s = append(s, fmt.Sprintf("nscratch %d words", hi.nscratch))
	//s = append(s, fmt.Sprintf("datasize %d %s", hi.datasize, []string{"csr", "words"}[hi.dataaccess]))
	//s = append(s, fmt.Sprintf("dataaccess %s(%d)", []string{"csr", "memory"}[hi.dataaccess], hi.dataaccess))
	//s = append(s, fmt.Sprintf("dataaddr 0x%x", hi.dataaddr))
	return strings.Join(s, "\n")
}

func (hi *hartInfo) examine() error {

	dbg := hi.dbg

	// select the hart
	_, err := dbg.SetCurrentHart(hi.info.ID)
	if err != nil {
		return err
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
	log.Info.Printf("hart%d: MXLEN %d", hi.info.ID, hi.info.MXLEN)

	// probe the access modes
	err = dbg.probeAccess()
	if err != nil {
		return err
	}

	// read the MISA value
	misa, err := dbg.RdCSR(rv.MISA, 0)
	if err != nil {
		return err
	}
	hi.info.MISA = uint(misa)
	log.Info.Printf("hart%d: MISA 0x%x", hi.info.ID, hi.info.MISA)

	// does MISA.mxl match our MXLEN?
	if rv.GetMxlMISA(hi.info.MISA, hi.info.MXLEN) != hi.info.MXLEN {
		return errors.New("MXLEN != misa.mxl")
	}

	// are we rv32e?
	if rv.CheckExtMISA(hi.info.MISA, 'e') {
		hi.info.Nregs = 16
	}

	// do we have supervisor mode?
	if rv.CheckExtMISA(hi.info.MISA, 's') {
		hi.info.SXLEN, err = dbg.getSXLEN(hi)
		if err != nil {
			return err
		}
	}

	// do we have user mode?
	if rv.CheckExtMISA(hi.info.MISA, 'u') {
		hi.info.UXLEN, err = dbg.getUXLEN(hi)
		if err != nil {
			return err
		}
	}

	// do we have hypervisor mode?
	if rv.CheckExtMISA(hi.info.MISA, 'h') {
		hi.info.HXLEN, err = dbg.getHXLEN(hi)
		if err != nil {
			return err
		}
	}

	// get the DXLEN value
	hi.info.DXLEN, err = dbg.getDXLEN(hi)
	if err != nil {
		return err
	}

	// get the FLEN value
	hi.info.FLEN, err = dbg.getFLEN()
	if err != nil {
		// ignore errors - we probably don't have floating point support.
		log.Info.Printf("hart%d: %v", hi.info.ID, err)
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
	mhartid, err := dbg.RdCSR(rv.MHARTID, 0)
	if err != nil {
		return err
	}
	hi.info.MHARTID = uint(mhartid)
	log.Info.Printf("hart%d: MHARTID %d", hi.info.ID, hi.info.MHARTID)

	/*
		// get hartinfo parameters
		x, err = dbg.rdDmi(hartinfo)
		if err != nil {
			return err
		}
		hi.nscratch = util.Bits(uint(x), 23, 20)
		hi.datasize = util.Bits(uint(x), 15, 12)
		hi.dataaccess = util.Bit(uint(x), 16)
		hi.dataaddr = util.Bits(uint(x), 11, 0)

		log.Info.Printf("hart%d: nscratch %d words", hi.info.ID, hi.nscratch)
		log.Info.Printf("hart%d: datasize %d %s", hi.info.ID, hi.datasize, []string{"csr", "words"}[hi.dataaccess])
		log.Info.Printf("hart%d: dataaccess %s(%d)", hi.info.ID, []string{"csr", "memory"}[hi.dataaccess], hi.dataaccess)
		log.Info.Printf("hart%d: dataaddr 0x%x", hi.info.ID, hi.dataaddr)
	*/

	// Now that we have the register lengths we can create the per-hart CSR decodes.
	hi.info.NewCsr().Setup()

	// Using the MISA extension bits setup the disassembler.
	hi.info.ISA, err = rvda.New(hi.info.MXLEN, hi.info.MISA)
	if err != nil {
		return err
	}
	log.Info.Printf("hart%d: disassembler ISA %s", hi.info.ID, hi.info.ISA)

	if !wasHalted {
		// resume the hart
		_, err := dbg.resume()
		if err != nil {
			return err
		}
		hi.info.State = rv.Running
	}

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
