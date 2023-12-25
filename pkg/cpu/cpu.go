package cpu

import (
	"github.com/exp625/gones/pkg/apu"
	"github.com/exp625/gones/pkg/bus"
	"github.com/exp625/gones/pkg/logger"
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
	// Status register
	P StatusRegister

	Bus bus.Bus

	APU *apu.APU

	Instructions [256]Instruction
	Mnemonics    map[uint8][2]string

	ClockCount           int64
	CycleCount           int
	IRQLine              bool
	IRQRequested         bool
	IRQLinePreviousCycle bool
	NMILine              bool
	NMIRequested         bool
	APUDMA               bool
	PPUDMA               bool
	PPUDMAPrepared       bool
	PPUDMAAddress        uint16

	CurrentInstruction Instruction
	CurrentLocation    uint16

	Logger logger.Loggable
}

func New() *CPU {
	c := &CPU{}
	c.generateInstructions()
	return c
}

// AddBus connects the CPU to the Bus
func (cpu *CPU) AddBus(bus bus.Bus) {
	cpu.Bus = bus
}

// AddAPU connects the CPU to the APU
func (cpu *CPU) AddAPU(apu *apu.APU) {
	cpu.APU = apu
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

	if cpu.CycleCount == 0 && cpu.CurrentInstruction.Length != 0 {
		cpu.CurrentInstruction.Execute(cpu.CurrentLocation, cpu.CurrentInstruction.Length)
		cpu.CurrentInstruction = Instruction{}

		if cpu.NMILine {
			cpu.NMIRequested = true
		}
		cpu.NMILine = false

		if cpu.IRQLinePreviousCycle {
			cpu.IRQRequested = true
		}
		cpu.IRQLinePreviousCycle = cpu.IRQLine

		if cpu.P.InterruptDisable() && cpu.IRQLine {
			cpu.IRQLine = false
		}
	}

	if cpu.CycleCount == 0 {
		if cpu.PPUDMA {
			if cpu.PPUDMAAddress&0xFF == 0xFF {
				cpu.PPUDMA = false
			}
			if !cpu.PPUDMAPrepared {
				cpu.CycleCount++
				if cpu.CycleCount%2 != 0 {
					cpu.CycleCount++
				}
				cpu.PPUDMAPrepared = true
			} else {
				cpu.Bus.CPUWrite(0x2004, cpu.Bus.CPURead(cpu.PPUDMAAddress))
				cpu.PPUDMAAddress++
				// Transfer takes one clock cycle
				cpu.CycleCount++
			}
		} else if cpu.APUDMA {
			cpu.CycleCount += 4
			cpu.APUDMA = false
			cpu.APU.DMA()
		} else {
			if cpu.NMIRequested {
				cpu.NMI()
				cpu.NMIRequested = false
				cpu.CycleCount += 7
			} else if cpu.IRQRequested && !cpu.P.InterruptDisable() {
				cpu.IRQ()
				cpu.IRQRequested = false
				cpu.CycleCount += 8
			} else {
				opcode := cpu.Bus.CPURead(cpu.PC)
				inst := cpu.Instructions[opcode]
				if inst.Length != 0 {
					cpu.log()
					loc, addCycle := inst.AddressMode(cpu.Bus.CPURead)
					cpu.CurrentInstruction = inst
					cpu.CurrentLocation = loc
					cpu.CycleCount += inst.ClockCycles + int(addCycle)
				}
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
	cpu.P = StatusRegister(FlagUnused | FlagBreak | FlagInterruptDisable)
	// Set stack pointer
	cpu.S = 0xFD
	// Load the program counter from the reset vector
	low := uint16(cpu.Bus.CPURead(ResetVector))
	high := uint16(cpu.Bus.CPURead(ResetVector + 1))
	cpu.PC = (high << 8) | low
	cpu.CurrentInstruction = Instruction{}
	cpu.CurrentLocation = 0
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
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(cpu.P)|FlagUnused)
	cpu.S--
	// Set Interrupt disable flag
	// We don't want another interrupt inside the interrupt handler
	cpu.P.SetInterruptDisable(true)
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
	cpu.Bus.CPUWrite(StackPage|uint16(cpu.S), uint8(cpu.P)|FlagUnused)
	cpu.S--
	// Set Interrupt disable flag
	// We don't want another interrupt inside the interrupt handler
	cpu.P.SetInterruptDisable(true)
	// Get pc from NMI vector and jump to location
	low := uint16(cpu.Bus.CPURead(NMIVector))
	high := uint16(cpu.Bus.CPURead(NMIVector + 1))
	cpu.PC = (high << 8) | low
}

func (cpu *CPU) log() {
	cpu.Logger.Log()
}
