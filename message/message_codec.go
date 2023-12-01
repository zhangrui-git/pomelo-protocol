package message

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
)

const (
	Request  byte = 0x00
	Notify   byte = 0x01
	Response byte = 0x02
	Push     byte = 0x03
)

const (
	CompressRouteMask uint8 = 0x01
	CompressGzipMask  uint8 = 0x10
	TypeMask          uint8 = 0xe
)

var (
	ErrWrongMessageType = errors.New("wrong message type")
)

var (
	routes = make(map[string]uint16)
	codes  = make(map[uint16]string)
)

func SetRouteDict(dict map[string]uint16) {
	for r, c := range dict {
		if _, ok := routes[r]; ok {
			continue
		}
		if _, ok := codes[c]; ok {
			continue
		}
		routes[r] = c
		codes[c] = r
	}
}

func Encode(m *Message) ([]byte, error) {
	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}

	var buf = make([]byte, 0)

	var flag byte = 0
	flag |= m.Type << 1
	if m.CompressGzip {
		flag |= CompressGzipMask
	}
	code, exist := routes[m.Route]
	if m.CompressRoute && exist {
		flag |= CompressRouteMask
	} else {
		m.CompressRoute = false
	}
	buf = append(buf, flag)

	if hasId(m.Type) {
		id := m.Id
		for {
			t := byte(id & 0x7f)
			if id >>= 7; id > 0 {
				t |= 0x80
			}
			buf = append(buf, t)
			if id == 0 {
				break
			}
		}
	}

	if hasRoute(m.Type) {
		if m.CompressRoute {
			buf = append(buf, byte(code>>8&0xff))
			buf = append(buf, byte(code&0xff))
		} else {
			buf = append(buf, byte(len(m.Route)))
			buf = append(buf, []byte(m.Route)...)
		}
	}

	if m.CompressGzip {
		df, err := deflate(m.Body)
		if err != nil {
			return nil, err
		}
		buf = append(buf, df...)
	} else {
		buf = append(buf, m.Body...)
	}

	return buf, nil
}

func Decode(data []byte) (*Message, error) {
	flag := data[0]
	mType := (flag & TypeMask) >> 1
	compressGzip := false
	if (flag&CompressGzipMask)>>4 > 0 {
		compressGzip = true
	}
	compressRoute := false
	if flag&CompressRouteMask > 0 {
		compressRoute = true
	}

	var offset uint = 1

	var id uint
	if hasId(mType) {
		var m uint = 0
		i := 0
		for {
			m = uint(data[offset])
			id += (m & 0x7f) << (7 * i)
			offset += 1
			i++
			if m < 128 {
				break
			}
		}
	}

	var route string
	if hasRoute(mType) {
		if compressRoute {
			code := binary.BigEndian.Uint16(data[offset : offset+2])
			offset += 2
			r, ok := codes[code]
			if ok {
				route = r
			}
		} else {
			rl := uint(data[offset])
			offset += 1
			route = string(data[offset : offset+rl])
			offset += rl
		}
	}

	body := data[offset:]
	if compressGzip {
		var err error
		body, err = inflate(body)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:          mType,
		CompressRoute: compressRoute,
		CompressGzip:  compressGzip,
		Id:            id,
		Route:         route,
		Body:          body,
	}, nil
}

func invalidType(mType uint8) bool {
	return mType != Request && mType != Notify && mType != Response && mType != Push
}

func hasId(mType uint8) bool {
	return mType == Request || mType == Response
}

func hasRoute(mType uint8) bool {
	return mType == Request || mType == Notify || mType == Push
}

func deflate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	if _, err := gzw.Write(data); err != nil {
		return nil, err
	}
	if err := gzw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func inflate(data []byte) ([]byte, error) {
	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if err = gzr.Close(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gzr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
