package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Paste struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Body    string    `json:"body"`
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
	np := Paste{
		Name:    c.PostForm("name"),
		Created: time.Now(),
		Body:    c.PostForm("body"),
	}
	all, err := ps.Add(np)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"pastes": all})
	}

	c.JSON(http.StatusOK, gin.H{"pastes": all})
}

func (ps *Pastes) ListPastes(c *gin.Context) {
	all, err := ps.List()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"pastes": all})
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
