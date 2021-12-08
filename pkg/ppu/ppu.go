package ppu

import (
	"github.com/exp625/gones/pkg/bus"
	"image/color"
)

type PPU struct {
	Bus bus.Bus

	ScanLine   uint16
	Position   uint16
	FrameCount uint64

	Palette [0x40][8]color.Color

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
	OAM [256]uint8

	PaletteRAM [0x20]uint8
}

func New() *PPU {
	p := &PPU{}
	p.GeneratePalette()
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
		if ppu.ppuctrl>>7&0x1 == 1 {
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

	for i := 0; i < 256; i++ {
		ppu.OAM[i] = 0
	}
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
			if !ppu.Bus.Debugging() {
				ppu.ppustatus &= 0b01111111
			}
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

func (ppu *PPU) DMAWrite(data uint8) {
	ppu.OAM[ppu.oamaddr] = data
	ppu.oamaddr++
}
