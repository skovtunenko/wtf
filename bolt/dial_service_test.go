package bolt_test

import (
	"reflect"
	"testing"

	"github.com/benbjohnson/wtf"
)

// Ensure dial can be created and retrieved.
func TestDialService_CreateDial(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DialService()

	dial := wtf.Dial{
		ID:    "XXX",
		Token: "YYY",
		Name:  "MY DIAL",
		Level: 50,
	}

	// Create new dial.
	if err := s.CreateDial(&dial); err != nil {
		t.Fatal(err)
	}

	// Retrieve dial and compare.
	other, err := s.Dial("XXX")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&dial, other) {
		t.Fatalf("unexpected dial: %#v", other)
	}
}

// Ensure dial validates the id.
func TestDialService_CreateDial_ErrDialIDRequired(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	if err := c.DialService().CreateDial(&wtf.Dial{ID: ""}); err != wtf.ErrDialIDRequired {
		t.Fatal(err)
	}
}

// Ensure duplicate dials are not allowed.
func TestDialService_CreateDial_ErrDialExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	if err := c.DialService().CreateDial(&wtf.Dial{ID: "X"}); err != nil {
		t.Fatal(err)
	}
	if err := c.DialService().CreateDial(&wtf.Dial{ID: "X"}); err != wtf.ErrDialExists {
		t.Fatal(err)
	}
}

// Ensure dial's level can be updated.
func TestDialService_SetLevel(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DialService()

	// Create new dials
	if err := s.CreateDial(&wtf.Dial{ID: "XXX", Token: "YYY", Level: 50}); err != nil {
		t.Fatal(err)
	} else if err := s.CreateDial(&wtf.Dial{ID: "AAA", Token: "BBB", Level: 80}); err != nil {
		t.Fatal(err)
	}

	// Update dial levels.
	if err := s.SetLevel("XXX", "YYY", 60); err != nil {
		t.Fatal(err)
	} else if err := s.SetLevel("AAA", "BBB", 10); err != nil {
		t.Fatal(err)
	}

	// Verify dial updated.
	if d, err := s.Dial("XXX"); err != nil {
		t.Fatal(err)
	} else if d.Level != 60 {
		t.Fatalf("unexpected dial #1 level: %f", d.Level)
	}

	// Verify dial 2 updated.
	if d, err := s.Dial("AAA"); err != nil {
		t.Fatal(err)
	} else if d.Level != 10 {
		t.Fatalf("unexpected dial #2 level: %f", d.Level)
	}
}

// Ensure dial level cannot be updated if token doesn't match.
func TestDialService_SetLevel_ErrUnauthorized(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DialService()

	// Create new dial.
	if err := s.CreateDial(&wtf.Dial{ID: "XXX", Token: "YYY", Level: 50}); err != nil {
		t.Fatal(err)
	}

	// Update dial level with wrong token.
	if err := s.SetLevel("XXX", "BAD_TOKEN", 60); err != wtf.ErrUnauthorized {
		t.Fatal(err)
	}
}

// Ensure error is returned if dial doesn't exist.
func TestDialService_SetLevel_ErrDialNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	if err := c.DialService().SetLevel("XXX", "", 50); err != wtf.ErrDialNotFound {
		t.Fatal(err)
	}
}
