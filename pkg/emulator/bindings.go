package emulator

import (
	"github.com/exp625/gones/pkg/controller"
	"github.com/hajimehoshi/ebiten/v2"
)

type Binding struct {
	Help       string
	BoundKey   *ebiten.Key
	DefaultKey ebiten.Key
	OnPressed  func()
	OnReleased func()
}

func (b *Binding) Key() ebiten.Key {
	if b.BoundKey == nil {
		return b.DefaultKey
	}
	return *b.BoundKey
}

type BindingGroup struct {
	Name     string
	Bindings []*Binding
}

func DefaultBindings(e *Emulator) []*BindingGroup {
	return []*BindingGroup{
		{"Emulator", []*Binding{
			{
				Help:       "Reset the emulator",
				DefaultKey: ebiten.KeyR,
				OnPressed:  e.ResetPressed,
			},
			{
				Help:       "Choose a NEW ROM to load into the emulator",
				DefaultKey: ebiten.KeyHome,
				OnPressed:  e.ShowScreenPressedFunc(OverlayROMChooser),
			},
			{
				Help:       "Show the sprites debug screen",
				DefaultKey: ebiten.KeyF6,
				OnPressed:  e.ShowScreenPressedFunc(OverlaySprites),
			},
			{
				Help:       "Show the key bindings screen",
				DefaultKey: ebiten.KeyF7,
				OnPressed:  e.ShowScreenPressedFunc(OverlayKeybindings),
			},
			{
				Help:       "Execute one CPU instruction if auto run mode is disabled",
				DefaultKey: ebiten.KeyEnter,
				OnPressed:  e.ExecuteOneCPUInstructionPressed,
			},
			{
				Help:       "Toggle auto run mode",
				DefaultKey: ebiten.KeySpace,
				OnPressed:  e.ToggleAutoRunModePressed,
			},
			{
				Help:       "Issue one master clock",
				DefaultKey: ebiten.KeyRight,
				OnPressed:  e.IssueMasterClockPressed,
			},
			{
				Help:       "Issue one CPU clock (equivalent to three master clocks)",
				DefaultKey: ebiten.KeyUp,
				OnPressed:  e.IssueCPUClockPressed,
			},
			{
				Help:       "enter a '0' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP0,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(0),
			},
			{
				Help:       "enter a '1' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP1,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(1),
			},
			{
				Help:       "enter a '2' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP2,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(2),
			},
			{
				Help:       "enter a '3' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP3,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(3),
			},
			{
				Help:       "enter a '4' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP4,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(4),
			},
			{
				Help:       "enter a '5' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP5,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(5),
			},
			{
				Help:       "enter a '6' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP6,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(6),
			},
			{
				Help:       "enter a '7' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP7,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(7),
			},
			{
				Help:       "enter a '8' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP8,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(8),
			},
			{
				Help:       "enter a '9' into the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyKP9,
				OnPressed:  e.AddDigitToRequestedNumberOfCyclesPressedFunc(9),
			},
			{
				Help:       "Clear the 'requested number of cycles' field",
				DefaultKey: ebiten.KeyEscape,
				OnPressed:  e.ClearRequestedNumberOfCyclesPressed,
			},
			{
				Help:       "Reset the emulator and set the program counter to 0x4000",
				DefaultKey: ebiten.KeyQ,
				OnPressed:  e.ResetForNESTestPressed,
			},
		}},
		{"Debug", []*Binding{
			{
				Help:       "Show the CPU debug screen",
				DefaultKey: ebiten.KeyF1,
				OnPressed:  e.ShowScreenPressedFunc(OverlayCPU),
			},
			{
				Help:       "Show the PPU debug screen",
				DefaultKey: ebiten.KeyF2,
				OnPressed:  e.ShowScreenPressedFunc(OverlayPPU),
			},
			{
				Help:       "Show the nametable debug screen",
				DefaultKey: ebiten.KeyF3,
				OnPressed:  e.ShowScreenPressedFunc(OverlayNametables),
			},
			{
				Help:       "Show the palette debug screen",
				DefaultKey: ebiten.KeyF4,
				OnPressed:  e.ShowScreenPressedFunc(OverlayPalettes),
			},
			{
				Help:       "Show the controller debug screen",
				DefaultKey: ebiten.KeyF5,
				OnPressed:  e.ShowScreenPressedFunc(OverlayControllers),
			},
			{
				Help:       "Enable logging",
				DefaultKey: ebiten.KeyL,
				OnPressed:  e.StartLoggingPressed,
			},
		}},
		{"Controller", []*Binding{
			{
				Help:       "NES standard controller 'A' button",
				DefaultKey: ebiten.KeyP,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonA),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonA),
			},
			{
				Help:       "NES standard controller 'B' button",
				DefaultKey: ebiten.KeyO,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonB),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonB),
			},
			{
				Help:       "NES standard controller 'SELECT' button",
				DefaultKey: ebiten.KeyG,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonSELECT),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonSELECT),
			},
			{
				Help:       "NES standard controller 'START' button",
				DefaultKey: ebiten.KeyH,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonSTART),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonSTART),
			},
			{
				Help:       "NES standard controller 'UP' button",
				DefaultKey: ebiten.KeyW,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonUP),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonUP),
			},
			{
				Help:       "NES standard controller 'DOWN' button",
				DefaultKey: ebiten.KeyS,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonDOWN),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonDOWN),
			},
			{
				Help:       "NES standard controller 'LEFT' button",
				DefaultKey: ebiten.KeyA,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonLEFT),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonLEFT),
			},
			{
				Help:       "NES standard controller 'RIGHT' button",
				DefaultKey: ebiten.KeyD,
				OnPressed:  e.ControllerButtonPressedFunc(controller.ButtonRIGHT),
				OnReleased: e.ControllerButtonReleasedFunc(controller.ButtonRIGHT),
			},
		}},
	}
}

