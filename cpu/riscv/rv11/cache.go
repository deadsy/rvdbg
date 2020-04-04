//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 Debug RAM Cache Functions

*/
//-----------------------------------------------------------------------------

package rv11

import (
	"fmt"
	"strings"

	"github.com/deadsy/rvda"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
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
			last = i
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
func (cache *ramCache) wr(i int, data uint32) {
	e := &cache.entry[i]
	if e.data != data {
		e.data = data
		e.wr = true
	}
}

// wrResume writes a "jal debugRomResume" to the cache.
func (cache *ramCache) wrResume(i int) {
	cache.wr(i, rv.InsJAL(rv.RegZero, uint(debugRomResume-(debugRamStart+(4*i)))))
}

//-----------------------------------------------------------------------------

// rd reads a value from the cache.
func (cache *ramCache) rd(i int) uint32 {
	return cache.entry[i].data
}

// read sets the read flag for the cache entry.
func (cache *ramCache) read(i uint) {
	cache.entry[i].rd = true
}

//-----------------------------------------------------------------------------

// flush runs the current set of write/read operations in the cache.
func (cache *ramCache) flush(exec bool) error {
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
