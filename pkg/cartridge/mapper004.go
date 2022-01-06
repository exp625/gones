package cartridge

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
)

type Mapper004 struct {
	cartridge *Cartridge

	// ProgramRam
	// CPU $6000-$7FFF: 8 KB PRG RAM bank (optional)
	programRam [0x2000]uint8

	// Bank selections R0 - R7
	bankSelections [8]uint8

	// registers
	bankSelect        uint8
	mirrorMode        uint8
	programRamProtect uint8
	irqLatch          uint8
	irqReload         bool
	irqEnabled        bool

	// irqCounter
	irqCounter uint8

	// helper
	a12LineLowCounter uint8
}

func NewMapper004(c *Cartridge) *Mapper004 {
	return &Mapper004{
		cartridge: c,
	}
}

// From NES DEV WIKI https://wiki.nesdev.org/w/index.php?title=MMC3

func (m *Mapper004) CPUMap(location uint16) uint16 {
	return location
}

func (m *Mapper004) CPURead(location uint16) uint8 {
	// PRG Banks
	// PRG map mode → 	$8000.D6 = 0 	$8000.D6 = 1
	// CPU Bank 		Value of MMC3 register
	// $8000-$9FFF 		R6 				(-2)
	// $A000-$BFFF 		R7 				R7
	// $C000-$DFFF 		(-2) 			R6
	// $E000-$FFFF 		(-1) 			(-1)
	// (-1) : the last bank
	// (-2) : the second last bank
	mapMode := m.bankSelect >> 6 & 0b1
	switch {
	case 0x6000 <= location && location <= 0x7FFF && m.programRamProtect>>7 == 1:
		return m.programRam[location-0x6000]
	case 0x8000 <= location && location <= 0x9FFF && mapMode == 0:
		return m.cartridge.PrgRom[uint32(location-0x8000)+uint32(0x2000)*uint32(m.bankSelections[6])]
	case 0xA000 <= location && location <= 0xBFFF && mapMode == 0:
		return m.cartridge.PrgRom[uint32(location-0xA000)+uint32(0x2000)*uint32(m.bankSelections[7])]
	case 0xC000 <= location && location <= 0xDFFF && mapMode == 0:
		return m.cartridge.PrgRom[uint32(location-0xC000)+uint32(0x2000)*uint32(m.cartridge.PrgRomSize*2-2)]
	case 0xE000 <= location && mapMode == 0:
		return m.cartridge.PrgRom[uint32(location-0xE000)+uint32(0x2000)*uint32(m.cartridge.PrgRomSize*2-1)]
	case 0x8000 <= location && location <= 0x9FFF && mapMode == 1:
		return m.cartridge.PrgRom[uint32(location-0x8000)+uint32(0x2000)*uint32(m.cartridge.PrgRomSize*2-2)]
	case 0xA000 <= location && location <= 0xBFFF && mapMode == 1:
		return m.cartridge.PrgRom[uint32(location-0xA000)+uint32(0x2000)*uint32(m.bankSelections[7])]
	case 0xC000 <= location && location <= 0xDFFF && mapMode == 1:
		return m.cartridge.PrgRom[uint32(location-0xC000)+uint32(0x2000)*uint32(m.bankSelections[6])]
	case 0xE000 <= location && mapMode == 1:
		return m.cartridge.PrgRom[uint32(location-0xE000)+uint32(0x2000)*uint32(m.cartridge.PrgRomSize*2-1)]
	}
	// Mapper was no responsible for the location
	return 0
}

