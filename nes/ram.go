package nes

type RAM struct {
	data [0x07FF]uint8
}

func (ram *RAM) Reset() {
	ram.data = [0x7FF]uint8{}
}

func (ram *RAM) Read(location uint16) uint8 {
	return ram.data[location]
}

func (ram *RAM) Write(location uint16, data uint8) {
	ram.data[location] = data
}
