// Code generated by running "go generate" in github.com/exp625/gones. DO NOT EDIT.

package apu

type ControlRegister uint8

func (C *ControlRegister) DMCEnable() bool {
	const bit = 1 << 4
	return *C&bit == bit
}

func (C *ControlRegister) SetDMCEnable(value bool) {
	const bit = uint8(1) << 4
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*C = ControlRegister((uint8(*C) & ^bit) | valueInt<<4)
}

func (C *ControlRegister) NoiseEnable() bool {
	const bit = 1 << 3
	return *C&bit == bit
}

func (C *ControlRegister) SetNoiseEnable(value bool) {
	const bit = uint8(1) << 3
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*C = ControlRegister((uint8(*C) & ^bit) | valueInt<<3)
}

func (C *ControlRegister) TriangleEnable() bool {
	const bit = 1 << 2
	return *C&bit == bit
}

func (C *ControlRegister) SetTriangleEnable(value bool) {
	const bit = uint8(1) << 2
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*C = ControlRegister((uint8(*C) & ^bit) | valueInt<<2)
}

func (C *ControlRegister) Pulse2Enable() bool {
	const bit = 1 << 1
	return *C&bit == bit
}

func (C *ControlRegister) SetPulse2Enable(value bool) {
	const bit = uint8(1) << 1
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*C = ControlRegister((uint8(*C) & ^bit) | valueInt<<1)
}

func (C *ControlRegister) Pulse1Enable() bool {
	const bit = 1 << 0
	return *C&bit == bit
}

func (C *ControlRegister) SetPulse1Enable(value bool) {
	const bit = uint8(1) << 0
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*C = ControlRegister((uint8(*C) & ^bit) | valueInt<<0)
}