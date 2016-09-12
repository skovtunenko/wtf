package internal_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benbjohnson/wtf"
	"github.com/benbjohnson/wtf/bolt/internal"
)

// Ensure dial can be marshaled and unmarshaled.
func TestMarshalDial(t *testing.T) {
	v := wtf.Dial{
		ID:      1,
		UserID:  2,
		Name:    "MYDIAL",
		Level:   10.2,
		ModTime: time.Now().UTC(),
	}

	var other wtf.Dial
	if buf, err := internal.MarshalDial(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.UnmarshalDial(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
