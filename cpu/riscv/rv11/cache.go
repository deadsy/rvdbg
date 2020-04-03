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
	data  uint32 // data for a debug ram word
	valid bool   // does the cache data match the debug ram?
	dirty bool   // does the cache data need to be written to the debug ram?
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
	isa, err := rvda.New(32, rvda.ExtI)
	if err != nil {
		return nil, err
	}
	cache.isa = isa
	// initialise the cache and debug ram
	cache.allDirty()
	err = cache.flush(false)
	return cache, err
}

// allDirty marks all cache entries as dirty.
func (cache *ramCache) allDirty() {
	for i := range cache.entry {
		cache.entry[i].dirty = true
	}
}

//-----------------------------------------------------------------------------

func (cache *ramCache) EntryString(idx int) string {
	e := &cache.entry[idx]

	flags := []rune{}
	if e.dirty {
		flags = append(flags, 'd')
	}
	if e.valid {
		flags = append(flags, 'v')
	}

	addr := cache.base + (4 * uint(idx))
	da := cache.isa.Disassemble(addr, uint(e.data))

	return fmt.Sprintf("%03x %09x %s %s", addr, e.data, string(flags), da.Assembly)
}

func (cache *ramCache) String() string {
	s := []string{}
	for i := range cache.entry {
		s = append(s, cache.EntryString(i))
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------

// wrOps returns dbus write operations for dirty cache entries.
func (cache *ramCache) wrOps() []dbusOp {
	op := []dbusOp{}
	for i := range cache.entry {
		e := &cache.entry[i]
		if e.dirty {
			op = append(op, dbusWr(uint(i), uint(e.data)))
		}
	}
	return op
}

// rdOps returns dbus read operations for invalid cache entries.
func (cache *ramCache) rdOps() []dbusOp {
	op := []dbusOp{}
	for i := range cache.entry {
		e := &cache.entry[i]
		if !e.valid {
			op = append(op, dbusRd(uint(i)))
		}
	}
	return op
}

//-----------------------------------------------------------------------------

// wr writes an instruction word to the cache.
func (cache *ramCache) wr(i int, data uint32) {
	e := &cache.entry[i]
	if e.data != data {
		e.data = data
		e.dirty = true
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

// invalid marks a cache entry as invalid.
func (cache *ramCache) invalid(i uint) {
	cache.entry[i].valid = false
}

// validate reads invalid cache entries from the debug target.
func (cache *ramCache) validate() error {
	ops := cache.rdOps()
	ops = append(ops, dbusEnd())
	data, err := cache.dbg.dbusOps(ops)
	if err != nil {
		return err
	}
	// previously invalid entries are now valid and clean
	k := 0
	for i := range cache.entry {
		e := &cache.entry[i]
		if !e.valid {
			e.data = uint32(data[k])
			k++
			e.dirty = false
			e.valid = true
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

// flush dirty cache entries to the debug target.
func (cache *ramCache) flush(exec bool) error {
	ops := cache.wrOps()
	if exec {
		if len(ops) >= 1 {
			// set the interrupt for the last operation
			ops[len(ops)-1] = ops[len(ops)-1].setInterrupt()
		} else {
			// no debug ram writes, set the interrupt in dmcontrol
			ops = append(ops, dbusWr(dmcontrol, haltNotification|debugInterrupt))
		}
	}
	ops = append(ops, dbusEnd())
	// run the operations
	_, err := cache.dbg.dbusOps(ops)
	// previously dirty entries are now clean and valid
	for i := range cache.entry {
		e := &cache.entry[i]
		if e.dirty {
			e.dirty = false
			e.valid = true
		}
	}
	return err
}

//-----------------------------------------------------------------------------
