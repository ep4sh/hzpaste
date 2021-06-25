package main

import (
	"log"
	"os"

	"github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/handlers"
	"github.com/ep4sh/hzpaste/internal/paste"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Println("Shutting down", "error:", err)
		os.Exit(1)
	}
}

func run() error {
	// Register gin defaul router with Logger and Recover middleware.
	route := gin.Default()

	// Register PasteList.
	pastes := paste.Storage{}

	// Register PGC cron task.
	handlers.SchedulePGCCron(&pastes)

	// Register handlers.
	route.GET("/pastes", handlers.ListPastesH(&pastes))
	route.GET("/pastes/:id", handlers.GetPasteH(&pastes))
	route.GET("/gc", handlers.PGCH(&pastes))
	route.POST("/pastes", handlers.AddPasteH(&pastes))
	route.DELETE("/killall", handlers.KillPastesH(&pastes))
	if err := route.Run(":8888"); err != nil {
		return err
	}
	return nil
}
