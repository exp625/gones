package cpu

import (
	"reflect"
	"runtime"
	"strings"
)

type AddressModeFunc func() (uint16, uint8, uint8) // location, data, additional cycles
type ExecuteFunc func(uint16, uint8, uint16)

type Instruction struct {
	AddressMode AddressModeFunc
	Execute     ExecuteFunc
	Length      uint16
	ClockCycles int
}

var Instructions [256]Instruction
var OpCodeMap map[uint8][2]string

func init() {
	// OpCode Matrix from http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf
	Instructions[0x00] = Instruction{IMP, BRK, 1, 7}
	Instructions[0x01] = Instruction{IDX, ORA, 2, 6}
	Instructions[0x02] = Instruction{}
	Instructions[0x03] = Instruction{}
	Instructions[0x04] = Instruction{}
	Instructions[0x05] = Instruction{ZP0, ORA, 2, 3}
	Instructions[0x06] = Instruction{ZP0, ASL, 2, 5}
	Instructions[0x07] = Instruction{}
	Instructions[0x08] = Instruction{IMP, PHP, 1, 3}
	Instructions[0x09] = Instruction{IMM, ORA, 2, 2}
	Instructions[0x0A] = Instruction{ACC, ASL, 1, 2}
	Instructions[0x0B] = Instruction{}
	Instructions[0x0C] = Instruction{}
	Instructions[0x0D] = Instruction{ABS, ORA, 3, 4}
	Instructions[0x0E] = Instruction{ABS, ASL, 3, 6}
	Instructions[0x0F] = Instruction{}

	Instructions[0x10] = Instruction{REL, BPL, 2, 2}
	Instructions[0x11] = Instruction{IZY, ORA, 2, 5}
	Instructions[0x12] = Instruction{}
	Instructions[0x13] = Instruction{}
	Instructions[0x14] = Instruction{}
	Instructions[0x15] = Instruction{ZPX, ORA, 2, 4}
	Instructions[0x16] = Instruction{ZPX, ASL, 2, 6}
	Instructions[0x17] = Instruction{}
	Instructions[0x18] = Instruction{IMP, CLC, 1, 2}
	Instructions[0x19] = Instruction{ABY, ORA, 3, 4}
	Instructions[0x1A] = Instruction{}
	Instructions[0x1B] = Instruction{}
	Instructions[0x1C] = Instruction{}
	Instructions[0x1D] = Instruction{ABX, ORA, 3, 4}
	Instructions[0x1E] = Instruction{ABX, ASL, 3, 7}
	Instructions[0x1F] = Instruction{}

	Instructions[0x20] = Instruction{ABS, JSR, 3, 6}
	Instructions[0x21] = Instruction{IDX, AND, 2, 6}
	Instructions[0x22] = Instruction{}
	Instructions[0x23] = Instruction{}
	Instructions[0x24] = Instruction{ZP0, BIT, 2, 3}
	Instructions[0x25] = Instruction{ZP0, AND, 2, 3}
	Instructions[0x26] = Instruction{ZP0, ROL, 2, 5}
	Instructions[0x27] = Instruction{}
	Instructions[0x28] = Instruction{IMP, PLP, 1, 4}
	Instructions[0x29] = Instruction{IMM, AND, 2, 2}
	Instructions[0x2A] = Instruction{ACC, ROL, 1, 2}
	Instructions[0x2B] = Instruction{}
	Instructions[0x2C] = Instruction{ABS, BIT, 3, 4}
	Instructions[0x2D] = Instruction{ABS, AND, 3, 4}
	Instructions[0x2E] = Instruction{ABS, ROL, 3, 6}
	Instructions[0x2F] = Instruction{}

	Instructions[0x30] = Instruction{REL, BMI, 2, 2}
	Instructions[0x31] = Instruction{IZY, AND, 2, 5}
	Instructions[0x32] = Instruction{}
	Instructions[0x33] = Instruction{}
	Instructions[0x34] = Instruction{}
	Instructions[0x35] = Instruction{ZPX, AND, 2, 4}
	Instructions[0x36] = Instruction{ZPX, ROL, 2, 6}
	Instructions[0x37] = Instruction{}
	Instructions[0x38] = Instruction{IMP, SEC, 1, 2}
	Instructions[0x39] = Instruction{ABY, AND, 3, 4}
	Instructions[0x3A] = Instruction{}
	Instructions[0x3B] = Instruction{}
	Instructions[0x3C] = Instruction{}
	Instructions[0x3D] = Instruction{ABX, AND, 3, 4}
	Instructions[0x3E] = Instruction{ABX, ROL, 3, 7}
	Instructions[0x3F] = Instruction{}

	Instructions[0x40] = Instruction{IMP, RTI, 1, 6}
	Instructions[0x41] = Instruction{IDX, EOR, 2, 6}
	Instructions[0x42] = Instruction{}
	Instructions[0x43] = Instruction{}
	Instructions[0x44] = Instruction{}
	Instructions[0x45] = Instruction{ZP0, EOR, 2, 3}
	Instructions[0x46] = Instruction{ZP0, LSR, 2, 5}
	Instructions[0x47] = Instruction{}
	Instructions[0x48] = Instruction{IMP, PHA, 1, 3}
	Instructions[0x49] = Instruction{IMM, EOR, 2, 2}
	Instructions[0x4A] = Instruction{ACC, LSR, 1, 2}
	Instructions[0x4B] = Instruction{}
	Instructions[0x4C] = Instruction{ABS, JMP, 3, 3}
	Instructions[0x4D] = Instruction{ABS, EOR, 3, 4}
	Instructions[0x4E] = Instruction{ABS, LSR, 3, 6}
	Instructions[0x4F] = Instruction{}

	Instructions[0x50] = Instruction{REL, BVC, 2, 2}
	Instructions[0x51] = Instruction{IZY, EOR, 2, 5}
	Instructions[0x52] = Instruction{}
	Instructions[0x53] = Instruction{}
	Instructions[0x54] = Instruction{}
	Instructions[0x55] = Instruction{ZPX, EOR, 2, 4}
	Instructions[0x56] = Instruction{ZPX, LSR, 2, 6}
	Instructions[0x57] = Instruction{}
	Instructions[0x58] = Instruction{IMP, CLI, 1, 2}
	Instructions[0x59] = Instruction{ABY, EOR, 3, 4}
	Instructions[0x5A] = Instruction{}
	Instructions[0x5B] = Instruction{}
	Instructions[0x5C] = Instruction{}
	Instructions[0x5D] = Instruction{ABX, EOR, 3, 4}
	Instructions[0x5E] = Instruction{ABX, LSR, 3, 7}
	Instructions[0x5F] = Instruction{}

	Instructions[0x60] = Instruction{IMP, RTS, 1, 6}
	Instructions[0x61] = Instruction{IDX, ADC, 2, 6}
	Instructions[0x62] = Instruction{}
	Instructions[0x63] = Instruction{}
	Instructions[0x64] = Instruction{}
	Instructions[0x65] = Instruction{ZP0, ADC, 2, 3}
	Instructions[0x66] = Instruction{ZP0, ROR, 2, 5}
	Instructions[0x67] = Instruction{}
	Instructions[0x68] = Instruction{IMP, PLA, 1, 4}
	Instructions[0x69] = Instruction{IMM, ADC, 2, 2}
	Instructions[0x6A] = Instruction{ACC, ROR, 1, 2}
	Instructions[0x6B] = Instruction{}
	Instructions[0x6C] = Instruction{IND, JMP, 3, 5}
	Instructions[0x6D] = Instruction{ABS, ADC, 3, 4}
	Instructions[0x6E] = Instruction{ABS, ROR, 3, 6}
	Instructions[0x6F] = Instruction{}

	Instructions[0x70] = Instruction{REL, BVS, 2, 2}
	Instructions[0x71] = Instruction{IZY, ADC, 2, 5}
	Instructions[0x72] = Instruction{}
	Instructions[0x73] = Instruction{}
	Instructions[0x74] = Instruction{}
	Instructions[0x75] = Instruction{ZPX, ADC, 2, 4}
	Instructions[0x76] = Instruction{ZPX, ROR, 2, 6}
	Instructions[0x77] = Instruction{}
	Instructions[0x78] = Instruction{IMP, SEI, 1, 2}
	Instructions[0x79] = Instruction{ABY, ADC, 3, 4}
	Instructions[0x7A] = Instruction{}
	Instructions[0x7B] = Instruction{}
	Instructions[0x7C] = Instruction{}
	Instructions[0x7D] = Instruction{ABX, ADC, 3, 4}
	Instructions[0x7E] = Instruction{ABX, ROR, 3, 7}
	Instructions[0x7F] = Instruction{}

	Instructions[0x80] = Instruction{}
	Instructions[0x81] = Instruction{IDX, STA, 2, 6}
	Instructions[0x82] = Instruction{}
	Instructions[0x83] = Instruction{}
	Instructions[0x84] = Instruction{ZP0, STY, 2, 3}
	Instructions[0x85] = Instruction{ZP0, STA, 2, 3}
	Instructions[0x86] = Instruction{ZP0, STX, 2, 3}
	Instructions[0x87] = Instruction{}
	Instructions[0x88] = Instruction{IMP, DEY, 1, 2}
	Instructions[0x89] = Instruction{}
	Instructions[0x8A] = Instruction{IMP, TXA, 1, 2}
	Instructions[0x8B] = Instruction{}
	Instructions[0x8C] = Instruction{ABS, STY, 3, 4}
	Instructions[0x8D] = Instruction{ABS, STA, 3, 4}
	Instructions[0x8E] = Instruction{ABS, STX, 3, 4}
	Instructions[0x8F] = Instruction{}

	Instructions[0x90] = Instruction{REL, BCC, 2, 2}
	Instructions[0x91] = Instruction{IZY, STA, 2, 6}
	Instructions[0x92] = Instruction{}
	Instructions[0x93] = Instruction{}
	Instructions[0x94] = Instruction{ZPX, STY, 2, 4}
	Instructions[0x95] = Instruction{ZPX, STA, 2, 4}
	Instructions[0x96] = Instruction{ZPY, STX, 2, 4}
	Instructions[0x97] = Instruction{}
	Instructions[0x98] = Instruction{IMP, TYA, 1, 2}
	Instructions[0x99] = Instruction{ABY, STA, 3, 5}
	Instructions[0x9A] = Instruction{IMP, TXS, 1, 2}
	Instructions[0x9B] = Instruction{}
	Instructions[0x9C] = Instruction{}
	Instructions[0x9D] = Instruction{ABX, STA, 3, 5}
	Instructions[0x9E] = Instruction{}
	Instructions[0x9F] = Instruction{}

	Instructions[0xA0] = Instruction{IMM, LDY, 2, 2}
	Instructions[0xA1] = Instruction{IDX, LDA, 2, 6}
	Instructions[0xA2] = Instruction{IMM, LDX, 2, 2}
	Instructions[0xA3] = Instruction{}
	Instructions[0xA4] = Instruction{ZP0, LDY, 2, 3}
	Instructions[0xA5] = Instruction{ZP0, LDA, 2, 3}
	Instructions[0xA6] = Instruction{ZP0, LDX, 2, 3}
	Instructions[0xA7] = Instruction{}
	Instructions[0xA8] = Instruction{IMP, TAY, 1, 2}
	Instructions[0xA9] = Instruction{IMM, LDA, 2, 2}
	Instructions[0xAA] = Instruction{IMP, TAX, 1, 2}
	Instructions[0xAB] = Instruction{}
	Instructions[0xAC] = Instruction{ABS, LDY, 3, 4}
	Instructions[0xAD] = Instruction{ABS, LDA, 3, 4}
	Instructions[0xAE] = Instruction{ABS, LDX, 3, 4}
	Instructions[0xAF] = Instruction{}

	Instructions[0xB0] = Instruction{REL, BCS, 2, 2}
	Instructions[0xB1] = Instruction{IZY, LDA, 2, 5}
	Instructions[0xB2] = Instruction{}
	Instructions[0xB3] = Instruction{}
	Instructions[0xB4] = Instruction{ZPX, LDY, 2, 4}
	Instructions[0xB5] = Instruction{ZPX, LDA, 2, 4}
	Instructions[0xB6] = Instruction{ZPY, LDX, 2, 4}
	Instructions[0xB7] = Instruction{}
	Instructions[0xB8] = Instruction{IMP, CLV, 1, 2}
	Instructions[0xB9] = Instruction{ABY, LDA, 3, 4}
	Instructions[0xBA] = Instruction{IMP, TSX, 1, 2}
	Instructions[0xBB] = Instruction{}
	Instructions[0xBC] = Instruction{ABX, LDY, 3, 4}
	Instructions[0xBD] = Instruction{ABX, LDA, 3, 4}
	Instructions[0xBE] = Instruction{ABY, LDX, 3, 4}
	Instructions[0xBF] = Instruction{}

	Instructions[0xC0] = Instruction{IMM, CPY, 2, 2}
	Instructions[0xC1] = Instruction{IDX, CMP, 2, 6}
	Instructions[0xC2] = Instruction{}
	Instructions[0xC3] = Instruction{}
	Instructions[0xC4] = Instruction{ZP0, CPY, 2, 3}
	Instructions[0xC5] = Instruction{ZP0, CMP, 2, 3}
	Instructions[0xC6] = Instruction{ZP0, DEC, 2, 5}
	Instructions[0xC7] = Instruction{}
	Instructions[0xC8] = Instruction{IMP, INY, 1, 2}
	Instructions[0xC9] = Instruction{IMM, CMP, 2, 2}
	Instructions[0xCA] = Instruction{IMP, DEX, 1, 2}
	Instructions[0xCB] = Instruction{}
	Instructions[0xCC] = Instruction{ABS, CPY, 3, 4}
	Instructions[0xCD] = Instruction{ABS, CMP, 3, 4}
	Instructions[0xCE] = Instruction{ABS, DEC, 3, 6}
	Instructions[0xCF] = Instruction{}

	Instructions[0xD0] = Instruction{REL, BNE, 2, 2}
	Instructions[0xD1] = Instruction{IZY, CMP, 2, 5}
	Instructions[0xD2] = Instruction{}
	Instructions[0xD3] = Instruction{}
	Instructions[0xD4] = Instruction{}
	Instructions[0xD5] = Instruction{ZPX, CMP, 2, 4}
	Instructions[0xD6] = Instruction{ZPX, DEC, 2, 6}
	Instructions[0xD7] = Instruction{}
	Instructions[0xD8] = Instruction{IMP, CLD, 1, 2}
	Instructions[0xD9] = Instruction{ABY, CMP, 3, 4}
	Instructions[0xDA] = Instruction{}
	Instructions[0xDB] = Instruction{}
	Instructions[0xDC] = Instruction{}
	Instructions[0xDD] = Instruction{ABX, CMP, 3, 4}
	Instructions[0xDE] = Instruction{ABX, DEC, 3, 7}
	Instructions[0xDF] = Instruction{}

	Instructions[0xE0] = Instruction{IMM, CPX, 2, 2}
	Instructions[0xE1] = Instruction{IDX, SBC, 2, 6}
	Instructions[0xE2] = Instruction{}
	Instructions[0xE3] = Instruction{}
	Instructions[0xE4] = Instruction{ZP0, CPX, 2, 3}
	Instructions[0xE5] = Instruction{ZP0, SBC, 2, 3}
	Instructions[0xE6] = Instruction{ZP0, INC, 2, 5}
	Instructions[0xE7] = Instruction{}
	Instructions[0xE8] = Instruction{IMP, INX, 1, 2}
	Instructions[0xE9] = Instruction{IMM, SBC, 2, 2}
	Instructions[0xEA] = Instruction{IMP, NOP, 1, 2}
	Instructions[0xEB] = Instruction{}
	Instructions[0xEC] = Instruction{ABS, CPX, 3, 4}
	Instructions[0xED] = Instruction{ABS, SBC, 3, 4}
	Instructions[0xEE] = Instruction{ABS, INC, 3, 6}
	Instructions[0xEF] = Instruction{}

	Instructions[0xF0] = Instruction{REL, BEQ, 2, 2}
	Instructions[0xF1] = Instruction{IZY, SBC, 2, 5}
	Instructions[0xF2] = Instruction{}
	Instructions[0xF3] = Instruction{}
	Instructions[0xF4] = Instruction{}
	Instructions[0xF5] = Instruction{ZPX, SBC, 2, 4}
	Instructions[0xF6] = Instruction{ZPX, INC, 2, 6}
	Instructions[0xF7] = Instruction{}
	Instructions[0xF8] = Instruction{IMP, SED, 1, 2}
	Instructions[0xF9] = Instruction{ABY, SBC, 3, 4}
	Instructions[0xFA] = Instruction{}
	Instructions[0xFB] = Instruction{}
	Instructions[0xFC] = Instruction{}
	Instructions[0xFD] = Instruction{ABX, SBC, 3, 4}
	Instructions[0xFE] = Instruction{ABX, INC, 3, 7}
	Instructions[0xFF] = Instruction{}

	// Generate debug map
	OpCodeMap = make(map[uint8][2]string, 0xFF)
	for i := range Instructions {
		if Instructions[i].Length == 0 {
			arr := [2]string{ "ERR" , "ERR"}
			OpCodeMap[uint8(i)] = arr
		} else {
			str1 := runtime.FuncForPC(reflect.ValueOf(Instructions[i].Execute).Pointer()).Name()
			str2 := runtime.FuncForPC(reflect.ValueOf(Instructions[i].AddressMode).Pointer()).Name()
			arr := [2]string{ str1[len(str1)-3:] , strings.Split(str2, ".")[2]}
			OpCodeMap[uint8(i)] = arr
		}
	}
}

