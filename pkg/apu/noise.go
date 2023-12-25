package apu

/*
Rate  $0 $1  $2  $3  $4  $5   $6   $7   $8   $9   $A   $B   $C    $D    $E    $F

	--------------------------------------------------------------------------

NTSC   4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068
PAL    4, 8, 14, 30, 60, 88, 118, 148, 188, 236, 354, 472, 708,  944, 1890, 3778
*/
var NoisePeriodTable = []uint16{4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068}

type NoiseChannel struct {
	GlobalRegister NoiseChannelGlobalRegister
	LengthRegister NoiseChannelLengthRegister
	PeriodRegister NoiseChannelPeriodRegister

	TimerPeriod  uint16
	TimerCounter uint16
}

func (N *NoiseChannel) Reset() {
}

func (N *NoiseChannel) GetValue() int16 {
	return 0
}

// Clock the Noise channel audio part
func (N *NoiseChannel) ClockAudio() {
	N.TimerCounter++

	if N.TimerCounter >= N.TimerPeriod {
		N.TimerPeriod = NoisePeriodTable[N.PeriodRegister.Period()]
		N.TimerCounter = 0
		N.timerClock()
	}
}

func (N *NoiseChannel) timerClock() {
}
