package nes

import "github.com/exp625/gones/nes/cartridge"

type Bus struct {
	RAM       *RAM
	CPU       *C
	PPU       *PPU
	APU       *APU
	Cartridge *cartridge.Cartridge
}

func (b *Bus) CPURead(location uint16) uint8 {
	switch {
	case location <= 0x1FFF:
		_, data := b.RAM.Read(location % 0x0800)
		return data
	case 0x2000 <= location && location <= 0x3FFF:
		_, data :=  b.PPU.Read(0x2000 + location%0x0008)
		return data
	case 0x4000 <= location && location <= 0x4017:
		// TODO: APU and I/O Registers
		return 0
	case 0x4018 <= location && location <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
		return 0
	case 0x4020 <= location:
		_, data :=  b.Cartridge.CPURead(location)
		return data
	default:
		panic("go is wrong")
	}
}

func (b *Bus) CPUWrite(location uint16, data uint8) {
	switch {
	case location <= 0x1FFF:
		b.RAM.Write(location%0x0800, data)
	case 0x2000 <= location && location <= 0x3FFF:
		b.PPU.Write(0x2000+location%0x0008, data)
	case 0x4000 <= location && location <= 0x4017:
		// TODO: APU and I/O Registers
	case 0x4018 <= location && location <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
	case 0x4020 <= location:
		b.Cartridge.CPUWrite(location, data)
	default:
		panic("go is wrong")
	}
}

func (b *Bus) Reset() {
	b.Cartridge.Reset()
	b.RAM.Reset()
	b.CPU.Reset()
	b.PPU.Reset()
	b.APU.Reset()
}
