package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

var (
	ErrNotFound = errors.New("Paste not found")
)

func main() {
	route := gin.Default()

	pastes := Pastes{}

	route.GET("/pastes", pastes.ListPastes)
	route.POST("/pastes", pastes.AddPaste)
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
	NewPaste.Created = time.Now()
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

// -------------------- HANDLERS
func (ps *Pastes) Add(np Paste) (*Paste, error) {
	ps.PasteList = append(ps.PasteList, np)
	return &np, nil
}

func (ps *Pastes) List() ([]Paste, error) {
	var EmptyPastes []Paste
	if ps.PasteList == nil {
		return EmptyPastes, ErrNotFound
	}
	return ps.PasteList, nil
}
