package emulator

import (
	"github.com/exp625/gones/pkg/controller"
	"github.com/hajimehoshi/ebiten/v2"
)

type KeyBinding struct {
	Help       string
	BoundKey   *ebiten.Key
	DefaultKey ebiten.Key
	Pressed    func(e *Emulator)
	Released   func(e *Emulator)
}
type KeyBindingGroup struct {
	Name     string
	Bindings []*KeyBinding
}

func (b *KeyBinding) Key() ebiten.Key {
	if b.BoundKey == nil {
		return b.DefaultKey
	}
	return *b.BoundKey
}

var Bindings []*KeyBindingGroup

func init() {
	Bindings = []*KeyBindingGroup{
		{"Emulator", []*KeyBinding{
			{
				Help:       "Reset the emulator",
				DefaultKey: ebiten.KeyR,
				Pressed:    ResetPressed,
			},
			{
				Help:       "Show the key bindings screen",
				DefaultKey: ebiten.KeyF6,
				Pressed:    ShowKeyBindingsScreenPressed,
			},
			{
				Help:       "Execute one CPU instruction if auto run mode is disabled",
				DefaultKey: ebiten.KeyEnter,
				Pressed:    ExecuteOneCPUInstructionPressed,
			},
			{
				Help:       "Toggle auto run mode",
				DefaultKey: ebiten.KeySpace,
				Pressed:    ToggleAutoRunModePressed,
			},
			{
				Help:       "Issue one master clock",
				DefaultKey: ebiten.KeyRight,
				Pressed:    IssueMasterClockPressed,
			},
			{
				Help:       "Issue one CPU clock (equivalent to three master clocks)",
				DefaultKey: ebiten.KeyUp,
				Pressed:    IssueCPUClockPressed,
			},
			{
				Help:       "Enter a 0 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP0,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(0),
			},
			{
				Help:       "Enter a 1 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP1,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(1),
			},
			{
				Help:       "Enter a 2 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP2,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(2),
			},
			{
				Help:       "Enter a 3 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP3,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(3),
			},
			{
				Help:       "Enter a 4 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP4,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(4),
			},
			{
				Help:       "Enter a 5 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP5,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(5),
			},
			{
				Help:       "Enter a 6 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP6,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(6),
			},
			{
				Help:       "Enter a 7 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP7,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(7),
			},
			{
				Help:       "Enter a 8 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP8,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(8),
			},
			{
				Help:       "Enter a 9 into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP9,
				Pressed:    AddDigitToRequestedNumberOfCyclesPressedFunc(9),
			},
			{
				Help:       "Clear the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyEscape,
				Pressed:    ClearRequestedNumberOfCyclesPressed,
			},
			{
				Help:       "Reset the emulator and set the program counter to 0x4000",
				DefaultKey: ebiten.KeyQ,
				Pressed:    ResetForNESTestPressed,
			},
		}},
		{"Debug", []*KeyBinding{
			{
				Help:       "Show the CPU debug screen",
				DefaultKey: ebiten.KeyF1,
				Pressed:    ShowScreenPressedFunc(ScreenDebugCPU),
			},
			{
				Help:       "Show the PPU debug screen",
				DefaultKey: ebiten.KeyF2,
				Pressed:    ShowScreenPressedFunc(ScreenDebugPPU),
			},
			{
				Help:       "Show the nametable debug screen",
				DefaultKey: ebiten.KeyF3,
				Pressed:    ShowScreenPressedFunc(ScreenDebugNametables),
			},
			{
				Help:       "Show the palette debug screen",
				DefaultKey: ebiten.KeyF4,
				Pressed:    ShowScreenPressedFunc(ScreenDebugPalettes),
			},
			{
				Help:       "Show the controller debug screen",
				DefaultKey: ebiten.KeyF5,
				Pressed:    ShowScreenPressedFunc(ScreenDebugController),
			},
			{
				Help:       "Enable logging",
				DefaultKey: ebiten.KeyL,
				Pressed: func(e *Emulator) {
					if e.LoggingEnabled {
						e.StopLogging()
					} else {
						e.StartLogging()
					}
				},
			},
		},
		},
		{"Controller", []*KeyBinding{
			{
				Help:       "NES standard controller 'A' button",
				DefaultKey: ebiten.KeyP,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonA)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonA)
				},
			},
			{
				Help:       "NES standard controller 'B' button",
				DefaultKey: ebiten.KeyO,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonB)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonB)
				},
			},
			{
				Help:       "NES standard controller 'SELECT' button",
				DefaultKey: ebiten.KeyG,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonSELECT)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonSELECT)
				},
			},
			{
				Help:       "NES standard controller 'START' button",
				DefaultKey: ebiten.KeyH,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonSTART)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonSTART)
				},
			},
			{
				Help:       "NES standard controller 'UP' button",
				DefaultKey: ebiten.KeyW,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonUP)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonUP)
				},
			},
			{
				Help:       "NES standard controller 'DOWN' button",
				DefaultKey: ebiten.KeyS,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonDOWN)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonDOWN)
				},
			},
			{
				Help:       "NES standard controller 'LEFT' button",
				DefaultKey: ebiten.KeyA,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonLEFT)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonLEFT)
				},
			},
			{
				Help:       "NES standard controller 'RIGHT' button",
				DefaultKey: ebiten.KeyD,
				Pressed: func(e *Emulator) {
					e.Controller1.Press(controller.ButtonRIGHT)
				},
				Released: func(e *Emulator) {
					e.Controller1.Release(controller.ButtonRIGHT)
				},
			},
		}},
	}
}

func ResetPressed(e *Emulator) {
	e.Reset()
}

func ShowKeyBindingsScreenPressed(e *Emulator) {
	if e.Screen == ScreenKeybindings {
		e.Screen = ScreenGame
	} else {
		e.Screen = ScreenKeybindings
	}
}

func ExecuteOneCPUInstructionPressed(e *Emulator) {
	if e.autoRunEnabled {
		return
	}
	if e.requestedSteps == 0 {
		e.requestedSteps = 1
	}
	for e.requestedSteps != 0 {
		e.Clock()
		e.Clock()
		e.Clock()
		for e.CPU.CycleCount != 0 {
			e.Clock()
			e.Clock()
			e.Clock()
		}
		e.requestedSteps--
	}
	e.requestedSteps = 0
}

func ToggleAutoRunModePressed(e *Emulator) {
	e.autoRunEnabled = !e.autoRunEnabled
}

func IssueMasterClockPressed(e *Emulator) {
	if !e.autoRunEnabled {
		e.Clock()
	}
}

func IssueCPUClockPressed(e *Emulator) {
	if !e.autoRunEnabled {
		e.Clock()
		e.Clock()
		e.Clock()
	}
}

func AddDigitToRequestedNumberOfCyclesPressedFunc(n int) func(e *Emulator) {
	return func(e *Emulator) {
		e.requestedSteps = e.requestedSteps*10 + n
	}
}

func ClearRequestedNumberOfCyclesPressed(e *Emulator) {
	e.requestedSteps = 0
}

func ResetForNESTestPressed(e *Emulator) {
	if !e.autoRunEnabled {
		e.Reset()
		e.CPU.PC = 0xC000
		e.CPU.P = 0x24
	}
}

func ShowScreenPressedFunc(screen Screen) func(e *Emulator) {
	return func(e *Emulator) {
		if e.Screen == screen {
			e.Screen = ScreenGame
		} else {
			e.Screen = screen
		}
	}
}