// ADC https://www.masswerk.at/6502/6502_instruction_set.html#ADC
// Add Memory to Accumulator with Carry
// A + M + C -> A, C
func ADC(location uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	// Perform calculation in 16 bit
	temp := uint16(CPU.A) + uint16(data) + uint16(carry)
	tempSigned := int16(int8(CPU.A)) + int16(int8(data)) + int16(int8(carry))
	// Store last 8th bits to A register
	CPU.A = uint8(temp & 0x00FF)
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Check if result is greater than 255. If true we have a carry
	CPU.Set(FlagCarry, temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	CPU.Set(FlagOverflow, tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	CPU.PC += length
}

// AND https://www.masswerk.at/6502/6502_instruction_set.html#AND
// AND Memory with Accumulator
// A AND M -> A
func AND(location uint16, data uint8, length uint16) {
	// AND Memory with Accumulator
	temp := CPU.A & data
	// Store result in A register
	CPU.A = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Advance program counter
	CPU.PC += length
}

// ASL https://www.masswerk.at/6502/6502_instruction_set.html#ASL
// Shift Left One Bit (Memory or Accumulator)
// C <- [76543210] <- 0
func ASL(location uint16, data uint8, length uint16) {
	// Get carry from data
	carry := (data >> 7) & 0x01
	// Shift one bit left
	temp := data << 1
	// Set carry
	CPU.Set(FlagCarry, carry == 1)
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	opcode := CPU.Bus.CPURead(CPU.PC)
	inst := Instructions[opcode]
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		CPU.A = temp
	} else {
		CPU.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	CPU.PC += length
}

// BCC https://www.masswerk.at/6502/6502_instruction_set.html#BCC
// Branch on Carry Clear
// branch on C = 0
func BCC(location uint16, data uint8, length uint16) {
	// Check if C = 0
	if !CPU.GetFlag(FlagCarry) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// BCS https://www.masswerk.at/6502/6502_instruction_set.html#BCS
// Branch on Carry Set
// branch on C = 1
func BCS(location uint16, data uint8, length uint16) {
	// Check if C = 1
	if CPU.GetFlag(FlagCarry) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// BEQ https://www.masswerk.at/6502/6502_instruction_set.html#BEQ
// Branch on Result Zero
// branch on Z = 1
func BEQ(location uint16, data uint8, length uint16) {
	// Check if Z = 1
	if CPU.GetFlag(FlagZero) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// BIT https://www.masswerk.at/6502/6502_instruction_set.html#BIT
// Test Bits in Memory with Accumulator
// bits 7 and 6 of operand are transferred to bit 7 and 6 of P (N,V);
// the zero-flag is set to the result of operand AND accumulator.
// A AND M, M7 -> N, M6 -> V
func BIT(location uint16, data uint8, length uint16) {
	// AND Memory with Accumulator
	temp := CPU.A & data
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Transfer bit 7 to N
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	// Transfer bit 6 to V
	CPU.Set(FlagOverflow, (data>>6)&0x01 == 1)
	// Advance program counter
	CPU.PC += length
}

// BMI https://www.masswerk.at/6502/6502_instruction_set.html#BMI
// Branch on Result negative
// branch on N = 1
func BMI(location uint16, data uint8, length uint16) {
	// Check if N = 1
	if CPU.GetFlag(FlagNegative) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// BNE https://www.masswerk.at/6502/6502_instruction_set.html#BMI
// Branch on Result not Zero
// branch on Z = 0
func BNE(location uint16, data uint8, length uint16) {
	// Check if Z = 0
	if !CPU.GetFlag(FlagZero) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// BPL https://www.masswerk.at/6502/6502_instruction_set.html#BMI
// Branch on Result positive
// branch on N = 0
func BPL(location uint16, data uint8, length uint16) {
	// Check if N = 0
	if !CPU.GetFlag(FlagNegative) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
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
func BRK(location uint16, data uint8, length uint16) {
	// Get current pc + 2
	pc := CPU.PC + 2
	// Store high bytes of pc to stack
	CPU.Bus.CPUWrite(StackPage | uint16(CPU.S), uint8((pc>>8)&0x00FF))
	CPU.S--
	// Store low bytes of pc to stack
	CPU.Bus.CPUWrite(StackPage | uint16(CPU.S), uint8(pc&0x00FF))
	CPU.S--
	// Set flags and store current pc onto stack
	// From https://wiki.nesdev.org/w/index.php?title=Status_flags
	// In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
	CPU.Bus.CPUWrite(StackPage | uint16(CPU.S), CPU.P | FlagBreak | FlagInterruptDisable)
	CPU.S--
	// Set Interrupt disable flag
	// We don't want another interrupt inside the interrupt handler
	CPU.Set(FlagInterruptDisable, true)
	// Get pc from IRQ/BRK vector and jump to location
	low := uint16(CPU.Bus.CPURead(IRQVector))
	high := uint16(CPU.Bus.CPURead(IRQVector + 1))
	CPU.PC = (high << 8) | low
}

// BVC https://www.masswerk.at/6502/6502_instruction_set.html#BVC
// Branch on Overflow clear
// branch on V = 0
func BVC(location uint16, data uint8, length uint16) {
	// Check if V = 0
	if !CPU.GetFlag(FlagOverflow) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// BVS https://www.masswerk.at/6502/6502_instruction_set.html#BVS
// Branch on Overflow set
// branch on V = 1
func BVS(location uint16, data uint8, length uint16) {
	// Check if V = 1
	if CPU.GetFlag(FlagOverflow) {
		// Taking a branch takes one additional cycle
		CPU.CycleCount++
		if (CPU.PC + length) & 0xFF00 != location & 0xFF00 {
			// Branching to a different pages takes one additional cycle
			CPU.CycleCount++
		}
		// Jump to new location
		CPU.PC = location
		return
	}
	// Advance program counter
	CPU.PC += length
}

// CLC https://www.masswerk.at/6502/6502_instruction_set.html#CLC
// Clear carry flag
// 0 -> C
func CLC(location uint16, data uint8, length uint16) {
	// Clear Flag
	CPU.Set(FlagCarry, false)
	// Advance program counter
	CPU.PC += length
}

// CLD https://www.masswerk.at/6502/6502_instruction_set.html#CLD
// Clear decimal mode
// 0 -> D
func CLD(location uint16, data uint8, length uint16) {
	// Clear Flag
	CPU.Set(FlagDecimal, false)
	// Advance program counter
	CPU.PC += length
}

// CLI https://www.masswerk.at/6502/6502_instruction_set.html#CLI
// Clear interrupt disable bit
// 0 -> I
func CLI(location uint16, data uint8, length uint16) {
	// Clear Flag
	CPU.Set(FlagInterruptDisable, false)
	// Advance program counter
	CPU.PC += length
}

// CLV https://www.masswerk.at/6502/6502_instruction_set.html#CLV
// Clear overflow flag
// 0 -> V
func CLV(location uint16, data uint8, length uint16) {
	// Clear Flag
	CPU.Set(FlagOverflow, false)
	// Advance program counter
	CPU.PC += length
}

// CMP https://www.masswerk.at/6502/6502_instruction_set.html#CMP
// Compare Memory with Accumulator
// A - M
func CMP(location uint16, data uint8, length uint16) {
	// A - M
	temp := CPU.A - data
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	CPU.Set(FlagCarry, CPU.A >= data)
	// Advance program counter
	CPU.PC += length
}

// CPX https://www.masswerk.at/6502/6502_instruction_set.html#CPX
// Compare Memory with X
// X - M
func CPX(location uint16, data uint8, length uint16) {
	// X - M
	temp := CPU.X - data
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	CPU.Set(FlagCarry, CPU.X >= data)
	// Advance program counter
	CPU.PC += length
}

// CPY https://www.masswerk.at/6502/6502_instruction_set.html#CPY
// Compare Memory with Y
// Y - M
func CPY(location uint16, data uint8, length uint16) {
	// Y - M
	temp := CPU.Y - data
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	CPU.Set(FlagCarry, CPU.Y >= data)
	// Advance program counter
	CPU.PC += length
}

// DEC https://www.masswerk.at/6502/6502_instruction_set.html#DEC
// Decrement Memory by One
// M - 1 -> M
func DEC(location uint16, data uint8, length uint16) {
	// M - 1 -> M
	temp := CPU.Bus.CPURead(location) - 1
	CPU.Bus.CPUWrite(location, temp)
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, temp == 0)
	// Advance program counter
	CPU.PC += length
}

// DEX https://www.masswerk.at/6502/6502_instruction_set.html#DEX
// Decrement X by One
// X - 1 -> X
func DEX(location uint16, data uint8, length uint16) {
	// X - 1 -> X
	temp := CPU.X - 1
	CPU.X = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, temp == 0)
	// Advance program counter
	CPU.PC += length
}

// DEY https://www.masswerk.at/6502/6502_instruction_set.html#DEY
// Decrement Y by One
// Y - 1 -> Y
func DEY(location uint16, data uint8, length uint16) {
	// Y - 1 -> Y
	temp := CPU.Y - 1
	CPU.Y = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, temp == 0)
	// Advance program counter
	CPU.PC += length
}

// EOR https://www.masswerk.at/6502/6502_instruction_set.html#EOR
// Exclusive-OR Memory with Accumulator
// A ^ M -> A
func EOR(location uint16, data uint8, length uint16) {
	// A ^ M -> A
	temp := CPU.A ^ CPU.Bus.CPURead(location)
	CPU.A = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, temp == 0)
	// Advance program counter
	CPU.PC += length
}

// INC https://www.masswerk.at/6502/6502_instruction_set.html#INC
// Increment Memory by One
// M + 1 -> M
func INC(location uint16, data uint8, length uint16) {
	// M + 1 -> M
	temp := CPU.Bus.CPURead(location) + 1
	CPU.Bus.CPUWrite(location, temp)
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Advance program counter
	CPU.PC += length
}

// INX https://www.masswerk.at/6502/6502_instruction_set.html#INX
// Increment X by One
// X + 1 -> X
func INX(location uint16, data uint8, length uint16) {
	// X + 1 -> X
	temp := CPU.X + 1
	CPU.X = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Advance program counter
	CPU.PC += length
}

// INY https://www.masswerk.at/6502/6502_instruction_set.html#INY
// Increment Y by One
// Y + 1 -> Y
func INY(location uint16, data uint8, length uint16) {
	// Y + 1 -> Y
	temp := CPU.Y + 1
	CPU.Y = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Advance program counter
	CPU.PC += length
}

// JMP https://www.masswerk.at/6502/6502_instruction_set.html#JMP
// Jump to New Location
func JMP(location uint16, data uint8, length uint16) {
	 // Jump to New Location
	 CPU.PC = location
}

// JSR https://www.masswerk.at/6502/6502_instruction_set.html#JSR
// Jump to New Location
// push (PC+2)
func JSR(location uint16, data uint8, length uint16) {
	// Get PC+2
	pc := CPU.PC + 2
	// Store high bytes of pc+2 to stack
	CPU.Bus.CPUWrite(StackPage| uint16(CPU.S), uint8((pc>>8)&0x00FF))
	CPU.S--
	// Store low bytes of pc+2 to stack
	CPU.Bus.CPUWrite(StackPage| uint16(CPU.S), uint8(pc&0x00FF))
	CPU.S--
	// Jump to New Location
	CPU.PC = location
}

// LDA https://www.masswerk.at/6502/6502_instruction_set.html#LDA
// Load Accumulator with Memory
// M -> A
func LDA(location uint16, data uint8, length uint16) {
	// M -> A
	CPU.A = data
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, data == 0)
	// Advance program counter
	CPU.PC += length
}

// LDX https://www.masswerk.at/6502/6502_instruction_set.html#LDX
// Load X with Memory
// M -> X
func LDX(location uint16, data uint8, length uint16) {
	// M -> X
	CPU.X = data
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, data == 0)
	// Advance program counter
	CPU.PC += length
}

// LDY https://www.masswerk.at/6502/6502_instruction_set.html#LDY
// Load Y with Memory
// M -> Y
func LDY(location uint16, data uint8, length uint16) {
	// M -> Y
	CPU.Y = data
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, data == 0)
	// Advance program counter
	CPU.PC += length
}

func LSR(location uint16, data uint8, length uint16) {
	temp := data >> 1
	CPU.Set(FlagCarry, data&0x01 == 1)
	CPU.Set(FlagNegative, false)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	opcode := CPU.Bus.CPURead(CPU.PC)
	inst := Instructions[opcode]
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		CPU.A = temp
	} else {
		CPU.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	CPU.PC += length
}

// NOP https://www.masswerk.at/6502/6502_instruction_set.html#NOP
// No Operation
func NOP(location uint16, data uint8, length uint16) {
	// Advance program counter
	CPU.PC += length
	return
}

// ORA https://www.masswerk.at/6502/6502_instruction_set.html#ORA
// OR Memory with Accumulator
// A | M -> A
func ORA(location uint16, data uint8, length uint16) {
	// A | M -> A
	temp := CPU.A | CPU.Bus.CPURead(location)
	CPU.A = temp
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Advance program counter
	CPU.PC += length
}

// PHA https://www.masswerk.at/6502/6502_instruction_set.html#PHA
// Push Accumulator on Stack
// push A
func PHA(location uint16, data uint8, length uint16) {
	// Push A onto stack
	CPU.Bus.CPUWrite(StackPage| uint16(CPU.S), CPU.A)
	CPU.S--
	// Advance program counter
	CPU.PC += length
}

// PHP https://www.masswerk.at/6502/6502_instruction_set.html#PHP
// Push Processor Status on Stack
// The status register will be pushed with the break flag and bit 5 set to 1.
// push P
func PHP(location uint16, data uint8, length uint16) {
	// push P with bit 5 and 6 set to 1
	CPU.Bus.CPUWrite(StackPage| uint16(CPU.S), CPU.P | FlagBreak | FlagUnused)
	CPU.S--
	// Advance program counter
	CPU.PC += length
}

// PLA https://www.masswerk.at/6502/6502_instruction_set.html#PLA
// Pull Accumulator from Stack
// pull A
func PLA(location uint16, data uint8, length uint16) {
	// pull A
	CPU.S++
	CPU.A = CPU.Bus.CPURead(StackPage | uint16(CPU.S))
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (CPU.A >>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, CPU.A == 0)
	// Advance program counter
	CPU.PC += length
}

// PLP https://www.masswerk.at/6502/6502_instruction_set.html#PLP
// Pull Processor Status  from Stack
// The status register will be pulled with the break flag and bit 5 ignored
// pull P
func PLP(location uint16, data uint8, length uint16) {
	// pull p
	CPU.S++
	temp := CPU.Bus.CPURead(StackPage | uint16(CPU.S))
	// Ignore bit 4 and 5 from Stack but keep the value of bit 4 and 5 on the PC
	// Only bit 4 and 5 | Value from Stack without bit 4 and 5
	CPU.P = (CPU.P & (FlagBreak | FlagUnused)) | temp & ^(FlagBreak|FlagUnused)
	// Advance program counter
	CPU.PC += length
}

// ROL https://www.masswerk.at/6502/6502_instruction_set.html#ROL
// Rotate One Bit Left (Memory or Accumulator)
// C <- [76543210] <- C
func ROL(location uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	// C <- [76543210] <- C
	temp := data << 1 + carry
	CPU.Set(FlagCarry, (data & 0x80) == 0x80)
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp&0x80) == 0x80)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	opcode := CPU.Bus.CPURead(CPU.PC)
	inst := Instructions[opcode]
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		CPU.A = temp
	} else {
		CPU.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	CPU.PC += length
}

