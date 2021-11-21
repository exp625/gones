package emulator

import (
	"github.com/exp625/gones/pkg/plz"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (e *Emulator) StartLogging() {
	name := time.Now().Format("log/2006-01-02_15-04-05_nes.log")
	e.EnsureLogDir(name)
	fo, err := os.Create(name)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		return
	}
	defer plz.Close(fo)

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	e.NES.Logger = f
	log.SetOutput(f)
	e.LoggingEnabled = true
}

func (e *Emulator) StopLogging() {
	e.LoggingEnabled = false
	log.SetOutput(os.Stdout)
	plz.Close(e.NES.Logger)
	e.NES.Logger = nil
}

func (e *Emulator) EnsureLogDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			log.Fatal(merr)
		}
	}
}
