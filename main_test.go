package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/exp625/gones/pkg/cpu"
	"github.com/exp625/gones/pkg/emulator"
	"github.com/exp625/gones/pkg/plz"

	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCPULogOutput(t *testing.T) {
	logFileName := time.Now().Format("log/2006-01-02_15-04-05_nes.log")
	nesFileName := "test/nestest.nes"

	e, err := emulator.New(nesFileName, false)
	if err != nil {
		t.Fatal(err)
	}
	e.CPU.PC = 0xC000
	e.CPU.P = cpu.FlagUnused | cpu.FlagInterruptDisable

	e.EnsureLogDir(logFileName)
	logFile, err := os.Create(logFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer plz.Close(logFile)

	// Run same number of instructions as in nestest.log
	for i := 0; i < 8992; i++ {
		e.Clock()
		e.Clock()
		e.Clock()
		for e.CPU.CycleCount != 0 {
			e.Clock()
			e.Clock()
			e.Clock()
		}
	}

	gotLogFile, err := os.Open(logFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer plz.Close(gotLogFile)

	expectedLogFile, err := os.Open("test/nestest.log")
	if err != nil {
		log.Fatal(err)
	}
	defer plz.Close(expectedLogFile)

	gotLogFileScanner := bufio.NewScanner(gotLogFile)
	expectedLogFileScanner := bufio.NewScanner(expectedLogFile)

	line := 1
	previouslyGotLine := ""
	previouslyExpectedLine := ""

	for gotLogFileScanner.Scan() {
		expectedLogFileScanner.Scan()

		gotLine := string(gotLogFileScanner.Bytes())
		expectedLine := string(expectedLogFileScanner.Bytes())

		if !bytes.Equal(gotLogFileScanner.Bytes(), expectedLogFileScanner.Bytes()) {
			lineLength := len(gotLine)
			if len(expectedLine) < lineLength {
				lineLength = len(expectedLine)
			}

			differences := make([]bool, lineLength)
			for i := 0; i < lineLength; i++ {
				gotCharacter := " "
				if i < len(gotLine) {
					gotCharacter = string(gotLine[i])
				}

				expectedCharacter := " "
				if i < len(expectedLine) {
					expectedCharacter = string(expectedLine[i])
				}

				differences[i] = gotCharacter == expectedCharacter
			}

			errorMessage := strings.Builder{}

			errorMessage.WriteString(fmt.Sprintf("\nEncountered invalid line: (%s:%d)\n\n", logFile.Name(), line))

			errorMessage.WriteString(fmt.Sprintf("%s\n", previouslyGotLine))
			errorMessage.WriteString(fmt.Sprintf("%s\n", gotLine))

			for _, d := range differences {
				if d {
					errorMessage.WriteString(" ")
				} else {
					errorMessage.WriteString("^")
				}
			}
			errorMessage.WriteString("\n\n")

			errorMessage.WriteString(fmt.Sprintf("%s\n", previouslyExpectedLine))
			errorMessage.WriteString(fmt.Sprintf("%s\n", expectedLine))

			t.Fatal(errorMessage.String())
		}
		line++
		previouslyGotLine = gotLine
		previouslyExpectedLine = expectedLine
	}
}
