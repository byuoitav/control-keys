package main

import (
	"net/http"

	"github.com/byuoitav/control-keys/codemap"
	"github.com/byuoitav/control-keys/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	port := ":8029"
	router := gin.Default()

	c := codemap.New()
	c.Start()
	h := handlers.New(c)

	// Functionality Endpoints
	router.GET("/:controlKey/getPreset", func(ctx *gin.Context) {
		h.GetPresetHandler(ctx)
	})
	router.GET("/:preset/getControlKey", func(ctx *gin.Context) {
		h.GetControlKeyHandler(ctx)
	})
	router.GET("/:room/refresh", func(ctx *gin.Context) {
		h.RefreshPresetKey(ctx)
	})
	router.GET("/status", func(ctx *gin.Context) {
		h.HealthCheck(ctx)
	})

	server := &http.Server{
		Addr:           port,
		Handler:        router,
		MaxHeaderBytes: 1024 * 10,
	}

	_ = server.ListenAndServe()
}
