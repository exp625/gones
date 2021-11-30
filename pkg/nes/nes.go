package nes

import (
	"fmt"
	"github.com/exp625/gones/pkg/apu"
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/cpu"
	"github.com/exp625/gones/pkg/plz"
	"github.com/exp625/gones/pkg/ppu"
	"github.com/exp625/gones/pkg/ram"
	"io"
)

// NES struct
type NES struct {
	APU  *apu.APU
	CPU  *cpu.CPU
	PPU  *ppu.PPU
	RAM  *ram.RAM
	VRAM *ram.RAM

	Cartridge *cartridge.Cartridge

	ClockTime       float64
	AudioSampleTime float64

	MasterClockCount uint64
	EmulatedTime     float64

	Logger io.ReadWriteCloser
}

// New creates a new NES instance
func New(clockTime float64, audioSampleTime float64) *NES {
	nes := &NES{
		ClockTime:       clockTime,
		AudioSampleTime: audioSampleTime,
		CPU:             cpu.New(),
		RAM:             &ram.RAM{},
		VRAM:            &ram.RAM{},
		PPU:             &ppu.PPU{},
		APU:             &apu.APU{},
	}
	nes.CPU.Bus = nes
	nes.PPU.Bus = nes
	return nes
}

// Clock will advance the master clock count by one.
// If the emulated time is greater than the time needed for one audio sample, the function returns true.
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

// Reset resets the NES to a known state
func (nes *NES) Reset() {
	nes.APU.Reset()
	nes.CPU.Reset()
	nes.PPU.Reset()
	nes.RAM.Reset()
	nes.Cartridge.Reset()
	nes.MasterClockCount = 0
	nes.EmulatedTime = 0
}

// InsertCartridge inserts the cartridge into the NES and resets the NES.
func (nes *NES) InsertCartridge(c *cartridge.Cartridge) {
	nes.Cartridge = c
	nes.Reset()
}

func (nes *NES) CPURead(location uint16) uint8 {
	mappedLocation := nes.Cartridge.CPUMapRead(location)
	switch {
	case mappedLocation <= 0x1FFF:
		_, data := nes.RAM.Read(mappedLocation % 0x0800)
		return data
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		_, data := nes.PPU.CPURead(mappedLocation)
		return data
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4017:
		// TODO: APU and I/O Registers
		return 0xFF
	case 0x4018 <= mappedLocation && mappedLocation <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
		return 0
	case 0x4020 <= mappedLocation:
		_, data := nes.Cartridge.CPURead(mappedLocation)
		return data
	default:
		panic("go is wrong")
	}
}

func (nes *NES) CPUWrite(location uint16, data uint8) {
	mappedLocation := nes.Cartridge.CPUMapWrite(location)
	switch {
	case mappedLocation <= 0x1FFF:
		nes.RAM.Write(mappedLocation%0x0800, data)
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		nes.PPU.CPUWrite(mappedLocation, data)
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4017:
		// TODO: APU and I/O Registers
	case 0x4018 <= mappedLocation && mappedLocation <= 0x401F:
		// TODO: APU and I/O functionality that is normally disabled
	case 0x4020 <= mappedLocation:
		nes.Cartridge.CPUWrite(mappedLocation, data)
	default:
		panic("go is wrong")
	}
}

func (nes *NES) PPURead(location uint16) uint8 {
	mappedLocation := nes.Cartridge.PPUMapRead(location)
	switch {
	case mappedLocation <= 0x1FFF:
		_, data := nes.Cartridge.PPURead(mappedLocation)
		return data
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		if mappedLocation >= 0x3000 {
			mappedLocation -= 0x1000
		}
		if nes.Cartridge.Mirroring() {
			// 1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
			nes.VRAM.Read((mappedLocation - 0x2000) % 0x800)
		} else {
			// 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
			if mappedLocation-0x2000 <= 0x800 {
				nes.VRAM.Read((mappedLocation - 0x2000) % 0x400)
			} else {
				nes.VRAM.Read((mappedLocation-0x2000)%0x400 + 0x400)
			}
		}
	}
	return 0
}

func (nes *NES) PPUWrite(location uint16, data uint8) {
	mappedLocation := nes.Cartridge.PPUMapRead(location)
	switch {
	case mappedLocation <= 0x1FFF:
		nes.Cartridge.PPUWrite(mappedLocation, data)
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		if mappedLocation > 0x2FFF {
			mappedLocation -= 0x1000
		}
		if nes.Cartridge.Mirroring() {
			// 1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
			nes.VRAM.Write((mappedLocation-0x2000)%0x800, data)
		} else {
			// 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
			if mappedLocation-0x2000 <= 0x800 {
				nes.VRAM.Write((mappedLocation-0x2000)%0x400, data)
			} else {
				nes.VRAM.Write((mappedLocation-0x2000)%0x400+0x400, data)
			}
		}
	}
}

func (nes *NES) Log() {
	if nes.Logger == nil {
		return
	}

	opCode := nes.CPURead(nes.CPU.PC)
	instruction := nes.CPU.Instructions[opCode]
	if instruction.Length == 0 {
		return
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
	ppuRegisters := fmt.Sprintf("PPU:%3d,%3d CYC:%d", nes.PPU.ScanLine, nes.PPU.Position, nes.CPU.ClockCount)

	logLine := fmt.Sprintf("%s  %s %s %s %s %s",
		programCounter,
		instructionBytes,
		instructionMnemonic,
		addressMnemonic,
		cpuRegisters,
		ppuRegisters,
	)

	plz.Just(fmt.Fprintln(nes.Logger, logLine))
}

func (nes *NES) addressMnemonic() string {
	opCode := nes.CPURead(nes.CPU.PC)
	instruction := nes.CPU.Instructions[opCode]
	addr, data, _ := instruction.AddressMode()
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
