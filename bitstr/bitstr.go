//-----------------------------------------------------------------------------
/*

Bit String

A package to operate on bit strings of arbitrary length.

Notes:

The least signifcant bit of the bit string is bit 0.
The transmit order is right to left (as per the string representation).
In the bitstream "head" bits are transmitted before "tail" bits.

That is:

[tail]...[body]...[head] <- Tx First

The bits are stored in bitSets (up to 64 bits) and a slice of bitSets forms
the bit string. The bitSets are stored in the slice head first. A head
operation works on the start of the slice, while a tail operation works on
the end of the slice.

*/
//-----------------------------------------------------------------------------

package bitstr

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

//-----------------------------------------------------------------------------

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

const setSize = 64
const zeroes = uint64(0)
const ones = uint64((1 << 64) - 1)

// bitSet stores 0 to 64 bits.
type bitSet struct {
	val uint64 // bits in this set
	n   int    // number of bits in this set
}

// newBitSet returns a bit set of 0 to 64 bits.
func newBitSet(val uint64, n int) bitSet {
	if n > setSize || n < 0 {
		panic("")
	}
	if n < 64 {
		val &= uint64((1 << n) - 1)
	}
	return bitSet{
		val: val,
		n:   n,
	}
}

func (bs *bitSet) dropHead(n int) {
	if n > setSize || n < 0 {
		panic("")
	}
	bs.val >>= n
	bs.n -= n
}

func (bs *bitSet) dropTail(n int) {
	if n > setSize || n < 0 {
		panic("")
	}
	bs.val &= uint64((1 << (bs.n - n)) - 1)
	bs.n -= n
}

func (bs *bitSet) String() string {
	if bs.n == 0 {
		return ""
	}
	fmtX := fmt.Sprintf("%%0%db", bs.n)
	return fmt.Sprintf(fmtX, bs.val)
}

//-----------------------------------------------------------------------------

// BitString is a bit string of arbitrary length.
type BitString struct {
	set []bitSet // bit sets
	n   int
}

// NewBitString returns a new 0 length bitstring.
func NewBitString() *BitString {
	return &BitString{}
}

// tail adds a bit set to the tail of the bit string.
func (b *BitString) tail(bs bitSet) *BitString {
	b.set = append(b.set, bs)
	b.n += bs.n
	return b
}

// Tail0 adds zero bits to the tail of the bit string.
func (b *BitString) Tail0(n int) *BitString {
	for n > 0 {
		k := min(n, setSize)
		b.tail(newBitSet(zeroes, k))
		n -= k
	}
	return b
}

// Tail1 adds one bits to the tail of the bit string.
func (b *BitString) Tail1(n int) *BitString {
	for n > 0 {
		k := min(n, setSize)
		b.tail(newBitSet(ones, k))
		n -= k
	}
	return b
}

// DropHead removes n bits from the head of the bit string.
func (b *BitString) DropHead(n int) *BitString {
	if n >= b.n {
		b.n = 0
		b.set = []bitSet{}
		return b
	}
	b.n -= n
	for i := range b.set {
		bs := &b.set[i]
		k := min(n, bs.n)
		bs.dropHead(k)
		n -= k
		if n == 0 {
			break
		}
	}
	return b
}

// DropTail removes n bits from the tail of the bit string.
func (b *BitString) DropTail(n int) *BitString {
	if n >= b.n {
		b.set = []bitSet{}
		b.n = 0
		return b
	}
	b.n -= n
	for i := len(b.set) - 1; i >= 0; i-- {
		bs := &b.set[i]
		k := min(n, bs.n)
		bs.dropTail(k)
		n -= k
		if n == 0 {
			break
		}
	}
	return b
}

// Copy returns a new copy of the bit string.
func (b *BitString) Copy() *BitString {
	x := NewBitString()
	for i := range b.set {
		x.tail(b.set[i])
	}
	return x
}

// Head adds a bit string to the head of a bit string.
func (b *BitString) Head(a *BitString) *BitString {
	return nil
}

// Tail adds a bit string to the tail of a bit string.
func (b *BitString) Tail(a *BitString) *BitString {
	return nil
}

// BitString returns a 1/0 string for the bit string.
func (b *BitString) BitString() string {
	s := []string{}
	for i := len(b.set) - 1; i >= 0; i-- {
		if b.set[i].n > 0 {
			s = append(s, b.set[i].String())
		}
	}
	return strings.Join(s, "")
}

func (b *BitString) String() string {
	return fmt.Sprintf("(%d) %s", b.n, b.BitString())
}

//-----------------------------------------------------------------------------

// Null returns an empty bit string.
func Null() *BitString {
	return NewBitString()
}

// Ones returns a bit string with n-one bits.
func Ones(n int) *BitString {
	return NewBitString().Tail1(n)
}

// Zeroes returns a bit string with n-zero bits.
func Zeroes(n int) *BitString {
	return NewBitString().Tail0(n)
}

// FromString returns a bit string from a 1/0 string.
func FromString(s string) (*BitString, error) {
	n := len(s)
	b := NewBitString()
	for n > 0 {
		j := len(s)
		k := min(n, setSize)
		x, err := strconv.ParseUint(s[j-k:j], 2, 64)
		if err != nil {
			return nil, err
		}
		b.tail(newBitSet(x, k))
		s = s[0 : j-k]
		n -= k
	}
	return b, nil
}

// Random returns a random bit string of n bits.
func Random(n int) *BitString {
	b := NewBitString()
	for n > 0 {
		k := min(n, setSize)
		b.tail(newBitSet(rand.Uint64(), k))
		n -= k
	}
	return b
}

//-----------------------------------------------------------------------------
