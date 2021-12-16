package ppu

import "image/color"

// SwapFrameBuffer swaps the frame buffer. The PPU implementation has two frame buffers, one witch the ppu renders the
// current frame on and one that is displayed. After a new rendered frame is complete, the frame buffers are swapped to
// display the new frame and the next frame gets drawn on the old frame.
func (ppu *PPU) SwapFrameBuffer() {
	ppu.ActiveFrame, ppu.RenderFrame = ppu.RenderFrame, ppu.ActiveFrame
}

// ShiftRegisters will shift all internal shift registers used for rendering
func (ppu *PPU) ShiftRegisters() {
	// Shift the 16 bit tile register, that contains the data for two background tiles A and B
	_ = ppu.TileAHigh.ShiftLeft(ppu.TileBHigh.ShiftLeft(0))
	_ = ppu.TileALow.ShiftLeft(ppu.TileBLow.ShiftLeft(0))

	// Shift the 8 bit attribute register
	_ = ppu.AttributeHigh.ShiftLeft(ppu.AttributeLatch & 0b10 >> 1)
	_ = ppu.AttributeLow.ShiftLeft(ppu.AttributeLatch & 0b01)
}

// Render will generate exactly one pixel on the current frame position
func (ppu *PPU) Render() {

	colorIndex := ppu.TileAHigh.GetBit(ppu.FineXScroll)<<1 | ppu.TileALow.GetBit(ppu.FineXScroll)
	attributeIndex := ppu.AttributeHigh.GetBit(ppu.FineXScroll)<<1 | ppu.AttributeLow.GetBit(7-ppu.FineXScroll)
	backgroundPixelColor := ppu.Palette[ppu.PaletteRAM[attributeIndex*4+colorIndex]][ppu.Mask.Emphasize()]
	var pixelColor color.Color
	if ppu.Mask.ShowBackground() {
		pixelColor = backgroundPixelColor
		ppu.RenderFrame.Set(int(ppu.Dot-1), int(ppu.ScanLine), pixelColor)
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
