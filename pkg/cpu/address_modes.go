package cpu

// ABS is for the Absolute addressing mode.
// Abbreviation: a
// Fetches the value from a 16-bit address anywhere in memory.
func (cpu *CPU) ABS(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	low := uint16(cpuRead(cpu.PC + 1))
	high := uint16(cpuRead(cpu.PC + 2))
	location := (high << 8) | low
	return location, cpuRead(location), 0
}

// ABX is for the Absolute indexed addressing mode
// Abbreviation: a,x
// val = PEEK(arg + X)
func (cpu *CPU) ABX(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Build location from high and low bits
	low := uint16(cpuRead(cpu.PC + 1))
	high := uint16(cpuRead(cpu.PC + 2))
	location := (high << 8) | low

	// Offset location by the value stored in the X register
	location += uint16(cpu.X)

	// high bits increased due to x offset
	// Only if instruction normally takes 4 clock cycles, this page cross adds another cycle
	if location&0xFF00 != (high<<8) && cpu.Instructions[cpuRead(cpu.PC)].ClockCycles == 4 {
		return location, cpuRead(location), 1
	}
	return location, cpuRead(location), 0
}

// ABY is for the Absolute indexed addressing mode
// Abbreviation: a,y
// val = PEEK(arg + Y)
func (cpu *CPU) ABY(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Build location from high and low bits
	low := uint16(cpuRead(cpu.PC + 1))
	high := uint16(cpuRead(cpu.PC + 2))
	location := (high << 8) | low

	// Offset location by the value stored in the X register
	location += uint16(cpu.Y)

	// high bits increased due to x offset
	// Only if instruction normally takes 4 clock cycles, this page cross adds another cycle
	if location&0xFF00 != (high<<8) && cpu.Instructions[cpuRead(cpu.PC)].ClockCycles == 4 {
		return location, cpuRead(location), 1
	}
	return location, cpuRead(location), 0
}

// ACC is for the Accumulator addressing mode.
// Abbreviation: A
// Many instructions can operate on the accumulator, e.g. LSR A. Some assemblers will treat no operand as an implicit A where applicable.
func (cpu *CPU) ACC(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	return 0, cpu.A, 0
}

// IDX is for the Indexed indirect addressing mode.
// Abbreviation: (d,x)
// val = PEEK(PEEK((arg + X) % 256) + PEEK((arg + X + 1) % 256) * 256)
func (cpu *CPU) IDX(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Build pointer from high and low bit
	offset := uint16(cpuRead(cpu.PC + 1))

	// Build location from high and low bits
	low := uint16(cpuRead(ZeroPage | (offset+uint16(cpu.X))&0x00FF))
	if offset&0x00FF == 0x00FF {
		// offset + 1 is on next page. However, we want low to wrap around and disallow page turn
		offset -= 0x0100
	}
	high := uint16(cpuRead(ZeroPage | (offset+uint16(cpu.X)+1)&0x00FF))

	location := (high << 8) | low

	return location, cpuRead(location), 0
}

// IMM is for the Immediate addressing mode.
// Abbreviation: #v
// Uses the 8-bit operand itself as the value for the operation, rather than fetching a value from a memory address.
func (cpu *CPU) IMM(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	location := cpu.PC + 1
	return location, cpuRead(location), 0
}

// IMP is for the Implicit addressing mode.
// Abbreviation: (none)
// Instructions like RTS or CLC have no address operand, the destination of results are implied.
func (cpu *CPU) IMP(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	return 0, 0, 0
}

// IND is for the Indirect addressing mode.
// Abbreviation: (a)
// The JMP instruction has a special indirect addressing mode that can jump to the address stored in a 16-bit pointer anywhere in memory.
func (cpu *CPU) IND(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Build pointer from high and low bits
	low := uint16(cpuRead(cpu.PC + 1))
	high := uint16(cpuRead(cpu.PC + 2))
	pointer := (high << 8) | low

	// Build location from high and low bits
	low = uint16(cpuRead(pointer))
	if pointer&0x00FF == 0x00FF {
		// pointer + 1 is on next page. However, we want low to wrap around and disallow page turn
		pointer -= 0x0100
	}
	high = uint16(cpuRead(pointer + 1))

	location := (high << 8) | low
	return location, cpuRead(location), 0
}

// IZY is for the Indirect indexed addressing mode.
// Abbreviation: (d),y
// val = PEEK(PEEK(arg) + PEEK((arg + 1) % 256) * 256 + Y)
func (cpu *CPU) IZY(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Build pointer from high and low bit
	offset := uint16(cpuRead(cpu.PC + 1))

	// Build location from high and low bits
	low := uint16(cpuRead(offset & 0x00FF))
	if offset&0x00FF == 0x00FF {
		// offset + 1 is on next page. However, we want low to wrap around and disallow page turn
		offset -= 0x0100
	}
	high := uint16(cpuRead(offset + 1&0x00FF))

	location := (high << 8) | low
	location += uint16(cpu.Y)

	// high bits increased due to x offset
	if location&0xFF00 != (high<<8) && cpu.Instructions[cpuRead(cpu.PC)].ClockCycles != 8 {
		return location, cpuRead(location), 1
	}
	return location, cpuRead(location), 0
}

// REL is for the Relative addressing mode.
// Abbreviation: label
// Branch instructions (e.g. BEQ, BCS) have a relative addressing mode that specifies an 8-bit signed offset relative to the current PC.
func (cpu *CPU) REL(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Read the address as 8-bit signed offset relative to the current PC
	offset := int8(cpuRead(cpu.PC + 1))
	opcode := cpuRead(cpu.PC)
	inst := cpu.Instructions[opcode]
	return uint16(int16(cpu.PC) + int16(offset) + int16(inst.Length)), 0, 0

}

// ZP0 is for the Zero page addressing mode.
// Abbreviation: d
// Fetches the value from an 8-bit address on the zero page.
func (cpu *CPU) ZP0(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(cpuRead(cpu.PC+1)) & 0x00FF
	return location, cpuRead(location), 0
}

// ZPX is for the Zero page indexed addressing mode.
// Abbreviation: d,x
// val = PEEK((arg + X) % 256)
func (cpu *CPU) ZPX(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF
	location := uint16(cpuRead(cpu.PC+1)) + uint16(cpu.X)
	if location > 0x00FF {
		// location is on next page. However, we want low to wrap around and disallow page turn
		location -= 0x0100
	}
	return location, cpuRead(location), 0
}

// ZPY is for the Zero page indexed addressing mode.
// Abbreviation: d,y
// val = PEEK((arg + Y) % 256)
func (cpu *CPU) ZPY(cpuRead func(location uint16) uint8) (uint16, uint8, uint8) {
	// Zero page is address range 0x0000 - 0x00FF

	location := uint16(cpuRead(cpu.PC+1)) + uint16(cpu.Y)
	if location > 0x00FF {
		// location is on next page. However, we want low to wrap around and disallow page turn
		location -= 0x0100
	}
	return location, cpuRead(location), 0
}
