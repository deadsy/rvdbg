//-----------------------------------------------------------------------------
/*

Memory Regions

Some device drivers (E.g. flash) break memory up into regions of different sizes.
This code provides a generic means of representing those regions.

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"

	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

func max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// Meta contains target/vendor specific region data.
type Meta interface {
	String() string
}

// Region is a contiguous region of memory.
type Region struct {
	name     string // name
	addrSize uint   // address size in bits
	addr     uint   // start address
	size     uint   // size in bytes
	end      uint   // end address of region
	meta     Meta   // target/vendor specific meta data
}

// NewRegion retuns a new memory region.
func NewRegion(name string, addr, size uint, meta Meta) *Region {
	return &Region{
		name:     name,
		addrSize: 32, // default to 32
		addr:     addr,
		size:     size,
		end:      addr + size - 1,
		meta:     meta,
	}
}

// SetAddrSize sets the region address size in bits.
func (r *Region) SetAddrSize(bits uint) {
	r.addrSize = bits
}

// Overlaps returns true if the regions overlap.
func (r *Region) Overlaps(x *Region) bool {
	return max(r.addr, x.addr) <= min(r.end, x.end)
}

// ColString returns a 4 string description of the memory region.
func (r *Region) ColString() []string {
	fmtAddr := util.UintFormat(r.addrSize)
	addrStr := fmt.Sprintf("%s %s", fmt.Sprintf(fmtAddr, r.addr), fmt.Sprintf(fmtAddr, r.end))
	metaStr := ""
	if r.meta != nil {
		metaStr = r.meta.String()
	}
	return []string{r.name, addrStr, util.MemSize(r.size), metaStr}
}

//-----------------------------------------------------------------------------
