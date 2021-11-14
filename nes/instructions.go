package nes

import (
	"log"
	"reflect"
	"runtime"
	"strings"
)

type AddressModeFunc func() (uint16, uint8)
type ExecuteFunc func(uint16, uint8)

type Instruction struct {
	AddressMode AddressModeFunc
	Execute     ExecuteFunc
	Length      int
	ClockCycles int
}

var Instructions [256]Instruction
var OpCodeMap map[uint8]string

func init() {
	Instructions[0x00] = Instruction{ImpliedAddressing, BRK, 1, 7}
	Instructions[0x01] = Instruction{IndexedIndirectAddressing, ORA, 2, 6}
	Instructions[0x02] = Instruction{}
	Instructions[0x03] = Instruction{}
	Instructions[0x04] = Instruction{}
	Instructions[0x05] = Instruction{ZeroPageAddressing, ORA, 2, 3}
	Instructions[0x06] = Instruction{ZeroPageAddressing, ASL, 2, 5}
	Instructions[0x07] = Instruction{}
	Instructions[0x08] = Instruction{ImpliedAddressing, PHP, 1, 3}
	Instructions[0x09] = Instruction{ImmediateAddress, ORA, 2, 2}
	Instructions[0x0A] = Instruction{AccumulatorAddressing, ASL, 1, 2}
	Instructions[0x0B] = Instruction{}
	Instructions[0x0C] = Instruction{}
	Instructions[0x0D] = Instruction{AbsoluteAddressing, ORA, 3, 4}
	Instructions[0x0E] = Instruction{AbsoluteAddressing, ASL, 3, 6}
	Instructions[0x0F] = Instruction{}

	Instructions[0x10] = Instruction{RelativeAddressing, BPL, 2, 2}
	Instructions[0x11] = Instruction{IndirectIndexedAddressing, ORA, 2, 5}
	Instructions[0x12] = Instruction{}
	Instructions[0x13] = Instruction{}
	Instructions[0x14] = Instruction{}
	Instructions[0x15] = Instruction{IndexedXZeroPageAddressing, ORA, 2, 4}
	Instructions[0x16] = Instruction{IndexedXZeroPageAddressing, ASL, 2, 6}
	Instructions[0x17] = Instruction{}
	Instructions[0x18] = Instruction{ImpliedAddressing, CLC, 1, 2}
	Instructions[0x19] = Instruction{IndexedYAbsoluteAddressing, ORA, 3, 4}
	Instructions[0x1A] = Instruction{}
	Instructions[0x1B] = Instruction{}
	Instructions[0x1C] = Instruction{}
	Instructions[0x1D] = Instruction{IndexedXAbsoluteAddressing, ORA, 3, 4}
	Instructions[0x1E] = Instruction{IndexedXAbsoluteAddressing, ASL, 3, 7}
	Instructions[0x1F] = Instruction{}

	Instructions[0x20] = Instruction{AbsoluteAddressing, JSR, 3, 6}
	Instructions[0x21] = Instruction{IndexedIndirectAddressing, AND, 2, 6}
	Instructions[0x22] = Instruction{}
	Instructions[0x23] = Instruction{}
	Instructions[0x24] = Instruction{ZeroPageAddressing, BIT, 2, 3}
	Instructions[0x25] = Instruction{ZeroPageAddressing, AND, 2, 3}
	Instructions[0x26] = Instruction{ZeroPageAddressing, ROL, 2, 5}
	Instructions[0x27] = Instruction{}
	Instructions[0x28] = Instruction{ImpliedAddressing, PLP, 1, 4}
	Instructions[0x29] = Instruction{ImmediateAddress, AND, 2, 2}
	Instructions[0x2A] = Instruction{AccumulatorAddressing, ROL, 1, 2}
	Instructions[0x2B] = Instruction{}
	Instructions[0x2C] = Instruction{AbsoluteAddressing, BIT, 3, 4}
	Instructions[0x2D] = Instruction{AbsoluteAddressing, AND, 3, 4}
	Instructions[0x2E] = Instruction{AbsoluteAddressing, ROL, 3, 6}
	Instructions[0x2F] = Instruction{}

	Instructions[0x30] = Instruction{RelativeAddressing, BMI, 2, 2}
	Instructions[0x31] = Instruction{IndirectIndexedAddressing, AND, 2, 5}
	Instructions[0x32] = Instruction{}
	Instructions[0x33] = Instruction{}
	Instructions[0x34] = Instruction{}
	Instructions[0x35] = Instruction{IndexedXZeroPageAddressing, AND, 2, 4}
	Instructions[0x36] = Instruction{IndexedXZeroPageAddressing, ROL, 2, 6}
	Instructions[0x37] = Instruction{}
	Instructions[0x38] = Instruction{ImpliedAddressing, SEC, 1, 2}
	Instructions[0x39] = Instruction{IndexedYAbsoluteAddressing, AND, 3, 4}
	Instructions[0x3A] = Instruction{}
	Instructions[0x3B] = Instruction{}
	Instructions[0x3C] = Instruction{}
	Instructions[0x3D] = Instruction{IndexedXAbsoluteAddressing, AND, 3, 4}
	Instructions[0x3E] = Instruction{IndexedXAbsoluteAddressing, ROL, 3, 7}
	Instructions[0x3F] = Instruction{}

	Instructions[0x40] = Instruction{ImpliedAddressing, RTI, 1, 6}
	Instructions[0x41] = Instruction{IndexedIndirectAddressing, EOR, 2, 6}
	Instructions[0x42] = Instruction{}
	Instructions[0x43] = Instruction{}
	Instructions[0x44] = Instruction{}
	Instructions[0x45] = Instruction{ZeroPageAddressing, EOR, 2, 3}
	Instructions[0x46] = Instruction{ZeroPageAddressing, LSR, 2, 5}
	Instructions[0x47] = Instruction{}
	Instructions[0x48] = Instruction{ImpliedAddressing, PHA, 1, 3}
	Instructions[0x49] = Instruction{ImmediateAddress, EOR, 2, 2}
	Instructions[0x4A] = Instruction{AccumulatorAddressing, LSR, 1, 2}
	Instructions[0x4B] = Instruction{}
	Instructions[0x4C] = Instruction{AbsoluteAddressing, JMP, 3, 3}
	Instructions[0x4D] = Instruction{AbsoluteAddressing, EOR, 3, 4}
	Instructions[0x4E] = Instruction{AbsoluteAddressing, LSR, 3, 6}
	Instructions[0x4F] = Instruction{}

	Instructions[0x50] = Instruction{RelativeAddressing, BVC, 2, 2}
	Instructions[0x51] = Instruction{IndirectIndexedAddressing, EOR, 2, 5}
	Instructions[0x52] = Instruction{}
	Instructions[0x53] = Instruction{}
	Instructions[0x54] = Instruction{}
	Instructions[0x55] = Instruction{IndexedXZeroPageAddressing, EOR, 2, 4}
	Instructions[0x56] = Instruction{IndexedXZeroPageAddressing, LSR, 2, 6}
	Instructions[0x57] = Instruction{}
	Instructions[0x58] = Instruction{ImpliedAddressing, CLI, 1, 2}
	Instructions[0x59] = Instruction{IndexedXAbsoluteAddressing, EOR, 3, 4}
	Instructions[0x5A] = Instruction{}
	Instructions[0x5B] = Instruction{}
	Instructions[0x5C] = Instruction{}
	Instructions[0x5D] = Instruction{IndexedXAbsoluteAddressing, EOR, 3, 4}
	Instructions[0x5E] = Instruction{IndexedXAbsoluteAddressing, LSR, 3, 7}
	Instructions[0x5F] = Instruction{}

	Instructions[0x60] = Instruction{ImpliedAddressing, RTS, 1, 6}
	Instructions[0x61] = Instruction{IndexedIndirectAddressing, ADC, 2, 6}
	Instructions[0x62] = Instruction{}
	Instructions[0x63] = Instruction{}
	Instructions[0x64] = Instruction{}
	Instructions[0x65] = Instruction{ZeroPageAddressing, ADC, 2, 3}
	Instructions[0x66] = Instruction{ZeroPageAddressing, ROR, 2, 5}
	Instructions[0x67] = Instruction{}
	Instructions[0x68] = Instruction{ImpliedAddressing, PLA, 1, 4}
	Instructions[0x69] = Instruction{ImmediateAddress, ADC, 2, 2}
	Instructions[0x6A] = Instruction{AccumulatorAddressing, ROR, 1, 2}
	Instructions[0x6B] = Instruction{}
	Instructions[0x6C] = Instruction{AbsoluteIndirect, JMP, 3, 5}
	Instructions[0x6D] = Instruction{AbsoluteAddressing, ADC, 3, 4}
	Instructions[0x6E] = Instruction{AbsoluteAddressing, ROR, 3, 6}
	Instructions[0x6F] = Instruction{}

	Instructions[0x70] = Instruction{RelativeAddressing, BVS, 2, 2}
	Instructions[0x71] = Instruction{IndirectIndexedAddressing, ADC, 2, 5}
	Instructions[0x72] = Instruction{}
	Instructions[0x73] = Instruction{}
	Instructions[0x74] = Instruction{}
	Instructions[0x75] = Instruction{IndexedXZeroPageAddressing, ADC, 2, 4}
	Instructions[0x76] = Instruction{IndexedXZeroPageAddressing, ROR, 2, 6}
	Instructions[0x77] = Instruction{}
	Instructions[0x78] = Instruction{ImpliedAddressing, SEI, 1, 2}
	Instructions[0x79] = Instruction{IndexedYAbsoluteAddressing, ADC, 3, 4}
	Instructions[0x7A] = Instruction{}
	Instructions[0x7B] = Instruction{}
	Instructions[0x7C] = Instruction{}
	Instructions[0x7D] = Instruction{IndexedXAbsoluteAddressing, ADC, 3, 4}
	Instructions[0x7E] = Instruction{IndexedXAbsoluteAddressing, ROR, 3, 7}
	Instructions[0x7F] = Instruction{}

	Instructions[0x80] = Instruction{}
	Instructions[0x81] = Instruction{IndexedIndirectAddressing, STA, 2, 6}
	Instructions[0x82] = Instruction{}
	Instructions[0x83] = Instruction{}
	Instructions[0x84] = Instruction{ZeroPageAddressing, STY, 2, 3}
	Instructions[0x85] = Instruction{ZeroPageAddressing, STA, 2, 3}
	Instructions[0x86] = Instruction{ZeroPageAddressing, STX, 2, 3}
	Instructions[0x87] = Instruction{}
	Instructions[0x88] = Instruction{ImpliedAddressing, DEY, 1, 2}
	Instructions[0x89] = Instruction{}
	Instructions[0x8A] = Instruction{ImpliedAddressing, TXA, 1, 2}
	Instructions[0x8B] = Instruction{}
	Instructions[0x8C] = Instruction{AbsoluteAddressing, STY, 3, 4}
	Instructions[0x8D] = Instruction{AbsoluteAddressing, STA, 3, 4}
	Instructions[0x8E] = Instruction{AbsoluteAddressing, STX, 3, 4}
	Instructions[0x8F] = Instruction{}

	Instructions[0x90] = Instruction{RelativeAddressing, BCC, 2, 2}
	Instructions[0x91] = Instruction{IndirectIndexedAddressing, STA, 2, 6}
	Instructions[0x92] = Instruction{}
	Instructions[0x93] = Instruction{}
	Instructions[0x94] = Instruction{IndexedXZeroPageAddressing, STY, 2, 4}
	Instructions[0x95] = Instruction{IndexedXZeroPageAddressing, STA, 2, 4}
	Instructions[0x96] = Instruction{IndexedXZeroPageAddressing, STX, 2, 4}
	Instructions[0x97] = Instruction{}
	Instructions[0x98] = Instruction{ImpliedAddressing, TYA, 1, 2}
	Instructions[0x99] = Instruction{IndexedYAbsoluteAddressing, STA, 3, 5}
	Instructions[0x9A] = Instruction{ImpliedAddressing, TXS, 1, 2}
	Instructions[0x9B] = Instruction{}
	Instructions[0x9C] = Instruction{}
	Instructions[0x9D] = Instruction{IndexedXAbsoluteAddressing, STA, 3, 5}
	Instructions[0x9E] = Instruction{}
	Instructions[0x9F] = Instruction{}

	Instructions[0xA0] = Instruction{ImmediateAddress, LDY, 2, 2}
	Instructions[0xA1] = Instruction{IndexedIndirectAddressing, LDA, 2, 6}
	Instructions[0xA2] = Instruction{ImmediateAddress, LDX, 2, 2}
	Instructions[0xA3] = Instruction{}
	Instructions[0xA4] = Instruction{ZeroPageAddressing, LDY, 2, 3}
	Instructions[0xA5] = Instruction{ZeroPageAddressing, LDA, 2, 3}
	Instructions[0xA6] = Instruction{ZeroPageAddressing, LDX, 2, 3}
	Instructions[0xA7] = Instruction{}
	Instructions[0xA8] = Instruction{ImpliedAddressing, TAY, 1, 2}
	Instructions[0xA9] = Instruction{ImmediateAddress, LDA, 2, 2}
	Instructions[0xAA] = Instruction{ImpliedAddressing, TAX, 1, 2}
	Instructions[0xAB] = Instruction{}
	Instructions[0xAC] = Instruction{AbsoluteAddressing, LDY, 3, 4}
	Instructions[0xAD] = Instruction{AbsoluteAddressing, LDA, 3, 4}
	Instructions[0xAE] = Instruction{AbsoluteAddressing, LDX, 3, 4}
	Instructions[0xAF] = Instruction{}

	Instructions[0xB0] = Instruction{RelativeAddressing, BCS, 2, 2}
	Instructions[0xB1] = Instruction{IndirectIndexedAddressing, LDA, 2, 5}
	Instructions[0xB2] = Instruction{}
	Instructions[0xB3] = Instruction{}
	Instructions[0xB4] = Instruction{IndexedXZeroPageAddressing, LDY, 2, 4}
	Instructions[0xB5] = Instruction{IndexedXZeroPageAddressing, LDA, 2, 4}
	Instructions[0xB6] = Instruction{IndexedYZeroPageAddressing, LDX, 2, 4}
	Instructions[0xB7] = Instruction{}
	Instructions[0xB8] = Instruction{ImpliedAddressing, CLV, 1, 2}
	Instructions[0xB9] = Instruction{IndexedYAbsoluteAddressing, LDA, 3, 4}
	Instructions[0xBA] = Instruction{ImpliedAddressing, TSX, 1, 2}
	Instructions[0xBB] = Instruction{}
	Instructions[0xBC] = Instruction{IndexedXAbsoluteAddressing, LDY, 3, 4}
	Instructions[0xBD] = Instruction{IndexedXAbsoluteAddressing, LDA, 3, 4}
	Instructions[0xBE] = Instruction{IndexedXAbsoluteAddressing, LDX, 3, 4}
	Instructions[0xBF] = Instruction{}

	Instructions[0xC0] = Instruction{ImmediateAddress, CPY, 2, 2}
	Instructions[0xC1] = Instruction{IndexedIndirectAddressing, CMP, 2, 6}
	Instructions[0xC2] = Instruction{}
	Instructions[0xC3] = Instruction{}
	Instructions[0xC4] = Instruction{ZeroPageAddressing, CPY, 2, 3}
	Instructions[0xC5] = Instruction{ZeroPageAddressing, CMP, 2, 3}
	Instructions[0xC6] = Instruction{ZeroPageAddressing, DEC, 2, 5}
	Instructions[0xC7] = Instruction{}
	Instructions[0xC8] = Instruction{ImpliedAddressing, INY, 1, 2}
	Instructions[0xC9] = Instruction{ImmediateAddress, CMP, 2, 3}
	Instructions[0xCA] = Instruction{ImpliedAddressing, DEX, 1, 2}
	Instructions[0xCB] = Instruction{}
	Instructions[0xCC] = Instruction{AbsoluteAddressing, CPY, 3, 4}
	Instructions[0xCD] = Instruction{AbsoluteAddressing, CMP, 3, 4}
	Instructions[0xCE] = Instruction{AbsoluteAddressing, DEC, 3, 6}
	Instructions[0xCF] = Instruction{}

	Instructions[0xD0] = Instruction{RelativeAddressing, BNE, 2, 2}
	Instructions[0xD1] = Instruction{IndirectIndexedAddressing, CMP, 2, 5}
	Instructions[0xD2] = Instruction{}
	Instructions[0xD3] = Instruction{}
	Instructions[0xD4] = Instruction{}
	Instructions[0xD5] = Instruction{IndexedXZeroPageAddressing, CMP, 2, 4}
	Instructions[0xD6] = Instruction{IndexedXZeroPageAddressing, DEC, 2, 6}
	Instructions[0xD7] = Instruction{}
	Instructions[0xD8] = Instruction{ImpliedAddressing, CLD, 1, 2}
	Instructions[0xD9] = Instruction{IndexedYAbsoluteAddressing, CMP, 3, 4}
	Instructions[0xDA] = Instruction{}
	Instructions[0xDB] = Instruction{}
	Instructions[0xDC] = Instruction{}
	Instructions[0xDD] = Instruction{IndexedXAbsoluteAddressing, CMP, 3, 4} // TODO: CMP oder CMD oder CMP??
	Instructions[0xDE] = Instruction{IndexedXAbsoluteAddressing, DEC, 3, 7}
	Instructions[0xDF] = Instruction{}

	Instructions[0xE0] = Instruction{ImmediateAddress, CPX, 2, 2}
	Instructions[0xE1] = Instruction{IndexedIndirectAddressing, SBC, 2, 6}
	Instructions[0xE2] = Instruction{}
	Instructions[0xE3] = Instruction{}
	Instructions[0xE4] = Instruction{ZeroPageAddressing, CPX, 2, 3}
	Instructions[0xE5] = Instruction{ZeroPageAddressing, SBC, 2, 3}
	Instructions[0xE6] = Instruction{ZeroPageAddressing, INC, 2, 5}
	Instructions[0xE7] = Instruction{}
	Instructions[0xE8] = Instruction{ImpliedAddressing, INX, 1, 2}
	Instructions[0xE9] = Instruction{ImmediateAddress, SBC, 2, 2}
	Instructions[0xEA] = Instruction{ImpliedAddressing, NOP, 1, 2}
	Instructions[0xEB] = Instruction{}
	Instructions[0xEC] = Instruction{AbsoluteAddressing, CPX, 3, 4}
	Instructions[0xED] = Instruction{AbsoluteAddressing, SBC, 3, 4}
	Instructions[0xEE] = Instruction{AbsoluteAddressing, INC, 3, 6}
	Instructions[0xEF] = Instruction{}

	Instructions[0xF0] = Instruction{RelativeAddressing, BEQ, 2, 2}
	Instructions[0xF1] = Instruction{IndirectIndexedAddressing, SBC, 2, 5}
	Instructions[0xF2] = Instruction{}
	Instructions[0xF3] = Instruction{}
	Instructions[0xF4] = Instruction{}
	Instructions[0xF5] = Instruction{IndexedXZeroPageAddressing, SBC, 2, 4}
	Instructions[0xF6] = Instruction{IndexedXZeroPageAddressing, INC, 2, 6}
	Instructions[0xF7] = Instruction{}
	Instructions[0xF8] = Instruction{ImpliedAddressing, SED, 1, 2}
	Instructions[0xF9] = Instruction{IndexedYAbsoluteAddressing, SBC, 3, 4}
	Instructions[0xFA] = Instruction{}
	Instructions[0xFB] = Instruction{}
	Instructions[0xFC] = Instruction{}
	Instructions[0xFD] = Instruction{IndexedXAbsoluteAddressing, SBC, 3, 4}
	Instructions[0xFE] = Instruction{IndexedXAbsoluteAddressing, INC, 3, 7}
	Instructions[0xFF] = Instruction{}

	OpCodeMap = make(map[uint8]string, 0xFF)
	for i := range Instructions {
		if Instructions[i].Length == 0 {
			OpCodeMap[uint8(i)] = "ERR"
		} else {
			str1 := runtime.FuncForPC(reflect.ValueOf(Instructions[i].Execute).Pointer()).Name()
						str2 := runtime.FuncForPC(reflect.ValueOf(Instructions[i].AddressMode).Pointer()).Name()
			OpCodeMap[uint8(i)] = str1[len(str1)-3:] + " (" + strings.Split(str2, ".")[2] + ") "
		}
	}
}

