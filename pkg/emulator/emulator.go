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

	WindowWidth  = 1200
	WindowHeight = 1000
)

var (
	top              int
	cpuText          *textutil.Text
	instructionsText *textutil.Text
	cartridgeText    *textutil.Text
	zeroPageText     *textutil.Text
	stackText        *textutil.Text
	ramText          *textutil.Text
)

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
		showsRAMPC: false,
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
	cpuText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 0, 0, 2)
	instructionsText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, -17, 200, 2)
	cartridgeText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 800-17, 200, 2)
	zeroPageText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 0-17, 370, 1)
	stackText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 400-17, 370, 1)
	ramText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 0-17, 620, 1)

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
	// Clear Text
	cpuText.Clear()
	instructionsText.Clear()
	zeroPageText.Clear()
	stackText.Clear()
	ramText.Clear()
	cartridgeText.Clear()

	// Show debug info
	if e.showsInfo {
		e.DrawCPU(cpuText)
		cpuText.Draw(screen)
	}
	if e.showsDebug {
		if !e.showsInfo {
			instructionsText.Position(0, 0)
			zeroPageText.Position(0, 170)
			stackText.Position(400, 170)
			ramText.Position(0, 420)
			cartridgeText.Position(800, 0)
		} else {
			instructionsText.Position(0, 200)
			zeroPageText.Position(0, 370)
			stackText.Position(400, 370)
			ramText.Position(0, 620)
			cartridgeText.Position(800, 200)
		}
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

	if e.showsPatternTables {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(0, 1000-(128*4))
		screen.DrawImage(e.DrawCHRROM(0), op)
		op.GeoM.Reset()
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(128*4+5, 1000-(128*4))
		screen.DrawImage(e.DrawCHRROM(1), op)
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
