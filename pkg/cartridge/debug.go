package cartridge

import "strings"

type Debugger interface {
	DebugDisplay(*strings.Builder)
}
