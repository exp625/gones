package ppu

type PPU struct {
	ScanLine   uint16
	Position   uint16
	FrameCount uint64
}

func (ppu *PPU) Clock() {

	if ppu.ScanLine < 261 {
		if ppu.Position < 340 {
			ppu.Position++
		} else {
			ppu.Position = 0
			ppu.ScanLine++
		}
	} else {
		ppu.Position = 0
		ppu.ScanLine = 0
		ppu.FrameCount++
		if ppu.FrameCount%2 != 0 {
			ppu.Position++
		}
	}

}

func (ppu *PPU) Reset() {
	ppu.ScanLine = 0
	ppu.Position = 0
	ppu.FrameCount = 0
}

func (ppu *PPU) Read(location uint16) (bool, uint8) {
	return true, 0
}

func (ppu *PPU) Write(location uint16, data uint8) {

}
