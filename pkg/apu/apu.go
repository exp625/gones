package apu

import (
	"io"
	"math"
	"math/rand"

	"github.com/exp625/gones/pkg/bus"
)

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

type APU struct {
	Cycle             uint64
	FrameCounterHalfs uint64
	FrameCounterReset int8
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

// AddBus connects the APU to the Bus
func (apu *APU) AddBus(bus bus.Bus) {
	apu.Bus = bus
}

// Clock the APU (this is only the logic part, not the audio part)
func (apu *APU) Clock() {
	apu.Cycle++
	apu.FrameCounterHalfs++
	// Clock channels
	apu.DMC.Clock()

	if apu.FrameCounterReset > 0 {
		apu.FrameCounterReset--
	}
	if apu.FrameCounterHalfs%2 == 0 && apu.FrameCounterReset == 0 {
		apu.FrameCounterHalfs = 0
		apu.FrameCounterReset = -1
	}

	if apu.FrameCounterHalfs == 3728*2+1 ||
		apu.FrameCounterHalfs == 7456*2+1 ||
		apu.FrameCounterHalfs == 11185*2+1 ||
		apu.FrameCounterHalfs == 14914*2+1 ||
		apu.FrameCounterHalfs == 18640*2+1 {
		// Envelopes & triangle's linear counter
	}
	if apu.FrameCounterHalfs == 7456*2+1 ||
		apu.FrameCounterHalfs == 14914*2+1 ||
		apu.FrameCounterHalfs == 18640*2+1 {
		// Length counters & sweep units
	}

	if apu.FrameCounterRegister.FiveFrameSequence() {
		if apu.FrameCounterHalfs >= 18641*2 {
			apu.FrameCounterHalfs = 0
		}
	}
	if !apu.FrameCounterRegister.FiveFrameSequence() {
		if !apu.FrameCounterRegister.DisableFrameIRQ() && apu.FrameCounterHalfs >= 14914*2 {
			apu.FrameInterrupt = true
		}
		if apu.FrameCounterHalfs >= 14915*2 {
			apu.FrameCounterHalfs = 0
		}
	}

	if apu.DMC.SampleBufferEmpty && apu.DMC.BytesRemainingCounter != 0 {
		apu.Bus.APUDMA()
	}

	// Interrupts
	if apu.FrameInterrupt || apu.DMCInterrupt {
		apu.Bus.IRQ()
	}
}

// ClockAudio will clock the audio part of the APU
func (apu *APU) ClockAudio() {
	// Clock channels
	apu.DMC.ClockAudio()
}

func (apu *APU) DMA() {
	// The sample buffer is filled with the next sample byte read from the current address, subject to whatever mapping hardware is present.
	apu.DMC.SampleBuffer = apu.Bus.CPURead(apu.DMC.AddressCounter)
	apu.DMC.SampleBufferEmpty = false
	// The address is incremented; if it exceeds $FFFF, it is wrapped around to $8000.
	apu.DMC.AddressCounter++
	if apu.DMC.AddressCounter > 0xFFFF {
		apu.DMC.AddressCounter = 0x8000
	}
	apu.DMC.BytesRemainingCounter--
	// The bytes remaining counter is decremented; if it becomes zero and the loop flag is set, the sample is restarted (see above); otherwise, if the bytes remaining counter becomes zero and the IRQ enabled flag is set, the interrupt flag is set.
	if apu.DMC.BytesRemainingCounter == 0 && apu.DMC.GlobalRegister.Loop() {
		// Sample address = %11AAAAAA.AA000000 = $C000 + (A * 64)
		apu.DMC.AddressCounter = 0xC000 + uint16(apu.DMC.SampleAddressRegister)*64
		// Sample length = %LLLL.LLLL0001 = (L * 16) + 1 bytes
		apu.DMC.BytesRemainingCounter = uint16(apu.DMC.SampleLengthRegister)*16 + 1
	} else if apu.DMC.BytesRemainingCounter == 0 && apu.DMC.GlobalRegister.IRQEnable() {
		apu.DMCInterrupt = true
	}
}

func (apu *APU) Reset() {
	apu.Cycle = 0
	apu.FrameCounterHalfs = 0
	apu.FrameInterrupt = false
	apu.DMCInterrupt = false
	apu.FrameCounterRegister = 0
	apu.FrameCounterRegister.SetDisableFrameIRQ(true)
	apu.ControlRegister = 0

	apu.DMC.Reset()
}

func (apu *APU) GetAudioSample() uint16 {
	/*
		output = pulse_out + tnd_out

		                            95.88
		pulse_out = ------------------------------------
		             (8128 / (pulse1 + pulse2)) + 100

		                                       159.79
		tnd_out = -------------------------------------------------------------
		                                    1
		           ----------------------------------------------------- + 100
		            (triangle / 8227) + (noise / 12241) + (dmc / 22638)
	*/
	pulseOut := 0.0
	if apu.Pulse1.GetValue() != 0 || apu.Pulse2.GetValue() != 0 {
		pulseOut = 95.88 / (8128/(float64(apu.Pulse1.GetValue())+float64(apu.Pulse2.GetValue())) + 100.0)
	}

	tnd_out := 0.0
	if apu.Triangle.GetValue() != 0 || apu.Noise.GetValue() != 0 || apu.DMC.GetValue() != 0 {
		tnd_out = 159.79 / ((1 / (float64(apu.Triangle.GetValue())/8227.0 + float64(apu.Noise.GetValue())/12241.0 + float64(apu.DMC.GetValue())/22638.0)) + 100.0)
	}
	out := uint16((pulseOut + tnd_out) * math.MaxInt16)
	return out
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
		//   TODO: Reading this register clears the frame interrupt flag (but not the DMC interrupt flag).
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
			// apu.Pulse1.GlobalRegister = PulseChannelGlobalRegister(data)
		case 0x4001:
			// apu.Pulse1.SweepRegister = PulseChannelSweepRegister(data)
		case 0x4002:
			// apu.Pulse1.TimerLow = PulseChannelTimerLowRegister(data)
		case 0x4003:
			// apu.Pulse1.TimerHigh = PulseChannelTimerHighRegister(data)
		case 0x4004:
			// apu.Pulse2.GlobalRegister = PulseChannelGlobalRegister(data)
		case 0x4005:
			// apu.Pulse2.SweepRegister = PulseChannelSweepRegister(data)
		case 0x4006:
			// apu.Pulse2.TimerLow = PulseChannelTimerLowRegister(data)
		case 0x4007:
			// apu.Pulse2.TimerHigh = PulseChannelTimerHighRegister(data)
		case 0x4008:
			// apu.Triangle.GlobalRegister = TriangleChannelGlobalRegister(data)
		case 0x4009:
			// Unused
		case 0x400A:
			// apu.Triangle.TimerLow = TriangleChannelTimerLowRegister(data)
		case 0x400B:
			// apu.Triangle.TimerHigh = TriangleChannelTimerHighRegister(data)
		case 0x400C:
			// apu.Noise.GlobalRegister = NoiseChannelGlobalRegister(data)
		case 0x400D:
			// Unused
		case 0x400E:
			// apu.Noise.PeriodRegister = NoiseChannelPeriodRegister(data)
		case 0x400F:
			// apu.Noise.LengthRegister = NoiseChannelLengthRegister(data)
		case 0x4010:
			apu.DMC.GlobalRegister = DMCChannelGlobalRegister(data)
			if !apu.DMC.GlobalRegister.IRQEnable() {
				apu.DMCInterrupt = false
			}
		case 0x4011:
			// The DMC output level is set to D, an unsigned value. If the timer is outputting a clock at the same time, the output level is occasionally not changed properly.
			apu.DMC.OutputLevelCounter = data & 0b01111111
		case 0x4012:
			apu.DMC.SampleAddressRegister = DMCChannelSampleAddressRegister(data)
		case 0x4013:
			apu.DMC.SampleLengthRegister = DMCChannelSampleLengthRegister(data)
		case 0x4014:
			panic("PPUDMA")
		case 0x4015:
			apu.ControlRegister = ControlRegister(data)
			apu.DMCInterrupt = false
			if !apu.ControlRegister.DMCEnable() {
				apu.DMC.BytesRemainingCounter = 0
			} else {
				apu.DMC.BytesRemainingCounter = uint16(apu.DMC.SampleLengthRegister)*16 + 1
				apu.DMC.AddressCounter = 0xC000 + uint16(apu.DMC.SampleAddressRegister)*64
			}
		case 0x4016:
			panic("Joystick")
		case 0x4017:
			apu.FrameCounterRegister = FrameCounterRegister(data)
			if apu.FrameCounterRegister.DisableFrameIRQ() {
				apu.FrameInterrupt = false
			}
			apu.FrameCounterReset = 3
			if apu.FrameCounterRegister.FiveFrameSequence() {
				// TODO: Clock quarter frame and half frame
			}
		}
		return
	}
	if location >= 0x4018 && location <= 0x401F {
		// APU test functionality that is normally disabled.
		return
	}

	panic("Incorrect CPUWrite on APU")
}
