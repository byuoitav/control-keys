package main

import (
	"net/http"

	controlkeys "github.com/byuoitav/control-keys"
	"github.com/byuoitav/control-keys/internal/handlers"
	"github.com/byuoitav/control-keys/internal/keys"
	"github.com/labstack/echo"
)

func main() {
	port := ":8029"
	router := echo.New()

	var ds controlkeys.DataService

	keys := keys.New()
	h := handlers.Handlers{
		KeyService:  keys,
		DataService: ds,
	}

	// Functionality Endpoints
	router.GET("/:key/getPreset", h.GetPresetHandler)
	router.GET("/:preset/getControlKey", h.GetControlKeyHandler)
	router.GET("/:room/refresh", h.RefreshPresetKey)
	router.GET("/status", h.HealthCheck)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	_ = router.StartServer(&server)
}
