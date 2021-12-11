package nes

import (
	"github.com/exp625/gones/pkg/apu"
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/controller"
	"github.com/exp625/gones/pkg/cpu"
	"github.com/exp625/gones/pkg/ppu"
	"github.com/exp625/gones/pkg/ram"
)

// NES struct
type NES struct {
	APU  *apu.APU
	CPU  *cpu.CPU
	PPU  *ppu.PPU
	RAM  *ram.RAM
	VRAM *ram.RAM

	Controller1 *controller.Controller
	Controller2 *controller.Controller

	Cartridge *cartridge.Cartridge

	ClockTime       float64
	AudioSampleTime float64

	MasterClockCount uint64
	EmulatedTime     float64
}

// New creates a new NES instance
func New(clockTime float64, audioSampleTime float64) *NES {
	nes := &NES{
		ClockTime:       clockTime,
		AudioSampleTime: audioSampleTime,
		CPU:             cpu.New(),
		RAM:             ram.New(),
		VRAM:            ram.New(),
		Controller1:     controller.New(),
		Controller2:     controller.New(),
		PPU:             ppu.New(),
		APU:             apu.New(),
	}

	// Wire everything up
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
	nes.VRAM.Reset()
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
	case mappedLocation == 0x4016:
		return nes.Controller1.SerialRead()
	case mappedLocation == 0x4017:
		return nes.Controller2.SerialRead()
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4015:
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
	case mappedLocation == 0x4014:
		nes.DMA(data)
	case mappedLocation == 0x4016:
		nes.Controller1.SetMode(data&0b1 == 0)
		nes.Controller2.SetMode(data&0b1 == 0)
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4015 || mappedLocation == 0x4017:
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
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3EFF:
		return nes.PPUReadRam(location)
	// $3F00-3FFF is not configurable, always mapped to the internal palette control.
	case 0x3F00 <= location:
		return nes.PPUReadPalette(location)
	}
	return 0
}

func (nes *NES) PPUReadPalette(location uint16) uint8 {
	// $3F00-$3F1F 	Palette RAM indexes
	// $3F20-$3FFF  Mirrors of $3F00-$3F1F
	mirroredLocation := (location)%0x0020 + 0x3F00
	// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C. Note that this goes for writing as well as reading.
	if mirroredLocation == 0x3F10 || mirroredLocation == 0x3F14 || mirroredLocation == 0x3F18 || mirroredLocation == 0x3F1C {
		mirroredLocation = 0x3F00
	}
	// Addresses $3F04/$3F08/$3F0C are mirrors of $3F00. Note that this goes for writing as well as reading.
	if mirroredLocation == 0x3F04 || mirroredLocation == 0x3F08 || mirroredLocation == 0x3F0C {
		mirroredLocation = 0x3F00
	}
	data := nes.PPU.PaletteRAM[mirroredLocation-0x3F00]
	return data

}

func (nes *NES) PPUReadRam(location uint16) uint8 {
	// $3000-$3EFF  -> 	Mirrors of $2000-$2EFF
	if 0x3000 <= location && location <= 0x3FFF {
		location -= 0x1000
	}
	if nes.Cartridge.Mirroring() {
		// 1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
		_, data := nes.VRAM.Read((location - 0x2000) % 0x800)
		return data
	} else {
		// 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
		if location-0x2000 < 0x800 {
			_, data := nes.VRAM.Read((location - 0x2000) % 0x400)
			return data
		} else {
			_, data := nes.VRAM.Read((location-0x2000)%0x400 + 0x400)
			return data
		}
	}
}

func (nes *NES) PPUWrite(location uint16, data uint8) {
	mappedLocation := nes.Cartridge.PPUMapRead(location)
	switch {
	case mappedLocation <= 0x1FFF:
		nes.Cartridge.PPUWrite(mappedLocation, data)
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3EFF:
		nes.PPUWriteRam(mappedLocation, data)
	// $3F00-3FFF is not configurable, always mapped to the internal palette control.
	case 0x3F00 <= location:
		nes.PPUWritePalette(mappedLocation, data)
	}
}

func (nes *NES) PPUWriteRam(location uint16, data uint8) {
	if 0x3000 <= location && location <= 0x3FFF {
		location -= 0x1000
	}
	if nes.Cartridge.Mirroring() {
		// 1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
		nes.VRAM.Write((location-0x2000)%0x800, data)
	} else {
		// 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
		if location-0x2000 < 0x800 {
			nes.VRAM.Write((location-0x2000)%0x400, data)
		} else {
			nes.VRAM.Write((location-0x2000)%0x400+0x400, data)
		}
	}
}

func (nes *NES) PPUWritePalette(location uint16, data uint8) {
	// $3F00-$3F1F 	Palette RAM indexes
	// $3F20-$3FFF  Mirrors of $3F00-$3F1F
	mirroredLocation := (location)%0x0020 + 0x3F00
	// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C. Note that this goes for writing as well as reading.
	if mirroredLocation == 0x3F10 || mirroredLocation == 0x3F14 || mirroredLocation == 0x3F18 || mirroredLocation == 0x3F1C {
		mirroredLocation -= 0x010
	}
	nes.PPU.PaletteRAM[mirroredLocation-0x3F00] = data

}

func (nes *NES) DMA(page uint8) {
	nes.CPU.DMA = true
	nes.CPU.DMAPrepared = false
	nes.CPU.DMAAddress = uint16(page) << 8
}

func (nes *NES) NMI() {
	nes.CPU.RequestNMI = true
}

func (nes *NES) IRQ() {
	nes.CPU.RequestIRQ = true
}
