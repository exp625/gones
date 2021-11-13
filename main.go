package main

import (
	"fmt"
	"github.com/exp625/godnes/nes"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"time"
)

const (
	AudioSampleRate =  44100
	PPUFrequency    = 5369318.0
	NESSampleTime   = 1.0 / AudioSampleRate
	NESClockTime    = 1.0 / PPUFrequency
)


// Start the main thread
func main() {
	pixelgl.Run(run)
}

// Emulator struct
type Emulator struct {
	*nes.NES
	autoRun bool
}

func run() {
	// Create Window
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create text atlas
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	statustext := text.New(pixel.V(0, 768 - atlas.LineHeight()), atlas)

	//Create NES
	emulator := &Emulator{NES: nes.NewNES(NESClockTime, NESSampleTime)}
	emulator.Reset()

	// Setup sound
	sr := beep.SampleRate(AudioSampleRate)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(Audio(emulator))

	// Render Loop
	for !win.Closed() {

		// Space will toggle the auto run mode
		if win.JustPressed(pixelgl.KeySpace) {
			emulator.autoRun = !emulator.autoRun
		}
		// Right Arrow Key issues one Master Clock
		if win.JustPressed(pixelgl.KeyRight) && !emulator.autoRun {
			emulator.Clock()
		}

		// R Key will reset the emulator
		if win.JustPressed(pixelgl.KeyR) {
			emulator.Reset()
		}

		// Display current state
		statustext.Clear()
		fmt.Fprintf(statustext, "Auto Run Mode: %t\n", emulator.autoRun)
		fmt.Fprintf(statustext, "Master Clock Count: %d\n", emulator.NES.MasterClockCount)
		win.Clear(colornames.Black)
		statustext.Draw(win, pixel.IM)

		// Update Frame
		win.Update()
	}

	// Cleanup
	speaker.Close()
}

// Audio Streamer
func Audio(emulator *Emulator) beep.Streamer {

	// The function gets called if the audio hardware request new audio samples. The length of the sample array indicates how many sample are requested.
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// If the emulator is set to auto run: Run the emulation until the time of one audio sample passed.
			if emulator.autoRun {
				for !emulator.Clock() {}

				// Get the audio sample for the APU
				sample := emulator.APU.GetAudioSample()
				samples[i][0] = sample
				samples[i][1] = sample
			} else {
				// No sound when auto run is false
				samples[i] = [2]float64{}
			}
		}
		return len(samples), true
	})
}