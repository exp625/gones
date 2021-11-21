package cpu

const (
	FlagCarry uint8 = 1 << iota
	FlagZero
	FlagInterruptDisable
	FlagDecimal
	FlagBreak
	FlagUnused // Unused
	FlagOverflow
	FlagNegative
)

func (cpu *CPU) Set(flag uint8, value bool) {
	if value {
		cpu.P = cpu.P | flag
	} else {
		cpu.P = cpu.P & ^flag
	}
}

func (cpu *CPU) Get(flag uint8) bool {
	return cpu.P&flag == flag
}
