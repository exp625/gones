package shift_register

type ShiftRegister8 struct {
	register uint8
}

func (r *ShiftRegister8) Set(value uint8) {
	r.register = value
}

func (r *ShiftRegister8) Get() uint8 {
	return r.register
}

func (r *ShiftRegister8) ShiftLeft(bit uint8) uint8 {
	ret := (r.register & 0b10000000) >> 7
	r.register = r.register<<1 | (bit & 0b1)
	return ret
}

func (r *ShiftRegister8) ShiftRight(bit uint8) uint8 {
	ret := r.register & 0b1
	r.register = r.register>>1 | (bit&0b1)<<7
	return ret
}

func (r *ShiftRegister8) SetBit(index uint8, value uint8) {
	bit := uint8(1) << index
	value = value & 0b1
	r.register = (r.register & ^bit) | value<<7
}

func (r *ShiftRegister8) GetBit(index uint8) uint8 {
	bit := uint8(1) << index
	if r.register&bit == bit {
		return 1
	}
	return 0
}
