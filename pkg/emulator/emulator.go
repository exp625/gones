package emulator

import (
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/nes"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"io/ioutil"
	"log"
	"time"
)

const (
	AudioSampleRate    = 44100
	PPUFrequency       = 5369318.0
	NESAudioSampleTime = 1.0 / AudioSampleRate
	NESClockTime       = 1.0 / PPUFrequency

	WindowWidth  = 1200
	WindowHeight = 1000
)

// Emulator struct
type Emulator struct {
	*nes.NES

	Window *pixelgl.Window

	autoRunEnabled bool
	LoggingEnabled bool

	showsDebug         bool
	showsInfo          bool
	showsPatternTables bool
	showsRAMPC         bool

	requestedSteps int
	autoRunCycles  int

	nanoSecondsSpentInAutoRun time.Duration
	autoRunStarted            time.Time
}

func New(romFile string, debug bool) (*Emulator, error) {
	bytes, err := ioutil.ReadFile(romFile)
	if err != nil {
		log.Fatal(err)
	}
	c := cartridge.Load(bytes)

	e := &Emulator{
		NES:        nes.New(NESClockTime, NESAudioSampleTime),
		showsDebug: debug,
		showsInfo:  debug,
		showsRAMPC: debug,
	}
	e.InsertCartridge(c)

	return e, nil
}

func (e *Emulator) Run() {
	// Create Window
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, WindowWidth, WindowHeight),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}
	e.Window = window

	// Set up text atlas
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	top := WindowHeight - atlas.LineHeight()*2

	// Set up text displays
	cpuText := text.New(pixel.V(0, top), atlas)
	instructionsText := text.New(pixel.V(0, top-200), atlas)
	cartridgeText := text.New(pixel.V(800, top-200), atlas)
	zeroPageText := text.New(pixel.V(0, top-370), atlas)
	stackText := text.New(pixel.V(400, top-370), atlas)
	ramText := text.New(pixel.V(0, top-620), atlas)

	// Set up sound
	sampleRate := beep.SampleRate(AudioSampleRate)
	if err = speaker.Init(sampleRate, sampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}
	defer speaker.Close()
	speaker.Play(e.Audio())

	// Render Loop
	for !e.Window.Closed() {
		e.Window.Clear(colornames.Black)

		// Measure time spent in auto run mode
		if e.autoRunEnabled {
			e.nanoSecondsSpentInAutoRun += time.Now().Sub(e.autoRunStarted)
		}
		e.autoRunStarted = time.Now()

		// Handle user input
		e.HandleInput()

		// Show debug info
		if e.showsInfo {
			cpuText.Clear()
			e.DrawCPU(cpuText)
			cpuText.Draw(e.Window, pixel.IM.Scaled(cpuText.Orig, 2))
		}
		if e.showsDebug {
			instructionsText.Clear()
			zeroPageText.Clear()
			stackText.Clear()
			ramText.Clear()
			cartridgeText.Clear()

			e.DrawInstructions(instructionsText)
			e.DrawZeroPage(zeroPageText)
			e.DrawStack(stackText)
			e.DrawRAM(ramText)
			e.DrawCartridge(cartridgeText)

			moved := pixel.IM
			if !e.showsInfo {
				moved = moved.Moved(pixel.V(0, 200))
			}
			instructionsText.Draw(e.Window, pixel.IM.Scaled(instructionsText.Orig, 2).Chained(moved))
			cartridgeText.Draw(e.Window, pixel.IM.Scaled(cartridgeText.Orig, 2).Chained(moved))
			zeroPageText.Draw(e.Window, moved)
			stackText.Draw(e.Window, moved)
			ramText.Draw(e.Window, moved)
		}
		if e.showsPatternTables {
			e.DrawCHRROM(0).Draw(e.Window, pixel.IM.Moved(pixel.V(256+5, 256+5)).Scaled(pixel.V(256+5, 256+5), 4))
			e.DrawCHRROM(1).Draw(e.Window, pixel.IM.Moved(pixel.V(256*3+10, 256+5)).Scaled(pixel.V(256*3+10, 256+5), 4))
		}

		// Update frame
		e.Window.Update()
	}
}

