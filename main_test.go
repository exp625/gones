package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/exp625/gones/nes"
	"github.com/exp625/gones/nes/cartridge"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCPULogOutput(t *testing.T) {
	b, err := ioutil.ReadFile("test/nestest.nes")
	if err != nil {
		t.Fatal(err)
	}

	cat := cartridge.LoadCartridge(b)

	logFile, err := os.Create(time.Now().Format("log/2006-01-02_15-04-05_nes.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()

	emulator := &Emulator{NES: nes.New(NESClockTime, NESSampleTime)}
	emulator.loggingEnabled = true
	emulator.InsertCartridge(cat)
	emulator.Reset()
	emulator.CPU.PC = 0xC000
	emulator.CPU.P = 0x24
	emulator.Logger = logFile

	sr := beep.SampleRate(AudioSampleRate)
	err = speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		panic(err)
	}
	defer speaker.Close()
	speaker.Play(Audio(emulator))

	for i := 0; i < 8992; i++ {
		emulator.Clock()
		emulator.Clock()
		emulator.Clock()
		for emulator.CPU.CycleCount != 0 {
			emulator.Clock()
			emulator.Clock()
			emulator.Clock()
		}
	}

	sf, err := os.Open(logFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer sf.Close()

	df, err := os.Open("test/nestest.log")
	if err != nil {
		log.Fatal(err)
	}
	defer df.Close()

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	line := 1
	previousLine := ""
	previousExpectedLine := ""

	for sscan.Scan() {
		currentLine := string(sscan.Bytes())
		dscan.Scan()
		currentExpectedLine := string(dscan.Bytes())
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			ll := len(currentLine)
			if len(currentExpectedLine) < ll {
				ll = len(currentExpectedLine)
			}
			diffs := make([]bool, ll)
			for i := 0; i < ll; i++ {
				cl := " "
				if i < len(currentLine) {
					cl = string(currentLine[i])
				}

				cel := " "
				if i < len(currentExpectedLine) {
					cel = string(currentExpectedLine[i])
				}
				diffs[i] = cl == cel
			}

			errorMessage := strings.Builder{}

			errorMessage.WriteString(fmt.Sprintf("\nEncountered invalid line: (%s:%d)\n\n", logFile.Name(), line))

			errorMessage.WriteString(fmt.Sprintf("%s\n", previousLine))
			errorMessage.WriteString(fmt.Sprintf("%s\n", currentLine))
			for _, d := range diffs {
				if d {
					errorMessage.WriteString(" ")
				} else {
					errorMessage.WriteString("^")
				}
			}
			errorMessage.WriteString("\n\n")

			errorMessage.WriteString(fmt.Sprintf("%s\n", previousExpectedLine))
			errorMessage.WriteString(fmt.Sprintf("%s\n", currentExpectedLine))
			t.Fatal(errorMessage.String())
		}
		line++
		previousLine = currentLine
		previousExpectedLine = currentExpectedLine
	}
}
