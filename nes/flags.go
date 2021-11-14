package nes

type Flag uint8

const (
	FlagIRQDisabled = 0x34
	StartLocation   = 0xFFFC
)

const (
	FlagCarry Flag = 1 << iota
	FlagZero
	FlagInterruptDisable
	FlagDecimal
	FlagB1 // Unused
	FlagB2 // Unused
	FlagOverflow
	FlagNegative
)

func (cpu *C) Set(flag Flag, value bool) {
	if value {
		cpu.P = cpu.P | uint8(flag)
	} else {
		cpu.P = cpu.P & ^uint8(flag)
	}
}

func (cpu *C) GetFlag(flag Flag) bool {
	return cpu.P&uint8(flag) == uint8(flag)
}
