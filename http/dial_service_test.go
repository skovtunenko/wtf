package http_test

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/benbjohnson/wtf"
	"github.com/benbjohnson/wtf/http"
	"github.com/benbjohnson/wtf/mock"
)

func TestDialService_Dial(t *testing.T) {
	t.Run("OK", testDialService_Dial)
	t.Run("NotFound", testDialService_Dial_NotFound)
	t.Run("ErrInternal", testDialService_Dial_ErrInternal)
}

// Ensure service can return a single dial.
func testDialService_Dial(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.DialHandler.DialService.DialFn = func(id wtf.DialID) (*wtf.Dial, error) {
		if id != "XXX" {
			t.Fatalf("unexpected id: %s", id)
		}
		return &wtf.Dial{ID: id, Name: "NAME", Level: 100, Token: "TOKEN", ModTime: Now}, nil
	}

	// Retrieve dial.
	d, err := c.DialService().Dial("XXX")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(d, &wtf.Dial{ID: "XXX", Name: "NAME", Level: 100, ModTime: Now}) {
		t.Fatalf("unexpected dial: %#v", d)
	}
}

// Ensure service handles fetching a non-existent dial.
func testDialService_Dial_NotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.DialHandler.DialService.DialFn = func(id wtf.DialID) (*wtf.Dial, error) {
		return nil, nil
	}

	// Retrieve dial.
	if d, err := c.DialService().Dial("NO_SUCH_DIAL"); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Fatal("expected nil dial")
	}
}

// Ensure service returns an internal error if an error occurs.
func testDialService_Dial_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.DialHandler.DialService.DialFn = func(id wtf.DialID) (*wtf.Dial, error) {
		return nil, errors.New("marker")
	}

	// Retrieve dial.
	if _, err := c.DialService().Dial("XXX"); err != wtf.ErrInternal {
		t.Fatal(err)
	} else if !strings.Contains(s.Handler.DialHandler.LogOutput.String(), "marker") {
		t.Fatalf("expected log output")
	}
}

func TestDialService_CreateDial(t *testing.T) {
	t.Run("OK", testDialService_CreateDial)
	t.Run("ErrDialRequired", testDialService_CreateDial_ErrDialRequired)
	t.Run("ErrDialIDRequired", testDialService_CreateDial_ErrDialIDRequired)
	t.Run("ErrDialExists", testDialService_CreateDial_ErrDialExists)
	t.Run("ErrInternal", testDialService_CreateDial_ErrInternal)
}

// Ensure service can create a new dial.
func testDialService_CreateDial(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.DialHandler.DialService.CreateDialFn = func(d *wtf.Dial) error {
		if !reflect.DeepEqual(d, &wtf.Dial{ID: "XXX", Token: "TOKEN", Name: "NAME", Level: 100}) {
			t.Fatalf("unexpected dial: %#v", d)
		}

		// Update mod time.
		d.ModTime = Now

		return nil
	}

	d := &wtf.Dial{ID: "XXX", Token: "TOKEN", Name: "NAME", Level: 100}

	// Create dial.
	err := c.DialService().CreateDial(d)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(d, &wtf.Dial{ID: "XXX", Token: "TOKEN", Name: "NAME", Level: 100, ModTime: Now}) {
		t.Fatalf("unexpected dial: %#v", d)
	}
}

// Ensure service returns an error if no dial is passed in.
func testDialService_CreateDial_ErrDialRequired(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	if err := c.DialService().CreateDial(nil); err != wtf.ErrDialRequired {
		t.Fatal(err)
	}
}

// Ensure service returns an error if dial id is passed blank.
func testDialService_CreateDial_ErrDialIDRequired(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	s.Handler.DialHandler.DialService.CreateDialFn = func(d *wtf.Dial) error {
		return wtf.ErrDialIDRequired
	}
	if err := c.DialService().CreateDial(&wtf.Dial{}); err != wtf.ErrDialIDRequired {
		t.Fatal(err)
	}
}

