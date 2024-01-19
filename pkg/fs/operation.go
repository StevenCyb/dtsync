package fs

import (
	"io"
	"os"
	"time"
)

// Operation provides FS operations.
type Operation struct{}

// NewOperation creates a new operation.
func NewOperation() OperationI {
	return &Operation{}
}

// Delete a file or directory (recursively).
func (o *Operation) Delete(path string) error {
	if exists := o.Exists(path); exists {
		return os.RemoveAll(path)
	}

	return os.ErrNotExist
}

// Copy a file or directory (recursively).
func (o *Operation) Copy(src, dst string) error {
	srcState, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcState.IsDir() {
		return os.Mkdir(dst, srcState.Mode())
	} else if !srcState.Mode().IsRegular() {
		return nil
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcState.Mode())
	if err != nil {
		return err
	}

	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	if err := os.Chtimes(dst, time.Now(), srcState.ModTime()); err != nil {
		return err
	}

	return err
}

// Exists checks if a file or directory exists and what type it is.
func (o *Operation) Exists(path string) bool {
	if state, err := os.Stat(path); err == nil {
		if state.IsDir() {
			return true
		}

		return true
	}

	return false
}

// Equal checks if two files are equal.
/* Returns:
* - one not exists -> false.
* - Directory==File || File==Directory -> false.
* - Directory==Directory -> true.
 * - File==File -> true if equal, false if not equal.
*/
func (o *Operation) Equal(src, dst string) bool {
	srcState, err := os.Stat(src)
	if err != nil {
		return false
	}

	dstStatus, err := os.Stat(dst)
	if err != nil {
		return false
	}

	if srcState.IsDir() {
		return srcState.IsDir() == dstStatus.IsDir() &&
			srcState.Mode().Perm().String() == dstStatus.Mode().Perm().String()
	}

	return srcState.Size() == dstStatus.Size() &&
		srcState.ModTime() == dstStatus.ModTime() &&
		srcState.Mode().Perm().String() == dstStatus.Mode().Perm().String()
}
