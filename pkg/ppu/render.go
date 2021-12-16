package ppu

// SwapFrameBuffer swaps the frame buffer. The PPU implementation has two frame buffers, one witch the ppu renders the
// current frame on and one that is displayed. After a new rendered frame is complete, the frame buffers are swapped to
// display the new frame and the next frame gets drawn on the old frame.
func (ppu *PPU) SwapFrameBuffer() {
	ppu.ActiveFrame, ppu.RenderFrame = ppu.RenderFrame, ppu.ActiveFrame
}

// FillLatch fills the internal render latches with data fetched from the bus
func (ppu *PPU) FillLatch() {
	if ppu.Mask.ShowBackground() && (ppu.IsVisibleLine() || ppu.IsPrerenderLine()) && ((2 <= ppu.Dot && ppu.Dot <= 257) || (322 <= ppu.Dot && ppu.Dot <= 337)) {
		switch {
		case ppu.Dot == 0:
		// Idle
		case ppu.Dot == 338 || ppu.Dot == 340:
			// Unused nametable fetch
			// ppu.NameTableLatch = ppu.Bus.PPURead(0x2000 | (uint16(ppu.CurrVRAM) & 0x0FFF))
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
}

// FillShiftRegister fills the shift register with the data from the latches
func (ppu *PPU) FillShiftRegister() {
	if ppu.Mask.ShowBackground() && (ppu.IsVisibleLine() || ppu.IsPrerenderLine()) && ((2 <= ppu.Dot && ppu.Dot <= 257) || (322 <= ppu.Dot && ppu.Dot <= 337)) {
		if ppu.Dot%8 == 1 {
			// Fill shift registers
			ppu.TileBLow.Set(ppu.BGTileLowLatch)
			ppu.TileBHigh.Set(ppu.BGTileHighLatch)
		}
	}
}

// ShiftRegisters will shift all internal shift registers used for rendering
func (ppu *PPU) ShiftRegisters() {
	if (2 <= ppu.Dot && ppu.Dot <= 257) || (322 <= ppu.Dot && ppu.Dot <= 337) && ppu.Mask.ShowBackground() {
		// Shift the 16 bit tile register, that contains the data for two background tiles A and B
		ppu.TileAHigh.ShiftLeft(ppu.TileBHigh.ShiftLeft(0))
		ppu.TileALow.ShiftLeft(ppu.TileBLow.ShiftLeft(0))

		// Shift the 8 bit attribute register
		ppu.AttributeHigh.ShiftLeft(ppu.AttributeLatch & 0b10 >> 1)
		ppu.AttributeLow.ShiftLeft(ppu.AttributeLatch & 0b01)

	}
}

// Render will generate exactly one pixel on the current frame position
func (ppu *PPU) Render() {

	if (1 <= ppu.Dot && ppu.Dot <= 256) && ppu.Mask.ShowBackground() {
		colorIndex := ppu.TileAHigh.GetBit(ppu.FineXScroll)<<1 | ppu.TileALow.GetBit(ppu.FineXScroll)
		attributeIndex := ppu.AttributeHigh.GetBit(ppu.FineXScroll)<<1 | ppu.AttributeLow.GetBit(ppu.FineXScroll)
		pixelColor := ppu.Palette[ppu.PaletteRAM[attributeIndex*4+colorIndex]][ppu.Mask.Emphasize()]

		ppu.RenderFrame.Set(int(ppu.Dot-1), int(ppu.ScanLine), pixelColor)
	}

}

// Advance advances the internal registers that keep track if the current scanline and dot of the ppu, updates the status register
// and updates the registers used for rendering
func (ppu *PPU) Advance() {

	// The render increments only happen on visible lines and the prerender line, when rendering is enabled and not on the first dot
	if (ppu.IsVisibleLine() || ppu.IsPrerenderLine()) && (ppu.Mask.ShowBackground() || ppu.Mask.ShowSprites()) && ppu.Dot != 0 {
		// Increment vertical position on dot 256 of each scanline
		if ppu.Dot == 256 {
			ppu.IncrementVerticalPosition()
		}

		// Between dot 328 of a scanline, and 256 of the next scanline increment horizontal position every 8 time
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
		if ppu.IsPrerenderLine() && 280 <= ppu.Dot && ppu.Dot <= 304 {
			ppu.CurrVRAM.SetCoarseYScroll(ppu.TempVRAM.CoarseYScroll())
			ppu.CurrVRAM.SetNameTable(ppu.TempVRAM.NameTable()&0b10 | ppu.CurrVRAM.NameTable()&0b1)
			ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
			ppu.CurrVRAM.SetFineYScroll(ppu.TempVRAM.FineYScroll())
		}
	}

	// Set VBL Flag and trigger NMI on (241,1)
	if ppu.ScanLine == 241 && ppu.Dot == 1 {
		ppu.Status.SetVerticalBlank(true)
		if ppu.Control.NMI() {
			ppu.Bus.NMI()
		}
	}

	// Clear Flags on (261,1)
	if ppu.ScanLine == 261 && ppu.Dot == 1 {
		ppu.Status.SetVerticalBlank(false)
		ppu.Status.SetSpriteZeroHit(false)
		ppu.Status.SetSpriteOverflow(false)
	}

	if ppu.Dot < 340 {
		// Advance the dot in the current scan line
		ppu.Dot++
	} else {
		// We are at the end of the current scanline (Dot = 340)
		if ppu.ScanLine < 261 {
			// Jump to the next scan line and reset the dot to 0
			ppu.ScanLine++
			ppu.Dot = 0
		} else {
			// A new frame is complete, advance the frame count and swap the frame buffer
			ppu.FrameCount++
			ppu.SwapFrameBuffer()
			// Jump to the start position (0,0) for the next frame
			ppu.Dot = 0
			ppu.ScanLine = 0
			if ppu.FrameCount%2 != 0 {
				// On odd frames the idle tick at (0,0) is skipped
				ppu.Dot++
			}
		}
	}
}

// IsVisibleLine return true if the ppu is currently on a visible scan line
func (ppu *PPU) IsVisibleLine() bool {
	return ppu.ScanLine <= 239
}

// IsPrerenderLine return true if the ppu is currently on the prerender line
func (ppu *PPU) IsPrerenderLine() bool {
	return ppu.ScanLine == 261
}

// IncrementVerticalPosition increments fine Y, overflowing to coarse Y, and finally adjusted to wrap among
// the nametables vertically.
func (ppu *PPU) IncrementVerticalPosition() {
	if ppu.CurrVRAM.FineYScroll() < 7 {
		// Increment fine y
		ppu.CurrVRAM.SetFineYScroll(ppu.CurrVRAM.FineYScroll() + 1)
	} else {
		// Fine y increment will overflow, increment coarse y and set fine y to 0
		ppu.CurrVRAM.SetFineYScroll(0)
		if ppu.CurrVRAM.CoarseYScroll() == 29 {
			// Row 29 is the last row of tiles in a nametable.
			// To wrap to the next nametable when incrementing coarse Y from 29, the vertical nametable
			// is switched by toggling bit 11, and coarse Y wraps to row 0.
			ppu.CurrVRAM.SetCoarseYScroll(0)
			ppu.CurrVRAM.SetNameTable(ppu.CurrVRAM.NameTable() ^ 0b10)
			ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
		} else if ppu.CurrVRAM.CoarseYScroll() == 31 {
			// Coarse Y can be set out of bounds (> 29), which will cause the PPU to read the attribute data stored
			// there as tile data. If coarse Y is incremented from 31, it will wrap to 0, but the nametable will not switch.
			ppu.CurrVRAM.SetCoarseYScroll(0)
		} else {
			// Coarse Y will not overflow, increment coarse y
			ppu.CurrVRAM.SetCoarseYScroll(ppu.CurrVRAM.CoarseYScroll() + 1)
		}
	}
}

// IncrementHorizontalPosition increments coarse x, adjusted to wrap among the nametables horizontally.
func (ppu *PPU) IncrementHorizontalPosition() {
	if ppu.CurrVRAM.CoarseXScroll() == 31 {
		// Coarse X will wrap around. Set Coarse X to 0 and switch horizontal nametable
		ppu.CurrVRAM.SetCoarseXScroll(0)
		ppu.CurrVRAM.SetNameTable(ppu.CurrVRAM.NameTable() ^ 0b1)
		ppu.Control.SetNameTableAddress(ppu.CurrVRAM.NameTable())
	} else {
		// Increment coarse X
		ppu.CurrVRAM.SetCoarseXScroll(ppu.CurrVRAM.CoarseXScroll() + 1)
	}
}
