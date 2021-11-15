package nes

type PPU struct {
}

func (ppu *PPU) Clock() {

}

func (ppu *PPU) Reset() {

}

func (ppu *PPU) Read(location uint16) (bool, uint8) {
	return true, 0
}

func (ppu *PPU) Write(location uint16, data uint8) {

}
