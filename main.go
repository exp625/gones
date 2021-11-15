package main

import (
	"fmt"
	"github.com/exp625/gones/nes"
	"github.com/exp625/gones/nes/cartridge"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	AudioSampleRate = 44100
	PPUFrequency    = 5369318.0
	NESSampleTime   = 1.0 / AudioSampleRate
	NESClockTime    = 1.0 / PPUFrequency
)

// Start the main thread
func main() {
	pixelgl.Run(run)
}

// Emulator struct
type Emulator struct {
	*nes.NES
	autoRun                   bool
	hideDebug                 bool
	hideInfo                  bool
	hidePatternTables         bool
	displayRamPC              bool
	autoRunCycles             int
	nanoSecondsSpentInAutoRun time.Duration
	autoRunStarted            time.Time
}

const (
	Width  = 1200
	Height = 1000
)

func run() {
	// Load Cartridge
	argsWithoutProg := os.Args[1:]
	if argsWithoutProg[0] == "" {
		log.Panic("No rom file provided")
	}
	bytes, err := ioutil.ReadFile(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	cat := cartridge.LoadCartridge(bytes)
	// Create Window
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, Width, Height),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create text atlas
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	cpuText := text.New(pixel.V(0, Height-atlas.LineHeight()*2), atlas)
	codeText := text.New(pixel.V(0, Height-atlas.LineHeight()*2-200), atlas)
	cardridgeText := text.New(pixel.V(800, Height-atlas.LineHeight()*2-200), atlas)
	zeroPageText := text.New(pixel.V(0, Height-atlas.LineHeight()*2-370), atlas)
	stackText := text.New(pixel.V(400, Height-atlas.LineHeight()*2-370), atlas)
	ramText := text.New(pixel.V(0, Height-atlas.LineHeight()*2-620), atlas)

	//Create NES
	emulator := &Emulator{NES: nes.New(NESClockTime, NESSampleTime)}
	emulator.hidePatternTables = true

	emulator.InsertCartridge(cat)


	emulator.Reset()

	// Setup sound
	sr := beep.SampleRate(AudioSampleRate)
	err = speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		panic(err)
	}
	defer speaker.Close()
	speaker.Play(Audio(emulator))

	// Render Loop
	for !win.Closed() {
		win.Clear(colornames.Black)

		handleInput(win, emulator)

		if !emulator.hideInfo {
			cpuText.Clear()
			DrawCPU(cpuText, emulator)
			cpuText.Draw(win, pixel.IM.Scaled(cpuText.Orig, 2))
		}

		if !emulator.hideDebug {
			codeText.Clear()
			zeroPageText.Clear()
			stackText.Clear()
			ramText.Clear()
			cardridgeText.Clear()

			DrawCode(codeText, emulator)
			DrawZeroPage(zeroPageText, emulator)
			DrawStack(stackText, emulator)
			DrawRAM(ramText, emulator)
			emulator.Bus.Cartridge.DebugDisplay(cardridgeText)

			moved := pixel.IM
			if emulator.hideInfo {
				moved = moved.Moved(pixel.V(0, 200))
			}
			codeText.Draw(win, pixel.IM.Scaled(codeText.Orig, 2).Chained(moved))
			cardridgeText.Draw(win, pixel.IM.Scaled(cardridgeText.Orig, 2).Chained(moved))
			zeroPageText.Draw(win, moved)
			stackText.Draw(win, moved)
			ramText.Draw(win, moved)
		}

		if !emulator.hidePatternTables {
			DrawCHRROM(emulator, 0).Draw(win, pixel.IM.Moved(pixel.V(256 + 5, 256 + 5)).Scaled(pixel.V(256 + 5, 256 + 5), 4))
			DrawCHRROM(emulator, 1).Draw(win, pixel.IM.Moved(pixel.V(256*3 + 10, 256 + 5)).Scaled(pixel.V(256*3 + 10, 256 + 5), 4))
		}



		// Update Frame
		win.Update()
	}
}

