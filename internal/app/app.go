package app

import (
	"fmt"

	"github.com/bnema/yank-that/internal/clipboard"
	"github.com/bnema/yank-that/internal/path"
)

// PathResolver defines the interface for path resolution.
type PathResolver interface {
	Resolve(input string) (string, error)
}

// App orchestrates the yank-that path use case.
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

// Run executes the yank-that path logic.
// It resolves the input path and copies it to the clipboard.
// Returns the copied path on success.
func (a *App) Run(input string) (string, error) {
	absPath, err := a.resolver.Resolve(input)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	copier, err := clipboard.FirstAvailable(a.copiers)
	if err != nil {
		return "", err
	}

	if err := copier.Copy(absPath); err != nil {
		return "", fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return absPath, nil
}
