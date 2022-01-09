package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"math"
)

const AxisThreshold = 0.2

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
	for _, group := range bindings.Groups {
		for _, binding := range group {
			key := binding.Key()
			if binding.OnPressed != nil {
				if inpututil.IsKeyJustPressed(key) {
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
					binding.OnPressed()
				} else if gamepads.gamepadAxis[id][axisId] && math.Abs(v) < AxisThreshold {
					binding.OnReleased()
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
