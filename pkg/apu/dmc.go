package apu

import (
	"github.com/exp625/gones/internal/shift_register"
)

/*
Rate   $0   $1   $2   $3   $4   $5   $6   $7   $8   $9   $A   $B   $C   $D   $E   $F
      ------------------------------------------------------------------------------
NTSC  428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106,  84,  72,  54
PAL   398, 354, 316, 298, 276, 236, 210, 198, 176, 148, 132, 118,  98,  78,  66,  50
*/

var PeriodTable = []uint16{398, 354, 316, 298, 276, 236, 210, 198, 176, 148, 132, 118, 98, 78, 66, 50}

type DMCChannel struct {
	GlobalRegister        DMCChannelGlobalRegister
	SampleAddressRegister DMCChannelSampleAddressRegister
	SampleLengthRegister  DMCChannelSampleLengthRegister

	// Timer
	TimerPeriod  uint16
	TimerCounter uint16

	// MemoryReader
	AddressCounter        uint16
	BytesRemainingCounter uint16
	SampleBuffer          uint8
	SampleBufferEmpty     bool

	// Output unit
	OutputShiftRegister  shift_register.ShiftRegister8
	BitsRemainingCounter uint8
	OutputLevelCounter   uint8
	SilenceFlag          bool
}

func (D *DMCChannel) Clock() {
	/*
		When an output cycle ends, a new cycle is started as follows:
		The bits-remaining counter is loaded with 8.
		If the sample buffer is empty, then the silence flag is set; otherwise, the silence flag is cleared and the sample buffer is emptied into the shift register.
	*/
	if D.BitsRemainingCounter == 0 {
		D.BitsRemainingCounter = 8
		if D.SampleBufferEmpty {
			D.SilenceFlag = true
		} else {
			D.OutputShiftRegister.Set(D.SampleBuffer)
			D.SampleBufferEmpty = true
		}
	}
}

// Clock the DMC channel audio part
func (D *DMCChannel) ClockAudio() {
	D.TimerCounter++

	if D.TimerCounter >= D.TimerPeriod {
		D.TimerPeriod = PeriodTable[D.GlobalRegister.FrequencyIndex()]
		D.TimerCounter = 0
		D.timerClock()
	}

}

func (D *DMCChannel) timerClock() {
	/*
		Output Unit:
		When the timer outputs a clock, the following actions occur in order:
		If the silence flag is clear, the output level changes based on bit 0 of the shift register. If the bit is 1, add 2; otherwise, subtract 2. But if adding or subtracting 2 would cause the output level to leave the 0-127 range, leave the output level unchanged. This means subtract 2 only if the current level is at least 2, or add 2 only if the current level is at most 125.
		The right shift register is clocked.
		As stated above, the bits-remaining counter is decremented. If it becomes zero, a new output cycle is started.
	*/
	if D.SilenceFlag {
		if D.OutputShiftRegister.GetBit(7) == 1 {
			if D.OutputLevelCounter <= 125 {
				D.OutputLevelCounter += 2
			}
		} else {
			if D.OutputLevelCounter >= 2 {
				D.OutputLevelCounter -= 2
			}
		}
	}

	D.OutputShiftRegister.ShiftLeft(0)
	D.BitsRemainingCounter--
}

func (D *DMCChannel) Reset() {
	D.OutputLevelCounter = 0
	D.TimerCounter = 0
	D.SilenceFlag = false
	D.OutputShiftRegister.Set(0)
	D.SampleBufferEmpty = true
	D.SampleBuffer = 0
	D.BitsRemainingCounter = 0
	D.SampleAddressRegister = 0
	D.SampleLengthRegister = 0
	D.BytesRemainingCounter = 0
}

func (D *DMCChannel) GetValue() uint8 {
	return D.OutputLevelCounter
}
