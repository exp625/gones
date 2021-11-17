package nes

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
	Instructions[0x59] = Instruction{ABX, EOR, 3, 4}
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
	Instructions[0x96] = Instruction{ZPX, STX, 2, 4}
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
	Instructions[0xBE] = Instruction{ABX, LDX, 3, 4}
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
	Instructions[0xDD] = Instruction{ABX, CMP, 3, 4} // TODO: CMP oder CMD oder CMP??
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

// TODO: Branch instruction add 1 cycle if a page is crossed

func ADC(location uint16, data uint8, length uint16) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	// Perform calculation in 16 bit
	temp := uint16(CPU.A) + uint16(data) + uint16(carry)

	// Store last 8 bits to A register
	CPU.A = uint8(temp & 0x00FF)

	// Bit 8 is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)

	// Check if last 8 bits are zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)

	// Check if result is greater thant 255. If true we have a carry
	CPU.Set(FlagCarry, temp > 255)

	// Positive Number + Positive Number = Negative Result -> Overflow
	// Negative Number + Negative Number = Positive Result -> Overflow
	// Positive Number + Negative Number = Either Result -> Cannot Overflow
	// Positive Number + Positive Number = Positive Result -> OK! No Overflow
	// Negative Number + Negative Number = Negative Result -> OK! NO Overflow
	// so V = ~(A^M) & (A^R)
	CPU.Set(FlagOverflow, ^(uint16(CPU.A)^uint16(data)) & (uint16(CPU.A)^temp) == 0)
	CPU.PC += length
}

