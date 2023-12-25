package cpu

// ADC https://www.masswerk.at/6502/6502_instruction_set.html#ADC
// Add Memory to Accumulator with Carry
// A + M + C -> A, C
func (cpu *CPU) ADC(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// Get carry
	var carry uint8
	if cpu.P.Carry() {
		carry = 1
	}
	// Perform calculation in 16 bit
	temp := uint16(cpu.A) + uint16(data) + uint16(carry)
	tempSigned := int16(int8(cpu.A)) + int16(int8(data)) + int16(int8(carry))
	// Store last 8th bits to A register
	cpu.A = uint8(temp & 0x00FF)
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Check if result is greater than 255. If true we have a carry
	cpu.P.SetCarry(temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	cpu.P.SetOverflow(tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	cpu.PC += length
}

// AND https://www.masswerk.at/6502/6502_instruction_set.html#AND
// AND Memory with Accumulator
// A AND M -> A
func (cpu *CPU) AND(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// AND Memory with Accumulator
	temp := cpu.A & data
	// Store result in A register
	cpu.A = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// ASL https://www.masswerk.at/6502/6502_instruction_set.html#ASL
// Shift Left One Bit (Memory or Accumulator)
// C <- [76543210] <- 0
func (cpu *CPU) ASL(location uint16, length uint16) {
	// Check if we have Accumulator addressing
	opcode := cpu.Bus.CPURead(cpu.PC)
	inst := cpu.Instructions[opcode]
	data := cpu.A
	if inst.ClockCycles != 2 {
		data = cpu.Bus.CPURead(location)
	}
	// Get carry from data
	carry := (data >> 7) & 0x01
	// Shift one bit left
	temp := data << 1
	// Set carry
	cpu.P.SetCarry(carry == 1)
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		cpu.A = temp
	} else {
		cpu.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	cpu.PC += length
}

// BCC https://www.masswerk.at/6502/6502_instruction_set.html#BCC
// Branch on Carry Clear
// branch on C = 0
func (cpu *CPU) BCC(location uint16, length uint16) {
	// Check if C = 0
	if !cpu.P.Carry() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BCS https://www.masswerk.at/6502/6502_instruction_set.html#BCS
// Branch on Carry Set
// branch on C = 1
func (cpu *CPU) BCS(location uint16, length uint16) {
	// Check if C = 1
	if cpu.P.Carry() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BEQ https://www.masswerk.at/6502/6502_instruction_set.html#BEQ
// Branch on Result Zero
// branch on Z = 1
func (cpu *CPU) BEQ(location uint16, length uint16) {
	// Check if Z = 1
	if cpu.P.Zero() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BIT https://www.masswerk.at/6502/6502_instruction_set.html#BIT
// Test Bits in Memory with Accumulator
// bits 7 and 6 of operand are transferred to bit 7 and 6 of P (N,V);
// the zero-flag is set to the result of operand AND accumulator.
// A AND M, M7 -> N, M6 -> V
func (cpu *CPU) BIT(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// AND Memory with Accumulator
	temp := cpu.A & data
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Transfer bit 7 to N
	cpu.P.SetNegative((data>>7)&0x01 == 1)
	// Transfer bit 6 to V
	cpu.P.SetOverflow((data>>6)&0x01 == 1)
	// Advance program counter
	cpu.PC += length
}

// BMI https://www.masswerk.at/6502/6502_instruction_set.html#BMI
// Branch on Result negative
// branch on N = 1
func (cpu *CPU) BMI(location uint16, length uint16) {
	// Check if N = 1
	if cpu.P.Negative() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BNE https://www.masswerk.at/6502/6502_instruction_set.html#BMI
// Branch on Result not Zero
// branch on Z = 0
func (cpu *CPU) BNE(location uint16, length uint16) {
	// Check if Z = 0
	if !cpu.P.Zero() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BPL https://www.masswerk.at/6502/6502_instruction_set.html#BMI
// Branch on Result positive
// branch on N = 0
func (cpu *CPU) BPL(location uint16, length uint16) {
	// Check if N = 0
	if !cpu.P.Negative() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BRK https://www.masswerk.at/6502/6502_instruction_set.html#BRK
// BRK initiates a software interrupt similar to a hardware
// interrupt (IRQ). The return address pushed to the stack is
// PC+2, providing an extra byte of spacing for a break mark
// (identifying a reason for the break.)
// The status register will be pushed to the stack with the break
// flag set to 1. However, when retrieved during RTI or by a PLP
// instruction, the break flag will be ignored.
// The interrupt disable flag is not set automatically.
func (cpu *CPU) BRK(uint16, uint16) {
	// Get current pc + 2
	pc := cpu.PC + 2
	// Store high bytes of pc to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8((pc>>8)&0x00FF))
	cpu.S--
	// Store low bytes of pc to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(pc&0x00FF))
	cpu.S--
	// Set flags and store current pc onto stack
	// From https://wiki.nesdev.org/w/index.php?title=Status_flags
	// In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(cpu.P)|FlagBreak|FlagInterruptDisable)
	cpu.S--
	// Set Interrupt disable flag
	// We don't want another interrupt inside the interrupt handler
	cpu.P.SetInterruptDisable(true)
	// Get pc from IRQ/BRK vector and jump to location
	low := uint16(cpu.Bus.CPURead(IRQVector))
	high := uint16(cpu.Bus.CPURead(IRQVector + 1))
	cpu.PC = (high << 8) | low
}

// BVC https://www.masswerk.at/6502/6502_instruction_set.html#BVC
// Branch on Overflow clear
// branch on V = 0
func (cpu *CPU) BVC(location uint16, length uint16) {
	// Check if V = 0
	if !cpu.P.Overflow() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// BVS https://www.masswerk.at/6502/6502_instruction_set.html#BVS
// Branch on Overflow set
// branch on V = 1
func (cpu *CPU) BVS(location uint16, length uint16) {
	// Check if V = 1
	if cpu.P.Overflow() {
		// Taking a branch takes one additional cycle
		cpu.CycleCount++
		if (cpu.PC+length)&0xFF00 != location&0xFF00 {
			// Branching to a different pages takes one additional cycle
			cpu.CycleCount++
		}
		// Jump to new location
		cpu.PC = location
		return
	}
	// Advance program counter
	cpu.PC += length
}

// CLC https://www.masswerk.at/6502/6502_instruction_set.html#CLC
// Clear carry flag
// 0 -> C
func (cpu *CPU) CLC(_ uint16, length uint16) {
	// Clear Flag
	cpu.P.SetCarry(false)
	// Advance program counter
	cpu.PC += length
}

// CLD https://www.masswerk.at/6502/6502_instruction_set.html#CLD
// Clear decimal mode
// 0 -> D
func (cpu *CPU) CLD(_ uint16, length uint16) {
	// Clear Flag
	cpu.P.SetDecimal(false)
	// Advance program counter
	cpu.PC += length
}

// CLI https://www.masswerk.at/6502/6502_instruction_set.html#CLI
// Clear interrupt disable bit
// 0 -> I
func (cpu *CPU) CLI(_ uint16, length uint16) {
	// Clear Flag
	cpu.P.SetInterruptDisable(false)
	// Advance program counter
	cpu.PC += length
	// Clear IRQ line
	cpu.IRQLinePreviousCycle = false
}

// CLV https://www.masswerk.at/6502/6502_instruction_set.html#CLV
// Clear overflow flag
// 0 -> V
func (cpu *CPU) CLV(_ uint16, length uint16) {
	// Clear Flag
	cpu.P.SetOverflow(false)
	// Advance program counter
	cpu.PC += length
}

// CMP https://www.masswerk.at/6502/6502_instruction_set.html#CMP
// Compare Memory with Accumulator
// A - M
func (cpu *CPU) CMP(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// A - M
	temp := cpu.A - data
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	cpu.P.SetZero(temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	cpu.P.SetCarry(cpu.A >= data)
	// Advance program counter
	cpu.PC += length
}

// CPX https://www.masswerk.at/6502/6502_instruction_set.html#CPX
// Compare Memory with X
// X - M
func (cpu *CPU) CPX(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// X - M
	temp := cpu.X - data
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	cpu.P.SetCarry(cpu.X >= data)
	// Advance program counter
	cpu.PC += length
}

// CPY https://www.masswerk.at/6502/6502_instruction_set.html#CPY
// Compare Memory with Y
// Y - M
func (cpu *CPU) CPY(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// Y - M
	temp := cpu.Y - data
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	cpu.P.SetCarry(cpu.Y >= data)
	// Advance program counter
	cpu.PC += length
}

// DEC https://www.masswerk.at/6502/6502_instruction_set.html#DEC
// Decrement Memory by One
// M - 1 -> M
func (cpu *CPU) DEC(location uint16, length uint16) {
	// M - 1 -> M
	temp := cpu.Bus.CPURead(location) - 1
	cpu.Bus.CPUWrite(location, temp)
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(temp == 0)
	// Advance program counter
	cpu.PC += length
}

// DEX https://www.masswerk.at/6502/6502_instruction_set.html#DEX
// Decrement X by One
// X - 1 -> X
func (cpu *CPU) DEX(_ uint16, length uint16) {
	// X - 1 -> X
	temp := cpu.X - 1
	cpu.X = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(temp == 0)
	// Advance program counter
	cpu.PC += length
}

// DEY https://www.masswerk.at/6502/6502_instruction_set.html#DEY
// Decrement Y by One
// Y - 1 -> Y
func (cpu *CPU) DEY(_ uint16, length uint16) {
	// Y - 1 -> Y
	temp := cpu.Y - 1
	cpu.Y = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(temp == 0)
	// Advance program counter
	cpu.PC += length
}

// EOR https://www.masswerk.at/6502/6502_instruction_set.html#EOR
// Exclusive-OR Memory with Accumulator
// A ^ M -> A
func (cpu *CPU) EOR(location uint16, length uint16) {
	// A ^ M -> A
	temp := cpu.A ^ cpu.Bus.CPURead(location)
	cpu.A = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(temp == 0)
	// Advance program counter
	cpu.PC += length
}

// INC https://www.masswerk.at/6502/6502_instruction_set.html#INC
// Increment Memory by One
// M + 1 -> M
func (cpu *CPU) INC(location uint16, length uint16) {
	// M + 1 -> M
	temp := cpu.Bus.CPURead(location) + 1
	cpu.Bus.CPUWrite(location, temp)
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// INX https://www.masswerk.at/6502/6502_instruction_set.html#INX
// Increment X by One
// X + 1 -> X
func (cpu *CPU) INX(_ uint16, length uint16) {
	// X + 1 -> X
	temp := cpu.X + 1
	cpu.X = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// INY https://www.masswerk.at/6502/6502_instruction_set.html#INY
// Increment Y by One
// Y + 1 -> Y
func (cpu *CPU) INY(_ uint16, length uint16) {
	// Y + 1 -> Y
	temp := cpu.Y + 1
	cpu.Y = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// JMP https://www.masswerk.at/6502/6502_instruction_set.html#JMP
// Jump to New Location
func (cpu *CPU) JMP(location uint16, _ uint16) {
	// Jump to New Location
	cpu.PC = location
}

// JSR https://www.masswerk.at/6502/6502_instruction_set.html#JSR
// Jump to New Location
// push (PC+2)
func (cpu *CPU) JSR(location uint16, _ uint16) {
	// Get PC+2
	pc := cpu.PC + 2
	// Store high bytes of pc+2 to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8((pc>>8)&0x00FF))
	cpu.S--
	// Store low bytes of pc+2 to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(pc&0x00FF))
	cpu.S--
	// Jump to New Location
	cpu.PC = location
}

// LDA https://www.masswerk.at/6502/6502_instruction_set.html#LDA
// Load Accumulator with Memory
// M -> A
func (cpu *CPU) LDA(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// M -> A
	cpu.A = data
	// Check if 8th bit is one
	cpu.P.SetNegative((data>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(data == 0)
	// Advance program counter
	cpu.PC += length
}

// LDX https://www.masswerk.at/6502/6502_instruction_set.html#LDX
// Load X with Memory
// M -> X
func (cpu *CPU) LDX(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// M -> X
	cpu.X = data
	// Check if 8th bit is one
	cpu.P.SetNegative((data>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(data == 0)
	// Advance program counter
	cpu.PC += length
}

// LDY https://www.masswerk.at/6502/6502_instruction_set.html#LDY
// Load Y with Memory
// M -> Y
func (cpu *CPU) LDY(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// M -> Y
	cpu.Y = data
	// Check if 8th bit is one
	cpu.P.SetNegative((data>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(data == 0)
	// Advance program counter
	cpu.PC += length
}

// LSR https://www.masswerk.at/6502/6502_instruction_set.html#LSR
// Shift One Bit Right (Memory or Accumulator)
// 0 -> [76543210] -> C
func (cpu *CPU) LSR(location uint16, length uint16) {
	// Check if we have Accumulator addressing
	opcode := cpu.Bus.CPURead(cpu.PC)
	inst := cpu.Instructions[opcode]
	data := cpu.A
	if inst.ClockCycles != 2 {
		data = cpu.Bus.CPURead(location)
	}
	temp := data >> 1
	cpu.P.SetCarry(data&0x01 == 1)
	cpu.P.SetNegative(false)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		cpu.A = temp
	} else {
		cpu.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	cpu.PC += length
}

// NOP https://www.masswerk.at/6502/6502_instruction_set.html#NOP
// No Operation
func (cpu *CPU) NOP(_ uint16, length uint16) {
	// Advance program counter
	cpu.PC += length
	return
}

// ORA https://www.masswerk.at/6502/6502_instruction_set.html#ORA
// OR Memory with Accumulator
// A | M -> A
func (cpu *CPU) ORA(location uint16, length uint16) {
	// A | M -> A
	temp := cpu.A | cpu.Bus.CPURead(location)
	cpu.A = temp
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Advance program counter
	cpu.PC += length
}

// PHA https://www.masswerk.at/6502/6502_instruction_set.html#PHA
// Push Accumulator on Stack
// push A
func (cpu *CPU) PHA(_ uint16, length uint16) {
	// Push A onto stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), cpu.A)
	cpu.S--
	// Advance program counter
	cpu.PC += length
}

// PHP https://www.masswerk.at/6502/6502_instruction_set.html#PHP
// Push Processor Status on Stack
// The status register will be pushed with the break flag and bit 5 set to 1.
// push P
func (cpu *CPU) PHP(_ uint16, length uint16) {
	// push P with bit 5 and 6 set to 1
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(cpu.P)|FlagBreak|FlagUnused)
	cpu.S--
	// Advance program counter
	cpu.PC += length
}

// PLA https://www.masswerk.at/6502/6502_instruction_set.html#PLA
// Pull Accumulator from Stack
// pull A
func (cpu *CPU) PLA(_ uint16, length uint16) {
	// pull A
	cpu.S++
	cpu.A = cpu.Bus.CPURead(StackPage | uint16(cpu.S))
	// Check if 8th bit is one
	cpu.P.SetNegative((cpu.A>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(cpu.A == 0)
	// Advance program counter
	cpu.PC += length
}

// PLP https://www.masswerk.at/6502/6502_instruction_set.html#PLP
// Pull Processor Status  from Stack
// The status register will be pulled with the break flag and bit 5 ignored
// pull P
func (cpu *CPU) PLP(_ uint16, length uint16) {
	// pull p
	cpu.S++
	temp := cpu.Bus.CPURead(StackPage | uint16(cpu.S))
	// Ignore bit 4 and 5 from Stack but keep the value of bit 4 and 5 on the P register
	// Only bit 4 and 5 | Value from Stack without bit 4 and 5
	cpu.P = StatusRegister((uint8(cpu.P) & (FlagBreak | FlagUnused)) | temp & ^(FlagBreak|FlagUnused))
	// Advance program counter
	cpu.PC += length
	// Clear IRQ line
	cpu.IRQRequested = false
}

// ROL https://www.masswerk.at/6502/6502_instruction_set.html#ROL
// Rotate One Bit Left (Memory or Accumulator)
// C <- [76543210] <- C
func (cpu *CPU) ROL(location uint16, length uint16) {
	// Check if we have Accumulator addressing
	opcode := cpu.Bus.CPURead(cpu.PC)
	inst := cpu.Instructions[opcode]
	data := cpu.A
	if inst.ClockCycles != 2 {
		data = cpu.Bus.CPURead(location)
	}
	// Get carry
	var carry uint8
	if cpu.P.Carry() {
		carry = 1
	}
	// C <- [76543210] <- C
	temp := data<<1 + carry
	cpu.P.SetCarry((data>>7)&0x01 == 1)
	// Check if 8th bit is one
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		cpu.A = temp
	} else {
		cpu.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	cpu.PC += length
}

// ROR https://www.masswerk.at/6502/6502_instruction_set.html#ROR
// Rotate One Bit Right (Memory or Accumulator)
// C -> [76543210] -> C
func (cpu *CPU) ROR(location uint16, length uint16) {
	// Check if we have Accumulator addressing
	opcode := cpu.Bus.CPURead(cpu.PC)
	inst := cpu.Instructions[opcode]
	data := cpu.A
	if inst.ClockCycles != 2 {
		data = cpu.Bus.CPURead(location)
	}
	// Get carry
	var carry uint8
	if cpu.P.Carry() {
		carry = 0x80
	}
	// C <- [76543210] <- C
	temp := data>>1 + carry
	cpu.P.SetCarry(data&0x01 == 1)
	// Check if result is zero
	cpu.P.SetNegative(temp&0x80 == 0x80)
	// Check if result is zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		cpu.A = temp
	} else {
		cpu.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	cpu.PC += length
}

// RTI https://www.masswerk.at/6502/6502_instruction_set.html#RTI
// Return from Interrupt
// The status register is pulled with the break flag and bit 5 ignored. Then PC is pulled from the stack.
// pull P, pull PC
func (cpu *CPU) RTI(uint16, uint16) {
	// pull P from stack
	cpu.S++
	status := cpu.Bus.CPURead(0x0100 + uint16(cpu.S))
	// Ignore bit 4 and 5 from Stack but keep the value of bit 4 and 5 on the PC
	// Only bit 4 and 5 | Value from Stack without bit 4 and 5
	cpu.P = StatusRegister((uint8(cpu.P) & (FlagBreak | FlagUnused)) | status & ^(FlagBreak|FlagUnused))
	// pull low bits of pc from stack
	cpu.S++
	low := uint16(cpu.Bus.CPURead(0x0100 + uint16(cpu.S)))
	// pull high bits of pc from stack
	cpu.S++
	high := uint16(cpu.Bus.CPURead(0x0100 + uint16(cpu.S)))
	// Set pc to pulled value
	pc := (high << 8) | low
	cpu.PC = pc
}

// RTS https://www.masswerk.at/6502/6502_instruction_set.html#RTS
// Return from Subroutine
// pull PC, PC+1 -> PC
func (cpu *CPU) RTS(uint16, uint16) {
	// pull low bits of pc from stack
	cpu.S++
	low := uint16(cpu.Bus.CPURead(0x0100 + uint16(cpu.S)))
	// pull high bits of pc from stack
	cpu.S++
	high := uint16(cpu.Bus.CPURead(0x0100 + uint16(cpu.S)))
	// Set pc to pulled value (PC+1)
	pc := (high << 8) | low
	// I don't know why???
	cpu.PC = pc + 1
}

// SBC https://www.masswerk.at/6502/6502_instruction_set.html#SBC
// Subtract Memory from Accumulator with Borrow
// A - M - ^C -> A
func (cpu *CPU) SBC(location uint16, length uint16) {
	// Get data from location determined by the address mode
	data := cpu.Bus.CPURead(location)
	// Get carry
	var carry uint8
	if cpu.P.Carry() {
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
	cpu.P.SetNegative((temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	cpu.P.SetZero((temp & 0x00FF) == 0)
	// Check if result is greater thant 255. If true we have a carry
	cpu.P.SetCarry(temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	cpu.P.SetOverflow(tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	cpu.PC += length
}

// SEC https://www.masswerk.at/6502/6502_instruction_set.html#SEC
// Set Carry flag
// 1 -> C
func (cpu *CPU) SEC(_ uint16, length uint16) {
	// Set Flag
	cpu.P.SetCarry(true)
	// Advance program counter
	cpu.PC += length
}

// SED https://www.masswerk.at/6502/6502_instruction_set.html#SED
// Set Decimal flag
// 1 -> D
func (cpu *CPU) SED(_ uint16, length uint16) {
	// Set Flag
	cpu.P.SetDecimal(true)
	// Advance program counter
	cpu.PC += length
}

// SEI https://www.masswerk.at/6502/6502_instruction_set.html#SEI
// Set Interrupt Disable Status
// 1 -> I
func (cpu *CPU) SEI(_ uint16, length uint16) {
	// Set Flag
	cpu.P.SetInterruptDisable(true)
	// Advance program counter
	cpu.PC += length
	// Clear IRQ line
	cpu.IRQLinePreviousCycle = false
}

// STA https://www.masswerk.at/6502/6502_instruction_set.html#STA
// Store Accumulator in Memory
// A -> M
func (cpu *CPU) STA(location uint16, length uint16) {
	// A -> M
	cpu.Bus.CPUWrite(location, cpu.A)
	// Advance program counter
	cpu.PC += length
}

// STX https://www.masswerk.at/6502/6502_instruction_set.html#STX
// Store X in Memory
// X -> M
func (cpu *CPU) STX(location uint16, length uint16) {
	// X -> M
	cpu.Bus.CPUWrite(location, cpu.X)
	// Advance program counter
	cpu.PC += length
}

// STY https://www.masswerk.at/6502/6502_instruction_set.html#STY
// Store Y in Memory
// Y -> M
func (cpu *CPU) STY(location uint16, length uint16) {
	// Y -> M
	cpu.Bus.CPUWrite(location, cpu.Y)
	// Advance program counter
	cpu.PC += length
}

// TAX https://www.masswerk.at/6502/6502_instruction_set.html#TAX
// Transfer Accumulator to Index X
// A -> X
func (cpu *CPU) TAX(_ uint16, length uint16) {
	// A -> X
	cpu.X = cpu.A
	// Check if 8th bit is one
	cpu.P.SetNegative((cpu.X>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(cpu.X == 0)
	// Advance program counter
	cpu.PC += length
}

// TAY https://www.masswerk.at/6502/6502_instruction_set.html#TAY
// Transfer Accumulator to Index Y
// A -> Y
func (cpu *CPU) TAY(_ uint16, length uint16) {
	cpu.Y = cpu.A
	// Check if 8th bit is one
	cpu.P.SetNegative((cpu.Y>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(cpu.Y == 0)
	// Advance program counter
	cpu.PC += length
}

// TSX https://www.masswerk.at/6502/6502_instruction_set.html#TSX
// Transfer stack pointer to x
// S -> X
func (cpu *CPU) TSX(_ uint16, length uint16) {
	// S -> X
	cpu.X = cpu.S
	// Check if 8th bit is one
	cpu.P.SetNegative((cpu.X>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(cpu.X == 0)
	// Advance program counter
	cpu.PC += length
}

// TXA https://www.masswerk.at/6502/6502_instruction_set.html#TXA
// Transfer X to A
// X -> A
func (cpu *CPU) TXA(_ uint16, length uint16) {
	// X -> A
	val := cpu.X
	cpu.A = val
	// Check if 8th bit is one
	cpu.P.SetNegative((val>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(val == 0)
	// Advance program counter
	cpu.PC += length
}

// TXS https://www.masswerk.at/6502/6502_instruction_set.html#TXS
// Transfer X to stack pointer
// X -> S
func (cpu *CPU) TXS(_ uint16, length uint16) {
	// X -> SP
	cpu.S = cpu.X
	// Advance program counter
	cpu.PC += length
}

// TYA https://www.masswerk.at/6502/6502_instruction_set.html#TYA
// Transfer Y to A
// Y -> A
func (cpu *CPU) TYA(_ uint16, length uint16) {
	// Y -> A
	cpu.A = cpu.Y
	// Check if 8th bit is one
	cpu.P.SetNegative((cpu.A>>7)&0x01 == 1)
	// Check if result is zero
	cpu.P.SetZero(cpu.A == 0)
	// Advance program counter
	cpu.PC += length
}
