package nes

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

	bus Bus
}

func (cpu *CPU) Clock() {

}

func (cpu *CPU) Reset() {
	// Set IRQ Disabled
	cpu.P = FlagIRQDisabled
	cpu.S = 0xFD

	// Load the program counter
	low := uint16(cpu.bus.CPURead(StartLocation))
	high := uint16(cpu.bus.CPURead(StartLocation + 1))
	cpu.PC = (high << 8) | low
}
