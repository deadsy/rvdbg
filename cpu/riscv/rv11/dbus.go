//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11
Debug Bus Access

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"fmt"
	"math/rand"

	"github.com/deadsy/rvdbg/bitstr"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
	"github.com/deadsy/rvdbg/util/log"
)

//-----------------------------------------------------------------------------
// dbus registers (5-bit address)

const ram0 = 0x00
const ram1 = 0x01
const ram2 = 0x02
const ram3 = 0x03
const ram4 = 0x04
const ram5 = 0x05
const ram6 = 0x06
const ram7 = 0x07
const ram8 = 0x08
const ram9 = 0x09
const ram10 = 0x0a
const ram11 = 0x0b
const ram12 = 0x0c
const ram13 = 0x0d
const ram14 = 0x0e
const ram15 = 0x0f

const dmcontrol = 0x10
const dminfo = 0x11
const authdata0 = 0x12
const authdata1 = 0x13
const serdata = 0x14
const serstatus = 0x15
const sbaddress0 = 0x16
const sbaddress1 = 0x17
const sbdata0 = 0x18
const sbdata1 = 0x19
const haltsum = 0x1b

// upper bits of each debug ram value (also in dmcontrol)
const haltNotification = (1 << 32)
const debugInterrupt = (1 << 33)

//-----------------------------------------------------------------------------
// dbus address locations

const debugRomStart = 0x800
const debugRomResume = debugRomStart + 4
const debugRomException = debugRomStart + 8
const debugRamStart = 0x400

//-----------------------------------------------------------------------------

func newDBUS() *soc.Device {
	return &soc.Device{
		Name: "DBUS",
		Peripherals: []soc.Peripheral{
			{
				Name:  "DBUS",
				Descr: "DBUS Registers",
				Registers: []soc.Register{
					{Offset: ram0, Size: 32, Name: "ram0", Descr: "Debug RAM 0"},
					{Offset: ram1, Size: 32, Name: "ram1", Descr: "Debug RAM 1"},
					{Offset: ram2, Size: 32, Name: "ram2", Descr: "Debug RAM 2"},
					{Offset: ram3, Size: 32, Name: "ram3", Descr: "Debug RAM 3"},
					{Offset: ram4, Size: 32, Name: "ram4", Descr: "Debug RAM 4"},
					{Offset: ram5, Size: 32, Name: "ram5", Descr: "Debug RAM 5"},
					{Offset: ram6, Size: 32, Name: "ram6", Descr: "Debug RAM 6"},
					{Offset: ram7, Size: 32, Name: "ram7", Descr: "Debug RAM 7"},
					{Offset: ram8, Size: 32, Name: "ram8", Descr: "Debug RAM 8"},
					{Offset: ram9, Size: 32, Name: "ram9", Descr: "Debug RAM 9"},
					{Offset: ram10, Size: 32, Name: "ram10", Descr: "Debug RAM 10"},
					{Offset: ram11, Size: 32, Name: "ram11", Descr: "Debug RAM 11"},
					{Offset: ram12, Size: 32, Name: "ram12", Descr: "Debug RAM 12"},
					{Offset: ram13, Size: 32, Name: "ram13", Descr: "Debug RAM 13"},
					{Offset: ram14, Size: 32, Name: "ram14", Descr: "Debug RAM 14"},
					{Offset: ram15, Size: 32, Name: "ram15", Descr: "Debug RAM 15"},
					{Offset: dmcontrol,
						Name:  "dmcontrol",
						Descr: "Control",
						Fields: []soc.Field{
							{Name: "interrupt", Msb: 33, Lsb: 33},
							{Name: "haltnot", Msb: 32, Lsb: 32},
							{Name: "buserror", Msb: 21, Lsb: 19},
							{Name: "serial", Msb: 18, Lsb: 16},
							{Name: "autoincrement", Msb: 15, Lsb: 15},
							{Name: "access", Msb: 14, Lsb: 12},
							{Name: "hartid", Msb: 11, Lsb: 2},
							{Name: "ndreset", Msb: 1, Lsb: 1},
							{Name: "fullreset", Msb: 0, Lsb: 0},
						},
					},
					{Offset: dminfo,
						Name:  "dminfo",
						Descr: "Info",
						Fields: []soc.Field{
							{Name: "abussize", Msb: 31, Lsb: 25},
							{Name: "serialcount", Msb: 24, Lsb: 21},
							{Name: "access128", Msb: 20, Lsb: 20},
							{Name: "access64", Msb: 19, Lsb: 19},
							{Name: "access32", Msb: 18, Lsb: 18},
							{Name: "access16", Msb: 17, Lsb: 17},
							{Name: "access8", Msb: 16, Lsb: 16},
							{Name: "dramsize", Msb: 15, Lsb: 10},
							{Name: "haltsum", Msb: 9, Lsb: 9},
							{Name: "hiversion", Msb: 7, Lsb: 6},
							{Name: "authenticated", Msb: 5, Lsb: 5},
							{Name: "authbusy", Msb: 4, Lsb: 4},
							{Name: "authtype", Msb: 3, Lsb: 2},
							{Name: "loversion", Msb: 1, Lsb: 0},
						},
					},
					{Offset: authdata0, Name: "authdata0", Descr: "Authentication Data"},
					{Offset: authdata1, Name: "authdata1", Descr: "Authentication Data"},
					{Offset: serdata, Name: "serdata", Descr: "Serial Data"},
					{Offset: serstatus, Name: "serstatus", Descr: "Serial Status"},
					{Offset: sbaddress0, Name: "sbaddress0", Descr: "System Bus Address 31:0"},
					{Offset: sbaddress1, Name: "sbaddress1", Descr: "System Bus Address 63:32"},
					{Offset: sbdata0, Name: "sbdata0", Descr: "System Bus Data 31:0"},
					{Offset: sbdata1, Name: "sbdata1", Descr: "System Bus Data 63:32"},
					{Offset: haltsum, Name: "haltsum", Descr: "Halt Notification Summary"},
				},
			},
		},
	}
}

