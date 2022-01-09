package emulator

import (
	"github.com/exp625/gones/pkg/controller"
	"github.com/exp625/gones/pkg/input"
)

func (e *Emulator) clearAllBindings() {
	for _, group := range e.Bindings.Groups {
		for _, binding := range group {
			binding.OnPressed = func() {}
			binding.OnReleased = func() {}
		}
	}
}

func (e *Emulator) registerAllBindings() {
	e.registerEmulatorBindings()
	e.registerControllerBindings()
	e.registerNumberHandler()
	e.registerDebugBindings()
}
func (e *Emulator) registerEmulatorBindings() {
	e.Bindings.Groups[input.Emulator][input.Reset].OnPressed = e.Reset
	e.Bindings.Groups[input.Emulator][input.Load].OnPressed = func() { e.ActiveOverlay = OverlayROMChooser }
	e.Bindings.Groups[input.Emulator][input.Save].OnPressed = e.SaveGame
	e.Bindings.Groups[input.Emulator][input.ScreenKeyBindings].OnPressed = func() { e.ActiveOverlay = OverlayKeybindings }
	e.Bindings.Groups[input.Emulator][input.ExecuteInstruction].OnPressed = e.executeOneCPUInstructionPressed
	e.Bindings.Groups[input.Emulator][input.Pause].OnPressed = func() { e.AutoRunEnabled = !e.AutoRunEnabled }
	e.Bindings.Groups[input.Emulator][input.ExecuteMasterClock].OnPressed = func() {
		if !e.AutoRunEnabled {
			e.Clock()
		}
	}
	e.Bindings.Groups[input.Emulator][input.ExecuteCPUClock].OnPressed = func() {
		if !e.AutoRunEnabled {
			e.Clock()
			e.Clock()
			e.Clock()
		}
	}
	e.Bindings.Groups[input.Emulator][input.Cancel].OnPressed = func() {
		e.RequestedSteps = 0
	}
}

func (e *Emulator) registerDebugBindings() {
	e.Bindings.Groups[input.Debug][input.ShowCPUDebug].OnPressed = func() { e.ActiveOverlay = OverlayCPU }
	e.Bindings.Groups[input.Debug][input.ShowPPUDebug].OnPressed = func() { e.ActiveOverlay = OverlayPPU }
	e.Bindings.Groups[input.Debug][input.ShowNametableDebug].OnPressed = func() { e.ActiveOverlay = OverlayNametables }
	e.Bindings.Groups[input.Debug][input.ShowPaletteDebug].OnPressed = func() { e.ActiveOverlay = OverlayPalettes }
	e.Bindings.Groups[input.Debug][input.ShowControllerDebug].OnPressed = func() { e.ActiveOverlay = OverlayControllers }
	e.Bindings.Groups[input.Debug][input.ShowSpriteDebug].OnPressed = func() { e.ActiveOverlay = OverlaySprites }
	e.Bindings.Groups[input.Debug][input.EnableLogging].OnPressed = func() {
		if e.Logger.LoggingEnabled() {
			e.Logger.StopLogging()
		} else {
			e.Logger.StartLogging()
		}
	}
}

func (e *Emulator) registerControllerBindings() {
	e.Bindings.Groups[input.Controller1][input.A].OnPressed = func() { e.Controller1.Press(controller.ButtonA) }
	e.Bindings.Groups[input.Controller1][input.B].OnPressed = func() { e.Controller1.Press(controller.ButtonB) }
	e.Bindings.Groups[input.Controller1][input.START].OnPressed = func() { e.Controller1.Press(controller.ButtonSTART) }
	e.Bindings.Groups[input.Controller1][input.SELECT].OnPressed = func() { e.Controller1.Press(controller.ButtonSELECT) }
	e.Bindings.Groups[input.Controller1][input.UP].OnPressed = func() { e.Controller1.Press(controller.ButtonUP) }
	e.Bindings.Groups[input.Controller1][input.DOWN].OnPressed = func() { e.Controller1.Press(controller.ButtonDOWN) }
	e.Bindings.Groups[input.Controller1][input.LEFT].OnPressed = func() { e.Controller1.Press(controller.ButtonLEFT) }
	e.Bindings.Groups[input.Controller1][input.RIGHT].OnPressed = func() { e.Controller1.Press(controller.ButtonRIGHT) }

	e.Bindings.Groups[input.Controller1][input.A].OnReleased = func() { e.Controller1.Release(controller.ButtonA) }
	e.Bindings.Groups[input.Controller1][input.B].OnReleased = func() { e.Controller1.Release(controller.ButtonB) }
	e.Bindings.Groups[input.Controller1][input.START].OnReleased = func() { e.Controller1.Release(controller.ButtonSTART) }
	e.Bindings.Groups[input.Controller1][input.SELECT].OnReleased = func() { e.Controller1.Release(controller.ButtonSELECT) }
	e.Bindings.Groups[input.Controller1][input.UP].OnReleased = func() { e.Controller1.Release(controller.ButtonUP) }
	e.Bindings.Groups[input.Controller1][input.DOWN].OnReleased = func() { e.Controller1.Release(controller.ButtonDOWN) }
	e.Bindings.Groups[input.Controller1][input.LEFT].OnReleased = func() { e.Controller1.Release(controller.ButtonLEFT) }
	e.Bindings.Groups[input.Controller1][input.RIGHT].OnReleased = func() { e.Controller1.Release(controller.ButtonRIGHT) }
}

func (e *Emulator) registerNumberHandler() {
	e.Bindings.NumberHandler = func(n int) {
		e.RequestedSteps = e.RequestedSteps*10 + n
	}
}

func (e *Emulator) executeOneCPUInstructionPressed() {
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
