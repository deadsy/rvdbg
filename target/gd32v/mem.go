//-----------------------------------------------------------------------------
/*

Memory Driver

This code implements the mem.Driver interface.

*/
//-----------------------------------------------------------------------------

package gd32v

import (
	"fmt"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/soc"
)

//-----------------------------------------------------------------------------

type memDriver struct {
	dbg rv.Debug
	dev *soc.Device
}

func newMemDriver(dbg rv.Debug, dev *soc.Device) *memDriver {
	return &memDriver{
		dbg: dbg,
		dev: dev,
	}
}

// GetAddressSize returns the address size in bits.
func (m *memDriver) GetAddressSize() uint {
	return m.dbg.GetAddressSize()
}

// LookupSymbol returns an address and size for a symbol.
func (m *memDriver) LookupSymbol(name string) (uint, uint, error) {
	p := m.dev.GetPeripheral(name)
	if p != nil {
		return p.Addr, p.Size, nil
	}
	return 0, 0, fmt.Errorf("symbol \"%s\" not found", name)
}

// RdMem reads n x width-bit values from memory.
func (m *memDriver) RdMem(width, addr, n uint) ([]uint, error) {
	return m.dbg.RdMem(width, addr, n)
}

// WrMem wirtes n x width-bit values to memory.
func (m *memDriver) WrMem(width, addr uint, val []uint) error {
	return m.dbg.WrMem(width, addr, val)
}

//-----------------------------------------------------------------------------
