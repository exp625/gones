package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/config"
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/cartridge"
	"github.com/exp625/gones/pkg/debugger"
	"github.com/exp625/gones/pkg/file_explorer"
	"github.com/exp625/gones/pkg/input"
	"github.com/exp625/gones/pkg/logger"
	"github.com/exp625/gones/pkg/nes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"io"
	"log"
	"os"
	"path/filepath"
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
	Bindings *input.Bindings

	AutoRunEnabled bool

	ActiveScreen Screen
	FileExplorer *file_explorer.FileExplorer

	RequestedSteps int
	AutoRunCycles  int

	NanoSecondsSpentInAutoRun time.Duration
	AutoRunStarted            time.Time

	AudioContext     *audio.Context
	Player           *audio.Player
	RemainingSamples []byte
}

func New(romFile string, debug bool) (*Emulator, error) {

	var directory string
	lastROMFile, ok := config.Get(config.LastROMFile)
	if ok {
		directory = filepath.Dir(lastROMFile)
	} else {
		_, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
	}

	explorer := file_explorer.New()
	if err := explorer.Select(directory); err != nil {
		return nil, err
	}

	e := &Emulator{
		NES:          nes.New(NESClockTime, NESAudioSampleTime),
		FileExplorer: explorer,
	}
	e.Bindings = input.GetBindings()
	e.Bindings.LoadCustomBindings()
	e.registerAllBindings()
	if debug {
		e.ChangeScreen(ScreenCPU)
	}

	if romFile != "" {
		bytes, err := os.ReadFile(romFile)
		if err != nil {
			return nil, err
		}
		c := cartridge.Load(bytes, e)
		if c == nil {
			return nil, fmt.Errorf("unsupported mapper")
		}
		e.InsertCartridge(c)
		e.LoadGame()
		e.ChangeScreen(ScreenGame)
	} else {
		e.ChangeScreen(OverlayROMChooser)
	}

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
		a := NESAudioSampleTime * 1000000000.0
		i := int(a)
		e.Player.SetBufferSize((time.Duration)(i))
		e.Player.SetVolume(0.1)
		e.Player.Play()
	}

	return nil
}

func (e *Emulator) Update() error {
	textutil.Update()
	input.HandleInput(e.Bindings)

	// Measure time spent in auto run mode
	if e.AutoRunEnabled {
		e.NanoSecondsSpentInAutoRun += time.Now().Sub(e.AutoRunStarted)
	}
	e.AutoRunStarted = time.Now()

	if e.ActiveScreen == OverlayROMChooser {
		if err := e.FileExplorer.Update(); err != nil {
			return err
		}
	}

	if e.FileExplorer.Ready {
		absolutePath, err := e.FileExplorer.Get()
		if err != nil {
			return err
		}
		file, err := os.Open(absolutePath)
		if err != nil {
			log.Println("could not open file: ", err.Error())
			return nil
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Println("could not close file: ", err.Error())
			}
		}()
		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Println("failed to read file: ", err.Error())
			return nil
		}
		c := cartridge.Load(bytes, e)
		if c == nil {
			return nil
		}
		e.InsertCartridge(c)
		e.LoadGame()
		e.Reset()

		e.ChangeScreen(ScreenGame)
		e.AutoRunEnabled = true
		if err := config.Set(config.LastROMFile, absolutePath); err != nil {
			log.Println("failed to set last ROM file in config: ", err.Error())
		}

	}

	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {

	e.DrawHeader(screen)

	switch e.ActiveScreen {
	case ScreenGame:
		e.DrawOverlayGame(screen)
	case ScreenCPU:
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
	case OverlayROMChooser:
		e.DrawROMChooser(screen)
	}

}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (e *Emulator) Log() {
	if e.Logger.LoggingEnabled() {
		e.Logger.LogLine(e.Debugger.LogCpu())
	}

}
