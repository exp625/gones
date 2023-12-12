package emulator

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

type HeaderEntry struct {
	Key    string
	Text   string
	Screen Screen
}

var (
	debuggerHeaderEntries = []HeaderEntry{
		{"ECS", "Close", -1},
		{"F1", "CPU", OverlayCPU},
		{"F2", "PPU", OverlayPPU},
		{"F3", "APU", OverlayAPU},
		{"F3", "Nametable", OverlayNametables},
		{"F4", "Palette", OverlayPalettes},
		{"F5", "Controller", OverlayControllers},
		{"F6", "Sprites", OverlaySprites},
	}
	settingsHeaderEntries = []HeaderEntry{
		{"ECS", "Close", -1},
		{"F1", "ROM Chooser", SettingROMChooser},
		{"F2", "Save/Load Game", SettingsSave},
		{"F3", "Key Bindings", SettingKeybindings},
		{"F4", "Audio", SettingAudio},
	}
)

func (e *Emulator) DrawHeader(entries []HeaderEntry, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	width, _ := ebiten.WindowSize()

	header := ebiten.NewImage(width, 20)
	header.Fill(color.White)
	screen.DrawImage(header, op)

	bg := ebiten.NewImage(width/len(entries), 20)
	bg.Fill(color.RGBA{B: 255, A: 255})

	for i, entry := range entries {
		headerEntry := textutil.New(basicfont.Face7x13, width, 20, (width/len(entries))*i, 3, 1)
		if e.ActiveScreen == entry.Screen {
			screen.DrawImage(bg, op)
			headerEntry.Color(color.White)
		} else {
			headerEntry.Color(color.Black)
		}
		op.GeoM.Translate(float64(width/(len(entries))), 0)

		plz.Just(fmt.Fprintf(headerEntry, " "+entry.Key+": "+entry.Text))
		headerEntry.Draw(screen)
	}
}
