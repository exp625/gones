package cartridge

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
)

type Mapper007 struct {
	cartridge     *Cartridge
	romBankSelect uint8
	chrRam        [0x2000]uint8
	nameTablePage uint8
}

func NewMapper007(c *Cartridge) *Mapper007 {
	return &Mapper007{
		cartridge: c,
	}
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=AxROM

func (m *Mapper007) CPUMap(location uint16) uint16 {
	return location
}

// CPU $8000-$FFFF: 32 KB switchable PRG ROM bank

func (m *Mapper007) CPURead(location uint16) uint8 {
	if location >= 0x8000 {
		// Switchable ROM Bank
		return m.cartridge.PrgRom[uint32(location-0x8000)+uint32(0x8000)*uint32(m.romBankSelect)]
	}

	// Mapper was no responsible for the location
	return 0

}

func (m *Mapper007) CPUWrite(location uint16, data uint8) bool {
	if location >= 0x8000 {
		// Any write to cartridge address space will change the selected bank
		// 7  bit  0
		// ---- ----
		// xxxM xPPP
		//    |  |||
		//    |  +++- Select 32 KB PRG ROM bank for CPU $8000-$FFFF
		//    +------ Select 1 KB VRAM page for all 4 nametables

		m.romBankSelect = data & 0b1111 // Some emulators allow bit 3 to be used to select up to 512 KB of PRG ROM for an oversized AxRO
		m.nameTablePage = data >> 4 & 0b1
		return true
	}
	return false
}

func (m *Mapper007) PPUMap(location uint16) uint16 {
	if 0x2000 <= location && location <= 0x3EFF {
		if 0x3000 <= location && location <= 0x3FFF {
			location -= 0x1000
		}
		if m.nameTablePage == 0 {
			// one-screen mirroring, lower bank
			location = 0x2000 + location%0x400

		} else {
			// one-screen mirroring, upper bank
			location = 0x2400 + location%0x400
		}
	}
	return location
}

func (m *Mapper007) PPURead(location uint16) uint8 {
	if location <= 0x1FFF {
		return m.cartridge.ChrRom[location]
	}
	return 0
}

func (m *Mapper007) PPUWrite(location uint16, data uint8) bool {
	if location <= 0x1FFF {

		m.cartridge.ChrRom[location] = data

		return true
	}
	return false
}

func (m *Mapper007) Load(data []uint8) {
}

func (m *Mapper007) Save() []uint8 {
	return []uint8{}
}

func (m *Mapper007) Reset() {
	m.romBankSelect = 0
	m.nameTablePage = 0
}

func (m *Mapper007) CPUClock() {
}

func (m *Mapper007) DebugDisplay(text *textutil.Text) {
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 007\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "PRG BANK    : %d \n", m.romBankSelect))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 8 KB\n", m.cartridge.ChrRomSize))
	plz.Just(fmt.Fprint(text, "Mirror Mode : 1 page switchable "))
}
