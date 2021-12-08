package debugger

import "github.com/exp625/gones/internal/textutil"

func (nes *Debugger) DrawCartridge(t *textutil.Text) {
	nes.Cartridge.Mapper.DebugDisplay(t)
}
