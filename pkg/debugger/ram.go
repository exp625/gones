package debugger

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"golang.org/x/image/colornames"
)

func (nes *Debugger) DrawZeroPage(t *textutil.Text) {
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
		plz.Just(fmt.Fprintf(t, "%02X ", nes.CPURead(uint16(i))))
	}
}

func (nes *Debugger) DrawStack(t *textutil.Text) {
	t.Color(colornames.White)
	plz.Just(fmt.Fprintf(t, "Stack: 0x%02X\n     ", nes.CPU.S))
	t.Color(colornames.Yellow)
	for i := 0; i <= 0xF; i++ {
		plz.Just(fmt.Fprintf(t, "%02X ", uint16(i)))
	}

	for i := 0x0100; i <= 0x01FF; i++ {
		if i%16 == 0 {
			t.Color(colornames.Yellow)
			plz.Just(fmt.Fprintf(t, "\n%04X ", uint16(i&0xFFF0)))
		}
		if nes.CPU.S == uint8(i) {
			t.Color(colornames.Green)
		} else {
			t.Color(colornames.White)
		}
		plz.Just(fmt.Fprintf(t, "%02X ", nes.CPURead(uint16(i))))
	}
}

func (nes *Debugger) DrawRAM(t *textutil.Text) {
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
			if nes.CPURead(uint16(x+y)) != 0x00 {
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
				plz.Just(fmt.Fprintf(t, "%02X ", nes.CPURead(uint16(x+y))))
			}
		}
	}
}
