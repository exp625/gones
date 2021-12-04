package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

var (
	headerEntries = [6][2]string{
		{"F1", "CPU Debug"},
		{"F2", "PPU Debug"},
		{"F3", "Nametable Debug"},
		{"F4", "Palette Debug"},
		{"F5", "Controller Debug"},
		{"F6", "Keybindings"},
	}
)

func (e *Emulator) DrawHeader(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	width, _ := ebiten.WindowSize()

	header := ebiten.NewImage(width, 20)
	header.Fill(color.White)
	screen.DrawImage(header, op)

	bg := ebiten.NewImage(width/len(headerEntries), 20)
	bg.Fill(color.RGBA{B: 255, A: 255})

	for i := range headerEntries {
		headerEntry := textutil.New(basicfont.Face7x13, width, 20, (width/len(headerEntries))*i, 3, 1)
		if e.ActiveOverlay == Overlay(i+1) {
			screen.DrawImage(bg, op)
			headerEntry.Color(color.White)
		} else {
			headerEntry.Color(color.Black)
		}
		op.GeoM.Translate(float64(width/(len(headerEntries))), 0)

		plz.Just(fmt.Fprintf(headerEntry, " "+headerEntries[i][0]+": "+headerEntries[i][1]))
		headerEntry.Draw(screen)
	}
}
