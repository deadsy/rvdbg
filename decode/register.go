//-----------------------------------------------------------------------------
/*

Peripheral Registers

*/
//-----------------------------------------------------------------------------

package decode

//-----------------------------------------------------------------------------

// Register is peripheral register.
type Register struct {
	Name   string
	Offset uint
	Size   uint
	Fset   FieldSet
	Descr  string
}

// RegisterSet is a set of registers.
type RegisterSet []Register

//-----------------------------------------------------------------------------
