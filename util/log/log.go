//-----------------------------------------------------------------------------
/*

Wrapper on standard logging.

*/
//-----------------------------------------------------------------------------

package log

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

//-----------------------------------------------------------------------------

// Writer is the log writer.
type Writer struct{}

// Logging types.
var (
	Info  = log.New(Writer{}, "INFO ", 0)
	Debug = log.New(Writer{}, "DEBUG ", 0)
	Error = log.New(Writer{}, "ERROR ", 0)
)

func (f Writer) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	log.Printf("%s:%d %s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}

//-----------------------------------------------------------------------------
