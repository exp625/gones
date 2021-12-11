package file_explorer

import "github.com/hajimehoshi/ebiten/v2"

// Common keys that represent a letter that appear in a filename (a-z and 0-9).
var fileNameKeys map[string]ebiten.Key

func init() {
	fileNameKeys = map[string]ebiten.Key{
		"A": ebiten.KeyA,
		"B": ebiten.KeyB,
		"C": ebiten.KeyC,
		"D": ebiten.KeyD,
		"E": ebiten.KeyE,
		"F": ebiten.KeyF,
		"G": ebiten.KeyG,
		"H": ebiten.KeyH,
		"I": ebiten.KeyI,
		"J": ebiten.KeyJ,
		"K": ebiten.KeyK,
		"L": ebiten.KeyL,
		"M": ebiten.KeyM,
		"N": ebiten.KeyN,
		"O": ebiten.KeyO,
		"P": ebiten.KeyP,
		"Q": ebiten.KeyQ,
		"R": ebiten.KeyR,
		"S": ebiten.KeyS,
		"T": ebiten.KeyT,
		"U": ebiten.KeyU,
		"V": ebiten.KeyV,
		"W": ebiten.KeyW,
		"X": ebiten.KeyX,
		"Y": ebiten.KeyY,
		"Z": ebiten.KeyZ,
		"0": ebiten.KeyDigit0,
		"1": ebiten.KeyDigit1,
		"2": ebiten.KeyDigit2,
		"3": ebiten.KeyDigit3,
		"4": ebiten.KeyDigit4,
		"5": ebiten.KeyDigit5,
		"6": ebiten.KeyDigit6,
		"7": ebiten.KeyDigit7,
		"8": ebiten.KeyDigit8,
		"9": ebiten.KeyDigit9,
	}
}
