package storage

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveFile(t *testing.T) {
	t.Run("ttl file", func(t *testing.T) {
		s := NewStorage(t.TempDir())
		data := []byte("hello")

		path, err := s.SaveFile(data, false)
		require.NoError(t, err)
		require.FileExists(t, path)

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, data, got)

		assert.Contains(t, path, filepath.Join("users"))
		assert.NotContains(t, path, "avatars")
	})

	t.Run("forever file", func(t *testing.T) {
		s := NewStorage(t.TempDir())
		data := []byte("avatar")

		path, err := s.SaveFile(data, true)
		require.NoError(t, err)
		require.FileExists(t, path)

		assert.Contains(t, path, filepath.Join("avatars"))
		assert.NotContains(t, path, filepath.Join("users"))
	})

	t.Run("path contains date parts", func(t *testing.T) {
		s := NewStorage(t.TempDir())
		now := time.Now()

		path, err := s.SaveFile([]byte("x"), false)
		require.NoError(t, err)

		assert.Contains(t, path, strconv.Itoa(now.Year()))
		assert.Contains(t, path, strconv.Itoa(int(now.Month())))
		assert.Contains(t, path, strconv.Itoa(now.Day()))
	})
}

func TestDelFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := NewStorage(t.TempDir())
		path, err := s.SaveFile([]byte("x"), false)
		require.NoError(t, err)

		require.NoError(t, s.DelFile(path))
		assert.NoFileExists(t, path)
	})

	t.Run("not existing", func(t *testing.T) {
		s := NewStorage(t.TempDir())
		err := s.DelFile(filepath.Join(t.TempDir(), "no-such"))
		require.Error(t, err)
	})
}

func TestRemoveOld(t *testing.T) {
	t.Run("removes old day dir", func(t *testing.T) {
		root := t.TempDir()
		s := NewStorage(root)

		old := time.Now().Add(-20 * 24 * time.Hour)
		oldDir := filepath.Join(
			s.usersDir(),
			strconv.Itoa(old.Year()),
			strconv.Itoa(int(old.Month())),
			strconv.Itoa(old.Day()),
			strconv.Itoa(old.Hour()),
		)
		require.NoError(t, os.MkdirAll(oldDir, 0o700))
		oldFile := filepath.Join(oldDir, "old-file")
		require.NoError(t, os.WriteFile(oldFile, []byte("old"), 0o600))

		// свежий файл через SaveFile
		freshPath, err := s.SaveFile([]byte("fresh"), false)
		require.NoError(t, err)

		// avatar не должен удаляться
		avatarPath, err := s.SaveFile([]byte("ava"), true)
		require.NoError(t, err)

		require.NoError(t, s.RemoveOld())

		assert.NoFileExists(t, oldFile)
		assert.FileExists(t, freshPath)
		assert.FileExists(t, avatarPath)
	})

	t.Run("keeps recent files", func(t *testing.T) {
		s := NewStorage(t.TempDir())
		path, err := s.SaveFile([]byte("recent"), false)
		require.NoError(t, err)

		require.NoError(t, s.RemoveOld())
		assert.FileExists(t, path)
	})
}
