// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package cpu

type StatusRegister uint8

func (S *StatusRegister) Negative() bool {
	const bit = 1 << 7
	return *S&bit == bit
}

func (S *StatusRegister) SetNegative(value bool) {
	const bit = uint8(1) << 7
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<7)
}

func (S *StatusRegister) Overflow() bool {
	const bit = 1 << 6
	return *S&bit == bit
}

func (S *StatusRegister) SetOverflow(value bool) {
	const bit = uint8(1) << 6
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<6)
}

func (S *StatusRegister) Unused() bool {
	const bit = 1 << 5
	return *S&bit == bit
}

func (S *StatusRegister) SetUnused(value bool) {
	const bit = uint8(1) << 5
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<5)
}

func (S *StatusRegister) Break() bool {
	const bit = 1 << 4
	return *S&bit == bit
}

func (S *StatusRegister) SetBreak(value bool) {
	const bit = uint8(1) << 4
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<4)
}

func (S *StatusRegister) Decimal() bool {
	const bit = 1 << 3
	return *S&bit == bit
}

func (S *StatusRegister) SetDecimal(value bool) {
	const bit = uint8(1) << 3
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<3)
}

func (S *StatusRegister) InterruptDisable() bool {
	const bit = 1 << 2
	return *S&bit == bit
}

func (S *StatusRegister) SetInterruptDisable(value bool) {
	const bit = uint8(1) << 2
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<2)
}

func (S *StatusRegister) Zero() bool {
	const bit = 1 << 1
	return *S&bit == bit
}

func (S *StatusRegister) SetZero(value bool) {
	const bit = uint8(1) << 1
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<1)
}

func (S *StatusRegister) Carry() bool {
	const bit = 1 << 0
	return *S&bit == bit
}

func (S *StatusRegister) SetCarry(value bool) {
	const bit = uint8(1) << 0
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<0)
}
