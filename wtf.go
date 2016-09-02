package wtf

import (
	"time"
)

// UserID represents a user identifier.
type UserID int

// User represents an authenticated user of the system.
type User struct {
	ID       UserID `json:"id"`
	Username string `json:"username"`
}

// DialID represents a dial identifier.
type DialID int

// Dial represents an adjustable WTF level associated with a user.
type Dial struct {
	ID      DialID    `json:"dialID"`
	UserID  UserID    `json:"userID"`
	Name    string    `json:"name,omitempty"`
	Level   float64   `json:"level"`
	ModTime time.Time `json:"modTime"`
}

// Authenticator represents a service for authenticating users.
type Authenticator interface {
	Authenticate(token string) (*User, error)
}

// Client creates a connection to the services.
type Client interface {
	Connect() Session
}

// Session represents authenticable connection to the services.
type Session interface {
	SetAuthToken(token string)
	DialService() DialService
}

// DialService represents a service for managing dials.
type DialService interface {
	Dial(id DialID) (*Dial, error)
	CreateDial(dial *Dial) error
	SetLevel(id DialID, level float64) error
}
