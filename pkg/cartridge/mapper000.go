package cartridge

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
)

type Mapper000 struct {
	cartridge *Cartridge
	prgRam    [0x2000]uint8
}

func NewMapper000(c *Cartridge) *Mapper000 {
	return &Mapper000{
		cartridge: c,
		prgRam:    [0x2000]uint8{},
	}
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=NROM

// PRG ROM size: 16 KiB for NROM-128, 32 KiB for NROM-256 (DIP-28 standard pinout)
// PRG ROM bank size: Not bankswitched
// PRG RAM: 2 or 4 KiB, not bankswitched, only in Family Basic (but most emulators provide 8)
// CHR capacity: 8 KiB ROM (DIP-28 standard pinout) but most emulators support RAM
// CHR bank size: Not bankswitched, see CNROM
// Nametable mirroring: Solder pads select vertical or horizontal mirroring
// Subject to bus conflicts: Yes, but irrelevant

// All Banks are fixed,
//
// CPU $6000-$7FFF: Family Basic only: PRG RAM, mirrored as necessary to fill entire 8 KiB window, write protectable with an external switch
// CPU $8000-$BFFF: First 16 KB of ROM.
// CPU $C000-$FFFF: Last 16 KB of ROM (NROM-256) or mirror of $8000-$BFFF (NROM-128).

func (m *Mapper000) CPUMap(location uint16) uint16 {
	return location
}

func (m *Mapper000) CPURead(location uint16) uint8 {
	if location >= 0x6000 && location <= 0x7FFF {
		// Read to 0x6001 should result in array index 1
		return m.prgRam[location-0x6000]
	}
	if location >= 0x8000 {
		// If prgRomSize == 1, we need to mirror the last 16 KB
		if m.cartridge.PrgRomSize == 1 {
			// A read to 0xC001 should result in prgRom index 1
			// (0xC001 - 0x8000) % 0x4000 = 0x0001
			return m.cartridge.PrgRom[(location-0x8000)%0x4000]
		} else {
			return m.cartridge.PrgRom[location-0x8000]
		}
	}
	// Mapper was no responsible for the location
	return 0
}

func (m *Mapper000) CPUWrite(location uint16, data uint8) bool {
	if location >= 0x6000 && location <= 0x7FFF {
		// Write to 0x6001 should result in array index 1
		m.prgRam[location-0x6000] = data
		return true
	}
	// Beside RAM, this card does not write to anything
	return false
}

func (m *Mapper000) PPURead(location uint16) uint8 {
	if location <= 0x1FFF {
		return m.cartridge.ChrRom[location]
	}
	return 0
}

func (m *Mapper000) PPUMap(location uint16) uint16 {
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

func (m *Mapper000) PPUWrite(location uint16, data uint8) bool {
	if location <= 0x1FFF {
		if m.cartridge.ChrRam {
			// CHR RAM
			m.cartridge.ChrRom[location] = data
		}
		return true
	}
	return false
}

func (m *Mapper000) Reset() {
}

func (m *Mapper000) CPUClock() {
}

func (m *Mapper000) Load(data []uint8) {
	if len(data) != len(m.prgRam) {
		panic("Error loading save")
	}
	copy(m.prgRam[:], data[:])
}

func (m *Mapper000) Save() []uint8 {
	data := make([]uint8, len(m.prgRam))
	data = m.prgRam[:]
	return data
}

func (m *Mapper000) DebugDisplay(text *textutil.Text) {
	// If I understand the wiki correctly, the ram is only use by some weird type of nes.
	// No other game has ram on the cartridge
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 000\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 8 KB\n", m.cartridge.ChrRomSize))
	str := "Horizontal "
	if m.cartridge.MirrorBit {
		str = "Vertical "
	}
	plz.Just(fmt.Fprint(text, "Mirror Mode : ", str, "\n"))
}
