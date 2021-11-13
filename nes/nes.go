package nes

// NES struct
type NES struct {
	MasterClockCount uint64
	ClockTime        float64
	AudioSampleTime  float64
	EmulatedTime     float64
	Bus *Bus
}

// NewNES creates a new NES instance
func NewNES(clockTime float64, audioSampleTime float64) *NES {
	ram := &RAM{}
	cpu := &CPU{}
	ppu := &PPU{}
	apu := &APU{}
	cat := &Catridge{}
	bus := &Bus{
		CPU: cpu,
		RAM: ram,
		PPU: ppu,
		APU: apu,
		Cartridge: cat,
	}

	return &NES{
		MasterClockCount: 0,
		ClockTime:        clockTime,
		AudioSampleTime:  audioSampleTime,
		EmulatedTime:     0,
		Bus: bus,
	}
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
