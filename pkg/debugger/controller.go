package debugger

import (
	"embed"
	"github.com/exp625/gones/pkg/controller"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"math"
)

//go:embed resources/*
var resourcesFS embed.FS

var (
	ControllerImage    *ebiten.Image
	ArrowPressedImage  *ebiten.Image
	PillPressedImage   *ebiten.Image
	CirclePressedImage *ebiten.Image
)

func init() {
	controllerImageReader, _ := resourcesFS.Open("resources/controller.png")
	arrowPressedImageReader, _ := resourcesFS.Open("resources/arrow_pressed.png")
	pillPressedImageReader, _ := resourcesFS.Open("resources/pill_pressed.png")
	circlePressedImageReader, _ := resourcesFS.Open("resources/circle_pressed.png")

	controllerImageDecoded, _, _ := image.Decode(controllerImageReader)
	arrowPressedImageDecoded, _, _ := image.Decode(arrowPressedImageReader)
	pillPressedImageDecoded, _, _ := image.Decode(pillPressedImageReader)
	circlePressedImageDecoded, _, _ := image.Decode(circlePressedImageReader)

	ControllerImage = ebiten.NewImageFromImage(controllerImageDecoded)
	ArrowPressedImage = ebiten.NewImageFromImage(arrowPressedImageDecoded)
	PillPressedImage = ebiten.NewImageFromImage(pillPressedImageDecoded)
	CirclePressedImage = ebiten.NewImageFromImage(circlePressedImageDecoded)
}


func (nes *Debugger) DrawController(port uint8) *ebiten.Image {

	width, height := ebiten.WindowSize()
	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := ebiten.NewImageFromImage(image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight}))

	xOffset := float64(width/2) - float64(ControllerImage.Bounds().Dx())/2
	yOffset := float64(height)/2 - float64(ControllerImage.Bounds().Dy())/2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(xOffset, yOffset)
	img.DrawImage(ControllerImage, op)

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
		if nes.Controller1.IsPressed(button.button) && port == 1 {
			img.DrawImage(button.image, button.options)
		}
		if nes.Controller2.IsPressed(button.button) && port == 2 {
			img.DrawImage(button.image, button.options)
		}
	}
	return img
}