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
	Control    ControlRegister
	Mask       MaskRegister
	Status     StatusRegister
	OamAddress uint8
	OamData    uint8
	ScrollX    uint8
	ScrollY    uint8
	Address    AddressRegister
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

	// Set VBL Flag and trigger NMI on line 241 dot 1
	if ppu.ScanLine == 241 && ppu.Position == 1 {
		ppu.Status = ppu.Status | 0b10000000
		if ppu.Control>>7&0x1 == 1 {
			ppu.Bus.NMI()
		}
	}

	// Rendering

	// Advance counters
	if ppu.Position < 340 {
		ppu.Position++
	} else {
		if ppu.ScanLine < 261 {
			ppu.ScanLine++
			ppu.Position = 0
		} else {
			ppu.Position = 0
			ppu.ScanLine = 0
			ppu.FrameCount++
			if ppu.FrameCount%2 != 0 {
				ppu.Position++
			}
		}
	}

	// Clear Flags on line 261 dot 1
	if ppu.ScanLine == 261 && ppu.Position == 1 {
		ppu.Status = 0b00000000
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
			ret := uint8(ppu.Status) | ppu.GenLatch&0b00011111
			ppu.Status = ppu.Status & 0b01100000
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
			ret := ppu.Bus.PPURead(uint16(ppu.Address))
			ppu.GenLatch = ret
			if location >= 0x3F00 {
				ppu.GenLatch = ppu.Bus.PPUReadRam(uint16(ppu.Address))
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
		ppu.Control = ControlRegister(data)
	case 1:
		ppu.Mask = MaskRegister(data)
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
			ppu.Address = AddressRegister((uint16(ppu.AddrLatch) << 8) | uint16(data)%0x3FFF)
			ppu.addrLatchWrite = false
		}
	case 7:
		ppu.Bus.PPUWrite(uint16(ppu.Address), data)
		if (ppu.Control >> 2 & 0x1) == 1 {
			ppu.Address += 32
		} else {
			ppu.Address++
		}
	}
}
