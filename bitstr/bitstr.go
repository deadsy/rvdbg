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

*/
//-----------------------------------------------------------------------------

package bitstr

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

func min(a, b uint) uint {
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
	n   uint   // number of bits in this set
}

// newBitSet returns a bit set of 0 to 64 bits.
func newBitSet(val uint64, n uint) bitSet {
	if n > setSize {
		panic("n > setSize")
	}
	if n < 64 {
		val &= uint64((1 << n) - 1)
	}
	return bitSet{
		val: val,
		n:   n,
	}
}

func (bs *bitSet) dropHead(n uint) {
	if n > setSize {
		panic("n > setSize")
	}
	bs.val >>= n
	bs.n -= n
}

func (bs *bitSet) dropTail(n uint) {
	if n > setSize {
		panic("n > setSize")
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
	n   uint
}

// NewBitString returns a new 0 length bitstring.
func NewBitString() *BitString {
	return &BitString{}
}

// Tail0 adds zero bits to the tail of the bit string.
func (b *BitString) Tail0(n uint) *BitString {
	for n > 0 {
		l := min(n, setSize)
		b.set = append(b.set, newBitSet(zeroes, l))
		n -= l
	}
	b.n += n
	return b
}

// Tail1 adds one bits to the tail of the bit string.
func (b *BitString) Tail1(n uint) *BitString {
	for n > 0 {
		l := min(n, setSize)
		b.set = append(b.set, newBitSet(ones, l))
		n -= l
	}
	b.n += n
	return b
}

// DropHead removes n bits from the head of the bit string.
func (b *BitString) DropHead(n uint) *BitString {
	if n >= b.n {
		b.n = 0
		b.set = []bitSet{}
		return b
	}
	b.n -= n
	for i := range b.set {
		bs := &b.set[i]
		l := min(n, bs.n)
		bs.dropHead(l)
		n -= l
		if n == 0 {
			break
		}
	}
	return b
}

// DropTail removes n bits from the tail of the bit string.
func (b *BitString) DropTail(n uint) *BitString {
	if n >= b.n {
		b.set = []bitSet{}
		b.n = 0
		return b
	}
	b.n -= n
	for i := len(b.set) - 1; i >= 0; i-- {
		bs := &b.set[i]
		l := min(n, bs.n)
		bs.dropTail(l)
		n -= l
		if n == 0 {
			break
		}
	}
	return b
}

func (b *BitString) String() string {
	s := []string{}
	for i := len(b.set) - 1; i >= 0; i-- {
		ss := b.set[i].String()
		if ss != "" {
			s = append(s, ss)
		}
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------

// Null returns an empty bit string.
func Null() *BitString {
	return NewBitString()
}

// Ones returns a bit string with n-one bits.
func Ones(n uint) *BitString {
	return NewBitString().Tail1(n)
}

// Zeroes returns a bit string with n-zero bits.
func Zeroes(n uint) *BitString {
	return NewBitString().Tail0(n)
}

//-----------------------------------------------------------------------------