func ADC(location uint16, data uint8) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	temp := CPU.A + data + carry
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagCarry, (data>>7)&0x01 == 1)
	CPU.Set(FlagOverflow, int(temp) < -128 || int(temp) > 127)
}

func AND(location uint16, data uint8) {
	temp := CPU.A & CPU.Bus.CPURead(location)
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func ASL(location uint16, data uint8) {
	newCarry := (data >> 7) & 0x01
	newData := data << 1
	CPU.Set(FlagCarry, newCarry == 1)
	CPU.Set(FlagNegative, (newData>>7)&0x01 == 1)
	CPU.Set(FlagZero, newData == 0)
}

func BCC(location uint16, data uint8) {
	if !CPU.GetFlag(FlagCarry) {
		CPU.PC = location
	}
}

func BCS(location uint16, data uint8) {
	if CPU.GetFlag(FlagCarry) {
		CPU.PC = location
	}
}

func BEQ(location uint16, data uint8) {
	if CPU.GetFlag(FlagZero) {
		CPU.PC = location
	}
}

func BIT(location uint16, data uint8) {
	CPU.Set(FlagZero, (CPU.A&data) == 1)
	CPU.Set(FlagOverflow, (data>>6)&0x01 == 1)
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
}

func BMI(location uint16, data uint8) {
	if CPU.GetFlag(FlagNegative) {
		CPU.PC = location
	}
}

func BNE(location uint16, data uint8) {
	if !CPU.GetFlag(FlagZero) {
		CPU.PC = location
	}
}

func BPL(location uint16, data uint8) {
	if !CPU.GetFlag(FlagNegative) {
		CPU.PC = location
	}
}

func BRK(location uint16, data uint8) {
	pc := CPU.PC
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8((pc>>8)&0x00FF))
	CPU.S--
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8(pc&0x00FF))
	CPU.S--
	CPU.Set(FlagBreak, true)
	CPU.Set(FlagUnused, true)
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), CPU.P)
	CPU.S--
	CPU.PC = location
	CPU.Set(FlagInterruptDisable, true)
	low := uint16(CPU.Bus.CPURead(0xFFFE))
	high := uint16(CPU.Bus.CPURead(0xFFFF))
	CPU.PC = (high << 8) | low
}

