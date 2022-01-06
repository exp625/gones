package cartridge

type Mapper interface {
	Debugger
	CPUMap(location uint16) uint16
	CPURead(location uint16) uint8
	CPUWrite(location uint16, data uint8) bool
	PPUMap(location uint16) uint16
	PPURead(location uint16) uint8
	PPUWrite(location uint16, data uint8) bool
	Reset()
	Scanline()
}
