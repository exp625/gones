package cpu

type Flag uint8

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

func (cpu *CPU6502) Set(flag uint8, value bool) {
	if value {
		cpu.P = cpu.P | flag
	} else {
		cpu.P = cpu.P & ^flag
	}
}

func (cpu *CPU6502) GetFlag(flag uint8) bool {
	return cpu.P&flag == flag
}
