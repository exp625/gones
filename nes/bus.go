package nes

type Bus struct {
	RAM       *RAM
	CPU       *CPU
	PPU       *PPU
	APU       *APU
	Cartridge *Catridge
}

func (b *Bus) CPURead(location uint16) uint8 {
	switch {
	case location <= 0x1FFF:
		return b.RAM.Read(location % 0x0800)
	case 0x2000 <= location && location <= 0x3FFF:
		return b.PPU.Read(0x2000 + location%0x0008)
	case 0x4000 <= location && location <= 0x4017:
		return 0 // TODO: APU and I/O Registers
	case 0x4018 <= location && location <= 0x401F:
		return 0 // TODO: APU and I/O functionality that is normally disabled
	case 0x4020 <= location:
		return b.Cartridge.Read(location)
	default:
		panic("go is wrong")
	}
}

func (b *Bus) CPUWrite(location uint16, data uint8) {

}

func (b *Bus) Reset() {
	b.RAM.Reset()
	b.CPU.Reset()
	b.PPU.Reset()
	b.APU.Reset()
}
