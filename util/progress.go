//-----------------------------------------------------------------------------
/*

Progress Indicator

*/
//-----------------------------------------------------------------------------

package util

import (
	"fmt"

	"github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// Progress contains the state for a progress indicator.
type Progress struct {
	ui       cli.USER // access to user interface
	scale    float32  // scaling constant, n to percentage
	percent  int      // current percentage
	progress string
}

// NewProgress returns the state for a new progress indicator.
func NewProgress(ui cli.USER, nmax int) *Progress {
	return &Progress{
		ui:    ui,
		scale: 100.0 / float32(nmax),
	}
}

// Erase erases the progress indication.
func (p *Progress) Erase() {
	n := len(p.progress)
	if n == 0 {
		return
	}
	s := make([]rune, 3*n)
	for i := 0; i < n; i++ {
		s[i] = '\b'
		s[n+i] = ' '
		s[(2*n)+i] = '\b'
	}
	p.ui.Put(string(s))
}

// Update updates the progress indication.
func (p *Progress) Update(n int) {
	percent := int(float32(n) * p.scale)
	if len(p.progress) == 0 || percent != p.percent {
		p.Erase()
		p.progress = fmt.Sprintf("%d%%", percent)
		p.percent = percent
		p.ui.Put(p.progress)
	}
}

//-----------------------------------------------------------------------------
