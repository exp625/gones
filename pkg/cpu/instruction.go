package cpu

type Instruction struct {
	AddressModeMnemonic string
	ExecuteMnemonic     string
	AddressMode         func() (uint16, uint8, uint8)
	Execute             func(uint16, uint8, uint16)
	Length              uint16
	ClockCycles         int
	Legal               bool
}
