package packet

import "errors"

const (
	Handshake    uint8 = 0x01
	HandshakeAck uint8 = 0x02
	Heartbeat    uint8 = 0x03
	Data         uint8 = 0x04
	Kick         uint8 = 0x05
)

const HeadLength = 4

var (
	ErrWrongPacketType = errors.New("wrong packet type")
)

func Encode(p *Packet) ([]byte, error) {
	if invalidType(p.Type) {
		return nil, ErrWrongPacketType
	}
	var length = 0
	if p.Data != nil {
		length = len(p.Data)
	}
	buf := make([]byte, HeadLength+length)
	buf[0] = p.Type
	buf[1] = byte(length >> 16 & 0xff)
	buf[2] = byte(length >> 8 & 0xff)
	buf[3] = byte(length & 0xff)
	copy(buf[HeadLength:], p.Data)
	return buf, nil
}

func Decode(d []byte) (*Packet, error) {
	pType := d[0]
	length := int(d[1])<<16 + int(d[2])<<8 + int(d[3])
	var data []byte
	if length > 0 {
		data = make([]byte, length)
		copy(data, d[HeadLength:HeadLength+length])
	}
	return &Packet{Type: pType, Data: data}, nil
}

func invalidType(pType uint8) bool {
	return pType != Handshake && pType != HandshakeAck && pType != Heartbeat && pType != Data && pType != Kick
}
