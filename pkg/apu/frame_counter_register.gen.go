// Code generated by running "go generate" in github.com/exp625/gones. DO NOT EDIT.

package apu

type FrameCounterRegister uint8

func (F *FrameCounterRegister) FiveFrameSequence() bool {
	const bit = 1 << 7
	return *F&bit == bit
}

func (F *FrameCounterRegister) SetFiveFrameSequence(value bool) {
	const bit = uint8(1) << 7
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*F = FrameCounterRegister((uint8(*F) & ^bit) | valueInt<<7)
}

func (F *FrameCounterRegister) DisableFrameIRQ() bool {
	const bit = 1 << 6
	return *F&bit == bit
}

func (F *FrameCounterRegister) SetDisableFrameIRQ(value bool) {
	const bit = uint8(1) << 6
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*F = FrameCounterRegister((uint8(*F) & ^bit) | valueInt<<6)
}
