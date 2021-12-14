package ppu

import (
	"github.com/exp625/gones/pkg/bus"
	"image/color"
)

type PPU struct {
	Bus bus.Bus

	ScanLine   uint16
	Dot        uint16
	FrameCount uint64

	Palette [0x40][8]color.Color

	// Registers
	Control    ControlRegister
	Mask       MaskRegister
	Status     StatusRegister
	OamAddress uint8
	OamData    uint8
	Data       uint8
	OamDma     uint8

	// Render Registers
	CurrentVRAMAddress   AddressRegister
	TemporaryVRAMAddress AddressRegister
	FineXScroll          uint8
	// Current write mode. 0 for first, 1 for second.
	WriteToggle uint8

	// Latch
	GenLatch  uint8
	ReadLatch uint8

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

	// Rendering

	if ppu.IsVisibleFrameOrPreRender() {
		// Increment vertical position on dot 256 of each scanline
		if ppu.Dot == 256 && ppu.Mask.ShowBackground() {
			ppu.IncrementVerticalPosition()
		}

		// Copy horizontal position on dot 257 of each scanline
		if ppu.Dot == 257 && ppu.Mask.ShowBackground() {
			ppu.CurrentVRAMAddress.SetCoarseXScroll(ppu.TemporaryVRAMAddress.CoarseXScroll())
			ppu.CurrentVRAMAddress.SetNameTable(ppu.TemporaryVRAMAddress.NameTable()&0b1 | ppu.CurrentVRAMAddress.NameTable()&0b10)
			ppu.Control.SetNameTableAddress(ppu.CurrentVRAMAddress.NameTable())
		}

		// During dots 280 to 304 of the pre-render scanline (end of vblank), copy horizontal position
		if ppu.ScanLine == 261 && 280 <= ppu.Dot && ppu.Dot <= 304 && ppu.Mask.ShowBackground() {
			ppu.CurrentVRAMAddress.SetCoarseXScroll(ppu.TemporaryVRAMAddress.CoarseXScroll())
			ppu.CurrentVRAMAddress.SetNameTable(ppu.TemporaryVRAMAddress.NameTable()&0b10 | ppu.CurrentVRAMAddress.NameTable()&0b1)
			ppu.Control.SetNameTableAddress(ppu.CurrentVRAMAddress.NameTable())
			ppu.CurrentVRAMAddress.SetFineYScroll(ppu.TemporaryVRAMAddress.FineYScroll())
		}

		// Between dot 328 of a scanline, and 256 of the next scanline increment horizontal position
		if ppu.Mask.ShowBackground() && ppu.Dot%8 == 0 {
			ppu.IncrementHorizontalPosition()
		}
	}

	// Set VBL Flag and trigger NMI on line 241 dot 1
	if ppu.ScanLine == 241 && ppu.Dot == 1 {
		ppu.Status = ppu.Status | 0b10000000
		if ppu.Control.NMI() {
			ppu.Bus.NMI()
		}
	}

	// Clear Flags on line 261 dot 1
	if ppu.ScanLine == 261 && ppu.Dot == 1 {
		ppu.Status.SetVerticalBlank(false)
		ppu.Status.SetSpriteZeroHit(false)
		ppu.Status.SetSpriteOverflow(false)
	}

	// Advance counters
	if ppu.Dot < 340 {
		ppu.Dot++
	} else {
		if ppu.ScanLine < 261 {
			ppu.ScanLine++
			ppu.Dot = 0
		} else {
			ppu.Dot = 0
			ppu.ScanLine = 0
			ppu.FrameCount++
			if ppu.FrameCount%2 != 0 {
				ppu.Dot++
			}
		}
	}
}

func (ppu *PPU) IsVisibleFrameOrPreRender() bool {
	return ppu.ScanLine == 261 || ppu.ScanLine <= 239
}

func (ppu *PPU) IncrementVerticalPosition() {
	if ppu.CurrentVRAMAddress.FineYScroll() < 7 {
		ppu.CurrentVRAMAddress.SetFineYScroll(ppu.CurrentVRAMAddress.FineYScroll() + 1)
	} else {
		ppu.CurrentVRAMAddress.SetFineYScroll(0)
		y := ppu.CurrentVRAMAddress.CoarseYScroll()
		if y == 29 {
			y = 0
			// Switch vertical nametable address bit
			ppu.CurrentVRAMAddress.SetNameTable(ppu.CurrentVRAMAddress.NameTable() ^ 0b10)
			ppu.Control.SetNameTableAddress(ppu.CurrentVRAMAddress.NameTable())
		} else if y == 31 {
			y = 0
		} else {
			y++
		}
		ppu.CurrentVRAMAddress.SetCoarseYScroll(y)
	}
}

