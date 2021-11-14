package main

import (
	"fmt"
	"github.com/exp625/gones/nes"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	AudioSampleRate = 44100
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

const (
	Width  = 1000
	Height = 1000
)

func run() {
	// Create Window
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, Width, Height),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create text atlas
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	statusText := text.New(pixel.V(0, Height-atlas.LineHeight()*2), atlas)
	displayText := text.New(pixel.V(0, 500), atlas)

	//Create NES
	emulator := &Emulator{NES: nes.New(NESClockTime, NESSampleTime)}

	argsWithoutProg := os.Args[1:]
	if argsWithoutProg[0] != "" {
		bytes, err := ioutil.ReadFile(argsWithoutProg[0])
		if err != nil {
			log.Fatal(err)
		}
		emulator.Bus.Cartridge.Load(bytes)
	}

	emulator.Reset()

	// Setup sound
	sr := beep.SampleRate(AudioSampleRate)
	err = speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		panic(err)
	}
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

		// Up Arrow Key issues three Master Clocks
		if win.JustPressed(pixelgl.KeyUp) && !emulator.autoRun {
			emulator.Clock()
			emulator.Clock()
			emulator.Clock()
		}

		// Enter Key one CPU instruction
		if win.JustPressed(pixelgl.KeyEnter) && !emulator.autoRun {
			emulator.Clock()
			emulator.Clock()
			emulator.Clock()
			for emulator.NES.Bus.CPU.CycleCount != 0 {
				emulator.Clock()
				emulator.Clock()
				emulator.Clock()
			}
		}

		// R Key will reset the emulator
		if win.JustPressed(pixelgl.KeyR) {
			emulator.Reset()
		}

		// Display current state
		statusText.Clear()
		displayText.Clear()
		fmt.Fprintf(statusText, "Auto Run Mode: %t\n", emulator.autoRun)
		fmt.Fprintf(statusText, "Master Clock Count: %d\n", emulator.NES.MasterClockCount)
		fmt.Fprintf(statusText, "CPU Clock Count: %d\n", emulator.NES.Bus.CPU.ClockCount)
		fmt.Fprintf(statusText, "Cycle count: %d\n", emulator.NES.Bus.CPU.CycleCount)
		fmt.Fprint(statusText, "\n")
		DrawCPU(statusText, emulator)
		DrawStack(displayText, emulator)
		fmt.Fprint(statusText, nes.OpCodeMap[emulator.Bus.CPURead(emulator.Bus.CPU.CurrentPC)], " ")
		for i := 1; emulator.Bus.CPU.CurrentInstruction.Length > i; i++ {
			fmt.Fprintf(statusText, "%02X ", emulator.Bus.CPURead(emulator.Bus.CPU.CurrentPC+uint16(i)))
		}
		fmt.Fprint(statusText, "\n")

		win.Clear(colornames.Black)
		statusText.Draw(win, pixel.IM.Scaled(statusText.Orig, 2))
		displayText.Draw(win, pixel.IM)
		// Update Frame
		win.Update()
	}

	// Cleanup
	speaker.Close()
}

func intbool(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

// Audio Streamer
func Audio(emulator *Emulator) beep.Streamer {
	// The function gets called if the audio hardware request new audio samples. The length of the sample array indicates how many sample are requested.
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// If the emulator is set to auto run: Run the emulation until the time of one audio sample passed.
			if emulator.autoRun {
				for !emulator.Clock() {
				}

				// Get the audio sample for the APU
				sample := emulator.Bus.APU.GetAudioSample()
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

func DrawCPU(statusText *text.Text, emulator *Emulator) {
	fmt.Fprintf(statusText, "Auto Run Mode: \t %t\n", emulator.autoRun)
	fmt.Fprintf(statusText, "Master Clock Count: \t %d\n", emulator.NES.MasterClockCount)
	fmt.Fprintf(statusText, "CPU Clock Count: \t %d\n", emulator.NES.Bus.CPU.ClockCount)
	fmt.Fprint(statusText, "\n")

	// Print Status
	fmt.Fprint(statusText, "CPU STATUS \t")

	arr := []string{"N", "V", "-", "B", "D", "I", "Z", "C"}
	for i := 7; i>=0; i-- {
		if emulator.Bus.CPU.GetFlag(1 << i) {
			statusText.Color = colornames.Green
		} else {
			statusText.Color = colornames.Red
		}
		fmt.Fprint(statusText, arr[i])
	}
	statusText.Color = colornames.White
	fmt.Fprint(statusText, "\n")
	fmt.Fprintf(statusText, "PC: 0x%02X\t", emulator.NES.Bus.CPU.PC)
	fmt.Fprintf(statusText, "A: 0x%02X\t", emulator.NES.Bus.CPU.A)
	fmt.Fprintf(statusText, "X: 0x%02X\t", emulator.NES.Bus.CPU.X)
	fmt.Fprintf(statusText, "Y: 0x%02X\t", emulator.NES.Bus.CPU.Y)
	fmt.Fprintf(statusText, "S: 0x%02X\t\n", emulator.NES.Bus.CPU.S)
}

func DrawRAM(statusText *text.Text, emulator *Emulator) {

}

func DrawStack(statusText *text.Text, emulator *Emulator) {
	fmt.Fprintf(statusText, "Stack: 0x%02X\n     ", emulator.NES.Bus.CPU.S)
	statusText.Color = colornames.Yellow
	for i:= 0; i<=0xF; i++ {
		fmt.Fprintf(statusText, "%02X ", uint16(i))
	}

	for i := 0x0100; i <= 0x01FF; i++ {
		if i%16 == 0 {
			statusText.Color = colornames.Yellow
			fmt.Fprintf(statusText, "\n%04X ", uint16(i & 0xFFF0))
		}
		if emulator.Bus.CPU.S == uint8(i) {
			statusText.Color = colornames.Green
		} else {
			statusText.Color = colornames.White
		}
		fmt.Fprintf(statusText, "%02X ", emulator.Bus.CPURead(uint16(i)))
	}

}
