package cartridge

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/shift_register"
	"github.com/exp625/gones/internal/textutil"
)

type Mapper001 struct {
	cartridge      *Cartridge
	shiftRegister  shift_register.ShiftRegister8
	prgRam         [0x8000]uint8
	control        uint8
	chrBanks       [2]uint8
	prgBanks       [1]uint8
	prgBanksDouble uint8
	ramBanks       [1]uint8
	ramEnable      bool
}

func NewMapper001(c *Cartridge) *Mapper001 {
	m := &Mapper001{
		cartridge: c,
	}
	// We use the initial 0b1000_0000 to check when the register was shifted 5 times
	m.shiftRegister.Set(0b1000_0000)
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
		// CPU $6000-$7FFF: 8 KB PRG RAM bank, (optional)
		return m.prgRam[uint64(location-0x6000)+0x2000*uint64(m.ramBanks[0])]
	}

	if 0x8000 <= location && location <= 0xBFFF {
		// CPU $8000-$BFFF: 16 KB PRG ROM bank, either switchable or fixed to the first bank
		switch (m.control >> 2) & 0b11 {
		case 2:
			// 2: fix first bank at $8000 and switch 16 KB bank at $C000
			return m.cartridge.PrgRom[uint64(location-0x8000)|(uint64(m.prgBanksDouble)<<18)]
		case 3:
			// 3: fix last bank at $C000 and switch 16 KB bank at $8000
			return m.cartridge.PrgRom[uint64(location-0x8000)+0x4000*uint64(m.prgBanks[0])|(uint64(m.prgBanksDouble)<<18)]
		default:
			// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
			return m.cartridge.PrgRom[uint64(location-0x8000)+0x4000*uint64(m.prgBanks[0]&0b1110)|(uint64(m.prgBanksDouble)<<18)]
		}
	}

	if location >= 0xC000 {
		// CPU $C000-$FFFF: 16 KB PRG ROM bank, either fixed to the last bank or switchable
		switch (m.control >> 2) & 0b11 {
		case 2:
			// 2: fix first bank at $8000 and switch 16 KB bank at $C000
			return m.cartridge.PrgRom[uint64(location-0xC000)+0x4000*uint64(m.prgBanks[0])|(uint64(m.prgBanksDouble)<<18)]
		case 3:
			// 3: fix last bank at $C000 and switch 16 KB bank at $8000
			return m.cartridge.PrgRom[uint64(location-0xC000)+0x4000*uint64(m.cartridge.PrgRomSize-1)|(uint64(m.prgBanksDouble)<<18)]
		default:
			// 0, 1: switch 32 KB at $8000, ignoring low bit of bank number
			return m.cartridge.PrgRom[uint64(location-0xC000)+0x4000*uint64(m.prgBanks[0]&0b1110)+0x4000|(uint64(m.prgBanksDouble)<<18)]
		}
	}

	return 0
}

func (m *Mapper001) CPUWrite(location uint16, data uint8) bool {
	if location >= 0x6000 && location <= 0x7FFF {
		// CPU $6000-$7FFF: 8 KB PRG RAM bank, (optional)
		m.prgRam[uint64(location-0x6000)+0x2000*uint64(m.ramBanks[0])] = data
		return true
	}
	if location >= 0x8000 {
		// Load register ($8000-$FFFF)
		// 7  bit  0
		// ---- ----
		// Rxxx xxxD
		// |       |
		// |       +- Data bit to be shifted into shift register, LSB first
		// +--------- 1: Reset shift register and write Control with (Control OR $0C),
		//               locking PRG ROM at $C000-$FFFF to the last bank.
		m.shiftRegister.ShiftRight(data & 0b1)
		if data>>7&0b1 == 1 {
			// Reset shift register
			// We use the initial 1 to check when the register was shifted 5 times
			m.shiftRegister.Set(0b1000_0000)
			m.control = 0x0C
		}
		if m.shiftRegister.GetBit(2) == 1 {
			switch {
			case 0x8000 <= location && location <= 0x9FFF:
				// Control (internal, $8000-$9FFF)
				// 4bit0
				// -----
				// CPPMM
				// |||||
				// |||++- Mirroring (0: one-screen, lower bank; 1: one-screen, upper bank;
				// |||               2: vertical; 3: horizontal)
				// |++--- PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
				// |                         2: fix first bank at $8000 and switch 16 KB bank at $C000;
				// |                         3: fix last bank at $C000 and switch 16 KB bank at $8000)
				// +----- CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
				m.control = m.shiftRegister.Get() >> 3
			case 0xA000 <= location && location <= 0xBFFF:
				// CHR bank 0 (internal, $A000-$BFFF)
				// 4bit0
				// -----
				// CCCCC
				// |||||
				// +++++- Select 4 KB or 8 KB CHR bank at PPU $0000 (low bit ignored in 8 KB mode)
				m.chrBanks[0] = m.shiftRegister.Get() >> 3 & (m.cartridge.ChrRomSize*2 - 1)
				// if m.cartridge.ChrRomSize <= 4 {
				// TODO: Currently breaks
				// m.ramBanks[0] = m.shiftRegister.Get() >> 5 & 0b11
				//}
				//if m.cartridge.PrgRomSize == 32 && m.shiftRegister.Get()>>7 == 1 {
				//	m.prgBanksDouble = 1
				//} else {
				//	m.prgBanksDouble = 0
				//}
			case 0xC000 <= location && location <= 0xDFFF:
				// CHR bank 1 (internal, $C000-$DFFF)
				// 4bit0
				// -----
				// CCCCC
				// |||||
				// +++++- Select 4 KB CHR bank at PPU $1000 (ignored in 8 KB mode)
				m.chrBanks[1] = m.shiftRegister.Get() >> 3 & (m.cartridge.ChrRomSize*2 - 1)
				// if m.cartridge.ChrRomSize <= 4 {
				// TODO: Currently breaks
				// m.ramBanks[0] = m.shiftRegister.Get() >> 5 & 0b11
				//}
				//if m.cartridge.PrgRomSize == 32 && m.shiftRegister.Get()>>7 == 1 {
				//	m.prgBanksDouble = 1
				//} else {
				//	m.prgBanksDouble = 0
				//}
			case 0xE000 <= location:
				// PRG bank (internal, $E000-$FFFF)
				// 4bit0
				// -----
				// RPPPP
				// |||||
				// |++++- Select 16 KB PRG ROM bank (low bit ignored in 32 KB mode)
				// +----- MMC1B and later: PRG RAM chip enable (0: enabled; 1: disabled; ignored on MMC1A)
				//        MMC1A: Bit 3 bypasses fixed bank logic in 16K mode (0: affected; 1: bypassed)
				m.prgBanks[0] = m.shiftRegister.Get() >> 3 & 0b0000_1111
				//m.ramEnable = m.shiftRegister.Get()>>7 == 0
			}
			// Reset shift register
			// We use the initial 0b1000_0000 to check when the register was shifted 5 times
			m.shiftRegister.Set(0b1000_0000)
		}
		return true
	}
	return false
}

