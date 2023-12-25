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
	CPU  *cpu.CPU
	APU  *apu.APU
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
		APU:             apu.New(),
		RAM:             ram.New(),
		VRAM:            ram.New(),
		Controller1:     controller.New(),
		Controller2:     controller.New(),
		PPU:             ppu.New(),
	}

	// Wire everything up
	nes.CPU.AddBus(nes)
	nes.CPU.AddAPU(nes.APU)
	nes.PPU.AddBus(nes)
	nes.APU.AddBus(nes)
	return nes
}

// Clock will advance the master clock count by one.
// If the emulated time is greater than the time needed for one audio sample, the function returns true.
func (nes *NES) Clock() bool {
	audioSampleReady := false

	// Advance master clock count
	nes.MasterClockCount++

	// Clock the PPU, APU and Cartridge
	nes.Cartridge.Clock()
	nes.PPU.Clock()

	// The NES CPU runs a one third of the frequency of the master clock
	if nes.MasterClockCount%3 == 0 {
		nes.APU.Clock()
		nes.CPU.Clock()
		nes.APU.ClockAudio()
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
	nes.Cartridge.Reset()
	nes.CPU.Reset()
	nes.APU.Reset()
	nes.PPU.Reset()
	nes.RAM.Reset()
	nes.VRAM.Reset()
	nes.MasterClockCount = 0
	nes.EmulatedTime = 0
}

// InsertCartridge inserts the cartridge into the NES and resets the NES.
func (nes *NES) InsertCartridge(c *cartridge.Cartridge) {
	nes.Cartridge = c
	nes.Reset()
}

func (nes *NES) CPUMap(location uint16) uint16 {
	return nes.Cartridge.CPUMap(location)
}

func (nes *NES) CPURead(location uint16) uint8 {
	mappedLocation := nes.CPUMap(location)
	switch {
	case mappedLocation <= 0x1FFF:
		data := nes.RAM.Read(mappedLocation % 0x0800)
		return data
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		data := nes.PPU.CPURead(mappedLocation)
		return data
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4014:
		// Open Bus
		return nes.APU.CPURead(mappedLocation)
	case mappedLocation == 0x4015:
		return nes.APU.CPURead(mappedLocation)
	case mappedLocation == 0x4016:
		return nes.Controller1.SerialRead()
	case mappedLocation == 0x4017:
		return nes.Controller2.SerialRead()
	case 0x4018 <= mappedLocation && mappedLocation <= 0x401F:
		// APU test functionality that is normally disabled
		return nes.APU.CPURead(mappedLocation)
	case 0x4020 <= mappedLocation:
		data := nes.Cartridge.CPURead(mappedLocation)
		return data
	default:
		panic("go is wrong")
	}
}

func (nes *NES) CPUWrite(location uint16, data uint8) {
	mappedLocation := nes.CPUMap(location)
	switch {
	case mappedLocation <= 0x1FFF:
		nes.RAM.Write(mappedLocation%0x0800, data)
	case 0x2000 <= mappedLocation && mappedLocation <= 0x3FFF:
		nes.PPU.CPUWrite(mappedLocation, data)
	case 0x4000 <= mappedLocation && mappedLocation <= 0x4013:
		nes.APU.CPUWrite(mappedLocation, data)
	case mappedLocation == 0x4014:
		nes.PPUDMA(data)
	case mappedLocation == 0x4015:
		nes.APU.CPUWrite(mappedLocation, data)
	case mappedLocation == 0x4016:
		nes.Controller1.SetMode(data&0b1 == 0)
		nes.Controller2.SetMode(data&0b1 == 0)
	case mappedLocation == 0x4017:
		nes.APU.CPUWrite(mappedLocation, data)
		// APU test functionality that is normally disabled
	case 0x4018 <= mappedLocation && mappedLocation <= 0x401F:
		nes.APU.CPUWrite(mappedLocation, data)
	case 0x4020 <= mappedLocation:
		nes.Cartridge.CPUWrite(mappedLocation, data)
	default:
		panic("go is wrong")
	}
}

func (nes *NES) PPUMap(location uint16) uint16 {
	return nes.Cartridge.PPUMap(location)
}

func (nes *NES) PPURead(location uint16) uint8 {
	location = nes.PPUMap(location)
	switch {
	case location <= 0x1FFF:
		data := nes.Cartridge.PPURead(location)
		return data
	case 0x2000 <= location && location <= 0x3EFF:
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
	return nes.VRAM.Read(location - 0x2000)
}

func (nes *NES) PPUWrite(location uint16, data uint8) {
	location = nes.PPUMap(location)
	switch {
	case location <= 0x1FFF:
		nes.Cartridge.PPUWrite(location, data)
	case 0x2000 <= location && location <= 0x3EFF:
		nes.PPUWriteRam(location, data)
	// $3F00-3FFF is not configurable, always mapped to the internal palette control.
	case 0x3F00 <= location:
		nes.PPUWritePalette(location, data)
	}
}

func (nes *NES) PPUWriteRam(location uint16, data uint8) {
	nes.VRAM.Write(location-0x2000, data)
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

// PPUDMA triggers a PPU PPUDMA request
func (nes *NES) PPUDMA(page uint8) {
	nes.CPU.PPUDMA = true
	nes.CPU.PPUDMAPrepared = false
	nes.CPU.PPUDMAAddress = uint16(page) << 8
}

// APUDMA triggers a APU PPUDMA request
func (nes *NES) APUDMA() {
	nes.CPU.APUDMA = true
}

// NMI triggers a NMI interrupt
func (nes *NES) NMI() {
	nes.CPU.RequestNMI = true
}

// IRQ triggers a IRQ interrupt
func (nes *NES) IRQ() {
	nes.CPU.RequestIRQ = true
}
