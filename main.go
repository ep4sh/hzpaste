package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
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

func main() {
	route := gin.Default()

	pastes := Pastes{}

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
		log.Fatalf("%+v\n", err)
	}

	c.JSON(200, gin.H{"pastes": all})
}

// -------------------- HANDLERS
func (ps *Pastes) Add(np Paste) (*Paste, error) {
	ps.PasteList = append(ps.PasteList, np)
	return &np, nil
}
