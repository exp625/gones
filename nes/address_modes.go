package nes

func ACC() (uint16, uint8, uint8) {
	return 0, CPU.A, 0
}

func IMM() (uint16, uint8, uint8) {
	location := CPU.PC + 1
	return location, CPU.Bus.CPURead(location), 0
}

func ZP0() (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(CPU.Bus.CPURead(CPU.PC+1)) & 0x00FF
	return location, CPU.Bus.CPURead(location), 0
}

func ABS() (uint16, uint8, uint8) {
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	location := (high << 8) | low
	return location, CPU.Bus.CPURead(location), 0
}

func ZPX() (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(CPU.Bus.CPURead(CPU.PC+1)) & 0x00FF
	return location, CPU.Bus.CPURead(location), 0
}

func ZPY() (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(CPU.Bus.CPURead(CPU.PC+1+uint16(CPU.Y))) & 0x00FF
	return location, CPU.Bus.CPURead(location), 0
}

func ABX() (uint16, uint8, uint8) {
	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	location := (high << 8) | low

	// Offset location by the value stored in the X register
	location += uint16(CPU.X)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		return location, CPU.Bus.CPURead(location), 1
	}
	return location, CPU.Bus.CPURead(location), 0
}

func ABY() (uint16, uint8, uint8) {
	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	location := (high << 8) | low

	// Offset location by the value stored in the X register
	location += uint16(CPU.Y)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		return location, CPU.Bus.CPURead(location), 1
	}
	return location, CPU.Bus.CPURead(location), 0
}

func IMP() (uint16, uint8, uint8) {
	return 0, 0, 0
}

func REL() (uint16, uint8, uint8) {
	// Read the address as 8-bit signed offset relative to the current PC
	offset := int8(CPU.Bus.CPURead(CPU.PC + 1))
	// Offset is negative
	return uint16(int16(CPU.PC) + int16(offset) + int16(CPU.CurrentInstruction.Length)), 0, 0

}

func IDX() (uint16, uint8, uint8) {
	// Build pointer from high and low bit
	offset := CPU.Bus.CPURead(CPU.PC + 1)

	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(ZeroPage | (uint16(offset) + uint16(CPU.X)) & 0x00FF))
	high := uint16(CPU.Bus.CPURead(ZeroPage | (uint16(offset) + uint16(CPU.X) + 1) & 0x00FF))

	location := (high << 8) | low

	return location, CPU.Bus.CPURead(location), 0
}

func IZY() (uint16, uint8, uint8) {
	// Build pointer from high and low bit
	offset := CPU.Bus.CPURead(CPU.PC + 1)

	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(uint16(offset) & 0x00FF))
	high := uint16(CPU.Bus.CPURead(uint16(offset + 1) & 0x00FF))

	location := (high << 8) | low
	location += uint16(CPU.Y)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		return location, CPU.Bus.CPURead(location), 1
	}
	return location, CPU.Bus.CPURead(location), 0
}

func IND() (uint16, uint8, uint8) {
	// Build pointer from high and low bits
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	pointer := (high << 8) | low

	// Build location from high and low bits
	low = uint16(CPU.Bus.CPURead(pointer))
	high = uint16(CPU.Bus.CPURead(pointer + 1))
	location := (high << 8) | low
	return location, CPU.Bus.CPURead(location), 0
}
