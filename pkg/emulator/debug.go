package emulator

import (
	"embed"
	"fmt"
	"github.com/exp625/gones/internal/textutil"
	"github.com/exp625/gones/pkg/plz"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	_ "image/png"
)

//go:embed resources/*
var resourcesFS embed.FS

var (
	ControllerImage    *ebiten.Image
	ArrowPressedImage  *ebiten.Image
	PillPressedImage   *ebiten.Image
	CirclePressedImage *ebiten.Image
)

func init() {
	controllerImageReader, _ := resourcesFS.Open("resources/controller.png")
	arrowPressedImageReader, _ := resourcesFS.Open("resources/arrow_pressed.png")
	pillPressedImageReader, _ := resourcesFS.Open("resources/pill_pressed.png")
	circlePressedImageReader, _ := resourcesFS.Open("resources/circle_pressed.png")

	ControllerImage, _, _ = ebitenutil.NewImageFromReader(controllerImageReader)
	ArrowPressedImage, _, _ = ebitenutil.NewImageFromReader(arrowPressedImageReader)
	PillPressedImage, _, _ = ebitenutil.NewImageFromReader(pillPressedImageReader)
	CirclePressedImage, _, _ = ebitenutil.NewImageFromReader(circlePressedImageReader)
}

func (e *Emulator) DrawCPU(t *textutil.Text) {
	plz.Just(fmt.Fprintf(t, "FPS: %0.2f \t Auto Run Mode: \t %t \t Logging Enabled: \t %t \n", ebiten.CurrentFPS(), e.autoRunEnabled, e.LoggingEnabled))
	plz.Just(fmt.Fprintf(t, "Master Clock Count: \t %d\n", e.NES.MasterClockCount))
	plz.Just(fmt.Fprintf(t, "CPU Clock Count: \t %d \t Requested: \t %d \n", e.CPU.ClockCount, e.requestedSteps))
	plz.Just(fmt.Fprintf(t, "Clock Cycles Per Second (during auto run): %0.2f/s\n",
		1000*1000*1000*float64(e.autoRunCycles)/(float64(e.nanoSecondsSpentInAutoRun))))
	plz.Just(fmt.Fprint(t, "\n"))

	// Print Status
	plz.Just(fmt.Fprint(t, "CPU STATUS \t"))

	arr := "CZIDB-VN"
	for i := 0; i < 8; i++ {
		if e.CPU.Get(1 << i) {
			t.Color(colornames.Green)
		} else {
			t.Color(colornames.Red)
		}
		plz.Just(fmt.Fprint(t, string(arr[i])))
	}
	t.Color(colornames.White)
	plz.Just(fmt.Fprintf(t, "\t%02X", e.CPU.P))
	plz.Just(fmt.Fprint(t, "\n"))
	plz.Just(fmt.Fprintf(t, "PC: 0x%02X\t", e.CPU.PC))
	plz.Just(fmt.Fprintf(t, "A: 0x%02X\t", e.CPU.A))
	plz.Just(fmt.Fprintf(t, "X: 0x%02X\t", e.CPU.X))
	plz.Just(fmt.Fprintf(t, "Y: 0x%02X\t", e.CPU.Y))
	plz.Just(fmt.Fprintf(t, "S: 0x%02X\t\n", e.CPU.S))
}

