package clipboard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Wayland implements Copier using wl-copy.
type Wayland struct{}

// NewWayland creates a new Wayland clipboard copier.
func NewWayland() *Wayland {
	return &Wayland{}
}

// Copy copies text to clipboard using wl-copy.
func (w *Wayland) Copy(text string) error {
	cmd := exec.Command("wl-copy", "--trim-newline")
	cmd.Stdin = strings.NewReader(text)
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wl-copy failed: %w", err)
	}
	return nil
}

// CopyBytes copies typed binary data to clipboard using wl-copy.
func (w *Wayland) CopyBytes(data []byte, contentType string) error {
	cmd := exec.Command("wl-copy", "--type", contentType)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wl-copy failed: %w", err)
	}
	return nil
}

// Available returns true if Wayland is available.
func (w *Wayland) Available() bool {
	// Check if WAYLAND_DISPLAY is set
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		return false
	}

	// Check if wl-copy is installed
	_, err := exec.LookPath("wl-copy")
	return err == nil
}

// Name returns the backend name.
func (w *Wayland) Name() string {
	return "wayland"
}
