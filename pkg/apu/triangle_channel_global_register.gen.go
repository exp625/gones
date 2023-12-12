// Code generated by running "go generate" in github.com/exp625/gones. DO NOT EDIT.

package apu

type TriangleChannelGlobalRegister uint8

func (T *TriangleChannelGlobalRegister) Control() bool {
	const bit = 1 << 7
	return *T&bit == bit
}

func (T *TriangleChannelGlobalRegister) SetControl(value bool) {
	const bit = uint8(1) << 7
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*T = TriangleChannelGlobalRegister((uint8(*T) & ^bit) | valueInt<<7)
}

func (T *TriangleChannelGlobalRegister) ReloadValue() uint8 {
	return uint8((*T >> 0) & 0x7f)
}

func (T *TriangleChannelGlobalRegister) SetReloadValue(value uint8) {
	const mask = uint8(((1 << 7) - 1) << 0)
	*T = TriangleChannelGlobalRegister((uint8(*T) & ^mask) | uint8(value)<<0)
}