package packet

type Packet struct {
	Type uint8
	Data []byte
}

func (p *Packet) Encode() ([]byte, error) {
	return Encode(p)
}
