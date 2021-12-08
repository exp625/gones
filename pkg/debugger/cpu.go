package debugger

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"golang.org/x/image/colornames"
)

func (nes *Debugger) DrawCPU(t *textutil.Text) {
	// Print Status
	plz.Just(fmt.Fprint(t, "CPU STATUS \t"))

	arr := "CZIDB-VN"
	for i := 0; i < 8; i++ {
		if nes.CPU.Get(1 << i) {
			t.Color(colornames.Green)
		} else {
			t.Color(colornames.Red)
		}
		plz.Just(fmt.Fprint(t, string(arr[i])))
	}
	t.Color(colornames.White)
	plz.Just(fmt.Fprintf(t, "\t%02X", nes.CPU.P))
	plz.Just(fmt.Fprint(t, "\n"))
	plz.Just(fmt.Fprintf(t, "PC: 0x%02X\t", nes.CPU.PC))
	plz.Just(fmt.Fprintf(t, "A: 0x%02X\t", nes.CPU.A))
	plz.Just(fmt.Fprintf(t, "X: 0x%02X\t", nes.CPU.X))
	plz.Just(fmt.Fprintf(t, "Y: 0x%02X\t", nes.CPU.Y))
	plz.Just(fmt.Fprintf(t, "S: 0x%02X\t\n", nes.CPU.S))
}

func (nes *Debugger) DrawInstructions(t *textutil.Text) {
	offset := uint16(0)
	if nes.CPU.CycleCount < 0 {
		plz.Just(fmt.Fprint(t, "ERR"))
		return
	}
	for j := 0; j <= 5; j++ {
		if j == 0 {
			t.Color(colornames.Yellow)
		}
		plz.Just(fmt.Fprintf(t, "%04X ", nes.CPU.PC+offset))
		inst := nes.CPU.Instructions[nes.CPURead(nes.CPU.PC+offset)]
		i := 0
		for ; i < int(inst.Length); i++ {
			plz.Just(fmt.Fprintf(t, "%02X ", nes.CPURead(nes.CPU.PC+offset+uint16(i))))
		}
		for ; i < 3; i++ {
			plz.Just(fmt.Fprint(t, "   "))
		}
		plz.Just(fmt.Fprint(t, "", nes.CPU.Mnemonics[nes.CPURead(nes.CPU.PC+offset)], " "))

		if inst.Length != 0 {
			addr, data, _ := inst.AddressMode(nes.CPURead)
			if j == 0 {
				// Display Address
				switch nes.CPU.Mnemonics[nes.CPURead(nes.CPU.PC)][1] {
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
					plz.Just(fmt.Fprintf(t, "$%02X,X @ %02X = %02X", nes.CPURead(nes.CPU.PC+1), addr, data))
				case "ZPY":
					plz.Just(fmt.Fprintf(t, "$%02X,Y @ %02X = %02X", nes.CPURead(nes.CPU.PC+1), addr, data))
				case "ZP0":
					plz.Just(fmt.Fprintf(t, "$%02X = %02X", addr&0x00FF, data))
				case "IDX":
					// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
					plz.Just(fmt.Fprintf(t, "($%02X,X) @ %02X = %04X = %02X", nes.CPURead(nes.CPU.PC+1), nes.CPURead(nes.CPU.PC+1)+nes.CPU.X, addr, data))
				case "IZY":
					// Second byte is added to register X -> result is a zero page address where the actual memory location is stored.
					plz.Just(fmt.Fprintf(t, "($%02X),Y = %04X @ %04X = %02X", nes.CPURead(nes.CPU.PC+1), addr-uint16(nes.CPU.Y), addr, data))
				case "IND":
					plz.Just(fmt.Fprintf(t, "($%02X%02X) = %04X", nes.CPURead(nes.CPU.PC+2), nes.CPURead(nes.CPU.PC+1), addr))
				case "ABX":
					plz.Just(fmt.Fprintf(t, "$%02X%02X,X @ %04X = %02X", nes.CPURead(nes.CPU.PC+2), nes.CPURead(nes.CPU.PC+1), addr, data))
				case "ABY":
					plz.Just(fmt.Fprintf(t, "$%02X%02X,Y @ %04X = %02X", nes.CPURead(nes.CPU.PC+2), nes.CPURead(nes.CPU.PC+1), addr, data))
				}
			}
			t.Color(colornames.White)
			offset += inst.Length
			plz.Just(fmt.Fprint(t, "\n"))
		}
	}
}
