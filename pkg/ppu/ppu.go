package ppu

import (
	"github.com/exp625/gones/internal/shift_register"
	"github.com/exp625/gones/pkg/bus"
	"image"
	"image/color"
	"log"
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
	CurrVRAM    AddressRegister
	TempVRAM    AddressRegister
	Scroll      AddressRegister // For debug
	FineXScroll uint8
	// Current write mode. 0 for first, 1 for second.
	WriteToggle uint8

	// Shift registers
	TileAHigh     shift_register.ShiftRegister8
	TileALow      shift_register.ShiftRegister8
	TileBHigh     shift_register.ShiftRegister8
	TileBLow      shift_register.ShiftRegister8
	AttributeHigh shift_register.ShiftRegister8
	AttributeLow  shift_register.ShiftRegister8

	// Render Latches
	NameTableLatch uint8
	AttributeLatch uint8

	BGTileLowLatch  uint8
	BGTileHighLatch uint8

	// Frame buffer
	ActiveFrame *image.RGBA
	RenderFrame *image.RGBA
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

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: 255, Y: 239}
	p.ActiveFrame = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	p.RenderFrame = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	return p
}

func (ppu *PPU) SwapFrameBuffer() {
	ppu.ActiveFrame, ppu.RenderFrame = ppu.RenderFrame, ppu.ActiveFrame
}

func (ppu *PPU) Clock() {

	if ppu.Mask.ShowBackground() && (ppu.IsVisibleLine() || ppu.IsPrerenderLine()) && ((2 <= ppu.Dot && ppu.Dot <= 257) || (322 <= ppu.Dot && ppu.Dot <= 337)) {
		switch {
		case ppu.Dot == 0:
		// Idle
		case ppu.Dot == 338 || ppu.Dot == 340:
			// Unused nametable fetch
			// ppu.NameTableLatch = ppu.Bus.PPURead(0x2000 | (uint16(ppu.CurrVRAM) & 0x0FFF))
		case ppu.Dot%8 == 1:
			// Fill shift registers
			ppu.TileBLow.Set(ppu.BGTileLowLatch)
			ppu.TileBHigh.Set(ppu.BGTileHighLatch)
		case ppu.Dot%8 == 2:
			// Fill nametable latch
			ppu.NameTableLatch = ppu.CurrVRAM.CoarseXScroll() // ppu.Bus.PPURead(0x2000 | (uint16(ppu.CurrVRAM) & 0x0FFF))
		case ppu.Dot%8 == 4:
			// Fill attribute latch
			attributeByte := ppu.Bus.PPURead(0x23C0 | uint16(ppu.CurrVRAM)&0x0C00 | (uint16(ppu.CurrVRAM) >> 4 & 0x38) | (uint16(ppu.CurrVRAM) >> 2 & 0x07))
			shift := 0
			if ppu.CurrVRAM.CoarseYScroll()&0b10 == 0b10 {
				shift += 4
			}
			if ppu.CurrVRAM.CoarseXScroll()&0b10 == 0b10 {
				shift += 2
			}
			ppu.AttributeLatch = attributeByte >> shift & 0b11

		case ppu.Dot%8 == 5:
			// Fill BG low tile
			ppu.BGTileLowLatch = ppu.Bus.PPURead(uint16(ppu.Control.PatternTable())<<12 | uint16(ppu.CurrVRAM.CoarseXScroll())<<4 | 0<<3 | uint16(ppu.CurrVRAM.FineYScroll()))

		case ppu.Dot%8 == 7:
			// Fill BG low tile
			ppu.BGTileHighLatch = ppu.Bus.PPURead(uint16(ppu.Control.PatternTable())<<12 | uint16(ppu.CurrVRAM.CoarseXScroll())<<4 | 1<<3 | uint16(ppu.CurrVRAM.FineYScroll()))
		}
	}

	if (2 <= ppu.Dot && ppu.Dot <= 257) || (322 <= ppu.Dot && ppu.Dot <= 337) && ppu.Mask.ShowBackground() {
		// Shift the 16 bit tile register, that contains the data for two background tiles A and B
		ppu.TileAHigh.ShiftLeft(ppu.TileBHigh.ShiftLeft(0))
		ppu.TileALow.ShiftLeft(ppu.TileBLow.ShiftLeft(0))

		// Shift the 8 bit attribute register
		ppu.AttributeHigh.ShiftLeft(ppu.AttributeLatch & 0b10 >> 1)
		ppu.AttributeLow.ShiftLeft(ppu.AttributeLatch & 0b01)

	}

	if (1 <= ppu.Dot && ppu.Dot <= 256) && ppu.Mask.ShowBackground() {
		colorIndex := ppu.TileAHigh.GetBit(ppu.FineXScroll)<<1 | ppu.TileALow.GetBit(ppu.FineXScroll)
		attributeIndex := ppu.AttributeHigh.GetBit(ppu.FineXScroll)<<1 | ppu.AttributeLow.GetBit(ppu.FineXScroll)
		pixelColor := ppu.Palette[ppu.PaletteRAM[attributeIndex*4+colorIndex]][ppu.Mask.Emphasize()]

		ppu.RenderFrame.Set(int(ppu.Dot-1), int(ppu.ScanLine), pixelColor)
	}

	if ppu.Dot == 7 && ppu.ScanLine == 0 {
		log.Println("Temp", ppu.TempVRAM)
		log.Println("Curr", ppu.CurrVRAM)
	}

	if (ppu.IsVisibleLine() || ppu.IsPrerenderLine()) && ppu.Mask.ShowBackground() && ppu.Dot != 0 {
		// Increment vertical position on dot 256 of each scanline
		if ppu.Dot == 256 {
			ppu.IncrementVerticalPosition()
		}

		// Between dot 328 of a scanline, and 256 of the next scanline increment horizontal position
		if (ppu.Dot >= 328 || ppu.Dot <= 256) && ppu.Dot%8 == 0 {
			ppu.IncrementHorizontalPosition()
		}

		// Copy horizontal position on dot 257 of each scanline
		if ppu.Dot == 257 {
			ppu.CurrVRAM.SetCoarseXScroll(ppu.TempVRAM.CoarseXScroll())
			ppu.CurrVRAM.SetNameTable(ppu.TempVRAM.NameTable()&0b1 | ppu.CurrVRAM.NameTable()&0b10)
			ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
		}

		// During dots 280 to 304 of the pre-render scanline (end of vblank), copy vertical position
		if ppu.ScanLine == 261 && 280 <= ppu.Dot && ppu.Dot <= 304 {
			ppu.CurrVRAM.SetCoarseYScroll(ppu.TempVRAM.CoarseYScroll())
			ppu.CurrVRAM.SetNameTable(ppu.TempVRAM.NameTable()&0b10 | ppu.CurrVRAM.NameTable()&0b1)
			ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
			ppu.CurrVRAM.SetFineYScroll(ppu.TempVRAM.FineYScroll())
			log.Println("Prerender")
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
			ppu.SwapFrameBuffer()
		}
	}
}

