//-----------------------------------------------------------------------------
/*

Bit string test functions.

*/
//-----------------------------------------------------------------------------

package bitstr

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

//-----------------------------------------------------------------------------

func repeatRune(r rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = r
	}
	return string(s)
}

//-----------------------------------------------------------------------------

const testString0 = "11011110101011011011111011101111"

//-----------------------------------------------------------------------------

func Test_BitString(t *testing.T) {

	var a, b, c *BitString
	var k int

	b = NewBitString().Tail0(2)
	if b.BitString() != "00" {
		t.Error("FAIL")
	}

	b = NewBitString().Tail1(5)
	if b.BitString() != "11111" {
		t.Error("FAIL")
	}

	b = NewBitString().Tail0(2).Tail1(3)
	if b.BitString() != "11100" {
		t.Error("FAIL")
	}

	b = NewBitString().Tail1(5).Tail0(7)
	if b.BitString() != "000000011111" {
		t.Error("FAIL")
	}

	k = 271
	b = NewBitString().Tail0(k)
	if b.BitString() != repeatRune('0', k) {
		t.Error("FAIL")
	}

	k = 1490
	b = NewBitString().Tail1(k)
	if b.BitString() != repeatRune('1', k) {
		t.Error("FAIL")
	}

	b = Null()
	if b.BitString() != "" {
		t.Error("FAIL")
	}

	b = Ones(3)
	if b.String() != "(3) 111" {
		t.Error("FAIL")
	}

	b = Zeroes(7)
	if b.String() != "(7) 0000000" {
		t.Error("FAIL")
	}

	b = Zeroes(7).Tail1(4)
	if b.String() != "(11) 11110000000" {
		t.Error("FAIL")
	}

	b = Ones(7).Tail0(4)
	if b.String() != "(11) 00001111111" {
		t.Error("FAIL")
	}

	b = FromString("11101")
	if b.String() != "(5) 11101" {
		t.Error("FAIL")
	}

	b = FromString("011111")
	if b.String() != "(6) 011111" {
		t.Error("FAIL")
	}

	b = Random(2017)
	a = FromString(b.BitString())
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	b = Random(1987)
	a = b.Copy()
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	a = Random(73)
	b = Random(89)
	c = a.Copy().Head(b)
	if a.BitString()+b.BitString() != c.BitString() {
		t.Error("FAIL")
	}

	a = Random(26)
	b = Random(128)
	c = a.Copy().Tail(b)
	if b.BitString()+a.BitString() != c.BitString() {
		t.Error("FAIL")
	}

	a = FromString("")
	if !bytes.Equal(a.GetBytes(), []byte{}) {
		t.Error("FAIL")
	}

	a = FromString("11111111")
	if !bytes.Equal(a.GetBytes(), []byte{255}) {
		t.Error("FAIL")
	}

	a = FromString(testString0)
	if !bytes.Equal(a.GetBytes(), []byte{0xef, 0xbe, 0xad, 0xde}) {
		t.Error("FAIL")
	}

	a = FromBytes([]byte{0xff}, 7)
	if a.String() != "(7) 1111111" {
		t.Error("FAIL")
	}

	a = FromBytes([]byte{0x1}, 7)
	if a.String() != "(7) 0000001" {
		t.Error("FAIL")
	}

	a = FromBytes([]byte{64}, 7)
	if a.String() != "(7) 1000000" {
		t.Error("FAIL")
	}

	a = FromString(testString0)
	b = FromBytes(a.GetBytes(), len(testString0))
	if a.BitString() != testString0 {
		t.Error("FAIL")
	}

	k = 2057
	a = Random(k)
	b = FromBytes(a.GetBytes(), k)
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	// random tails
	rand.Seed(1)
	a = NewBitString()
	for i := 0; i < 500; i++ {
		a = a.Tail(Random(rand.Int() % 197))
	}
	b = FromBytes(a.GetBytes(), a.Length())
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	// random heads
	rand.Seed(1)
	a = NewBitString()
	for i := 0; i < 500; i++ {
		a = a.Head(Random(rand.Int() % 1709))
	}
	b = FromBytes(a.GetBytes(), a.Length())
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	// random tail/head
	rand.Seed(1)
	a = NewBitString()
	for i := 0; i < 500; i++ {
		a = a.Tail(Random(rand.Int() % 197))
	}
	for i := 0; i < 100; i++ {
		a = a.Head(Random(rand.Int() % 1709))
	}

	b = FromBytes(a.GetBytes(), a.Length())

	fmt.Printf("%d %d\n", a.Length(), b.Length())

	if a.String() != b.String() {
		t.Error("FAIL")
	}

	// random head/tail
	rand.Seed(1)
	a = NewBitString()
	for i := 0; i < 500; i++ {
		a = a.Head(Random(rand.Int() % 197))
	}
	for i := 0; i < 500; i++ {
		a = a.Tail(Random(rand.Int() % 1709))
	}
	b = FromBytes(a.GetBytes(), a.Length())
	if a.String() != b.String() {
		t.Error("FAIL")
	}

}

//-----------------------------------------------------------------------------
