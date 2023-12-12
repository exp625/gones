package apu

import (
	"github.com/exp625/gones/pkg/bus"
	"io"
	"math"
	"math/rand"
)

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

type APU struct {
	Cycle             uint64
	FrameCounterHalfs uint64
	Bus               bus.Bus

	// Registers
	ControlRegister      ControlRegister
	FrameCounterRegister FrameCounterRegister

	// Channels
	Pulse1   PulseChannel
	Pulse2   PulseChannel
	Triangle TriangleChannel
	Noise    NoiseChannel
	DMC      DMCChannel

	FrameInterrupt bool
	DMCInterrupt   bool
}

func New() *APU {
	return &APU{}
}

func (apu *APU) Clock() {
	apu.Cycle++
	if apu.Cycle%3 == 0 {
		apu.FrameCounterHalfs++
	}
	if apu.FrameCounterHalfs > 14914*2 {
		//apu.Bus.IRQ()
	}
	if apu.FrameCounterHalfs == 14915*2 {
		apu.FrameCounterHalfs = 0
	}
}

func (apu *APU) Reset() {
	apu.Cycle = 0
	apu.FrameCounterHalfs = 0
	apu.FrameInterrupt = false
	apu.DMCInterrupt = false

	apu.Pulse1.Reset()
	apu.Pulse2.Reset()
	apu.Triangle.Reset()
	apu.Noise.Reset()
	apu.DMC.Reset()
}

func (apu *APU) GetAudioSample() int16 {
	const freq = 880.0
	sinVal := (float64)(apu.Cycle) / float64(6*14915) * freq * 2 * math.Pi / 60.0
	val := (math.Sin(sinVal))*0.5*float64(math.MaxInt16) + float64(math.MaxInt16)*0.5

	/*
		Ist auch falsch jetzt, leider:

					MaxUint16    |            -
				                 |           / \
			        MaxUint16/2  |  --------/---\--------------------
				                 |               \
				    0            |_____________________________________________
	*/

	return int16(val)
}

// CPURead performs a read operation coming from the cpu bus
func (apu *APU) CPURead(location uint16) uint8 {
	if location >= 0x4000 && location <= 0x4014 {
		// Open Bus behavior:
		// https://www.nesdev.org/wiki/Open_bus_behavior
		return uint8(rand.Intn(0xFF))
	}
	if location == 0x4015 {
		// https://www.nesdev.org/wiki/APU#Status_($4015)
		value := uint8(apu.ControlRegister)
		// IF-D NT21
		if apu.DMCInterrupt {
			value |= 1 << 7
		}
		if apu.FrameInterrupt {
			value |= 1 << 6
		}
		//   N/T/2/1 will read as 1 if the corresponding length counter has not been halted through either expiring or a write of 0 to the corresponding bit. For the triangle channel, the status of the linear counter is irrelevant.
		//   D will read as 1 if the DMC bytes remaining is more than 0.
		//   Reading this register clears the frame interrupt flag (but not the DMC interrupt flag).
		apu.FrameInterrupt = false
		//   TODO: If an interrupt flag was set at the same moment of the read, it will read back as 1 but it will not be cleared.
		//   TODO: This register is internal to the CPU and so the external CPU data bus is disconnected when reading it. Therefore the returned value cannot be seen by external devices and the value does not affect open bus.
		//   TODO: Bit 5 is open bus. Because the external bus is disconnected when reading $4015, the open bus value comes from the last cycle that did not read $4015.
		value &= 0b11011111 // Just clear for now
		return value
	}
	if location >= 0x4018 && location <= 0x401F {
		// TODO: APU and I/O functionality that is normally disabled
		return 0
	}
	panic("Incorrect CPURead on APU")
}

// CPUWrite performs a write operation coming from the cpu bus
func (apu *APU) CPUWrite(location uint16, data uint8) {
	if location >= 0x4000 && location <= 0x4017 {
		switch location {
		case 0x4000:
			apu.Pulse1.GlobalRegister = PulseChannelGlobalRegister(data)
		case 0x4001:
			apu.Pulse1.SweepRegister = PulseChannelSweepRegister(data)
		case 0x4002:
			apu.Pulse1.TimerLow = PulseChannelTimerLowRegister(data)
		case 0x4003:
			apu.Pulse1.TimerHigh = PulseChannelTimerHighRegister(data)
		case 0x4004:
			apu.Pulse2.GlobalRegister = PulseChannelGlobalRegister(data)
		case 0x4005:
			apu.Pulse2.SweepRegister = PulseChannelSweepRegister(data)
		case 0x4006:
			apu.Pulse2.TimerLow = PulseChannelTimerLowRegister(data)
		case 0x4007:
			apu.Pulse2.TimerHigh = PulseChannelTimerHighRegister(data)
		case 0x4008:
			apu.Triangle.GlobalRegister = TriangleChannelGlobalRegister(data)
		case 0x4009:
			// Unused
		case 0x400A:
			apu.Triangle.TimerLow = TriangleChannelTimerLowRegister(data)
		case 0x400B:
			apu.Triangle.TimerHigh = TriangleChannelTimerHighRegister(data)
		case 0x400C:
			apu.Noise.GlobalRegister = NoiseChannelGlobalRegister(data)
		case 0x400D:
			// Unused
		case 0x400E:
			apu.Noise.PeriodRegister = NoiseChannelPeriodRegister(data)
		case 0x400F:
			apu.Noise.LengthRegister = NoiseChannelLengthRegister(data)
		case 0x4010:
			apu.DMC.GlobalRegister = DMCChannelGlobalRegister(data)
		case 0x4011:
			apu.DMC.DirectLoadRegister = DMCChannelDirectLoadRegister(data)
		case 0x4012:
			apu.DMC.SampleAddressRegister = DMCChannelSampleAddressRegister(data)
		case 0x4013:
			apu.DMC.SampleLengthRegister = DMCChannelSampleLengthRegister(data)
		case 0x4014:
			panic("DMA")
		case 0x4015:
			apu.ControlRegister = ControlRegister(data)
		case 0x4016:
			panic("Joystick")
		case 0x4017:
			apu.FrameCounterRegister = FrameCounterRegister(data)
		}
		return
	}
	if location >= 0x4018 && location <= 0x401F {
		// APU test functionality that is normally disabled.
		return
	}

	panic("Incorrect CPUWrite on APU")
}
