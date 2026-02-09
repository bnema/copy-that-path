package aliases

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	startMarker = "# >>> yank-that aliases >>>"
	endMarker   = "# <<< yank-that aliases <<<"
)

type Installer struct {
	homeDir func() (string, error)
}

func NewInstaller() *Installer {
	return &Installer{homeDir: os.UserHomeDir}
}

func (i *Installer) Install(shell string) ([]string, error) {
	home, err := i.homeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to detect home directory: %w", err)
	}

	switch shell {
	case "", "all":
		updates := make([]string, 0, 3)
		for _, targetShell := range []string{"bash", "zsh", "fish"} {
			result, installErr := i.installShell(home, targetShell)
			if installErr != nil {
				return nil, installErr
			}
			updates = append(updates, result)
		}
		return updates, nil
	case "bash", "zsh", "fish":
		result, installErr := i.installShell(home, shell)
		if installErr != nil {
			return nil, installErr
		}
		return []string{result}, nil
	default:
		return nil, fmt.Errorf("unsupported shell %q (use bash, zsh, fish, or all)", shell)
	}
}

func (i *Installer) installShell(home, shell string) (string, error) {
	var path string
	var snippet string

	switch shell {
	case "bash":
		path = filepath.Join(home, ".bashrc")
		snippet = bashSnippet()
	case "zsh":
		path = filepath.Join(home, ".zshrc")
		snippet = zshSnippet()
	case "fish":
		path = filepath.Join(home, ".config", "fish", "config.fish")
		snippet = fishSnippet()
	default:
		return "", errors.New("unsupported shell")
	}

	changed, err := appendSnippet(path, snippet)
	if err != nil {
		return "", fmt.Errorf("failed to install aliases for %s: %w", shell, err)
	}

	if changed {
		return fmt.Sprintf("installed %s aliases in %s", shell, path), nil
	}
	return fmt.Sprintf("%s aliases already installed in %s", shell, path), nil
}

func appendSnippet(path, snippet string) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}

	content := ""
	data, err := os.ReadFile(path)
	if err == nil {
		content = string(data)
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	hasStart := strings.Contains(content, startMarker)
	hasEnd := strings.Contains(content, endMarker)
	if hasStart && hasEnd {
		return false, nil
	}
	if hasStart || hasEnd {
		return false, fmt.Errorf("corrupt alias block in %s: found only one marker; remove the partial block and re-run", path)
	}

	newContent := content
	if len(newContent) > 0 && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += "\n" + snippet

	perm := os.FileMode(0o644)
	if info, statErr := os.Stat(path); statErr == nil {
		perm = info.Mode().Perm()
	}

	if err := os.WriteFile(path, []byte(newContent), perm); err != nil {
		return false, err
	}

	return true, nil
}

func bashSnippet() string {
	return strings.Join([]string{
		startMarker,
		"alias ytp='yt path'",
		"alias ytc='yt content'",
		"alias ytg='yt github'",
		endMarker,
		"",
	}, "\n")
}

func zshSnippet() string {
	return strings.Join([]string{
		startMarker,
		"alias ytp='yt path'",
		"alias ytc='yt content'",
		"alias ytg='yt github'",
		endMarker,
		"",
	}, "\n")
}

func fishSnippet() string {
	return strings.Join([]string{
		startMarker,
		"alias ytp 'yt path'",
		"alias ytc 'yt content'",
		"alias ytg 'yt github'",
		endMarker,
		"",
	}, "\n")
}
