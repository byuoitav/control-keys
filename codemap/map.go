package codemap

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
)

type CodeMap struct {
	controlKeys map[string]Preset
	m           sync.RWMutex
}

func New() *CodeMap {
	return &CodeMap{
		controlKeys: make(map[string]Preset),
		m:           sync.RWMutex{},
	}
}

func (c *CodeMap) Start() {
	go c.refreshMap()
}

// Preset struct
type Preset struct {
	RoomID     string
	PresetName string
}

func generateMap() (map[string]Preset, error) {
	//Query the DB for all of the UIConfigs
	uiConfigs, er := db.GetDB().GetAllUIConfigs()
	if er != nil {
		log.L.Errorf("error: %s", er)
		return nil, er
	}
	//create a map for Room/Preset
	m := make(map[string]Preset)
	for r := range uiConfigs {
		for p := range uiConfigs[r].Presets {
			code := generateCode()
			_, exists := m[code]
			for exists {
				code = generateCode()
				_, exists = m[code]
			}
			m[code] = Preset{
				RoomID:     uiConfigs[r].ID,
				PresetName: uiConfigs[r].Presets[p].Name,
			}
		}
	}

	return m, nil
}

func generateCode() string {
	min := 0
	max := 1000000
	code := strconv.Itoa(rand.Intn(max - min))
	//prepend it with zeros so every number selected is still a 6 digit number (i.e 1234 --> 001234)
	code = "000000" + code
	code = string(code[len(code)-6:])
	return code
}

// GetPresetFromMap function
func (c *CodeMap) GetPresetFromMap(code string) Preset {
	c.m.RLock()
	defer c.m.RUnlock()

	toReturn, ok := c.controlKeys[code]
	if !ok {
		return Preset{}
	}

	return toReturn
}

func (c *CodeMap) GetControlKeyFromPreset(preset Preset) string {
	c.m.RLock()
	defer c.m.RUnlock()

	for key, value := range c.controlKeys {
		if value == preset {
			return key
		}
	}

	return ""
}

func (c *CodeMap) RefreshControlKey(roomID string) bool {
	c.m.Lock()
	defer c.m.Unlock()

	var roomKeys []string

	for k, v := range c.controlKeys {
		if v.RoomID == roomID {
			roomKeys = append(roomKeys, k)
		}
	}

	if len(roomKeys) == 0 {
		// gonna assume it's not a valid preset
		return false
	}

	for _, key := range roomKeys {
		code := generateCode()
		_, exists := c.controlKeys[code]
		for exists {
			code = generateCode()
			_, exists = c.controlKeys[code]
		}

		preset := c.controlKeys[key]
		delete(c.controlKeys, key)
		c.controlKeys[code] = preset
	}

	return true
}

func (c *CodeMap) refreshMap() {
	newKeys, err := generateMap()
	for err != nil {
		time.Sleep(60 * time.Second)
		newKeys, err = generateMap()
	}
	c.m.Lock()
	c.controlKeys = newKeys
	c.m.Unlock()

	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		newKeys, err := generateMap()
		for err != nil {
			ticker.Reset(1 * time.Hour)
			time.Sleep(60 * time.Second)
			newKeys, err = generateMap()
		}
		c.m.Lock()
		c.controlKeys = newKeys
		c.m.Unlock()
	}
}
