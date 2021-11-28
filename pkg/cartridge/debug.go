package cartridge

import (
	"github.com/exp625/gones/internal/textutil"
)

type Debugger interface {
	DebugDisplay(t *textutil.Text)
}
