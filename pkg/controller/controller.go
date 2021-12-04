package controller

type Button uint8

const (
	ButtonA Button = 1 << iota
	ButtonB
	ButtonSELECT
	ButtonSTART
	ButtonUP
	ButtonDOWN
	ButtonLEFT
	ButtonRIGHT
)

type Controller struct {
	Buttons         uint8
	register        uint8
	serialMode      bool
	serialReadCount uint8
}

func (c *Controller) SetMode(mode bool) {
	if c.serialMode == true && mode == false {
		c.register = c.Buttons
		c.serialReadCount = 0
	}
	c.serialMode = mode
}

func (c *Controller) SerialRead() uint8 {
	bit := uint8(1)
	if c.serialMode && c.serialReadCount < 8 {
		bit = (c.register >> c.serialReadCount) & 0b1
		c.serialReadCount++
	}
	return bit
}
func (c *Controller) Press(b Button) {
	c.Buttons = c.Buttons | uint8(b)
}

func (c *Controller) Release(b Button) {
	c.Buttons = c.Buttons & ^uint8(b)
}

func (c *Controller) IsPressed(b Button) bool {
	return c.Buttons&uint8(b) == uint8(b)
}
