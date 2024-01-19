package fs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

var (
	// ErrScannerAtEnd is returned when the scanner is at the end.
	ErrScannerAtEnd = errors.New("scanner is at the end")
	// ErrShadowScanStopped is returned when the scanner is stopped.
	ErrShadowScanStopped = errors.New("scanner already stopped")
)

// ScannerCallback is the callback function for the scanner.
// It is called when a file or directory is found.
// When the callback returns an error, the scanner stops.
type ScannerCallback = func(srcPath, dstPath string) error

// ShadowScan provides FS scanning functionality.
// The callback is called with the src and dst path.
type ShadowScan struct {
	stop bool
}

// NewShadowScan creates a new scanner.
func NewShadowScan() ShadowScanI {
	return &ShadowScan{}
}

// Start starts the scanner.
func (s *ShadowScan) Start(srcRootPath, dstRootPath string, fileCallback, dirCallback ScannerCallback) chan error {
	errChan := make(chan error)

	if s.stop {
		errChan <- ErrShadowScanStopped

		return errChan
	}

	go func() {
		err := fs.WalkDir(os.DirFS(srcRootPath), ".", func(srcPath string, dirEntry fs.DirEntry, err error) error {
			switch {
			case errors.Is(err, unix.ENOENT):
				return nil
			case err != nil:
				return err
			case s.stop:
				return ErrShadowScanStopped
			}

			dstPath := filepath.Join(
				dstRootPath,
				srcPath,
			)

			if srcPath != "." {
				srcPath = filepath.Join(
					srcRootPath,
					srcPath,
				)
			} else {
				srcPath = srcRootPath
			}

			if dirEntry.IsDir() {
				return dirCallback(srcPath, dstPath)
			}

			return fileCallback(srcPath, dstPath)
		})

		if err != nil {
			errChan <- err
		} else {
			errChan <- ErrScannerAtEnd
		}
	}()

	return errChan
}

// Stop stops the scanner.
func (s *ShadowScan) Stop() {
	s.stop = true
}