func (m *Mapper004) CPUWrite(location uint16, data uint8) bool {
	switch {
	case 0x6000 <= location && location <= 0x7FFF && m.programRamProtect>>7 == 1:
		m.programRam[location-0x6000] = data
	case 0x8000 <= location && location <= 0x9FFF && location%2 == 0:
		// Bank select ($8000-$9FFE, even)
		// 7  bit  0
		// ---- ----
		// CPMx xRRR
		// |||   |||
		// |||   +++- Specify which bank register to update on next write to Bank Data register
		// |||          000: R0: Select 2 KB CHR bank at PPU $0000-$07FF (or $1000-$17FF)
		// |||          001: R1: Select 2 KB CHR bank at PPU $0800-$0FFF (or $1800-$1FFF)
		// |||          010: R2: Select 1 KB CHR bank at PPU $1000-$13FF (or $0000-$03FF)
		// |||          011: R3: Select 1 KB CHR bank at PPU $1400-$17FF (or $0400-$07FF)
		// |||          100: R4: Select 1 KB CHR bank at PPU $1800-$1BFF (or $0800-$0BFF)
		// |||          101: R5: Select 1 KB CHR bank at PPU $1C00-$1FFF (or $0C00-$0FFF)
		// |||          110: R6: Select 8 KB PRG ROM bank at $8000-$9FFF (or $C000-$DFFF)
		// |||          111: R7: Select 8 KB PRG ROM bank at $A000-$BFFF
		// ||+------- Nothing on the MMC3, see MMC6
		// |+-------- PRG ROM bank mode (0: $8000-$9FFF swappable,
		// |                                $C000-$DFFF fixed to second-last bank;
		// |                             1: $C000-$DFFF swappable,
		// |                                $8000-$9FFF fixed to second-last bank)
		// +--------- CHR A12 inversion (0: two 2 KB banks at $0000-$0FFF,
		//                                  four 1 KB banks at $1000-$1FFF;
		//                               1: two 2 KB banks at $1000-$1FFF,
		//                                  four 1 KB banks at $0000-$0FFF)
		m.bankSelect = data
	case 0x8000 <= location && location <= 0x9FFF && location%2 == 1:
		// Bank data ($8001-$9FFF, odd)
		// 7  bit  0
		// ---- ----
		// DDDD DDDD
		// |||| ||||
		// ++++-++++- New bank value, based on last value written to Bank select register (mentioned above)
		switch m.bankSelect & 0b111 {
		case 0b000, 0b001:
			// R0 and R1 ignore the bottom bit
			m.bankSelections[m.bankSelect&0b111] = data & 0b11111110
		case 0b110, 0b111:
			// R6 and R7 will ignore the top two bits
			m.bankSelections[m.bankSelect&0b111] = data & 0b00111111
		default:
			m.bankSelections[m.bankSelect&0b111] = data
		}
	case 0xA000 <= location && location <= 0xBFFF && location%2 == 0:
		// Mirroring ($A000-$BFFE, even)
		// 7  bit  0
		// ---- ----
		// xxxx xxxM
		//         |
		//         +- Nametable mirroring (0: vertical; 1: horizontal)
		m.mirrorMode = data & 0b1
	case 0xA000 <= location && location <= 0xBFFF && location%2 == 1:
		// PRG RAM protect ($A001-$BFFF, odd)
		// 7  bit  0
		// ---- ----
		// RWXX xxxx
		// ||||
		// ||++------ Nothing on the MMC3, see MMC6
		// |+-------- Write protection (0: allow writes; 1: deny writes)
		// +--------- PRG RAM chip enable (0: disable; 1: enable)
		m.programRamProtect = data
	case 0xC000 <= location && location <= 0xDFFF && location%2 == 0:
		// IRQ latch ($C000-$DFFE, even)
		// 7  bit  0
		// ---- ----
		// DDDD DDDD
		// |||| ||||
		// ++++-++++- IRQ latch value
		m.irqLatch = data
	case 0xC000 <= location && location <= 0xDFFF && location%2 == 1:
		// IRQ reload ($C001-$DFFF, odd)
		// 7  bit  0
		// ---- ----
		// xxxx xxxx
		m.irqReload = true
	case 0xE000 <= location && location%2 == 0:
		// IRQ disable ($E000-$FFFE, even)
		// 7  bit  0
		// ---- ----
		// xxxx xxxx
		m.irqEnabled = false
	case 0xE000 <= location && location%2 == 1:
		// IRQ enable ($E001-$FFFF, odd)
		// 7  bit  0
		// ---- ----
		// xxxx xxxx
		m.irqEnabled = true
	default:
		return false
	}
	return true
}

func (m *Mapper004) Clock() {
	if m.a12LineLowCounter < 3 {
		m.a12LineLowCounter++
	}
}

