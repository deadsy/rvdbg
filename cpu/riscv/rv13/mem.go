//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.13 Memory Operations

*/
//-----------------------------------------------------------------------------

package rv13

//-----------------------------------------------------------------------------
// read memory

func (dbg *Debug) RdMem8(addr, n uint) ([]uint8, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem8(addr, n)
}

func (dbg *Debug) RdMem16(addr, n uint) ([]uint16, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem16(addr, n)
}

func (dbg *Debug) RdMem32(addr, n uint) ([]uint32, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem32(addr, n)
}

func (dbg *Debug) RdMem64(addr, n uint) ([]uint64, error) {
	if n == 0 {
		return nil, nil
	}
	return dbg.hart[dbg.hartid].rdMem64(addr, n)
}

//-----------------------------------------------------------------------------
// write memory

func (dbg *Debug) WrMem8(addr uint, val []uint8) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem8(addr, val)
}

func (dbg *Debug) WrMem16(addr uint, val []uint16) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem16(addr, val)
}

func (dbg *Debug) WrMem32(addr uint, val []uint32) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem32(addr, val)
}

func (dbg *Debug) WrMem64(addr uint, val []uint64) error {
	if len(val) == 0 {
		return nil
	}
	return dbg.hart[dbg.hartid].wrMem64(addr, val)
}

//-----------------------------------------------------------------------------
