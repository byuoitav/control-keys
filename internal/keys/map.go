package keys

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	controlkeys "github.com/byuoitav/control-keys"
)

const (
	_maxSize = (999999 - 100000) / 4
)

type Map struct {
	controlGroups map[string]controlGroup
	keys          map[controlGroup]string

	sync.RWMutex
}

func New() *Map {
	return &Map{
		controlGroups: make(map[string]controlGroup),
		keys:          make(map[controlGroup]string),
	}
}

type controlGroup struct {
	Room         string
	ControlGroup string
}

func (m *Map) Rebuild(ctx context.Context, cgs []controlkeys.ControlGroup) error {
	if len(cgs) > _maxSize {
		return errors.New("can't support that many controlGroups")
	}

	m.Lock()
	defer m.Unlock()

	// keep all of the current controlGroups keys
	controlGroups := make(map[string]controlGroup, len(m.controlGroups))
	for key, v := range m.controlGroups {
		if !strings.HasPrefix(key, "old-") {
			controlGroups["old-"+key] = v
		}
	}

	genKey := func() string {
		for {
			key := newKey()
			if _, ok := controlGroups[key]; ok {
				continue
			}

			if _, ok := controlGroups["old-"+key]; ok {
				continue
			}

			return key
		}
	}

	keys := make(map[controlGroup]string, len(m.keys))
	for _, cg := range cgs {
		key := genKey()
		tmp := controlGroup{
			Room:         cg.Room,
			ControlGroup: cg.ControlGroup,
		}

		controlGroups[key] = tmp
		keys[tmp] = key
	}

	m.controlGroups = controlGroups
	m.keys = keys
	return nil
}

func (m *Map) ControlGroup(ctx context.Context, key string) (controlkeys.ControlGroup, bool) {
	m.RLock()
	defer m.RUnlock()

	cg, ok := m.controlGroups[key]
	if !ok {
		// check if this key is a previous key
		cg, ok = m.controlGroups["old-"+key]
	}

	return cg.convert(), ok
}

func (m *Map) Key(ctx context.Context, cg controlkeys.ControlGroup) (string, bool) {
	m.RLock()
	defer m.RUnlock()

	tmp := controlGroup{
		Room:         cg.Room,
		ControlGroup: cg.ControlGroup,
	}

	key, ok := m.keys[tmp]
	return key, ok
}

func (m *Map) Refresh(ctx context.Context, cg controlkeys.ControlGroup) error {
	if len(m.controlGroups)+1 > _maxSize {
		return errors.New("can't support that many controlGroups")
	}

	// find a valid key
	key := newKey()
	for {
		if _, ok := m.ControlGroup(ctx, key); !ok {
			break
		}

		key = newKey()
	}

	tmp := controlGroup{
		Room:         cg.Room,
		ControlGroup: cg.ControlGroup,
	}

	m.Lock()
	defer m.Unlock()

	m.controlGroups[key] = tmp
	m.keys[tmp] = key

	return nil
}

func (c controlGroup) convert() controlkeys.ControlGroup {
	return controlkeys.ControlGroup{
		Room:         c.Room,
		ControlGroup: c.ControlGroup,
	}
}

// newKey generates a 6 digit code
func newKey() string {
	min := 100000
	max := 999999
	return strconv.Itoa(rand.Intn(max-min) + min)
}