func handleInput(win *pixelgl.Window, emulator *Emulator) {
	if emulator.autoRun {
		emulator.nanoSecondsSpentInAutoRun += time.Now().Sub(emulator.autoRunStarted)
	}
	emulator.autoRunStarted = time.Now()

	// Space will toggle the auto run mode
	if win.JustPressed(pixelgl.KeySpace) {
		emulator.autoRun = !emulator.autoRun
	}

	// D Key toggles display of hideDebug hideInfo
	if win.JustPressed(pixelgl.KeyD) {
		emulator.hideDebug = !emulator.hideDebug
	}

	// I Key toggles display of hideInfo
	if win.JustPressed(pixelgl.KeyI) {
		emulator.hideInfo = !emulator.hideInfo
	}

	// P Key toggles display of patternTables
	if win.JustPressed(pixelgl.KeyP) {
		emulator.hidePatternTables = !emulator.hidePatternTables
	}

	// Right Arrow Key issues one Master Clock
	if win.JustPressed(pixelgl.KeyRight) && !emulator.autoRun {
		emulator.Clock()
	}

	// Toggle rom pc
	if win.JustPressed(pixelgl.KeyX) && !emulator.autoRun {
		emulator.displayRamPC = !emulator.displayRamPC
	}

	// Up Arrow Key issues three Master Clocks
	if win.JustPressed(pixelgl.KeyUp) && !emulator.autoRun {
		emulator.Clock()
		emulator.Clock()
		emulator.Clock()
	}

	// Enter Key one CPU instruction
	if win.JustPressed(pixelgl.KeyEnter) && !emulator.autoRun {
		emulator.Clock()
		emulator.Clock()
		emulator.Clock()
		for emulator.NES.Bus.CPU.CycleCount != 0 {
			emulator.Clock()
			emulator.Clock()
			emulator.Clock()
		}
	}

	// R Key will reset the emulator
	if win.JustPressed(pixelgl.KeyR) {
		emulator.Reset()
	}
}

func intbool(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

// Audio Streamer
func Audio(emulator *Emulator) beep.Streamer {
	// The function gets called if the audio hardware request new audio samples. The length of the sample array indicates how many sample are requested.
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// If the emulator is set to auto run: Run the emulation until the time of one audio sample passed.
			if emulator.autoRun {
				for !emulator.Clock() {
					emulator.autoRunCycles++
				}

				// Get the audio sample for the APU
				sample := emulator.Bus.APU.GetAudioSample()
				samples[i][0] = sample
				samples[i][1] = sample
			} else {
				// No sound when auto run is false
				samples[i] = [2]float64{}
			}
		}
		return len(samples), true
	})
}

func DrawCPU(statusText *text.Text, emulator *Emulator) {
	fmt.Fprintf(statusText, "Auto Run Mode: \t %t\n", emulator.autoRun)
	fmt.Fprintf(statusText, "Master Clock Count: \t %d\n", emulator.NES.MasterClockCount)
	fmt.Fprintf(statusText, "CPU Clock Count: \t %d\n", emulator.NES.Bus.CPU.ClockCount)
	fmt.Fprintf(statusText, "Clock Cycles Per Second (during auto run): %0.2f/s\n",
		1000*1000*1000*float64(emulator.autoRunCycles)/(float64(emulator.nanoSecondsSpentInAutoRun)),
	)
	fmt.Fprint(statusText, "\n")

	// Print Status
	fmt.Fprint(statusText, "CPU STATUS \t")

	arr := "CZIDB-VN"
	for i := 0; i < 8; i++ {
		if emulator.Bus.CPU.GetFlag(1 << i) {
			statusText.Color = colornames.Green
		} else {
			statusText.Color = colornames.Red
		}
		fmt.Fprint(statusText, string(arr[i]))
	}
	statusText.Color = colornames.White
	fmt.Fprint(statusText, "\n")
	fmt.Fprintf(statusText, "PC: 0x%02X\t", emulator.NES.Bus.CPU.PC)
	fmt.Fprintf(statusText, "A: 0x%02X\t", emulator.NES.Bus.CPU.A)
	fmt.Fprintf(statusText, "X: 0x%02X\t", emulator.NES.Bus.CPU.X)
	fmt.Fprintf(statusText, "Y: 0x%02X\t", emulator.NES.Bus.CPU.Y)
	fmt.Fprintf(statusText, "S: 0x%02X\t\n", emulator.NES.Bus.CPU.S)
}

func DrawCode(statusText *text.Text, emulator *Emulator) {
	offset := uint16(0)
	for j := 0; j <= 5; j++ {
		if j == 0 {
			statusText.Color = colornames.Yellow
		}
		if emulator.displayRamPC {
			fmt.Fprintf(statusText, "0x%04X ", (emulator.Bus.CPU.CurrentPC+offset - 0x8000) % 0x4000 * uint16(emulator.Bus.Cartridge.PrgRomSize) + 0x0010)
		} else {
			fmt.Fprintf(statusText, "0x%04X ", emulator.Bus.CPU.CurrentPC+offset)
		}

		fmt.Fprintf(statusText, "0x%02X ", emulator.Bus.CPURead(emulator.Bus.CPU.CurrentPC+offset))
		fmt.Fprint(statusText, "[",nes.OpCodeMap[emulator.Bus.CPURead(emulator.Bus.CPU.CurrentPC+offset)], "] ")
		inst := nes.Instructions[emulator.Bus.CPURead(emulator.Bus.CPU.CurrentPC+offset)]
		if inst.Length != 0 {
			addr, _ := inst.AddressMode()
			if j == 0 {
				fmt.Fprintf(statusText, "(0x%04X) ", addr)
			} else {
				fmt.Fprint(statusText, "         ")
			}
			for i := 1; i < int(inst.Length); i++ {
				fmt.Fprintf(statusText, "%02X ", emulator.Bus.CPURead(emulator.Bus.CPU.CurrentPC+offset+uint16(i)))
			}
			statusText.Color = colornames.White
			offset += uint16(inst.Length)
			fmt.Fprint(statusText, "\n")
		}
	}
}

