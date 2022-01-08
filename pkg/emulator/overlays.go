package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

type Overlay int

const (
	OverlayGame Overlay = iota
	OverlayCPU
	OverlayPPU
	OverlayNametables
	OverlayPalettes
	OverlayControllers
	OverlaySprites
	OverlayKeybindings
	OverlayROMChooser
)

func (e *Emulator) DrawOverlayGame(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 5)
	op.GeoM.Scale(4, 4)
	screen.DrawImage(ebiten.NewImageFromImage(e.PPU.ActiveFrame), op)
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
	width, height := ebiten.WindowSize()
	text := textutil.New(basicfont.Face7x13, width, height, 4, 24, 2)
	for _, group := range e.Bindings {
		plz.Just(text.WriteString(fmt.Sprintf("%s: \n", group.Name)))
		for _, binding := range group.Bindings {
			plz.Just(text.WriteString(fmt.Sprintf("    %s: %s\n", binding.Key().String(), binding.Help)))
		}
	}
	text.Draw(screen)
}

func (e *Emulator) DrawROMChooser(screen *ebiten.Image) {
	gray := color.Gray{Y: 20}
	screen.Fill(gray)

	const pad = 20
	lines := []string{
		fmt.Sprintf("%s\n", e.FileExplorer.Directory),
		"<LEFT> to go to the parent directory\n",
		"<RIGHT> to go into the selected directory\n",
		"<UP>/<DOWN> to browse through the current directory\n",
		"<ENTER> to choose the currently selected file/directory\n",
		"<A-Z>/<0-9> to quickly selected a file/directory starting with the letter/number",
	}
	text := textutil.New(basicfont.Face7x13, screen.Bounds().Dx()-2*pad, len(lines)*basicfont.Face7x13.Height+2*pad, pad, pad, 1)
	for _, line := range lines {
		plz.Just(text.WriteString(line))
	}
	text.Draw(screen)

	img := ebiten.NewImage(screen.Bounds().Dx()-2*pad, screen.Bounds().Dy()-(len(lines)*basicfont.Face7x13.Height)-3*pad)
	img.Fill(color.Gray{Y: 70})

	e.FileExplorer.Draw(img)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pad, float64(len(lines)*basicfont.Face7x13.Height+2*pad))

	screen.DrawImage(img, op)
}
