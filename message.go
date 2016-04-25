package lead

import (
	"encoding/binary"
	"io"
)

type message struct {
	magic1   uint32
	magic2   uint16
	command  uint16
	value    uint8
	checksum uint8
	magic3   uint16
}

func (m *message) check() {
	var c = uint8(0xe4)
	c += uint8(m.magic1 >> 24)
	c += uint8(m.magic1 >> 16)
	c += uint8(m.magic1 >> 8)
	c += uint8(m.magic1)
	c += uint8(m.magic2 >> 8)
	c += uint8(m.magic2)
	c += uint8(m.command >> 8)
	c += uint8(m.command)
	c += uint8(m.value)
	m.checksum = c
}

func (m *message) writeTo(w io.Writer) error {
	m.check()
	return binary.Write(w, binary.BigEndian, m)
}
