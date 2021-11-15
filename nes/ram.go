package nes

type RAM struct {
	Data [0x0800]uint8
}

func (ram *RAM) Reset() {
	ram.Data = [0x0800]uint8{}
}

func (ram *RAM) Read(location uint16) (bool, uint8) {
	return true, ram.Data[location]
}

func (ram *RAM) Write(location uint16, data uint8) bool {
	ram.Data[location] = data
	return true
}
