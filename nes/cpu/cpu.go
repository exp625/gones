package cpu

import (
	"github.com/exp625/gones/nes/bus"
)

var CPU *CPU6502

const (
	ZeroPage    uint16 = 0x0000
	StackPage   uint16 = 0x0100
	NMIVector   uint16 = 0xFFFA
	ResetVector uint16 = 0xFFFC
	IRQVector   uint16 = 0xFFFE
)

func init() {
	CPU = &CPU6502{}
}

type CPU6502 struct {
	// Accumulator
	A uint8
	// Index X
	X uint8
	// Index Y
	Y uint8
	// Program Counter
	PC uint16
	// Stack Pointer
	S uint8
	// Status Register
	P uint8

	Bus        bus.Bus
	ClockCount int64
	CycleCount int
}

func (cpu *CPU6502) Clock() {
	cpu.ClockCount++
	if cpu.CycleCount == 0 {
		opcode := cpu.Bus.CPURead(cpu.PC)
		inst := Instructions[opcode]
		if inst.Length != 0 {
			cpu.Bus.Log()
			loc, data, addCycle := inst.AddressMode()
			inst.Execute(loc, data, inst.Length)
			cpu.CycleCount += inst.ClockCycles + int(addCycle)
		}
	}
	cpu.CycleCount--
}

func (cpu *CPU6502) Reset() {
	// Reset takes 6 clock cycles
	cpu.ClockCount = 0
	cpu.CycleCount = 6
	// Set Registers to zero
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	// Set status flags
	cpu.P = FlagUnused | FlagBreak | FlagInterruptDisable
	// Set stack pointer
	cpu.S = 0xFD
	// Load the program counter from the reset vector
	low := uint16(cpu.Bus.CPURead(ResetVector))
	high := uint16(cpu.Bus.CPURead(ResetVector + 1))
	cpu.PC = (high << 8) | low
}

func (cpu *CPU6502) IRQ() {

}

func (cpu *CPU6502) NMI() {

}
