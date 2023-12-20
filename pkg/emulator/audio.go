package emulator

func (e *Emulator) Read(buf []byte) (int, error) {
	// The function gets called if the audio hardware request new audio samples.
	// The length of the sample array indicates how many sample are requested.

	const TotalNumberOfBytesForAllAudioSamplesInOneFrame = 4 * AudioSampleRate / 60
	const TotalNumberOfAudioSamplesInOneFrame = AudioSampleRate / 60

	// Force the maximum sample time to be 0,016 s = 1/60

	if len(e.RemainingSamples) > 0 {
		n := copy(buf, e.RemainingSamples)
		e.RemainingSamples = e.RemainingSamples[n:]
		if n >= TotalNumberOfBytesForAllAudioSamplesInOneFrame {
			return n, nil
		}
	}

	newBuffer := make([]byte, TotalNumberOfBytesForAllAudioSamplesInOneFrame)
	for i := 0; i < TotalNumberOfAudioSamplesInOneFrame; i++ {
		var sample uint16
		if e.AutoRunEnabled {
			for {
				audioSampleReady := e.Clock()
				e.AutoRunCycles++
				if audioSampleReady {
					break
				}
			}

			// Get the audio sample for the APU
			sample = e.APU.GetAudioSample()
		} else {
			// No sound when auto run is false
			sample = 0
		}

		// Store the sample in the buffer
		newBuffer[4*i] = byte(sample)
		newBuffer[4*i+1] = byte(sample >> 8)
		newBuffer[4*i+2] = byte(sample)
		newBuffer[4*i+3] = byte(sample >> 8)
	}

	n := copy(buf, newBuffer)
	e.RemainingSamples = newBuffer[n:]

	return TotalNumberOfBytesForAllAudioSamplesInOneFrame, nil
}
