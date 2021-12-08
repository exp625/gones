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
	Control      uint8
	Mask         uint8
	Status       uint8
	OamAddress   uint8
	OamData      uint8
	ScrollX      uint8
	ScrollY      uint8
	scrollTemp   uint8
	scrollWrite  bool
	Address      uint16
	addressWrite bool
	addressTemp  uint8
	Data         uint8
	OamDma       uint8

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
		ppu.Status |= 0b10000000
		if ppu.Control>>7&0x1 == 1 {
			ppu.Bus.NMI()
		}
	}

	if ppu.ScanLine == 261 && ppu.Position == 1 {
		ppu.Status &= 0b00011111
	}
}

func (ppu *PPU) Reset() {
	ppu.ScanLine = 0
	ppu.Position = 0
	ppu.FrameCount = 0

	ppu.Control = 0
	ppu.Mask = 0
	ppu.Status = 0
	ppu.OamAddress = 0
	ppu.OamData = 0
	ppu.ScrollX = 0
	ppu.ScrollY = 0
	ppu.scrollTemp = 0
	ppu.scrollWrite = false
	ppu.Address = 0
	ppu.addressWrite = false
	ppu.addressTemp = 0
	ppu.Data = 0
	ppu.OamDma = 0

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
			ret := ppu.Status
			ppu.Status &= 0b01111111
			return true, ret
		case 3:
			return true, ppu.OamAddress
		case 4:
			return true, ppu.OAM[ppu.OamData]
		case 5:
			return true, 0
		case 6:
			return true, ppu.Data
		case 7:
			ret := ppu.Bus.PPURead(ppu.Address)
			//if (ppu.Control >> 2 & 0x1) == 1 {
			//	ppu.Address += 32
			//} else {
			//	ppu.Address++
			//}
			return true, ret
		}
	}
	return false, 0
}

func (ppu *PPU) CPUWrite(location uint16, data uint8) {
	switch (location - 0x2000) % 0x8 {
	case 0:
		ppu.Control = data
	case 1:
		ppu.Mask = data
	case 2:
		// Do nothing
	case 3:
		// Do nothing
	case 4:
		ppu.OAM[ppu.OamAddress] = data
		ppu.OamAddress++
	case 5:
		if !ppu.scrollWrite {
			ppu.scrollTemp = data
			ppu.scrollWrite = true
		} else {
			ppu.ScrollX = ppu.scrollTemp
			ppu.ScrollY = data
			ppu.scrollWrite = false
		}
	case 6:
		if !ppu.addressWrite {
			ppu.addressTemp = data
			ppu.addressWrite = true
		} else {
			ppu.Address = (uint16(ppu.addressTemp) << 8) | uint16(data)
			ppu.addressWrite = false
		}
	case 7:
		ppu.Bus.PPUWrite(ppu.Address, data)
		if (ppu.Control >> 2 & 0x1) == 1 {
			ppu.Address += 32
		} else {
			ppu.Address++
		}
	}
}

func (ppu *PPU) DMAWrite(data uint8) {
	ppu.OAM[ppu.OamAddress] = data
	ppu.OamAddress++
}
