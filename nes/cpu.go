package nes

var CPU *C

const (
	ZeroPage uint16 = 0x0100
	NMIVector uint16 = 0xFFFA
	ResetVector uint16 = 0xFFFC
	IRQVector uint16 = 0xFFFE
)

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
	NES *NES
	ClockCount         int64
	CycleCount         int
	CurrentInstruction Instruction
	CurrentPC          uint16
}

func (cpu *C) Clock() {
	cpu.ClockCount++
	if cpu.CycleCount == 0  {
		opcode := cpu.Bus.CPURead(cpu.PC)
		i := Instructions[opcode]
		cpu.CurrentInstruction = i
		cpu.CurrentPC = cpu.PC
		if cpu.CurrentInstruction.Length != 0 {
			// Execute Instruction
			cpu.NES.Log()
			loc, data, addCycle := cpu.CurrentInstruction.AddressMode()
			cpu.CurrentInstruction.Execute(loc, data, cpu.CurrentInstruction.Length)
			cpu.CycleCount += cpu.CurrentInstruction.ClockCycles + int(addCycle)
		}
	}
	cpu.CycleCount--
}

func (cpu *C) Reset() {
	// Reset takes 6 clock cycles
	cpu.ClockCount = 0
	cpu.CycleCount = 6
	// Set IRQ Disabled
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
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