func AND(location uint16, data uint8, length uint16) {
	temp := CPU.A & data
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func ASL(location uint16, data uint8, length uint16) {
	newCarry := (data >> 7) & 0x01
	newData := data << 1
	CPU.Set(FlagCarry, newCarry == 1)
	CPU.Set(FlagNegative, (newData>>7)&0x01 == 1)
	CPU.Set(FlagZero, newData == 0)
	CPU.PC += length
}

func BCC(location uint16, data uint8, length uint16) {
	if !CPU.GetFlag(FlagCarry) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BCS(location uint16, data uint8, length uint16) {
	if CPU.GetFlag(FlagCarry) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BEQ(location uint16, data uint8, length uint16) {
	if CPU.GetFlag(FlagZero) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BIT(location uint16, data uint8, length uint16) {
	temp := CPU.A & data

	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagOverflow, (data>>6)&0x01 == 1)
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.PC += length
}

func BMI(location uint16, data uint8, length uint16) {
	if CPU.GetFlag(FlagNegative) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BNE(location uint16, data uint8, length uint16) {
	if !CPU.GetFlag(FlagZero) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BPL(location uint16, data uint8, length uint16) {
	if !CPU.GetFlag(FlagNegative) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BRK(location uint16, data uint8, length uint16) {
	pc := CPU.PC + length
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8((pc>>8)&0x00FF))
	CPU.S--
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8(pc&0x00FF))
	CPU.S--
	CPU.Set(FlagBreak, true)
	CPU.Set(FlagUnused, true)
	CPU.Set(FlagInterruptDisable, true)
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), CPU.P)
	CPU.S--
	CPU.Set(FlagBreak, false)

	low := uint16(CPU.Bus.CPURead(0xFFFE))
	high := uint16(CPU.Bus.CPURead(0xFFFF))
	CPU.PC = (high << 8) | low
}

func BVC(location uint16, data uint8, length uint16) {
	if !CPU.GetFlag(FlagOverflow) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func BVS(location uint16, data uint8, length uint16) {
	if CPU.GetFlag(FlagOverflow) {
		CPU.PC = location
		CPU.CycleCount++
		return
	}
	CPU.PC += length
}

func CLC(location uint16, data uint8, length uint16) {
	CPU.Set(FlagCarry, false)
	CPU.PC += length
}

func CLD(location uint16, data uint8, length uint16) {
	CPU.Set(FlagDecimal, false)
	CPU.PC += length
}

func CLI(location uint16, data uint8, length uint16) {
	CPU.Set(FlagInterruptDisable, false)
	CPU.PC += length
}

func CLV(location uint16, data uint8, length uint16) {
	CPU.Set(FlagOverflow, false)
	CPU.PC += length
}

func CMP(location uint16, data uint8, length uint16) {
	temp := CPU.A - data
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp & 0x00FF == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	CPU.Set(FlagCarry, CPU.A >= data)
	CPU.PC += length
}

func CPX(location uint16, data uint8, length uint16) {
	temp := CPU.X - data
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagCarry, CPU.GetFlag(FlagZero) && !CPU.GetFlag(FlagNegative))
	CPU.PC += length
}

func CPY(location uint16, data uint8, length uint16) {
	temp := CPU.Y - data
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagCarry, CPU.GetFlag(FlagZero) && !CPU.GetFlag(FlagNegative))
	CPU.PC += length
}

func DEC(location uint16, data uint8, length uint16) {
	temp := CPU.Bus.CPURead(location) - 1
	CPU.Bus.CPUWrite(location, temp)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func DEX(location uint16, data uint8, length uint16) {
	temp := CPU.X - 1
	CPU.X = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func DEY(location uint16, data uint8, length uint16) {
	temp := CPU.Y - 1
	CPU.Y = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func EOR(location uint16, data uint8, length uint16) {
	temp := CPU.A ^ CPU.Bus.CPURead(location)
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func INC(location uint16, data uint8, length uint16) {
	temp := CPU.Bus.CPURead(location) + 1
	CPU.Bus.CPUWrite(location, temp)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func INX(location uint16, data uint8, length uint16) {
	temp := CPU.X + 1
	CPU.X = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func INY(location uint16, data uint8, length uint16) {
	temp := CPU.Y + 1
	CPU.Y = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func JMP(location uint16, data uint8, length uint16) {
	 CPU.PC = location
}

func JSR(location uint16, data uint8, length uint16) {
	pc := CPU.PC + length
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8((pc>>8)&0x00FF))
	CPU.S--
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8(pc&0x00FF))
	CPU.S--
	CPU.PC = location
}

func LDA(location uint16, data uint8, length uint16) {
	CPU.A = data
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.Set(FlagZero, data == 0)
	CPU.PC += length
}

func LDX(location uint16, data uint8, length uint16) {
	CPU.X = data
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.Set(FlagZero, data == 0)
	CPU.PC += length
}

func LDY(location uint16, data uint8, length uint16) {
	CPU.Y = data
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.Set(FlagZero, data == 0)
	CPU.PC += length
}

func LSR(location uint16, data uint8, length uint16) {
	temp := data >> 1
	CPU.Set(FlagCarry, data&0x01 == 1)
	CPU.Set(FlagNegative, false)
	CPU.Set(FlagZero, temp == 0)
	if CPU.CurrentInstruction.ClockCycles == 2 {
		// Accumulator Addressing
		CPU.A = temp
	} else {
		CPU.Bus.CPUWrite(location, temp)
	}
	CPU.PC += length
}

func NOP(location uint16, data uint8, length uint16) {
	CPU.PC += length
	return
}

func ORA(location uint16, data uint8, length uint16) {

	temp := CPU.A | CPU.Bus.CPURead(location)
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func PHA(location uint16, data uint8, length uint16) {
	temp := CPU.A
	CPU.Bus.CPUWrite(ZeroPage | uint16(CPU.S), temp)
	CPU.S--
	CPU.PC += length
}

func PHP(location uint16, data uint8, length uint16) {
	temp := CPU.P
	temp = temp | FlagBreak | FlagUnused
	CPU.Bus.CPUWrite(ZeroPage | uint16(CPU.S), temp)
	CPU.S--
	CPU.PC += length
}

func PLA(location uint16, data uint8, length uint16) {
	CPU.S++
	temp := CPU.Bus.CPURead(ZeroPage | uint16(CPU.S))
	CPU.A = temp

	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length
}

func PLP(location uint16, data uint8, length uint16) {
	CPU.S++
	temp := CPU.Bus.CPURead(ZeroPage | uint16(CPU.S))
	// Ignore bit 4 and 5 from Stack but keep the value of bit 4 and 5 on the PC
	// Only bit 4 and 5 | Value from Stack without bit 4 and 5
	CPU.P = (CPU.P & (FlagBreak | FlagUnused)) | temp & ^(FlagBreak | FlagUnused)
	CPU.PC += length

}

func ROL(location uint16, data uint8, length uint16) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}

	temp := data << 1 + carry

	CPU.Set(FlagCarry, (data & 0x80) == 0x80)
	CPU.Set(FlagNegative, (temp&0x80) == 0x80)
	CPU.Set(FlagZero, temp == 0)
	CPU.PC += length

}

func ROR(location uint16, data uint8, length uint16) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 0x80
	}

	temp := data >> 1 + carry
	CPU.Set(FlagCarry, data & 0x01 == 1)
	CPU.Set(FlagNegative, temp & 0x80 == 0x80)
	CPU.Set(FlagZero, temp == 0)

	if CPU.CurrentInstruction.ClockCycles == 2 {
		// Accumulator Addressing
		CPU.A = temp
	} else {
		CPU.Bus.CPUWrite(location, temp)
	}
	CPU.PC += length

}

func RTI(location uint16, data uint8, length uint16) {
	CPU.S++
	status := CPU.Bus.CPURead(0x0100 + uint16(CPU.S))

	status = status & 0b11001111
	CPU.P = status
	CPU.S++
	low := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	CPU.S++
	high := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))

	pc := (high << 8) | low
	CPU.PC = pc
}

func RTS(location uint16, data uint8, length uint16) {
	CPU.S++
	low := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	CPU.S++
	high := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	pc := (high << 8) | low
	CPU.PC = pc
}

func SBC(location uint16, data uint8, length uint16) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	dataInverse := ^data

	// Perform calculation in 16 bit
	temp := uint16(CPU.A) + uint16(dataInverse) + uint16(carry)

	// Store last 8 bits to A register
	CPU.A = uint8(temp & 0x00FF)

	// Bit 8 is one
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)

	// Check if last 8 bits are zero
	CPU.Set(FlagZero, (temp & 0x00FF) == 0)

	// Check if result is greater thant 255. If true we have a carry
	CPU.Set(FlagCarry, temp > 255)

	// Positive Number + Positive Number = Negative Result -> Overflow
	// Negative Number + Negative Number = Positive Result -> Overflow
	// Positive Number + Negative Number = Either Result -> Cannot Overflow
	// Positive Number + Positive Number = Positive Result -> OK! No Overflow
	// Negative Number + Negative Number = Negative Result -> OK! NO Overflow
	// so V = ~(A^M) & (A^R)
	CPU.Set(FlagOverflow, ^(uint16(CPU.A)^uint16(dataInverse)) & (uint16(CPU.A)^temp) == 0)
	CPU.PC += length
}

