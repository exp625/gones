package emulator

import (
	"github.com/exp625/gones/pkg/controller"
	"github.com/hajimehoshi/ebiten/v2"
)

type Binding struct {
	key        ebiten.Key
	defaultKey ebiten.Key
	Pressed    func(e *Emulator)
	Released   func(e *Emulator)
}

func (b *Binding) Key() ebiten.Key {
	if b.key == 0 {
		return b.defaultKey
	}
	return b.key
}

var Bindings []*Binding

func init() {
	Bindings = []*Binding{
		{
			defaultKey: ebiten.KeyR,
			Pressed: func(e *Emulator) {
				e.Reset()
			}},
		{
			defaultKey: ebiten.KeyF1,
			Pressed: func(e *Emulator) {
				if e.Screen == ScreenDebugCPU {
					e.Screen = ScreenGame
				} else {
					e.Screen = ScreenDebugCPU
				}
			}},
		{
			defaultKey: ebiten.KeyF2,
			Pressed: func(e *Emulator) {
				if e.Screen == ScreenDebugPPU {
					e.Screen = ScreenGame
				} else {
					e.Screen = ScreenDebugPPU
				}
			}},
		{
			defaultKey: ebiten.KeyF3,
			Pressed: func(e *Emulator) {
				if e.Screen == ScreenDebugNametables {
					e.Screen = ScreenGame
				} else {
					e.Screen = ScreenDebugNametables
				}
			}},
		{
			defaultKey: ebiten.KeyF4,
			Pressed: func(e *Emulator) {
				if e.Screen == ScreenDebugPalettes {
					e.Screen = ScreenGame
				} else {
					e.Screen = ScreenDebugPalettes
				}
			}},
		{
			defaultKey: ebiten.KeyF5,
			Pressed: func(e *Emulator) {
				if e.Screen == ScreenDebugController {
					e.Screen = ScreenGame
				} else {
					e.Screen = ScreenDebugController
				}
			}},
		{
			defaultKey: ebiten.KeyL,
			Pressed: func(e *Emulator) {
				if e.LoggingEnabled {
					e.StopLogging()
				} else {
					e.StartLogging()
				}
			}},
		{
			defaultKey: ebiten.KeyEnter,
			Pressed: func(e *Emulator) {
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
			}},
		{
			defaultKey: ebiten.KeySpace,
			Pressed: func(e *Emulator) {
				e.autoRunEnabled = !e.autoRunEnabled
			}},
		{
			defaultKey: ebiten.KeyRight,
			Pressed: func(e *Emulator) {
				if !e.autoRunEnabled {
					e.Clock()
				}
			}},
		{
			defaultKey: ebiten.KeyUp,
			Pressed: func(e *Emulator) {
				if !e.autoRunEnabled {
					e.Clock()
					e.Clock()
					e.Clock()
				}
			}},
		{
			defaultKey: ebiten.KeyKP0,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 0
			}},
		{
			defaultKey: ebiten.KeyKP1,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 1
			}},
		{
			defaultKey: ebiten.KeyKP2,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 2
			}},
		{
			defaultKey: ebiten.KeyKP3,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 3
			}},
		{
			defaultKey: ebiten.KeyKP4,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 4
			}},
		{
			defaultKey: ebiten.KeyKP5,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 5
			}},
		{
			defaultKey: ebiten.KeyKP6,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 6
			}},
		{
			defaultKey: ebiten.KeyKP7,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 7
			}},
		{
			defaultKey: ebiten.KeyKP8,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 8
			}},
		{
			defaultKey: ebiten.KeyKP9,
			Pressed: func(e *Emulator) {
				e.requestedSteps = e.requestedSteps*10 + 9
			}},
		{
			defaultKey: ebiten.KeyEscape,
			Pressed: func(e *Emulator) {
				e.requestedSteps = 0
			}},
		{
			defaultKey: ebiten.KeyQ,
			Pressed: func(e *Emulator) {
				if !e.autoRunEnabled {
					e.Reset()
					e.CPU.PC = 0xC000
					e.CPU.P = 0x24
				}
			}},
		{
			defaultKey: ebiten.KeyP,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonA)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonA)
			}},
		{
			defaultKey: ebiten.KeyO,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonB)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonB)
			}},
		{
			defaultKey: ebiten.KeyG,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonSELECT)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonSELECT)
			}},
		{
			defaultKey: ebiten.KeyH,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonSTART)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonSTART)
			}},
		{
			defaultKey: ebiten.KeyW,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonUP)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonUP)
			}},
		{
			defaultKey: ebiten.KeyS,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonDOWN)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonDOWN)
			}},
		{
			defaultKey: ebiten.KeyA,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonLEFT)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonLEFT)
			}},
		{
			defaultKey: ebiten.KeyD,
			Pressed: func(e *Emulator) {
				e.Controller1.Press(controller.ButtonRIGHT)
			},
			Released: func(e *Emulator) {
				e.Controller1.Release(controller.ButtonRIGHT)
			}},
	}
}
