package cpu

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
	location := uint16(CPU.Bus.CPURead(CPU.PC+1)) + uint16(CPU.X)
	if location > 0x00FF {
		// location is on next page. However, we want low to warp around and disallow page turn
		location -= 0x0100
	}
	return location, CPU.Bus.CPURead(location), 0
}

func ZPY() (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF

	location := uint16(CPU.Bus.CPURead(CPU.PC+1)) + uint16(CPU.Y)
	if location > 0x00FF {
		// location is on next page. However, we want low to warp around and disallow page turn
		location -= 0x0100
	}
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
	// Only if instruction normally takes 4 clock cycles, this page cross adds another cycle
	if location&0xFF00 != (high << 8) && Instructions[CPU.Bus.CPURead(CPU.PC)].ClockCycles == 4 {
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
	// Only if instruction normally takes 4 clock cycles, this page cross adds another cycle
	if location&0xFF00 != (high << 8) && Instructions[CPU.Bus.CPURead(CPU.PC)].ClockCycles == 4 {
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
	opcode := CPU.Bus.CPURead(CPU.PC)
	inst := Instructions[opcode]
	return uint16(int16(CPU.PC) + int16(offset) + int16(inst.Length)), 0, 0

}

func IDX() (uint16, uint8, uint8) {
	// Build pointer from high and low bit
	offset := uint16(CPU.Bus.CPURead(CPU.PC + 1))

	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(ZeroPage | (offset + uint16(CPU.X)) & 0x00FF))
	if offset & 0x00FF == 0x00FF {
		// offset + 1 is on next page. However, we want low to warp around and disallow page turn
		offset -= 0x0100
	}
	high := uint16(CPU.Bus.CPURead(ZeroPage | (offset + uint16(CPU.X) + 1) & 0x00FF))

	location := (high << 8) | low

	return location, CPU.Bus.CPURead(location), 0
}

func IZY() (uint16, uint8, uint8) {
	// Build pointer from high and low bit
	offset := uint16(CPU.Bus.CPURead(CPU.PC + 1))

	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(offset & 0x00FF))
	if offset & 0x00FF == 0x00FF {
		// offset + 1 is on next page. However, we want low to warp around and disallow page turn
		offset -= 0x0100
	}
	high := uint16(CPU.Bus.CPURead(offset + 1 & 0x00FF))

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
	if pointer & 0x00FF == 0x00FF {
		// pointer + 1 is on next page. However, we want low to warp around and disallow page turn
		pointer -= 0x0100
	}
	high = uint16(CPU.Bus.CPURead(pointer + 1))


	location := (high << 8) | low
	return location, CPU.Bus.CPURead(location), 0
}
