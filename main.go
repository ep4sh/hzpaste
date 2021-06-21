package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound           = errors.New("Paste not found")
	ErrInvalidID          = errors.New("Product UUID is invalid")
	ErrNotEnoughDataForGC = errors.New("Not enough data for GC")
)

const (
	GCItems = 5
	GCDays  = 3
)

type Paste struct {
	ID      string    `json:"id"`
	Name    string    `json:"name" binding:"required"`
	Created time.Time `json:"created"`
	Body    string    `json:"body" binding:"required"`
}

type Pastes struct {
	PasteList []Paste
}

func main() {
	route := gin.Default()

	pastes := Pastes{}

	route.GET("/pastes", pastes.ListPastes)
	route.GET("/pastes/:id", pastes.GetPaste)
	route.GET("/gc", pastes.GC)
	route.POST("/pastes", pastes.AddPaste)
	route.DELETE("/killall", pastes.KillPastes)
	route.Run(":8888")
}

// -------------------- HTTP HANDLERS
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

func (ps *Pastes) ListPastes(c *gin.Context) {
	all, err := ps.List()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pastes": all})
}

func (ps *Pastes) GetPaste(c *gin.Context) {
	id := c.Param("id")
	paste, err := ps.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"paste": paste})
}

func (ps *Pastes) KillPastes(c *gin.Context) {
	clearAll, err := ps.Kill()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"pastes": clearAll})
}

func (ps *Pastes) GC(c *gin.Context) {
	gcItems, err := ps.GCRun()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gc_items": gcItems})
}

// -------------------- HANDLERS
func (ps *Pastes) Add(np Paste) (*Paste, error) {
	ps.PasteList = append(ps.PasteList, np)
	return &np, nil
}

func (ps *Pastes) List() ([]Paste, error) {
	if ps.PasteList == nil {
		return nil, ErrNotFound
	}
	return ps.PasteList, nil
}

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

func (ps *Pastes) Kill() ([]Paste, error) {
	if ps.PasteList == nil {
		return nil, ErrNotFound
	}
	ps.PasteList = nil
	return ps.PasteList, nil
}

func (ps *Pastes) GCRun() (int, error) {
	if len(ps.PasteList) <= GCItems {
		return 0, ErrNotEnoughDataForGC
	}

	today := time.Now()
	gcDate := today.Add(-24 * time.Hour * GCDays)

	if gcDate.After(ps.PasteList[GCItems].Created) {
		log.Println("Cleanup old data..", ps.PasteList[GCItems])
		ps.PasteList = ps.PasteList[GCItems:]
	}

	return GCItems, nil
}
