package fs

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.RemoveAll("test_delete")
	})
	assert.NoError(t, os.Mkdir("test_delete", 0o755))
	assert.NoError(t, os.Mkdir("test_delete/dir", 0o755))
	createTestFile(t, "test_delete/dir/file.txt", 0x755, time.Now(), []byte{})
	createTestFile(t, "test_delete/file.txt", 0x755, time.Now(), []byte{})

	operation := NewOperation()

	t.Run("ExistingFile", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, operation.Delete("test_delete/file.txt"))
	})

	t.Run("ExistingDirectory", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, operation.Delete("test_delete/dir"))
	})

	t.Run("NonExistingFile", func(t *testing.T) {
		t.Parallel()

		err := operation.Delete("test_delete/not_exists_file.txt")
		assert.Error(t, err)
		assert.Equal(t, os.ErrNotExist, err)
	})

	t.Run("NonExistingDirectory", func(t *testing.T) {
		t.Parallel()

		err := operation.Delete("test_delete/not_exists_dir")
		assert.Error(t, err)
		assert.Equal(t, os.ErrNotExist, err)
	})
}

func TestCopy(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.RemoveAll("test_create")
	})
	assert.NoError(t, os.Mkdir("test_create", 0o755))
	assert.NoError(t, os.Mkdir("test_create/a", 0o755))
	assert.NoError(t, os.Mkdir("test_create/b", 0o755))
	assert.NoError(t, os.Mkdir("test_create/a/dir", 0o755))
	createTestFile(t, "test_create/a/dir/file.txt", 0o755, time.Now(), []byte{})
	createTestFile(t, "test_create/a/file.txt", 0o755, time.Now(), []byte{})

	operation := NewOperation()

	t.Run("ExistingFile", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, operation.Copy("test_create/a/file.txt", "test_create/b/file.txt"))

		srcState, err := os.Stat("test_create/a/file.txt")
		assert.NoError(t, err)
		dstState, err := os.Stat("test_create/b/file.txt")
		assert.NoError(t, err)

		assert.Equal(t, srcState.Mode(), dstState.Mode())
		assert.Equal(t, srcState.ModTime(), dstState.ModTime())
	})

	t.Run("ExistingDirectory", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, operation.Copy("test_create/a/dir", "test_create/b/dir"))

		srcState, err := os.Stat("test_create/a/dir")
		assert.NoError(t, err)
		dstState, err := os.Stat("test_create/b/dir")
		assert.NoError(t, err)

		assert.Equal(t, srcState.Mode(), dstState.Mode())
	})

	t.Run("NonExistingFile", func(t *testing.T) {
		t.Parallel()

		err := operation.Copy("test_delete/not_exists_file.txt", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("NonExistingDirectory", func(t *testing.T) {
		t.Parallel()

		err := operation.Copy("test_delete/not_exists_dir", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestExists(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.RemoveAll("test_exists")
	})
	assert.NoError(t, os.Mkdir("test_exists", 0o755))
	assert.NoError(t, os.Mkdir("test_exists/dir", 0o755))
	createTestFile(t, "test_exists/file.txt", 0x755, time.Now(), []byte{})

	operation := NewOperation()

	t.Run("ExistingFile", func(t *testing.T) {
		t.Parallel()

		exists := operation.Exists("test_exists/file.txt")
		assert.True(t, exists)
	})

	t.Run("ExistingDirectory", func(t *testing.T) {
		t.Parallel()

		exists := operation.Exists("test_exists/dir")
		assert.True(t, exists)
	})

	t.Run("NonExistingFile", func(t *testing.T) {
		t.Parallel()

		exists := operation.Exists("test_exists/not_exists_file.txt")
		assert.False(t, exists)
	})

	t.Run("NonExistingDirectory", func(t *testing.T) {
		t.Parallel()

		exists := operation.Exists("test_exists/not_exists_dir")
		assert.False(t, exists)
	})
}

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.RemoveAll("test_equal")
	})
	assert.NoError(t, os.Mkdir("test_equal", 0o755))
	assert.NoError(t, os.Mkdir("test_equal/a", 0o755))
	assert.NoError(t, os.Mkdir("test_equal/b", 0o755))
	assert.NoError(t, os.Mkdir("test_equal/c", 0o777))
	assert.NoError(t, os.Chmod("test_equal/c", 0o777)) // need to set mode manually

	now := time.Now()
	createTestFile(t, "test_equal/a.txt", 0o755, now, []byte("test"))
	createTestFile(t, "test_equal/a_.txt", 0o755, now, []byte("test"))
	createTestFile(t, "test_equal/b.txt", 0o755, now, []byte("test123"))
	createTestFile(t, "test_equal/c.txt", 0o755, now.Add(time.Hour), []byte("test"))
	createTestFile(t, "test_equal/d.txt", 0o777, now, []byte("test"))

	operation := NewOperation()

	t.Run("FileEqual", func(t *testing.T) {
		t.Parallel()

		assert.True(t, operation.Equal("test_equal/a.txt", "test_equal/a_.txt"))
	})

	t.Run("FileWithDifferentSize", func(t *testing.T) {
		t.Parallel()

		assert.False(t, operation.Equal("test_equal/a.txt", "test_equal/b.txt"))
	})

	t.Run("FileWithDifferentModTime", func(t *testing.T) {
		t.Parallel()

		assert.False(t, operation.Equal("test_equal/a.txt", "test_equal/c.txt"))
	})

	t.Run("FileWithDifferentMod", func(t *testing.T) {
		t.Parallel()

		assert.False(t, operation.Equal("test_equal/a.txt", "test_equal/d.txt"))
	})

	t.Run("FileWithNotExisting", func(t *testing.T) {
		t.Parallel()

		assert.False(t, operation.Equal("test_equal/a.txt", "test_equal/not_exists.txt"))
	})

	t.Run("DirectoryEqual", func(t *testing.T) {
		t.Parallel()

		assert.True(t, operation.Equal("test_equal/a", "test_equal/b"))
	})

	t.Run("DirectoryWithDifferentMod", func(t *testing.T) {
		t.Parallel()

		assert.False(t, operation.Equal("test_equal/a", "test_equal/c"))
	})

	t.Run("DirectoryWithNotExisting", func(t *testing.T) {
		t.Parallel()

		assert.False(t, operation.Equal("test_equal/a", "test_equal/not_exists"))
	})
}

func createTestFile(t *testing.T, path string, mod fs.FileMode, modTime time.Time, content []byte) {
	t.Helper()

	destination, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mod)
	assert.NoError(t, err)

	_, err = destination.Write(content)
	assert.NoError(t, err)

	defer destination.Close()

	assert.NoError(t, os.Chtimes(path, modTime, modTime))
	assert.NoError(t, os.Chmod(path, mod))
}
