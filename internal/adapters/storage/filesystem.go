package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ravro_dcrpt/internal/ports"
)

// FileSystemService implements the StorageService interface
type FileSystemService struct{}

// NewFileSystemService creates a new filesystem service
func NewFileSystemService() ports.StorageService {
	return &FileSystemService{}
}

// ReadFile reads a file and returns its content
func (s *FileSystemService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes content to a file
func (s *FileSystemService) WriteFile(path string, data []byte) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := s.CreateDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// CreateDir creates a directory if it doesn't exist
func (s *FileSystemService) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// ListFiles lists files in a directory matching a pattern
func (s *FileSystemService) ListFiles(dir string, pattern string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file matches pattern
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if matched {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// FileExists checks if a file exists
func (s *FileSystemService) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ExtractZip extracts a zip file to a destination
func (s *FileSystemService) ExtractZip(zipPath string, destPath string) error {
	// Open zip file
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	// Create destination directory
	if err := s.CreateDir(destPath); err != nil {
		return err
	}

	// Extract files
	for _, f := range r.File {
		err := s.extractZipFile(f, destPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from a zip archive
func (s *FileSystemService) extractZipFile(f *zip.File, destPath string) error {
	// Construct destination path
	destFilePath := filepath.Join(destPath, f.Name)

	// Check for ZipSlip vulnerability
	if !strings.HasPrefix(destFilePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", f.Name)
	}

	// Create directories if needed
	if f.FileInfo().IsDir() {
		return os.MkdirAll(destFilePath, f.Mode())
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
		return err
	}

	// Open source file
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create destination file
	destFile, err := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, rc)
	return err
}

// CopyFile copies a file from source to destination
func (s *FileSystemService) CopyFile(src string, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination directory
	destDir := filepath.Dir(dst)
	if err := s.CreateDir(destDir); err != nil {
		return err
	}

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Sync to ensure data is written
	return destFile.Sync()
}

// OpenFile opens a file for reading
func (s *FileSystemService) OpenFile(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// CreateFile creates a file for writing
func (s *FileSystemService) CreateFile(path string) (io.WriteCloser, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := s.CreateDir(dir); err != nil {
		return nil, err
	}
	return os.Create(path)
}