func SEC(location uint16, data uint8, length uint16) {
	CPU.Set(FlagCarry, true)
	CPU.PC += length
}

func SED(location uint16, data uint8, length uint16) {
	CPU.Set(FlagDecimal, true)
	CPU.PC += length
}

func SEI(location uint16, data uint8, length uint16) {
	CPU.Set(FlagInterruptDisable, true)
	CPU.PC += length
}

func STA(location uint16, data uint8, length uint16) {
	CPU.Bus.CPUWrite(location, CPU.A)
	CPU.PC += length
}

func STX(location uint16, data uint8, length uint16) {
	CPU.Bus.CPUWrite(location, CPU.X)
	CPU.PC += length
}

func STY(location uint16, data uint8, length uint16) {
	CPU.Bus.CPUWrite(location, CPU.Y)
	CPU.PC += length
}

func TAX(location uint16, data uint8, length uint16) {
	val := CPU.A
	CPU.X = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
	CPU.PC += length
}

func TAY(location uint16, data uint8, length uint16) {
	val := CPU.A
	CPU.Y = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
	CPU.PC += length
}

func TSX(location uint16, data uint8, length uint16) {
	val := CPU.S
	CPU.X = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
	CPU.PC += length
}

func TXA(location uint16, data uint8, length uint16) {
	val := CPU.X
	CPU.A = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
	CPU.PC += length
}

func TXS(location uint16, data uint8, length uint16) {
	val := CPU.X
	CPU.S = val
	CPU.PC += length
}

func TYA(location uint16, data uint8, length uint16) {
	val := CPU.Y
	CPU.A = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
	CPU.PC += length
}
