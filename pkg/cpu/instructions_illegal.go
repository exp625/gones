package cpu

// ALR https://www.masswerk.at/6502/6502_instruction_set.html#ALR
// AND oper + LSR
// A AND oper, 0 -> [76543210] -> C
func (cpu *CPU) ALR(_ uint16, data uint8, length uint16) {
	// AND Memory with Operand
	temp := cpu.A & data
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.Set(FlagZero, (temp&0x00FF) == 0)
	// Set Carry flag
	cpu.Set(FlagCarry, (temp&0x01) == 1)
	// Advance program counter
	cpu.PC += length

}

// ANC https://www.masswerk.at/6502/6502_instruction_set.html#ANC
// AND oper + set C as ASL
// A AND oper, bit(7) -> C
func (cpu *CPU) ANC(_ uint16, data uint8, length uint16) {
	// AND Memory with Operand
	temp := cpu.A & data
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.Set(FlagZero, (temp&0x00FF) == 0)
	// Set Carry flag
	cpu.Set(FlagCarry, (temp>>7)&0x01 == 1)
	// Advance program counter
	cpu.PC += length
}

// ANE  https://www.masswerk.at/6502/6502_instruction_set.html#ANE
// Unstable / Unimplemented
func (cpu *CPU) ANE(_ uint16, _ uint8, _ uint16) {}

// ARR https://www.masswerk.at/6502/6502_instruction_set.html#ARR
// Unused by NES Games
func (cpu *CPU) ARR(_ uint16, _ uint8, _ uint16) {}

// DCP https://www.masswerk.at/6502/6502_instruction_set.html#DCP
// DEC oper + CMP oper
// M - 1 -> M, A - M
func (cpu *CPU) DCP(location uint16, _ uint8, length uint16) {
	// M - 1 -> M
	temp1 := cpu.Bus.CPURead(location) - 1
	cpu.Bus.CPUWrite(location, temp1)
	// A - M
	temp2 := cpu.A - temp1
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp2>>7)&0x01 == 1)
	cpu.Set(FlagZero, temp2 == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	cpu.Set(FlagCarry, cpu.A >= temp1)
	// Advance program counter
	cpu.PC += length
}

// ILLNOP https://www.masswerk.at/6502/6502_instruction_set.html#NOPs
// Instructions effecting in 'no operations' in various address modes. Operands are ignored
func (cpu *CPU) ILLNOP(_ uint16, _ uint8, length uint16) {
	// Advance program counter
	cpu.PC += length
	return
}

