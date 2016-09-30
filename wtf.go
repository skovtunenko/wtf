package wtf

import (
	"time"
)

// DialID represents a dial identifier.
type DialID string

// Dial represents an adjustable WTF level associated with a user.
// A user-defined token is used to loosely authenticate update requests.
type Dial struct {
	ID      DialID    `json:"dialID"`
	Token   string    `json:"-"`
	Name    string    `json:"name,omitempty"`
	Level   float64   `json:"level"`
	ModTime time.Time `json:"modTime"`
}

// Client creates a connection to the services.
type Client interface {
	DialService() DialService
}

// DialService represents a service for managing dials.
type DialService interface {
	Dial(id DialID) (*Dial, error)
	CreateDial(dial *Dial) error
	SetLevel(id DialID, token string, level float64) error
}