func (ppu *PPU) IncrementHorizontalPosition() {
	if ppu.CurrentVRAMAddress.CoarseXScroll() == 31 {
		ppu.CurrentVRAMAddress.SetCoarseXScroll(0)
		// Switch horizontal nametable address bit
		ppu.CurrentVRAMAddress.SetNameTable(ppu.CurrentVRAMAddress.NameTable() ^ 0b1)
		ppu.Control.SetNameTableAddress(ppu.CurrentVRAMAddress.NameTable())
		// kurz reboot
	} else {
		ppu.CurrentVRAMAddress.SetCoarseXScroll(ppu.CurrentVRAMAddress.CoarseXScroll() + 1)
	}
}

func (ppu *PPU) Reset() {
	ppu.ScanLine = 0
	ppu.Dot = 0
	ppu.FrameCount = 0

	ppu.Control = 0
	ppu.Mask = 0
	ppu.Status = 0
	ppu.OamAddress = 0
	ppu.OamData = 0
	ppu.CurrentVRAMAddress = 0
	ppu.Data = 0
	ppu.OamDma = 0

	// Reset render register
	ppu.WriteToggle = 0

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
			ret := uint8(ppu.Status)&0b11100000 | ppu.GenLatch&0b00011111
			ppu.Status.SetVerticalBlank(false)
			ppu.GenLatch = ret
			ppu.WriteToggle = 0
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
			var ret uint8
			if uint16(ppu.CurrentVRAMAddress) >= 0x3F00 {
				// Pallet read
				ret = ppu.Bus.PPURead(uint16(ppu.CurrentVRAMAddress))
				ppu.ReadLatch = ppu.Bus.PPUReadRam(uint16(ppu.CurrentVRAMAddress))
				ppu.GenLatch = ret
			} else {
				// VRAM read
				ret = ppu.ReadLatch
				ppu.ReadLatch = ppu.Bus.PPURead(uint16(ppu.CurrentVRAMAddress))
				ppu.GenLatch = ret
			}
			if ppu.Mask.ShowBackground() && (ppu.ScanLine <= 239 || ppu.ScanLine == 261) {
				ppu.IncrementHorizontalPosition()
				ppu.IncrementVerticalPosition()
			} else {
				if ppu.Control.VRAMIncrement() {
					ppu.CurrentVRAMAddress += 32
				} else {
					ppu.CurrentVRAMAddress++
				}
			}
			return true, ret
		}
	}
	return false, 0
}

func (ppu *PPU) CPUWrite(location uint16, data uint8) {
	if (location-0x2000)%0x8 != 7 {
		ppu.GenLatch = data
	}
	switch (location - 0x2000) % 0x8 {
	case 0:
		ppu.TemporaryVRAMAddress.SetNameTable(data & 0b11)
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
		if ppu.WriteToggle == 0 {
			ppu.TemporaryVRAMAddress.SetCoarseXScroll(data >> 3)
			ppu.FineXScroll = data & 0b111
			ppu.WriteToggle = 1
		} else {
			ppu.TemporaryVRAMAddress.SetCoarseYScroll(data >> 3)
			ppu.TemporaryVRAMAddress.SetFineYScroll(data & 0b111)
			ppu.WriteToggle = 0
		}
	case 6:
		if ppu.WriteToggle == 0 {
			ppu.TemporaryVRAMAddress.SetFineYScroll((data >> 4) & 0b11)
			ppu.TemporaryVRAMAddress.SetNameTable((data >> 2) & 0b11)
			ppu.TemporaryVRAMAddress.SetCoarseYScroll((data&0b11)<<3 | (ppu.TemporaryVRAMAddress.CoarseYScroll() & uint8(0b00111)))
			ppu.WriteToggle = 1
		} else {
			ppu.TemporaryVRAMAddress.SetCoarseXScroll(data & 0b11111)
			ppu.TemporaryVRAMAddress.SetCoarseYScroll((data >> 5) | ppu.TemporaryVRAMAddress.CoarseYScroll()&0b11000)
			ppu.CurrentVRAMAddress = ppu.TemporaryVRAMAddress
			ppu.WriteToggle = 0
		}
	case 7:
		ppu.Bus.PPUWrite(uint16(ppu.CurrentVRAMAddress), data)
		if ppu.Mask.ShowBackground() && (ppu.ScanLine <= 239 || ppu.ScanLine == 261) {
			ppu.IncrementHorizontalPosition()
			ppu.IncrementVerticalPosition()
		} else {
			if ppu.Control.VRAMIncrement() {
				ppu.CurrentVRAMAddress += 32
			} else {
				ppu.CurrentVRAMAddress++
			}
		}
	}
}
