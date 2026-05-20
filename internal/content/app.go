package content

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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

	if contentType, ok := imageContentType(content); ok {
		binaryCopier, err := a.firstAvailableBinaryCopier()
		if err != nil {
			return "", err
		}
		if err := binaryCopier.CopyBytes(content, contentType); err != nil {
			return "", fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		return absPath, nil
	}

	copier, err := clipboard.FirstAvailable(a.copiers)
	if err != nil {
		return "", err
	}

	if err := copier.Copy(string(content)); err != nil {
		return "", fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return absPath, nil
}

func (a *App) firstAvailableBinaryCopier() (clipboard.BinaryCopier, error) {
	for _, copier := range a.copiers {
		if !copier.Available() {
			continue
		}
		binaryCopier, ok := copier.(clipboard.BinaryCopier)
		if !ok {
			continue
		}
		return binaryCopier, nil
	}

	copier, err := clipboard.FirstAvailable(a.copiers)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("clipboard backend %q cannot copy binary data", copier.Name())
}

func imageContentType(content []byte) (string, bool) {
	contentType := http.DetectContentType(content)
	if strings.HasPrefix(contentType, "image/") {
		return contentType, true
	}

	return "", false
}
