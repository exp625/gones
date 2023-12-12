package apu

type TriangleChannel struct {
	GlobalRegister TriangleChannelGlobalRegister
	TimerLow       TriangleChannelTimerLowRegister
	TimerHigh      TriangleChannelTimerHighRegister
}

func (T *TriangleChannel) Reset() {

}

func (T *TriangleChannel) GetValue() int16 {
	return 0
}
