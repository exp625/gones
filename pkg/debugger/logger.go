package debugger

import (
	"fmt"
)

func (nes *Debugger) LogCpu() string {

	opCode := nes.CPURead(nes.CPU.PC)
	instruction := nes.CPU.Instructions[opCode]
	if instruction.Length == 0 {
		return "ERR"
	}

	legalPrefix := " "
	if !instruction.Legal {
		legalPrefix = "*"
	}

	ib := [3]string{}
	i := 0
	for ; i < int(instruction.Length); i++ {
		ib[i] = fmt.Sprintf("%02X", nes.CPU.Bus.CPURead(nes.CPU.PC+uint16(i)))
	}
	for ; i < 3; i++ {
		ib[i] = "  "
	}

	programCounter := fmt.Sprintf("%04X", nes.CPU.PC)
	instructionBytes := fmt.Sprintf("%s %s %s", ib[0], ib[1], ib[2])
	instructionMnemonic := fmt.Sprintf("%s%s", legalPrefix, instruction.ExecuteMnemonic)
	addressMnemonic := nes.addressMnemonic()
	cpuRegisters := fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", nes.CPU.A, nes.CPU.X, nes.CPU.Y, nes.CPU.P, nes.CPU.S)
	ppuRegisters := fmt.Sprintf("PPU:%3d,%3d CYC:%d", nes.PPU.ScanLine, nes.PPU.Dot, nes.CPU.ClockCount)

	logLine := fmt.Sprintf("%s  %s %s %s %s %s",
		programCounter,
		instructionBytes,
		instructionMnemonic,
		addressMnemonic,
		cpuRegisters,
		ppuRegisters,
	)
	return logLine
}

func (nes *Debugger) addressMnemonic() string {
	opCode := nes.CPURead(nes.CPU.PC)
	instruction := nes.CPU.Instructions[opCode]
	addr, _ := instruction.AddressMode(nes.CPURead)
	data := nes.CPURead(addr)
	switch instruction.AddressModeMnemonic {
	case "REL":
		return fmt.Sprintf("$%04X                      ", addr)
	case "ABS":
		if addr <= 0x4020 || !instruction.Legal {
			return fmt.Sprintf("$%04X = %02X                 ", addr, data)
		} else {
			return fmt.Sprintf("$%04X                      ", addr)
		}
	case "IMM":
		return fmt.Sprintf("#$%02X                       ", data)
	case "IMP":
		return fmt.Sprint("                           ")
	case "ACC":
		return fmt.Sprint("A                          ")
	case "ZPX":
		return fmt.Sprintf("$%02X,X @ %02X = %02X            ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
	case "ZPY":
		return fmt.Sprintf("$%02X,Y @ %02X = %02X            ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
	case "ZP0":
		return fmt.Sprintf("$%02X = %02X                   ", addr&0x00FF, data)
	case "IDX":
		// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
		return fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X   ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), nes.CPU.Bus.CPURead(nes.CPU.PC+1)+nes.CPU.X, addr, data)
	case "IZY":
		// The second byte of the instruction points to a memory location in zero page -> content is added to Y register -> result is low order byte of the effective address
		return fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr-uint16(nes.CPU.Y), addr, data)
	case "IND":
		return fmt.Sprintf("($%02X%02X) = %04X             ", nes.CPU.Bus.CPURead(nes.CPU.PC+2), nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr)
	case "ABX":
		return fmt.Sprintf("$%02X%02X,X @ %04X = %02X        ", nes.CPU.Bus.CPURead(nes.CPU.PC+2), nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
	case "ABY":
		return fmt.Sprintf("$%02X%02X,Y @ %04X = %02X        ", nes.CPU.Bus.CPURead(nes.CPU.PC+2), nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
	default:
		return fmt.Sprint("                           ")
	}
}
