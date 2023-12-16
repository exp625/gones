package templates

type APUPulseChannelGlobalRegister struct {
	Duty           uint8 `bitfield:"2"`
	LoopEnvelope   bool  `bitfield:"1"`
	ConstantVolume bool  `bitfield:"1"`
	EnvelopePeriod uint8 `bitfield:"4"`
}

type APUPulseChannelSweepRegister struct {
	Enable     bool  `bitfield:"1"`
	Period     uint8 `bitfield:"3"`
	Negate     bool  `bitfield:"1"`
	ShiftCount uint8 `bitfield:"3"`
}

type APUPulseChannelTimerLowRegister struct {
	TimerLow uint8 `bitfield:"8"`
}

type APUPulseChannelTimerHighRegister struct {
	LengthCounterLoad uint8 `bitfield:"5"`
	TimerHigh         uint8 `bitfield:"3"`
}

type APUTriangleChannelGlobalRegister struct {
	Control     bool  `bitfield:"1"`
	ReloadValue uint8 `bitfield:"7"`
}

type APUTriangleChannelTimerLowRegister struct {
	TimerLow uint8 `bitfield:"8"`
}

type APUTriangleChannelTimerHighRegister struct {
	LengthCounterLoad uint8 `bitfield:"5"`
	TimerHigh         uint8 `bitfield:"3"`
}

type APUNoiseChannelGlobalRegister struct {
	_              uint8 `bitfield:"2"`
	LoopEnvelope   bool  `bitfield:"1"`
	ConstantVolume bool  `bitfield:"1"`
	EnvelopePeriod uint8 `bitfield:"4"`
}

type APUNoiseChannelPeriodRegister struct {
	Loop   bool  `bitfield:"1"`
	_      uint8 `bitfield:"3"`
	Period uint8 `bitfield:"4"`
}

type APUNoiseChannelLengthRegister struct {
	LengthCounterLoad uint8 `bitfield:"5"`
	_                 uint8 `bitfield:"3"`
}

type APUDMCChannelGlobalRegister struct {
	IRQEnable      bool  `bitfield:"1"`
	Loop           bool  `bitfield:"1"`
	_              uint8 `bitfield:"2"`
	FrequencyIndex uint8 `bitfield:"4"`
}

type APUDMCChannelSampleAddressRegister struct {
	SampleAddress uint8 `bitfield:"8"`
}

type APUDMCChannelSampleLengthRegister struct {
	SampleLength uint8 `bitfield:"8"`
}

type APUControlRegister struct {
	_              uint8 `bitfield:"3"`
	DMCEnable      bool  `bitfield:"1"`
	NoiseEnable    bool  `bitfield:"1"`
	TriangleEnable bool  `bitfield:"1"`
	Pulse2Enable   bool  `bitfield:"1"`
	Pulse1Enable   bool  `bitfield:"1"`
}

type APUFrameCounterRegister struct {
	Mode            bool  `bitfield:"1"`
	DisableFrameIRQ bool  `bitfield:"1"`
	_               uint8 `bitfield:"6"`
}
