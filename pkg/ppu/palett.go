package ppu

import (
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
