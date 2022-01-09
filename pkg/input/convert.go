package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"strings"
)

func keyNameToKeyCode(name string) (ebiten.Key, bool) {
	switch strings.ToLower(name) {
	case "0":
		return ebiten.Key0, true
	case "1":
		return ebiten.Key1, true
	case "2":
		return ebiten.Key2, true
	case "3":
		return ebiten.Key3, true
	case "4":
		return ebiten.Key4, true
	case "5":
		return ebiten.Key5, true
	case "6":
		return ebiten.Key6, true
	case "7":
		return ebiten.Key7, true
	case "8":
		return ebiten.Key8, true
	case "9":
		return ebiten.Key9, true
	case "a":
		return ebiten.KeyA, true
	case "b":
		return ebiten.KeyB, true
	case "c":
		return ebiten.KeyC, true
	case "d":
		return ebiten.KeyD, true
	case "e":
		return ebiten.KeyE, true
	case "f":
		return ebiten.KeyF, true
	case "g":
		return ebiten.KeyG, true
	case "h":
		return ebiten.KeyH, true
	case "i":
		return ebiten.KeyI, true
	case "j":
		return ebiten.KeyJ, true
	case "k":
		return ebiten.KeyK, true
	case "l":
		return ebiten.KeyL, true
	case "m":
		return ebiten.KeyM, true
	case "n":
		return ebiten.KeyN, true
	case "o":
		return ebiten.KeyO, true
	case "p":
		return ebiten.KeyP, true
	case "q":
		return ebiten.KeyQ, true
	case "r":
		return ebiten.KeyR, true
	case "s":
		return ebiten.KeyS, true
	case "t":
		return ebiten.KeyT, true
	case "u":
		return ebiten.KeyU, true
	case "v":
		return ebiten.KeyV, true
	case "w":
		return ebiten.KeyW, true
	case "x":
		return ebiten.KeyX, true
	case "y":
		return ebiten.KeyY, true
	case "z":
		return ebiten.KeyZ, true
	case "alt":
		return ebiten.KeyAlt, true
	case "altleft":
		return ebiten.KeyAltLeft, true
	case "altright":
		return ebiten.KeyAltRight, true
	case "apostrophe":
		return ebiten.KeyApostrophe, true
	case "arrowdown":
		return ebiten.KeyArrowDown, true
	case "arrowleft":
		return ebiten.KeyArrowLeft, true
	case "arrowright":
		return ebiten.KeyArrowRight, true
	case "arrowup":
		return ebiten.KeyArrowUp, true
	case "backquote":
		return ebiten.KeyBackquote, true
	case "backslash":
		return ebiten.KeyBackslash, true
	case "backspace":
		return ebiten.KeyBackspace, true
	case "bracketleft":
		return ebiten.KeyBracketLeft, true
	case "bracketright":
		return ebiten.KeyBracketRight, true
	case "capslock":
		return ebiten.KeyCapsLock, true
	case "comma":
		return ebiten.KeyComma, true
	case "contextmenu":
		return ebiten.KeyContextMenu, true
	case "control":
		return ebiten.KeyControl, true
	case "controlleft":
		return ebiten.KeyControlLeft, true
	case "controlright":
		return ebiten.KeyControlRight, true
	case "delete":
		return ebiten.KeyDelete, true
	case "digit0":
		return ebiten.KeyDigit0, true
	case "digit1":
		return ebiten.KeyDigit1, true
	case "digit2":
		return ebiten.KeyDigit2, true
	case "digit3":
		return ebiten.KeyDigit3, true
	case "digit4":
		return ebiten.KeyDigit4, true
	case "digit5":
		return ebiten.KeyDigit5, true
	case "digit6":
		return ebiten.KeyDigit6, true
	case "digit7":
		return ebiten.KeyDigit7, true
	case "digit8":
		return ebiten.KeyDigit8, true
	case "digit9":
		return ebiten.KeyDigit9, true
	case "down":
		return ebiten.KeyDown, true
	case "end":
		return ebiten.KeyEnd, true
	case "enter":
		return ebiten.KeyEnter, true
	case "equal":
		return ebiten.KeyEqual, true
	case "escape":
		return ebiten.KeyEscape, true
	case "f1":
		return ebiten.KeyF1, true
	case "f2":
		return ebiten.KeyF2, true
	case "f3":
		return ebiten.KeyF3, true
	case "f4":
		return ebiten.KeyF4, true
	case "f5":
		return ebiten.KeyF5, true
	case "f6":
		return ebiten.KeyF6, true
	case "f7":
		return ebiten.KeyF7, true
	case "f8":
		return ebiten.KeyF8, true
	case "f9":
		return ebiten.KeyF9, true
	case "f10":
		return ebiten.KeyF10, true
	case "f11":
		return ebiten.KeyF11, true
	case "f12":
		return ebiten.KeyF12, true
	case "graveaccent":
		return ebiten.KeyGraveAccent, true
	case "home":
		return ebiten.KeyHome, true
	case "insert":
		return ebiten.KeyInsert, true
	case "kp0":
		return ebiten.KeyKP0, true
	case "kp1":
		return ebiten.KeyKP1, true
	case "kp2":
		return ebiten.KeyKP2, true
	case "kp3":
		return ebiten.KeyKP3, true
	case "kp4":
		return ebiten.KeyKP4, true
	case "kp5":
		return ebiten.KeyKP5, true
	case "kp6":
		return ebiten.KeyKP6, true
	case "kp7":
		return ebiten.KeyKP7, true
	case "kp8":
		return ebiten.KeyKP8, true
	case "kp9":
		return ebiten.KeyKP9, true
	case "kpdecimal":
		return ebiten.KeyKPDecimal, true
	case "kpdivide":
		return ebiten.KeyKPDivide, true
	case "kpenter":
		return ebiten.KeyKPEnter, true
	case "kpequal":
		return ebiten.KeyKPEqual, true
	case "kpmultiply":
		return ebiten.KeyKPMultiply, true
	case "kpsubtract":
		return ebiten.KeyKPSubtract, true
	case "left":
		return ebiten.KeyLeft, true
	case "leftbracket":
		return ebiten.KeyLeftBracket, true
	case "menu":
		return ebiten.KeyMenu, true
	case "meta":
		return ebiten.KeyMeta, true
	case "metaleft":
		return ebiten.KeyMetaLeft, true
	case "metaright":
		return ebiten.KeyMetaRight, true
	case "minus":
		return ebiten.KeyMinus, true
	case "numlock":
		return ebiten.KeyNumLock, true
	case "numpad0":
		return ebiten.KeyNumpad0, true
	case "numpad1":
		return ebiten.KeyNumpad1, true
	case "numpad2":
		return ebiten.KeyNumpad2, true
	case "numpad3":
		return ebiten.KeyNumpad3, true
	case "numpad4":
		return ebiten.KeyNumpad4, true
	case "numpad5":
		return ebiten.KeyNumpad5, true
	case "numpad6":
		return ebiten.KeyNumpad6, true
	case "numpad7":
		return ebiten.KeyNumpad7, true
	case "numpad8":
		return ebiten.KeyNumpad8, true
	case "numpad9":
		return ebiten.KeyNumpad9, true
	case "numpadadd":
		return ebiten.KeyNumpadAdd, true
	case "numpaddecimal":
		return ebiten.KeyNumpadDecimal, true
	case "numpaddivide":
		return ebiten.KeyNumpadDivide, true
	case "numpadenter":
		return ebiten.KeyNumpadEnter, true
	case "numpadequal":
		return ebiten.KeyNumpadEqual, true
	case "numpadmultiply":
		return ebiten.KeyNumpadMultiply, true
	case "numpadsubtract":
		return ebiten.KeyNumpadSubtract, true
	case "pagedown":
		return ebiten.KeyPageDown, true
	case "pageup":
		return ebiten.KeyPageUp, true
	case "pause":
		return ebiten.KeyPause, true
	case "period":
		return ebiten.KeyPeriod, true
	case "printscreen":
		return ebiten.KeyPrintScreen, true
	case "quote":
		return ebiten.KeyQuote, true
	case "right":
		return ebiten.KeyRight, true
	case "rightbracket":
		return ebiten.KeyRightBracket, true
	case "scrolllock":
		return ebiten.KeyScrollLock, true
	case "semicolon":
		return ebiten.KeySemicolon, true
	case "shift":
		return ebiten.KeyShift, true
	case "shiftleft":
		return ebiten.KeyShiftLeft, true
	case "shiftright":
		return ebiten.KeyShiftRight, true
	case "slash":
		return ebiten.KeySlash, true
	case "space":
		return ebiten.KeySpace, true
	case "tab":
		return ebiten.KeyTab, true
	case "up":
		return ebiten.KeyUp, true
	}
	return 0, false
}
