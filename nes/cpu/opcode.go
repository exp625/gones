package cpu

import (
	"reflect"
	"runtime"
	"strings"
)

func init() {
	// OpCode Matrix from http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf
	Instructions[0x00] = Instruction{IMP, BRK, 1, 7}
	Instructions[0x01] = Instruction{IDX, ORA, 2, 6}
	Instructions[0x02] = Instruction{} // JAM
	Instructions[0x03] = Instruction{IDX, SLO, 2, 8}
	Instructions[0x04] = Instruction{ZP0, ILLNOP, 2, 3}
	Instructions[0x05] = Instruction{ZP0, ORA, 2, 3}
	Instructions[0x06] = Instruction{ZP0, ASL, 2, 5}
	Instructions[0x07] = Instruction{ZP0, SLO, 2, 5}
	Instructions[0x08] = Instruction{IMP, PHP, 1, 3}
	Instructions[0x09] = Instruction{IMM, ORA, 2, 2}
	Instructions[0x0A] = Instruction{ACC, ASL, 1, 2}
	Instructions[0x0B] = Instruction{IMM, ANC, 2, 2}
	Instructions[0x0C] = Instruction{ABS, ILLNOP, 3, 4}
	Instructions[0x0D] = Instruction{ABS, ORA, 3, 4}
	Instructions[0x0E] = Instruction{ABS, ASL, 3, 6}
	Instructions[0x0F] = Instruction{ABS, SLO, 3, 6}

	Instructions[0x10] = Instruction{REL, BPL, 2, 2}
	Instructions[0x11] = Instruction{IZY, ORA, 2, 5}
	Instructions[0x12] = Instruction{} // JAM
	Instructions[0x13] = Instruction{IZY, SLO, 2, 8}
	Instructions[0x14] = Instruction{ZPX, ILLNOP, 2, 4}
	Instructions[0x15] = Instruction{ZPX, ORA, 2, 4}
	Instructions[0x16] = Instruction{ZPX, ASL, 2, 6}
	Instructions[0x17] = Instruction{ZPX, SLO, 2, 6}
	Instructions[0x18] = Instruction{IMP, CLC, 1, 2}
	Instructions[0x19] = Instruction{ABY, ORA, 3, 4}
	Instructions[0x1A] = Instruction{IMP, ILLNOP, 1, 2}
	Instructions[0x1B] = Instruction{ABY, SLO, 3, 7}
	Instructions[0x1C] = Instruction{ABX, ILLNOP, 3, 4}
	Instructions[0x1D] = Instruction{ABX, ORA, 3, 4}
	Instructions[0x1E] = Instruction{ABX, ASL, 3, 7}
	Instructions[0x1F] = Instruction{ABX, SLO, 3, 7}

	Instructions[0x20] = Instruction{ABS, JSR, 3, 6}
	Instructions[0x21] = Instruction{IDX, AND, 2, 6}
	Instructions[0x22] = Instruction{} // JAM
	Instructions[0x23] = Instruction{IDX, RLA, 2, 8}
	Instructions[0x24] = Instruction{ZP0, BIT, 2, 3}
	Instructions[0x25] = Instruction{ZP0, AND, 2, 3}
	Instructions[0x26] = Instruction{ZP0, ROL, 2, 5}
	Instructions[0x27] = Instruction{ZP0, RLA, 2, 5}
	Instructions[0x28] = Instruction{IMP, PLP, 1, 4}
	Instructions[0x29] = Instruction{IMM, AND, 2, 2}
	Instructions[0x2A] = Instruction{ACC, ROL, 1, 2}
	Instructions[0x2B] = Instruction{IMM, ANC, 2, 2}
	Instructions[0x2C] = Instruction{ABS, BIT, 3, 4}
	Instructions[0x2D] = Instruction{ABS, AND, 3, 4}
	Instructions[0x2E] = Instruction{ABS, ROL, 3, 6}
	Instructions[0x2F] = Instruction{ABS, RLA, 3, 6}

	Instructions[0x30] = Instruction{REL, BMI, 2, 2}
	Instructions[0x31] = Instruction{IZY, AND, 2, 5}
	Instructions[0x32] = Instruction{} // JAM
	Instructions[0x33] = Instruction{IZY, RLA, 2, 8}
	Instructions[0x34] = Instruction{ZPX, ILLNOP, 2, 4}
	Instructions[0x35] = Instruction{ZPX, AND, 2, 4}
	Instructions[0x36] = Instruction{ZPX, ROL, 2, 6}
	Instructions[0x37] = Instruction{ZPX, RLA, 2, 6}
	Instructions[0x38] = Instruction{IMP, SEC, 1, 2}
	Instructions[0x39] = Instruction{ABY, AND, 3, 4}
	Instructions[0x3A] = Instruction{IMP, ILLNOP, 1, 2}
	Instructions[0x3B] = Instruction{ABY, RLA, 3, 7}
	Instructions[0x3C] = Instruction{ABX, ILLNOP, 3, 4}
	Instructions[0x3D] = Instruction{ABX, AND, 3, 4}
	Instructions[0x3E] = Instruction{ABX, ROL, 3, 7}
	Instructions[0x3F] = Instruction{ABX, RLA, 3, 7}

	Instructions[0x40] = Instruction{IMP, RTI, 1, 6}
	Instructions[0x41] = Instruction{IDX, EOR, 2, 6}
	Instructions[0x42] = Instruction{} // JAM
	Instructions[0x43] = Instruction{IDX, SRE, 2, 8}
	Instructions[0x44] = Instruction{ZP0, ILLNOP, 2, 3}
	Instructions[0x45] = Instruction{ZP0, EOR, 2, 3}
	Instructions[0x46] = Instruction{ZP0, LSR, 2, 5}
	Instructions[0x47] = Instruction{ZP0, SRE, 2, 5}
	Instructions[0x48] = Instruction{IMP, PHA, 1, 3}
	Instructions[0x49] = Instruction{IMM, EOR, 2, 2}
	Instructions[0x4A] = Instruction{ACC, LSR, 1, 2}
	Instructions[0x4B] = Instruction{IMM, ALR, 2, 2}
	Instructions[0x4C] = Instruction{ABS, JMP, 3, 3}
	Instructions[0x4D] = Instruction{ABS, EOR, 3, 4}
	Instructions[0x4E] = Instruction{ABS, LSR, 3, 6}
	Instructions[0x4F] = Instruction{ABS, SRE, 3, 6}

	Instructions[0x50] = Instruction{REL, BVC, 2, 2}
	Instructions[0x51] = Instruction{IZY, EOR, 2, 5}
	Instructions[0x52] = Instruction{} // JAM
	Instructions[0x53] = Instruction{IZY, SRE, 2, 8}
	Instructions[0x54] = Instruction{ZPX, ILLNOP, 2, 4}
	Instructions[0x55] = Instruction{ZPX, EOR, 2, 4}
	Instructions[0x56] = Instruction{ZPX, LSR, 2, 6}
	Instructions[0x57] = Instruction{ZPX, SRE, 2, 6}
	Instructions[0x58] = Instruction{IMP, CLI, 1, 2}
	Instructions[0x59] = Instruction{ABY, EOR, 3, 4}
	Instructions[0x5A] = Instruction{IMP, ILLNOP, 1, 2}
	Instructions[0x5B] = Instruction{ABY, SRE, 3, 7}
	Instructions[0x5C] = Instruction{ABX, ILLNOP, 3, 4}
	Instructions[0x5D] = Instruction{ABX, EOR, 3, 4}
	Instructions[0x5E] = Instruction{ABX, LSR, 3, 7}
	Instructions[0x5F] = Instruction{ABX, SRE, 3, 7}

	Instructions[0x60] = Instruction{IMP, RTS, 1, 6}
	Instructions[0x61] = Instruction{IDX, ADC, 2, 6}
	Instructions[0x62] = Instruction{} // JAM
	Instructions[0x63] = Instruction{IDX, RRA, 2, 8}
	Instructions[0x64] = Instruction{ZP0, ILLNOP, 2, 3}
	Instructions[0x65] = Instruction{ZP0, ADC, 2, 3}
	Instructions[0x66] = Instruction{ZP0, ROR, 2, 5}
	Instructions[0x67] = Instruction{ZP0, RRA, 2, 5}
	Instructions[0x68] = Instruction{IMP, PLA, 1, 4}
	Instructions[0x69] = Instruction{IMM, ADC, 2, 2}
	Instructions[0x6A] = Instruction{ACC, ROR, 1, 2}
	Instructions[0x6B] = Instruction{IMM, ARR, 2, 2}
	Instructions[0x6C] = Instruction{IND, JMP, 3, 5}
	Instructions[0x6D] = Instruction{ABS, ADC, 3, 4}
	Instructions[0x6E] = Instruction{ABS, ROR, 3, 6}
	Instructions[0x6F] = Instruction{ABS, RRA, 3, 6}

	Instructions[0x70] = Instruction{REL, BVS, 2, 2}
	Instructions[0x71] = Instruction{IZY, ADC, 2, 5}
	Instructions[0x72] = Instruction{} // JAM
	Instructions[0x73] = Instruction{IZY, RRA, 2, 8}
	Instructions[0x74] = Instruction{ZPX, ILLNOP, 2, 4}
	Instructions[0x75] = Instruction{ZPX, ADC, 2, 4}
	Instructions[0x76] = Instruction{ZPX, ROR, 2, 6}
	Instructions[0x77] = Instruction{ZPX, RRA, 2, 6}
	Instructions[0x78] = Instruction{IMP, SEI, 1, 2}
	Instructions[0x79] = Instruction{ABY, ADC, 3, 4}
	Instructions[0x7A] = Instruction{IMP, ILLNOP, 1, 2}
	Instructions[0x7B] = Instruction{ABY, RRA, 3, 7}
	Instructions[0x7C] = Instruction{ABX, ILLNOP, 3, 4}
	Instructions[0x7D] = Instruction{ABX, ADC, 3, 4}
	Instructions[0x7E] = Instruction{ABX, ROR, 3, 7}
	Instructions[0x7F] = Instruction{ABX, RRA, 3, 7}

	Instructions[0x80] = Instruction{IMM, ILLNOP, 2, 2}
	Instructions[0x81] = Instruction{IDX, STA, 2, 6}
	Instructions[0x82] = Instruction{IMM, ILLNOP, 2, 2}
	Instructions[0x83] = Instruction{IDX, SAX, 2, 6}
	Instructions[0x84] = Instruction{ZP0, STY, 2, 3}
	Instructions[0x85] = Instruction{ZP0, STA, 2, 3}
	Instructions[0x86] = Instruction{ZP0, STX, 2, 3}
	Instructions[0x87] = Instruction{ZP0, SAX, 2, 3}
	Instructions[0x88] = Instruction{IMP, DEY, 1, 2}
	Instructions[0x89] = Instruction{IMM, ILLNOP, 2, 2}
	Instructions[0x8A] = Instruction{IMP, TXA, 1, 2}
	Instructions[0x8B] = Instruction{IMM, ANE, 2, 2} // Highly unstable
	Instructions[0x8C] = Instruction{ABS, STY, 3, 4}
	Instructions[0x8D] = Instruction{ABS, STA, 3, 4}
	Instructions[0x8E] = Instruction{ABS, STX, 3, 4}
	Instructions[0x8F] = Instruction{ABS, SAX, 3, 4}

	Instructions[0x90] = Instruction{REL, BCC, 2, 2}
	Instructions[0x91] = Instruction{IZY, STA, 2, 6}
	Instructions[0x92] = Instruction{} // JAM
	Instructions[0x93] = Instruction{IZY, SHA, 2, 6}
	Instructions[0x94] = Instruction{ZPX, STY, 2, 4}
	Instructions[0x95] = Instruction{ZPX, STA, 2, 4}
	Instructions[0x96] = Instruction{ZPY, STX, 2, 4}
	Instructions[0x97] = Instruction{ZPY, SAX, 2, 4}
	Instructions[0x98] = Instruction{IMP, TYA, 1, 2}
	Instructions[0x99] = Instruction{ABY, STA, 3, 5}
	Instructions[0x9A] = Instruction{IMP, TXS, 1, 2}
	Instructions[0x9B] = Instruction{ABY, TAS, 3, 5}
	Instructions[0x9C] = Instruction{ABX, SHY, 3, 5}
	Instructions[0x9D] = Instruction{ABX, STA, 3, 5}
	Instructions[0x9E] = Instruction{ABY, SHX, 3, 5}
	Instructions[0x9F] = Instruction{ABY, SHA, 3, 5}

	Instructions[0xA0] = Instruction{IMM, LDY, 2, 2}
	Instructions[0xA1] = Instruction{IDX, LDA, 2, 6}
	Instructions[0xA2] = Instruction{IMM, LDX, 2, 2}
	Instructions[0xA3] = Instruction{IDX, LAX, 2, 6}
	Instructions[0xA4] = Instruction{ZP0, LDY, 2, 3}
	Instructions[0xA5] = Instruction{ZP0, LDA, 2, 3}
	Instructions[0xA6] = Instruction{ZP0, LDX, 2, 3}
	Instructions[0xA7] = Instruction{ZP0, LAX, 2, 3}
	Instructions[0xA8] = Instruction{IMP, TAY, 1, 2}
	Instructions[0xA9] = Instruction{IMM, LDA, 2, 2}
	Instructions[0xAA] = Instruction{IMP, TAX, 1, 2}
	Instructions[0xAB] = Instruction{IMM, LXA, 2, 2} // Highly unstable
	Instructions[0xAC] = Instruction{ABS, LDY, 3, 4}
	Instructions[0xAD] = Instruction{ABS, LDA, 3, 4}
	Instructions[0xAE] = Instruction{ABS, LDX, 3, 4}
	Instructions[0xAF] = Instruction{ABS, LAX, 3, 4}

	Instructions[0xB0] = Instruction{REL, BCS, 2, 2}
	Instructions[0xB1] = Instruction{IZY, LDA, 2, 5}
	Instructions[0xB2] = Instruction{} // JAM
	Instructions[0xB3] = Instruction{IZY, LAX, 2, 5}
	Instructions[0xB4] = Instruction{ZPX, LDY, 2, 4}
	Instructions[0xB5] = Instruction{ZPX, LDA, 2, 4}
	Instructions[0xB6] = Instruction{ZPY, LDX, 2, 4}
	Instructions[0xB7] = Instruction{ZPY, LAX, 2, 4}
	Instructions[0xB8] = Instruction{IMP, CLV, 1, 2}
	Instructions[0xB9] = Instruction{ABY, LDA, 3, 4}
	Instructions[0xBA] = Instruction{IMP, TSX, 1, 2}
	Instructions[0xBB] = Instruction{ABY, LAS, 3, 4}
	Instructions[0xBC] = Instruction{ABX, LDY, 3, 4}
	Instructions[0xBD] = Instruction{ABX, LDA, 3, 4}
	Instructions[0xBE] = Instruction{ABY, LDX, 3, 4}
	Instructions[0xBF] = Instruction{ABY, LAX, 3, 4}

	Instructions[0xC0] = Instruction{IMM, CPY, 2, 2}
	Instructions[0xC1] = Instruction{IDX, CMP, 2, 6}
	Instructions[0xC2] = Instruction{IMM, ILLNOP, 2, 2}
	Instructions[0xC3] = Instruction{IDX, DCP, 2, 8}
	Instructions[0xC4] = Instruction{ZP0, CPY, 2, 3}
	Instructions[0xC5] = Instruction{ZP0, CMP, 2, 3}
	Instructions[0xC6] = Instruction{ZP0, DEC, 2, 5}
	Instructions[0xC7] = Instruction{ZP0, DCP, 2, 5}
	Instructions[0xC8] = Instruction{IMP, INY, 1, 2}
	Instructions[0xC9] = Instruction{IMM, CMP, 2, 2}
	Instructions[0xCA] = Instruction{IMP, DEX, 1, 2}
	Instructions[0xCB] = Instruction{IMM, SBX, 2, 2}
	Instructions[0xCC] = Instruction{ABS, CPY, 3, 4}
	Instructions[0xCD] = Instruction{ABS, CMP, 3, 4}
	Instructions[0xCE] = Instruction{ABS, DEC, 3, 6}
	Instructions[0xCF] = Instruction{ABS, DCP, 3, 6}

	Instructions[0xD0] = Instruction{REL, BNE, 2, 2}
	Instructions[0xD1] = Instruction{IZY, CMP, 2, 5}
	Instructions[0xD2] = Instruction{} // JAM
	Instructions[0xD3] = Instruction{IZY, DCP, 2, 8}
	Instructions[0xD4] = Instruction{ZPX, ILLNOP, 2, 4}
	Instructions[0xD5] = Instruction{ZPX, CMP, 2, 4}
	Instructions[0xD6] = Instruction{ZPX, DEC, 2, 6}
	Instructions[0xD7] = Instruction{ZPX, DCP, 2, 6}
	Instructions[0xD8] = Instruction{IMP, CLD, 1, 2}
	Instructions[0xD9] = Instruction{ABY, CMP, 3, 4}
	Instructions[0xDA] = Instruction{IMP, ILLNOP, 1, 2}
	Instructions[0xDB] = Instruction{ABY, DCP, 3, 7}
	Instructions[0xDC] = Instruction{ABX, ILLNOP, 3, 4}
	Instructions[0xDD] = Instruction{ABX, CMP, 3, 4}
	Instructions[0xDE] = Instruction{ABX, DEC, 3, 7}
	Instructions[0xDF] = Instruction{ABX, DCP, 3, 7}

	Instructions[0xE0] = Instruction{IMM, CPX, 2, 2}
	Instructions[0xE1] = Instruction{IDX, SBC, 2, 6}
	Instructions[0xE2] = Instruction{IMM, ILLNOP, 2, 2}
	Instructions[0xE3] = Instruction{IDX, ISC, 2, 8}
	Instructions[0xE4] = Instruction{ZP0, CPX, 2, 3}
	Instructions[0xE5] = Instruction{ZP0, SBC, 2, 3}
	Instructions[0xE6] = Instruction{ZP0, INC, 2, 5}
	Instructions[0xE7] = Instruction{ZP0, ISC, 2, 5}
	Instructions[0xE8] = Instruction{IMP, INX, 1, 2}
	Instructions[0xE9] = Instruction{IMM, SBC, 2, 2}
	Instructions[0xEA] = Instruction{IMP, NOP, 1, 2}
	Instructions[0xEB] = Instruction{IMM, USBC, 2, 2}
	Instructions[0xEC] = Instruction{ABS, CPX, 3, 4}
	Instructions[0xED] = Instruction{ABS, SBC, 3, 4}
	Instructions[0xEE] = Instruction{ABS, INC, 3, 6}
	Instructions[0xEF] = Instruction{ABS, ISC, 3, 6}

	Instructions[0xF0] = Instruction{REL, BEQ, 2, 2}
	Instructions[0xF1] = Instruction{IZY, SBC, 2, 5}
	Instructions[0xF2] = Instruction{} // JAM
	Instructions[0xF3] = Instruction{IZY, ISC, 2, 8}
	Instructions[0xF4] = Instruction{ZPX, ILLNOP, 2, 4}
	Instructions[0xF5] = Instruction{ZPX, SBC, 2, 4}
	Instructions[0xF6] = Instruction{ZPX, INC, 2, 6}
	Instructions[0xF7] = Instruction{ZPX, ISC, 2, 6}
	Instructions[0xF8] = Instruction{IMP, SED, 1, 2}
	Instructions[0xF9] = Instruction{ABY, SBC, 3, 4}
	Instructions[0xFA] = Instruction{IMP, ILLNOP, 1, 2}
	Instructions[0xFB] = Instruction{ABY, ISC, 3, 7}
	Instructions[0xFC] = Instruction{ABX, ILLNOP, 3, 4}
	Instructions[0xFD] = Instruction{ABX, SBC, 3, 4}
	Instructions[0xFE] = Instruction{ABX, INC, 3, 7}
	Instructions[0xFF] = Instruction{ABX, ISC, 3, 7}

	// Generate debug map
	OpCodeMap = make(map[uint8][2]string, 0xFF)
	for i := range Instructions {
		if Instructions[i].Length == 0 {
			arr := [2]string{"ERR", "ERR"}
			OpCodeMap[uint8(i)] = arr
		} else {
			str1 := runtime.FuncForPC(reflect.ValueOf(Instructions[i].Execute).Pointer()).Name()
			str2 := runtime.FuncForPC(reflect.ValueOf(Instructions[i].AddressMode).Pointer()).Name()
			arr := [2]string{strings.Split(str1, ".")[2], strings.Split(str2, ".")[2]}
			OpCodeMap[uint8(i)] = arr
		}
	}
}
