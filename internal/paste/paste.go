package paste

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

// Predefined consts identify PGC behaviour.
const (
	GCItems    = 5
	GCDays     = 3
	GCSchedule = "@hourly"
)

// Predefined errors identify expected failure conditions.
var (
	ErrNotFound           = errors.New("Paste not found")
	ErrInvalidID          = errors.New("Paste UUID is invalid")
	ErrNotEnoughDataForGC = errors.New("Not enough data for GC")
)

type Storage struct {
	pasteList []Paste
}

// Add adds a Paste to the slice of Pastes (pasteList).
func (s *Storage) Add(np Paste) (*Paste, error) {
	s.pasteList = append(s.pasteList, np)
	return &np, nil
}

// List gets all Pastes from the slice of Pastes.
func (s *Storage) List() ([]Paste, error) {
	if s.pasteList == nil {
		return nil, ErrNotFound
	}
	return s.pasteList, nil
}

// Get retrieve a single paste identified by id.
func (s *Storage) Get(id string) (*Paste, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	for i, p := range s.pasteList {
		if p.ID == id {
			return &s.pasteList[i], nil

		}
	}
	return nil, ErrNotFound
}

// Kill purges list of Pastes.
func (s *Storage) Kill() ([]Paste, error) {
	if s.pasteList == nil {
		return nil, ErrNotFound
	}
	s.pasteList = nil
	return s.pasteList, nil
}

// PGCRun checks if Pastes' creation dates are older than `GCDays`
// and removes `GCItems` of pastes from the slice of Pastes.
func (s *Storage) PGCRun() (int, error) {
	if len(s.pasteList) <= GCItems {
		return 0, ErrNotEnoughDataForGC
	}

	today := time.Now()
	gcDate := today.Add(-24 * time.Hour * GCDays)

	if gcDate.After(s.pasteList[GCItems].Created) {
		log.Println("PGC found obsolete data: ", GCItems)
		s.pasteList = s.pasteList[GCItems:]
	}

	return GCItems, nil
}
