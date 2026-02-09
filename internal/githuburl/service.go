package githuburl

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bnema/yank-that/internal/clipboard"
	"github.com/bnema/yank-that/internal/path"
	"github.com/cli/go-gh/v2/pkg/repository"
)

type PathResolver interface {
	Resolve(input string) (string, error)
}

type Service struct {
	resolver    PathResolver
	copiers     []clipboard.Copier
	repoCurrent func() (repository.Repository, error)
	repoRoot    func() (string, error)
	currentRef  func() (string, error)
}

func NewDefault() *Service {
	return &Service{
		resolver:    path.NewResolver(),
		copiers:     []clipboard.Copier{clipboard.NewWayland(), clipboard.NewX11()},
		repoCurrent: repository.Current,
		repoRoot: func() (string, error) {
			return gitOutput("rev-parse", "--show-toplevel")
		},
		currentRef: func() (string, error) {
			ref, err := gitOutput("symbolic-ref", "--short", "HEAD")
			if err == nil {
				return ref, nil
			}
			return gitOutput("rev-parse", "HEAD")
		},
	}
}

func (s *Service) Run(input string) (string, error) {
	builtURL, err := s.BuildURL(input)
	if err != nil {
		return "", err
	}

	copier, err := s.findCopier()
	if err != nil {
		return "", err
	}

	if err := copier.Copy(builtURL); err != nil {
		return "", fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return builtURL, nil
}

func (s *Service) BuildURL(input string) (string, error) {
	pathInput, line, hasLine, err := parseInput(input)
	if err != nil {
		return "", err
	}

	absPath, err := s.resolver.Resolve(pathInput)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	if info.IsDir() {
		return "", errors.New("path points to a directory, expected a file")
	}

	repo, err := s.repoCurrent()
	if err != nil {
		return "", fmt.Errorf("failed to detect repository: %w", err)
	}

	root, err := s.repoRoot()
	if err != nil {
		return "", fmt.Errorf("failed to detect repository root: %w", err)
	}

	ref, err := s.currentRef()
	if err != nil {
		return "", fmt.Errorf("failed to detect git ref: %w", err)
	}

	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to compute repository-relative path: %w", err)
	}
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return "", errors.New("path is outside repository root")
	}

	host := repo.Host
	if host == "" {
		host = "github.com"
	}

	encodedRef := encodePath(filepath.ToSlash(ref))
	encodedRelPath := encodePath(filepath.ToSlash(relPath))
	builtURL := fmt.Sprintf("https://%s/%s/%s/blob/%s/%s", host, repo.Owner, repo.Name, encodedRef, encodedRelPath)
	if hasLine {
		builtURL += fmt.Sprintf("#L%d", line)
	}

	return builtURL, nil
}

func (s *Service) findCopier() (clipboard.Copier, error) {
	for _, c := range s.copiers {
		if c.Available() {
			return c, nil
		}
	}
	return nil, errors.New("no clipboard backend available (need wl-copy for Wayland or xclip for X11)")
}

func parseInput(input string) (string, int, bool, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", 0, false, errors.New("input cannot be empty")
	}

	idx := strings.LastIndex(input, ":")
	if idx == -1 {
		return input, 0, false, nil
	}

	// Skip Windows drive-letter colons (e.g. C:\file.txt)
	if idx == 1 && len(input) > 2 && (input[0] >= 'A' && input[0] <= 'Z' || input[0] >= 'a' && input[0] <= 'z') && (input[2] == '\\' || input[2] == '/') {
		return input, 0, false, nil
	}

	pathPart := input[:idx]
	linePart := input[idx+1:]
	if pathPart == "" || linePart == "" {
		return input, 0, false, nil
	}

	line, err := strconv.Atoi(linePart)
	if err != nil || line <= 0 {
		return input, 0, false, nil
	}

	return pathPart, line, true, nil
}

func encodePath(value string) string {
	parts := strings.Split(value, "/")
	encoded := make([]string, 0, len(parts))
	for _, part := range parts {
		encoded = append(encoded, url.PathEscape(part))
	}
	return strings.Join(encoded, "/")
}

func gitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			return "", fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
		}
		return "", fmt.Errorf("git %s failed: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(string(out)), nil
}
