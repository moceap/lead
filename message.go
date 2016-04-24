package lead

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

type message struct {
	Magic1  uint32
	Magic2  uint16
	Command uint16
	Value   uint8
	Check   uint8
	End     uint16
}

func (m *message) check() {
	var c = uint8(0xe4)
	c += uint8(m.Magic1 >> 24)
	c += uint8(m.Magic1 >> 16)
	c += uint8(m.Magic1 >> 8)
	c += uint8(m.Magic1)
	c += uint8(m.Magic2 >> 8)
	c += uint8(m.Magic2)
	c += uint8(m.Command >> 8)
	c += uint8(m.Command)
	c += uint8(m.Value)
	m.Check = c
}

func (m *message) print() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m)
	fmt.Println(hex.Dump(buf.Bytes()))
}

func (m *message) writeTo(w io.Writer) error {
	m.check()
	return binary.Write(w, binary.BigEndian, m)
}