func BVC(location uint16, data uint8) {
	if !CPU.GetFlag(FlagOverflow) {
		CPU.PC = location
	}
}

func BVS(location uint16, data uint8) {
	if CPU.GetFlag(FlagOverflow) {
		CPU.PC = location
	}
}

func CLC(location uint16, data uint8) {
	CPU.Set(FlagCarry, false)
}

func CLD(location uint16, data uint8) {
	CPU.Set(FlagDecimal, false)
}

func CLI(location uint16, data uint8) {
	CPU.Set(FlagInterruptDisable, false)
}

func CLV(location uint16, data uint8) {
	CPU.Set(FlagOverflow, false)
}

func CMP(location uint16, data uint8) {
	temp := CPU.A - data
	log.Println(CPU.A, data)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	// From Wiki: After SBC or CMP, this flag will be set if no borrow was the result, or alternatively a "greater than or equal" result.
	CPU.Set(FlagCarry, CPU.A >= data)
}

func CPX(location uint16, data uint8) {
	temp := CPU.X - data
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagCarry, CPU.GetFlag(FlagZero) && !CPU.GetFlag(FlagNegative))
}

func CPY(location uint16, data uint8) {
	temp := CPU.Y - data
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagCarry, CPU.GetFlag(FlagZero) && !CPU.GetFlag(FlagNegative))
}

