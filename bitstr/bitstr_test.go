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

	b0 := NewBitString().Tail0(2)
	if b0.String() != "00" {
		t.Error("FAIL")
	}

	b1 := NewBitString().Tail1(5)
	if b1.String() != "11111" {
		t.Error("FAIL")
	}

	b0 = b0.Tail1(3)
	if b0.String() != "11100" {
		t.Error("FAIL")
	}

	b1 = b1.Tail0(7)
	if b1.String() != "000000011111" {
		t.Error("FAIL")
	}

	b0 = NewBitString().Tail0(271)
	if b0.String() != repeatRune('0', 271) {
		t.Error("FAIL")
	}

	b1 = NewBitString().Tail1(1490)
	if b1.String() != repeatRune('1', 1490) {
		t.Error("FAIL")
	}

}

//-----------------------------------------------------------------------------
