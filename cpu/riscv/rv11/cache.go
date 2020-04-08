//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 Debug RAM Cache Functions

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
)

//-----------------------------------------------------------------------------

type cacheEntry struct {
	data uint32 // data for a debug ram word
	wr   bool   // write the data to the debug ram
	rd   bool   // read this data from the debug ram
}

type ramCache struct {
	dbg   *Debug       // pointer back to parent debugger
	isa   *rvda.ISA    // RV32i disassembler
	base  uint         // base address of debug ram
	entry []cacheEntry // cache entries
}

// newCache returns a debug RAM cache.
func (dbg *Debug) newCache(base, entries uint) (*ramCache, error) {
	cache := &ramCache{
		dbg:   dbg,
		base:  base,
		entry: make([]cacheEntry, entries),
	}
	// add a disassembler for decoding instructions in the cache
	isa, err := rvda.New(64, rvda.ExtI|rvda.ExtF|rvda.ExtD)
	if err != nil {
		return nil, err
	}
	cache.isa = isa
	// sync the cache and target debug ram
	err = cache.reset()
	if err != nil {
		return nil, err
	}
	return cache, nil
}

//-----------------------------------------------------------------------------

// clearAll clears wr/rd flags from all cache enties.
func (cache *ramCache) clearAll() {
	for i := range cache.entry {
		cache.entry[i].wr = false
		cache.entry[i].rd = false
	}
}

//-----------------------------------------------------------------------------

// entryString returns the display string for a cache entry.
func (cache *ramCache) entryString(idx int) string {
	e := &cache.entry[idx]
	flags := [2]rune{'.', '.'}
	if e.wr {
		flags[0] = 'w'
	}
	if e.rd {
		flags[1] = 'r'
	}
	addr := cache.base + (4 * uint(idx))
	da := cache.isa.Disassemble(addr, uint(e.data))
	return fmt.Sprintf("%2d: %03x %08x %s %s", idx, addr, e.data, string(flags[:]), da.Assembly)
}

func (cache *ramCache) String() string {
	s := []string{}
	for i := range cache.entry {
		s = append(s, cache.entryString(i))
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------

// ops returns dbus write/read operations for cache entries.
func (cache *ramCache) ops(exec bool) []dbusOp {
	op := []dbusOp{}
	// writes
	last := -1
	for i := range cache.entry {
		e := &cache.entry[i]
		if e.wr {
			op = append(op, dbusWr(uint(i), uint(e.data)))
			last = len(op) - 1
		}
	}
	// mark the last write to execute instructions
	if exec {
		if last >= 0 {
			// set the interrupt for the last write operation
			op[last] = op[last].setInterrupt()
		} else {
			// no debug ram writes, set the interrupt in dmcontrol
			op = append(op, dbusWr(dmcontrol, haltNotification|debugInterrupt))
		}
	}
	// reads
	for i := range cache.entry {
		e := &cache.entry[i]
		if e.rd {
			op = append(op, dbusRd(uint(i)))
		}
	}
	// final operation to read the last value/status
	op = append(op, dbusEnd())
	return op
}

//-----------------------------------------------------------------------------

// wr writes an instruction word to the cache.
func (cache *ramCache) wr32(i int, data uint32) {
	e := &cache.entry[i]
	if e.data != data {
		e.data = data
		e.wr = true
	}
}

// wr64 writes a 64-bit word to the cache.
func (cache *ramCache) wr64(i int, data uint64) {
	cache.wr32(i, uint32(data))
	cache.wr32(i+1, uint32(data>>32))
}

// wrResume writes a "jal debugRomResume" to the cache.
func (cache *ramCache) wrResume(i int) {
	cache.wr32(i, rv.InsJAL(rv.RegZero, uint(debugRomResume-(debugRamStart+(4*i)))))
}

// rv64Addr adds cache code to setup a 64-bit address in S0.
func (cache *ramCache) rv64Addr(addr uint) {
	if addr&util.Upper32 == 0 {
		// use slots 0,4
		cache.wr32(0, rv.InsLWU(rv.RegS0, ramAddr(4), rv.RegZero))
		cache.wr32(4, uint32(addr))
	} else {
		// use slots 0,4,5
		cache.wr32(0, rv.InsLD(rv.RegS0, ramAddr(4), rv.RegZero))
		cache.wr64(4, uint64(addr))
	}
}

//-----------------------------------------------------------------------------

// rd32 reads a 32-bit value from the cache.
func (cache *ramCache) rd32(i int) uint32 {
	return cache.entry[i].data
}

// rd64 reads a 64-bit value from the cache.
func (cache *ramCache) rd64(i int) uint64 {
	l := uint64(cache.entry[i].data)
	h := uint64(cache.entry[i+1].data)
	return (h << 32) | l
}

// read sets the read flag for the cache entry.
func (cache *ramCache) read(i uint) {
	cache.entry[i].rd = true
}

//-----------------------------------------------------------------------------

// flush runs the current set of write/read operations in the cache.
func (cache *ramCache) flush(exec bool) error {
	// The last word of debug ram indicates exceptions.
	ex := uint(len(cache.entry) - 1)
	if exec {
		cache.read(ex)
	}
	// run the operations
	data, err := cache.dbg.dbusOps(cache.ops(exec))
	if err != nil {
		return err
	}
	// data[] has the read data
	k := 0
	for i := range cache.entry {
		e := &cache.entry[i]
		if e.rd {
			e.data = uint32(data[k])
			k++
		}
	}
	// clear all cache flags
	cache.clearAll()
	// check for exceptions
	if exec && cache.entry[ex].data != 0 {
		return errors.New("exception")
	}
	return nil
}

//-----------------------------------------------------------------------------

// reset and sync the cache and debug ram state.
func (cache *ramCache) reset() error {
	for i := range cache.entry {
		e := &cache.entry[i]
		e.data = 0xdeadbeef
		e.rd = false
		e.wr = true
	}
	return cache.flush(false)
}

//-----------------------------------------------------------------------------
