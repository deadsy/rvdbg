//-----------------------------------------------------------------------------
/*

RISC-V Debugger 0.11
Debug Bus Access

*/
//-----------------------------------------------------------------------------

package rv11

//-----------------------------------------------------------------------------
// hart selection

/*

// setHartSelect sets the hart select value in a dmcontrol value.
func setHartSelect(x uint32, id int) uint32 {
	x &= ^uint32(hartselhi | hartsello)
	hi := ((id >> 10) << 6) & hartselhi
	lo := (id << 16) & hartsello
	return x | uint32(hi) | uint32(lo)
}

// getHartSelect gets the hart select value from a dmcontrol value.
func getHartSelect(x uint32) int {
	return int((util.Bits(uint(x), 15, 6) << 10) | util.Bits(uint(x), 25, 16))
}

*/

// selectHart sets the dmcontrol hartsel value.
func (dbg *Debug) selectHart(id int) error {

	/*

		x, err := dbg.rdDmi(dmcontrol)
		if err != nil {
			return err
		}
		x = setHartSelect(x, id)
		return dbg.wrDmi(dmcontrol, x)

	*/

	return nil
}

//-----------------------------------------------------------------------------
