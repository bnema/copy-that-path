package aliases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstaller_Install_All(t *testing.T) {
	tmp := t.TempDir()
	installer := NewInstaller()
	installer.homeDir = func() (string, error) {
		return tmp, nil
	}

	updates, err := installer.Install("all")
	require.NoError(t, err)
	require.Len(t, updates, 3)

	bashBytes, readErr := os.ReadFile(filepath.Join(tmp, ".bashrc"))
	require.NoError(t, readErr)
	assert.Contains(t, string(bashBytes), "alias ytp='yt path'")

	zshBytes, readErr := os.ReadFile(filepath.Join(tmp, ".zshrc"))
	require.NoError(t, readErr)
	assert.Contains(t, string(zshBytes), "alias ytc='yt content'")

	fishBytes, readErr := os.ReadFile(filepath.Join(tmp, ".config", "fish", "config.fish"))
	require.NoError(t, readErr)
	assert.Contains(t, string(fishBytes), "alias ytg 'yt github'")
}

func TestInstaller_Install_IsIdempotent(t *testing.T) {
	tmp := t.TempDir()
	installer := NewInstaller()
	installer.homeDir = func() (string, error) {
		return tmp, nil
	}

	_, err := installer.Install("bash")
	require.NoError(t, err)

	updates, err := installer.Install("bash")
	require.NoError(t, err)
	require.Len(t, updates, 1)
	assert.Contains(t, updates[0], "already installed")
}

func TestInstaller_Install_InvalidShell(t *testing.T) {
	tmp := t.TempDir()
	installer := NewInstaller()
	installer.homeDir = func() (string, error) {
		return tmp, nil
	}

	_, err := installer.Install("tcsh")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}
