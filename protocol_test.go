package protocol

import (
	"encoding/binary"
	"testing"
)

var request = &Message{
	MsgId:         12345678,
	MsgType:       MsgTypeRequest,
	CompressRoute: false,
	CompressGzip:  false,
	Route:         []byte("user.login"),
	MsgBody:       []byte("{\"id\":8888,\"password\":\"balabala\"}"),
}

var notify = &Message{
	MsgType:       MsgTypeNotify,
	CompressRoute: false,
	CompressGzip:  false,
	Route:         []byte("mailbox.new"),
	MsgBody:       []byte("{\"from\": 9999, \"content\":\"hi tom\"}"),
}
var response = &Message{
	MsgId:         12345678,
	MsgType:       MsgTypeResponse,
	CompressRoute: false,
	CompressGzip:  false,
	MsgBody:       []byte("{\"code\":0,\"msg\":\"success\"}"),
}

var push = &Message{
	MsgType:       MsgTypePush,
	CompressRoute: true,
	CompressGzip:  false,
	Route:         make([]byte, 2), // two byte spaces
	MsgBody:       []byte("{}"),
}

func TestRequestData(t *testing.T) {
	dataPkg := &Package{
		PkgType: PkgTypeData,
		Msg:     request,
	}
	buffer := Encode(dataPkg)
	t.Logf("%b\n", buffer)
	decode := Decode(buffer)
	t.Logf("%+v\n", decode)
	t.Logf("%+v\n", decode.Msg)
	t.Logf("%s %s\n", decode.Msg.Route, decode.Msg.MsgBody)
}

func TestNotifyData(t *testing.T) {
	dataPkg := &Package{
		PkgType: PkgTypeData,
		Msg:     notify,
	}
	buffer := Encode(dataPkg)
	t.Logf("%b\n", buffer)
	decode := Decode(buffer)
	t.Logf("%+v\n", decode)
	t.Logf("%+v\n", decode.Msg)
	t.Logf("%s %s\n", decode.Msg.Route, decode.Msg.MsgBody)
}

func TestResponseData(t *testing.T) {
	dataPkg := &Package{
		PkgType: PkgTypeData,
		Msg:     response,
	}
	buffer := Encode(dataPkg)
	t.Logf("%b\n", buffer)
	decode := Decode(buffer)
	t.Logf("%+v\n", decode)
	t.Logf("%+v\n", decode.Msg)
	t.Logf("%s %s\n", decode.Msg.Route, decode.Msg.MsgBody)
}

func TestPushData(t *testing.T) {
	dataPkg := &Package{
		PkgType: PkgTypeData,
		Msg:     push,
	}
	binary.BigEndian.PutUint16(dataPkg.Msg.Route, 65535)

	buffer := Encode(dataPkg)
	t.Logf("%b\n", buffer)
	decode := Decode(buffer)
	t.Logf("%+v\n", decode)
	t.Logf("%+v\n", decode.Msg)
	t.Logf("%d %s\n", binary.BigEndian.Uint16(decode.Msg.Route), decode.Msg.MsgBody)
}

func TestEncode(t *testing.T) {
	kickPkg := &Package{
		PkgType: PkgTypeKick,
		Msg:     nil,
	}
	buffer2 := Encode(kickPkg)
	t.Logf("%b\n", buffer2)
	decode2 := Decode(buffer2)
	t.Logf("%+v\n", decode2)
	t.Logf("%+v\n", decode2.Msg)

	handshakePkg := *kickPkg
	handshakePkg.PkgType = PkgTypeHandshake
	buffer3 := Encode(&handshakePkg)
	t.Logf("%b\n", buffer3)
	decode3 := Decode(buffer3)
	t.Logf("%+v\n", decode3)
	t.Logf("%+v\n", decode3.Msg)

	handshakeAckPkg := *kickPkg
	handshakeAckPkg.PkgType = PkgTypeHandshakeAck
	buffer4 := Encode(&handshakeAckPkg)
	t.Logf("%b\n", buffer4)
	decode4 := Decode(buffer4)
	t.Logf("%+v\n", decode4)
	t.Logf("%+v\n", decode4.Msg)

	heartbeatPkg := *kickPkg
	heartbeatPkg.PkgType = PkgTypeHeartbeat
	buffer5 := Encode(&heartbeatPkg)
	t.Logf("%b\n", buffer5)
	decode5 := Decode(buffer5)
	t.Logf("%+v\n", decode5)
	t.Logf("%+v\n", decode5.Msg)
}

func TestEncodeMsgId(t *testing.T) {
	buffer := make([]byte, 10)
	offset := encodeMsgId(0, buffer, 0)
	t.Logf("%b\n", buffer)
	t.Log(offset)
}