func (m *Mapper004) PPUMap(location uint16) uint16 {
	if location>>12&0b1 == 1 && m.a12LineLowCounter != 3 {
		m.a12LineLowCounter = 0
	}

	if location>>12&0b1 == 1 && m.a12LineLowCounter == 3 {
		m.a12LineLowCounter = 0
		// Rising edge
		// If the IRQ counter is zero and IRQs are enabled ($E001), an IRQ is triggered. The "alternate revision"
		// checks the IRQ counter transition 1→0, whether from decrementing or reloading.
		if m.irqCounter == 0 && m.irqEnabled {
			m.cartridge.Bus.IRQ()
		}

		// When the IRQ is clocked (filtered A12 0→1), the counter value is checked - if zero or the reload flag
		// is true, it's reloaded with the IRQ latched value at $C000; otherwise, it decrements.
		if m.irqCounter == 0 || m.irqReload {
			m.irqCounter = m.irqLatch
			m.irqReload = false
		} else {
			m.irqCounter--
		}
	}

	if 0x2000 <= location && location <= 0x3EFF {
		if 0x3000 <= location && location <= 0x3FFF {
			location -= 0x1000
		}
		if m.mirrorMode&0b1 == 1 {
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

func (m *Mapper004) PPURead(location uint16) uint8 {
	// CHR Banks
	// CHR map mode → $8000.D7 = 0 	$8000.D7 = 1
	// PPU Bank 	  Value of MMC3 register
	// $0000-$03FF 	  R0 			R2
	// $0400-$07FF 	  R3
	// $0800-$0BFF 	  R1		 	R4
	// $0C00-$0FFF 	  R5
	// $1000-$13FF 	  R2		 	R0
	// $1400-$17FF 	  R3
	// $1800-$1BFF 	  R4 			R1
	// $1C00-$1FFF 	  R5

	mapMode := m.bankSelect >> 7 & 0b1
	switch {
	case location <= 0x03FF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location)+uint32(0x0400)*uint32(m.bankSelections[0])]
	case 0x0400 <= location && location <= 0x07FF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x0400)+uint32(0x0400)*uint32(m.bankSelections[0]+1)]
	case 0x0800 <= location && location <= 0x0BFF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x0800)+uint32(0x0400)*uint32(m.bankSelections[1])]
	case 0x0C00 <= location && location <= 0x0FFF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x0C00)+uint32(0x0400)*uint32(m.bankSelections[1]+1)]
	case 0x1000 <= location && location <= 0x13FF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x1000)+uint32(0x0400)*uint32(m.bankSelections[2])]
	case 0x1400 <= location && location <= 0x17FF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x1400)+uint32(0x0400)*uint32(m.bankSelections[3])]
	case 0x1800 <= location && location <= 0x1BFF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x1800)+uint32(0x0400)*uint32(m.bankSelections[4])]
	case 0x1C00 <= location && location <= 0x1FFF && mapMode == 0:
		return m.cartridge.ChrRom[uint32(location-0x1C00)+uint32(0x0400)*uint32(m.bankSelections[5])]
	case location <= 0x03FF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location)+uint32(0x0400)*uint32(m.bankSelections[2])]
	case 0x0400 <= location && location <= 0x07FF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x400)+uint32(0x0400)*uint32(m.bankSelections[3])]
	case 0x0800 <= location && location <= 0x0BFF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x0800)+uint32(0x0400)*uint32(m.bankSelections[4])]
	case 0x0C00 <= location && location <= 0x0FFF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x0C00)+uint32(0x0400)*uint32(m.bankSelections[5])]
	case 0x1000 <= location && location <= 0x13FF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x1000)+uint32(0x0400)*uint32(m.bankSelections[0])]
	case 0x1400 <= location && location <= 0x17FF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x1400)+uint32(0x0400)*uint32(m.bankSelections[0]+1)]
	case 0x1800 <= location && location <= 0x1BFF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x1800)+uint32(0x0400)*uint32(m.bankSelections[1])]
	case 0x1C00 <= location && location <= 0x1FFF && mapMode == 1:
		return m.cartridge.ChrRom[uint32(location-0x1C00)+uint32(0x0400)*uint32(m.bankSelections[1]+1)]

	}
	// Mapper was not responsible for the location
	return 0
}

func (m *Mapper004) PPUWrite(location uint16, data uint8) bool {
	// No support for CHR Ram
	return false
}

func (m *Mapper004) Reset() {
	m.bankSelections = [8]uint8{}
	m.bankSelect = 0
	m.irqEnabled = false
	m.irqCounter = 0
	m.irqLatch = 0
	m.irqReload = false
	m.a12LineLowCounter = 0
}

func (m *Mapper004) DebugDisplay(text *textutil.Text) {
	plz.Just(fmt.Fprint(text, "Cartridge with Mapper 002\n"))
	plz.Just(fmt.Fprintf(text, "PRG ROM Size: %d * 16 KB\n", m.cartridge.PrgRomSize))
	plz.Just(fmt.Fprintf(text, "PRG BANK    : %d \n", m.bankSelect))
	plz.Just(fmt.Fprintf(text, "CHR ROM Size: %d * 8 KB\n", m.cartridge.ChrRomSize))
	plz.Just(fmt.Fprintf(text, "PRG Mode: %d \n", m.bankSelect>>6&0b1))
	plz.Just(fmt.Fprintf(text, "CHR Mode: %d \n", m.bankSelect>>7))
	plz.Just(fmt.Fprintf(text, "R0: %d\n", m.bankSelections[0]))
	plz.Just(fmt.Fprintf(text, "R1: %d\n", m.bankSelections[1]))
	plz.Just(fmt.Fprintf(text, "R2: %d\n", m.bankSelections[2]))
	plz.Just(fmt.Fprintf(text, "R3: %d\n", m.bankSelections[3]))
	plz.Just(fmt.Fprintf(text, "R4: %d\n", m.bankSelections[4]))
	plz.Just(fmt.Fprintf(text, "R5: %d\n", m.bankSelections[5]))
	plz.Just(fmt.Fprintf(text, "R6: %d\n", m.bankSelections[6]))
	plz.Just(fmt.Fprintf(text, "R7: %d\n", m.bankSelections[7]))
	plz.Just(fmt.Fprintf(text, "IRQ: %t\n", m.irqEnabled))
}