func (e *Emulator) ResetPressed() {
	e.Reset()
}

func (e *Emulator) ExecuteOneCPUInstructionPressed() {
	if e.AutoRunEnabled {
		return
	}
	if e.RequestedSteps == 0 {
		e.RequestedSteps = 1
	}
	for e.RequestedSteps != 0 {
		e.Clock()
		e.Clock()
		e.Clock()
		for e.CPU.CycleCount != 0 {
			e.Clock()
			e.Clock()
			e.Clock()
		}
		e.RequestedSteps--
	}
	e.RequestedSteps = 0
}

func (e *Emulator) ToggleAutoRunModePressed() {
	e.AutoRunEnabled = !e.AutoRunEnabled
}

func (e *Emulator) IssueMasterClockPressed() {
	if !e.AutoRunEnabled {
		e.Clock()
	}
}

func (e *Emulator) IssueCPUClockPressed() {
	if !e.AutoRunEnabled {
		e.Clock()
		e.Clock()
		e.Clock()
	}
}

func (e *Emulator) AddDigitToRequestedNumberOfCyclesPressedFunc(n int) func() {
	return func() {
		e.RequestedSteps = e.RequestedSteps*10 + n
	}
}

func (e *Emulator) ClearRequestedNumberOfCyclesPressed() {
	e.RequestedSteps = 0
}

func (e *Emulator) ResetForNESTestPressed() {
	if !e.AutoRunEnabled {
		e.Reset()
		e.CPU.PC = 0xC000
		e.CPU.P = 0x24
	}
}

func (e *Emulator) ShowScreenPressedFunc(screen Overlay) func() {
	return func() {
		if e.ActiveOverlay == screen {
			e.ActiveOverlay = OverlayGame
		} else {
			e.ActiveOverlay = screen
		}
	}
}

func (e *Emulator) StartLoggingPressed() {
	if e.Logger.LoggingEnabled() {
		e.Logger.StopLogging()
	} else {
		e.Logger.StartLogging()
	}
}

func (e *Emulator) ControllerButtonPressedFunc(b controller.Button) func() {
	return func() {
		e.Controller1.Press(b)
	}
}

func (e *Emulator) ControllerButtonReleasedFunc(b controller.Button) func() {
	return func() {
		e.Controller1.Release(b)
	}
}