// ROR https://www.masswerk.at/6502/6502_instruction_set.html#ROR
// Rotate One Bit Right (Memory or Accumulator)
// C -> [76543210] -> C
func ROR(location uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 0x80
	}
	// C <- [76543210] <- C
	temp := data >> 1 + carry
	CPU.Set(FlagCarry, data & 0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagNegative, temp & 0x80 == 0x80)
	// Check if result is zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Check if we need to store the result in memory or in the A register
	opcode := CPU.Bus.CPURead(CPU.PC)
	inst := Instructions[opcode]
	if inst.ClockCycles == 2 {
		// Accumulator Addressing
		CPU.A = temp
	} else {
		CPU.Bus.CPUWrite(location, temp)
	}
	// Advance program counter
	CPU.PC += length
}

// RTI https://www.masswerk.at/6502/6502_instruction_set.html#RTI
// Return from Interrupt
// The status register is pulled with the break flag and bit 5 ignored. Then PC is pulled from the stack.
// pull P, pull PC
func RTI(location uint16, data uint8, length uint16) {
	// pull P from stack
	CPU.S++
	status := CPU.Bus.CPURead(0x0100 + uint16(CPU.S))
	// Ignore bit 4 and 5 from Stack but keep the value of bit 4 and 5 on the PC
	// Only bit 4 and 5 | Value from Stack without bit 4 and 5
	CPU.P = (CPU.P & (FlagBreak | FlagUnused)) | status & ^(FlagBreak|FlagUnused)
	// pull low bits of pc from stack
	CPU.S++
	low := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	// pull high bits of pc from stack
	CPU.S++
	high := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	// Set pc to pulled value
	pc := (high << 8) | low
	CPU.PC = pc
}

