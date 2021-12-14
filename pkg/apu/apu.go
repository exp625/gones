package apu

import "github.com/exp625/gones/pkg/bus"

type APU struct {
	Cycle uint64

	Bus bus.Bus
}

func New() *APU {
	return &APU{}
}

func (apu *APU) Clock() {
	apu.Cycle++
	if apu.Cycle > 14914*2 {
		apu.Bus.IRQ()
	}
	if apu.Cycle == 14915*2 {
		apu.Cycle = 0
	}
}

func (apu *APU) Reset() {
}

func (apu *APU) GetAudioSample() int32 {
	return 0
}
