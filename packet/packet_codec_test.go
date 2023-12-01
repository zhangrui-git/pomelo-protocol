package packet

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var handshake = &Packet{Type: Handshake, Data: nil}
	var handshakeAck = &Packet{Type: HandshakeAck, Data: nil}
	var heartbeat = &Packet{Type: Heartbeat, Data: nil}
	var kick = &Packet{Type: Kick, Data: nil}
	var data = &Packet{Type: Data, Data: []byte("{\"nickname\": \"google\", \"fight\": false, \"level\": 100}")}

	type args struct {
		p *Packet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Handshake", args: args{handshake}, wantErr: false},
		{name: "HandshakeAck", args: args{handshakeAck}, wantErr: false},
		{name: "Heartbeat", args: args{heartbeat}, wantErr: false},
		{name: "Kick", args: args{kick}, wantErr: false},
		{name: "Data", args: args{data}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decode, err := Decode(got)
			if err != nil {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", decode)
			t.Logf("%s", string(decode.Data))
			if decode.Type != tt.args.p.Type || !bytes.Equal(decode.Data, tt.args.p.Data) {
				t.Errorf("Decode error")
				return
			}
		})
	}
}
