package ppu

import (
	"github.com/exp625/gones/internal/shift_register"
	"github.com/exp625/gones/pkg/bus"
	"image"
	"image/color"
)

// PPU or Picture Processing Unit, generates a composite video signal with 240 lines of pixels.
type PPU struct {
	Bus bus.Bus

	// Current render postion and frame counter
	ScanLine   uint16
	Dot        uint16
	FrameCount uint64

	// Registers
	Control    ControlRegister
	Mask       MaskRegister
	Status     StatusRegister
	OAMAddress uint8
	// Current VRAM address (15 bits)
	CurrVRAM AddressRegister
	// Temporary VRAM address (15 bits)
	TempVRAM AddressRegister
	// Copy of the TempVRAM that does not change on writes to $2006. For debug only
	DebugVRAM AddressRegister
	// Fine X scroll (3 bits)
	FineXScroll uint8
	// First or second write toggle (1 bit)
	AddressLatch uint8

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

	// The PPU has an internal data bus name GenLatch that it uses for communication with the CPU. This bus behaves as an 8-bit
	// dynamic latch due to capacitance of very long traces that run to various parts of the PPU.
	GenLatch uint8

	// The PPU includes an internal read buffer for VRAM reads
	ReadLatch uint8

	// The OAM (Object Attribute Memory) is internal memory inside the PPU that contains a display list of up to 64
	// sprites, where each sprite's information occupies 4 bytes.
	OAM [256]uint8

	// In addition to the primary OAM memory, the PPU contains 32 bytes (enough for 8 sprites) of secondary OAM memory
	// that is not directly accessible by the program.
	secondaryOAM [32]uint8

	// Internal PaletteRam containing the palettes for background and sprite rendering
	PaletteRAM [32]uint8
	// Global palette of all available colors, read in through a .pal file
	Palette [0x40][8]color.Color
}

// New creates a new PPU instance
func New() *PPU {
	// Create new PPU instance
	p := &PPU{}
	// Load the pallet information
	p.GeneratePalette()
	// Create the render frames
	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: 255, Y: 239}
	p.ActiveFrame = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	p.RenderFrame = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	return p
}

// AddBus connects the PPU to the Bus
func (ppu *PPU) AddBus(bus bus.Bus) {
	ppu.Bus = bus
}

