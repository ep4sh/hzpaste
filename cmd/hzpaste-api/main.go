package main

import (
	"fmt"

	_ "github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/docs"
	"github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/handlers"
	"github.com/ep4sh/hzpaste/cmd/hzpaste-api/internal/hzconfig"
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

func setupRouter() *gin.Engine {
	// Register gin defaul router with Logger and Recover middleware.
	router := gin.Default()

	// Register PasteList.
	pastes := paste.Storage{}

	// Register PGC cron task.
	handlers.SchedulePGCCron(&pastes)

	// Register handlers.
	router.GET("/pastes", handlers.ListPastesH(&pastes))
	router.GET("/pastes/:id", handlers.GetPasteH(&pastes))
	router.GET("/gc", handlers.PGCH(&pastes))
	router.GET("/ping", handlers.PingH(&pastes))
	router.POST("/pastes", handlers.AddPasteH(&pastes))
	router.DELETE("/killall", handlers.KillPastesH(&pastes))

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func main() {
	config := hzconfig.InitConfig()
	endpoint := fmt.Sprintf("%s:%s", config.Host, config.Port)
	router := setupRouter()
	router.Run(endpoint)
}
