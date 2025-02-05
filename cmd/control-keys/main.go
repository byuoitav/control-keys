package main

import (
	"net/http"

	"github.com/byuoitav/control-keys/codemap"
	"github.com/byuoitav/control-keys/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	port := ":8029"
	router := gin.Default()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Debug("Starting application")

	c := codemap.New(logger)
	c.Start()
	h := handlers.New(c)

	// Functionality Endpoints
	router.GET("/:param/*subpath", func(ctx *gin.Context) {
		subpath := ctx.Param("subpath") // This will capture everything after the first segment
		logger.Debug("Received request", zap.String("subpath", subpath))

		switch subpath {
		case "/getPreset":
			logger.Debug("Handling getPreset")
			h.GetPresetHandler(ctx)
		case "/getControlKey":
			logger.Debug("Handling getControlKey")
			h.GetControlKeyHandler(ctx)
		case "/refresh":
			logger.Debug("Handling refresh")
			h.RefreshPresetKey(ctx)
		default:
			logger.Debug("Invalid endpoint")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid endpoint"})
		}
	})
	router.GET("/status", func(ctx *gin.Context) {
		logger.Debug("Handling status check")
		h.HealthCheck(ctx)
	})

	server := &http.Server{
		Addr:           port,
		Handler:        router,
		MaxHeaderBytes: 1024 * 10,
	}

	logger.Debug("Starting server", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
