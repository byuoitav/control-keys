package handlers

import (
	"net/http"
	"strings"

	"github.com/byuoitav/control-keys/codemap"
	"github.com/gin-gonic/gin"
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

// GetPresetHandler The endpoint to get the preset from the map
func (h *Handler) GetPresetHandler(context *gin.Context) {
	controlKey := context.Param("controlKey")
	preset := h.c.GetPresetFromMap(controlKey)
	if preset == (codemap.Preset{}) {
		context.JSON(http.StatusNotFound, gin.H{"message": "The preset was not found for this control key"})
		return
	}
	context.JSON(http.StatusOK, preset)
}

func (h *Handler) GetControlKeyHandler(context *gin.Context) {
	presetParam := context.Param("preset")
	presetParts := strings.SplitN(presetParam, " ", 2)
	preset := codemap.Preset{
		RoomID:     presetParts[0],
		PresetName: presetParts[1],
	}

	key := h.c.GetControlKeyFromPreset(preset)
	if key == "" {
		context.JSON(http.StatusNotFound, gin.H{"message": "The control key was not found for this preset"})
		return
	}

	context.JSON(http.StatusOK, ControlKey{ControlKey: key})
}

func (h *Handler) RefreshPresetKey(context *gin.Context) {
	roomID := context.Param("room")

	ok := h.c.RefreshControlKey(roomID)

	if !ok {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid preset"})
		return
	}

	context.Status(http.StatusOK)
}

func (h *Handler) HealthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Healthy!"})
}
