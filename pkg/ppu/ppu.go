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
	Control    uint8
	Mask       uint8
	Status     uint8
	OamAddress uint8
	OamData    uint8
	ScrollX    uint8
	ScrollY    uint8
	Address    uint16
	Data       uint8
	OamDma     uint8

	// Latch
	GenLatch       uint8
	AddrLatch      uint8
	addrLatchWrite bool

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
	ppu.AddrLatch = 0
	ppu.addrLatchWrite = false
	ppu.Address = 0
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
			return true, ppu.GenLatch
		case 1:
			return true, ppu.GenLatch
		case 2:
			ret := ppu.Status&0b11100000 | ppu.GenLatch&0b00011111
			ppu.Status &= 0b01111111
			ppu.GenLatch = ret
			return true, ret
		case 3:
			return true, ppu.GenLatch
		case 4:
			ret := ppu.OAM[ppu.OamAddress]
			ppu.GenLatch = ret
			return true, ret
		case 5:
			return true, ppu.GenLatch
		case 6:
			return true, ppu.GenLatch
		case 7:
			ret := ppu.Bus.PPURead(ppu.Address)
			ppu.GenLatch = ret
			if location >= 0x3F00 {
				ppu.GenLatch = ppu.Bus.PPUReadRam(ppu.Address)
			}
			if (ppu.Control >> 2 & 0x1) == 1 {
				ppu.Address += 32
			} else {
				ppu.Address++
			}
			return true, ret
		}
	}
	return false, 0
}

func (ppu *PPU) CPUWrite(location uint16, data uint8) {
	ppu.GenLatch = data
	switch (location - 0x2000) % 0x8 {
	case 0:
		ppu.Control = data
	case 1:
		ppu.Mask = data
	case 2:
		// Do nothing
	case 3:
		ppu.OamAddress = data
	case 4:
		ppu.OAM[ppu.OamAddress] = data
		ppu.OamAddress++
	case 5:
		if !ppu.addrLatchWrite {
			ppu.AddrLatch = data
			ppu.addrLatchWrite = true
		} else {
			ppu.ScrollX = ppu.AddrLatch
			ppu.ScrollY = data
			ppu.addrLatchWrite = false
		}
	case 6:
		if !ppu.addrLatchWrite {
			ppu.AddrLatch = data
			ppu.addrLatchWrite = true
		} else {
			ppu.Address = (uint16(ppu.AddrLatch) << 8) | uint16(data)%0x3FFF
			ppu.addrLatchWrite = false
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
