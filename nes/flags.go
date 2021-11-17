package nes

type Flag uint8

const (
	FlagIRQDisabled = 0x34
	StartLocation   = 0xFFFC
)

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

func (cpu *C) Set(flag uint8, value bool) {
	if value {
		cpu.P = cpu.P | flag
	} else {
		cpu.P = cpu.P & ^flag
	}
}

func (cpu *C) GetFlag(flag uint8) bool {
	return cpu.P&flag == flag
}