// RTS https://www.masswerk.at/6502/6502_instruction_set.html#RTS
// Return from Subroutine
// pull PC, PC+1 -> PC
func RTS(location uint16, data uint8, length uint16) {
	// pull low bits of pc from stack
	CPU.S++
	low := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	// pull high bits of pc from stack
	CPU.S++
	high := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	// Set pc to pulled value (PC+1)
	pc := (high << 8) | low
	// I don't know why???
	CPU.PC = pc + 1
}

// SBC https://www.masswerk.at/6502/6502_instruction_set.html#SBC
// Subtract Memory from Accumulator with Borrow
// A - M - ^C -> A
func SBC(location uint16, data uint8, length uint16) {
	// Get carry
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	// A - M - ^C = A + ^M + C
	dataInverse := ^data
	// Perform calculation in 16 bit
	temp := uint16(CPU.A) + uint16(dataInverse) + uint16(carry)
	tempSigned := int16(int8(CPU.A)) + int16(int8(dataInverse)) + int16(int8(carry))
	// Store last 8th bits to A register
	CPU.A = uint8(temp & 0x00FF)
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	// Check if last 8th bits are zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)
	// Check if result is greater thant 255. If true we have a carry
	CPU.Set(FlagCarry, temp > 255)
	// http://www.6502.org/tutorials/vflag.html
	// V indicates whether the result of an addition or subtraction is outside the range -128 to 127, i.e. whether there is a twos complement overflow
	CPU.Set(FlagOverflow, tempSigned < -128 || tempSigned > 127)
	// Advance program counter
	CPU.PC += length
}

