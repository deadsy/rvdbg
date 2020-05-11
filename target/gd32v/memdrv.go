//-----------------------------------------------------------------------------
/*

Memory Driver

This code implements the mem.Driver interface.

*/
//-----------------------------------------------------------------------------

package gd32v

import (
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/mem"
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

// GetDefaultRegion returns a default memory region.
func (m *memDriver) GetDefaultRegion() *mem.Region {
	return mem.NewRegion("", 0, 0x100, nil)
}

// LookupSymbol returns an address and size for a symbol.
func (m *memDriver) LookupSymbol(name string) *mem.Region {
	p, err := m.dev.GetPeripheral(name)
	if err != nil {
		return nil
	}
	return mem.NewRegion(name, p.Addr, p.Size, nil)
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
