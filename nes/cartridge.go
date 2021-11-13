package nes

type Catridge struct {
	data [0xFFFF]uint8
}

func (c *Catridge) Reset() {
	c.data = [0xFFFF]uint8{}
	c.data[0xFFFC] = 0x00
	c.data[0xFFFD] = 0x00
}

func (c *Catridge) Read(location uint16) uint8 {
	return c.data[location]
}

func (c *Catridge) Write(location uint16, data uint8) {
	return
}