func DEC(location uint16, data uint8) {
	temp := CPU.Bus.CPURead(location) - 1
	CPU.Bus.CPUWrite(location, temp)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func DEX(location uint16, data uint8) {
	temp := CPU.X - 1
	CPU.X = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func DEY(location uint16, data uint8) {
	temp := CPU.Y - 1
	CPU.Y = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func EOR(location uint16, data uint8) {
	temp := CPU.A ^ CPU.Bus.CPURead(location)
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func INC(location uint16, data uint8) {
	temp := CPU.Bus.CPURead(location) + 1
	CPU.Bus.CPUWrite(location, temp)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func INX(location uint16, data uint8) {
	temp := CPU.X + 1
	CPU.X = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func INY(location uint16, data uint8) {
	temp := CPU.Y + 1
	CPU.Y = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func JMP(location uint16, data uint8) {
	CPU.PC = location - uint16(CPU.CurrentInstruction.Length)
}

func JSR(location uint16, data uint8) {
	pc := CPU.PC + uint16(CPU.CurrentInstruction.Length)
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8((pc>>8)&0x00FF))
	CPU.S--
	CPU.Bus.CPUWrite(0x0100+uint16(CPU.S), uint8(pc&0x00FF))
	CPU.S--
	CPU.PC = location - uint16(CPU.CurrentInstruction.Length)
}

func LDA(location uint16, data uint8) {
	CPU.A = data
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.Set(FlagZero, data == 0)
}

func LDX(location uint16, data uint8) {
	CPU.X = data
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.Set(FlagZero, data == 0)
}

func LDY(location uint16, data uint8) {
	CPU.Y = data
	CPU.Set(FlagNegative, (data>>7)&0x01 == 1)
	CPU.Set(FlagZero, data == 0)
}

func LSR(location uint16, data uint8) {
	temp := uint16(data) << 7
	CPU.Bus.CPUWrite(location, uint8(temp>>8))
	CPU.Set(FlagCarry, uint8(temp>>7)&0x01 == 1)
	CPU.Set(FlagNegative, false)
	CPU.Set(FlagZero, temp == 0)
}

func NOP(location uint16, data uint8) {
	return
}

func ORA(location uint16, data uint8) {

	temp := CPU.A | CPU.Bus.CPURead(location)
	CPU.A = temp
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func PHA(location uint16, data uint8) {
	temp := CPU.A
	CPU.Bus.CPUWrite(0x0100|uint16(CPU.S), temp)
	CPU.S--
}

func PHP(location uint16, data uint8) {
	temp := CPU.P
	temp = temp & 0b00110000
	CPU.Bus.CPUWrite(0x0100|uint16(CPU.S), temp)
	CPU.S--
}

func PLA(location uint16, data uint8) {
	temp := CPU.Bus.CPURead(0x0100 | uint16(CPU.S))
	CPU.A = temp
	CPU.S++
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func PLP(location uint16, data uint8) {
	temp := CPU.Bus.CPURead(0x0100 | uint16(CPU.S))
	temp = temp | 0b11001111
	CPU.S++
}

func ROL(location uint16, data uint8) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}

	newCarry := (data >> 7) & 0x01
	newData := (data << 1) & carry
	CPU.Set(FlagCarry, newCarry == 1)
	CPU.Set(FlagNegative, (newData>>7)&0x01 == 1)
	CPU.Set(FlagZero, newData == 0)
}

func ROR(location uint16, data uint8) {
	carry := uint16(0x0000)
	if CPU.GetFlag(FlagCarry) {
		carry = uint16(0x8000)
	}

	temp := uint16(data)<<7 | carry
	CPU.Bus.CPUWrite(location, uint8(temp>>8))
	CPU.Set(FlagCarry, uint8(temp>>7)&0x01 == 1)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
}

func RTI(location uint16, data uint8) {
	status := CPU.Bus.CPURead(0x0100 + uint16(CPU.S))
	CPU.S++
	status = status & 0b11001111
	CPU.P = status
	low := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	CPU.S++
	high := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	CPU.S++
	pc := (high << 8) | low
	CPU.PC = pc
}

func RTS(location uint16, data uint8) {
	CPU.S++
	low := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	CPU.S++
	high := uint16(CPU.Bus.CPURead(0x0100 + uint16(CPU.S)))
	pc := (high << 8) | low
	CPU.PC = pc - uint16(CPU.CurrentInstruction.Length)
}

func SBC(location uint16, data uint8) {
	var carry uint8
	if CPU.GetFlag(FlagCarry) {
		carry = 1
	}
	carry = ^carry
	temp := CPU.A - data - carry
	CPU.A = temp
	CPU.Set(FlagCarry, uint8(temp>>7)&0x01 == 1)
	CPU.Set(FlagNegative, (temp>>7)&0x01 == 1)
	CPU.Set(FlagZero, temp == 0)
	CPU.Set(FlagOverflow, int(temp) < -128 || int(temp) > 127)
}

func SEC(location uint16, data uint8) {
	CPU.Set(FlagCarry, true)
}

func SED(location uint16, data uint8) {
	CPU.Set(FlagDecimal, true)
}

func SEI(location uint16, data uint8) {
	CPU.Set(FlagInterruptDisable, true)
}

func STA(location uint16, data uint8) {
	CPU.Bus.CPUWrite(location, CPU.A)
}

func STX(location uint16, data uint8) {
	CPU.Bus.CPUWrite(location, CPU.X)
}

func STY(location uint16, data uint8) {
	CPU.Bus.CPUWrite(location, CPU.Y)
}

func TAX(location uint16, data uint8) {
	val := CPU.A
	CPU.X = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
}

func TAY(location uint16, data uint8) {
	val := CPU.A
	CPU.Y = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
}

func TSX(location uint16, data uint8) {
	val := CPU.S
	CPU.X = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
}

func TXA(location uint16, data uint8) {
	val := CPU.X
	CPU.A = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
}

func TXS(location uint16, data uint8) {
	val := CPU.X
	CPU.S = val
}

func TYA(location uint16, data uint8) {
	val := CPU.Y
	CPU.A = val
	CPU.Set(FlagNegative, (val>>7)&0x01 == 1)
	CPU.Set(FlagZero, val == 0)
}
