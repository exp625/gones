package ppu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"io/ioutil"
	"log"
)

// Pallet created using https://bisqwit.iki.fi/utils/nespalette.php
func (ppu *PPU) generatePallet() {
	bytes, err := ioutil.ReadFile("./palette.pal")
	if err != nil {
		log.Fatal(err)
	}
	ppu.Pallet = [0x40][8]color.Color{}
	pointer := 0
	for emphasis := 0; emphasis < 8; emphasis++ {
		for index := 0; index < 0x40; index++ {
			ppu.Pallet[index][emphasis] = color.RGBA{R: bytes[pointer], G: bytes[pointer+1], B: bytes[pointer+2], A: 255}
			pointer += 3
		}
	}
}

func (ppu *PPU) DrawPalette() *ebiten.Image {
	width := 16
	height := 4 * 8

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for emphasis := 0; emphasis < 8; emphasis++ {
		for index := 0; index < 0x40; index++ {
			img.Set(index % 0x10, emphasis * 4 + index / 0x10, ppu.Pallet[index][emphasis])
		}
	}
	return ebiten.NewImageFromImage(img)
}
