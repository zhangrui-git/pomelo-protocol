package pomelo_protocol

const (
	PkgHeadBytes = 4
)

const (
	PkgTypeHandshake    = 1
	PkgTypeHandshakeAck = 2
	PkgTypeHeartbeat    = 3
	PkgTypeData         = 4
	PkgTypeKick         = 5
)

const (
	MsgTypeRequest  uint8 = 0
	MsgTypeNotify   uint8 = 1
	MsgTypeResponse uint8 = 2
	MsgTypePush     uint8 = 3
)

const (
	MsgFlagBytes      = 1
	MsgRouteCodeBytes = 2
	MsgRouteLenBytes  = 1
)

const (
	MsgCompressGzipEncodeMask uint8 = 0x10
	MsgCompressRouteMask      uint8 = 0x01
	MsgTypeMask               uint8 = 0xe
)

type Package struct {
	PkgType uint8
	Msg     *Message
}

type Message struct {
	MsgId         uint
	MsgType       uint8
	CompressRoute bool
	CompressGzip  bool
	Route         []byte
	MsgBody       []byte
}

func msgEncode(m *Message) []byte {
	idBytes := 0
	if msgHasId(m.MsgType) {
		idBytes = msgIdBytes(m.MsgId)
	}
	msgLen := MsgFlagBytes + idBytes
	if msgHasRoute(m.MsgType) {
		if m.CompressRoute {
			msgLen += MsgRouteCodeBytes
		} else {
			msgLen += MsgRouteLenBytes + len(m.Route)
		}
	}
	msgLen += len(m.MsgBody)
	buffer := make([]byte, msgLen)

	offset := 0

	offset = encodeMsgFlag(m.MsgType, m.CompressRoute, m.CompressGzip, buffer, offset)
	if msgHasId(m.MsgType) {
		offset = encodeMsgId(m.MsgId, buffer, offset)
	}
	if msgHasRoute(m.MsgType) {
		offset = encodeMsgRoute(m.Route, m.CompressRoute, buffer, offset)
	}
	if len(m.MsgBody) > 0 {
		offset = encodeMsgBody(m.MsgBody, buffer, offset)
	}
	return buffer
}

func msgDecode(data []byte) *Message {
	offset := 0

	flag := data[offset]
	offset += 1

	compressGzip := false
	if (flag&MsgCompressGzipEncodeMask)>>4 > 0 {
		compressGzip = true
	}
	msgType := (flag & MsgTypeMask) >> 1
	compressRoute := false
	if flag&MsgCompressRouteMask > 0 {
		compressRoute = true
	}

	var msgId uint = 0

	if msgHasId(msgType) {
		var m uint = 0
		i := 0
		for {
			m = uint(data[offset])
			msgId += (m & 0x7f) << (7 * i)
			offset += 1
			i++
			if m < 128 {
				break
			}
		}
	}

	var route []byte
	if msgHasRoute(msgType) {
		if compressRoute {
			route = make([]byte, 2)
			route[0] = data[offset]
			route[1] = data[offset+1]
			offset += MsgRouteCodeBytes
		} else {
			routeLen := int(data[offset])
			offset += MsgRouteLenBytes
			route = make([]byte, routeLen)
			if routeLen > 0 {
				copyByteSlice(route, 0, data, offset, routeLen)
			}
			offset += routeLen
		}
	}

	msgBodyLen := len(data) - offset
	msgBody := make([]byte, msgBodyLen)
	copyByteSlice(msgBody, 0, data, offset, msgBodyLen)

	return &Message{
		MsgId:         msgId,
		MsgType:       msgType,
		CompressRoute: compressRoute,
		CompressGzip:  compressGzip,
		Route:         route,
		MsgBody:       msgBody,
	}
}

func Encode(p *Package) []byte {
	var body []byte
	bodyLen := 0
	if p.Msg != nil {
		body = msgEncode(p.Msg)
		bodyLen = len(body)
	}
	buffer := make([]byte, PkgHeadBytes+bodyLen)

	buffer[0] = p.PkgType & 0xff
	buffer[1] = byte(bodyLen >> 16 & 0xff)
	buffer[2] = byte(bodyLen >> 8 & 0xff)
	buffer[3] = byte(bodyLen & 0xff)

	if bodyLen > 0 {
		copyByteSlice(buffer, 4, body, 0, bodyLen)
	}

	return buffer
}

func Decode(data []byte) *Package {
	pkgType := data[0]

	var msg *Message
	bodyLen := int(data[1])<<16 | int(data[2])<<8 | int(data[3])
	if bodyLen > 0 {
		body := make([]byte, bodyLen)
		copyByteSlice(body, 0, data, 4, bodyLen)
		msg = msgDecode(body)
	}

	return &Package{PkgType: pkgType, Msg: msg}
}

func msgHasId(msgType uint8) bool {
	return msgType == MsgTypeRequest || msgType == MsgTypeResponse
}

func msgHasRoute(msgType uint8) bool {
	return msgType == MsgTypeRequest || msgType == MsgTypeNotify || msgType == MsgTypePush
}

func encodeMsgFlag(msgType uint8, compressRoute bool, compressGzip bool, buffer []byte, offset int) int {
	buffer[offset] = msgType << 1
	if compressRoute {
		buffer[offset] |= MsgCompressRouteMask
	}
	if compressGzip {
		buffer[offset] |= MsgCompressGzipEncodeMask
	}
	return offset + MsgFlagBytes
}

func encodeMsgId(id uint, buffer []byte, offset int) int {
	for {
		t := id & 0x7f
		if id >>= 7; id > 0 {
			t |= 0x80
		}
		buffer[offset] = byte(t)
		offset += 1
		if id == 0 {
			break
		}
	}
	return offset
}

func encodeMsgRoute(route []byte, compressRoute bool, buffer []byte, offset int) int {
	if compressRoute {
		buffer[offset] = route[0]
		offset += 1
		buffer[offset] = route[1]
		offset += 1
	} else {
		if len(route) > 0 {
			buffer[offset] = byte(len(route))
			offset += 1
			copyByteSlice(buffer, offset, []byte(route), 0, len(route))
			offset += len(route)
		} else {
			buffer[offset] = 0
			offset += 1
		}
	}
	return offset
}

func encodeMsgBody(msgBody []byte, buffer []byte, offset int) int {
	copyByteSlice(buffer, offset, []byte(msgBody), 0, len(msgBody))
	return offset + len(msgBody)
}

func msgIdBytes(id uint) int {
	l := 0
	for {
		l++
		id >>= 7
		if id == 0 {
			break
		}
	}
	return l
}

func copyByteSlice(dest []byte, dOffset int, src []byte, sOffset, length int) {
	for i := 0; i < length; i++ {
		dest[dOffset+i] = src[sOffset+i]
	}
}
