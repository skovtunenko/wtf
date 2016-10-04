package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/benbjohnson/wtf"
	"github.com/julienschmidt/httprouter"
)

// DialHandler represents an HTTP API handler for dials.
type DialHandler struct {
	*httprouter.Router

	DialService wtf.DialService

	Logger *log.Logger
}

// NewDialHandler returns a new instance of DialHandler.
func NewDialHandler() *DialHandler {
	h := &DialHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.POST("/api/dials", h.handlePostDial)
	h.GET("/api/dials/:id", h.handleGetDial)
	h.PATCH("/api/dials/:id", h.handlePatchDial)
	return h
}

// handleGetDial handles requests to fetch a single dial.
func (h *DialHandler) handleGetDial(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	// Find dial by ID.
	d, err := h.DialService.Dial(wtf.DialID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if d == nil {
		NotFound(w)
	} else {
		encodeJSON(w, &getDialResponse{Dial: d}, h.Logger)
	}
}

type getDialResponse struct {
	Dial *wtf.Dial `json:"dial,omitempty"`
	Err  string    `json:"err,omitempty"`
}

// handleGetDial handles requests to create a new dial.
func (h *DialHandler) handlePostDial(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request.
	var req postDialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	d := req.Dial
	d.Token = req.Token
	d.ModTime = time.Time{}

	// Create dial.
	switch err := h.DialService.CreateDial(d); err {
	case nil:
		encodeJSON(w, &postDialResponse{Dial: d}, h.Logger)
	case wtf.ErrDialRequired, wtf.ErrDialIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case wtf.ErrDialExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type postDialRequest struct {
	Dial  *wtf.Dial `json:"dial,omitempty"`
	Token string    `json:"token,omitempty"`
}

type postDialResponse struct {
	Dial *wtf.Dial `json:"dial,omitempty"`
	Err  string    `json:"err,omitempty"`
}

// handlePatchDial handles requests to update a dial level.
func (h *DialHandler) handlePatchDial(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request.
	var req patchDialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	// Create dial.
	switch err := h.DialService.SetLevel(req.ID, req.Token, req.Level); err {
	case nil:
		encodeJSON(w, &patchDialResponse{}, h.Logger)
	case wtf.ErrDialNotFound:
		Error(w, err, http.StatusNotFound, h.Logger)
	case wtf.ErrUnauthorized:
		Error(w, err, http.StatusUnauthorized, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type patchDialRequest struct {
	ID    wtf.DialID `json:"id"`
	Token string     `json:"token"`
	Level float64    `json:"level"`
}

type patchDialResponse struct {
	Err string `json:"err,omitempty"`
}

// Ensure service implements interface.
var _ wtf.DialService = &DialService{}

// DialService represents an HTTP implementation of wtf.DialService.
type DialService struct {
	URL *url.URL
}

// Dial returns a dial by id.
func (s *DialService) Dial(id wtf.DialID) (*wtf.Dial, error) {
	u := *s.URL
	u.Path = "/api/dials/" + url.QueryEscape(string(id))

	// Execute request.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody getDialResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Err != "" {
		return nil, wtf.Error(respBody.Err)
	}
	return respBody.Dial, nil
}

// CreateDial creates a new dial.
func (s *DialService) CreateDial(d *wtf.Dial) error {
	// Validate arguments.
	if d == nil {
		return wtf.ErrDialRequired
	}

	u := *s.URL
	u.Path = "/api/dials"

	// Save token.
	token := d.Token

	// Encode request body.
	reqBody, err := json.Marshal(postDialRequest{Dial: d, Token: token})
	if err != nil {
		return err
	}

	// Execute request.
	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody postDialResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Err != "" {
		return wtf.Error(respBody.Err)
	}

	// Copy returned dial.
	*d = *respBody.Dial
	d.Token = token

	return nil
}

// SetLevel sets the level of an existing dial.
func (s *DialService) SetLevel(id wtf.DialID, token string, level float64) error {
	u := *s.URL
	u.Path = "/api/dials/" + url.QueryEscape(string(id))

	// Encode request body.
	reqBody, err := json.Marshal(patchDialRequest{ID: id, Token: token, Level: level})
	if err != nil {
		return err
	}

	// Create request.
	req, err := http.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	// Execute request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody postDialResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Err != "" {
		return wtf.Error(respBody.Err)
	}

	return nil
}
