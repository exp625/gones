// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package ppu

type StatusRegister uint8

func (S *StatusRegister) VerticalBlank() bool {
	const bit = 1 << 7
	return *S&bit == bit
}

func (S *StatusRegister) SetVerticalBlank(value bool) {
	const bit = uint8(1) << 7
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<7)
}

func (S *StatusRegister) SpriteZeroHit() bool {
	const bit = 1 << 6
	return *S&bit == bit
}

func (S *StatusRegister) SetSpriteZeroHit(value bool) {
	const bit = uint8(1) << 6
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<6)
}

func (S *StatusRegister) SpriteOverflow() bool {
	const bit = 1 << 5
	return *S&bit == bit
}

func (S *StatusRegister) SetSpriteOverflow(value bool) {
	const bit = uint8(1) << 5
	valueInt := uint8(0)
	if value {
		valueInt = 1
	}

	*S = StatusRegister((uint8(*S) & ^bit) | valueInt<<5)
}
