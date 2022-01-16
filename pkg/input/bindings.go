package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
)

type GroupName string
type BindingName string

const (
	Emulator     GroupName = "Emulator"
	Debug                  = "Debug"
	Controller1            = "Controller 1"
	Controller2            = "Controller 2"
	FileExplorer           = "File Explorer"

	Reset  BindingName = "Reset"
	Load               = "Load"
	Save               = "Save"
	Pause              = "Pause"
	Cancel             = "Cancel"

	ScreenKeyBindings   = "Key Bindings Screen"
	ExecuteInstruction  = "Execute Instruction"
	ExecuteCPUClock     = "ExecuteCPUClock"
	ExecuteMasterClock  = "ExecuteMasterClock"
	ShowCPUDebug        = "Show CPU Debug"
	ShowPPUDebug        = "Show PPU Debug"
	ShowNametableDebug  = "Show Nametable Debug"
	ShowPaletteDebug    = "Show Palette Debug"
	ShowSpriteDebug     = "Show Sprite Debug"
	ShowControllerDebug = "Show Controller Debug"
	EnableLogging       = "EnableLogging"

	Select            = "Select"
	OpenFolder        = "OpenFolder"
	ParentFolder      = "ParentFolder"
	MoveSelectionUp   = "MoveSelectionUp"
	MoveSelectionDown = "MoveSelectionDown"

	A      = "A"
	B      = "B"
	UP     = "UP"
	DOWN   = "DOWN"
	LEFT   = "LEFT"
	RIGHT  = "RIGHT"
	START  = "START"
	SELECT = "SELECT"
)

type Binding struct {
	Help                       string
	BoundKey                   ebiten.Key
	HasBoundKey                bool
	DefaultKey                 ebiten.Key
	DefaultControllerButton    ebiten.StandardGamepadButton
	HasDefaultControllerButton bool
	BoundControllerButton      ebiten.StandardGamepadButton
	HasBoundControllerButton   bool
	DefaultControllerAxis      ebiten.StandardGamepadAxis
	DefaultControllerAxisSign  float64
	HasDefaultControllerAxis   bool
	BoundControllerAxis        ebiten.StandardGamepadAxis
	BoundControllerAxisSign    float64
	HAsBoundControllerAxis     bool
	OnPressed                  func()
	OnReleased                 func()
}

type Bindings struct {
	Groups        BindingGroups
	NumberHandler func(int)
	TextHandler   func(rune)
}

func (b *Binding) Key() ebiten.Key {
	if !b.HasBoundKey {
		return b.DefaultKey
	}
	return b.BoundKey
}

func (b *Binding) ControllerButton() (bool, ebiten.StandardGamepadButton) {
	if b.HasBoundControllerButton {
		return true, b.BoundControllerButton
	}
	if b.HasDefaultControllerButton {
		return true, b.DefaultControllerButton
	}
	return false, 0
}

func (b *Binding) ControllerAxis() (bool, ebiten.StandardGamepadAxis, float64) {
	if b.HAsBoundControllerAxis {
		return true, b.BoundControllerAxis, b.BoundControllerAxisSign
	}
	if b.HasDefaultControllerAxis {
		return true, b.DefaultControllerAxis, b.DefaultControllerAxisSign
	}
	return false, 0, 0
}

type BindingGroup map[BindingName]*Binding
type BindingGroups map[GroupName]BindingGroup

