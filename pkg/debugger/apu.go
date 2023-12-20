package debugger

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"golang.org/x/image/colornames"
)

func (nes *Debugger) DrawAPU(t *textutil.Text) {
	// Print Status
	if !nes.APU.FrameCounterRegister.DisableFrameIRQ() {
		t.Color(colornames.Green)
	} else {
		t.Color(colornames.Red)
	}
	plz.Just(fmt.Fprint(t, "APU Frame Interrupt Enable \n"))
	if nes.APU.FrameInterrupt {
		t.Color(colornames.Green)
	} else {
		t.Color(colornames.Red)
	}
	plz.Just(fmt.Fprint(t, "APU Frame Interrupt \n"))
	t.Color(colornames.White)
	if nes.APU.FrameCounterRegister.FiveFrameSequence() {
		plz.Just(fmt.Fprintf(t, "Sequence: 5\n"))
	} else {
		plz.Just(fmt.Fprintf(t, "Sequence: 4\n"))
	}

	plz.Just(fmt.Fprintf(t, "Counter: %d\n", nes.APU.FrameCounterHalfs))

}

func (nes *Debugger) DrawAPUDMC(t *textutil.Text) {
	// Print Status
	if nes.APU.ControlRegister.DMCEnable() {
		t.Color(colornames.Green)
	} else {
		t.Color(colornames.Red)
	}
	plz.Just(fmt.Fprint(t, "DMC Channel \n"))
	t.Color(colornames.White)
	plz.Just(fmt.Fprintf(t, "Loop: %t\n", nes.APU.DMC.GlobalRegister.Loop()))
	plz.Just(fmt.Fprintf(t, "Interrup: %t\n", nes.APU.DMC.GlobalRegister.IRQEnable()))
	plz.Just(fmt.Fprintf(t, "Address: 0x%02X\n", nes.APU.DMC.SampleAddressRegister))
	plz.Just(fmt.Fprintf(t, "Length: 0x%02X\n", nes.APU.DMC.SampleLengthRegister))
	plz.Just(fmt.Fprintf(t, "BytesRemaing: 0x%02X\n", nes.APU.DMC.BytesRemainingCounter))
	plz.Just(fmt.Fprintf(t, "SampleBuffer: 0x%02X\n", nes.APU.DMC.SampleBuffer))
	plz.Just(fmt.Fprintf(t, "Output: 0x%02X\n", nes.APU.DMC.OutputLevelCounter))
}
