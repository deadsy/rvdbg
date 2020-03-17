//-----------------------------------------------------------------------------
/*

Color Strings

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

const colorRed = "\033[0;31m"
const colorGreen = "\033[0;32m"
const colorBlue = "\033[0;34m"
const colorNone = "\033[0m"

//-----------------------------------------------------------------------------

// RedString returns a red string.
func RedString(s string) string {
	return colorRed + s + colorNone
}

// GreenString returns a green string.
func GreenString(s string) string {
	return colorGreen + s + colorNone
}

//-----------------------------------------------------------------------------
