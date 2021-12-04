package controller

const (
	ButtonA uint8 = 1 << iota
	ButtonB
	ButtonSELECT
	ButtonSTART
	ButtonUP
	ButtonDOWN
	ButtonLEFT
	ButtonRIGHT
)

type Controller struct {
	Buttons    uint8
	SerialMode bool
	register   uint8
	serialReadCount uint8
}

func (c *Controller) SetMode(mode bool) {
	if c.SerialMode == true && mode == false {
		c.register = c.Buttons
		c.serialReadCount = 0
	}
	c.SerialMode = mode
}

func (c *Controller) SerialRead() uint8 {
	bit := uint8(1)
	if c.SerialMode && c.serialReadCount < 8 {
		bit = (c.register >> c.serialReadCount) & 0b1
		c.serialReadCount++
	}
	return bit
}
func (c *Controller) Press(b uint8) {
	c.Buttons = c.Buttons | b
}

func (c *Controller) Release(b uint8) {
	c.Buttons = c.Buttons & ^b
}

func (c *Controller) IsPressed(b uint8) bool {
	return c.Buttons&b == b
}
