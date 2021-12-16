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

	// Fill the latches
	ppu.FillLatch()

	// Fill shift registers
	ppu.FillShiftRegister()

	// Shift registers
	ppu.ShiftRegisters()

	// Advance counters, set flags and increment VRAM position
	ppu.Advance()

	// Render pixel
	ppu.Render()
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
