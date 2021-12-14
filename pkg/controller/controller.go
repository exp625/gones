package controller

import "github.com/exp625/gones/internal/shift_register"

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
	Buttons    uint8
	register   shift_register.ShiftRegister8
	serialMode bool
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) SetMode(mode bool) {
	if c.serialMode == true && mode == false {
		c.register.Set(c.Buttons)
	}
	c.serialMode = mode
}

func (c *Controller) SerialRead() uint8 {
	if c.serialMode {
		return c.register.ShiftRight(1)
	}
	return 1
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
