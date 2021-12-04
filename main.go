package main

import (
	"github.com/exp625/gones/pkg/emulator"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"os"
)

func main() {
	romFile := "test/nestest.nes"
	if len(os.Args) == 2 {
		romFile = os.Args[1]
	}

	e, err := emulator.New(romFile, true)
	if err != nil {
		log.Fatal("failed to set up emulator: ", err)
	}

	ebiten.SetWindowTitle("GoNES")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(emulator.WindowWidth, emulator.WindowHeight)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
