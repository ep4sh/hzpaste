package paste

import "time"

// Paste is an item we store.
type Paste struct {
	ID      string    `json:"id"`
	Name    string    `json:"name" binding:"required"`
	Created time.Time `json:"created"`
	Body    string    `json:"body" binding:"required"`
}
