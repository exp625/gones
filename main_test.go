package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/pkg/cpu"
	"github.com/exp625/gones/pkg/emulator"
	"github.com/exp625/gones/pkg/logger"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// TestCPULogOutput runs the CPU Output against the nestest.log
func TestCPULogOutput(t *testing.T) {
	logFileName := time.Now().Format("log/2006-01-02_15-04-05_nes.log")
	nesFileName := "test/nestest.nes"

	e, err := emulator.New(nesFileName, false)
	if err != nil {
		t.Fatal(err)
	}
	e.CPU.PC = 0xC000
	e.CPU.P = cpu.StatusRegister(cpu.FlagUnused | cpu.FlagInterruptDisable)

	fileLogger := &logger.FileLogger{}
	fileLogger.Name = logFileName
	e.Logger = fileLogger
	e.Logger.StartLogging()
	logFile, err := os.Create(logFileName)

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
	e.Close()
}

// TestVRAMAccess test PPU palette RAM read/write and mirroring test
func TestPPUPaletteRam(t *testing.T) {
	nesFileName := "test/palette_ram.nes"
	e, err := emulator.New(nesFileName, false)
	if err != nil {
		t.Fatal(err)
	}
	// Run emulator until the infinite loop is reached
	for e.CPU.PC != 0xE412 {
		e.Clock()
	}
	// Get the result
	result := e.CPURead(0x00F0)
	switch result {
	case 1:
		// Test passed
		return
	case 2:
		t.Fatal("Palette read shouldn't be buffered like other VRAM")
	case 3:
		t.Fatal("Palette write/read doesn't work")
	case 4:
		t.Fatal("Palette should be mirrored within $3f00-$3fff")
	case 5:
		t.Fatal("Write to $10 should be mirrored at $00")
	case 6:
		t.Fatal("Write to $00 should be mirrored at $10")
	default:
		t.Fatal("TestPPUPaletteRam failed")
	}
	e.Close()
}

// TestSpriteRam tests sprite RAM access via $2003, $2004, and $4014
func TestSpriteRam(t *testing.T) {
	nesFileName := "test/sprite_ram.nes"
	e, err := emulator.New(nesFileName, false)
	if err != nil {
		t.Fatal(err)
	}
	// Run emulator until the infinite loop is reached
	for e.CPU.PC != 0xE467 {
		e.Clock()
	}
	// Get the result
	result := e.CPURead(0x00F0)
	switch result {
	case 1:
		// Test passed
		return
	case 2:
		t.Fatal("Basic read/write doesn't work")
	case 3:
		t.Fatal("Address should increment on $2004 write")
	case 4:
		t.Fatal("Address should not increment on $2004 read")
	case 5:
		t.Fatal("Third sprite bytes should be masked with $e3 on read ")
	case 6:
		t.Fatal("$4014 DMA copy doesn't work at all")
	case 7:
		t.Fatal("$4014 DMA copy should start at value in $2003 and wrap")
	case 8:
		t.Fatal("$4014 DMA copy should leave value in $2003 intact")
	default:
		t.Fatal("TestVRAMAccess failed")
	}
	e.Close()
}

// TestVBLFlag tests if the VBL flag ($2002.7) is cleared by the PPU around 2270 CPU clocks after NMI occurs.
func TestVBLFlag(t *testing.T) {
	nesFileName := "test/vbl_clear_time.nes"
	e, err := emulator.New(nesFileName, false)
	if err != nil {
		t.Fatal(err)
	}
	// Run emulator until the infinite loop is reached
	for e.CPU.PC != 0xE3B3 {
		e.Clock()
	}
	// Get the result
	result := e.CPURead(0x00F0)
	switch result {
	case 1:
		// Test passed
		return
	case 2:
		t.Fatal(" VBL flag cleared too soon")
	case 3:
		t.Fatal("VBL flag cleared too late")
	default:
		t.Fatal("TestVRAMAccess failed")
	}
	e.Close()
}

// TestVRAMAccess tests PPU VRAM read/write and internal read buffer operation
func TestVRAMAccess(t *testing.T) {
	nesFileName := "test/vram_access.nes"
	e, err := emulator.New(nesFileName, false)
	if err != nil {
		t.Fatal(err)
	}
	// Run emulator until the infinite loop is reached
	for e.CPU.PC != 0xE48D {
		e.Clock()
	}
	// Get the result
	result := e.CPURead(0x00F0)
	switch result {
	case 1:
		// Test passed
		return
	case 2:
		t.Fatal("VRAM reads should be delayed in a buffer")
	case 3:
		t.Fatal("Basic Write/read doesn't work")
	case 4:
		t.Fatal("Read buffer shouldn't be affected by VRAM write")
	case 5:
		t.Fatal("Read buffer shouldn't be affected by palette write")
	case 6:
		t.Fatal("Palette read should also read VRAM into read buffer")
	case 7:
		t.Fatal("Shadow VRAM read unaffected by palette transparent color mirroring")
	default:
		t.Fatal("TestVRAMAccess failed")
	}
	e.Close()
}
