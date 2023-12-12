package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

type Screen int

const (
	ScreenGame Screen = iota
	OverlayCPU
	OverlayPPU
	OverlayAPU
	OverlayNametables
	OverlayPalettes
	OverlayControllers
	OverlaySprites
	SettingKeybindings
	SettingAudio
	SettingROMChooser
	SettingsSave
)

func (e *Emulator) ChangeScreen(screen Screen) {
	// If no cartridge is loaded, don't allow debug screens
	if e.NES.Cartridge == nil && screen != SettingROMChooser && screen != SettingKeybindings && screen != SettingAudio && screen != SettingsSave && screen != ScreenGame {
		return
	}

	e.clearAllBindings()
	e.ActiveScreen = screen
	switch screen {
	case SettingROMChooser:
		e.registerSettingsBindings()
		e.registerFileExplorerBindings()
	case SettingsSave:
		e.registerSettingsBindings()
		e.registerFileExplorerBindings()
	case SettingKeybindings:
		e.registerSettingsBindings()
		e.registerInputBindings()
	case SettingAudio:
		e.registerSettingsBindings()
		e.registerAudioBindings()
	case ScreenGame:
		e.registerEmulatorBindings()
	default:
		e.registerDebugBindings()
	}
}

func (e *Emulator) DrawOverlayGame(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(4, 4)
	op.GeoM.Translate((float64(screen.Bounds().Dx())-(256*4))/2, 0)
	if e.NES.Cartridge != nil {
		screen.Fill(e.PPU.Palette[e.PPU.PaletteRAM[0]][e.PPU.Mask.Emphasize()])
		screen.DrawImage(ebiten.NewImageFromImage(e.PPU.ActiveFrame), op)
	} else {
		screen.Fill(color.Black)
		noGameText := textutil.New(basicfont.Face7x13, screen.Bounds().Dx(), screen.Bounds().Dy(), (screen.Bounds().Dx())/2-basicfont.Face7x13.Advance*14, (screen.Bounds().Dy())/2-basicfont.Face7x13.Height, 2)
		plz.Just(fmt.Fprint(noGameText, "No Game Loaded"))
		noGameText.Draw(screen)
	}
}

func (e *Emulator) DrawOverlayCPU(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	cpuText := textutil.New(basicfont.Face7x13, width, height, 4, 24, 2)
	instructionsText := textutil.New(basicfont.Face7x13, width, height, 4, 220, 2)
	cartridgeText := textutil.New(basicfont.Face7x13, width, height, 800, 400, 1)
	zeroPageText := textutil.New(basicfont.Face7x13, width, height, 4, 400, 1)
	stackText := textutil.New(basicfont.Face7x13, width, height, 400, 400, 1)
	ramText := textutil.New(basicfont.Face7x13, width, height, 4, 640, 1)
	plz.Just(fmt.Fprintf(cpuText, "FPS: %0.2f \t Auto Run Mode: \t %t \t Logging Enabled: \t %t \n", ebiten.CurrentFPS(), e.AutoRunEnabled, e.Logger.LoggingEnabled()))
	plz.Just(fmt.Fprintf(cpuText, "Master CPUClock Count: \t %d\n", e.MasterClockCount))
	plz.Just(fmt.Fprintf(cpuText, "CPU CPUClock Count: \t %d \t Requested: \t %d \n", e.CPU.ClockCount, e.RequestedSteps))
	plz.Just(fmt.Fprintf(cpuText, "CPUClock Cycles Per Second (during auto run): %0.2f/s\n\n",
		1000*1000*1000*float64(e.AutoRunCycles)/(float64(e.NanoSecondsSpentInAutoRun))))
	e.Debugger.DrawCPU(cpuText)
	cpuText.Draw(screen)
	e.Debugger.DrawInstructions(instructionsText)
	instructionsText.Draw(screen)
	e.Debugger.DrawZeroPage(zeroPageText)
	zeroPageText.Draw(screen)
	e.Debugger.DrawStack(stackText)
	stackText.Draw(screen)
	e.Debugger.DrawRAM(ramText)
	ramText.Draw(screen)
	e.Debugger.DrawCartridge(cartridgeText)
	cartridgeText.Draw(screen)
}

func (e *Emulator) DrawOverlayPPU(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	ppuText := textutil.New(basicfont.Face7x13, width, height, 400, 256*2+40, 1)
	oamText := textutil.New(basicfont.Face7x13, width, height, 4, 256*2+40, 1)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.Debugger.DrawPatternTable(0), op)
	op.GeoM.Reset()
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(256*2, 20)
	screen.DrawImage(e.Debugger.DrawPatternTable(1), op)
	e.Debugger.DrawOAM(oamText)
	oamText.Draw(screen)
	e.Debugger.DrawPPUInfo(ppuText)
	ppuText.Draw(screen)
	op.GeoM.Reset()
	op.GeoM.Scale(64, 64)
	op.GeoM.Translate(0, float64(height-64*2))
	screen.DrawImage(e.Debugger.DrawPalettes(), op)
}

func (e *Emulator) DrawOverlayAPU(screen *ebiten.Image) {

}

func (e *Emulator) DrawOverlayNametables(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.Debugger.DrawNametableInColor(0), op)
	op.GeoM.Reset()
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(256*2, 20)
	screen.DrawImage(e.Debugger.DrawNametableInColor(1), op)
	op.GeoM.Reset()
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(0, 240*2+20)
	screen.DrawImage(e.Debugger.DrawNametableInColor(2), op)
	op.GeoM.Reset()
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(256*2, 240*2+20)
	screen.DrawImage(e.Debugger.DrawNametableInColor(3), op)
	op.GeoM.Reset()
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.Debugger.DrawScrollWindow(), op)
}

func (e *Emulator) DrawOverlayPalettes(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(64, 30)
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.Debugger.DrawLoadedPalette(), op)
}

func (e *Emulator) DrawOverlayControllers(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(e.Debugger.DrawController(1), op)
}

func (e *Emulator) DrawOverlaySprites(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 40)
	screen.DrawImage(e.Debugger.DrawOAMSprites(), op)
}

func (e *Emulator) DrawOverlayKeybindings(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.Bindings.Draw(), op)

}

func (e *Emulator) DrawOverlayAudioSettings(screen *ebiten.Image) {

}

func (e *Emulator) DrawSAVEChooser(screen *ebiten.Image) {

}

func (e *Emulator) DrawROMChooser(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.FileExplorer.Draw(), op)
}
