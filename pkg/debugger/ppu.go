package debugger

import (
"fmt"
"github.com/exp625/gones/internal/plz"
"github.com/exp625/gones/internal/textutil"
"github.com/hajimehoshi/ebiten/v2"
"golang.org/x/image/colornames"
"image"
"image/color"
)

func (nes *Debugger) DrawPatternTable(table int) *ebiten.Image {
	width := 128
	height := 128

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// DCBA98 76543210
	// ---------------
	// 0HRRRR CCCCPTTT
	// |||||| |||||+++- T: Fine Y offset, the row number within a tile
	// |||||| ||||+---- P: Bit plane (0: "lower"; 1: "upper")
	// |||||| ++++----- 6502: Tile column
	// ||++++---------- R: Tile row
	// |+-------------- H: Half of sprite table (0: "left"; 1: "right")
	// +--------------- 0: Pattern table is at $0000-$1FFF

	for row := 0; row < 16; row++ {
		for tile := 0; tile < 16; tile++ {
			for tileY := 0; tileY < 8; tileY++ {
				addressPlane0 := uint16(table<<12 | row<<8 | tile<<4 | 0<<3 | tileY)
				addressPlane1 := uint16(table<<12 | row<<8 | tile<<4 | 1<<3 | tileY)
				plane0 := nes.PPURead(addressPlane0)
				plane1 := nes.PPURead(addressPlane1)

				colorEmphasis := nes.PPU.Mask >> 5
				pallets := [4]color.Color{
					nes.PPU.Palette[nes.PPURead(0x3F00)%0x40][colorEmphasis],
					nes.PPU.Palette[nes.PPURead(0x3F01)%0x40][colorEmphasis],
					nes.PPU.Palette[nes.PPURead(0x3F02)%0x40][colorEmphasis],
					nes.PPU.Palette[nes.PPURead(0x3F03)%0x40][colorEmphasis],
				}
				if nes.PPURead(0x3F00)%0x40 == 0 && nes.PPURead(0x3F01)%0x40 == 0 && nes.PPURead(0x3F02)%0x40 == 0 && nes.PPURead(0x3F03)%0x40 == 0 {
					// Pallet no initialized
					pallets = [4]color.Color{
						color.RGBA{0, 0, 0, 255},
						color.RGBA{100, 100, 100, 255},
						color.RGBA{150, 150, 150, 255},
						color.RGBA{255, 255, 255, 255},
					}
				}

				for tileX := 0; tileX < 8; tileX++ {
					colorIndex := ((plane1 >> (7 - tileX)) & 0x01 << 1) | (plane0>>(7-tileX))&0x01
					c := pallets[colorIndex]

					imgX := tile*8 + tileX
					imgY := row*8 + tileY
					img.Set(imgX, imgY, c)
				}
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func (nes *Debugger) DrawPPUInfo(t *textutil.Text) {
	plz.Just(fmt.Fprintf(t, "PPUCTRL: %02X \tGenerate NMI: %t  \t\t\t", nes.PPU.Control, (nes.PPU.Control>>7)&0x1 == 1))
	if (nes.PPU.Control>>6)&0x1 == 0 {
		plz.Just(fmt.Fprintf(t, "PPU Type: Master \t"))
	} else {
		plz.Just(fmt.Fprintf(t, "PPU Type: Slave \t"))
	}
	if (nes.PPU.Control>>5)&0x1 == 0 {
		plz.Just(fmt.Fprintf(t, "Sprite Size: 8x8 \n"))
	} else {
		plz.Just(fmt.Fprintf(t, "Sprite Size: 16x8 \n"))
	}
	if (nes.PPU.Control>>4)&0x1 == 0 {
		plz.Just(fmt.Fprintf(t, "\t\t\t\tPattern Table Addres: 0x0000 \t\t"))
	} else {
		plz.Just(fmt.Fprintf(t, "\t\t\t\tPattern Table Addres: 0x1000 \t\t"))
	}
	if (nes.PPU.Control>>3)&0x1 == 0 {
		plz.Just(fmt.Fprintf(t, "Sprite Table Addres: 0x0000 \n"))
	} else {
		plz.Just(fmt.Fprintf(t, "Sprite Table Addres: 0x1000 \n"))
	}
	if (nes.PPU.Control>>2)&0x1 == 0 {
		plz.Just(fmt.Fprintf(t, "\t\t\t\tVRAM Increment: 1 \t\t\t\t"))
	} else {
		plz.Just(fmt.Fprintf(t, "\t\t\t\tVRAM Increment: 32 \t\t\t\t"))
	}
	switch nes.PPU.Control & 0b11 {
	case 0:
		plz.Just(fmt.Fprintf(t, "Base Nametable: 0x2000 \n\n"))
	case 1:
		plz.Just(fmt.Fprintf(t, "Base Nametable: 0x2400 \n\n"))
	case 2:
		plz.Just(fmt.Fprintf(t, "Base Nametable: 0x2800 \n\n"))
	case 3:
		plz.Just(fmt.Fprintf(t, "Base Nametable: 0x2C00 \n\n"))
	}
	plz.Just(fmt.Fprintf(t, "PPUMASK: %02X \tEmphasize BGR: %03b \t\t\t\tShow Sprites: %t \n\t\t\t\tShow Background: %t \t\t\tGreyscale: %t\n \t\t\t\t\tSprite leftmost 8px: %t \t\tBackground leftmost 8px: %t \n\n",
		nes.PPU.Mask, (nes.PPU.Control>>5)&0b111, (nes.PPU.Control>>4)&0x1 == 1, (nes.PPU.Control>>3)&0x1 == 1, nes.PPU.Control&0x1 == 1, (nes.PPU.Control>>2)&0x1 == 1, (nes.PPU.Control>>1)&0x1 == 1))
	plz.Just(fmt.Fprintf(t, "PPUSTATUS: %02X \tVBLANK: %t \tSprite 0 Hit: %t \t Sprite Overflow: %t \n",
		nes.PPU.Status, (nes.PPU.Status>>7)&0x1 == 1, (nes.PPU.Status>>6)&0x1 == 1, (nes.PPU.Status>>5)&0x1 == 1))
	plz.Just(fmt.Fprintf(t, "\t\t\t\tX-Scroll %02X \tY-Scroll: %02X \n",
		nes.PPU.ScrollX, nes.PPU.ScrollY))

}

func (nes *Debugger) DrawOAM(t *textutil.Text) {
	t.Color(colornames.White)
	plz.Just(fmt.Fprintf(t, "OAM: 0x%02X\n   ", nes.PPU.OamAddress))
	t.Color(colornames.Yellow)
	for i := 0; i <= 0xF; i++ {
		plz.Just(fmt.Fprintf(t, "%02X ", uint16(i)))
	}

	for i := 0x00; i <= 0xFF; i++ {
		if i%16 == 0 {
			t.Color(colornames.Yellow)
			plz.Just(fmt.Fprintf(t, "\n%02X ", uint16(i&0xF0)))
		}
		if nes.PPU.OamAddress == uint8(i) {
			t.Color(colornames.Green)
		} else {
			t.Color(colornames.White)
		}
		plz.Just(fmt.Fprintf(t, "%02X ", nes.PPU.OAM[uint16(i)]))
	}
}

func (nes *Debugger) DrawPalettes() *ebiten.Image {
	width := 16
	height := 2

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Background Palettes
	for palette := uint16(0); palette < 4; palette++ {
		for index := uint16(0); index < 4; index++ {
			if index == 0 {
				img.Set(int(palette*4+index), 0, nes.PPU.Palette[nes.PPU.Bus.PPURead(0x3F00)][0])
			} else {
				img.Set(int(palette*4+index), 0, nes.PPU.Palette[nes.PPU.Bus.PPURead(0x3F00+index+palette*4)%0x40][0])
			}

		}
	}

	// Sprite Palettes
	for palette := uint16(0); palette < 8; palette++ {
		for index := uint16(0); index < 0x40; index++ {
			if index == 0 {
				img.Set(int(palette*4+index), 1, nes.PPU.Palette[nes.PPU.Bus.PPURead(0x3F00)][0])
			} else {
				img.Set(int(palette*4+index), 1, nes.PPU.Palette[nes.PPU.Bus.PPURead(0x3F10+index+palette*4)%0x40][0])
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func (nes *Debugger) DrawLoadedPalette() *ebiten.Image {
	width := 16
	height := 4 * 8

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for emphasis := 0; emphasis < 8; emphasis++ {
		for index := 0; index < 0x40; index++ {
			img.Set(index%0x10, emphasis*4+index/0x10, nes.PPU.Palette[index][emphasis])
		}
	}
	return ebiten.NewImageFromImage(img)
}

func (nes *Debugger) DrawNametableInBW(table int) *ebiten.Image {
	width := 32 * 8  // 256
	height := 30 * 8 // 240

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	for row := uint16(0); row < 30; row++ {
		for tile := uint16(0); tile < 32; tile++ {
			// Get tile byte
			const nameTableBaseAddress = 0x2000
			nameTableOffset := uint16(table * 0x400)
			tileIndex := row*32 + tile
			tileByte := uint16(nes.PPU.Bus.PPURead(nameTableBaseAddress+nameTableOffset+tileIndex))
			// RRRRCCCC
			// Background pattern table address (0: $0000; 1: $1000)
			backgroundTable := uint16(nes.PPU.Control >> 4 & 0x1)

			// Get tile byte
			// DCBA98 76543210
			// ---------------
			// 0HRRRR CCCCPTTT
			// |||||| |||||+++- T: Fine Y offset, the row number within a tile
			// |||||| ||||+---- P: Bit plane (0: "lower"; 1: "upper")
			// |||||| ++++----- C: Tile column
			// ||++++---------- R: Tile row
			// |+-------------- H: Half of sprite table (0: "left"; 1: "right")
			// +--------------- 0: Pattern table is at $0000-$1FFF
			for tileY := uint16(0); tileY < 8; tileY++ {
				addressPlane0 := backgroundTable<<12 | tileByte<<4 | 0<<3 | tileY
				addressPlane1 := backgroundTable<<12 | tileByte<<4 | 1<<3 | tileY
				plane0 := nes.PPU.Bus.PPURead(addressPlane0)
				plane1 := nes.PPU.Bus.PPURead(addressPlane1)

				for tileX := uint16(0); tileX < 8; tileX++ {
					var c color.Color
					if (plane0>>(7-tileX))&0x01 == 1 && (plane1>>(7-tileX))&0x01 == 1 {
						c = color.White
					} else if (plane1>>(7-tileX))&0x01 == 1 {
						c = color.Gray16{Y: 0xAAAA}
					} else if (plane0>>(7-tileX))&0x01 == 1 {
						c = color.Gray16{Y: 0x5555}
					} else {
						c = color.Gray16{Y: 0x1111}
					}
					imgX := int(tile*8 + tileX)
					imgY := int(row*8 + tileY)
					img.Set(imgX, imgY, c)
				}
			}
		}
	}

	return ebiten.NewImageFromImage(img)
}

func (nes *Debugger) DrawNametableInColor(table int) *ebiten.Image {
	width := 32 * 8  // 256
	height := 30 * 8 // 240

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	colorEmphasis := nes.PPU.Mask >> 5
	// greyscale := nes.PPU.Mask & 0x1

	for row := uint16(0); row < 30; row++ {
		for tile := uint16(0); tile < 32; tile++ {
			// Get tile byte
			const nameTableBaseAddress = 0x2000
			const attributeTableBaseAddress = 0x23C0
			nameTableOffset := uint16(table * 0x400)
			tileIndex := row*32 + tile
			tileByte := uint16(nes.PPU.Bus.PPURead(nameTableBaseAddress+nameTableOffset+tileIndex))
			// RRRRCCCC
			// Background pattern table address (0: $0000; 1: $1000)
			backgroundTable := uint16(nes.PPU.Control >> 4 & 0x1)

			// Get assigned attributeByte
			attributeIndex := (tile / 4) + (row/4)*8
			// tileByte = attributeIndex
			attributeByte := uint16(nes.PPU.Bus.PPURead(attributeTableBaseAddress+nameTableOffset+attributeIndex))
			attribute := uint16(0)

			switch {
			case tile%4 <= 1 && row%4 <= 1:
				attribute = attributeByte & 0b11
			case tile%4 >= 2 && row%4 <= 1:
				attribute = attributeByte >> 2 & 0b11
			case tile%4 <= 1 && row%4 >= 2:
				attribute = attributeByte >> 4 & 0b11
			case tile%4 >= 2 && row%4 >= 2:
				attribute = attributeByte >> 6 & 0b11
			}
			// tileByte = attribute
			// Get tile byte
			// DCBA98 76543210
			// ---------------
			// 0HRRRR CCCCPTTT
			// |||||| |||||+++- T: Fine Y offset, the row number within a tile
			// |||||| ||||+---- P: Bit plane (0: "lower"; 1: "upper")
			// |||||| ++++----- C: Tile column
			// ||++++---------- R: Tile row
			// |+-------------- H: Half of sprite table (0: "left"; 1: "right")
			// +--------------- 0: Pattern table is at $0000-$1FFF
			for tileY := uint16(0); tileY < 8; tileY++ {
				addressPlane0 := backgroundTable<<12 | tileByte<<4 | 0<<3 | tileY
				addressPlane1 := backgroundTable<<12 | tileByte<<4 | 1<<3 | tileY
				plane0 := nes.PPU.Bus.PPURead(addressPlane0)
				plane1 := nes.PPU.Bus.PPURead(addressPlane1)

				for tileX := uint16(0); tileX < 8; tileX++ {
					colorIndex := uint16(((plane1 >> (7 - tileX)) & 0x01 << 1) | (plane0>>(7-tileX))&0x01)
					c := nes.PPU.Palette[nes.PPU.Bus.PPURead(0x3F00+attribute*4+colorIndex)][colorEmphasis]

					imgX := int(tile*8 + tileX)
					imgY := int(row*8 + tileY)
					img.Set(imgX, imgY, c)
				}
			}
		}
	}

	return ebiten.NewImageFromImage(img)
}

func (nes *Debugger) DrawOAMSprites() *ebiten.Image {
	colorEmphasis := nes.PPU.Mask >> 5
	width := 32 * 8  // 256
	height := 30 * 8 // 240

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	spriteTable := uint16(0x0000)
	if nes.PPU.Control>>3&0x1 == 1 {
		spriteTable = 0x1000
	}

	for sprite := 0; sprite < 64; sprite++ {
		spriteYPos := nes.PPU.OAM[sprite]
		spriteIndex := uint16(nes.PPU.OAM[sprite+1])
		spriteAttributes := uint16(nes.PPU.OAM[sprite+2])
		spriteXPos := nes.PPU.OAM[sprite+3]

		paletteIndex := spriteAttributes&0b11 + 4
		flipHorizontal := spriteAttributes>>6&0b1 == 1
		flipVertical := spriteAttributes>>7&0b1 == 1

		// Get tile byte
		for tileY := uint16(0); tileY < 8; tileY++ {
			tileYFlipped := tileY
			if flipVertical {
				tileYFlipped = 7 - tileYFlipped
			}
			addressPlane0 := spriteTable<<12 | spriteIndex<<4 | 0<<3 | tileYFlipped
			addressPlane1 := spriteTable<<12 | spriteIndex<<4 | 1<<3 | tileYFlipped
			plane0 := nes.PPURead(addressPlane0)
			plane1 := nes.PPURead(addressPlane1)

			for tileX := uint16(0); tileX < 8; tileX++ {
				tileXFlipped := tileX
				if flipHorizontal {
					tileXFlipped = 7 - tileXFlipped
				}
				colorIndex := uint16(((plane1 >> (7 - tileXFlipped)) & 0x01 << 1) | (plane0>>(7-tileXFlipped))&0x01)
				c := nes.PPU.Palette[nes.PPURead(0x3F00+(paletteIndex*4+colorIndex))][colorEmphasis]

				imgX := int(spriteXPos + uint8(tileX))
				imgY := int(spriteYPos + uint8(tileY))
				img.Set(imgX, imgY, c)
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func (nes *Debugger) DrawGameDebug() *ebiten.Image {
	width := 32 * 8  // 256
	height := 30 * 8 // 240
	table := 0
	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	imgBackgroundColor := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	imgBackgroundSprite := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	imgBackgroundTiles := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	imgSprites := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	spriteTable := uint16(0x0000)
	if nes.PPU.Control>>3&0x1 == 1 {
		spriteTable = 0x1000
	}
	colorEmphasis := nes.PPU.Mask >> 5

	for row := uint16(0); row < 30; row++ {
		for tile := uint16(0); tile < 32; tile++ {
			// Get tile byte
			const nameTableBaseAddress = 0x2000
			const attributeTableBaseAddress = 0x23C0
			nameTableOffset := uint16(table * 0x400)
			tileIndex := row*32 + tile
			tileByte := uint16(nes.PPURead(nameTableBaseAddress+nameTableOffset+tileIndex))
			// RRRRCCCC
			// Background pattern table address (0: $0000; 1: $1000)
			backgroundTable := uint16(nes.PPU.Control >> 4 & 0x1)

			// Get assigned attributeByte
			attributeIndex := (tile / 4) + (row/4)*8
			// tileByte = attributeIndex
			attributeByte := uint16(nes.PPURead(attributeTableBaseAddress+nameTableOffset+attributeIndex))
			attribute := uint16(0)

			switch {
			case tile%4 <= 1 && row%4 <= 1:
				attribute = attributeByte & 0b11
			case tile%4 >= 2 && row%4 <= 1:
				attribute = attributeByte >> 2 & 0b11
			case tile%4 <= 1 && row%4 >= 2:
				attribute = attributeByte >> 4 & 0b11
			case tile%4 >= 2 && row%4 >= 2:
				attribute = attributeByte >> 6 & 0b11
			}
			// tileByte = attribute
			// Get tile byte
			// DCBA98 76543210
			// ---------------
			// 0HRRRR CCCCPTTT
			// |||||| |||||+++- T: Fine Y offset, the row number within a tile
			// |||||| ||||+---- P: Bit plane (0: "lower"; 1: "upper")
			// |||||| ++++----- C: Tile column
			// ||++++---------- R: Tile row
			// |+-------------- H: Half of sprite table (0: "left"; 1: "right")
			// +--------------- 0: Pattern table is at $0000-$1FFF
			for tileY := uint16(0); tileY < 8; tileY++ {
				addressPlane0 := backgroundTable<<12 | tileByte<<4 | 0<<3 | tileY
				addressPlane1 := backgroundTable<<12 | tileByte<<4 | 1<<3 | tileY
				plane0 := nes.PPURead(addressPlane0)
				plane1 := nes.PPURead(addressPlane1)

				for tileX := uint16(0); tileX < 8; tileX++ {
					colorIndex := uint16(((plane1 >> (7 - tileX)) & 0x01 << 1) | (plane0>>(7-tileX))&0x01)
					c := nes.PPU.Palette[nes.PPURead(0x3F00+attribute*4+colorIndex)][colorEmphasis]
					if colorIndex == 0 {
						c = color.RGBA{}
					}
					imgX := int(tile*8 + tileX)
					imgY := int(row*8 + tileY)
					imgBackgroundTiles.Set(imgX, imgY, c)
				}
			}
		}
	}
	for sprite := 0; sprite < 64; sprite++ {
		spriteYPos := nes.PPU.OAM[sprite]
		spriteIndex := uint16(nes.PPU.OAM[sprite+1])
		spriteAttributes := uint16(nes.PPU.OAM[sprite+2])
		spriteXPos := nes.PPU.OAM[sprite+3]

		paletteIndex := spriteAttributes&0b11 + 4
		flipHorizontal := spriteAttributes>>6&0b1 == 1
		flipVertical := spriteAttributes>>7&0b1 == 1
		priority := spriteAttributes>>5&0b1 == 0

		// Get tile byte
		for tileY := uint16(0); tileY < 8; tileY++ {
			tileYFlipped := tileY
			if flipVertical {
				tileYFlipped = 7 - tileYFlipped
			}
			addressPlane0 := spriteTable<<12 | spriteIndex<<4 | 0<<3 | tileYFlipped
			addressPlane1 := spriteTable<<12 | spriteIndex<<4 | 1<<3 | tileYFlipped
			plane0 := nes.PPURead(addressPlane0)
			plane1 := nes.PPURead(addressPlane1)

			for tileX := uint16(0); tileX < 8; tileX++ {
				tileXFlipped := tileX
				if flipHorizontal {
					tileXFlipped = 7 - tileXFlipped
				}
				colorIndex := uint16(((plane1 >> (7 - tileXFlipped)) & 0x01 << 1) | (plane0>>(7-tileXFlipped))&0x01)
				c := nes.PPU.Palette[nes.PPURead(0x3F00+(paletteIndex*4+colorIndex))][colorEmphasis]
				if colorIndex == 0 {
					c = color.RGBA{}
				}
				imgX := int(spriteXPos + uint8(tileX))
				imgY := int(spriteYPos + uint8(tileY))
				if priority {
					imgSprites.Set(imgX, imgY, c)
				} else {
					imgBackgroundSprite.Set(imgX, imgY, c)
				}
			}
		}
	}

	img := ebiten.NewImageFromImage(imgBackgroundColor)
	img.Fill(nes.PPU.Palette[nes.PPURead(0x3F00)][colorEmphasis])
	img.DrawImage(ebiten.NewImageFromImage(imgBackgroundSprite), nil)
	img.DrawImage(ebiten.NewImageFromImage(imgBackgroundTiles), nil)
	img.DrawImage(ebiten.NewImageFromImage(imgSprites), nil)
	return ebiten.NewImageFromImage(img)
}