func (e *Emulator) DrawInstructions(t *textutil.Text) {
	offset := uint16(0)
	if e.CPU.CycleCount < 0 {
		plz.Just(fmt.Fprint(t, "ERR"))
		return
	}
	for j := 0; j <= 5; j++ {
		if j == 0 {
			t.Color(colornames.Yellow)
		}
		plz.Just(fmt.Fprintf(t, "%04X ", e.CPU.PC+offset))
		inst := e.CPU.Instructions[e.CPURead(e.CPU.PC+offset)]
		i := 0
		for ; i < int(inst.Length); i++ {
			plz.Just(fmt.Fprintf(t, "%02X ", e.CPURead(e.CPU.PC+offset+uint16(i))))
		}
		for ; i < 3; i++ {
			plz.Just(fmt.Fprint(t, "   "))
		}
		plz.Just(fmt.Fprint(t, "", e.CPU.Mnemonics[e.CPURead(e.CPU.PC+offset)], " "))

		if inst.Length != 0 {
			addr, data, _ := inst.AddressMode()
			if j == 0 {
				// Display Address
				switch e.CPU.Mnemonics[e.CPURead(e.CPU.PC)][1] {
				case "REL":
					plz.Just(fmt.Fprintf(t, "$%04X", addr))
				case "ABS":
					if addr <= 0x1FFF {
						plz.Just(fmt.Fprintf(t, "$%04X = %02X                  ", addr, data))
					} else {
						plz.Just(fmt.Fprintf(t, "$%04X                       ", addr))
					}
				case "ACC":
					plz.Just(fmt.Fprint(t, "A"))
				case "IMM":
					plz.Just(fmt.Fprintf(t, "#$%02X", data))
				case "ZPX":
					plz.Just(fmt.Fprintf(t, "$%02X,X @ %02X = %02X", e.CPURead(e.CPU.PC+1), addr, data))
				case "ZPY":
					plz.Just(fmt.Fprintf(t, "$%02X,Y @ %02X = %02X", e.CPURead(e.CPU.PC+1), addr, data))
				case "ZP0":
					plz.Just(fmt.Fprintf(t, "$%02X = %02X", addr&0x00FF, data))
				case "IDX":
					// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
					plz.Just(fmt.Fprintf(t, "($%02X,X) @ %02X = %04X = %02X", e.CPURead(e.CPU.PC+1), e.CPURead(e.CPU.PC+1)+e.CPU.X, addr, data))
				case "IZY":
					// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
					plz.Just(fmt.Fprintf(t, "($%02X),Y = %04X @ %04X = %02X", e.CPURead(e.CPU.PC+1), addr-uint16(e.CPU.Y), addr, data))
				case "IND":
					plz.Just(fmt.Fprintf(t, "($%02X%02X) = %04X", e.CPURead(e.CPU.PC+2), e.CPURead(e.CPU.PC+1), addr))
				case "ABX":
					plz.Just(fmt.Fprintf(t, "$%02X%02X,X @ %04X = %02X", e.CPURead(e.CPU.PC+2), e.CPURead(e.CPU.PC+1), addr, data))
				case "ABY":
					plz.Just(fmt.Fprintf(t, "$%02X%02X,Y @ %04X = %02X", e.CPURead(e.CPU.PC+2), e.CPURead(e.CPU.PC+1), addr, data))
				}
			}
			t.Color(colornames.White)
			offset += inst.Length
			plz.Just(fmt.Fprint(t, "\n"))
		}
	}
}

func (e *Emulator) DrawZeroPage(t *textutil.Text) {

	t.Color(colornames.White)
	plz.Just(fmt.Fprint(t, "Zero Page:\n     "))
	t.Color(colornames.Yellow)
	for i := 0; i <= 0xF; i++ {
		plz.Just(fmt.Fprintf(t, "%02X ", uint16(i)))
	}

	for i := 0x0000; i <= 0x00FF; i++ {
		if i%16 == 0 {
			t.Color(colornames.Yellow)
			plz.Just(fmt.Fprintf(t, "\n%04X ", uint16(i&0xFFF0)))
		}
		t.Color(colornames.White)
		plz.Just(fmt.Fprintf(t, "%02X ", e.CPURead(uint16(i))))
	}
}

func (e *Emulator) DrawStack(t *textutil.Text) {

	t.Color(colornames.White)
	plz.Just(fmt.Fprintf(t, "Stack: 0x%02X\n     ", e.CPU.S))
	t.Color(colornames.Yellow)
	for i := 0; i <= 0xF; i++ {
		plz.Just(fmt.Fprintf(t, "%02X ", uint16(i)))
	}

	for i := 0x0100; i <= 0x01FF; i++ {
		if i%16 == 0 {
			t.Color(colornames.Yellow)
			plz.Just(fmt.Fprintf(t, "\n%04X ", uint16(i&0xFFF0)))
		}
		if e.CPU.S == uint8(i) {
			t.Color(colornames.Green)
		} else {
			t.Color(colornames.White)
		}
		plz.Just(fmt.Fprintf(t, "%02X ", e.CPURead(uint16(i))))
	}
}

func (e *Emulator) DrawRAM(t *textutil.Text) {
	t.Color(colornames.White)
	plz.Just(fmt.Fprint(t, "Ram Content:\n     "))
	t.Color(colornames.Yellow)
	for i := 0; i <= 0xF; i++ {
		plz.Just(fmt.Fprintf(t, "%02X ", uint16(i)))
	}

	for x := 0x0200; x <= 0x07FF; x += 0x10 {
		// Check if this "row" of memory has anything other than 0x00 in it
		var hasContent bool
		for y := 0; y <= 15; y++ {
			if e.CPURead(uint16(x+y)) != 0x00 {
				hasContent = true
				break
			}
		}
		// Display the "row" of memory iuf
		if hasContent {
			t.Color(colornames.Yellow)
			plz.Just(fmt.Fprintf(t, "\n%04X ", uint16(x&0xFFF0)))
			t.Color(colornames.White)
			for y := 0; y <= 15; y++ {
				plz.Just(fmt.Fprintf(t, "%02X ", e.CPURead(uint16(x+y))))
			}
		}
	}
}

func (e *Emulator) DrawCartridge(t *textutil.Text) {
	e.Cartridge.Mapper.DebugDisplay(t)
}
