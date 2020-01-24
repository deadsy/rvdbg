//-----------------------------------------------------------------------------
/*

Bit string test functions.

*/
//-----------------------------------------------------------------------------

package bitstr

import (
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

func Test_BitString(t *testing.T) {

	var a, b *BitString

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

	b = NewBitString().Tail0(271)
	if b.BitString() != repeatRune('0', 271) {
		t.Error("FAIL")
	}

	b = NewBitString().Tail1(1490)
	if b.BitString() != repeatRune('1', 1490) {
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

	b, _ = FromString("11101")
	if b.String() != "(5) 11101" {
		t.Error("FAIL")
	}

	b, _ = FromString("011111")
	if b.String() != "(6) 011111" {
		t.Error("FAIL")
	}

	b = Random(2017)
	a, _ = FromString(b.BitString())
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	b = Random(1987)
	a = b.Copy()
	if a.String() != b.String() {
		t.Error("FAIL")
	}

	/*


	   x = bits.from_random(10)
	   y = bits.from_random(10)
	   z = x.copy().head(y)
	   self.assertEqual(x.bit_str() + y.bit_str(), z.bit_str())


	   x = bits.from_random(10)
	   y = bits.from_random(10)
	   z = x.copy().tail(y)
	   self.assertEqual(y.bit_str() + x.bit_str(), z.bit_str())


	   x = bits.from_tuple('11111111')
	   y = x.get_bytes()
	   self.assertEqual(len(y), 1)
	   self.assertEqual(y[0], 255)


	   x = bits.from_tuple('11011110101011011011111011101111')
	   y = x.get_bytes()
	   self.assertEqual(len(y), 4)
	   self.assertEqual(y[0], 0xef)
	   self.assertEqual(y[1], 0xbe)
	   self.assertEqual(y[2], 0xad)
	   self.assertEqual(y[3], 0xde)


	   x = bits.bits().set_bytes((0xff,), 7)
	   self.assertEqual(str(x), '(7) 1111111')
	   x = bits.bits().set_bytes((0x1,), 7)
	   self.assertEqual(str(x), '(7) 0000001')
	   x = bits.bits().set_bytes((64,), 7)
	   self.assertEqual(str(x), '(7) 1000000')


	   str0 = '11011110101011011011111011101111'
	   x = bits.from_tuple(str0)
	   y = bits.bits().set_bytes(x.get_bytes(), len(str0))
	   self.assertEqual(y.bit_str(), str0)


	*/

}

//-----------------------------------------------------------------------------
