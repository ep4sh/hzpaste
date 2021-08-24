package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ep4sh/hzpaste/internal/paste"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
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

// AddPasteH godoc
// @Summary add a paste
// @Description add new paste
// @Accept  json
// @Produce  json
// @Success 200 {object} paste.Paste "new_paste"
// @Router /pastes [post]
// AddPasteH returns a handler that validates the body of a request to create
// a new paste.
func AddPasteH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var NewPaste paste.Paste
		if err := c.ShouldBindJSON(&NewPaste); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		NewPaste.ID = uuid.New().String()
		NewPaste.Created = time.Now()

		np, err := p.Add(NewPaste)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"new_paste": np})
	}

	return gin.HandlerFunc(fn)
}

// ListPastesH godoc
// @Summary shows all pastes
// @Description get pastes list
// @Accept  json
// @Produce  json
// @Success 200 {array} paste.Paste "pastes"
// @Router /pastes [get]
// ListPastesH returns a handler that shows all pastes.
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

// GetPasteH godoc
// @Summary finds a single paste identified by an ID in the request URL.
// @Description get paste by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path string true "Paste ID"
// @Success 200 {object} paste.Paste "paste"
// @Router /pastes/{id} [get]
// GetPasteH returns a handler that finds a single paste identified by an
// ID in the request URL.
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

// KillPastesH godoc
// @Summary purges all pastes
// @Description returns empty slice for Pastes
// @Produce  json
// @Success 200 {array} paste.Paste "pastes"
// @Router /pastes [delete]
// KillPastesH returns a handler that pugres all pastes by clearing the
// pastes list.
func KillPastesH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		_, err := p.Kill()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"pastes": "null"})
	}
	return gin.HandlerFunc(fn)
}

// PGCH godoc
// @Summary triggers obsolete data collection and removes it
// @Description returns count of the removed data
// @Produce  json
// @Success 200 {array} int "gc"
// @Router /gc [get]
// PGCH returns a handler that triggers obsolete data collection.
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

// PingH godoc
// @Summary provides health checks
// @Description provides health checks
// @Accept  json
// @Produce  json
// @Success 200 {object} bool "OK"
// @Router /ping [get]
// GetPasteH returns a handler that provides health check
func PingH(p *paste.Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ready, err := p.Ping()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ready": ready})
	}
	return gin.HandlerFunc(fn)
}
