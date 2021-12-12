package templates

type CPUStatusRegister struct {
	Negative         bool `bitfield:"1"`
	Overflow         bool `bitfield:"1"`
	Unused           bool `bitfield:"1"`
	Break            bool `bitfield:"1"`
	Decimal          bool `bitfield:"1"`
	InterruptDisable bool `bitfield:"1"`
	Zero             bool `bitfield:"1"`
	Carry            bool `bitfield:"1"`
}
