package message

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	SetRouteDict(map[string]uint16{"mail.sendTo": 1})
	var request = &Message{
		Type:          Request,
		CompressRoute: false,
		CompressGzip:  true,
		Id:            100001,
		Route:         "user.login",
		Body:          []byte("{\"id\": 888,\"password\":\"abc\"}"),
	}
	var push = &Message{
		Type:          Push,
		CompressRoute: true,
		CompressGzip:  false,
		Id:            100001,
		Route:         "mail.sendTo",
		Body:          []byte("{\"id\": 888,\"password\":\"abc\"}"),
	}
	type args struct {
		m *Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Request", args: args{m: request}, wantErr: false},
		{name: "Push", args: args{m: push}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decode, err := Decode(got)
			if err != nil {
				t.Errorf("Dncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", decode)
			t.Logf("%s", string(decode.Body))
			if decode.Type != tt.args.m.Type || !bytes.Equal(decode.Body, tt.args.m.Body) {
				t.Errorf("Decode error")
				return
			}
		})
	}
}
