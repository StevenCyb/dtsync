package fs

// OperationI is the interface for the operations.
type OperationI interface {
	// Delete a file or directory (recursively).
	Delete(path string) error
	// Copy a file or directory (recursively).
	Copy(src, dst string) error
	// Exists checks if a file or directory exists and what type it is.
	Exists(path string) bool
	// Equal checks if two files are equal.
	/* Returns:
	 * - Directory==Directory -> true.
	 * - Directory==File || File==Directory -> false.
	 * - File==File -> true if equal, false if not equal.
	 *   Check includes: creation time, modify time, size,
	 */
	Equal(src, dst string) bool
}

// ShadowScanI is the interface for the FS scanning library.
type ShadowScanI interface {
	// Start a scanning process for a given root path calling the callbacks with
	// the src path and des path. It returns fs.ErrScannerAtEnd when reaching the end.
	Start(rootPath, dstRootPath string, fileCallback, dirCallback ScannerCallback) chan error
	// Stop all scanning processes.
	// Once stopped, the instance cant be reused.
	Stop()
}
