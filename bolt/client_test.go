package bolt_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/benbjohnson/wtf/bolt"
)

// Now is the mocked current time for testing.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Client is a test wrapper for bolt.Client.
type Client struct {
	*bolt.Client
}

// NewClient returns a new instance of Client pointing at a temporary file.
func NewClient() *Client {
	// Generate temporary filename.
	f, err := ioutil.TempFile("", "wtf-bolt-client-")
	if err != nil {
		panic(err)
	}
	f.Close()

	// Create client wrapper.
	c := &Client{
		Client: bolt.NewClient(),
	}
	c.Path = f.Name()
	c.Now = func() time.Time { return Now }

	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	defer os.Remove(c.Path)
	return c.Client.Close()
}
