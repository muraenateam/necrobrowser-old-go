package server

import (
	"github.com/muraenateam/necrobrowser/log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"github.com/muraenateam/necrobrowser/core"
)

func SetupRouter(options *core.Options) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	if *options.Debug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	if *options.Debug {
		router.Use(gin.Logger())
	}

	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Middleware
	authToken, err := uuid.FromString(*options.AuthToken)
	if err != nil {
		log.Fatal("%v", err)
	}
	router.Use(AuthMiddleware(authToken))

	// Share options
	router.Use(OptionsMiddleware(options))

	// Routes
	router.POST("/instrument/:token", NewBrowserHandler)
	router.POST("/instrumentKnown/:token", KnownBrowserHandler)
	router.GET("/status/:token", StatusHandler)

	return router
}
