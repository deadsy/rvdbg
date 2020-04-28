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

type Progress struct {
	ui       cli.USER
	scale    float32
	mask     int
	progress string
}

func NewProgress(ui cli.USER, nmax int) *Progress {
	n := nmax / 100
	p2 := 1
	mask := 0
	for p2 <= n {
		p2 *= 2
		mask = (mask << 1) | 1
	}
	return &Progress{
		ui:    ui,
		scale: 100.0 / float32(nmax),
		mask:  mask,
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
	if n&p.mask == 0 {
		p.Erase()
		p.progress = fmt.Sprintf("%d%%", int(float32(n)*p.scale))
		p.ui.Put(p.progress)
	}
}

//-----------------------------------------------------------------------------
