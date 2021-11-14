package nes

type RAM struct {
	Data [0x07FF]uint8
}

func (ram *RAM) Reset() {
	ram.Data = [0x07FF]uint8{}
	ram.Data[0x00FF] = 0xDD
}

func (ram *RAM) Read(location uint16) uint8 {
	return ram.Data[location]
}

func (ram *RAM) Write(location uint16, data uint8) {
	ram.Data[location] = data
}
