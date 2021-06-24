package main

import (
	"log"
	"os"

	"github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Register PGC cron.
	handlers.SchedulePGCCron()
	if err := run(); err != nil {
		log.Println("Shutting down", "error:", err)
		os.Exit(1)
	}
}

func run() error {
	// Register gin defaul router with Logger and Recover middleware.
	route := gin.Default()

	// Register handlers.
	route.GET("/pastes", handlers.ListPastes)
	route.GET("/pastes/:id", handlers.GetPaste)
	route.GET("/gc", handlers.PGC)
	route.POST("/pastes", handlers.AddPaste)
	route.DELETE("/killall", handlers.KillPastes)
	if err := route.Run(":8888"); err != nil {
		return err
	}
	return nil
}
