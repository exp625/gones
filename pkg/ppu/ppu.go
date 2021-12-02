package ppu

import (
	"github.com/exp625/gones/pkg/bus"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type PPU struct {
	Bus bus.Bus

	ScanLine   uint16
	Position   uint16
	FrameCount uint64

	Pallet [0x40][8]color.Color

	// Registers
	ppuctrl        uint8
	ppumask        uint8
	ppustatus      uint8
	oamaddr        uint8
	oamdata        uint8
	ppuscrollx     uint8
	ppuscrolly     uint8
	ppuscrolltemp  uint8
	ppuscrollwrite bool
	ppuaddr        uint16
	ppuaddrwrite   bool
	ppuaddrtemp    uint8
	ppudata        uint8
	oamdma         uint8

	// oam
	OAM [0xFF]byte
}

func New() *PPU {
	p := &PPU{}
	p.generatePallet()
	return p
}

func (ppu *PPU) Clock() {
	if ppu.ScanLine < 261 {
		if ppu.Position < 340 {
			ppu.Position++
		} else {
			ppu.Position = 0
			ppu.ScanLine++
		}
	} else {
		ppu.Position = 0
		ppu.ScanLine = 0
		ppu.FrameCount++
		if ppu.FrameCount%2 != 0 {
			ppu.Position++
		}
	}

	if ppu.ScanLine == 241 && ppu.Position == 1 {
		ppu.ppustatus |= 0b10000000
		if ppu.ppuctrl >> 7 & 0x1 == 1 {
			ppu.Bus.NMI()
		}
	}

	if ppu.ScanLine == 261 && ppu.Position == 1 {
		ppu.ppustatus &= 0b00011111
	}
}

func (ppu *PPU) Reset() {
	ppu.ScanLine = 0
	ppu.Position = 0
	ppu.FrameCount = 0

	ppu.ppuctrl = 0
	ppu.ppumask = 0
	ppu.ppustatus = 0
	ppu.oamaddr = 0
	ppu.oamdata = 0
	ppu.ppuscrollx = 0
	ppu.ppuscrolly = 0
	ppu.ppuscrolltemp = 0
	ppu.ppuscrollwrite = false
	ppu.ppuaddr = 0
	ppu.ppuaddrwrite = false
	ppu.ppuaddrtemp = 0
	ppu.ppudata = 0
	ppu.oamdma = 0
}

func (ppu *PPU) CPURead(location uint16) (bool, uint8) {
	if location >= 0x2000 && location <= 0x3FFF {
		switch (location - 0x2000) % 0x8 {
		case 0:
			return true, 0
		case 1:
			return true, 0
		case 2:
			ret := ppu.ppustatus
			ppu.ppustatus &= 0b01111111
			return true, ret
		case 3:
			return true, ppu.oamaddr
		case 4:
			return true, ppu.OAM[ppu.oamdata]
		case 5:
			return true, 0
		case 6:
			return true, ppu.ppudata
		case 7:
			return true, ppu.Bus.PPURead(ppu.ppuaddr)

		}
	}
	return false, 0
}

func (ppu *PPU) CPUWrite(location uint16, data uint8) {
	switch (location - 0x2000) % 0x8 {
	case 0:
		ppu.ppuctrl = data
	case 1:
		ppu.ppumask = data
	case 2:
		// Do nothing
	case 3:
		// Do nothing
	case 4:
		ppu.OAM[ppu.oamaddr] = data
		ppu.oamaddr++
	case 5:
		if !ppu.ppuscrollwrite {
			ppu.ppuscrolltemp = data
			ppu.ppuscrollwrite = true
		} else {
			ppu.ppuscrollx = ppu.ppuscrolltemp
			ppu.ppuscrolly = data
			ppu.ppuscrollwrite = false
		}
	case 6:
		if !ppu.ppuaddrwrite {
			ppu.ppuaddrtemp = data
			ppu.ppuaddrwrite = true
		} else {
			ppu.ppuaddr = (uint16(ppu.ppuaddrtemp) << 8) | uint16(data)
			ppu.ppuaddrwrite = false
		}
	case 7:
		ppu.Bus.PPUWrite(ppu.ppuaddr, data)
		if (ppu.ppuctrl >> 2 & 0x1) == 1 {
			ppu.ppuaddr += 32
		} else {
			ppu.ppuaddr++
		}
	}
}

func (ppu *PPU) DrawNametableInBW(table int) *ebiten.Image {
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
			tileByte := uint16(ppu.Bus.PPURead(nameTableBaseAddress + nameTableOffset + tileIndex))
			// RRRRCCCC
			// Background pattern table address (0: $0000; 1: $1000)
			backgroundTable := uint16(ppu.ppuctrl >> 4 & 0x1)

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
				plane0 := ppu.Bus.PPURead(addressPlane0)
				plane1 := ppu.Bus.PPURead(addressPlane1)

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
func (ppu *PPU) DrawNametableInColor(table int) *ebiten.Image {
	width := 32 * 8  // 256
	height := 30 * 8 // 240

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	colorEmpasis := ppu.ppumask >> 5
	// greyscale := ppu.ppumask & 0x1
	universalBackgroundColor := ppu.Pallet[ppu.Bus.PPURead(0x3F00) % 0x40][colorEmpasis]
	pallets := [4][4]color.Color{{
		universalBackgroundColor,
		ppu.Pallet[ppu.Bus.PPURead(0x3F01) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F02) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F03) % 0x40][colorEmpasis],
	},{
		universalBackgroundColor,
		ppu.Pallet[ppu.Bus.PPURead(0x3F05) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F06) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F07) % 0x40][colorEmpasis],
	},{
		universalBackgroundColor,
		ppu.Pallet[ppu.Bus.PPURead(0x3F09) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F0A) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F0B) % 0x40][colorEmpasis],
	},{
		universalBackgroundColor,
		ppu.Pallet[ppu.Bus.PPURead(0x3F0D) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F0E) % 0x40][colorEmpasis],
		ppu.Pallet[ppu.Bus.PPURead(0x3F0F) % 0x40][colorEmpasis],
	}}

	for row := uint16(0); row < 30; row++ {
		for tile := uint16(0); tile < 32; tile++ {
			// Get tile byte
			const nameTableBaseAddress = 0x2000
			const attributeTableBaseAddress = 0x23C0
			nameTableOffset := uint16(table * 0x400)
			tileIndex := row*32 + tile
			tileByte := uint16(ppu.Bus.PPURead(nameTableBaseAddress + nameTableOffset + tileIndex))
			// RRRRCCCC
			// Background pattern table address (0: $0000; 1: $1000)
			backgroundTable := uint16(ppu.ppuctrl >> 4 & 0x1)

			// Get assigned pallet
			pallet := ppu.Bus.PPURead((tile / 2) + (row / 2) * 8)
			palletIndex := uint8(0)
			switch {
			case tile % 2 == 0 && row % 2 == 0:
				palletIndex = pallet & 0x2
			case tile % 2 == 1 && row % 2 == 0:
				palletIndex = pallet >> 2 & 0x2
			case tile % 2 == 0 && row % 2 == 1:
				palletIndex = pallet >> 4 & 0x2
			case tile % 2 == 1 && row % 2 == 1:
				palletIndex = pallet >> 6 & 0x2
			}


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
				plane0 := ppu.Bus.PPURead(addressPlane0)
				plane1 := ppu.Bus.PPURead(addressPlane1)

				for tileX := uint16(0); tileX < 8; tileX++ {
					colorIndex := ((plane1>>(7-tileX))&0x01 << 1) | (plane0>>(7-tileX))&0x01
					c := pallets[palletIndex][colorIndex]

					imgX := int(tile*8 + tileX)
					imgY := int(row*8 + tileY)
					img.Set(imgX, imgY, c)
				}
			}
		}
	}

	return ebiten.NewImageFromImage(img)
}
