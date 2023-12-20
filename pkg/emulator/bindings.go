package emulator

import (
	"github.com/exp625/gones/pkg/controller"
	"github.com/exp625/gones/pkg/input"
)

func (e *Emulator) clearAllBindings() {
	for _, group := range e.Bindings.Groups {
		for _, binding := range group {
			binding.OnPressed = nil
			binding.OnReleased = nil
		}
	}
	e.Bindings.GlobalHandler = nil
	e.Bindings.TextHandler = nil
	e.Bindings.NumberHandler = nil
	e.Bindings.RepeatKeys = false

}

func (e *Emulator) registerEmulatorBindings() {

	e.Bindings.Groups[input.Emulator][input.Reset].OnPressed = e.Reset
	e.Bindings.Groups[input.Emulator][input.Load].OnPressed = func() { e.ChangeScreen(SettingROMChooser) }
	e.Bindings.Groups[input.Emulator][input.Save].OnPressed = e.SaveGame
	e.Bindings.Groups[input.Emulator][input.ScreenKeyBindings].OnPressed = func() { e.ChangeScreen(SettingKeybindings) }
	e.Bindings.Groups[input.Emulator][input.ExecuteInstruction].OnPressed = e.executeOneCPUInstructionPressed
	e.Bindings.Groups[input.Emulator][input.Pause].OnPressed = func() { e.AutoRunEnabled = !e.AutoRunEnabled }
	e.Bindings.Groups[input.Emulator][input.OpenDebug].OnPressed = func() { e.ChangeScreen(OverlayCPU) }
	e.Bindings.Groups[input.Emulator][input.OpenSettings].OnPressed = func() { e.ChangeScreen(SettingKeybindings) }
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
	e.Bindings.Groups[input.Debug][input.ShowCPUDebug].OnPressed = func() { e.ChangeScreen(OverlayCPU) }
	e.Bindings.Groups[input.Debug][input.ShowPPUDebug].OnPressed = func() { e.ChangeScreen(OverlayPPU) }
	e.Bindings.Groups[input.Debug][input.ShowAPUDebug].OnPressed = func() { e.ChangeScreen(OverlayAPU) }
	e.Bindings.Groups[input.Debug][input.ShowNametableDebug].OnPressed = func() { e.ChangeScreen(OverlayNametables) }
	e.Bindings.Groups[input.Debug][input.ShowPaletteDebug].OnPressed = func() { e.ChangeScreen(OverlayPalettes) }
	e.Bindings.Groups[input.Debug][input.ShowControllerDebug].OnPressed = func() { e.ChangeScreen(OverlayControllers) }
	e.Bindings.Groups[input.Debug][input.ShowSpriteDebug].OnPressed = func() { e.ChangeScreen(OverlaySprites) }
	e.Bindings.Groups[input.Debug][input.LeaveDebug].OnPressed = func() { e.ChangeScreen(ScreenGame) }
	e.Bindings.Groups[input.Debug][input.EnableLogging].OnPressed = func() {
		if e.Logger.LoggingEnabled() {
			e.Logger.StopLogging()
		} else {
			e.Logger.StartLogging()
		}
	}
	e.Bindings.Groups[input.Debug][input.ForceC000].OnPressed = func() {
		if !e.AutoRunEnabled {
			e.CPU.PC = 0xC000
		}
	}
}

func (e *Emulator) registerSettingsBindings() {
	e.Bindings.Groups[input.Settings][input.ROMChooser].OnPressed = func() { e.ChangeScreen(SettingROMChooser) }
	e.Bindings.Groups[input.Settings][input.SAVEChooser].OnPressed = func() { e.ChangeScreen(SettingsSave) }
	e.Bindings.Groups[input.Settings][input.Keybindings].OnPressed = func() { e.ChangeScreen(SettingKeybindings) }
	e.Bindings.Groups[input.Settings][input.AudioSettings].OnPressed = func() { e.ChangeScreen(SettingAudio) }
	e.Bindings.Groups[input.Settings][input.Cancel].OnPressed = func() { e.ChangeScreen(ScreenGame) }
}

func (e *Emulator) registerAudioBindings() {

}