func GetBindings() *Bindings {
	return &Bindings{
		Groups: BindingGroups{
			Emulator: BindingGroup{
				Reset: &Binding{
					Help:       "Reset the emulator",
					DefaultKey: ebiten.KeyR,
				},
				Load: &Binding{
					Help:       "Choose a NEW ROM to load into the emulator",
					DefaultKey: ebiten.KeyHome,
				},
				Save: &Binding{
					Help:       "Save Game",
					DefaultKey: ebiten.KeyF12,
				},
				ScreenKeyBindings: &Binding{
					Help:       "Show the key bindings screen",
					DefaultKey: ebiten.KeyF7,
				},
				ExecuteInstruction: &Binding{
					Help:       "Execute one CPU instruction if auto run mode is disabled",
					DefaultKey: ebiten.KeyEnter,
				},
				Pause: &Binding{
					Help:       "Pause emulator",
					DefaultKey: ebiten.KeySpace,
				},
				ExecuteMasterClock: &Binding{
					Help:       "Issue one master clock",
					DefaultKey: ebiten.KeyRight,
				},
				ExecuteCPUClock: &Binding{
					Help:       "Issue one CPU clock (equivalent to three master clocks)",
					DefaultKey: ebiten.KeyUp,
				},
				Cancel: &Binding{
					Help:       "Clear the 'requested number of cycles' field",
					DefaultKey: ebiten.KeyEscape,
				},
			},
			Debug: BindingGroup{
				ShowCPUDebug: &Binding{
					Help:       "Show the CPU debug screen",
					DefaultKey: ebiten.KeyF1,
				},
				ShowPPUDebug: &Binding{
					Help:       "Show the PPU debug screen",
					DefaultKey: ebiten.KeyF2,
				},
				ShowNametableDebug: &Binding{
					Help:       "Show the nametable debug screen",
					DefaultKey: ebiten.KeyF3,
				},
				ShowPaletteDebug: &Binding{
					Help:       "Show the palette debug screen",
					DefaultKey: ebiten.KeyF4,
				},
				ShowControllerDebug: &Binding{
					Help:       "Show the controller debug screen",
					DefaultKey: ebiten.KeyF5,
				},
				ShowSpriteDebug: &Binding{
					Help:       "Show the sprites debug screen",
					DefaultKey: ebiten.KeyF6,
				},
				EnableLogging: &Binding{
					Help:       "Enable logging",
					DefaultKey: ebiten.KeyL,
				},
			},
			Controller1: BindingGroup{
				A: &Binding{
					Help:                       "NES standard controller 'A' button",
					DefaultKey:                 ebiten.KeyP,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonRightBottom,
				},
				B: &Binding{
					Help:                       "NES standard controller 'B' button",
					DefaultKey:                 ebiten.KeyO,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonRightRight,
				},
				SELECT: &Binding{
					Help:                       "NES standard controller 'SELECT' button",
					DefaultKey:                 ebiten.KeyG,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonCenterLeft,
				},
				START: &Binding{
					Help:                       "NES standard controller 'START' button",
					DefaultKey:                 ebiten.KeyH,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonCenterRight,
				},
				UP: &Binding{
					Help:                       "NES standard controller 'UP' button",
					DefaultKey:                 ebiten.KeyW,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonLeftTop,
					HasDefaultControllerAxis:   true,
					DefaultControllerAxis:      ebiten.StandardGamepadAxisLeftStickVertical,
					DefaultControllerAxisSign:  -1,
				},
				DOWN: &Binding{
					Help:                       "NES standard controller 'DOWN' button",
					DefaultKey:                 ebiten.KeyS,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonLeftBottom,
					HasDefaultControllerAxis:   true,
					DefaultControllerAxis:      ebiten.StandardGamepadAxisLeftStickVertical,
					DefaultControllerAxisSign:  +1,
				},
				LEFT: &Binding{
					Help:                       "NES standard controller 'LEFT' button",
					DefaultKey:                 ebiten.KeyA,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonLeftLeft,
					HasDefaultControllerAxis:   true,
					DefaultControllerAxis:      ebiten.StandardGamepadAxisLeftStickHorizontal,
					DefaultControllerAxisSign:  -1,
				},
				RIGHT: &Binding{
					Help:                       "NES standard controller 'RIGHT' button",
					DefaultKey:                 ebiten.KeyD,
					HasDefaultControllerButton: true,
					DefaultControllerButton:    ebiten.StandardGamepadButtonLeftRight,
					HasDefaultControllerAxis:   true,
					DefaultControllerAxis:      ebiten.StandardGamepadAxisLeftStickHorizontal,
					DefaultControllerAxisSign:  +1,
				},
			},
			FileExplorer: BindingGroup{
				Select: &Binding{
					Help:       "Open the selected file",
					DefaultKey: ebiten.KeyEnter,
				},
				OpenFolder: &Binding{
					Help:       "Open the selected folder",
					DefaultKey: ebiten.KeyArrowRight,
				},
				ParentFolder: &Binding{
					Help:       "Open the parent folder",
					DefaultKey: ebiten.KeyArrowLeft,
				},
				MoveSelectionUp: &Binding{
					Help:       "Move the selected on file up",
					DefaultKey: ebiten.KeyArrowUp,
				},
				MoveSelectionDown: &Binding{
					Help:       "Move the selected on file down",
					DefaultKey: ebiten.KeyArrowDown,
				},
			},
		},
		NumberHandler: func(int) {

		},
		TextHandler: func(rune) {

		},
	}
}

func LoadCustomBindings(file os.File, bindings Bindings) {

}
