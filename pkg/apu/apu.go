package apu

import (
	"github.com/exp625/gones/pkg/bus"
	"io"
	"math"
)

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

type APU struct {
	Cycle             uint64
	FrameCounterHalfs uint64
	Bus               bus.Bus
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