func (e *Emulator) registerControllerBindings() {
	e.Bindings.Groups[input.Controller][input.A].OnPressed = func() { e.Controller1.Press(controller.ButtonA) }
	e.Bindings.Groups[input.Controller][input.B].OnPressed = func() { e.Controller1.Press(controller.ButtonB) }
	e.Bindings.Groups[input.Controller][input.START].OnPressed = func() { e.Controller1.Press(controller.ButtonSTART) }
	e.Bindings.Groups[input.Controller][input.SELECT].OnPressed = func() { e.Controller1.Press(controller.ButtonSELECT) }
	e.Bindings.Groups[input.Controller][input.UP].OnPressed = func() { e.Controller1.Press(controller.ButtonUP) }
	e.Bindings.Groups[input.Controller][input.DOWN].OnPressed = func() { e.Controller1.Press(controller.ButtonDOWN) }
	e.Bindings.Groups[input.Controller][input.LEFT].OnPressed = func() { e.Controller1.Press(controller.ButtonLEFT) }
	e.Bindings.Groups[input.Controller][input.RIGHT].OnPressed = func() { e.Controller1.Press(controller.ButtonRIGHT) }

	e.Bindings.Groups[input.Controller][input.A].OnReleased = func() { e.Controller1.Release(controller.ButtonA) }
	e.Bindings.Groups[input.Controller][input.B].OnReleased = func() { e.Controller1.Release(controller.ButtonB) }
	e.Bindings.Groups[input.Controller][input.START].OnReleased = func() { e.Controller1.Release(controller.ButtonSTART) }
	e.Bindings.Groups[input.Controller][input.SELECT].OnReleased = func() { e.Controller1.Release(controller.ButtonSELECT) }
	e.Bindings.Groups[input.Controller][input.UP].OnReleased = func() { e.Controller1.Release(controller.ButtonUP) }
	e.Bindings.Groups[input.Controller][input.DOWN].OnReleased = func() { e.Controller1.Release(controller.ButtonDOWN) }
	e.Bindings.Groups[input.Controller][input.LEFT].OnReleased = func() { e.Controller1.Release(controller.ButtonLEFT) }
	e.Bindings.Groups[input.Controller][input.RIGHT].OnReleased = func() { e.Controller1.Release(controller.ButtonRIGHT) }
}

func (e *Emulator) registerNumberHandler() {
	e.Bindings.NumberHandler = func(n int) {
		e.RequestedSteps = e.RequestedSteps*10 + n
	}
}

func (e *Emulator) registerFileExplorerBindings() {
	e.Bindings.Groups[input.FileExplorer][input.Select].OnPressed = func() { e.FileExplorer.Ready = true }
	e.Bindings.Groups[input.FileExplorer][input.MoveSelectionUp].OnPressed = func() { e.FileExplorer.Selected -= 1 }
	e.Bindings.Groups[input.FileExplorer][input.MoveSelectionDown].OnPressed = func() { e.FileExplorer.Selected += 1 }
	e.Bindings.Groups[input.FileExplorer][input.OpenFolder].OnPressed = e.FileExplorer.OpenFolder
	e.Bindings.Groups[input.FileExplorer][input.ParentFolder].OnPressed = e.FileExplorer.CloseFolder
	e.Bindings.Groups[input.Emulator][input.Cancel].OnPressed = func() { e.ChangeScreen(ScreenGame) }
	e.Bindings.TextHandler = e.FileExplorer.TextInput
	e.Bindings.RepeatKeys = true
}

func (e *Emulator) registerInputBindings() {

	e.Bindings.Groups[input.FileExplorer][input.MoveSelectionUp].OnPressed = e.Bindings.MoveSelectionUp
	e.Bindings.Groups[input.FileExplorer][input.MoveSelectionDown].OnPressed = e.Bindings.MoveSelectionDown
	e.Bindings.Groups[input.FileExplorer][input.OpenFolder].OnPressed = e.Bindings.Select
	e.Bindings.Groups[input.FileExplorer][input.ParentFolder].OnPressed = e.Bindings.Deselect
	e.Bindings.Groups[input.Emulator][input.Cancel].OnPressed = func() {
		e.Bindings.SaveCustomBindings()
		e.ChangeScreen(ScreenGame)
	}
	e.Bindings.Groups[input.Emulator][input.Reset].OnPressed = e.Bindings.Reset
	e.Bindings.Groups[input.FileExplorer][input.Select].OnPressed = e.Bindings.ChangeKeybinding
	e.Bindings.GlobalHandler = e.Bindings.ApplyKeybinding

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