// ISC https://www.masswerk.at/6502/6502_instruction_set.html#ISC
// INC oper + SBC oper
// M + 1 -> M, A - M - C -> A
// Also ISB for logging
func (cpu *CPU) ISC(location uint16, _ uint8, length uint16) {
	// M + 1 -> M
	temp1 := cpu.Bus.CPURead(location) + 1
	cpu.Bus.CPUWrite(location, temp1)

	// Get carry
	var carry uint8
	if cpu.Get(FlagCarry) {
		carry = 1
	}
	// A - M - ^C = A + ^M + C
	dataInverse := ^temp1
	// Perform calculation in 16 bit
	temp := uint16(cpu.A) + uint16(dataInverse) + uint16(carry)
	tempSigned := int16(int8(cpu.A)) + int16(int8(dataInverse)) + int16(int8(carry))
	// Store last 8th bits to A register
	cpu.A = uint8(temp & 0x00FF)
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	cpu.Set(FlagZero, (temp&0x00FF) == 0)
	// Check if result is greater thant 255. If true we have a carry
	cpu.Set(FlagCarry, temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	cpu.Set(FlagOverflow, tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	cpu.PC += length
}

// LAS  https://www.masswerk.at/6502/6502_instruction_set.html#LAS
// LDA/TSX oper
// M AND SP -> A, X, SP
func (cpu *CPU) LAS(_ uint16, data uint8, length uint16) {
	// M AND SP -> A, X, SP
	temp := data & cpu.S
	cpu.A = temp
	cpu.X = temp
	cpu.S = temp
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.Set(FlagZero, temp == 0)
	// Advance program counter
	cpu.PC += length
}

// LAX  https://www.masswerk.at/6502/6502_instruction_set.html#LAX
// LDA oper + LDX oper
// M -> A -> X
func (cpu *CPU) LAX(_ uint16, data uint8, length uint16) {
	// M -> A -> X
	cpu.A = data
	cpu.X = data
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (data>>7)&0x01 == 1)
	// Check if result is zero
	cpu.Set(FlagZero, data == 0)
	// Advance program counter
	cpu.PC += length
}

// LXA  https://www.masswerk.at/6502/6502_instruction_set.html#LXA
// Unstable / Unimplemented
func (cpu *CPU) LXA(_ uint16, _ uint8, _ uint16) {}

// RLA  https://www.masswerk.at/6502/6502_instruction_set.html#RLA
// ROL oper + AND oper
// M = C <- [76543210] <- C, A AND M -> A
func (cpu *CPU) RLA(location uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if cpu.Get(FlagCarry) {
		carry = 1
	}
	// M = C <- [76543210] <- C
	temp1 := data<<1 + carry
	cpu.Bus.CPUWrite(location, temp1)
	cpu.Set(FlagCarry, (data>>7)&0x01 == 1)
	// A AND M -> A
	temp2 := cpu.A & temp1
	// Store result in A register
	cpu.A = temp2
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp2>>7)&0x01 == 1)
	// Check if result is zero
	cpu.Set(FlagZero, (temp2&0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// RRA  https://www.masswerk.at/6502/6502_instruction_set.html#RRA
// ROR oper + ADC oper
// M = C -> [76543210] -> C, A + M + C -> A, C
func (cpu *CPU) RRA(location uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if cpu.Get(FlagCarry) {
		carry = 0x80
	}
	// M = C -> [76543210] -> C
	temp1 := data>>1 + carry
	cpu.Bus.CPUWrite(location, temp1)
	cpu.Set(FlagCarry, data&0x01 == 1)

	// A + M + C -> A, C
	// Get carry
	if cpu.Get(FlagCarry) {
		carry = 1
	}
	// Perform calculation in 16 bit
	temp := uint16(cpu.A) + uint16(temp1) + uint16(carry)
	tempSigned := int16(int8(cpu.A)) + int16(int8(temp1)) + int16(int8(carry))
	// Store last 8th bits to A register
	cpu.A = uint8(temp & 0x00FF)
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	cpu.Set(FlagZero, (temp&0x00FF) == 0)
	// Check if result is greater than 255. If true we have a carry
	cpu.Set(FlagCarry, temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	cpu.Set(FlagOverflow, tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	cpu.PC += length
}

// SAX  https://www.masswerk.at/6502/6502_instruction_set.html#SAX
// A and X are put on the bus at the same time (resulting effectively in an AND operation) and stored in M
// A AND X -> M
func (cpu *CPU) SAX(location uint16, _ uint8, length uint16) {
	// A AND X -> M
	temp := cpu.A & cpu.X
	cpu.Bus.CPUWrite(location, temp)
	// Advance program counter
	cpu.PC += length
}

// SBX  https://www.masswerk.at/6502/6502_instruction_set.html#SBX
// CMP and DEX at once, sets flags like CMP
// (A AND X) - oper -> X
func (cpu *CPU) SBX(_ uint16, data uint8, length uint16) {
	// (A AND X) - oper -> X
	temp := (cpu.A & cpu.X) - data
	cpu.X = temp
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	cpu.Set(FlagZero, temp == 0)
	// Sets flag like CMP
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	cpu.Set(FlagCarry, cpu.A >= data)
	// Advance program counter
	cpu.PC += length
}

// SHA  https://www.masswerk.at/6502/6502_instruction_set.html#SHA
// Unstable / Unimplemented
func (cpu *CPU) SHA(_ uint16, _ uint8, _ uint16) {}

// SHX  https://www.masswerk.at/6502/6502_instruction_set.html#SHX
// Unstable / Unimplemented
func (cpu *CPU) SHX(_ uint16, _ uint8, _ uint16) {}

// SHY  https://www.masswerk.at/6502/6502_instruction_set.html#SHY
// Unstable / Unimplemented
func (cpu *CPU) SHY(_ uint16, _ uint8, _ uint16) {}

// SRE  https://www.masswerk.at/6502/6502_instruction_set.html#SRE
// LSR oper + EOR oper
// M = 0 -> [76543210] -> C, A EOR M -> A
func (cpu *CPU) SRE(location uint16, data uint8, length uint16) {
	carry := data & 0x01
	cpu.Set(FlagCarry, carry == 1)
	temp1 := data >> 1
	cpu.Bus.CPUWrite(location, temp1)
	temp2 := cpu.A ^ temp1
	cpu.A = temp2
	// Advance program counter
	cpu.PC += length
}

// SLO  https://www.masswerk.at/6502/6502_instruction_set.html#SLO
// ASL oper + ORA oper
// M = C <- [76543210] <- 0, A OR M -> A
func (cpu *CPU) SLO(location uint16, data uint8, length uint16) {
	// M = C <- [76543210] <- 0
	// Get carry from data
	carry := (data >> 7) & 0x01
	// Shift one bit left
	temp1 := data << 1
	cpu.Bus.CPUWrite(location, temp1)
	// Set carry
	cpu.Set(FlagCarry, carry == 1)

	// A OR M -> A
	temp2 := cpu.A | temp1
	cpu.A = temp2
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp2>>7)&0x01 == 1)
	// Check if result is zero
	cpu.Set(FlagZero, (temp2&0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// TAS  https://www.masswerk.at/6502/6502_instruction_set.html#TAS
// Unstable / Unimplemented
func (cpu *CPU) TAS(_ uint16, _ uint8, _ uint16) {}

// USBC  https://www.masswerk.at/6502/6502_instruction_set.html#USBC
// SBC oper + NOP
// A - M - C -> A
func (cpu *CPU) USBC(_ uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if cpu.Get(FlagCarry) {
		carry = 1
	}
	// A - M - ^C = A + ^M + C
	dataInverse := ^data
	// Perform calculation in 16 bit
	temp := uint16(cpu.A) + uint16(dataInverse) + uint16(carry)
	tempSigned := int16(int8(cpu.A)) + int16(int8(dataInverse)) + int16(int8(carry))
	// Store last 8th bits to A register
	cpu.A = uint8(temp & 0x00FF)
	// Check if 8th bit is one
	cpu.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	cpu.Set(FlagZero, (temp&0x00FF) == 0)
	// Check if result is greater thant 255. If true we have a carry
	cpu.Set(FlagCarry, temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	cpu.Set(FlagOverflow, tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	cpu.PC += length
}
