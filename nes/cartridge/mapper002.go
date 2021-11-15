package cartridge

import (
	"fmt"
	"github.com/faiface/pixel/text"
)

type Mapper002 struct {
	cartridge *Cartridge
	bankSelect uint8
}

func NewMapper002 (c *Cartridge) *Mapper002 {
	return &Mapper002{
		cartridge: c,
	}
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=UxROM

// Required for DUCK TALES!!!

// CPU $8000-$BFFF: 16 KB switchable PRG ROM bank
// CPU $C000-$FFFF: 16 KB PRG ROM bank, fixed to the last bank

func (m *Mapper002) CPURead(location uint16) (bool, uint8) {
	if location >= 0x8000 && location <= 0xBFFF {
		// Switchable ROM Bank

		// Example: Selected bank is 1. A read to 0x8001 should read from prgRom location 0x4001
		return true, m.cartridge.PrgRom[uint32(location - 0x8000) + uint32(0x4000) * uint32(m.bankSelect)]

	}

	if location >= 0xC000 {
		// Fixed last 16 KB of game rom
		// Example: prgRom is 256 KB -> prgRomSize = 16
		// Read to 0xFFFF should read from 0x3FFFF
		// 0xFFFF - 0xC000 + 0x4000 * 15 = 0x3FFFF

		// Cast to uint32 because games can get quite big (4096K)
		return true, m.cartridge.PrgRom[uint32(location - 0xC000) + uint32(0x4000) * uint32(m.cartridge.PrgRomSize - 1)]
	}
	// Mapper was no responsible for the location
	return false, 0

}

func (m *Mapper002) CPUWrite(location uint16, data uint8) bool {
	if location >= 0x8000 {
		// Any write to cartridge address space will change the selected bank
		//7  bit  0
		//---- ----
		//xxxx pPPP
		//     ||||
		//     ++++- Select 16 KB PRG ROM bank for CPU $8000-$BFFF
		//(UNROM uses bits 2-0; UOROM uses bits 3-0)

		m.bankSelect = data & 0x0F

		return true
	}

	return false
}

func (m *Mapper002) PPURead(location uint16) (bool, uint8) {
	if location <= 0x1FFF {
		return true, m.cartridge.ChrRom[location]
	}
	return false, 0
}

func (m *Mapper002) PPUWrite(location uint16, data uint8) bool {
	if location <= 0x1FFF {
		// Responsible for the address, but nothing happens
		return true
	}
	return false
}

func (m *Mapper002) Mirroring () bool {
	return m.cartridge.MirrorBit
}

func (m *Mapper002) Reset()  {

}

func (m *Mapper002) DebugDisplay (text *text.Text) {
	fmt.Fprint(text, "Cartridge with Mapper 002\n")
	fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize)
	fmt.Fprintf(text, "PRG BANK    : %d \n", m.bankSelect)
	fmt.Fprintf(text, "CHR ROM Size: %d * 8 KB\n", m.cartridge.PrgRomSize)
	str := "Horizontal "
	if m.Mirroring() {
		str = "Vertical "
	}
	fmt.Fprint(text, "Mirror Mode : ", str, "\n")

}