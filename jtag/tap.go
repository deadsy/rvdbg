//-----------------------------------------------------------------------------
/*

Precanned JTAG TAP State Transitions

*/
//-----------------------------------------------------------------------------

package jtag

import "github.com/deadsy/rvdbg/bitstr"

//-----------------------------------------------------------------------------

// ToIdle : any state -> run-test/idle
var ToIdle = bitstr.FromString("011111")

// IdleToIRshift : run-test/idle -> shift-ir
var IdleToIRshift = bitstr.FromString("0011")

// IdleToDRshift : run-test/idle -> shift-dr
var IdleToDRshift = bitstr.FromString("001")

// MaxIdle is the maximum number of additional TCK cycles we will stay
// in the run-test-idle state after scanning IR/DR.
const MaxIdle = 16

// ShiftToIdle : shift-x -> run-test/idle
var ShiftToIdle = [MaxIdle + 1]*bitstr.BitString{
	bitstr.FromString("011"),                 // + 0 cycles
	bitstr.FromString("0011"),                // + 1 cycles
	bitstr.FromString("00011"),               // + 2 cycles
	bitstr.FromString("000011"),              // + 3 cycles
	bitstr.FromString("0000011"),             // + 4 cycles
	bitstr.FromString("00000011"),            // + 5 cycles
	bitstr.FromString("000000011"),           // + 6 cycles
	bitstr.FromString("0000000011"),          // + 7 cycles
	bitstr.FromString("00000000011"),         // + 8 cycles
	bitstr.FromString("000000000011"),        // + 9 cycles
	bitstr.FromString("0000000000011"),       // + 10 cycles
	bitstr.FromString("00000000000011"),      // + 11 cycles
	bitstr.FromString("000000000000011"),     // + 12 cycles
	bitstr.FromString("0000000000000011"),    // + 13 cycles
	bitstr.FromString("00000000000000011"),   // + 14 cycles
	bitstr.FromString("000000000000000011"),  // + 15 cycles
	bitstr.FromString("0000000000000000011"), // + 16 cycles
}

//-----------------------------------------------------------------------------
