package handlers

import (
	"net/http"
	"strings"

	"github.com/byuoitav/control-keys/codemap"
	"github.com/labstack/echo"
)

type Handler struct {
	c *codemap.CodeMap
}

type ControlKey struct {
	ControlKey string
}

func New(c *codemap.CodeMap) *Handler {
	return &Handler{
		c: c,
	}
}

//GetPresetHandler The endpoint to get the preset from the map
func (h *Handler) GetPresetHandler(context echo.Context) error {
	controlKey := context.Param("controlKey")
	preset := h.c.GetPresetFromMap(controlKey)
	if preset == (codemap.Preset{}) {
		return context.JSON(http.StatusNotFound, "The preset was not found for this control key")
	}
	return context.JSON(http.StatusOK, preset)
}

func (h *Handler) GetControlKeyHandler(context echo.Context) error {
	presetParam := context.Param("preset")
	presetParts := strings.SplitN(presetParam, " ", 2)
	preset := codemap.Preset{
		RoomID:     presetParts[0],
		PresetName: presetParts[1],
	}

	key := h.c.GetControlKeyFromPreset(preset)
	if key == "" {
		return context.JSON(http.StatusNotFound, "The control key was not found for this preset")
	}

	return context.JSON(http.StatusOK, ControlKey{ControlKey: key})
}

func (h *Handler) RefreshPresetKey(context echo.Context) error {
	roomID := context.Param("room")

	ok := h.c.RefreshControlKey(roomID)

	if !ok {
		return context.JSON(http.StatusNotFound, "Invalid preset")
	}

	return context.NoContent(http.StatusOK)
}

func (h *Handler) HealthCheck(context echo.Context) error {
	return context.JSON(http.StatusOK, "Healthy!")
}
