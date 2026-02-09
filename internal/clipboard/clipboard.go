package clipboard

import "errors"

// Copier defines the interface for clipboard operations.
type Copier interface {
	// Copy copies the given text to the system clipboard.
	Copy(text string) error

	// Available returns true if this clipboard backend is available.
	Available() bool

	// Name returns the name of this clipboard backend.
	Name() string
}

// FirstAvailable returns the first available Copier from the list,
// or an error if none are available.
func FirstAvailable(copiers []Copier) (Copier, error) {
	for _, c := range copiers {
		if c.Available() {
			return c, nil
		}
	}
	return nil, errors.New("no clipboard backend available (need wl-copy for Wayland or xclip for X11)")
}
