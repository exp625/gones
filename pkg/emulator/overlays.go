package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/controller"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"math"
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
)

func (e *Emulator) DrawOverlayGame(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 20)
	op.GeoM.Scale(4, 4)
	screen.DrawImage(e.PPU.DrawGameDebug(), op)
}

func (e *Emulator) DrawOverlayCPU(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	cpuText := textutil.New(basicfont.Face7x13, width, height, 4, 24, 2)
	instructionsText := textutil.New(basicfont.Face7x13, width, height, 4, 220, 2)
	cartridgeText := textutil.New(basicfont.Face7x13, width, height, 800, 400, 1)
	zeroPageText := textutil.New(basicfont.Face7x13, width, height, 4, 400, 1)
	stackText := textutil.New(basicfont.Face7x13, width, height, 400, 400, 1)
	ramText := textutil.New(basicfont.Face7x13, width, height, 4, 640, 1)
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

func (e *Emulator) DrawOverlayPPU(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	ppuText := textutil.New(basicfont.Face7x13, width, height, 400, 256*2+40, 1)
	oamText := textutil.New(basicfont.Face7x13, width, height, 4, 256*2+40, 1)
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
	op.GeoM.Translate(0, float64(height-64*2))
	screen.DrawImage(e.PPU.DrawPalettes(), op)
}

func (e *Emulator) DrawOverlayNametables(screen *ebiten.Image) {
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
}

func (e *Emulator) DrawOverlayPalettes(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(64, 30)
	op.GeoM.Translate(0, 20)
	screen.DrawImage(e.NES.PPU.DrawLoadedPalette(), op)
}

func (e *Emulator) DrawOverlayControllers(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()

	xOffset := float64(width/2) - float64(ControllerImage.Bounds().Dx())/2
	yOffset := float64(height)/2 - float64(ControllerImage.Bounds().Dy())/2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(xOffset, yOffset)
	screen.DrawImage(ControllerImage, op)

	buttonAOptions := &ebiten.DrawImageOptions{}
	buttonAOptions.GeoM.Translate(804+xOffset, 252+yOffset)

	buttonBOptions := &ebiten.DrawImageOptions{}
	buttonBOptions.GeoM.Translate(672+xOffset, 252+yOffset)

	buttonSELECTOptions := &ebiten.DrawImageOptions{}
	buttonSELECTOptions.GeoM.Translate(362+xOffset, 286+yOffset)

	buttonSTARTOptions := &ebiten.DrawImageOptions{}
	buttonSTARTOptions.GeoM.Translate(502+xOffset, 286+yOffset)

	buttonUPOptions := &ebiten.DrawImageOptions{}
	buttonUPOptions.GeoM.Translate(151+xOffset, 167+yOffset)

	buttonDOWNOptions := &ebiten.DrawImageOptions{}
	buttonDOWNOptions.GeoM.Translate(-float64(ArrowPressedImage.Bounds().Dx())/2, -float64(ArrowPressedImage.Bounds().Dy())/2)
	buttonDOWNOptions.GeoM.Rotate(math.Pi)
	buttonDOWNOptions.GeoM.Translate(float64(ArrowPressedImage.Bounds().Dx())/2, float64(ArrowPressedImage.Bounds().Dy())/2)
	buttonDOWNOptions.GeoM.Translate(151+xOffset, 296+yOffset)

	buttonLEFTOptions := &ebiten.DrawImageOptions{}
	buttonLEFTOptions.GeoM.Translate(-float64(ArrowPressedImage.Bounds().Dx())/2, -float64(ArrowPressedImage.Bounds().Dy())/2)
	buttonLEFTOptions.GeoM.Rotate(-math.Pi / 2)
	buttonLEFTOptions.GeoM.Translate(float64(ArrowPressedImage.Bounds().Dx())/2, float64(ArrowPressedImage.Bounds().Dy())/2)
	buttonLEFTOptions.GeoM.Translate(88+xOffset, 231+yOffset)

	buttonRIGHTOptions := &ebiten.DrawImageOptions{}
	buttonRIGHTOptions.GeoM.Translate(-float64(ArrowPressedImage.Bounds().Dx())/2, -float64(ArrowPressedImage.Bounds().Dy())/2)
	buttonRIGHTOptions.GeoM.Rotate(math.Pi / 2)
	buttonRIGHTOptions.GeoM.Translate(float64(ArrowPressedImage.Bounds().Dx())/2, float64(ArrowPressedImage.Bounds().Dy())/2)
	buttonRIGHTOptions.GeoM.Translate(215+xOffset, 231+yOffset)

	for _, button := range []struct {
		button  controller.Button
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
}

func (e *Emulator) DrawOverlaySprites(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 20)
	op.GeoM.Scale(4, 4)
	screen.DrawImage(e.PPU.DrawOAMSprites(), op)
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
