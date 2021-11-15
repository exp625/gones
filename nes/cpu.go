package nes

import (
	"fmt"
)

var CPU *C

func init() {
	CPU = &C{}
}

type C struct {
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

	Bus *Bus

	ClockCount         int64
	CycleCount         int
	CurrentInstruction Instruction
	CurrentPC          uint16
}

func (cpu *C) Clock() {
	if cpu.CycleCount == 0 && cpu.CurrentInstruction.Length != 0 {

		fmt.Printf("0x%04X ", cpu.CurrentPC)


		fmt.Printf( "0x%02X ", cpu.Bus.CPURead(cpu.CurrentPC))
		fmt.Print( "[",OpCodeMap[cpu.Bus.CPURead(cpu.CurrentPC)], "] ")
		inst := cpu.CurrentInstruction
		addr, _ := inst.AddressMode()
		fmt.Printf("(0x%04X) ", addr)

		for i := 1; i < int(inst.Length); i++ {
			fmt.Printf("%02X ", cpu.Bus.CPURead(cpu.Bus.CPU.CurrentPC+uint16(i)))
		}
		fmt.Print("\n")
		// Execute Instruction
		loc, data := cpu.CurrentInstruction.AddressMode()
		cpu.CurrentInstruction.Execute(loc, data, cpu.CurrentInstruction.Length)
		cpu.CycleCount = cpu.CurrentInstruction.ClockCycles
	}
	cpu.ClockCount++
	cpu.CycleCount--
	if cpu.CycleCount == 0 {
		// Execution Complete. Load next Instruction
		opcode := cpu.Bus.CPURead(cpu.PC)
		i := Instructions[opcode]
		cpu.CurrentInstruction = i
		cpu.CurrentPC = cpu.PC
	}
}

func (cpu *C) Reset() {
	cpu.ClockCount = 0
	cpu.CycleCount = 0
	// Set IRQ Disabled
	cpu.P = FlagIRQDisabled
	cpu.S = 0xFD

	// Load the program counter
	low := uint16(cpu.Bus.CPURead(StartLocation))
	high := uint16(cpu.Bus.CPURead(StartLocation + 1))
	cpu.PC = (high << 8) | low

	// Debug
	opcode := cpu.Bus.CPURead(cpu.PC)
	i := Instructions[opcode]
	cpu.CurrentInstruction = i
	cpu.CurrentPC = cpu.PC
}

func (cpu *C) IRQ() {

}

func (cpu *C) NMI() {

}
