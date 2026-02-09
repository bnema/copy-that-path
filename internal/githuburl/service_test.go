package githuburl

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bnema/yank-that/internal/clipboard"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type staticResolver struct {
	value string
	err   error
}

func (s *staticResolver) Resolve(input string) (string, error) {
	return s.value, s.err
}

func TestService_BuildURL_FileOnly(t *testing.T) {
	repoRoot := t.TempDir()
	filePath := filepath.Join(repoRoot, "README.md")
	require.NoError(t, os.WriteFile(filePath, []byte("hi"), 0o644))

	s := &Service{
		resolver: &staticResolver{value: filePath},
		repoCurrent: func() (repository.Repository, error) {
			return repository.Repository{Host: "github.com", Owner: "bnema", Name: "yank-that"}, nil
		},
		repoRoot:   func() (string, error) { return repoRoot, nil },
		currentRef: func() (string, error) { return "main", nil },
	}

	actual, err := s.BuildURL("README.md")
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/bnema/yank-that/blob/main/README.md", actual)
}

func TestService_BuildURL_WithLine(t *testing.T) {
	repoRoot := t.TempDir()
	filePath := filepath.Join(repoRoot, "internal", "app", "app.go")
	require.NoError(t, os.MkdirAll(filepath.Dir(filePath), 0o755))
	require.NoError(t, os.WriteFile(filePath, []byte("package app"), 0o644))

	s := &Service{
		resolver: &staticResolver{value: filePath},
		repoCurrent: func() (repository.Repository, error) {
			return repository.Repository{Host: "github.com", Owner: "bnema", Name: "yank-that"}, nil
		},
		repoRoot:   func() (string, error) { return repoRoot, nil },
		currentRef: func() (string, error) { return "feature/ref", nil },
	}

	actual, err := s.BuildURL("internal/app/app.go:42")
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/bnema/yank-that/blob/feature/ref/internal/app/app.go#L42", actual)
}

func TestService_BuildURL_PathOutsideRepo(t *testing.T) {
	repoRoot := t.TempDir()
	otherRoot := t.TempDir()
	filePath := filepath.Join(otherRoot, "README.md")
	require.NoError(t, os.WriteFile(filePath, []byte("hi"), 0o644))

	s := &Service{
		resolver: &staticResolver{value: filePath},
		repoCurrent: func() (repository.Repository, error) {
			return repository.Repository{Host: "github.com", Owner: "bnema", Name: "yank-that"}, nil
		},
		repoRoot:   func() (string, error) { return repoRoot, nil },
		currentRef: func() (string, error) { return "main", nil },
	}

	_, err := s.BuildURL("README.md")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "outside repository root")
}

func TestService_Run_CopiesBuiltURL(t *testing.T) {
	repoRoot := t.TempDir()
	filePath := filepath.Join(repoRoot, "README.md")
	require.NoError(t, os.WriteFile(filePath, []byte("hi"), 0o644))

	copier := clipboard.NewMockCopier(t)
	copier.EXPECT().Available().Return(true)
	copier.EXPECT().Copy("https://github.com/bnema/yank-that/blob/main/README.md").Return(nil)

	s := &Service{
		resolver: &staticResolver{value: filePath},
		copiers:  []clipboard.Copier{copier},
		repoCurrent: func() (repository.Repository, error) {
			return repository.Repository{Host: "github.com", Owner: "bnema", Name: "yank-that"}, nil
		},
		repoRoot:   func() (string, error) { return repoRoot, nil },
		currentRef: func() (string, error) { return "main", nil },
	}

	actual, err := s.Run("README.md")
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/bnema/yank-that/blob/main/README.md", actual)
}

func TestService_Run_NoClipboardBackend(t *testing.T) {
	repoRoot := t.TempDir()
	filePath := filepath.Join(repoRoot, "README.md")
	require.NoError(t, os.WriteFile(filePath, []byte("hi"), 0o644))

	copier := clipboard.NewMockCopier(t)
	copier.EXPECT().Available().Return(false)

	s := &Service{
		resolver: &staticResolver{value: filePath},
		copiers:  []clipboard.Copier{copier},
		repoCurrent: func() (repository.Repository, error) {
			return repository.Repository{Host: "github.com", Owner: "bnema", Name: "yank-that"}, nil
		},
		repoRoot:   func() (string, error) { return repoRoot, nil },
		currentRef: func() (string, error) { return "main", nil },
	}

	_, err := s.Run("README.md")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no clipboard backend available")
}

func TestService_Run_CopyFails(t *testing.T) {
	repoRoot := t.TempDir()
	filePath := filepath.Join(repoRoot, "README.md")
	require.NoError(t, os.WriteFile(filePath, []byte("hi"), 0o644))

	copier := clipboard.NewMockCopier(t)
	copier.EXPECT().Available().Return(true)
	copier.EXPECT().Copy("https://github.com/bnema/yank-that/blob/main/README.md").Return(errors.New("copy failed"))

	s := &Service{
		resolver: &staticResolver{value: filePath},
		copiers:  []clipboard.Copier{copier},
		repoCurrent: func() (repository.Repository, error) {
			return repository.Repository{Host: "github.com", Owner: "bnema", Name: "yank-that"}, nil
		},
		repoRoot:   func() (string, error) { return repoRoot, nil },
		currentRef: func() (string, error) { return "main", nil },
	}

	_, err := s.Run("README.md")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to copy to clipboard")
}
