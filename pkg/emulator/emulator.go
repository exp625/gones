package emulator

import (
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/nes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
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

var (
	normalFont font.Face
	top        int
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Emulator struct
type Emulator struct {
	*nes.NES

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

	audioContext     *audio.Context
	player           *audio.Player
	remainingSamples []byte
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

func (e *Emulator) Init() error {

	// Setup Audio
	if e.audioContext == nil {
		e.audioContext = audio.NewContext(AudioSampleRate)
	}
	if e.player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		e.player, err = e.audioContext.NewPlayer(e)
		if err != nil {
			return err
		}
		e.player.Play()
	}

	top = 0

	// Set up text displays
	/*
		cpuText := text.New(pixel.V(0, top), atlas)
		instructionsText := text.New(pixel.V(0, top-200), atlas)
		cartridgeText := text.New(pixel.V(800, top-200), atlas)
		zeroPageText := text.New(pixel.V(0, top-370), atlas)
		stackText := text.New(pixel.V(400, top-370), atlas)
		ramText := text.New(pixel.V(0, top-620), atlas)
	*/
	return nil

}

func (e *Emulator) Update() error {

	// Handle input
	e.HandleInput()

	// Measure time spent in auto run mode
	if e.autoRunEnabled {
		e.nanoSecondsSpentInAutoRun += time.Now().Sub(e.autoRunStarted)
	}
	e.autoRunStarted = time.Now()

	// Handle user input
	e.HandleInput()

	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {

	// Show debug info
	if e.showsInfo {
		text.Draw(screen, e.DrawCPU(), normalFont, 0, 0, color.White)
	}
	if e.showsDebug {
		if !e.showsInfo {
			top = 200
		} else {
			top = 0
		}
		text.Draw(screen, e.DrawInstructions(), normalFont, 0, 200-top, color.White)
		text.Draw(screen, e.DrawZeroPage(), normalFont, 0, 370-top, color.White)
		text.Draw(screen, e.DrawStack(), normalFont, 400, 370-top, color.White)
		text.Draw(screen, e.DrawRAM(), normalFont, 0, 620-top, color.White)
		text.Draw(screen, e.DrawCartridge(), normalFont, 800, 200-top, color.White)

	}
	/*
		if e.showsPatternTables {
			e.DrawCHRROM(0).Draw(e.Window, pixel.IM.Moved(pixel.V(256+5, 256+5)).Scaled(pixel.V(256+5, 256+5), 4))
			e.DrawCHRROM(1).Draw(e.Window, pixel.IM.Moved(pixel.V(256*3+10, 256+5)).Scaled(pixel.V(256*3+10, 256+5), 4))
		}
	*/
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func (e *Emulator) HandleInput() {
	// 'R' resets the emulator
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		e.Reset()
	}

	// 'D' toggles the display of debug info
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		e.showsDebug = !e.showsDebug
	}

	// 'I' toggles the display of info
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		e.showsInfo = !e.showsInfo
	}

	// 'P' toggles the display of pattern tables
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		e.showsPatternTables = !e.showsPatternTables
	}

	// 'X' toggles the RAM program counter (TODO: ???)
	if inpututil.IsKeyJustPressed(ebiten.KeyX) && !e.autoRunEnabled {
		e.showsRAMPC = !e.showsRAMPC
	}

	// 'L' toggles logging
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		if e.LoggingEnabled {
			e.StopLogging()
		} else {
			e.StartLogging()
		}
	}

	// 'Enter' executes one CPU instruction if auto run is disabled
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !e.autoRunEnabled {
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		e.autoRunEnabled = !e.autoRunEnabled
	}

	// 'Right Arrow' issues one Master Clock
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && !e.autoRunEnabled {
		e.Clock()
	}

	// 'Up Arrow' issues three Master Clocks e.g. one CPU clock
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && !e.autoRunEnabled {
		e.Clock()
		e.Clock()
		e.Clock()
	}

	// The numpad number keys add the requested digit to the end of the number of steps that will be requested
	if inpututil.IsKeyJustPressed(ebiten.KeyKP0) {
		e.requestedSteps = e.requestedSteps*10 + 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP1) {
		e.requestedSteps = e.requestedSteps*10 + 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP2) {
		e.requestedSteps = e.requestedSteps*10 + 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP3) {
		e.requestedSteps = e.requestedSteps*10 + 3
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP4) {
		e.requestedSteps = e.requestedSteps*10 + 4
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP5) {
		e.requestedSteps = e.requestedSteps*10 + 5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP6) {
		e.requestedSteps = e.requestedSteps*10 + 6
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP7) {
		e.requestedSteps = e.requestedSteps*10 + 7
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP8) {
		e.requestedSteps = e.requestedSteps*10 + 8
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP9) {
		e.requestedSteps = e.requestedSteps*10 + 9
	}

	// 'Escape' clears the number of requested steps
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		e.requestedSteps = 0
	}

	// 'Q' sets the program counter to 0x4000
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && !e.autoRunEnabled {
		e.Reset()
		e.CPU.PC = 0xC000
		e.CPU.P = 0x24
	}
}

// Audio Streamer
func (e *Emulator) Read(buf []byte) (int, error) {
	// The function gets called if the audio hardware request new audio samples.
	// The length of the sample array indicates how many sample are requested.

	if len(e.remainingSamples) > 0 {
		n := copy(buf, e.remainingSamples)
		e.remainingSamples = e.remainingSamples[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	for i := 0; i < len(buf)/4; i++ {
		if e.autoRunEnabled {
			for !e.Clock() {
				e.autoRunCycles++
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
		e.remainingSamples = buf[n:]
		return n, nil
	}
	return len(buf), nil
}

func (e *Emulator) Close() error {
	return nil
}
