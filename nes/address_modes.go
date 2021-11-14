package nes

func AccumulatorAddressing() (uint16, uint8) {
	return 0, CPU.A
}

func ImmediateAddress() (uint16, uint8) {
	location := CPU.PC + 1
	return location, CPU.Bus.CPURead(location)
}

func ZeroPageAddressing() (uint16, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(CPU.Bus.CPURead(CPU.PC+1)) & 0x00FF
	return location, CPU.Bus.CPURead(location)
}

func AbsoluteAddressing() (uint16, uint8) {
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	location := (high << 8) | low
	return location, CPU.Bus.CPURead(location)
}

func IndexedXZeroPageAddressing() (uint16, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(CPU.Bus.CPURead(CPU.PC+1)) & 0x00FF
	return location, CPU.Bus.CPURead(location)
}

func IndexedYZeroPageAddressing() (uint16, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(CPU.Bus.CPURead(CPU.PC+1+uint16(CPU.Y))) & 0x00FF
	return location, CPU.Bus.CPURead(location)
}

func IndexedXAbsoluteAddressing() (uint16, uint8) {
	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	location := (high << 8) | low

	// Offset location by the value stored in the X register
	location += uint16(CPU.X)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		CPU.CycleCount++
	}
	return location, CPU.Bus.CPURead(location)
}

func IndexedYAbsoluteAddressing() (uint16, uint8) {
	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	location := (high << 8) | low

	// Offset location by the value stored in the X register
	location += uint16(CPU.Y)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		CPU.CycleCount++
	}
	return location, CPU.Bus.CPURead(location)
}

func ImpliedAddressing() (uint16, uint8) {
	return 0, 0
}

func RelativeAddressing() (uint16, uint8) {
	// Read the address as 8-bit signed offset relative to the current PC
	offset := CPU.Bus.CPURead(CPU.PC + 1)
	// Offset is negative
	if offset&0x80 == 0x80 {
		return CPU.PC - uint16(offset), 0
	} else {
		return CPU.PC + uint16(offset), 0
	}

}

func IndexedIndirectAddressing() (uint16, uint8) {
	// Build pointer from high and low bit
	pointer := CPU.PC + 1

	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead((pointer + uint16(CPU.X)) & 0x00FF))
	high := uint16(CPU.Bus.CPURead((pointer + uint16(CPU.X) + 1) & 0x00FF))

	location := (high << 8) | low
	location += uint16(CPU.Y)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		CPU.CycleCount++
	}

	return location, CPU.Bus.CPURead(location)
}

func IndirectIndexedAddressing() (uint16, uint8) {
	// Build pointer from high and low bit
	pointer := CPU.PC + 1

	// Build location from high and low bits
	low := uint16(CPU.Bus.CPURead(pointer & 0x00FF))
	high := uint16(CPU.Bus.CPURead((pointer + 1) & 0x00FF))

	location := (high << 8) | low
	location += uint16(CPU.Y)

	// high bits increased due to x offset
	if location&0xFF00 != (high << 8) {
		CPU.CycleCount++
	}

	return location, CPU.Bus.CPURead(location)
}

func AbsoluteIndirect() (uint16, uint8) {
	// Build pointer from high and low bits
	low := uint16(CPU.Bus.CPURead(CPU.PC + 1))
	high := uint16(CPU.Bus.CPURead(CPU.PC + 2))
	pointer := (high << 8) | low

	// Build location from high and low bits
	low = uint16(CPU.Bus.CPURead(pointer))
	high = uint16(CPU.Bus.CPURead(pointer + 1))
	location := (high << 8) | low
	return location, CPU.Bus.CPURead(location)
}
