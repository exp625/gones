package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/pkg/cpu"
	"github.com/exp625/gones/pkg/emulator"
	"github.com/exp625/gones/pkg/logger"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type TestConfig struct {
	Tests []struct {
		Rom     string `json:"rom"`
		Frames  int    `json:"frames"`
		Output  string `json:"output"`
		Results []struct {
			Code    int    `json:"code"`
			Pass    bool   `json:"pass"`
			Message string `json:"message"`
		} `json:"results"`
	} `json:"tests"`
}

// TestCPULogOutput runs the CPU Output against the nestest.log
func TestCPULogOutput(t *testing.T) {
	logFileName := time.Now().Format("log/2006-01-02_15-04-05_nes.log")
	nesFileName := "test/nes-test-roms/other/nestest.nes"

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

	expectedLogFile, err := os.Open("test/nes-test-roms/other/nestest.log")
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

func TestFromConfig(t *testing.T) {
	var config TestConfig
	configFile, err := os.Open("./test/test.json")
	if err != nil {
		t.Fatal(err)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &config)

	for _, test := range config.Tests {
		rom := strings.Split(test.Rom, "/")
		testname := fmt.Sprintf("Test %s", strings.Split(rom[len(rom)-1], ".")[0])
		t.Run(testname, func(t *testing.T) {
			e, err := emulator.New("./test/"+test.Rom, false)
			if err != nil {
				t.Fatal(err)
			}

			// Run emulator until the infinite loop is reached
			for e.PPU.FrameCount != uint64(test.Frames) {
				e.Clock()
			}

			// Get the result
			numberStr := strings.Replace(test.Output, "0x", "", -1)
			numberStr = strings.Replace(numberStr, "0X", "", -1)
			resultLocation, err := strconv.ParseUint(numberStr, 16, 64)
			if err != nil {
				t.Fatal(err)
			}

			resultValue := e.CPURead(uint16(resultLocation))

			for _, result := range test.Results {
				if resultValue == uint8(result.Code) {
					if result.Pass {
						// Test passed
						return
					} else {
						t.Fatal(result.Message)
						return
					}
				}
			}
			t.Fatal("Could not match result for test")

		})
	}
}