// Ensure service returns an error if dial already exists.
func testDialService_CreateDial_ErrDialExists(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	s.Handler.DialHandler.DialService.CreateDialFn = func(d *wtf.Dial) error {
		return wtf.ErrDialExists
	}
	if err := c.DialService().CreateDial(&wtf.Dial{ID: "XXX"}); err != wtf.ErrDialExists {
		t.Fatal(err)
	}
}

// Ensure service returns an error if an internal error occurs.
func testDialService_CreateDial_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	s.Handler.DialHandler.DialService.CreateDialFn = func(d *wtf.Dial) error {
		return errors.New("marker")
	}

	if err := c.DialService().CreateDial(&wtf.Dial{ID: "XXX"}); err != wtf.ErrInternal {
		t.Fatal(err)
	} else if !strings.Contains(s.Handler.DialHandler.LogOutput.String(), "marker") {
		t.Fatalf("expected log output")
	}
}

func TestDialService_SetLevel(t *testing.T) {
	t.Run("OK", testDialService_SetLevel)
	t.Run("ErrDialNotFound", testDialService_SetLevel_ErrDialNotFound)
	t.Run("ErrUnauthorized", testDialService_SetLevel_ErrUnauthorized)
	t.Run("ErrInternal", testDialService_SetLevel_ErrInternal)
}

// Ensure service can set the level of an existing dial.
func testDialService_SetLevel(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.DialHandler.DialService.SetLevelFn = func(id wtf.DialID, token string, level float64) error {
		if id != "XXX" {
			t.Fatalf("unexpected dial id: %s", id)
		} else if token != "TOKEN" {
			t.Fatalf("unexpected token: %s", token)
		} else if level != 100 {
			t.Fatalf("unexpected level: %v", level)
		}
		return nil
	}

	// Set dial level.
	err := c.DialService().SetLevel("XXX", "TOKEN", 100)
	if err != nil {
		t.Fatal(err)
	}
}

// Ensure service returns an error if the dial doesn't exist.
func testDialService_SetLevel_ErrDialNotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	s.Handler.DialHandler.DialService.SetLevelFn = func(id wtf.DialID, token string, level float64) error {
		return wtf.ErrDialNotFound
	}
	if err := c.DialService().SetLevel("XXX", "", 100); err != wtf.ErrDialNotFound {
		t.Fatal(err)
	}
}

// Ensure service returns an error if the user is not authorized.
func testDialService_SetLevel_ErrUnauthorized(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	s.Handler.DialHandler.DialService.SetLevelFn = func(id wtf.DialID, token string, level float64) error {
		return wtf.ErrUnauthorized
	}
	if err := c.DialService().SetLevel("XXX", "", 100); err != wtf.ErrUnauthorized {
		t.Fatal(err)
	}
}

// Ensure service returns an error if an internal error occurs.
func testDialService_SetLevel_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()
	s.Handler.DialHandler.DialService.SetLevelFn = func(id wtf.DialID, token string, level float64) error {
		return errors.New("marker")
	}

	if err := c.DialService().SetLevel("XXX", "", 100); err != wtf.ErrInternal {
		t.Fatal(err)
	} else if !strings.Contains(s.Handler.DialHandler.LogOutput.String(), "marker") {
		t.Fatalf("expected log output")
	}
}

// DialHandler represents a test wrapper for http.DialHandler.
type DialHandler struct {
	*http.DialHandler

	DialService mock.DialService
	LogOutput   bytes.Buffer
}

// NewDialHandler returns a new instance of DialHandler.
func NewDialHandler() *DialHandler {
	h := &DialHandler{DialHandler: http.NewDialHandler()}
	h.DialHandler.DialService = &h.DialService
	h.Logger = log.New(VerboseWriter(&h.LogOutput), "", log.LstdFlags)
	return h
}
