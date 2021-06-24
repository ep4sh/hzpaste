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
	GCSchedule = "*/10 * * * *"
)

// Predefined errors identify expected failure conditions.
var (
	ErrNotFound           = errors.New("Paste not found")
	ErrInvalidID          = errors.New("Product UUID is invalid")
	ErrNotEnoughDataForGC = errors.New("Not enough data for GC")
)

type Storage struct {
	PasteList []Paste
}

// Add adds a Paste to the slice of Pastes (PasteList).
func (s *Storage) Add(np Paste) (*Paste, error) {
	s.PasteList = append(s.PasteList, np)
	return &np, nil
}

// List gets all Pastes from the slice of Pastes.
func (s *Storage) List() ([]Paste, error) {
	if s.PasteList == nil {
		return nil, ErrNotFound
	}
	return s.PasteList, nil
}

// Get retrieve a single paste identified by id.
func (s *Storage) Get(id string) (*Paste, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	for i, p := range s.PasteList {
		if p.ID == id {
			return &s.PasteList[i], nil

		}
	}
	return nil, ErrNotFound
}

// Kill purges list of Pastes.
func (s *Storage) Kill() ([]Paste, error) {
	if s.PasteList == nil {
		return nil, ErrNotFound
	}
	s.PasteList = nil
	return s.PasteList, nil
}

// PGCRun checks if Pastes' creation dates are older than `GCDays`
// and removes `GCItems` of pastes from the slice of Pastes.
func (s *Storage) PGCRun() (int, error) {
	if len(s.PasteList) <= GCItems {
		return 0, ErrNotEnoughDataForGC
	}

	today := time.Now()
	gcDate := today.Add(-24 * time.Hour * GCDays)

	if gcDate.After(s.PasteList[GCItems].Created) {
		log.Println("PGC found obsolete data: ", GCItems)
		s.PasteList = s.PasteList[GCItems:]
	}

	return GCItems, nil
}