// SEC https://www.masswerk.at/6502/6502_instruction_set.html#SEC
// Set Carry flag
// 1 -> C
func SEC(location uint16, data uint8, length uint16) {
	// Set Flag
	CPU.Set(FlagCarry, true)
	// Advance program counter
	CPU.PC += length
}

// SED https://www.masswerk.at/6502/6502_instruction_set.html#SED
// Set Decimal flag
// 1 -> D
func SED(location uint16, data uint8, length uint16) {
	// Set Flag
	CPU.Set(FlagDecimal, true)
	// Advance program counter
	CPU.PC += length
}

// SEI https://www.masswerk.at/6502/6502_instruction_set.html#SEI
// Set Interrupt Disable Status
// 1 -> I
func SEI(location uint16, data uint8, length uint16) {
	// Set Flag
	CPU.Set(FlagInterruptDisable, true)
	// Advance program counter
	CPU.PC += length
}

// STA https://www.masswerk.at/6502/6502_instruction_set.html#STA
// Store Accumulator in Memory
// A -> M
func STA(location uint16, data uint8, length uint16) {
	// A -> M
	CPU.Bus.CPUWrite(location, CPU.A)
	// Advance program counter
	CPU.PC += length
}

// STX https://www.masswerk.at/6502/6502_instruction_set.html#STX
// Store X in Memory
// X -> M
func STX(location uint16, data uint8, length uint16) {
	// X -> M
	CPU.Bus.CPUWrite(location, CPU.X)
	// Advance program counter
	CPU.PC += length
}

