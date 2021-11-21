package main

import (
	"github.com/exp625/gones/pkg/emulator"
	"github.com/faiface/pixel/pixelgl"
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

	pixelgl.Run(e.Run)
}
