package controlkeys

import "context"

type DataService interface {
	// ControlGroups returns all of the ControlGroups this service should generate codes for.
	ControlGroups(context.Context) ([]ControlGroup, error)

	ControlGroupsInRoom(context.Context, string) ([]ControlGroup, error)
}

type ControlGroup struct {
	Room         string `json:"room"`
	ControlGroup string `json:"controlGroup"`
}

type KeyService interface {
	// ControlGroup returns the ControlGroup that the key is associated with.
	// If the key does not map to a ControlGroup, it returns false.
	ControlGroup(context.Context, string) (ControlGroup, bool)

	// Key returns the current key for the given controlGroup.
	// If no key exists, it returns false.
	Key(context.Context, ControlGroup) (string, bool)

	// Rebuild rebuilds the map with the given list of controlGroups
	Rebuild(context.Context, []ControlGroup) error

	// Refresh refreshes the key for the given controlGroup
	Refresh(context.Context, ControlGroup) error
}