//-----------------------------------------------------------------------------
// hart selection

const hartid = uint(((1 << 10) - 1) << 2)

// setHartSelect sets the hartid field in a dmcontrol value.
func setHartSelect(x uint, id uint) uint {
	return (x & ^hartid) | ((id << 2) & hartid)
}

// getHartSelect gets the hart select value from a dmcontrol value.
func getHartSelect(x uint) uint {
	return util.Bits(x, 11, 2)
}

// selectHart sets the dmcontrol hartsel value.
func (dbg *Debug) selectHart(id int) error {
	x, err := dbg.rdDbus(dmcontrol)
	if err != nil {
		return err
	}
	x = setHartSelect(x, uint(id))
	return dbg.wrDbus(dmcontrol, x)
}

//-----------------------------------------------------------------------------

// dbus operations
const opIgnore = 0
const opRd = 1
const opWr = 2

// dbus errors
const opOk = 0
const opFail = 2
const opBusy = 3
const opMask = (1 << 2) - 1

type dbusOp uint

// dbusRd returns a dbus read operation.
func dbusRd(addr uint) dbusOp {
	return dbusOp((addr << 36) | opRd)
}

// dbusWr returns a dbus write operation.
func dbusWr(addr, data uint) dbusOp {
	return dbusOp((addr << 36) | (data << 2) | opWr)
}

// dbusEnd returns a dbus no-op, typically used to clock out a final data value.
func dbusEnd() dbusOp {
	return dbusOp(opIgnore)
}

// setInterrupt sets the interrupt and halt notification bits in a dbus operation.
func (x dbusOp) setInterrupt() dbusOp {
	return x | ((haltNotification | debugInterrupt) << 2)
}

// isRead returns true if this dbus operation is a read.
func (x dbusOp) isRead() bool {
	return (x & opMask) == opRd
}

