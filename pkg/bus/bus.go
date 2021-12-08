package bus

type Bus interface {
	CPURead(location uint16) uint8
	CPUWrite(location uint16, data uint8)
	PPURead(location uint16) uint8
	PPUWrite(location uint16, data uint8)
	DMA(page uint8)
	DMAWrite(data uint8)
	NMI()
	IRQ()
}
