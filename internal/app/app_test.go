package app

import (
	"errors"
	"testing"

	"github.com/bnema/yank-that/internal/clipboard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockResolver struct {
	path string
	err  error
}

func (m *mockResolver) Resolve(input string) (string, error) {
	return m.path, m.err
}

func TestApp_Run(t *testing.T) {
	t.Run("copies resolved path to clipboard", func(t *testing.T) {
		resolver := &mockResolver{path: "/home/user/test.txt"}
		copier := clipboard.NewMockCopier(t)
		copier.EXPECT().Available().Return(true)
		copier.EXPECT().Copy("/home/user/test.txt").Return(nil)

		app := New(resolver, []clipboard.Copier{copier})
		result, err := app.Run("test.txt")

		require.NoError(t, err)
		assert.Equal(t, "/home/user/test.txt", result)
	})

	t.Run("returns error when resolver fails", func(t *testing.T) {
		resolver := &mockResolver{err: errors.New("resolve error")}
		copier := clipboard.NewMockCopier(t)

		app := New(resolver, []clipboard.Copier{copier})
		_, err := app.Run("test.txt")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to resolve path")
	})

	t.Run("returns error when no clipboard backend available", func(t *testing.T) {
		resolver := &mockResolver{path: "/home/user/test.txt"}
		copier := clipboard.NewMockCopier(t)
		copier.EXPECT().Available().Return(false)

		app := New(resolver, []clipboard.Copier{copier})
		_, err := app.Run("test.txt")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no clipboard backend available")
	})

	t.Run("returns error when copy fails", func(t *testing.T) {
		resolver := &mockResolver{path: "/home/user/test.txt"}
		copier := clipboard.NewMockCopier(t)
		copier.EXPECT().Available().Return(true)
		copier.EXPECT().Copy("/home/user/test.txt").Return(errors.New("copy failed"))

		app := New(resolver, []clipboard.Copier{copier})
		_, err := app.Run("test.txt")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to copy to clipboard")
	})

	t.Run("uses first available copier", func(t *testing.T) {
		resolver := &mockResolver{path: "/tmp/file"}
		unavailable := clipboard.NewMockCopier(t)
		unavailable.EXPECT().Available().Return(false)
		available := clipboard.NewMockCopier(t)
		available.EXPECT().Available().Return(true)
		available.EXPECT().Copy("/tmp/file").Return(nil)

		app := New(resolver, []clipboard.Copier{unavailable, available})
		result, err := app.Run("file")

		require.NoError(t, err)
		assert.Equal(t, "/tmp/file", result)
	})
}
