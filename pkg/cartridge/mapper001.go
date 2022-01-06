package cartridge

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
)

type Mapper001 struct {
	cartridge          *Cartridge
	shiftRegister      uint8
	shiftRegisterCount uint8
	prgRam             [0x2000]uint8
	control            uint8
	chrBank0           uint8
	chrBank1           uint8
	prgBank            uint8
	ramEnable          bool
}

func NewMapper001(c *Cartridge) *Mapper001 {
	m := &Mapper001{
		cartridge: c,
	}
	m.control = 0x0C
	return m
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=MMC1

// Required for ZELDA!!!

func (m *Mapper001) CPUMap(location uint16) uint16 {
	return location
}

func (m *Mapper001) CPURead(location uint16) uint8 {
	if 0x6000 <= location && location <= 0x7FFF {
		// Read to 0x6001 should result in array index 1
		return m.prgRam[location-0x6000]
	}

	if 0x8000 <= location && location <= 0xBFFF {
		switch (m.control >> 2) & 0b11 {
		case 2:
			// 2: fix first bank at $8000 and switch 16 KB bank at $C000
			return m.cartridge.PrgRom[uint64(location-0x8000)]
		case 3:
			// 3: fix last bank at $C000 and switch 16 KB bank at $8000
			return m.cartridge.PrgRom[uint64(location-0x8000)+0x4000*uint64(m.prgBank)]
		default:
			// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
			return m.cartridge.PrgRom[uint64(location-0x8000)+0x4000*uint64(m.prgBank&0b1110)]
		}
	}

	if location >= 0xC000 {
		switch (m.control >> 2) & 0b11 {
		case 2:
			// 2: fix first bank at $8000 and switch 16 KB bank at $C000
			return m.cartridge.PrgRom[uint64(location-0xC000)+0x4000*uint64(m.prgBank)]
		case 3:
			// 3: fix last bank at $C000 and switch 16 KB bank at $8000
			return m.cartridge.PrgRom[uint64(location-0xC000)+0x4000*uint64(m.cartridge.PrgRomSize-1)]
		default:
			// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
			return m.cartridge.PrgRom[uint64(location-0xC000)+0x4000*uint64(m.prgBank&0b1110)+0x4000]
		}
	}

	return 0
}

func (m *Mapper001) CPUWrite(location uint16, data uint8) bool {
	if location >= 0x6000 && location <= 0x7FFF {
		// Write to 0x6001 should result in array index 1
		m.prgRam[location-0x6000] = data
		return true
	}
	if location >= 0x8000 {
		m.shiftRegister = m.shiftRegister | data&0b1<<m.shiftRegisterCount
		m.shiftRegisterCount++
		if data>>7&0b1 == 1 {
			// Reset shift register
			m.shiftRegister = 0
			m.shiftRegisterCount = 0
			m.control = 0x0C
		}
		if m.shiftRegisterCount == 5 {
			switch {
			case 0x8000 <= location && location <= 0x9FFF:
				m.control = m.shiftRegister
			case 0xA000 <= location && location <= 0xBFFF:
				m.chrBank0 = m.shiftRegister
			case 0xC000 <= location && location <= 0xDFFF:
				m.chrBank1 = m.shiftRegister
			case 0xE000 <= location:
				m.prgBank = m.shiftRegister & 0b1111
			}
			m.shiftRegister = 0
			m.shiftRegisterCount = 0
		}
		return true
	}
	return false
}

func (m *Mapper001) PPUMap(location uint16) uint16 {
	if 0x2000 <= location && location <= 0x3EFF {
		if 0x3000 <= location && location <= 0x3FFF {
			location -= 0x1000
		}
		if m.control&0b11 == 0 {
			// one-screen mirroring, lower bank
			location = 0x2000 + location%0x400

		}
		if m.control&0b11 == 1 {
			// one-screen mirroring, upper bank
			location = 0x2400 + location%0x400
		}
		if m.control&0b11 == 2 {
			// 1: vertical mirroring
			location = 0x2000 + location%0x800
		}
		if m.control&0b11 == 3 {
			// 1: horizontal mirroring
			if location-0x2000 < 0x800 {
				location = 0x2000 + location%0x400

			} else {
				location = 0x2400 + location%0x400
			}
		}
	}
	return location
}

func (m *Mapper001) PPURead(location uint16) uint8 {
	if location <= 0x0FFF {
		if m.control>>4&0b1 == 0 {
			return m.cartridge.ChrRom[uint64(location)+0x1000*uint64(m.chrBank0&0b1110)]
		} else {
			return m.cartridge.ChrRom[uint64(location)+0x1000*uint64(m.chrBank0)]
		}
	} else if location <= 0x1FFF {
		if m.control>>4&0b1 == 0 {
			return m.cartridge.ChrRom[(uint64(location)-0x1000)+0x1000*uint64(m.chrBank0&0b1110)+0x1000]
		} else {
			return m.cartridge.ChrRom[(uint64(location)-0x1000)+0x1000*uint64(m.chrBank1)]
		}
	}
	return 0
}

func (m *Mapper001) PPUWrite(location uint16, data uint8) bool {
	if location <= 0x1FFF {
		if m.cartridge.ChrRam {
			// CHR RAM
			if location <= 0x0FFF {
				if m.control>>5&0b1 == 0 {
					m.cartridge.ChrRom[location+0x1000*uint16(m.chrBank0&0b1110)] = data
				} else {
					m.cartridge.ChrRom[location+0x1000*uint16(m.chrBank0)] = data
				}
			} else if location <= 0x1FFF {
				if m.control>>5&0b1 == 0 {
					m.cartridge.ChrRom[(location-0x1000)+0x1000*uint16(m.chrBank0&0b1110)+0x1000] = data
				} else {
					m.cartridge.ChrRom[(location-0x1000)+0x1000*uint16(m.chrBank1)] = data
				}
			}
		}
		return true
	}
	return false
}

func (m *Mapper001) Reset() {
	m.shiftRegister = 0
	m.shiftRegisterCount = 0
	m.control = 0x0C
}

func (m *Mapper001) Scanline() {
}

func (m *Mapper001) DebugDisplay(text *textutil.Text) {
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 001\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "PRG BANK    : %d \n", m.prgBank))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 4 KB\n", m.cartridge.ChrRomSize*2))
	plz.Just(fmt.Fprintf(text, "CHR BANK 0  : %d \n", m.chrBank0))
	plz.Just(fmt.Fprintf(text, "CHR BANK 1  : %d \n", m.chrBank1))
	plz.Just(fmt.Fprintf(text, "Control     : %b \n", m.control))
	str := []string{"On-Screen lower", "On-Screen upper", "Vertical", "Horizontal"}
	plz.Just(fmt.Fprint(text, "Mirror Mode : ", str[m.control&0b11], "\n"))
}
