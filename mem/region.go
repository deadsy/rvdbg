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
	"strings"

	cli "github.com/deadsy/go-cli"
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
	Name     string // name
	Addr     uint   // start address
	Size     uint   // size in bytes
	addrSize uint   // address size in bits
	meta     Meta   // target/vendor specific meta data
}

// NewRegion returns a new memory region.
func NewRegion(name string, addr, size uint, meta Meta) *Region {
	return &Region{
		Name:     name,
		Addr:     addr,
		Size:     size,
		addrSize: 32, // default to 32
		meta:     meta,
	}
}

// SetAddrSize sets the region address size in bits.
func (r *Region) SetAddrSize(bits uint) {
	r.addrSize = bits
}

// SetSize sets the region size in bytes.
func (r *Region) SetSize(size uint) {
	r.Size = size
}

// Overlaps returns true if the regions overlap.
func (r *Region) Overlaps(x *Region) bool {
	rEnd := r.Addr + r.Size - 1
	xEnd := x.Addr + x.Size - 1
	return max(r.Addr, x.Addr) <= min(rEnd, xEnd)
}

// Align32 adjusts a region for 32-bit alignment.
func (r *Region) Align32() {
	end := r.Addr + r.Size
	// round down start address to 32-bit byte boundary
	r.Addr &= ^uint(3)
	// round up end address to 32-bit byte boundary
	end += 3
	end &= ^uint(3)
	// adjust size
	r.Size = end - r.Addr
}

// ColString returns a 4 string description of the memory region.
func (r *Region) ColString() []string {
	fmtAddr := util.UintFormat(r.addrSize)
	end := r.Addr + r.Size - 1
	addrStr := fmt.Sprintf("%s %s", fmt.Sprintf(fmtAddr, r.Addr), fmt.Sprintf(fmtAddr, end))
	metaStr := ""
	if r.meta != nil {
		metaStr = r.meta.String()
	}
	return []string{r.Name, addrStr, util.MemSize(r.Size), metaStr}
}

func (r *Region) String() string {
	s := r.ColString()
	return strings.Join(s, " ")
}

//-----------------------------------------------------------------------------
// argument processing

// RegionDriver has the methods to process the command line arguments for memory regions.
type RegionDriver interface {
	GetDefaultRegion() *Region        // get a default region
	GetAddressSize() uint             // get address size in bits
	LookupSymbol(name string) *Region // lookup the address of a symbol
}

// FileRegionArg converts filename and memory region arguments.
func FileRegionArg(drv RegionDriver, args []string) (string, *Region, error) {
	err := cli.CheckArgc(args, []int{2, 3})
	if err != nil {
		return "", nil, err
	}

	// args[0] is the filename
	name := args[0]

	// the remaining arguments define the memory region
	r, err := RegionArg(drv, args[1:])
	if err != nil {
		return "", nil, err
	}
	return name, r, nil
}

// RegionArg converts command line arguments to a memory region.
func RegionArg(drv RegionDriver, args []string) (*Region, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return nil, err
	}

	defRegion := drv.GetDefaultRegion()

	if len(args) == 0 {
		return defRegion, nil
	}

	// lookup the first argument as a symbol
	r := drv.LookupSymbol(args[0])
	if r != nil {
		if len(args) == 2 {
			// don't take the symbol size, use the argument
			n, err := cli.UintArg(args[1], [2]uint{1, 0x100000000}, 16)
			if err != nil {
				return nil, err
			}
			r.SetSize(n)
		}
		return r, nil
	}

	// get the address
	maxAddr := uint((1 << drv.GetAddressSize()) - 1)
	addr, err := cli.UintArg(args[0], [2]uint{0, maxAddr}, 16)
	if err != nil {
		return nil, err
	}

	if len(args) == 1 {
		return NewRegion("", addr, defRegion.Size, nil), nil
	}

	// get the size
	n, err := cli.UintArg(args[1], [2]uint{1, 0x100000000}, 16)
	if err != nil {
		return nil, err
	}
	return NewRegion("", addr, n, nil), nil
}

//-----------------------------------------------------------------------------
