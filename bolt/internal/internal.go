package internal

import (
	"time"

	"github.com/benbjohnson/wtf"
	"github.com/gogo/protobuf/proto"
)

//go:generate protoc --gogo_out=. internal.proto

// MarshalDial encodes a dial to binary format.
func MarshalDial(d *wtf.Dial) ([]byte, error) {
	return proto.Marshal(&Dial{
		ID:      string(d.ID),
		Token:   d.Token,
		Name:    d.Name,
		Level:   d.Level,
		ModTime: d.ModTime.UnixNano(),
	})
}

// UnmarshalDial decodes a dial from a binary data.
func UnmarshalDial(data []byte, d *wtf.Dial) error {
	var pb Dial
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	d.ID = wtf.DialID(pb.ID)
	d.Token = pb.Token
	d.Name = pb.Name
	d.Level = pb.Level
	d.ModTime = time.Unix(0, pb.ModTime).UTC()

	return nil
}
