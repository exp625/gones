package main

import (
	"github.com/exp625/gones/pkg/emulator"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"os"
)

//go:generate go run gen/cmd/main.go

func main() {

	romFile := ""
	debug := false
	if len(os.Args) == 2 {
		romFile = os.Args[1]
		// Maybe add flag for debug
	}

	e, err := emulator.New(romFile, debug)
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
