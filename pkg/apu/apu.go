package apu

type APU struct {
}

func New() *APU {
	return &APU{}
}

func (apu *APU) Clock() {
}

func (apu *APU) Reset() {
}

func (apu *APU) GetAudioSample() int32 {
	return 0
}
