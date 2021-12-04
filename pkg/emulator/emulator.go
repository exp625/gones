package emulator

import (
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/nes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	WindowWidth  = 256 * 4
	WindowHeight = 240*4 + 20
)

var (
	cpuText          *textutil.Text
	instructionsText *textutil.Text
	cartridgeText    *textutil.Text
	zeroPageText     *textutil.Text
	stackText        *textutil.Text
	ramText          *textutil.Text
	ppuText          *textutil.Text
	oamText          *textutil.Text
)

const (
	ScreenGame = iota
	ScreenDebugCPU
	ScreenDebugPPU
	ScreenDebugNametables
	ScreenDebugPalettes
)

// Emulator struct
type Emulator struct {
	*nes.NES

	autoRunEnabled bool
	LoggingEnabled bool

	Screen int

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
		NES: nes.New(NESClockTime, NESAudioSampleTime),
	}
	if debug {
		e.Screen = ScreenDebugCPU
	}

	e.InsertCartridge(c)
	err = e.Init()
	if err != nil {
		log.Fatal(err)
	}
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

	// Set up text displays
	cpuText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 24, 2)
	instructionsText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 220, 2)
	cartridgeText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 800, 400, 1)
	zeroPageText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 400, 1)
	stackText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 400, 400, 1)
	ramText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 640, 1)
	ppuText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 400, 256*2+40, 1)
	oamText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 256*2+40, 1)

	return nil
}

func (e *Emulator) Update() error {
	textutil.Update()
	// Handle input
	e.HandleInput()

	// Measure time spent in auto run mode
	if e.autoRunEnabled {
		e.nanoSecondsSpentInAutoRun += time.Now().Sub(e.autoRunStarted)
	}
	e.autoRunStarted = time.Now()

	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	e.DrawHeader(screen)

	// Clear Text
	cpuText.Clear()
	instructionsText.Clear()
	zeroPageText.Clear()
	stackText.Clear()
	ramText.Clear()
	cartridgeText.Clear()
	ppuText.Clear()
	oamText.Clear()

	// Show debug info

	if e.Screen == ScreenDebugCPU {
		e.DrawCPU(cpuText)
		cpuText.Draw(screen)
		e.DrawInstructions(instructionsText)
		instructionsText.Draw(screen)
		e.DrawZeroPage(zeroPageText)
		zeroPageText.Draw(screen)
		e.DrawStack(stackText)
		stackText.Draw(screen)
		e.DrawRAM(ramText)
		ramText.Draw(screen)
		e.DrawCartridge(cartridgeText)
		cartridgeText.Draw(screen)
	}

	if e.Screen == ScreenDebugPPU {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(0, 20)
		screen.DrawImage(e.PPU.DrawPatternTable(0), op)
		op.GeoM.Reset()
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(256*2, 20)
		screen.DrawImage(e.PPU.DrawPatternTable(1), op)
		e.PPU.DrawPPUInfo(ppuText)
		ppuText.Draw(screen)
		e.PPU.DrawOAM(oamText)
		oamText.Draw(screen)
		op.GeoM.Reset()
		op.GeoM.Scale(64, 64)
		op.GeoM.Translate(0, WindowHeight-64*2)
		screen.DrawImage(e.PPU.DrawPalettes(), op)
	}

	if e.Screen == ScreenDebugNametables {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(0, 20)
		screen.DrawImage(e.NES.PPU.DrawNametableInColor(0), op)
		op.GeoM.Reset()
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(256*2, 20)
		screen.DrawImage(e.NES.PPU.DrawNametableInColor(1), op)
		op.GeoM.Reset()
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(0, 240*2 + 20)
		screen.DrawImage(e.NES.PPU.DrawNametableInColor(2), op)
		op.GeoM.Reset()
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(256*2, 240*2 + 20)
		screen.DrawImage(e.NES.PPU.DrawNametableInColor(3), op)
	}

	if e.Screen == ScreenDebugPalettes {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(64, 30)
		op.GeoM.Translate(0, 20)
		screen.DrawImage(e.NES.PPU.DrawLoadedPalette(), op)
	}
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func (e *Emulator) HandleInput() {
	// 'R' resets the emulator
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		e.Reset()
	}

	// 'F1' toggles the display of debug info
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		if e.Screen == ScreenDebugCPU {
			e.Screen = ScreenGame
		} else {
			e.Screen = ScreenDebugCPU
		}
	}

	// 'P' toggles the display of pattern tables
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		if e.Screen == ScreenDebugPPU {
			e.Screen = ScreenGame
		} else {
			e.Screen = ScreenDebugPPU
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		if e.Screen == ScreenDebugNametables {
			e.Screen = ScreenGame
		} else {
			e.Screen = ScreenDebugNametables
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		if e.Screen == ScreenDebugPalettes {
			e.Screen = ScreenGame
		} else {
			e.Screen = ScreenDebugPalettes
		}
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

	// Force the maximum sample time to be 0,016 s = 1/60
	buf = make([]byte, AudioSampleRate*4/60)

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
