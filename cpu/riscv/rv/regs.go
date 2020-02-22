//-----------------------------------------------------------------------------
/*

RISC-V Registers

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// Register numbers.
const (
	RegZero = iota // 0: zero
	RegRa          // 1: return address
	RegSp          // 2: stack pointer
	RegGp          // 3: global pointer
	RegTp          // 4: thread pointer
	RegT0          // 5:
	RegT1          // 6:
	RegT2          // 7:
	RegS0          // 8: frame pointer
	RegS1          // 9:
	RegA0          // 10: syscall 0
	RegA1          // 11: syscall 1
	RegA2          // 12: syscall 2
	RegA3          // 13: syscall 3
	RegA4          // 14: syscall 4
	RegA5          // 15: syscall 5
	RegA6          // 16: syscall 6
	RegA7          // 17: syscall 7
	RegS2          // 18:
	RegS3          // 19:
	RegS4          // 20:
	RegS5          // 21:
	RegS6          // 22:
	RegS7          // 23:
	RegS8          // 24:
	RegS9          // 25:
	RegS10         // 26:
	RegS11         // 27:
	RegT3          // 28:
	RegT4          // 29:
	RegT5          // 30:
	RegT6          // 31:
)

//-----------------------------------------------------------------------------
