package apu

type PulseChannel struct {
	GlobalRegister PulseChannelGlobalRegister
	SweepRegister  PulseChannelSweepRegister
	TimerLow       PulseChannelTimerLowRegister
	TimerHigh      PulseChannelTimerHighRegister
}

func (P *PulseChannel) Reset() {

}

func (P *PulseChannel) GetValue() uint8 {
	return 0
}
