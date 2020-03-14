//-----------------------------------------------------------------------------
/*

SoC Device

*/
//-----------------------------------------------------------------------------

package soc

//-----------------------------------------------------------------------------

// CPU provides high-level CPU information.
type CPU struct {
}

// Device is the top-level device description.
type Device struct {
	Vendor      string
	Name        string
	Descr       string
	Version     string
	CPU         *CPU
	Interrupts  []Interrupt
	Peripherals []Peripheral
}

//-----------------------------------------------------------------------------
