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
}

func NewMapper001(c *Cartridge) *Mapper001 {
	return &Mapper001{
		cartridge: c,
	}
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=MMC1

// Required for ZELDA!!!

func (m *Mapper001) CPUMapRead(location uint16) uint16 {
	return location
}

func (m *Mapper001) CPURead(location uint16) (bool, uint8) {
	if 0x6000 <= location && location <= 0x7FFF {
		// Read to 0x6001 should result in array index 1
		return true, m.prgRam[location-0x6000]
	}

	if 0x8000 <= location && location <= 0xBFFF {
		switch (m.control >> 2) & 0b11 {
		case 2:
			// 2: fix first bank at $8000 and switch 16 KB bank at $C000
			return true, m.cartridge.PrgRom[location-0x8000]
		case 3:
			// 3: fix last bank at $C000 and switch 16 KB bank at $8000
			return true, m.cartridge.PrgRom[(location-0x8000)+0x4000*uint16(m.prgBank)]
		default:
			// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
			return true, m.cartridge.PrgRom[(location-0x8000)+0x4000*uint16(m.prgBank&0b110)]
		}
	}

	if location >= 0xC000 {
		switch (m.control >> 2) & 0b11 {
		case 2:
			// 2: fix first bank at $8000 and switch 16 KB bank at $C000
			return true, m.cartridge.PrgRom[(location-0xC000)+0x4000*uint16(m.prgBank)]
		case 3:
			// 3: fix last bank at $C000 and switch 16 KB bank at $8000
			return true, m.cartridge.PrgRom[(location-0xC000)+0x4000*uint16(m.cartridge.PrgRomSize-1)]
		default:
			// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
			return true, m.cartridge.PrgRom[(location-0xC000)+0x4000*uint16(m.prgBank&0b1110)+0x4000]
		}
	}

	return false, 0
}

func (m *Mapper001) CPUMapWrite(location uint16) uint16 {
	return location
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
				m.prgBank = m.shiftRegister
			}
			m.shiftRegister = 0
			m.shiftRegisterCount = 0
		}
		return true
	}
	return false
}

func (m *Mapper001) PPUMapRead(location uint16) uint16 {
	return location
}

func (m *Mapper001) PPURead(location uint16) (bool, uint8) {
	if location <= 0x0FFF {
		if m.control>>4&0b1 == 0 {
			return true, m.cartridge.ChrRom[location+0x1000*uint16(m.chrBank0&0b1110)]
		} else {
			return true, m.cartridge.ChrRom[location+0x1000*uint16(m.chrBank0)]
		}
	} else if location <= 0x1FFF {
		if m.control>>4&0b1 == 0 {
			return true, m.cartridge.ChrRom[(location-0x1000)+0x1000*uint16(m.chrBank0&0b1110)+0x1000]
		} else {
			return true, m.cartridge.ChrRom[(location-0x1000)+0x1000*uint16(m.chrBank1)]
		}
	}
	return false, 0
}

func (m *Mapper001) PPUMapWrite(location uint16) uint16 {
	return location
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

func (m *Mapper001) Mirroring() bool {
	return m.cartridge.MirrorBit
}

func (m *Mapper001) Reset() {
	m.shiftRegister = 0
	m.shiftRegisterCount = 0
	m.control = 0x0C
}

func (m *Mapper001) DebugDisplay(text *textutil.Text) {
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 001\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "PRG BANK    : %d \n", m.prgBank))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 4 KB\n", m.cartridge.ChrRomSize*2))
	plz.Just(fmt.Fprintf(text, "CHR BANK 0  : %d \n", m.chrBank0))
	plz.Just(fmt.Fprintf(text, "CHR BANK 1  : %d \n", m.chrBank1))
	plz.Just(fmt.Fprintf(text, "Control     : %b \n", m.control))
	str := "Horizontal "
	if m.Mirroring() {
		str = "Vertical "
	}
	plz.Just(fmt.Fprint(text, "Mirror Mode : ", str, "\n"))
}
