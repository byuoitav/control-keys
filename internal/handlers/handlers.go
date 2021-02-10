package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	controlkeys "github.com/byuoitav/control-keys"
	"github.com/labstack/echo"
)

type Handlers struct {
	KeyService  controlkeys.KeyService
	DataService controlkeys.DataService
}

type keyResponse struct {
	ControlKey string `json:"ControlKey"`
}

type presetResponse struct {
	RoomID     string `json:"RoomID"`
	PresetName string `json:"PresetName"`
}

//GetPresetHandler The endpoint to get the preset from the map
func (h *Handlers) GetPresetHandler(c echo.Context) error {
	controlKey := c.Param("key")

	cg, ok := h.KeyService.ControlGroup(c.Request().Context(), controlKey)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, presetResponse{
		RoomID:     cg.Room,
		PresetName: cg.ControlGroup,
	})
}

func (h *Handlers) GetControlKeyHandler(c echo.Context) error {
	split := strings.SplitN(c.Param("preset"), " ", 2)
	cg := controlkeys.ControlGroup{
		Room:         split[0],
		ControlGroup: split[1],
	}

	key, ok := h.KeyService.Key(c.Request().Context(), cg)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, keyResponse{ControlKey: key})
}

func (h *Handlers) RefreshPresetKey(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	cgs, err := h.DataService.ControlGroupsInRoom(ctx, c.Param("room"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get control groups for this room: %s", err))
	}

	for _, cg := range cgs {
		if err := h.KeyService.Refresh(ctx, cg); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to refresh controlGroup: %s", err))
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) HealthCheck(context echo.Context) error {
	return context.JSON(http.StatusOK, "Healthy!")
}
