package wtf

// General errors.
const (
	ErrUnauthorized = Error("unauthorized")
	ErrInternal     = Error("internal error")
)

// Dial errors.
const (
	ErrDialRequired   = Error("dial required")
	ErrDialNotFound   = Error("dial not found")
	ErrDialExists     = Error("dial already exists")
	ErrDialIDRequired = Error("dial id required")
)

// Error represents a WTF error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
