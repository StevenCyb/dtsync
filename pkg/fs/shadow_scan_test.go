package fs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShadowScan(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.RemoveAll("test_shadow_scan")
	})
	assert.NoError(t, os.Mkdir("test_shadow_scan", 0o755))
	assert.NoError(t, os.Mkdir("test_shadow_scan/a", 0o755))
	assert.NoError(t, os.Mkdir("test_shadow_scan/b", 0o755))
	assert.NoError(t, os.Mkdir("test_shadow_scan/a/b", 0o755))
	createTestFile(t, "test_shadow_scan/a/hello.txt", 0x755, time.Now(), []byte{})
	createTestFile(t, "test_shadow_scan/b/hello.txt", 0x755, time.Now(), []byte{})
	createTestFile(t, "test_shadow_scan/b/world.txt", 0x755, time.Now(), []byte{})
	createTestFile(t, "test_shadow_scan/a/b/some.txt", 0x755, time.Now(), []byte{})

	t.Run("Normal", func(t *testing.T) {
		t.Parallel()

		scanner := NewShadowScan()
		foundedFiles := map[string]string{}
		foundedDirectories := map[string]string{}

		errChan := scanner.Start("test_shadow_scan", "dest",
			func(srcPath, dstPath string) error {
				foundedFiles[srcPath] = dstPath

				return nil
			},
			func(srcPath, dstPath string) error {
				foundedDirectories[srcPath] = dstPath

				return nil
			},
		)

		err := <-errChan
		assert.Error(t, err)
		assert.Equal(t, ErrScannerAtEnd, err)
		assert.Equal(t, map[string]string{
			"test_shadow_scan/a/b/some.txt": "dest/a/b/some.txt", "test_shadow_scan/a/hello.txt": "dest/a/hello.txt",
			"test_shadow_scan/b/hello.txt": "dest/b/hello.txt", "test_shadow_scan/b/world.txt": "dest/b/world.txt",
		}, foundedFiles)
		assert.Equal(t, map[string]string{
			"test_shadow_scan": "dest", "test_shadow_scan/a": "dest/a",
			"test_shadow_scan/a/b": "dest/a/b", "test_shadow_scan/b": "dest/b",
		}, foundedDirectories)
	})

	t.Run("Stopped", func(t *testing.T) {
		t.Parallel()

		scanner := NewShadowScan()
		errChan := scanner.Start("test_shadow_scan", "dest",
			func(srcPath, dstPath string) error {
				time.Sleep(time.Second)

				return nil
			},
			func(srcPath, dstPath string) error {
				time.Sleep(time.Second)

				return nil
			},
		)

		time.Sleep(time.Second)
		scanner.Stop()

		err := <-errChan
		assert.Error(t, err)
		assert.Equal(t, ErrShadowScanStopped, err)
	})
}
