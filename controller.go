package lead

import (
	"net"
	"time"
)

const (
	connectTimeout = 5 * time.Second
	setRed         = 0x848
	setGreen       = 0x849
	setBlue        = 0x84a
	setBrightness  = 0x84c
)

var defaultMessage = message{
	Magic1: 0x5599f13d,
	Magic2: 0x0200,
	End:    0xaaaa,
}

type Controller struct {
	Address string
	Serial  string
	Model   string

	conn net.Conn
}

func (c *Controller) lazyConnect() error {
	if c.conn == nil {
		conn, err := net.DialTimeout("tcp", c.Address, connectTimeout)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}

func (c *Controller) SetBrightness(b float64) error {
	if err := c.lazyConnect(); err != nil {
		return err
	}

	msg := defaultMessage
	msg.Command = setBrightness
	msg.Value = clamp(b, 0x3f)

	return msg.writeTo(c.conn)
}

func (c *Controller) SetRGB(r, g, b float64) error {
	if err := c.lazyConnect(); err != nil {
		return err
	}

	msg := defaultMessage
	msg.Command = setRed
	msg.Value = clamp(r, 0xff)

	if err := msg.writeTo(c.conn); err != nil {
		return err
	}

	msg.Command = setGreen
	msg.Value = clamp(g, 0xff)
	msg.check()

	if err := msg.writeTo(c.conn); err != nil {
		return err
	}

	msg.Command = setBlue
	msg.Value = clamp(b, 0xff)
	msg.check()

	return msg.writeTo(c.conn)
}

func clamp(v float64, r uint8) uint8 {
	t := int(v * float64(r))
	if t < 0 {
		return 0
	}
	if t > int(r) {
		return r
	}
	return uint8(t)
}
