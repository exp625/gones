package cartridge

import "github.com/faiface/pixel/text"

type Debugger interface {
	DebugDisplay(*text.Text)
}
