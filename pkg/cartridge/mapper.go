package cartridge

type Mapper interface {
	Debugger
	CPURead(location uint16) (bool, uint8)
	CPUWrite(location uint16, data uint8) bool
	PPURead(location uint16) (bool, uint8)
	PPUWrite(location uint16, data uint8) bool
	Mirroring() bool
	Reset()
}
