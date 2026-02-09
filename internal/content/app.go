package content

import (
	"errors"
	"fmt"
	"os"

	"github.com/bnema/yank-that/internal/clipboard"
	"github.com/bnema/yank-that/internal/path"
)

// PathResolver defines the interface for path resolution.
type PathResolver interface {
	Resolve(input string) (string, error)
}

// App orchestrates the yank-that content use case.
type App struct {
	resolver PathResolver
	copiers  []clipboard.Copier
}

// New creates a new App with the given dependencies.
func New(resolver PathResolver, copiers []clipboard.Copier) *App {
	return &App{
		resolver: resolver,
		copiers:  copiers,
	}
}

// NewDefault creates a new App with default production configuration.
func NewDefault() *App {
	return New(
		path.NewResolver(),
		[]clipboard.Copier{
			clipboard.NewWayland(),
			clipboard.NewX11(),
		},
	)
}

// maxFileSize is the largest file we will read into memory (10 MB).
const maxFileSize = 10 * 1024 * 1024

// Run resolves the file path, reads file content, and copies it to clipboard.
func (a *App) Run(input string) (string, error) {
	absPath, err := a.resolver.Resolve(input)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	if info.Size() > maxFileSize {
		return "", fmt.Errorf("file too large (%d bytes, max %d)", info.Size(), maxFileSize)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	copier, err := a.findCopier()
	if err != nil {
		return "", err
	}

	if err := copier.Copy(string(content)); err != nil {
		return "", fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return absPath, nil
}

// findCopier returns the first available clipboard backend.
func (a *App) findCopier() (clipboard.Copier, error) {
	for _, c := range a.copiers {
		if c.Available() {
			return c, nil
		}
	}
	return nil, errors.New("no clipboard backend available (need wl-copy for Wayland or xclip for X11)")
}
