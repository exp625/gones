package nes

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
		// Execute Instruction
		cpu.CurrentInstruction.Execute(cpu.CurrentInstruction.AddressMode())
		cpu.CycleCount = cpu.CurrentInstruction.ClockCycles
	}
	cpu.ClockCount++
	cpu.CycleCount--
	if cpu.CycleCount == 0 {
		// Execution Complete. Load next Instruction
		cpu.PC += uint16(cpu.CurrentInstruction.Length)
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