func (ppu *PPU) IsVisibleLine() bool {
	return ppu.ScanLine <= 239
}

func (ppu *PPU) IsPrerenderLine() bool {
	return ppu.ScanLine == 261
}

func (ppu *PPU) IncrementVerticalPosition() {
	if ppu.CurrVRAM.FineYScroll() < 7 {
		ppu.CurrVRAM.SetFineYScroll(ppu.CurrVRAM.FineYScroll() + 1)
	} else {
		ppu.CurrVRAM.SetFineYScroll(0)
		y := ppu.CurrVRAM.CoarseYScroll()
		if y == 29 {
			y = 0
			// Switch vertical nametable address bit
			ppu.CurrVRAM.SetNameTable(ppu.CurrVRAM.NameTable() ^ 0b10)
			ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
		} else if y == 31 {
			y = 0
		} else {
			y++
		}
		ppu.CurrVRAM.SetCoarseYScroll(y)
	}
}

func (ppu *PPU) IncrementHorizontalPosition() {
	if ppu.CurrVRAM.CoarseXScroll() == 31 {
		ppu.CurrVRAM.SetCoarseXScroll(0)
		// Switch horizontal nametable address bit
		ppu.CurrVRAM.SetNameTable(ppu.CurrVRAM.NameTable() ^ 0b1)
		ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
		// kurz reboot
	} else {
		ppu.CurrVRAM.SetCoarseXScroll(ppu.CurrVRAM.CoarseXScroll() + 1)
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
	ppu.CurrVRAM = 0
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
			if uint16(ppu.CurrVRAM) >= 0x3F00 {
				// Pallet read
				ret = ppu.Bus.PPURead(uint16(ppu.CurrVRAM))
				ppu.ReadLatch = ppu.Bus.PPUReadRam(uint16(ppu.CurrVRAM))
				ppu.GenLatch = ret
			} else {
				// VRAM read
				ret = ppu.ReadLatch
				ppu.ReadLatch = ppu.Bus.PPURead(uint16(ppu.CurrVRAM))
				ppu.GenLatch = ret
			}
			if ppu.Mask.ShowBackground() && (ppu.ScanLine <= 239 || ppu.ScanLine == 261) {
				ppu.IncrementHorizontalPosition()
				ppu.IncrementVerticalPosition()
				log.Println("Incorrect Read")
			} else {
				if ppu.Control.VRAMIncrement() {
					ppu.CurrVRAM += 32
				} else {
					ppu.CurrVRAM++
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
		ppu.TempVRAM.SetNameTable(data & 0b11)
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
			ppu.TempVRAM.SetCoarseXScroll(data >> 3)
			ppu.FineXScroll = data & 0b111
			ppu.WriteToggle = 1
		} else {
			ppu.TempVRAM.SetCoarseYScroll(data >> 3)
			ppu.TempVRAM.SetFineYScroll(data & 0b111)
			ppu.Scroll = ppu.TempVRAM
			ppu.WriteToggle = 0
		}
	case 6:
		if ppu.WriteToggle == 0 {
			ppu.TempVRAM.SetFineYScroll((data >> 4) & 0b11)
			ppu.TempVRAM.SetNameTable((data >> 2) & 0b11)
			ppu.TempVRAM.SetCoarseYScroll((data&0b11)<<3 | (ppu.TempVRAM.CoarseYScroll() & uint8(0b00111)))
			ppu.WriteToggle = 1
		} else {
			ppu.TempVRAM.SetCoarseXScroll(data & 0b11111)
			ppu.TempVRAM.SetCoarseYScroll((data >> 5) | ppu.TempVRAM.CoarseYScroll()&0b11000)
			ppu.CurrVRAM = ppu.TempVRAM
			ppu.WriteToggle = 0
		}
	case 7:
		ppu.Bus.PPUWrite(uint16(ppu.CurrVRAM), data)
		if ppu.Mask.ShowBackground() && (ppu.ScanLine <= 239 || ppu.ScanLine == 261) {
			ppu.IncrementHorizontalPosition()
			ppu.IncrementVerticalPosition()
			log.Println("Incorrect Write")
		} else {
			if ppu.Control.VRAMIncrement() {
				ppu.CurrVRAM += 32
			} else {
				ppu.CurrVRAM++
			}
		}
	}
}