// STY https://www.masswerk.at/6502/6502_instruction_set.html#STY
// Store Y in Memory
// Y -> M
func STY(location uint16, data uint8, length uint16) {
	// Y -> M
	CPU.Bus.CPUWrite(location, CPU.Y)
	// Advance program counter
	CPU.PC += length
}

// TAX https://www.masswerk.at/6502/6502_instruction_set.html#TAX
// Transfer Accumulator to Index X
// A -> X
func TAX(location uint16, data uint8, length uint16) {
	// A -> X
	CPU.X = CPU.A
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (CPU.X>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, CPU.X == 0)
	// Advance program counter
	CPU.PC += length
}

// TAY https://www.masswerk.at/6502/6502_instruction_set.html#TAY
// Transfer Accumulator to Index Y
// A -> Y
func TAY(location uint16, data uint8, length uint16) {
	CPU.Y = CPU.A
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (CPU.Y>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, CPU.Y == 0)
	// Advance program counter
	CPU.PC += length
}

// TSX https://www.masswerk.at/6502/6502_instruction_set.html#TSX
// Transfer stack pointer to x
// S -> X
func TSX(location uint16, data uint8, length uint16) {
	// S -> X
	CPU.X = CPU.S
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (CPU.X>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, CPU.X == 0)
	// Advance program counter
	CPU.PC += length
}

// TXA https://www.masswerk.at/6502/6502_instruction_set.html#TXA
// Transfer X to A
// X -> A
func TXA(location uint16, data uint8, length uint16) {
	// X -> A
	val := CPU.X
	CPU.A = val
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, val == 0)
	// Advance program counter
	CPU.PC += length
}

// TXS https://www.masswerk.at/6502/6502_instruction_set.html#TXS
// Transfer X to stack pointer
// X -> S
func TXS(location uint16, data uint8, length uint16) {
	// X -> SP
	CPU.S = CPU.X
	// Advance program counter
	CPU.PC += length
}

// TYA https://www.masswerk.at/6502/6502_instruction_set.html#TYA
// Transfer Y to A
// Y -> A
func TYA(location uint16, data uint8, length uint16) {
	// Y -> A
	CPU.A = CPU.Y
	// Check if 8th bit is one
	CPU.Set(FlagNegative, (CPU.A>>7)&0x01 == 1)
	// Check if result is zero
	CPU.Set(FlagZero, CPU.A == 0)
	// Advance program counter
	CPU.PC += length
}
