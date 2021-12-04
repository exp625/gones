package emulator

import (
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/controller"
	"github.com/exp625/gones/pkg/nes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/basicfont"
	"io/ioutil"
	"log"
	"math"
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
	controllerText   *textutil.Text
	keybindingsText  *textutil.Text
)

type Screen int

const (
	ScreenGame Screen = iota
	ScreenDebugCPU
	ScreenDebugPPU
	ScreenDebugNametables
	ScreenDebugPalettes
	ScreenDebugController
	ScreenKeybindings
)

// Emulator struct
type Emulator struct {
	*nes.NES

	KeyBindings []*KeyBindingGroup

	autoRunEnabled bool
	LoggingEnabled bool

	Screen Screen

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
		NES:         nes.New(NESClockTime, NESAudioSampleTime),
		KeyBindings: Bindings,
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
	controllerText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 24, 1)
	keybindingsText = textutil.New(basicfont.Face7x13, WindowWidth, WindowHeight, 4, 24, 1)

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

func (e *Emulator) HandleInput() {
	for _, group := range e.KeyBindings {
		for _, binding := range group.Bindings {
			key := binding.Key()
			if binding.Pressed != nil {
				if inpututil.IsKeyJustPressed(key) {
					binding.Pressed(e)
				}
			}
			if binding.Released != nil {
				if inpututil.IsKeyJustReleased(key) {
					binding.Released(e)
				}
			}
		}
	}
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
	controllerText.Clear()
	keybindingsText.Clear()

	// Show debug info
	switch e.Screen {
	case ScreenDebugCPU:
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
	case ScreenDebugPPU:
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
	case ScreenDebugNametables:
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
		op.GeoM.Translate(0, 240*2+20)
		screen.DrawImage(e.NES.PPU.DrawNametableInColor(2), op)
		op.GeoM.Reset()
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(256*2, 240*2+20)
		screen.DrawImage(e.NES.PPU.DrawNametableInColor(3), op)
	case ScreenDebugPalettes:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(64, 30)
		op.GeoM.Translate(0, 20)
		screen.DrawImage(e.NES.PPU.DrawLoadedPalette(), op)
	case ScreenDebugController:
		offset := WindowHeight/2 - float64(ControllerImage.Bounds().Dy())/2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, WindowHeight/2-float64(ControllerImage.Bounds().Dy())/2)
		screen.DrawImage(ControllerImage, op)

		buttonAOptions := &ebiten.DrawImageOptions{}
		buttonAOptions.GeoM.Translate(804, 252+offset)

		buttonBOptions := &ebiten.DrawImageOptions{}
		buttonBOptions.GeoM.Translate(672, 252+offset)

		buttonSELECTOptions := &ebiten.DrawImageOptions{}
		buttonSELECTOptions.GeoM.Translate(362, 286+offset)

		buttonSTARTOptions := &ebiten.DrawImageOptions{}
		buttonSTARTOptions.GeoM.Translate(502, 286+offset)

		buttonUPOptions := &ebiten.DrawImageOptions{}
		buttonUPOptions.GeoM.Translate(151, 167+offset)

		buttonDOWNOptions := &ebiten.DrawImageOptions{}
		buttonDOWNOptions.GeoM.Translate(-float64(ArrowPressedImage.Bounds().Dx())/2, -float64(ArrowPressedImage.Bounds().Dy())/2)
		buttonDOWNOptions.GeoM.Rotate(math.Pi)
		buttonDOWNOptions.GeoM.Translate(float64(ArrowPressedImage.Bounds().Dx())/2, float64(ArrowPressedImage.Bounds().Dy())/2)
		buttonDOWNOptions.GeoM.Translate(151, 296+offset)

		buttonLEFTOptions := &ebiten.DrawImageOptions{}
		buttonLEFTOptions.GeoM.Translate(-float64(ArrowPressedImage.Bounds().Dx())/2, -float64(ArrowPressedImage.Bounds().Dy())/2)
		buttonLEFTOptions.GeoM.Rotate(-math.Pi / 2)
		buttonLEFTOptions.GeoM.Translate(float64(ArrowPressedImage.Bounds().Dx())/2, float64(ArrowPressedImage.Bounds().Dy())/2)
		buttonLEFTOptions.GeoM.Translate(88, 231+offset)

		buttonRIGHTOptions := &ebiten.DrawImageOptions{}
		buttonRIGHTOptions.GeoM.Translate(-float64(ArrowPressedImage.Bounds().Dx())/2, -float64(ArrowPressedImage.Bounds().Dy())/2)
		buttonRIGHTOptions.GeoM.Rotate(math.Pi / 2)
		buttonRIGHTOptions.GeoM.Translate(float64(ArrowPressedImage.Bounds().Dx())/2, float64(ArrowPressedImage.Bounds().Dy())/2)
		buttonRIGHTOptions.GeoM.Translate(215, 231+offset)

		for _, button := range []struct {
			button  uint8
			image   *ebiten.Image
			options *ebiten.DrawImageOptions
		}{
			{controller.ButtonA, CirclePressedImage, buttonAOptions},
			{controller.ButtonB, CirclePressedImage, buttonBOptions},
			{controller.ButtonSELECT, PillPressedImage, buttonSELECTOptions},
			{controller.ButtonSTART, PillPressedImage, buttonSTARTOptions},
			{controller.ButtonUP, ArrowPressedImage, buttonUPOptions},
			{controller.ButtonDOWN, ArrowPressedImage, buttonDOWNOptions},
			{controller.ButtonLEFT, ArrowPressedImage, buttonLEFTOptions},
			{controller.ButtonRIGHT, ArrowPressedImage, buttonRIGHTOptions},
		} {
			if e.Controller1.IsPressed(button.button) {
				screen.DrawImage(button.image, button.options)
			}
		}
	case ScreenKeybindings:
		e.DrawKeybindings(keybindingsText)
		keybindingsText.Draw(screen)
	}
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
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
