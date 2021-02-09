package main

import (
	"net/http"

	"github.com/byuoitav/control-keys/codemap"
	"github.com/byuoitav/control-keys/handlers"
	"github.com/labstack/echo"
)

func main() {
	port := ":8029"
	router := echo.New()

	c := codemap.New()
	c.Start()
	h := handlers.New(c)

	// Functionality Endpoints
	router.GET("/:controlKey/getPreset", h.GetPresetHandler)
	router.GET("/:preset/getControlKey", h.GetControlKeyHandler)
	router.GET("/:preset/refresh", h.RefreshPresetKey)
	router.GET("/status", h.HealthCheck)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	_ = router.StartServer(&server)
}
