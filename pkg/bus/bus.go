package bus

type Bus interface {
	CPURead(location uint16) uint8
	CPUWrite(location uint16, data uint8)
	PPURead(location uint16) uint8
	PPUWrite(location uint16, data uint8)
	Log()
	NMI()
	IRQ()
}