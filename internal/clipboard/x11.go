package clipboard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// X11 implements Copier using xclip.
type X11 struct{}

// NewX11 creates a new X11 clipboard copier.
func NewX11() *X11 {
	return &X11{}
}

// Copy copies text to clipboard using xclip.
func (x *X11) Copy(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xclip failed: %w", err)
	}
	return nil
}

// CopyBytes copies typed binary data to clipboard using xclip.
func (x *X11) CopyBytes(data []byte, contentType string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-t", contentType)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xclip failed: %w", err)
	}
	return nil
}

// Available returns true if X11 is available.
func (x *X11) Available() bool {
	// Check if DISPLAY is set
	if os.Getenv("DISPLAY") == "" {
		return false
	}

	// Check if xclip is installed
	_, err := exec.LookPath("xclip")
	return err == nil
}

// Name returns the backend name.
func (x *X11) Name() string {
	return "x11"
}
