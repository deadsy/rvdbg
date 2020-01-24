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
	"strings"
)

//-----------------------------------------------------------------------------

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// bytesToUint64 converts an 8-byte slice to a uint64.
func bytesToUint64(b []byte) uint64 {
	_ = b[7] // bounds check hint
	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
}

// stringToUint64 converts a 0/1 string to a uint64.
func stringToUint64(s string) uint64 {
	var val uint64
	for i := range s {
		val <<= 1
		if s[i] == '1' {
			val |= 1
		}
	}
	return val
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
	if n < setSize {
		val &= uint64((1 << n) - 1)
	}
	return bitSet{
		val: val,
		n:   n,
	}
}

// dropHead drops n-bits from the head of a bit set.
func (bs *bitSet) dropHead(n int) {
	if n > setSize || n < 0 {
		panic("")
	}
	bs.val >>= n
	bs.n -= n
}

// dropTail drops n-bits from the tail of a bit set.
func (bs *bitSet) dropTail(n int) {
	if n > setSize || n < 0 {
		panic("")
	}
	bs.val &= uint64((1 << (bs.n - n)) - 1)
	bs.n -= n
}

// genByte generates an 8-byte slice from an input bit set and stores left-over bits.
func (bs *bitSet) genByte(in bitSet) []byte {
	if in.n == 0 {
		return nil
	}
	if bs.n+in.n >= setSize {
		// generate bytes
		buf := make([]byte, 8)
		val := bs.val | (in.val << bs.n)
		for i := range buf {
			buf[i] = byte(val >> (i * 8))
		}
		// store the left over bits
		k := setSize - bs.n
		bs.val = in.val >> k
		bs.n = in.n - k
		return buf
	}
	// store the input bits
	bs.val |= (in.val << bs.n)
	bs.n += in.n
	return nil
}

// flushByte flushes out any remaining bytes in the bit set.
func (bs *bitSet) flushByte() []byte {
	if bs.n == 0 {
		return nil
	}
	buf := make([]byte, (bs.n+7)>>3)
	for i := range buf {
		buf[i] = byte(bs.val >> (i * 8))
	}
	bs.val = 0
	bs.n = 0
	return buf
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
	if bs.n > 0 {
		b.set = append(b.set, bs)
		b.n += bs.n
	}
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
	b.set = append(a.set, b.set...)
	b.n += a.n
	return b
}

// Tail adds a bit string to the tail of a bit string.
func (b *BitString) Tail(a *BitString) *BitString {
	b.set = append(b.set, a.set...)
	b.n += a.n
	return b
}

// GetBytes returns a byte slice for the bit string.
func (b *BitString) GetBytes() []byte {
	buf := make([]byte, 0, (b.n+7)>>3)
	state := &bitSet{}
	for i := range b.set {
		buf = append(buf, state.genByte(b.set[i])...)
	}
	buf = append(buf, state.flushByte()...)
	return buf
}

// Length returns the length of the bit string.
func (b *BitString) Length() int {
	return b.n
}

// Split splits a bit string into []uint using the number of bits in the input slice.
func (b *BitString) Split(in []int) []uint {
	x := make([]uint, len(in))
	// TODO
	return x
}

func (b *BitString) String() string {
	s := []string{}
	for i := len(b.set) - 1; i >= 0; i-- {
		if b.set[i].n > 0 {
			s = append(s, b.set[i].String())
		}
	}
	return strings.Join(s, "")
}

// LengthBits returns a length/bits string for the bit string.
func (b *BitString) LengthBits() string {
	return fmt.Sprintf("(%d) %s", b.n, b.String())
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

// FromBytes returns a bit string from a byte slice.
func FromBytes(s []byte, n int) *BitString {
	// sanity check
	k := len(s)
	if n > k*8 {
		panic("")
	}
	b := NewBitString()
	i := 0
	// 8 bytes at a time
	for n >= setSize {
		val := bytesToUint64(s[i : i+8])
		b.tail(newBitSet(val, setSize))
		i += 8
		k -= 8
		n -= setSize
	}
	// left over bits
	if k > 0 {
		var val uint64
		for j := 0; j < k; j++ {
			val |= uint64(s[i+j]) << (j * 8)
		}
		b.tail(newBitSet(val, n))
	}
	return b
}

// FromString returns a bit string from a 1/0 string.
func FromString(s string) *BitString {
	n := len(s)
	b := NewBitString()
	for n > 0 {
		k := min(n, setSize)
		val := stringToUint64(s[n-k : n])
		b.tail(newBitSet(val, k))
		n -= k
	}
	return b
}

// Random returns a random bit string with n bits.
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
