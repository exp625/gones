package nes

import (
	"fmt"
	"github.com/exp625/gones/nes/apu"
	"github.com/exp625/gones/nes/cartridge"
	"github.com/exp625/gones/nes/cpu"
	"github.com/exp625/gones/nes/ppu"
	"github.com/exp625/gones/nes/ram"
	"os"
)

// NES struct
type NES struct {
	MasterClockCount uint64
	ClockTime        float64
	AudioSampleTime  float64
	EmulatedTime     float64

	Logger *os.File

	RAM       *ram.RAM
	CPU       *cpu.CPU6502
	PPU       *ppu.PPU
	APU       *apu.APU
	Cartridge *cartridge.Cartridge
}

// New creates a new NES instance
func New(clockTime float64, audioSampleTime float64) *NES {
	c := cpu.CPU
	ram := &ram.RAM{}
	ppu := &ppu.PPU{}
	apu := &apu.APU{}

	nes := &NES{
		MasterClockCount: 0,
		ClockTime:        clockTime,
		AudioSampleTime:  audioSampleTime,
		EmulatedTime:     0,
		RAM:              ram,
		CPU:              c,
		PPU:              ppu,
		APU:              apu,
	}
	c.Bus = nes

	return nes
}

// Reset resets the NES to a know state
func (nes *NES) Reset() {
	nes.MasterClockCount = 0
	nes.EmulatedTime = 0
	nes.RAM.Reset()
	nes.CPU.Reset()
	nes.PPU.Reset()
	nes.APU.Reset()
	nes.Cartridge.Reset()
}

// Clock will advance the master clock count by on. If the emulated time is greater than the
// time needed for one audio sample, the function returns true.
func (nes *NES) Clock() bool {
	audioSampleReady := false

	// Advance master clock count
	nes.MasterClockCount++

	// Clock the PPU and APU
	nes.PPU.Clock()
	nes.APU.Clock()

	// The NES CPU runs a one third of the frequency of the master clock
	if nes.MasterClockCount%3 == 0 {
		nes.CPU.Clock()
	}

	// Add the time for one master clock cycle to the emulated time.
	nes.EmulatedTime += nes.ClockTime
	// If the emulated time is greater than the time needed for one audio sample:
	// Reset the emulated time and set the audioSampleReady flag to true
	if nes.EmulatedTime >= nes.AudioSampleTime {
		nes.EmulatedTime -= nes.AudioSampleTime
		audioSampleReady = true
	}

	// Return if an audio sample is ready
	return audioSampleReady
}

func (nes *NES) InsertCartridge(cat *cartridge.Cartridge) {
	nes.Cartridge = cat
}

func (nes *NES) CPURead(location uint16) uint8 {
	switch {
	case location <= 0x1FFF:
		_, data := nes.RAM.Read(location % 0x0800)
		return data
	case 0x2000 <= location && location <= 0x3FFF:
		_, data := nes.PPU.Read(0x2000 + location%0x0008)
		return data
	case 0x4000 <= location && location <= 0x4017:
		// TODO: APU and I/O Registers
		return 0
	case 0x4018 <= location && location <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
		return 0
	case 0x4020 <= location:
		_, data := nes.Cartridge.CPURead(location)
		return data
	default:
		panic("go is wrong")
	}
}

func (nes *NES) CPUWrite(location uint16, data uint8) {
	switch {
	case location <= 0x1FFF:
		nes.RAM.Write(location%0x0800, data)
	case 0x2000 <= location && location <= 0x3FFF:
		nes.PPU.Write(0x2000+location%0x0008, data)
	case 0x4000 <= location && location <= 0x4017:
		// TODO: APU and I/O Registers
	case 0x4018 <= location && location <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
	case 0x4020 <= location:
		nes.Cartridge.CPUWrite(location, data)
	default:
		panic("go is wrong")
	}
}

func (nes *NES) PPURead(location uint16) uint8 {
	switch {
	case location <= 0x1FFF:
		_, data := nes.Cartridge.PPURead(location)
		return data
	}
	return 0
}

func (nes *NES) PPUWrite(location uint16, data uint8) {
	switch {
	case location <= 0x1FFF:
		nes.Cartridge.PPUWrite(location, data)
	}

}

func (nes *NES) Log() {
	if nes.Logger != nil {
		// Build log line
		opcode := nes.CPURead(nes.CPU.PC)
		inst := cpu.Instructions[opcode]
		if inst.Length == 0 {
			return
		}
		logLine := fmt.Sprintf("%04X  ", nes.CPU.PC)
		logLine += fmt.Sprintf("%02X ", nes.CPU.Bus.CPURead(nes.CPU.PC))
		i := 1
		for ; i < int(inst.Length); i++ {
			logLine += fmt.Sprintf("%02X ", nes.CPU.Bus.CPURead(nes.CPU.PC+uint16(i)))
		}
		for ; i < 3; i++ {
			logLine += "   "
		}
		logLine += fmt.Sprint(" ", cpu.OpCodeMap[nes.CPU.Bus.CPURead(nes.CPU.PC)][0], " ")
		addr, data, _ := inst.AddressMode()
		// Display Address
		opCode := cpu.OpCodeMap[nes.CPU.Bus.CPURead(nes.CPU.PC)]
		switch opCode[1] {
		case "REL":
			logLine += fmt.Sprintf("$%04X                       ", addr)
		case "ABS":
			if addr <= 0x1FFF {
				logLine += fmt.Sprintf("$%04X = %02X                  ", addr, data)
			} else {
				logLine += fmt.Sprintf("$%04X                       ", addr)
			}
		case "IMM":
			logLine += fmt.Sprintf("#$%02X                        ", data)
		case "IMP":
			logLine += fmt.Sprint("                            ")
		case "ACC":
			logLine += fmt.Sprint("A                           ")
		case "ZPX":
			logLine += fmt.Sprintf("$%02X,X @ %02X = %02X             ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
		case "ZPY":
			logLine += fmt.Sprintf("$%02X,Y @ %02X = %02X             ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
		case "ZP0":
			logLine += fmt.Sprintf("$%02X = %02X                    ", addr&0x00FF, data)
		case "IDX":
			// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
			logLine += fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X    ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), nes.CPU.Bus.CPURead(nes.CPU.PC+1)+nes.CPU.X, addr, data)
		case "IZY":
			// The second byte of the instruction points to a memory location in zero page -> content is added to Y register -> result is low order byte of the effective address
			logLine += fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X  ", nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr-uint16(nes.CPU.Y), addr, data)
		case "IND":
			logLine += fmt.Sprintf("($%02X%02X) = %04X              ", nes.CPU.Bus.CPURead(nes.CPU.PC+2), nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr)
		case "ABX":
			logLine += fmt.Sprintf("$%02X%02X,X @ %04X = %02X         ", nes.CPU.Bus.CPURead(nes.CPU.PC+2), nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
		case "ABY":
			logLine += fmt.Sprintf("$%02X%02X,Y @ %04X = %02X         ", nes.CPU.Bus.CPURead(nes.CPU.PC+2), nes.CPU.Bus.CPURead(nes.CPU.PC+1), addr, data)
		default:
			logLine += fmt.Sprint("                            ")
		}

		// Add current CPU Status
		logLine += fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X ", nes.CPU.A, nes.CPU.X, nes.CPU.Y, nes.CPU.P, nes.CPU.S)
		logLine += fmt.Sprintf("PPU:%3d,%3d ", nes.PPU.ScanLine, nes.PPU.Position)
		logLine += fmt.Sprintf("CYC:%d", nes.CPU.ClockCount)
		fmt.Fprintln(nes.Logger, logLine)
	}

}
