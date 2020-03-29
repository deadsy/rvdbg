//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11 Debug RAM Cache Functions

*/
//-----------------------------------------------------------------------------

package rv11

//-----------------------------------------------------------------------------

type cacheEntry struct {
	data  uint // data for a debug ram word
	valid bool // does the cache data match the debug ram?
	dirty bool // does the cache data need to be written to the debug ram?
}

type ramCache struct {
	dbg   *Debug       // pointer back to parent debugger
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
	cache.allDirty()
	err := cache.flush()
	return cache, err
}

// allDirty marks all cache entries as dirty.
func (cache *ramCache) allDirty() {
	for i := range cache.entry {
		cache.entry[i].dirty = true
	}
}

//-----------------------------------------------------------------------------

// flushOps returns dbus write operations for dirty cache entries.
func (cache *ramCache) flushOps() []dbusOp {
	op := []dbusOp{}
	for i := range cache.entry {
		e := &cache.entry[i]
		if e.dirty {
			op = append(op, dbusWr(cache.base+uint(i), e.data))
		}
	}
	return op
}

// flush dirty cache entries to the debug target.
func (cache *ramCache) flush() error {
	ops := cache.flushOps()
	ops = append(ops, dbusEnd())
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