func DrawRAM(statusText *text.Text, emulator *Emulator) {
	statusText.Color = colornames.White
	fmt.Fprint(statusText, "Ram Content:\n     ")
	statusText.Color = colornames.Yellow
	for i := 0; i <= 0xF; i++ {
		fmt.Fprintf(statusText, "%02X ", uint16(i))
	}

	for x := 0x0200; x <= 0x07FF; x += 0x10 {
		// Check if this "row" of memory has anything other than 0x00 in it
		var hasContent bool
		for y := 0; y <= 15; y++ {
			if emulator.Bus.CPURead(uint16(x+y)) != 0x00 {
				hasContent = true
				break
			}
		}
		// Display the "row" of memory iuf
		if hasContent {
			statusText.Color = colornames.Yellow
			fmt.Fprintf(statusText, "\n%04X ", uint16(x&0xFFF0))
			statusText.Color = colornames.White
			for y := 0; y <= 15; y++ {
				fmt.Fprintf(statusText, "%02X ", emulator.Bus.CPURead(uint16(x+y)))
			}
		}
	}
}

func DrawZeroPage(statusText *text.Text, emulator *Emulator) {
	statusText.Color = colornames.White
	fmt.Fprint(statusText, "Zero Page:\n     ")
	statusText.Color = colornames.Yellow
	for i := 0; i <= 0xF; i++ {
		fmt.Fprintf(statusText, "%02X ", uint16(i))
	}

	for i := 0x0000; i <= 0x00FF; i++ {
		if i%16 == 0 {
			statusText.Color = colornames.Yellow
			fmt.Fprintf(statusText, "\n%04X ", uint16(i&0xFFF0))
		}
		statusText.Color = colornames.White
		fmt.Fprintf(statusText, "%02X ", emulator.Bus.CPURead(uint16(i)))
	}
}

func DrawStack(statusText *text.Text, emulator *Emulator) {
	statusText.Color = colornames.White
	fmt.Fprintf(statusText, "Stack: 0x%02X\n     ", emulator.NES.Bus.CPU.S)
	statusText.Color = colornames.Yellow
	for i := 0; i <= 0xF; i++ {
		fmt.Fprintf(statusText, "%02X ", uint16(i))
	}

	for i := 0x0100; i <= 0x01FF; i++ {
		if i%16 == 0 {
			statusText.Color = colornames.Yellow
			fmt.Fprintf(statusText, "\n%04X ", uint16(i&0xFFF0))
		}
		if emulator.Bus.CPU.S == uint8(i) {
			statusText.Color = colornames.Green
		} else {
			statusText.Color = colornames.White
		}
		fmt.Fprintf(statusText, "%02X ", emulator.Bus.CPURead(uint16(i)))
	}

}

func DrawCHRROM (emulator *Emulator, table int) *pixel.Sprite{
	width := 128
	height := 128

	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	//DCBA98 76543210
	//---------------
	//	0HRRRR CCCCPTTT
	//  |||||| |||||+++- T: Fine Y offset, the row number within a tile
	//  |||||| ||||+---- P: Bit plane (0: "lower"; 1: "upper")
	//  |||||| ++++----- C: Tile column
	//  ||++++---------- R: Tile row
	//  |+-------------- H: Half of sprite table (0: "left"; 1: "right")
	//  +--------------- 0: Pattern table is at $0000-$1FFF

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
				for tileY := 0; tileY < 8; tileY++ {

					addressPlane0 := uint16(table<<12 | y<<8 | x<<4 | 0 << 3 | tileY)
					addressPlane1 := uint16(table<<12 | y<<8 | x<<4 | 1 << 3 | tileY)
					plane0 := emulator.Bus.PPURead(addressPlane0)
					plane1 := emulator.Bus.PPURead(addressPlane1)

					for tileX := 0; tileX < 8; tileX++ {

						if (plane0 >> (7 - tileX)) & 0x01 == 1 && (plane1 >> (7 - tileX)) & 0x01 ==1 {
							img.Set(x * 8 + tileX, y * 8 + tileY, color.White)
						} else if (plane1 >> (7 - tileX)) & 0x01 ==1 {
							img.Set(x * 8 + tileX, y * 8 + tileY, color.Gray16{0xAAAA})
						} else if (plane0 >> (7 - tileX)) & 0x01 ==1 {
							img.Set(x * 8 + tileX, y * 8 + tileY, color.Gray16{0x5555})
						} else {
							img.Set(x * 8 + tileX, y * 8 + tileY, color.Gray16{0x1111})
						}

					}
				}
			}
		}
	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds())

}