// dbusOps runs a set of dbus operations and returns any read data.
func (dbg *Debug) dbusOps(ops []dbusOp) ([]uint, error) {
	data := []uint{}

	// select dbus
	err := dbg.wrIR(irDbus)
	if err != nil {
		return nil, err
	}

	read := false
	for i := 0; i < len(ops); i++ {
		op := ops[i]
		// run the operation
		tdo, err := dbg.dev.RdWrDR(bitstr.FromUint(uint(op), dbg.drDbusLength), dbg.idle)
		if err != nil {
			return nil, err
		}
		x := tdo.Split([]int{dbg.drDbusLength})[0]
		// check the result
		result := x & opMask
		if result != opOk {
			// clear error condition
			dbg.wrDtmcontrol(dbusreset)
			// re-select dbus
			dbg.wrIR(irDbus)
			if result == opBusy {
				// auto-adjust timing
				log.Info.Printf("increment idle timing %d->%d cycles", dbg.idle, dbg.idle+1)
				dbg.idle++
				if dbg.idle > jtag.MaxIdle {
					return nil, fmt.Errorf("dbus operation error %d", result)
				}
				// redo the operation
				i--
				continue
			} else {
				return nil, fmt.Errorf("dbus operation error %d", result)
			}
		}
		// get the read data
		if read {
			data = append(data, (x>>2)&util.Mask34)
		}
		// setup the next read
		read = op.isRead()
	}
	return data, nil
}

//-----------------------------------------------------------------------------
// dbus read/write

// rdDbus reads a dbus register.
func (dbg *Debug) rdDbus(addr uint) (uint, error) {
	ops := []dbusOp{
		dbusRd(addr),
		dbusEnd(),
	}
	data, err := dbg.dbusOps(ops)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

// wrDbus writes a dbus register.
func (dbg *Debug) wrDbus(addr, data uint) error {
	ops := []dbusOp{
		dbusWr(addr, data),
		dbusEnd(),
	}
	_, err := dbg.dbusOps(ops)
	return err
}

// rmwDbus read/modify/write a dbus register.
func (dbg *Debug) rmwDbus(addr, mask, bits uint) error {
	// read
	x, err := dbg.rdDbus(addr)
	if err != nil {
		return err
	}
	// modify
	x &= ^mask
	x |= bits
	// write
	return dbg.wrDbus(addr, x)
}

// setDbus sets bits in a dbus register.
func (dbg *Debug) setDbus(addr, bits uint) error {
	return dbg.rmwDbus(addr, bits, bits)
}

// clrDbus clears bits in a dbus register.
func (dbg *Debug) clrDbus(addr, bits uint) error {
	return dbg.rmwDbus(addr, bits, 0)
}

//-----------------------------------------------------------------------------
// decode/display dbus registers

type dbusDriver struct {
	dbg *Debug
}

func (drv *dbusDriver) GetAddressSize() uint {
	return drv.dbg.abits
}

func (drv *dbusDriver) GetRegisterSize(r *soc.Register) uint {
	return 34
}

func (drv *dbusDriver) Rd(width, addr uint) (uint, error) {
	return drv.dbg.rdDbus(addr)
}

func (dbg *Debug) dbusDump() (string, error) {
	p := dbg.dbusDevice.GetPeripheral("DBUS")
	drv := &dbusDriver{dbg}
	return p.Display(drv, nil, true), nil
}

//-----------------------------------------------------------------------------

// testBuffers tests dbus r/w buffers.
func (dbg *Debug) testBuffers(addr, n uint) error {

	// random write values
	wr := make([]uint, n)
	for i := range wr {
		wr[i] = uint(rand.Uint32())
	}

	// write to dbus registers
	for i := range wr {
		err := dbg.wrDbus(addr+uint(i), wr[i])
		if err != nil {
			return err
		}
	}

	// read back from dbus registers
	for i := range wr {
		x, err := dbg.rdDbus(addr + uint(i))
		if err != nil {
			return err
		}
		if x != (wr[i] & util.Mask32) {
			return fmt.Errorf("w/r mismatch at 0x%x", addr+uint(i))
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
