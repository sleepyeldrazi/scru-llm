// Package workspace provides workspace and sandbox management.
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Workspace represents a task workspace
type Workspace struct {
	Path string
}

// Create creates a new workspace at the given path
func Create(path string) error {
	// Create directory structure
	dirs := []string{
		path,
		filepath.Join(path, "spec"),
		filepath.Join(path, "tests"),
		filepath.Join(path, "src"),
		filepath.Join(path, "artifacts"),
		filepath.Join(path, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// Open opens an existing workspace
func Open(path string) (*Workspace, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("workspace not found: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("workspace path is not a directory: %s", path)
	}

	return &Workspace{Path: path}, nil
}

// GetPath returns the full path for a subdirectory
func (w *Workspace) GetPath(subdir string) string {
	return filepath.Join(w.Path, subdir)
}

// WriteFile writes a file to the workspace
func (w *Workspace) WriteFile(subdir, filename string, data []byte) error {
	path := filepath.Join(w.Path, subdir, filename)
	return os.WriteFile(path, data, 0644)
}

// ReadFile reads a file from the workspace
func (w *Workspace) ReadFile(subdir, filename string) ([]byte, error) {
	path := filepath.Join(w.Path, subdir, filename)
	return os.ReadFile(path)
}

// Delete removes the workspace
func (w *Workspace) Delete() error {
	return os.RemoveAll(w.Path)
}

// Exists checks if a workspace exists
func Exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
