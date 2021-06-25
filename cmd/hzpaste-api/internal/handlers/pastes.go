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

// SchedulePGCCron schedules PGCRun() as cron job according to the GCSchedule.
func SchedulePGCCron(ps *paste.Storage) {
	log.Println("Registering new PGC cron job")
	cron := cron.New()
	cron.AddFunc(paste.GCSchedule, func() {
		log.Println("Scheduling PGC via cron")
		gcItems, err := ps.PGCRun()
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
func AddPasteH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var NewPaste paste.Paste
		if err := c.ShouldBindJSON(&NewPaste); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		NewPaste.ID = uuid.New().String()
		NewPaste.Created = time.Now().Add(-24 * time.Hour * 10)

		all, err := p.Add(NewPaste)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"pastes": all})
	}

	return gin.HandlerFunc(fn)
}

// ListPastes shows all pastes.
func ListPastesH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		all, err := p.List()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"pastes": all})
	}
	return gin.HandlerFunc(fn)
}

// GetPaste finds a single paste identified by an ID in the request URL.
func GetPasteH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		paste, err := p.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"paste": paste})
	}
	return gin.HandlerFunc(fn)
}

// KillPastes pugres all pastes.
func KillPastesH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		clearAll, err := p.Kill()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"pastes": clearAll})
	}
	return gin.HandlerFunc(fn)
}

// TODO: Remove handler, add goroutine.
// PGC (PasteGC) triggers obsolete data collection.
func PGCH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gcItems, err := p.PGCRun()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"gc_items": gcItems})
	}
	return gin.HandlerFunc(fn)
}
