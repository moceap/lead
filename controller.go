package lead

import (
	"errors"
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
	magic1: 0x5599f13d, // No idea what this is
	magic2: 0x0200,     // nor this
	magic3: 0xaaaa,     // seriously
}

// A Controller is a Lead Energy WiFi LED controller.
type Controller struct {
	address string
	serial  string
	model   string
	conn    net.Conn
}

// NewController returns a new controller object for the given address. See
// also Discover() to return a list of Controllers on a given LAN segment.
func NewController(address string) *Controller {
	return &Controller{
		address: address,
	}
}

// Address returns the address (ip:port) of the LED controller.
func (c *Controller) Address() string { return c.address }

// Model returns the model number for the WiFi controller. This is set only
// if the Controller is created via Discover().
func (c *Controller) Model() string { return c.model }

// Serial returns the serial number for the WiFi controller. This is set
// only if the Controller is created via Discover().
func (c *Controller) Serial() string { return c.serial }

func (c *Controller) lazyConnect() error {
	if c.conn == nil {
		conn, err := net.DialTimeout("tcp", c.address, connectTimeout)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}

// SetBrightness sets the brightness factor. The range of
// valid values is 0 through 63 inclusive.
func (c *Controller) SetBrightness(b int) error {
	if err := c.lazyConnect(); err != nil {
		return err
	}

	if b < 0 || b > 0x3F {
		return errors.New("value out of range")
	}

	msg := defaultMessage
	msg.command = setBrightness
	msg.value = uint8(b)

	return msg.writeTo(c.conn)
}

// SetRGB sets the color. The range of valid values for r, g and b is 0
// through 255, inclusive.
func (c *Controller) SetRGB(r, g, b int) error {
	if err := c.lazyConnect(); err != nil {
		return err
	}

	if r < 0 || r > 0xFF || g < 0 || g > 0xFF || b < 0 || b > 0xFF {
		return errors.New("value out of range")
	}

	msg := defaultMessage
	msg.command = setRed
	msg.value = uint8(r)

	if err := msg.writeTo(c.conn); err != nil {
		return err
	}

	msg.command = setGreen
	msg.value = uint8(g)
	msg.check()

	if err := msg.writeTo(c.conn); err != nil {
		return err
	}

	msg.command = setBlue
	msg.value = uint8(b)
	msg.check()

	return msg.writeTo(c.conn)
}

// Close closes the connection to the LED controller.
func (c *Controller) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}
