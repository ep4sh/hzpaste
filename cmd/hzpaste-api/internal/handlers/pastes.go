package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ep4sh/hzpaste/internal/paste"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron"
)

// Pastes  defines handlers related to pastes. It hold app state
// needed by the handler methods
var Pastes paste.Storage

// SchedulePGCCron schedules PGCRun() as cron job according to the GCSchedule.
func SchedulePGCCron() {
	log.Println("Registering new PGC cron job")
	cron := cron.New()
	cron.AddFunc(paste.GCSchedule, func() {
		log.Println("Scheduling PGC via cron")
		gcItems, err := Pastes.PGCRun()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("PGC cleaned up obsolete data: ", gcItems)
	})
	cron.Start()
}

// -------------------- HTTP HANDLERS.

// AddPaste validates the body of a request to create a new paste.
func AddPaste(c *gin.Context) {
	var NewPaste paste.Paste
	if err := c.ShouldBindJSON(&NewPaste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	NewPaste.ID = uuid.New().String()
	NewPaste.Created = time.Now().Add(-24 * time.Hour * 10)

	all, err := Pastes.Add(NewPaste)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"pastes": all})
}

// ListPastes shows all pastes.
func ListPastes(c *gin.Context) {
	all, err := Pastes.List()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pastes": all})
}

// GetPaste finds a single paste identified by an ID in the request URL.
func GetPaste(c *gin.Context) {
	id := c.Param("id")
	paste, err := Pastes.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"paste": paste})
}

// KillPastes pugres all pastes.
func KillPastes(c *gin.Context) {
	clearAll, err := Pastes.Kill()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"pastes": clearAll})
}

// TODO: Remove handler, add goroutine.
// PGC (PasteGC) triggers obsolete data collection.
func PGC(c *gin.Context) {
	gcItems, err := Pastes.PGCRun()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gc_items": gcItems})
}
