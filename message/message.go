package message

type Message struct {
	Type          uint8
	CompressRoute bool
	CompressGzip  bool
	Id            uint
	Route         string
	Body          []byte
}

func (m *Message) Encode() ([]byte, error) {
	return Encode(m)
}