func (e *Emulator) HandleInput() {
	// 'R' resets the emulator
	if e.Window.JustPressed(pixelgl.KeyR) {
		e.Reset()
	}

	// 'D' toggles the display of debug info
	if e.Window.JustPressed(pixelgl.KeyD) {
		e.showsDebug = !e.showsDebug
	}

	// 'I' toggles the display of info
	if e.Window.JustPressed(pixelgl.KeyI) {
		e.showsInfo = !e.showsInfo
	}

	// 'P' toggles the display of pattern tables
	if e.Window.JustPressed(pixelgl.KeyP) {
		e.showsPatternTables = !e.showsPatternTables
	}

	// 'X' toggles the RAM program counter (TODO: ???)
	if e.Window.JustPressed(pixelgl.KeyX) && !e.autoRunEnabled {
		e.showsRAMPC = !e.showsRAMPC
	}

	// 'L' toggles logging
	if e.Window.JustPressed(pixelgl.KeyL) {
		if e.LoggingEnabled {
			e.StopLogging()
		} else {
			e.StartLogging()
		}
	}

	// 'Enter' executes one CPU instruction if auto run is disabled
	if e.Window.JustPressed(pixelgl.KeyEnter) && !e.autoRunEnabled {
		if e.requestedSteps == 0 {
			e.requestedSteps = 1
		}
		for e.requestedSteps != 0 {
			e.Clock()
			e.Clock()
			e.Clock()
			for e.CPU.CycleCount != 0 {
				e.Clock()
				e.Clock()
				e.Clock()
			}
			e.requestedSteps--
		}
		e.requestedSteps = 0
	}

	// 'Space' toggles the auto run mode
	if e.Window.JustPressed(pixelgl.KeySpace) {
		e.autoRunEnabled = !e.autoRunEnabled
	}

	// 'Right Arrow' issues one Master Clock
	if e.Window.JustPressed(pixelgl.KeyRight) && !e.autoRunEnabled {
		e.Clock()
	}

	// 'Up Arrow' issues three Master Clocks e.g. one CPU clock
	if e.Window.JustPressed(pixelgl.KeyUp) && !e.autoRunEnabled {
		e.Clock()
		e.Clock()
		e.Clock()
	}

	// The numpad number keys add the requested digit to the end of the number of steps that will be requested
	if e.Window.JustPressed(pixelgl.KeyKP0) {
		e.requestedSteps = e.requestedSteps*10 + 0
	}
	if e.Window.JustPressed(pixelgl.KeyKP1) {
		e.requestedSteps = e.requestedSteps*10 + 1
	}
	if e.Window.JustPressed(pixelgl.KeyKP2) {
		e.requestedSteps = e.requestedSteps*10 + 2
	}
	if e.Window.JustPressed(pixelgl.KeyKP3) {
		e.requestedSteps = e.requestedSteps*10 + 3
	}
	if e.Window.JustPressed(pixelgl.KeyKP4) {
		e.requestedSteps = e.requestedSteps*10 + 4
	}
	if e.Window.JustPressed(pixelgl.KeyKP5) {
		e.requestedSteps = e.requestedSteps*10 + 5
	}
	if e.Window.JustPressed(pixelgl.KeyKP6) {
		e.requestedSteps = e.requestedSteps*10 + 6
	}
	if e.Window.JustPressed(pixelgl.KeyKP7) {
		e.requestedSteps = e.requestedSteps*10 + 7
	}
	if e.Window.JustPressed(pixelgl.KeyKP8) {
		e.requestedSteps = e.requestedSteps*10 + 8
	}
	if e.Window.JustPressed(pixelgl.KeyKP9) {
		e.requestedSteps = e.requestedSteps*10 + 9
	}

	// 'Escape' clears the number of requested steps
	if e.Window.JustPressed(pixelgl.KeyEscape) {
		e.requestedSteps = 0
	}

	// 'Q' sets the program counter to 0x4000
	if e.Window.JustPressed(pixelgl.KeyQ) && !e.autoRunEnabled {
		e.Reset()
		e.CPU.PC = 0xC000
		e.CPU.P = 0x24
	}
}

// Audio Streamer
func (e *Emulator) Audio() beep.Streamer {
	// The function gets called if the audio hardware request new audio samples.
	// The length of the sample array indicates how many sample are requested.
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// If the emulator is set to auto run:
			// Run the emulation until the time of one audio sample passed.
			if e.autoRunEnabled {
				for !e.Clock() {
					e.autoRunCycles++
				}

				// Get the audio sample for the APU
				sample := e.APU.GetAudioSample()
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
