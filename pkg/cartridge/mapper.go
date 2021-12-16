package cartridge

type Mapper interface {
	Debugger
	CPUMapRead(location uint16) uint16
	CPURead(location uint16) uint8
	CPUMapWrite(location uint16) uint16
	CPUWrite(location uint16, data uint8) bool
	PPUMapRead(location uint16) uint16
	PPURead(location uint16) uint8
	PPUMapWrite(location uint16) uint16
	PPUWrite(location uint16, data uint8) bool
	Mirroring() bool
	Reset()
}
