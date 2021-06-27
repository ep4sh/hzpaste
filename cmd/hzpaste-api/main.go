package main

import (
	"log"
	"os"

	_ "github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/docs"
	"github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/handlers"
	"github.com/ep4sh/hzpaste/internal/paste"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @contact.name API Support
// @contact.url http://ep4sh.cc
// @contact.email ep4sh2k@gm[a]il.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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

	// use ginSwagger middleware to serve the API docs
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := route.Run(":8888"); err != nil {
		return err
	}
	return nil
}
