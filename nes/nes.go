package nes

import (
	"fmt"
	"github.com/exp625/gones/nes/cartridge"
	"os"
)

// NES struct
type NES struct {
	MasterClockCount uint64
	ClockTime        float64
	AudioSampleTime  float64
	EmulatedTime     float64
	Bus              *Bus
	Logger			 *os.File
}

// New creates a new NES instance
func New(clockTime float64, audioSampleTime float64) *NES {
	c := CPU
	ram := &RAM{}
	ppu := &PPU{}
	apu := &APU{}
	bus := &Bus{
		CPU:       c,
		RAM:       ram,
		PPU:       ppu,
		APU:       apu,
	}
	bus.CPU.Bus = bus

	nes := &NES{
		MasterClockCount: 0,
		ClockTime:        clockTime,
		AudioSampleTime:  audioSampleTime,
		EmulatedTime:     0,
		Bus:              bus,
	}
	bus.CPU.NES = nes
	return nes
}

// Reset resets the NES to a know state
func (nes *NES) Reset() {
	nes.Bus.Reset()
	nes.MasterClockCount = 0
	nes.EmulatedTime = 0
}

// Clock will advance the master clock count by on. If the emulated time is greater than the
// time needed for one audio sample, the function returns true.
func (nes *NES) Clock() bool {
	audioSampleReady := false

	// Advance master clock count
	nes.MasterClockCount++

	// Clock the PPU and APU
	nes.Bus.PPU.Clock()
	nes.Bus.APU.Clock()

	// The NES CPU runs a one third of the frequency of the master clock
	if nes.MasterClockCount%3 == 0 {
		nes.Bus.CPU.Clock()
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

func (nes *NES) InsertCartridge(cat *cartridge.Cartridge)  {
	nes.Bus.Cartridge = cat
}

func (nes *NES) Log () {
	if nes.Logger != nil && nes.Bus.CPU.CurrentInstruction.Length != 0{
		// Build log line
		cpu := nes.Bus.CPU
		inst := cpu.CurrentInstruction
		logLine := fmt.Sprintf("%04X  ", cpu.CurrentPC)
		logLine += fmt.Sprintf( "%02X ", cpu.Bus.CPURead(cpu.CurrentPC))
		i := 1
		for ; i < int(inst.Length); i++ {
			logLine += fmt.Sprintf("%02X ", cpu.Bus.CPURead(cpu.CurrentPC+uint16(i)))
		}
		for ; i < 3; i++ {
			logLine += "   "
		}
		logLine += fmt.Sprint( " ",OpCodeMap[cpu.Bus.CPURead(cpu.CurrentPC)][0], " ")
		addr, data, _ := inst.AddressMode()
		// Display Address
		opCode := OpCodeMap[cpu.Bus.CPURead(cpu.CurrentPC)]
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
			logLine += fmt.Sprintf("$%02X,X @ %02X = %02X             ", addr & 0x00FF, (addr & 0x00FF) + uint16(cpu.X), data)
		case "ZPY":
			logLine += fmt.Sprintf("$%02X,Y @ %02X = %02X             ", addr & 0x00FF, (addr & 0x00FF) + uint16(cpu.Y), data)
		case "ZP0":
			logLine += fmt.Sprintf("$%02X = %02X                    ", addr & 0x00FF, data)
		case "IDX":
			// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
			logLine += fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X    ", cpu.Bus.CPURead(cpu.CurrentPC + 1), cpu.Bus.CPURead(cpu.CurrentPC + 1) + cpu.X, addr, data)
		case "IZY":
			// The second byte of the instruction points to a memory location in zero page -> content is added to Y register -> result is low order byte of the effective address
			logLine += fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X  ", cpu.Bus.CPURead(cpu.CurrentPC + 1), addr - uint16(cpu.Y), addr, data)
		case "IND":
			logLine += fmt.Sprintf("($%02X%02X) = %04X              ",cpu.Bus.CPURead(cpu.CurrentPC + 2), cpu.Bus.CPURead(cpu.CurrentPC + 1), addr )
		case "ABX":
			logLine += fmt.Sprintf("$%02X%02X,X @ %04X = %02X         ", cpu.Bus.CPURead(cpu.CurrentPC + 2), cpu.Bus.CPURead(cpu.CurrentPC + 1), addr, data)
		case "ABY":
			logLine += fmt.Sprintf("$%02X%02X,Y @ %04X = %02X         ", cpu.Bus.CPURead(cpu.CurrentPC + 2), cpu.Bus.CPURead(cpu.CurrentPC + 1), addr, data)
		default:
			logLine += fmt.Sprint("                            ")
		}

		// Add current CPU Status
		logLine += fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X ", cpu.A, cpu.X, cpu.Y, cpu.P, cpu.S)
		logLine += fmt.Sprintf("PPU:%3d,%3d ", nes.Bus.PPU.ScanLine, nes.Bus.PPU.Position)
		logLine += fmt.Sprintf("CYC:%d", cpu.ClockCount)
		fmt.Fprintln(nes.Logger, logLine)
	}

}
