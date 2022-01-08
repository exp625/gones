package emulator

import (
	"encoding/hex"
	"github.com/exp625/gones/internal/plz"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (e *Emulator) SaveGame() {

	saveName := time.Now().Format("saves/2006-01-02_15-04-05.") + hex.EncodeToString(e.Cartridge.Identifier[:]) + ".save"
	ensureSaveDir(saveName)
	fo, err := os.Create(saveName)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		return
	}
	_, err = fo.Write(e.Cartridge.Save())
	if err != nil {
		log.Fatal("error saving", err)
		return
	}
	log.Println("Game saved")
	defer plz.Close(fo)
}

func (e *Emulator) LoadGame() {
	entries, err := os.ReadDir("saves/")
	if err != nil {
		entries = make([]os.DirEntry, 0)
	}
	var newestSave os.DirEntry
	var newestTime time.Time
	for _, entry := range entries {
		nameArr := strings.Split(entry.Name(), ".")
		if nameArr[1] == hex.EncodeToString(e.Cartridge.Identifier[:]) {
			entryTime, err := time.Parse("2006-01-02_15-04-05", nameArr[0])
			if err != nil {
				continue
			}
			if newestTime.Before(entryTime) {
				newestSave = entry
				newestTime = entryTime
			}
		}
	}
	if newestSave != nil {
		saveFileBytes, err := os.ReadFile("saves/" + newestSave.Name())
		if err != nil {
			log.Print("error opening savefile", err)
			return
		}
		e.Cartridge.Load(saveFileBytes)
		log.Println("Found save file from ", newestTime.String())
	}

}

func ensureSaveDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			log.Fatal(merr)
		}
	}
}
