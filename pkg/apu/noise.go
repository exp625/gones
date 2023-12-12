package apu

type NoiseChannel struct {
	GlobalRegister NoiseChannelGlobalRegister
	LengthRegister NoiseChannelLengthRegister
	PeriodRegister NoiseChannelPeriodRegister
}

func (N *NoiseChannel) Reset() {

}

func (N *NoiseChannel) GetValue() int16 {
	return 0
}
