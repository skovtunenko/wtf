package bolt

import (
	"encoding/binary"
	"time"

	"github.com/benbjohnson/wtf"
	"github.com/boltdb/bolt"
)

// Session represents an authenticable connection to the database.
type Session struct {
	db  *bolt.DB
	now time.Time

	// Authentication
	authenticator wtf.Authenticator
	authToken     string
	user          *wtf.User

	// Services
	dialService DialService
}

// newSession returns a new instance of Session attached to db.
func newSession(db *bolt.DB) *Session {
	s := &Session{db: db}
	s.dialService.session = s
	return s
}

// SetAuthToken sets token as the authentication token for the session.
func (s *Session) SetAuthToken(token string) {
	s.authToken = token
}

// Authenticate returns the current authenticate user.
func (s *Session) Authenticate() (*wtf.User, error) {
	// Return user if already authenticated.
	if s.user != nil {
		return s.user, nil
	}

	// Authenticate using token.
	u, err := s.authenticator.Authenticate(s.authToken)
	if err != nil {
		return nil, err
	}

	// Cache authenticated user.
	s.user = u

	return u, nil
}

// DialService returns a dial service associated with this session.
func (s *Session) DialService() wtf.DialService { return &s.dialService }

// itob returns an 8-byte big-endian encoded byte slice of v.
//
// This function is typically used for encoding integer IDs to byte slices
// so that they can be used as BoltDB keys.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
