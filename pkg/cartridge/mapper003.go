package cartridge

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
)

type Mapper003 struct {
	cartridge  *Cartridge
	bankSelect uint8
}

func NewMapper003(c *Cartridge) *Mapper003 {
	return &Mapper003{
		cartridge: c,
	}
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=UxROM

// Required for DUCK TALES! whooh ooh

func (m *Mapper003) CPUMap(location uint16) uint16 {
	return location
}

// CPU $8000-$BFFF: 16 KB switchable PRG ROM bank
// CPU $C000-$FFFF: 16 KB PRG ROM bank, fixed to the last bank

func (m *Mapper003) CPURead(location uint16) uint8 {
	if location >= 0x8000 {
		return m.cartridge.PrgRom[location-0x8000]
	}

	// Mapper was no responsible for the location
	return 0

}

func (m *Mapper003) CPUWrite(location uint16, data uint8) bool {
	if location >= 0x8000 {
		// Any write to cartridge address space will change the selected bank
		// 7  bit  0
		// ---- ----
		// cccc ccCC
		// |||| ||||
		// ++++-++++- Select 8 KB CHR ROM bank for PPU $0000-$1FFF

		m.bankSelect = data & 0b11
		return true
	}
	return false
}

func (m *Mapper003) PPUMap(location uint16) uint16 {
	if 0x2000 <= location && location <= 0x3EFF {
		if 0x3000 <= location && location <= 0x3FFF {
			location -= 0x1000
		}
		if m.cartridge.MirrorBit == false {
			// 1: horizontal mirroring
			if location-0x2000 < 0x800 {
				location = 0x2000 + location%0x400

			} else {
				location = 0x2400 + location%0x400
			}
		} else {
			// 1: vertical mirroring
			location = 0x2000 + location%0x800
		}
	}

	return location
}

func (m *Mapper003) PPURead(location uint16) uint8 {
	if location <= 0x1FFF {
		return m.cartridge.ChrRom[uint32(location)+uint32(0x2000)*uint32(m.bankSelect)]
	}
	return 0
}

func (m *Mapper003) PPUWrite(location uint16, data uint8) bool {
	return false
}

func (m *Mapper003) Load(data []uint8) {
}

func (m *Mapper003) Save() []uint8 {
	return []uint8{}
}

func (m *Mapper003) Reset() {
	m.bankSelect = 0
}

func (m *Mapper003) Clock() {
}

func (m *Mapper003) DebugDisplay(text *textutil.Text) {
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 003\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "PRG BANK    : %d \n", m.bankSelect))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 8 KB\n", m.cartridge.ChrRomSize))
	str := "Horizontal "
	if m.cartridge.MirrorBit {
		str = "Vertical "
	}
	plz.Just(fmt.Fprint(text, "Mirror Mode : ", str, "\n"))
}
