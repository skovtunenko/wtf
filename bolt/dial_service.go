package bolt

import (
	"github.com/benbjohnson/wtf"
	"github.com/benbjohnson/wtf/bolt/internal"
)

// Ensure DialService implements wtf.DialService.
var _ wtf.DialService = &DialService{}

// DialService represents a service for managing dials.
type DialService struct {
	client *Client
}

// Dial returns a dial by ID.
func (s *DialService) Dial(id wtf.DialID) (*wtf.Dial, error) {
	// Start read-only transaction.
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal dial.
	var d wtf.Dial
	if v := tx.Bucket([]byte("Dials")).Get([]byte(id)); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalDial(v, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// CreateDial creates a new dial.
func (s *DialService) CreateDial(d *wtf.Dial) error {
	// Require object and id.
	if d == nil {
		return wtf.ErrDialRequired
	} else if d.ID == "" {
		return wtf.ErrDialIDRequired
	}

	// Start read-write transaction.
	tx, err := s.client.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify dial doesn't already exist.
	b := tx.Bucket([]byte("Dials"))
	if v := b.Get([]byte(d.ID)); v != nil {
		return wtf.ErrDialExists
	}

	// Update modified time.
	d.ModTime = s.client.Now()

	// Marshal and insert record.
	if v, err := internal.MarshalDial(d); err != nil {
		return err
	} else if err := b.Put([]byte(d.ID), v); err != nil {
		return err
	}

	return tx.Commit()
}

// SetLevel sets the current WTF level for a dial.
func (s *DialService) SetLevel(id wtf.DialID, token string, level float64) error {
	// Start read-write transaction.
	tx, err := s.client.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("Dials"))

	// Find and unmarshal record.
	var d wtf.Dial
	if v := b.Get([]byte(id)); v == nil {
		return wtf.ErrDialNotFound
	} else if err := internal.UnmarshalDial(v, &d); err != nil {
		return err
	}

	// Only update if token matches.
	if d.Token != token {
		return wtf.ErrUnauthorized
	}

	// Update dial level.
	d.Level = level
	d.ModTime = s.client.Now()

	// Marshal and insert record.
	if v, err := internal.MarshalDial(&d); err != nil {
		return err
	} else if err := b.Put([]byte(d.ID), v); err != nil {
		return err
	}

	return tx.Commit()
}