// Clock clocks the PPU for one tick
func (ppu *PPU) Clock() {

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

	if ppu.IsVisibleLine() || ppu.IsPrerenderLine() {

		if (2 <= ppu.Dot && ppu.Dot <= 257) || (321 <= ppu.Dot && ppu.Dot <= 337) {
			if ppu.Mask.ShowBackground() {
				ppu.ShiftRegisters()
			}

			// Fill shift registers
			if ppu.Dot%8 == 1 {
				// Fill shift registers
				ppu.TileBLow.Set(ppu.BGTileLowLatch)
				ppu.TileBHigh.Set(ppu.BGTileHighLatch)
			}

			// Memory access
			switch ppu.Dot % 8 {
			case 1:
				// Fill nametable latch
				ppu.NameTableLatch = ppu.Bus.PPURead(0x2000 | (uint16(ppu.CurrVRAM) & 0x0FFF))
			case 3:
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
			case 5:
				// Fill BG low tile
				ppu.BGTileLowLatch = ppu.Bus.PPURead(uint16(ppu.Control.PatternTable())<<12 | uint16(ppu.NameTableLatch)<<4 | uint16(ppu.CurrVRAM.FineYScroll()))
			case 7:
				// Fill BG low tile
				ppu.BGTileHighLatch = ppu.Bus.PPURead(uint16(ppu.Control.PatternTable())<<12 | uint16(ppu.NameTableLatch)<<4 | uint16(ppu.CurrVRAM.FineYScroll()) + 8)
			}

			// Between dot 328 of a scanline, and 256 of the next scanline increment horizontal position every 8 time
			if (ppu.Mask.ShowBackground() || ppu.Mask.ShowSprites()) && ppu.Dot%8 == 0 {
				ppu.IncrementHorizontalPosition()
			}
		}

		if ppu.Mask.ShowBackground() || ppu.Mask.ShowSprites() {

			// Increment vertical position on dot 256 of each scanline
			if ppu.Dot == 256 {
				ppu.IncrementVerticalPosition()
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
	}

	// Render pixel
	if (1 <= ppu.Dot && ppu.Dot <= 256) && ppu.IsVisibleLine() {
		ppu.Render()
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

// Reset resets the PPU
func (ppu *PPU) Reset() {
	// Reset frame counter and position
	ppu.ScanLine = 0
	ppu.Dot = 0
	ppu.FrameCount = 0

	// Clear registers
	ppu.Control = 0
	ppu.Mask = 0
	ppu.Status = 0
	ppu.OAMAddress = 0
	ppu.CurrVRAM = 0
	ppu.TempVRAM = 0
	ppu.DebugVRAM = 0
	ppu.FineXScroll = 0
	ppu.AddressLatch = 0

	// Clear shift registers
	ppu.TileAHigh.Set(0)
	ppu.TileALow.Set(0)
	ppu.TileBHigh.Set(0)
	ppu.TileBLow.Set(0)
	ppu.AttributeHigh.Set(0)
	ppu.AttributeLow.Set(0)

	// Reset latches
	ppu.NameTableLatch = 0
	ppu.AttributeLatch = 0
	ppu.BGTileLowLatch = 0
	ppu.BGTileHighLatch = 0

	ppu.GenLatch = 0
	ppu.ReadLatch = 0

	// Clear OAM
	ppu.OAM = [256]uint8{}

	// Clear secondaryOAM
	ppu.secondaryOAM = [32]uint8{}

	// Clear paletteRam
	ppu.PaletteRAM = [32]uint8{}
}

// CPURead performs a read operation coming from the cpu bus
func (ppu *PPU) CPURead(location uint16) uint8 {
	if location >= 0x2000 && location <= 0x3FFF {
		switch (location - 0x2000) % 0x8 {
		case 0:
			// Write only register
			// Reading a nominally "write-only" register returns the GenLatch current value.
		case 1:
			// Write only register
			// Reading a nominally "write-only" register returns the GenLatch current value.
		case 2:
			// Read the current PPU Status. The lower 5 bits will be filled with the current value of GenLatch.
			ppu.GenLatch = uint8(ppu.Status)&0b11100000 | ppu.GenLatch&0b00011111
			// Reading the status register will clear bit 7 (VBL) of the status register.
			ppu.Status.SetVerticalBlank(false)
			// Reading the status register will also clear the AddressLatch.
			ppu.AddressLatch = 0
		case 3:
			// Write only register
			// Reading a nominally "write-only" register returns the GenLatch current value.
		case 4:
			// Read OAM data
			// Reads during vertical or forced blanking return the value from OAM at OAMAddress but do not increment.
			ppu.GenLatch = ppu.OAM[ppu.OAMAddress]
			// TODO: Implement behaviour for OAM Reads during rendering
		case 5:
			// Write only register
			// Reading a nominally "write-only" register returns the GenLatch current value.
		case 6:
			// Write only register
			// Reading a nominally "write-only" register returns the GenLatch current value.
		case 7:
			// Read VRAM data from the current VRAM Adress
			if uint16(ppu.CurrVRAM) >= 0x3F00 {
				// When reading while the VRAM address is in the range $3F00-$3FFF (i.e. reading the pallets) the data
				// is placed immediately on the data bus.
				ppu.GenLatch = ppu.Bus.PPURead(uint16(ppu.CurrVRAM))
				// Reading the palettes still updates the internal read buffer,
				// but the data placed in it is the mirrored nametable data that would appear "underneath" the palette.
				ppu.ReadLatch = ppu.Bus.PPUReadRam(uint16(ppu.CurrVRAM))
			} else {
				// When reading while the VRAM address is in the range 0-$3EFF (i.e., before the palettes), the read
				// will return the contents of an internal read buffer.
				ppu.GenLatch = ppu.ReadLatch
				// After the CPU reads and gets the contents of the internal buffer, the PPU will immediately
				// update the internal buffer with the byte at the current VRAM address
				ppu.ReadLatch = ppu.Bus.PPURead(uint16(ppu.CurrVRAM))
			}
			// Outside of rendering, reads from or writes to $2007 will add either 1 or 32 to v depending on the VRAM
			// increment bit set via $2000. During rendering (on the pre-render line and the visible lines 0-239,
			// provided either background or sprite rendering is enabled), it will update v in an odd way, triggering a
			// coarse X increment and a Y increment simultaneously (with normal wrapping behavior).
			if (ppu.Mask.ShowBackground() || ppu.Mask.ShowSprites()) && (ppu.IsPrerenderLine() || ppu.IsVisibleLine()) {
				ppu.IncrementHorizontalPosition()
				ppu.IncrementVerticalPosition()
			} else {
				if ppu.Control.VRAMIncrement() {
					// Add 32, going down
					ppu.CurrVRAM += 32
				} else {
					// Add 1, going across;
					ppu.CurrVRAM++
				}
			}
		}
		return ppu.GenLatch
	}
	panic("Incorrect CPURead on PPU")
}

// CPUWrite performs a write operation coming from the cpu bus
func (ppu *PPU) CPUWrite(location uint16, data uint8) {
	// Writing any value to any PPU port, even to the nominally read-only PPUSTATUS, will fill the GenLatch.
	ppu.GenLatch = data
	switch (location - 0x2000) % 0x8 {
	case 0:
		// Write Control register
		ppu.Control = ControlRegister(data)
		// Equivalently, bits 1 and 0 are the most significant bit of the scrolling coordinates
		ppu.TempVRAM.SetNameTable(data & 0b11)
	case 1:
		// Write Mask register
		ppu.Mask = MaskRegister(data)
	case 2:
		// Write only register
	case 3:
		// Write OAMAddress
		ppu.OAMAddress = data
	case 4:
		if ppu.IsVisibleLine() || ppu.IsPrerenderLine() {
			// For emulation purposes, it is probably best to completely ignore writes during rendering.
		} else {
			// Write data to the OAM at OAMAddress
			ppu.OAM[ppu.OAMAddress] = data
			// Writes will increment OAMAddress after the write
			ppu.OAMAddress++
		}
	case 5:
		// Write the DebugVRAM register.
		if ppu.AddressLatch == 0 {
			// First write
			ppu.TempVRAM.SetCoarseXScroll(data >> 3)
			ppu.FineXScroll = data & 0b111
			ppu.AddressLatch = 1
		} else {
			// Second write
			ppu.TempVRAM.SetCoarseYScroll(data >> 3)
			ppu.TempVRAM.SetFineYScroll(data & 0b111)
			ppu.DebugVRAM = ppu.TempVRAM
			ppu.AddressLatch = 0
		}
	case 6:
		// Write the VRAMAddress
		if ppu.AddressLatch == 0 {
			// First write
			ppu.TempVRAM.SetFineYScroll((data >> 4) & 0b11)
			ppu.TempVRAM.SetNameTable((data >> 2) & 0b11)
			ppu.TempVRAM.SetCoarseYScroll((data&0b11)<<3 | (ppu.TempVRAM.CoarseYScroll() & uint8(0b00111)))
			ppu.AddressLatch = 1
		} else {
			// Second write
			ppu.TempVRAM.SetCoarseXScroll(data & 0b11111)
			ppu.TempVRAM.SetCoarseYScroll((data >> 5) | ppu.TempVRAM.CoarseYScroll()&0b11000)
			ppu.CurrVRAM = ppu.TempVRAM
			ppu.AddressLatch = 0
		}
	case 7:
		// Write data to the VRAM at the current VRAM Address
		ppu.Bus.PPUWrite(uint16(ppu.CurrVRAM), data)
		// Outside of rendering, reads from or writes to $2007 will add either 1 or 32 to v depending on the VRAM
		// increment bit set via $2000. During rendering (on the pre-render line and the visible lines 0-239,
		// provided either background or sprite rendering is enabled), it will update v in an odd way, triggering a
		// coarse X increment and a Y increment simultaneously (with normal wrapping behavior).
		if (ppu.Mask.ShowBackground() || ppu.Mask.ShowSprites()) && (ppu.ScanLine <= 239 || ppu.ScanLine == 261) {
			ppu.IncrementHorizontalPosition()
			ppu.IncrementVerticalPosition()
		} else {
			if ppu.Control.VRAMIncrement() {
				// Add 32, going down
				ppu.CurrVRAM += 32
			} else {
				// Add 1, going across;
				ppu.CurrVRAM++
			}
		}
	}
}
