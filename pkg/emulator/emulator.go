package emulator

import (
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/debugger"
	"github.com/exp625/gones/pkg/logger"
	"github.com/exp625/gones/pkg/nes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"os"
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

// Emulator struct
type Emulator struct {
	*nes.NES
	Debugger *debugger.Debugger
	Logger   logger.Logger
	Bindings []*BindingGroup

	AutoRunEnabled bool

	ActiveOverlay Overlay

	RequestedSteps int
	AutoRunCycles  int

	NanoSecondsSpentInAutoRun time.Duration
	AutoRunStarted            time.Time

	AudioContext     *audio.Context
	Player           *audio.Player
	RemainingSamples []byte
}

func New(romFile string, debug bool) (*Emulator, error) {
	bytes, err := os.ReadFile(romFile)
	if err != nil {
		return nil, err
	}
	c := cartridge.Load(bytes)

	e := &Emulator{
		NES: nes.New(NESClockTime, NESAudioSampleTime),
	}
	e.Bindings = DefaultBindings(e)
	if debug {
		e.ActiveOverlay = OverlayCPU
	}

	e.InsertCartridge(c)
	if err := e.Init(); err != nil {
		return nil, err
	}
	e.Debugger = debugger.New(e.NES)
	e.Logger = &logger.FileLogger{}
	e.CPU.Logger = e
	return e, nil
}

func (e *Emulator) Close() error {
	return e.Player.Close()
}

func (e *Emulator) Init() error {
	// Setup Audio
	if e.AudioContext == nil {
		if audio.CurrentContext() == nil {
			e.AudioContext = audio.NewContext(AudioSampleRate)
		} else {
			e.AudioContext = audio.CurrentContext()
		}
	}
	if e.Player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		e.Player, err = e.AudioContext.NewPlayer(e)
		if err != nil {
			return err
		}
		e.Player.Play()
	}

	return nil
}

func (e *Emulator) Update() error {
	textutil.Update()
	e.HandleInput()

	// Measure time spent in auto run mode
	if e.AutoRunEnabled {
		e.NanoSecondsSpentInAutoRun += time.Now().Sub(e.AutoRunStarted)
	}
	e.AutoRunStarted = time.Now()

	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {

	e.DrawHeader(screen)

	switch e.ActiveOverlay {
	case OverlayGame:
		e.DrawOverlayGame(screen)
	case OverlayCPU:
		e.DrawOverlayCPU(screen)
	case OverlayPPU:
		e.DrawOverlayPPU(screen)
	case OverlayNametables:
		e.DrawOverlayNametables(screen)
	case OverlayPalettes:
		e.DrawOverlayPalettes(screen)
	case OverlayControllers:
		e.DrawOverlayControllers(screen)
	case OverlaySprites:
		e.DrawOverlaySprites(screen)
	case OverlayKeybindings:
		e.DrawOverlayKeybindings(screen)
	}

}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (e *Emulator) HandleInput() {
	for _, group := range e.Bindings {
		for _, binding := range group.Bindings {
			key := binding.Key()
			if binding.OnPressed != nil {
				if inpututil.IsKeyJustPressed(key) {
					binding.OnPressed()
				}
			}
			if binding.OnReleased != nil {
				if inpututil.IsKeyJustReleased(key) {
					binding.OnReleased()
				}
			}
		}
	}
}

func (e *Emulator) Log() {
	if e.Logger.LoggingEnabled() {
		e.Logger.LogLine(e.Debugger.LogCpu())
	}

}
