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
