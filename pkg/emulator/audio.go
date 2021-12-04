package emulator

func (e *Emulator) Read(buf []byte) (int, error) {
	// The function gets called if the audio hardware request new audio samples.
	// The length of the sample array indicates how many sample are requested.

	// Force the maximum sample time to be 0,016 s = 1/60
	buf = make([]byte, AudioSampleRate*4/60)

	if len(e.RemainingSamples) > 0 {
		n := copy(buf, e.RemainingSamples)
		e.RemainingSamples = e.RemainingSamples[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	for i := 0; i < len(buf)/4; i++ {
		if e.AutoRunEnabled {
			for !e.Clock() {
				e.AutoRunCycles++
			}

			// Get the audio sample for the APU
			sample := e.APU.GetAudioSample()
			buf[4*i] = byte(sample)
			buf[4*i+1] = byte(sample >> 8)
			buf[4*i+2] = byte(sample)
			buf[4*i+3] = byte(sample >> 8)
		} else {
			// No sound when auto run is false
			sample := 0
			buf[4*i] = byte(sample)
			buf[4*i+1] = byte(sample >> 8)
			buf[4*i+2] = byte(sample)
			buf[4*i+3] = byte(sample >> 8)
		}
	}

	if origBuf != nil {
		n := copy(origBuf, buf)
		e.RemainingSamples = buf[n:]
		return n, nil
	}
	return len(buf), nil
}
