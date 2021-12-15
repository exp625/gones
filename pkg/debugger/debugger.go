package debugger

import (
	"github.com/exp625/gones/pkg/nes"
)

// Debugger struct
type Debugger struct {
	*nes.NES
}

// New creates a new NES instance
func New(nes *nes.NES) *Debugger {
	debugger := &Debugger{
		NES: nes,
	}
	return debugger
}

func (nes *Debugger) CPURead(location uint16) uint8 {
	mappedLocation := nes.Cartridge.CPUMapRead(location)
	switch {
	case mappedLocation <= 0x1FFF:
		// TODO:
		return nes.RAM.Data[mappedLocation%0x0800]
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		// TODO:
		switch (location - 0x2000) % 0x8 {
		case 0:
			return 0
		case 1:
			return 0
		case 2:
			return uint8(nes.PPU.Status)
		case 3:
			return 0
		case 4:
			return nes.PPU.OAM[nes.PPU.OamAddress]
		case 5:
			return 0
		case 6:
			return 0
		case 7:
			return nes.PPURead(uint16(nes.PPU.CurrVRAM))
		default:
			panic("go is wrong")
		}
	case mappedLocation == 0x4016:
		// Debugger should not read advance the shift register
		// return nes.Controller1.SerialRead()
		return 1
	case mappedLocation == 0x4017:
		// Debugger should not read advance the shift register
		// return nes.Controller2.SerialRead()
		return 1
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4015:
		// TODO: APU and I/O Registers
		return 0xFF
	case 0x4018 <= mappedLocation && mappedLocation <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
		return 0
	case 0x4020 <= mappedLocation:
		_, data := nes.Cartridge.CPURead(mappedLocation)
		return data
	default:
		panic("go is wrong")
	}
}

func (nes *Debugger) PPURead(location uint16) uint8 {
	mappedLocation := nes.Cartridge.PPUMapRead(location)
	switch {
	case mappedLocation <= 0x1FFF:
		_, data := nes.Cartridge.PPURead(mappedLocation)
		return data
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3EFF:
		// $3000-$3EFF  -> 	Mirrors of $2000-$2EFF
		if 0x3000 <= mappedLocation && mappedLocation <= 0x3EFF {
			mappedLocation -= 0x1000
		}
		if nes.Cartridge.Mirroring() {
			// 1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
			_, data := nes.VRAM.Read((mappedLocation - 0x2000) % 0x800)
			return data
		} else {
			// 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
			if mappedLocation-0x2000 < 0x800 {
				_, data := nes.VRAM.Read((mappedLocation - 0x2000) % 0x400)
				return data
			} else {
				_, data := nes.VRAM.Read((mappedLocation-0x2000)%0x400 + 0x400)
				return data
			}
		}
	// $3F00-3FFF is not configurable, always mapped to the internal palette control.
	case 0x3F00 <= location:
		// $3F00-$3F1F 	Palette RAM indexes
		// $3F20-$3FFF  Mirrors of $3F00-$3F1F
		if 0x3F00 <= mappedLocation && mappedLocation <= 0x3FFF {
			mirroredLocation := (mappedLocation)%0x0020 + 0x3F00
			// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C. Note that this goes for writing as well as reading.
			if mirroredLocation == 0x3F10 || mirroredLocation == 0x3F14 || mirroredLocation == 0x3F18 || mirroredLocation == 0x3F1C {
				mirroredLocation = 0x3F00
			}
			// Addresses $3F04/$3F08/$3F0C are mirrors of $3F00. Note that this goes for writing as well as reading.
			if mirroredLocation == 0x3F04 || mirroredLocation == 0x3F08 || mirroredLocation == 0x3F0C {
				mirroredLocation = 0x3F00
			}
			data := nes.PPU.PaletteRAM[mirroredLocation-0x3F00]
			return data
		}
	}
	return 0
}
