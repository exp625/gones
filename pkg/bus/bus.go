package bus

type Bus interface {
	CPUMap(location uint16) uint16
	CPURead(location uint16) uint8
	CPUWrite(location uint16, data uint8)
	PPUMap(location uint16) uint16
	PPURead(location uint16) uint8
	PPUReadRam(location uint16) uint8
	PPUReadPalette(location uint16) uint8
	PPUWrite(location uint16, data uint8)
	PPUWriteRam(location uint16, data uint8)
	PPUWritePalette(location uint16, data uint8)
	DMA(page uint8)
	NMI()
	IRQ()
}
