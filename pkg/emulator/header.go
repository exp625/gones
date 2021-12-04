package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/plz"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

var (
	headerEntries = [6][2]string{
		{"F1", "CPU Debug"},
		{"F2", "PPU Debug"},
		{"F3", "Nametable Debug"},
		{"F4", "Pallet Debug"},
		{"F5", "Controller Debug"},
		{"F6", ""},
	}
	header = ebiten.NewImage(WindowWidth, 20)
	bg     = ebiten.NewImage(WindowWidth/len(headerEntries), 20)
	op     = &ebiten.DrawImageOptions{}
)

func init() {
	header.Fill(color.White)
	bg.Fill(color.RGBA{B: 255, A: 255})
}

func (e *Emulator) DrawHeader(screen *ebiten.Image) {
	op.GeoM.Reset()
	screen.DrawImage(header, op)
	for i := range headerEntries {
		headerEntry := textutil.New(basicfont.Face7x13, WindowWidth, 20, (WindowWidth/len(headerEntries))*i, 3, 1)
		if e.Screen == i+1 {
			screen.DrawImage(bg, op)
			headerEntry.Color(color.White)
		} else {
			headerEntry.Color(color.Black)
		}
		op.GeoM.Translate(float64(WindowWidth/(len(headerEntries))), 0)

		plz.Just(fmt.Fprintf(headerEntry, " "+headerEntries[i][0]+": "+headerEntries[i][1]))
		headerEntry.Draw(screen)
	}
}
