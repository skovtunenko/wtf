package http_test

import (
	"github.com/benbjohnson/wtf/http"
)

// Handler represents a test wrapper for http.Handler.
type Handler struct {
	*http.Handler

	DialHandler *DialHandler
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler:     &http.Handler{},
		DialHandler: NewDialHandler(),
	}
	h.Handler.DialHandler = h.DialHandler.DialHandler
	return h
}
