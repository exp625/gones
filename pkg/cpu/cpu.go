package cpu

import (
	"github.com/exp625/gones/pkg/bus"
)

type CPU struct {
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

	Bus bus.Bus

	Instructions [256]Instruction
	Mnemonics    map[uint8][2]string

	ClockCount int64
	CycleCount int
}

func New() *CPU {
	c := &CPU{}
	c.generateInstructions()
	return c
}

const (
	ZeroPage    uint16 = 0x0000
	StackPage          = 0x0100
	NMIVector          = 0xFFFA
	ResetVector        = 0xFFFC
	IRQVector          = 0xFFFE
)

func (cpu *CPU) Clock() {
	cpu.ClockCount++
	if cpu.CycleCount == 0 {
		opcode := cpu.Bus.CPURead(cpu.PC)
		inst := cpu.Instructions[opcode]
		if inst.Length != 0 {
			cpu.Bus.Log()
			loc, data, addCycle := inst.AddressMode()
			inst.Execute(loc, data, inst.Length)
			cpu.CycleCount += inst.ClockCycles + int(addCycle)
		}
	}
	cpu.CycleCount--
}

func (cpu *CPU) Reset() {
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

func (cpu *CPU) IRQ() {
}

func (cpu *CPU) NMI() {
}
