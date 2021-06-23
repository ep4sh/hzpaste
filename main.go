package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// Predefined errors identify expected failure conditions.
var (
	ErrNotFound           = errors.New("Paste not found")
	ErrInvalidID          = errors.New("Product UUID is invalid")
	ErrNotEnoughDataForGC = errors.New("Not enough data for GC")
)

// Predefined consts identify PGC behaviour.
const (
	GCItems    = 5
	GCDays     = 3
	GCSchedule = "* * * * *"
)

// Paste is an item we store.
type Paste struct {
	ID      string    `json:"id"`
	Name    string    `json:"name" binding:"required"`
	Created time.Time `json:"created"`
	Body    string    `json:"body" binding:"required"`
}

// Paste is an collection of items we store.
type Pastes struct {
	PasteList []Paste
}

func main() {
	// Register gin defaul router with Logger and Recover middleware.
	route := gin.Default()

	// Register PasteList.
	pastes := Pastes{}
	// Register PGC cron.
	pastes.CheckPGCCron()

	// Register handlers.
	route.GET("/pastes", pastes.ListPastes)
	route.GET("/pastes/:id", pastes.GetPaste)
	route.GET("/gc", pastes.PGC)
	route.POST("/pastes", pastes.AddPaste)
	route.DELETE("/killall", pastes.KillPastes)
	route.Run(":8888")
}

// CheckPGCCron schedules PGCRun() as cron job according to the GCSchedule.
func (ps *Pastes) CheckPGCCron() {
	log.Println("Registering new PGC cron job")
	c := cron.New()
	c.AddFunc(GCSchedule, func() {
		log.Println("Scheduling PGC via cron")
		gcItems, err := ps.PGCRun()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("PGC cleaned up obsolete data: ", gcItems)
	})
	c.Start()
}

// -------------------- HTTP HANDLERS.

// AddPaste validates the body of a request to create a new paste.
func (ps *Pastes) AddPaste(c *gin.Context) {
	var NewPaste Paste
	if err := c.ShouldBindJSON(&NewPaste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	NewPaste.ID = uuid.New().String()
	NewPaste.Created = time.Now().Add(-24 * time.Hour * 10)

	all, err := ps.Add(NewPaste)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"pastes": all})
}

// ListPastes shows all pastes.
func (ps *Pastes) ListPastes(c *gin.Context) {
	all, err := ps.List()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pastes": all})
}

// GetPaste finds a single paste identified by an ID in the request URL.
func (ps *Pastes) GetPaste(c *gin.Context) {
	id := c.Param("id")
	paste, err := ps.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"paste": paste})
}

// KillPastes pugres all pastes.
func (ps *Pastes) KillPastes(c *gin.Context) {
	clearAll, err := ps.Kill()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"pastes": clearAll})
}

// TODO: Remove handler, add goroutine.
// PGC (PasteGC) triggers obsolete data collection.
func (ps *Pastes) PGC(c *gin.Context) {
	gcItems, err := ps.PGCRun()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gc_items": gcItems})
}

// -------------------- Paste HANDLERS
// Add adds a Paste to the slice of Pastes (PasteList).
func (ps *Pastes) Add(np Paste) (*Paste, error) {
	ps.PasteList = append(ps.PasteList, np)
	return &np, nil
}

// List gets all Pastes from the slice of Pastes.
func (ps *Pastes) List() ([]Paste, error) {
	if ps.PasteList == nil {
		return nil, ErrNotFound
	}
	return ps.PasteList, nil
}

// Get retrieve a single paste identified by id.
func (ps *Pastes) Get(id string) (*Paste, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	for i, p := range ps.PasteList {
		if p.ID == id {
			return &ps.PasteList[i], nil

		}
	}
	return nil, ErrNotFound
}

// Kill purges list of Pastes.
func (ps *Pastes) Kill() ([]Paste, error) {
	if ps.PasteList == nil {
		return nil, ErrNotFound
	}
	ps.PasteList = nil
	return ps.PasteList, nil
}

// PGCRun checks if Pastes' creation dates are older than `GCDays`
// and removes `GCItems` of pastes from the slice of Pastes.
func (ps *Pastes) PGCRun() (int, error) {
	if len(ps.PasteList) <= GCItems {
		return 0, ErrNotEnoughDataForGC
	}

	today := time.Now()
	gcDate := today.Add(-24 * time.Hour * GCDays)

	if gcDate.After(ps.PasteList[GCItems].Created) {
		log.Println("PGC found obsolete data: ", GCItems)
		ps.PasteList = ps.PasteList[GCItems:]
	}

	return GCItems, nil
}
