package ports

import "io"

// StorageService defines the interface for file system operations
type StorageService interface {
	// ReadFile reads a file and returns its content
	ReadFile(path string) ([]byte, error)

	// WriteFile writes content to a file
	WriteFile(path string, data []byte) error

	// CreateDir creates a directory if it doesn't exist
	CreateDir(path string) error

	// ListFiles lists files in a directory matching a pattern
	ListFiles(dir string, pattern string) ([]string, error)

	// FileExists checks if a file exists
	FileExists(path string) bool

	// ExtractZip extracts a zip file to a destination
	ExtractZip(zipPath string, destPath string) error

	// CopyFile copies a file from source to destination
	CopyFile(src string, dst string) error

	// OpenFile opens a file for reading
	OpenFile(path string) (io.ReadCloser, error)

	// CreateFile creates a file for writing
	CreateFile(path string) (io.WriteCloser, error)
}