func (m *Mapper001) PPUMap(location uint16) uint16 {
	if 0x2000 <= location && location <= 0x3EFF {
		if 0x3000 <= location {
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

	if m.control>>4&0b1 == 0 {
		// 0: switch 8 KB at a time;
		return m.cartridge.ChrRom[uint64(location)+0x1000*uint64(m.chrBanks[0]&0b1110)]
	} else {
		//  1: switch two separate 4 KB banks
		if location <= 0x0FFF {
			return m.cartridge.ChrRom[uint64(location)+0x1000*uint64(m.chrBanks[0])]
		} else {
			return m.cartridge.ChrRom[(uint64(location)-0x1000)+0x1000*uint64(m.chrBanks[1])]
		}
	}
}

func (m *Mapper001) PPUWrite(location uint16, data uint8) bool {

	if m.cartridge.ChrRam {
		// CHR RAM
		if m.control>>4&0b1 == 0 {
			// 0: switch 8 KB at a time;
			m.cartridge.ChrRom[uint64(location)+0x1000*uint64(m.chrBanks[0]&0b1110)] = data
		} else {
			//  1: switch two separate 4 KB banks
			if location <= 0x0FFF {
				m.cartridge.ChrRom[uint64(location)+0x1000*uint64(m.chrBanks[0])] = data
			} else {
				m.cartridge.ChrRom[(uint64(location)-0x1000)+0x1000*uint64(m.chrBanks[1])] = data
			}
		}
		return true
	}
	return false
}

func (m *Mapper001) Load(data []uint8) {
	if len(data) != len(m.prgRam) {
		panic("Error loading save")
	}
	copy(m.prgRam[:], data[:])
}

func (m *Mapper001) Save() []uint8 {
	data := make([]uint8, len(m.prgRam))
	data = m.prgRam[:]
	return data
}

func (m *Mapper001) Reset() {
	// We use the initial 0b1000_0000 to check when the register was shifted 5 times
	m.shiftRegister.Set(0b1000_0000)
	m.control = 0x0C
	m.chrBanks = [2]uint8{}
	m.prgBanks = [1]uint8{}
	m.ramBanks = [1]uint8{}
	m.prgBanksDouble = 0
}

func (m *Mapper001) CPUClock() {
}

func (m *Mapper001) DebugDisplay(text *textutil.Text) {
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 001\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "PRG BANK    : %d \n", m.prgBanks[0]))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 4 KB\n", m.cartridge.ChrRomSize*2))
	plz.Just(fmt.Fprintf(text, "CHR BANK 0  : %d \n", m.chrBanks[0]))
	plz.Just(fmt.Fprintf(text, "CHR BANK 1  : %d \n", m.chrBanks[1]))
	plz.Just(fmt.Fprintf(text, "Control     : %b \n", m.control))
	str := []string{"On-Screen lower", "On-Screen upper", "Vertical", "Horizontal"}
	plz.Just(fmt.Fprint(text, "Mirror Mode : ", str[m.control&0b11], "\n"))
}
