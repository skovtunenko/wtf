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
	s := c.Connect().DialService()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*wtf.User, error) {
		return &wtf.User{ID: 123}, nil
	}

	dial := wtf.Dial{
		Name:  "MY DIAL",
		Level: 50,
	}

	// Create new dial.
	if err := s.CreateDial(&dial); err != nil {
		t.Fatal(err)
	} else if dial.ID != 1 {
		t.Fatalf("unexpected id: %d", dial.ID)
	} else if dial.UserID != 123 {
		t.Fatalf("unexpected user id: %d", dial.UserID)
	}

	// Retrieve dial and compare.
	other, err := s.Dial(1)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&dial, other) {
		t.Fatalf("unexpected dial: %#v", other)
	}
}

// Ensure dial's level can be updated.
func TestDialService_SetLevel(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().DialService()

	// Create new dials
	if err := s.CreateDial(&wtf.Dial{Level: 50}); err != nil {
		t.Fatal(err)
	} else if err := s.CreateDial(&wtf.Dial{Level: 80}); err != nil {
		t.Fatal(err)
	}

	// Update dial levels.
	if err := s.SetLevel(1, 60); err != nil {
		t.Fatal(err)
	} else if err := s.SetLevel(2, 10); err != nil {
		t.Fatal(err)
	}

	// Verify dial 1 updated.
	if d, err := s.Dial(1); err != nil {
		t.Fatal(err)
	} else if d.Level != 60 {
		t.Fatalf("unexpected dial #1 level: %f", d.Level)
	}

	// Verify dial 2 updated.
	if d, err := s.Dial(2); err != nil {
		t.Fatal(err)
	} else if d.Level != 10 {
		t.Fatalf("unexpected dial #2 level: %f", d.Level)
	}
}

// Ensure error is returned if an unauthorized user updates the level.
func TestDialService_SetLevel_ErrUnauthorized(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	// Connect in one session and create dial.
	session0 := c.Connect()
	if err := session0.DialService().CreateDial(&wtf.Dial{Level: 50}); err != nil {
		t.Fatal(err)
	}

	// Connect in a different session with a different user and attempt to update.
	c.Authenticator.AuthenticateFn = func(token string) (*wtf.User, error) {
		return &wtf.User{ID: 100000}, nil
	}
	session1 := c.Connect()
	if err := session1.DialService().SetLevel(1, 20); err != wtf.ErrUnauthorized {
		t.Fatalf("unexpected error: %s", err)
	}
}
