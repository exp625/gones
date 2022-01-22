package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"math"
)

const AxisThreshold = 0.5

var gamepads Gamepad

func init() {
	gamepads = Gamepad{}
}

type Gamepad struct {
	gamepadIDsBuf []ebiten.GamepadID
	gamepadIDs    map[ebiten.GamepadID]struct{}
	gamepadAxis   map[ebiten.GamepadID]map[float64]bool
}

func HandleInput(bindings *Bindings) {
	handleControllerEvents()
	handleKeyboardInputs(bindings)
	handleControllerInputs(bindings)

}

func handleKeyboardInputs(bindings *Bindings) {

	if bindings.NumberHandler != nil {
		for number, keys := range numberKeys {
			for _, key := range keys {
				if isKeyJustPressed(key, bindings.RepeatKeys) {
					bindings.NumberHandler(number)
				}
			}
		}
	}

	if bindings.TextHandler != nil {
		for rune, keys := range textKeys {
			for _, key := range keys {
				if isKeyJustPressed(key, bindings.RepeatKeys) {
					bindings.TextHandler(rune)
				}
			}
		}
	}

	if bindings.GlobalHandler != nil {
		for key := ebiten.KeyA; key <= ebiten.KeyMeta; key++ {
			if isKeyJustPressed(key, bindings.RepeatKeys) {
				bindings.GlobalHandler(key, -1, -1, 0)
			}
		}
		for id := range gamepads.gamepadIDs {
			if !ebiten.IsStandardGamepadLayoutAvailable(id) {
				continue
			}
			for btn := ebiten.StandardGamepadButtonRightBottom; btn <= ebiten.StandardGamepadButtonMax; btn++ {
				if inpututil.IsStandardGamepadButtonJustPressed(id, btn) {
					bindings.GlobalHandler(-1, btn, -1, 0)
				}
			}
			for axis := ebiten.StandardGamepadAxis(0); axis <= ebiten.StandardGamepadAxis(3); axis++ {
				v := ebiten.StandardGamepadAxisValue(id, axis)
				if math.Abs(v) >= AxisThreshold {
					sign := math.Signbit(v)
					if sign {
						bindings.GlobalHandler(-1, -1, axis, -1)
					} else {
						bindings.GlobalHandler(-1, -1, axis, 1)
					}

				}
			}
		}
	}

	for _, group := range bindings.Groups {
		for _, binding := range group {
			key := binding.Key()
			if binding.OnPressed != nil {
				if isKeyJustPressed(key, bindings.RepeatKeys) {
					binding.OnPressed()
				}
			}
			if binding.OnReleased != nil {
				if inpututil.IsKeyJustReleased(key) {
					binding.OnReleased()
				}
			}
		}
	}
}

func isKeyJustPressed(key ebiten.Key, repeating bool) bool {
	const (
		delay    = 15
		interval = 3
	)
	if !repeating {
		return inpututil.IsKeyJustPressed(key)
	} else {
		d := inpututil.KeyPressDuration(key)
		return d == 1 || (d >= delay && (d-delay)%interval == 0)
	}
}

func handleControllerInputs(bindings *Bindings) {
	for id := range gamepads.gamepadIDs {
		if !ebiten.IsStandardGamepadLayoutAvailable(id) {
			continue
		}
		for _, group := range bindings.Groups {
			for _, binding := range group {
				hasButton, button := binding.ControllerButton()
				hasAxis, axis, sign := binding.ControllerAxis()
				if !hasButton && !hasAxis {
					continue
				}
				if binding.OnPressed != nil {
					if hasButton && inpututil.IsStandardGamepadButtonJustPressed(id, button) {
						binding.OnPressed()
					}
				}
				if binding.OnReleased != nil {
					if hasButton && inpututil.IsStandardGamepadButtonJustReleased(id, button) {
						binding.OnReleased()
					}
				}
				v := ebiten.StandardGamepadAxisValue(id, axis)
				axisId := float64(axis+1) * sign
				if math.Abs(v) >= AxisThreshold && math.Signbit(sign) == math.Signbit(v) && hasAxis && !gamepads.gamepadAxis[id][axisId] {
					gamepads.gamepadAxis[id][axisId] = true
					if binding.OnPressed != nil {
						binding.OnPressed()
					}
				} else if gamepads.gamepadAxis[id][axisId] && math.Abs(v) < AxisThreshold {
					if binding.OnReleased != nil {
						binding.OnReleased()
					}
					gamepads.gamepadAxis[id][axisId] = false
				}
			}
		}
	}
}

func handleControllerEvents() {
	if gamepads.gamepadIDs == nil {
		gamepads.gamepadIDs = map[ebiten.GamepadID]struct{}{}
		gamepads.gamepadAxis = map[ebiten.GamepadID]map[float64]bool{}
	}
	gamepads.gamepadIDsBuf = inpututil.AppendJustConnectedGamepadIDs(gamepads.gamepadIDsBuf[:0])
	for _, id := range gamepads.gamepadIDsBuf {
		log.Printf("gamepad connected: id: %d, SDL ID: %s", id, ebiten.GamepadSDLID(id))
		gamepads.gamepadIDs[id] = struct{}{}
		gamepads.gamepadAxis[id] = map[float64]bool{}
	}
	for id := range gamepads.gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			log.Printf("gamepad disconnected: id: %d, SDL ID: %s", id, ebiten.GamepadSDLID(id))
			delete(gamepads.gamepadIDs, id)
		}
	}
}
