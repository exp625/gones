package apu

type DMCChannel struct {
	GlobalRegister        DMCChannelGlobalRegister
	SampleAddressRegister DMCChannelSampleAddressRegister
	SampleLengthRegister  DMCChannelSampleLengthRegister
	DirectLoadRegister    DMCChannelDirectLoadRegister
}

func (D *DMCChannel) Reset() {

}

func (D *DMCChannel) GetValue() int16 {
	return 0
}
