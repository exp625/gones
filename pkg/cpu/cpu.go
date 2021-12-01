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
	RequestNMI bool
	RequestIRQ bool
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
		if cpu.RequestNMI {
			cpu.NMI()
			cpu.RequestNMI = false
			cpu.CycleCount += 7
		} else if cpu.RequestIRQ && !cpu.Get(FlagInterruptDisable) {
			cpu.IRQ()
			cpu.RequestIRQ = false
			cpu.CycleCount += 8
		} else {
			if inst.Length != 0 {
				cpu.Bus.Log()
				loc, data, addCycle := inst.AddressMode()
				inst.Execute(loc, data, inst.Length)
				cpu.CycleCount += inst.ClockCycles + int(addCycle)
			}
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
	// Get current pc
	pc := cpu.PC
	// Store high bytes of pc to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8((pc>>8)&0x00FF))
	cpu.S--
	// Store low bytes of pc to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(pc&0x00FF))
	cpu.S--
	// Set flags and store current pc onto stack
	// From https://wiki.nesdev.org/w/index.php?title=Status_flags
	// In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), cpu.P|FlagInterruptDisable)
	cpu.S--
	// Set Interrupt disable flag
	// We don't want another interrupt inside the interrupt handler
	cpu.Set(FlagInterruptDisable, true)
	// Get pc from IRQ/BRK vector and jump to location
	low := uint16(cpu.Bus.CPURead(IRQVector))
	high := uint16(cpu.Bus.CPURead(IRQVector + 1))
	cpu.PC = (high << 8) | low
}

func (cpu *CPU) NMI() {
	// Get current pc
	pc := cpu.PC
	// Store high bytes of pc to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8((pc>>8)&0x00FF))
	cpu.S--
	// Store low bytes of pc to stack
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(pc&0x00FF))
	cpu.S--
	// Set flags and store current pc onto stack
	// From https://wiki.nesdev.org/w/index.php?title=Status_flags
	// In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), cpu.P|FlagInterruptDisable)
	cpu.S--
	// Set Interrupt disable flag
	// We don't want another interrupt inside the interrupt handler
	cpu.Set(FlagInterruptDisable, true)
	// Get pc from NMI vector and jump to location
	low := uint16(cpu.Bus.CPURead(NMIVector))
	high := uint16(cpu.Bus.CPURead(NMIVector + 1))
	cpu.PC = (high << 8) | low
}
