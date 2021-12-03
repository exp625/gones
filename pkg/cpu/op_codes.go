package cpu

// Op Code Matrix taken from http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf
func (cpu *CPU) generateInstructions() {
	cpu.Instructions = [256]Instruction{}
	cpu.Mnemonics = make(map[uint8][2]string, 256)

	cpu.Instructions[0x00] = Instruction{"IMP", "BRK", cpu.IMP, cpu.BRK, 1, 7, true}
	cpu.Instructions[0x01] = Instruction{"IDX", "ORA", cpu.IDX, cpu.ORA, 2, 6, true}
	cpu.Instructions[0x02] = Instruction{}
	cpu.Instructions[0x03] = Instruction{"IDX", "SLO", cpu.IDX, cpu.SLO, 2, 8, false}
	cpu.Instructions[0x04] = Instruction{"ZP0", "NOP", cpu.ZP0, cpu.ILLNOP, 2, 3, false}
	cpu.Instructions[0x05] = Instruction{"ZP0", "ORA", cpu.ZP0, cpu.ORA, 2, 3, true}
	cpu.Instructions[0x06] = Instruction{"ZP0", "ASL", cpu.ZP0, cpu.ASL, 2, 5, true}
	cpu.Instructions[0x07] = Instruction{"ZP0", "SLO", cpu.ZP0, cpu.SLO, 2, 5, false}
	cpu.Instructions[0x08] = Instruction{"IMP", "PHP", cpu.IMP, cpu.PHP, 1, 3, true}
	cpu.Instructions[0x09] = Instruction{"IMM", "ORA", cpu.IMM, cpu.ORA, 2, 2, true}
	cpu.Instructions[0x0A] = Instruction{"ACC", "ASL", cpu.ACC, cpu.ASL, 1, 2, true}
	cpu.Instructions[0x0B] = Instruction{"IMM", "ANC", cpu.IMM, cpu.ANC, 2, 2, false}
	cpu.Instructions[0x0C] = Instruction{"ABS", "NOP", cpu.ABS, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0x0D] = Instruction{"ABS", "ORA", cpu.ABS, cpu.ORA, 3, 4, true}
	cpu.Instructions[0x0E] = Instruction{"ABS", "ASL", cpu.ABS, cpu.ASL, 3, 6, true}
	cpu.Instructions[0x0F] = Instruction{"ABS", "SLO", cpu.ABS, cpu.SLO, 3, 6, false}

	cpu.Instructions[0x10] = Instruction{"REL", "BPL", cpu.REL, cpu.BPL, 2, 2, true}
	cpu.Instructions[0x11] = Instruction{"IZY", "ORA", cpu.IZY, cpu.ORA, 2, 5, true}
	cpu.Instructions[0x12] = Instruction{}
	cpu.Instructions[0x13] = Instruction{"IZY", "SLO", cpu.IZY, cpu.SLO, 2, 8, false}
	cpu.Instructions[0x14] = Instruction{"ZPX", "NOP", cpu.ZPX, cpu.ILLNOP, 2, 4, false}
	cpu.Instructions[0x15] = Instruction{"ZPX", "ORA", cpu.ZPX, cpu.ORA, 2, 4, true}
	cpu.Instructions[0x16] = Instruction{"ZPX", "ASL", cpu.ZPX, cpu.ASL, 2, 6, true}
	cpu.Instructions[0x17] = Instruction{"ZPX", "SLO", cpu.ZPX, cpu.SLO, 2, 6, false}
	cpu.Instructions[0x18] = Instruction{"IMP", "CLC", cpu.IMP, cpu.CLC, 1, 2, true}
	cpu.Instructions[0x19] = Instruction{"ABY", "ORA", cpu.ABY, cpu.ORA, 3, 4, true}
	cpu.Instructions[0x1A] = Instruction{"IMP", "NOP", cpu.IMP, cpu.ILLNOP, 1, 2, false}
	cpu.Instructions[0x1B] = Instruction{"ABY", "SLO", cpu.ABY, cpu.SLO, 3, 7, false}
	cpu.Instructions[0x1C] = Instruction{"ABX", "NOP", cpu.ABX, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0x1D] = Instruction{"ABX", "ORA", cpu.ABX, cpu.ORA, 3, 4, true}
	cpu.Instructions[0x1E] = Instruction{"ABX", "ASL", cpu.ABX, cpu.ASL, 3, 7, true}
	cpu.Instructions[0x1F] = Instruction{"ABX", "SLO", cpu.ABX, cpu.SLO, 3, 7, false}

	cpu.Instructions[0x20] = Instruction{"ABS", "JSR", cpu.ABS, cpu.JSR, 3, 6, true}
	cpu.Instructions[0x21] = Instruction{"IDX", "AND", cpu.IDX, cpu.AND, 2, 6, true}
	cpu.Instructions[0x22] = Instruction{}
	cpu.Instructions[0x23] = Instruction{"IDX", "RLA", cpu.IDX, cpu.RLA, 2, 8, false}
	cpu.Instructions[0x24] = Instruction{"ZP0", "BIT", cpu.ZP0, cpu.BIT, 2, 3, true}
	cpu.Instructions[0x25] = Instruction{"ZP0", "AND", cpu.ZP0, cpu.AND, 2, 3, true}
	cpu.Instructions[0x26] = Instruction{"ZP0", "ROL", cpu.ZP0, cpu.ROL, 2, 5, true}
	cpu.Instructions[0x27] = Instruction{"ZP0", "RLA", cpu.ZP0, cpu.RLA, 2, 5, false}
	cpu.Instructions[0x28] = Instruction{"IMP", "PLP", cpu.IMP, cpu.PLP, 1, 4, true}
	cpu.Instructions[0x29] = Instruction{"IMM", "AND", cpu.IMM, cpu.AND, 2, 2, true}
	cpu.Instructions[0x2A] = Instruction{"ACC", "ROL", cpu.ACC, cpu.ROL, 1, 2, true}
	cpu.Instructions[0x2B] = Instruction{"IMM", "ANC", cpu.IMM, cpu.ANC, 2, 2, false}
	cpu.Instructions[0x2C] = Instruction{"ABS", "BIT", cpu.ABS, cpu.BIT, 3, 4, true}
	cpu.Instructions[0x2D] = Instruction{"ABS", "AND", cpu.ABS, cpu.AND, 3, 4, true}
	cpu.Instructions[0x2E] = Instruction{"ABS", "ROL", cpu.ABS, cpu.ROL, 3, 6, true}
	cpu.Instructions[0x2F] = Instruction{"ABS", "RLA", cpu.ABS, cpu.RLA, 3, 6, false}

	cpu.Instructions[0x30] = Instruction{"REL", "BMI", cpu.REL, cpu.BMI, 2, 2, true}
	cpu.Instructions[0x31] = Instruction{"IZY", "AND", cpu.IZY, cpu.AND, 2, 5, true}
	cpu.Instructions[0x32] = Instruction{}
	cpu.Instructions[0x33] = Instruction{"IZY", "RLA", cpu.IZY, cpu.RLA, 2, 8, false}
	cpu.Instructions[0x34] = Instruction{"ZPX", "NOP", cpu.ZPX, cpu.ILLNOP, 2, 4, false}
	cpu.Instructions[0x35] = Instruction{"ZPX", "AND", cpu.ZPX, cpu.AND, 2, 4, true}
	cpu.Instructions[0x36] = Instruction{"ZPX", "ROL", cpu.ZPX, cpu.ROL, 2, 6, true}
	cpu.Instructions[0x37] = Instruction{"ZPX", "RLA", cpu.ZPX, cpu.RLA, 2, 6, false}
	cpu.Instructions[0x38] = Instruction{"IMP", "SEC", cpu.IMP, cpu.SEC, 1, 2, true}
	cpu.Instructions[0x39] = Instruction{"ABY", "AND", cpu.ABY, cpu.AND, 3, 4, true}
	cpu.Instructions[0x3A] = Instruction{"IMP", "NOP", cpu.IMP, cpu.ILLNOP, 1, 2, false}
	cpu.Instructions[0x3B] = Instruction{"ABY", "RLA", cpu.ABY, cpu.RLA, 3, 7, false}
	cpu.Instructions[0x3C] = Instruction{"ABX", "NOP", cpu.ABX, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0x3D] = Instruction{"ABX", "AND", cpu.ABX, cpu.AND, 3, 4, true}
	cpu.Instructions[0x3E] = Instruction{"ABX", "ROL", cpu.ABX, cpu.ROL, 3, 7, true}
	cpu.Instructions[0x3F] = Instruction{"ABX", "RLA", cpu.ABX, cpu.RLA, 3, 7, false}

	cpu.Instructions[0x40] = Instruction{"IMP", "RTI", cpu.IMP, cpu.RTI, 1, 6, true}
	cpu.Instructions[0x41] = Instruction{"IDX", "EOR", cpu.IDX, cpu.EOR, 2, 6, true}
	cpu.Instructions[0x42] = Instruction{}
	cpu.Instructions[0x43] = Instruction{"IDX", "SRE", cpu.IDX, cpu.SRE, 2, 8, false}
	cpu.Instructions[0x44] = Instruction{"ZP0", "NOP", cpu.ZP0, cpu.ILLNOP, 2, 3, false}
	cpu.Instructions[0x45] = Instruction{"ZP0", "EOR", cpu.ZP0, cpu.EOR, 2, 3, true}
	cpu.Instructions[0x46] = Instruction{"ZP0", "LSR", cpu.ZP0, cpu.LSR, 2, 5, true}
	cpu.Instructions[0x47] = Instruction{"ZP0", "SRE", cpu.ZP0, cpu.SRE, 2, 5, false}
	cpu.Instructions[0x48] = Instruction{"IMP", "PHA", cpu.IMP, cpu.PHA, 1, 3, true}
	cpu.Instructions[0x49] = Instruction{"IMM", "EOR", cpu.IMM, cpu.EOR, 2, 2, true}
	cpu.Instructions[0x4A] = Instruction{"ACC", "LSR", cpu.ACC, cpu.LSR, 1, 2, true}
	cpu.Instructions[0x4B] = Instruction{"IMM", "ALR", cpu.IMM, cpu.ALR, 2, 2, false}
	cpu.Instructions[0x4C] = Instruction{"ABS", "JMP", cpu.ABS, cpu.JMP, 3, 3, true}
	cpu.Instructions[0x4D] = Instruction{"ABS", "EOR", cpu.ABS, cpu.EOR, 3, 4, true}
	cpu.Instructions[0x4E] = Instruction{"ABS", "LSR", cpu.ABS, cpu.LSR, 3, 6, true}
	cpu.Instructions[0x4F] = Instruction{"ABS", "SRE", cpu.ABS, cpu.SRE, 3, 6, false}

	cpu.Instructions[0x50] = Instruction{"REL", "BVC", cpu.REL, cpu.BVC, 2, 2, true}
	cpu.Instructions[0x51] = Instruction{"IZY", "EOR", cpu.IZY, cpu.EOR, 2, 5, true}
	cpu.Instructions[0x52] = Instruction{}
	cpu.Instructions[0x53] = Instruction{"IZY", "SRE", cpu.IZY, cpu.SRE, 2, 8, false}
	cpu.Instructions[0x54] = Instruction{"ZPX", "NOP", cpu.ZPX, cpu.ILLNOP, 2, 4, false}
	cpu.Instructions[0x55] = Instruction{"ZPX", "EOR", cpu.ZPX, cpu.EOR, 2, 4, true}
	cpu.Instructions[0x56] = Instruction{"ZPX", "LSR", cpu.ZPX, cpu.LSR, 2, 6, true}
	cpu.Instructions[0x57] = Instruction{"ZPX", "SRE", cpu.ZPX, cpu.SRE, 2, 6, false}
	cpu.Instructions[0x58] = Instruction{"IMP", "CLI", cpu.IMP, cpu.CLI, 1, 2, true}
	cpu.Instructions[0x59] = Instruction{"ABY", "EOR", cpu.ABY, cpu.EOR, 3, 4, true}
	cpu.Instructions[0x5A] = Instruction{"IMP", "NOP", cpu.IMP, cpu.ILLNOP, 1, 2, false}
	cpu.Instructions[0x5B] = Instruction{"ABY", "SRE", cpu.ABY, cpu.SRE, 3, 7, false}
	cpu.Instructions[0x5C] = Instruction{"ABX", "NOP", cpu.ABX, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0x5D] = Instruction{"ABX", "EOR", cpu.ABX, cpu.EOR, 3, 4, true}
	cpu.Instructions[0x5E] = Instruction{"ABX", "LSR", cpu.ABX, cpu.LSR, 3, 7, true}
	cpu.Instructions[0x5F] = Instruction{"ABX", "SRE", cpu.ABX, cpu.SRE, 3, 7, false}

	cpu.Instructions[0x60] = Instruction{"IMP", "RTS", cpu.IMP, cpu.RTS, 1, 6, true}
	cpu.Instructions[0x61] = Instruction{"IDX", "ADC", cpu.IDX, cpu.ADC, 2, 6, true}
	cpu.Instructions[0x62] = Instruction{}
	cpu.Instructions[0x63] = Instruction{"IDX", "RRA", cpu.IDX, cpu.RRA, 2, 8, false}
	cpu.Instructions[0x64] = Instruction{"ZP0", "NOP", cpu.ZP0, cpu.ILLNOP, 2, 3, false}
	cpu.Instructions[0x65] = Instruction{"ZP0", "ADC", cpu.ZP0, cpu.ADC, 2, 3, true}
	cpu.Instructions[0x66] = Instruction{"ZP0", "ROR", cpu.ZP0, cpu.ROR, 2, 5, true}
	cpu.Instructions[0x67] = Instruction{"ZP0", "RRA", cpu.ZP0, cpu.RRA, 2, 5, false}
	cpu.Instructions[0x68] = Instruction{"IMP", "PLA", cpu.IMP, cpu.PLA, 1, 4, true}
	cpu.Instructions[0x69] = Instruction{"IMM", "ADC", cpu.IMM, cpu.ADC, 2, 2, true}
	cpu.Instructions[0x6A] = Instruction{"ACC", "ROR", cpu.ACC, cpu.ROR, 1, 2, true}
	cpu.Instructions[0x6B] = Instruction{"IMM", "ARR", cpu.IMM, cpu.ARR, 2, 2, false}
	cpu.Instructions[0x6C] = Instruction{"IND", "JMP", cpu.IND, cpu.JMP, 3, 5, true}
	cpu.Instructions[0x6D] = Instruction{"ABS", "ADC", cpu.ABS, cpu.ADC, 3, 4, true}
	cpu.Instructions[0x6E] = Instruction{"ABS", "ROR", cpu.ABS, cpu.ROR, 3, 6, true}
	cpu.Instructions[0x6F] = Instruction{"ABS", "RRA", cpu.ABS, cpu.RRA, 3, 6, false}

	cpu.Instructions[0x70] = Instruction{"REL", "BVS", cpu.REL, cpu.BVS, 2, 2, true}
	cpu.Instructions[0x71] = Instruction{"IZY", "ADC", cpu.IZY, cpu.ADC, 2, 5, true}
	cpu.Instructions[0x72] = Instruction{}
	cpu.Instructions[0x73] = Instruction{"IZY", "RRA", cpu.IZY, cpu.RRA, 2, 8, false}
	cpu.Instructions[0x74] = Instruction{"ZPX", "NOP", cpu.ZPX, cpu.ILLNOP, 2, 4, false}
	cpu.Instructions[0x75] = Instruction{"ZPX", "ADC", cpu.ZPX, cpu.ADC, 2, 4, true}
	cpu.Instructions[0x76] = Instruction{"ZPX", "ROR", cpu.ZPX, cpu.ROR, 2, 6, true}
	cpu.Instructions[0x77] = Instruction{"ZPX", "RRA", cpu.ZPX, cpu.RRA, 2, 6, false}
	cpu.Instructions[0x78] = Instruction{"IMP", "SEI", cpu.IMP, cpu.SEI, 1, 2, true}
	cpu.Instructions[0x79] = Instruction{"ABY", "ADC", cpu.ABY, cpu.ADC, 3, 4, true}
	cpu.Instructions[0x7A] = Instruction{"IMP", "NOP", cpu.IMP, cpu.ILLNOP, 1, 2, false}
	cpu.Instructions[0x7B] = Instruction{"ABY", "RRA", cpu.ABY, cpu.RRA, 3, 7, false}
	cpu.Instructions[0x7C] = Instruction{"ABX", "NOP", cpu.ABX, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0x7D] = Instruction{"ABX", "ADC", cpu.ABX, cpu.ADC, 3, 4, true}
	cpu.Instructions[0x7E] = Instruction{"ABX", "ROR", cpu.ABX, cpu.ROR, 3, 7, true}
	cpu.Instructions[0x7F] = Instruction{"ABX", "RRA", cpu.ABX, cpu.RRA, 3, 7, false}

	cpu.Instructions[0x80] = Instruction{"IMM", "NOP", cpu.IMM, cpu.ILLNOP, 2, 2, false}
	cpu.Instructions[0x81] = Instruction{"IDX", "STA", cpu.IDX, cpu.STA, 2, 6, true}
	cpu.Instructions[0x82] = Instruction{"IMM", "NOP", cpu.IMM, cpu.ILLNOP, 2, 2, false}
	cpu.Instructions[0x83] = Instruction{"IDX", "SAX", cpu.IDX, cpu.SAX, 2, 6, false}
	cpu.Instructions[0x84] = Instruction{"ZP0", "STY", cpu.ZP0, cpu.STY, 2, 3, true}
	cpu.Instructions[0x85] = Instruction{"ZP0", "STA", cpu.ZP0, cpu.STA, 2, 3, true}
	cpu.Instructions[0x86] = Instruction{"ZP0", "STX", cpu.ZP0, cpu.STX, 2, 3, true}
	cpu.Instructions[0x87] = Instruction{"ZP0", "SAX", cpu.ZP0, cpu.SAX, 2, 3, false}
	cpu.Instructions[0x88] = Instruction{"IMP", "DEY", cpu.IMP, cpu.DEY, 1, 2, true}
	cpu.Instructions[0x89] = Instruction{"IMM", "NOP", cpu.IMM, cpu.ILLNOP, 2, 2, false}
	cpu.Instructions[0x8A] = Instruction{"IMP", "TXA", cpu.IMP, cpu.TXA, 1, 2, true}
	cpu.Instructions[0x8B] = Instruction{"IMM", "ANE", cpu.IMM, cpu.ANE, 2, 2, false}
	cpu.Instructions[0x8C] = Instruction{"ABS", "STY", cpu.ABS, cpu.STY, 3, 4, true}
	cpu.Instructions[0x8D] = Instruction{"ABS", "STA", cpu.ABS, cpu.STA, 3, 4, true}
	cpu.Instructions[0x8E] = Instruction{"ABS", "STX", cpu.ABS, cpu.STX, 3, 4, true}
	cpu.Instructions[0x8F] = Instruction{"ABS", "SAX", cpu.ABS, cpu.SAX, 3, 4, false}

	cpu.Instructions[0x90] = Instruction{"REL", "BCC", cpu.REL, cpu.BCC, 2, 2, true}
	cpu.Instructions[0x91] = Instruction{"IZY", "STA", cpu.IZY, cpu.STA, 2, 6, true}
	cpu.Instructions[0x92] = Instruction{}
	cpu.Instructions[0x93] = Instruction{"IZY", "SHA", cpu.IZY, cpu.SHA, 2, 6, false}
	cpu.Instructions[0x94] = Instruction{"ZPX", "STY", cpu.ZPX, cpu.STY, 2, 4, true}
	cpu.Instructions[0x95] = Instruction{"ZPX", "STA", cpu.ZPX, cpu.STA, 2, 4, true}
	cpu.Instructions[0x96] = Instruction{"ZPY", "STX", cpu.ZPY, cpu.STX, 2, 4, true}
	cpu.Instructions[0x97] = Instruction{"ZPY", "SAX", cpu.ZPY, cpu.SAX, 2, 4, false}
	cpu.Instructions[0x98] = Instruction{"IMP", "TYA", cpu.IMP, cpu.TYA, 1, 2, true}
	cpu.Instructions[0x99] = Instruction{"ABY", "STA", cpu.ABY, cpu.STA, 3, 5, true}
	cpu.Instructions[0x9A] = Instruction{"IMP", "TXS", cpu.IMP, cpu.TXS, 1, 2, true}
	cpu.Instructions[0x9B] = Instruction{"ABY", "TAS", cpu.ABY, cpu.TAS, 3, 5, false}
	cpu.Instructions[0x9C] = Instruction{"ABX", "SHY", cpu.ABX, cpu.SHY, 3, 5, false}
	cpu.Instructions[0x9D] = Instruction{"ABX", "STA", cpu.ABX, cpu.STA, 3, 5, true}
	cpu.Instructions[0x9E] = Instruction{"ABY", "SHX", cpu.ABY, cpu.SHX, 3, 5, false}
	cpu.Instructions[0x9F] = Instruction{"ABY", "SHA", cpu.ABY, cpu.SHA, 3, 5, false}

	cpu.Instructions[0xA0] = Instruction{"IMM", "LDY", cpu.IMM, cpu.LDY, 2, 2, true}
	cpu.Instructions[0xA1] = Instruction{"IDX", "LDA", cpu.IDX, cpu.LDA, 2, 6, true}
	cpu.Instructions[0xA2] = Instruction{"IMM", "LDX", cpu.IMM, cpu.LDX, 2, 2, true}
	cpu.Instructions[0xA3] = Instruction{"IDX", "LAX", cpu.IDX, cpu.LAX, 2, 6, false}
	cpu.Instructions[0xA4] = Instruction{"ZP0", "LDY", cpu.ZP0, cpu.LDY, 2, 3, true}
	cpu.Instructions[0xA5] = Instruction{"ZP0", "LDA", cpu.ZP0, cpu.LDA, 2, 3, true}
	cpu.Instructions[0xA6] = Instruction{"ZP0", "LDX", cpu.ZP0, cpu.LDX, 2, 3, true}
	cpu.Instructions[0xA7] = Instruction{"ZP0", "LAX", cpu.ZP0, cpu.LAX, 2, 3, false}
	cpu.Instructions[0xA8] = Instruction{"IMP", "TAY", cpu.IMP, cpu.TAY, 1, 2, true}
	cpu.Instructions[0xA9] = Instruction{"IMM", "LDA", cpu.IMM, cpu.LDA, 2, 2, true}
	cpu.Instructions[0xAA] = Instruction{"IMP", "TAX", cpu.IMP, cpu.TAX, 1, 2, true}
	cpu.Instructions[0xAB] = Instruction{"IMM", "LXA", cpu.IMM, cpu.LXA, 2, 2, false}
	cpu.Instructions[0xAC] = Instruction{"ABS", "LDY", cpu.ABS, cpu.LDY, 3, 4, true}
	cpu.Instructions[0xAD] = Instruction{"ABS", "LDA", cpu.ABS, cpu.LDA, 3, 4, true}
	cpu.Instructions[0xAE] = Instruction{"ABS", "LDX", cpu.ABS, cpu.LDX, 3, 4, true}
	cpu.Instructions[0xAF] = Instruction{"ABS", "LAX", cpu.ABS, cpu.LAX, 3, 4, false}

	cpu.Instructions[0xB0] = Instruction{"REL", "BCS", cpu.REL, cpu.BCS, 2, 2, true}
	cpu.Instructions[0xB1] = Instruction{"IZY", "LDA", cpu.IZY, cpu.LDA, 2, 5, true}
	cpu.Instructions[0xB2] = Instruction{}
	cpu.Instructions[0xB3] = Instruction{"IZY", "LAX", cpu.IZY, cpu.LAX, 2, 5, false}
	cpu.Instructions[0xB4] = Instruction{"ZPX", "LDY", cpu.ZPX, cpu.LDY, 2, 4, true}
	cpu.Instructions[0xB5] = Instruction{"ZPX", "LDA", cpu.ZPX, cpu.LDA, 2, 4, true}
	cpu.Instructions[0xB6] = Instruction{"ZPY", "LDX", cpu.ZPY, cpu.LDX, 2, 4, true}
	cpu.Instructions[0xB7] = Instruction{"ZPY", "LAX", cpu.ZPY, cpu.LAX, 2, 4, false}
	cpu.Instructions[0xB8] = Instruction{"IMP", "CLV", cpu.IMP, cpu.CLV, 1, 2, true}
	cpu.Instructions[0xB9] = Instruction{"ABY", "LDA", cpu.ABY, cpu.LDA, 3, 4, true}
	cpu.Instructions[0xBA] = Instruction{"IMP", "TSX", cpu.IMP, cpu.TSX, 1, 2, true}
	cpu.Instructions[0xBB] = Instruction{"ABY", "LAS", cpu.ABY, cpu.LAS, 3, 4, false}
	cpu.Instructions[0xBC] = Instruction{"ABX", "LDY", cpu.ABX, cpu.LDY, 3, 4, true}
	cpu.Instructions[0xBD] = Instruction{"ABX", "LDA", cpu.ABX, cpu.LDA, 3, 4, true}
	cpu.Instructions[0xBE] = Instruction{"ABY", "LDX", cpu.ABY, cpu.LDX, 3, 4, true}
	cpu.Instructions[0xBF] = Instruction{"ABY", "LAX", cpu.ABY, cpu.LAX, 3, 4, false}

	cpu.Instructions[0xC0] = Instruction{"IMM", "CPY", cpu.IMM, cpu.CPY, 2, 2, true}
	cpu.Instructions[0xC1] = Instruction{"IDX", "CMP", cpu.IDX, cpu.CMP, 2, 6, true}
	cpu.Instructions[0xC2] = Instruction{"IMM", "NOP", cpu.IMM, cpu.ILLNOP, 2, 2, false}
	cpu.Instructions[0xC3] = Instruction{"IDX", "DCP", cpu.IDX, cpu.DCP, 2, 8, false}
	cpu.Instructions[0xC4] = Instruction{"ZP0", "CPY", cpu.ZP0, cpu.CPY, 2, 3, true}
	cpu.Instructions[0xC5] = Instruction{"ZP0", "CMP", cpu.ZP0, cpu.CMP, 2, 3, true}
	cpu.Instructions[0xC6] = Instruction{"ZP0", "DEC", cpu.ZP0, cpu.DEC, 2, 5, true}
	cpu.Instructions[0xC7] = Instruction{"ZP0", "DCP", cpu.ZP0, cpu.DCP, 2, 5, false}
	cpu.Instructions[0xC8] = Instruction{"IMP", "INY", cpu.IMP, cpu.INY, 1, 2, true}
	cpu.Instructions[0xC9] = Instruction{"IMM", "CMP", cpu.IMM, cpu.CMP, 2, 2, true}
	cpu.Instructions[0xCA] = Instruction{"IMP", "DEX", cpu.IMP, cpu.DEX, 1, 2, true}
	cpu.Instructions[0xCB] = Instruction{"IMM", "SBX", cpu.IMM, cpu.SBX, 2, 2, false}
	cpu.Instructions[0xCC] = Instruction{"ABS", "CPY", cpu.ABS, cpu.CPY, 3, 4, true}
	cpu.Instructions[0xCD] = Instruction{"ABS", "CMP", cpu.ABS, cpu.CMP, 3, 4, true}
	cpu.Instructions[0xCE] = Instruction{"ABS", "DEC", cpu.ABS, cpu.DEC, 3, 6, true}
	cpu.Instructions[0xCF] = Instruction{"ABS", "DCP", cpu.ABS, cpu.DCP, 3, 6, false}

	cpu.Instructions[0xD0] = Instruction{"REL", "BNE", cpu.REL, cpu.BNE, 2, 2, true}
	cpu.Instructions[0xD1] = Instruction{"IZY", "CMP", cpu.IZY, cpu.CMP, 2, 5, true}
	cpu.Instructions[0xD2] = Instruction{}
	cpu.Instructions[0xD3] = Instruction{"IZY", "DCP", cpu.IZY, cpu.DCP, 2, 8, false}
	cpu.Instructions[0xD4] = Instruction{"ZPX", "NOP", cpu.ZPX, cpu.ILLNOP, 2, 4, false}
	cpu.Instructions[0xD5] = Instruction{"ZPX", "CMP", cpu.ZPX, cpu.CMP, 2, 4, true}
	cpu.Instructions[0xD6] = Instruction{"ZPX", "DEC", cpu.ZPX, cpu.DEC, 2, 6, true}
	cpu.Instructions[0xD7] = Instruction{"ZPX", "DCP", cpu.ZPX, cpu.DCP, 2, 6, false}
	cpu.Instructions[0xD8] = Instruction{"IMP", "CLD", cpu.IMP, cpu.CLD, 1, 2, true}
	cpu.Instructions[0xD9] = Instruction{"ABY", "CMP", cpu.ABY, cpu.CMP, 3, 4, true}
	cpu.Instructions[0xDA] = Instruction{"IMP", "NOP", cpu.IMP, cpu.ILLNOP, 1, 2, false}
	cpu.Instructions[0xDB] = Instruction{"ABY", "DCP", cpu.ABY, cpu.DCP, 3, 7, false}
	cpu.Instructions[0xDC] = Instruction{"ABX", "NOP", cpu.ABX, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0xDD] = Instruction{"ABX", "CMP", cpu.ABX, cpu.CMP, 3, 4, true}
	cpu.Instructions[0xDE] = Instruction{"ABX", "DEC", cpu.ABX, cpu.DEC, 3, 7, true}
	cpu.Instructions[0xDF] = Instruction{"ABX", "DCP", cpu.ABX, cpu.DCP, 3, 7, false}

	cpu.Instructions[0xE0] = Instruction{"IMM", "CPX", cpu.IMM, cpu.CPX, 2, 2, true}
	cpu.Instructions[0xE1] = Instruction{"IDX", "SBC", cpu.IDX, cpu.SBC, 2, 6, true}
	cpu.Instructions[0xE2] = Instruction{"IMM", "NOP", cpu.IMM, cpu.ILLNOP, 2, 2, false}
	cpu.Instructions[0xE3] = Instruction{"IDX", "ISB", cpu.IDX, cpu.ISC, 2, 8, false}
	cpu.Instructions[0xE4] = Instruction{"ZP0", "CPX", cpu.ZP0, cpu.CPX, 2, 3, true}
	cpu.Instructions[0xE5] = Instruction{"ZP0", "SBC", cpu.ZP0, cpu.SBC, 2, 3, true}
	cpu.Instructions[0xE6] = Instruction{"ZP0", "INC", cpu.ZP0, cpu.INC, 2, 5, true}
	cpu.Instructions[0xE7] = Instruction{"ZP0", "ISB", cpu.ZP0, cpu.ISC, 2, 5, false}
	cpu.Instructions[0xE8] = Instruction{"IMP", "INX", cpu.IMP, cpu.INX, 1, 2, true}
	cpu.Instructions[0xE9] = Instruction{"IMM", "SBC", cpu.IMM, cpu.SBC, 2, 2, true}
	cpu.Instructions[0xEA] = Instruction{"IMP", "NOP", cpu.IMP, cpu.NOP, 1, 2, true}
	cpu.Instructions[0xEB] = Instruction{"IMM", "SBC", cpu.IMM, cpu.USBC, 2, 2, false}
	cpu.Instructions[0xEC] = Instruction{"ABS", "CPX", cpu.ABS, cpu.CPX, 3, 4, true}
	cpu.Instructions[0xED] = Instruction{"ABS", "SBC", cpu.ABS, cpu.SBC, 3, 4, true}
	cpu.Instructions[0xEE] = Instruction{"ABS", "INC", cpu.ABS, cpu.INC, 3, 6, true}
	cpu.Instructions[0xEF] = Instruction{"ABS", "ISB", cpu.ABS, cpu.ISC, 3, 6, false}

	cpu.Instructions[0xF0] = Instruction{"REL", "BEQ", cpu.REL, cpu.BEQ, 2, 2, true}
	cpu.Instructions[0xF1] = Instruction{"IZY", "SBC", cpu.IZY, cpu.SBC, 2, 5, true}
	cpu.Instructions[0xF2] = Instruction{}
	cpu.Instructions[0xF3] = Instruction{"IZY", "ISB", cpu.IZY, cpu.ISC, 2, 8, false}
	cpu.Instructions[0xF4] = Instruction{"ZPX", "NOP", cpu.ZPX, cpu.ILLNOP, 2, 4, false}
	cpu.Instructions[0xF5] = Instruction{"ZPX", "SBC", cpu.ZPX, cpu.SBC, 2, 4, true}
	cpu.Instructions[0xF6] = Instruction{"ZPX", "INC", cpu.ZPX, cpu.INC, 2, 6, true}
	cpu.Instructions[0xF7] = Instruction{"ZPX", "ISB", cpu.ZPX, cpu.ISC, 2, 6, false}
	cpu.Instructions[0xF8] = Instruction{"IMP", "SED", cpu.IMP, cpu.SED, 1, 2, true}
	cpu.Instructions[0xF9] = Instruction{"ABY", "SBC", cpu.ABY, cpu.SBC, 3, 4, true}
	cpu.Instructions[0xFA] = Instruction{"IMP", "NOP", cpu.IMP, cpu.ILLNOP, 1, 2, false}
	cpu.Instructions[0xFB] = Instruction{"ABY", "ISB", cpu.ABY, cpu.ISC, 3, 7, false}
	cpu.Instructions[0xFC] = Instruction{"ABX", "NOP", cpu.ABX, cpu.ILLNOP, 3, 4, false}
	cpu.Instructions[0xFD] = Instruction{"ABX", "SBC", cpu.ABX, cpu.SBC, 3, 4, true}
	cpu.Instructions[0xFE] = Instruction{"ABX", "INC", cpu.ABX, cpu.INC, 3, 7, true}
	cpu.Instructions[0xFF] = Instruction{"ABX", "ISB", cpu.ABX, cpu.ISC, 3, 7, false}

	for i := 0; i <= 0xFF; i++ {
		cpu.Mnemonics[uint8(i)] = [2]string{cpu.Instructions[i].ExecuteMnemonic, cpu.Instructions[i].AddressModeMnemonic}
	}
}