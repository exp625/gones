package ppu

import (
	_ "embed"
	"image/color"
)

//go:embed palette.pal
var palette []byte

// GeneratePalette generates a palette for the PPU
// The palette is created using https://bisqwit.iki.fi/utils/nespalette.php
func (ppu *PPU) GeneratePalette() {
	ppu.Palette = [0x40][8]color.Color{}
	pointer := 0
	for emphasis := 0; emphasis < 8; emphasis++ {
		for index := 0; index < 0x40; index++ {
			ppu.Palette[index][emphasis] = color.RGBA{
				R: palette[pointer],
				G: palette[pointer+1],
				B: palette[pointer+2],
				A: 255,
			}
			pointer += 3
		}
	}
}
