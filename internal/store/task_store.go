// Package store provides persistent storage for tasks and sprints.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sleepyeldrazi/scru-llm/internal/errors"
	"github.com/sleepyeldrazi/scru-llm/internal/types"
)

// TaskStore defines the interface for task storage
type TaskStore interface {
	// Get retrieves a task by ID
	Get(id string) (*types.Task, error)

	// Save persists a task
	Save(task *types.Task) error

	// Delete removes a task
	Delete(id string) error

	// List returns all tasks, optionally filtered by status
	List(status ...types.TaskStatus) ([]*types.Task, error)

	// ListPending returns tasks awaiting processing
	ListPending() ([]*types.Task, error)

	// ListActive returns tasks currently in a sprint
	ListActive() ([]*types.Task, error)
}

// FileTaskStore implements TaskStore using the filesystem
type FileTaskStore struct {
	mu       sync.RWMutex
	basePath string
}

// NewFileTaskStore creates a new file-based task store
func NewFileTaskStore(basePath string) (*FileTaskStore, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create task store directory: %w", err)
	}

	return &FileTaskStore{
		basePath: basePath,
	}, nil
}

// taskPath returns the filesystem path for a task file
func (s *FileTaskStore) taskPath(id string) string {
	return filepath.Join(s.basePath, fmt.Sprintf("%s.json", id))
}

// Get retrieves a task by ID
func (s *FileTaskStore) Get(id string) (*types.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.taskPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", errors.ErrTaskNotFound, id)
		}
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	var task types.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// Save persists a task
func (s *FileTaskStore) Save(task *types.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update timestamp
	task.UpdatedAt = time.Now().UTC()

	// Marshal to JSON
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Write atomically using temp file
	path := s.taskPath(task.ID)
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename task file: %w", err)
	}

	return nil
}

// Delete removes a task
func (s *FileTaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.taskPath(id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", errors.ErrTaskNotFound, id)
		}
		return fmt.Errorf("failed to delete task file: %w", err)
	}

	return nil
}

// List returns all tasks, optionally filtered by status
func (s *FileTaskStore) List(status ...types.TaskStatus) ([]*types.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Read all task files
	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task directory: %w", err)
	}

	var tasks []*types.Task
	statusFilter := make(map[types.TaskStatus]bool)
	for _, s := range status {
		statusFilter[s] = true
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(s.basePath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue // Skip unreadable files
		}

		var task types.Task
		if err := json.Unmarshal(data, &task); err != nil {
			continue // Skip invalid files
		}

		// Apply status filter if provided
		if len(statusFilter) > 0 && !statusFilter[task.Status] {
			continue
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// ListPending returns tasks awaiting processing
func (s *FileTaskStore) ListPending() ([]*types.Task, error) {
	return s.List(types.TaskStatusPending)
}

// ListActive returns tasks currently in a sprint
func (s *FileTaskStore) ListActive() ([]*types.Task, error) {
	return s.List(
		types.TaskStatusSpecSprint,
		types.TaskStatusTestSprint,
		types.TaskStatusImplementationSprint,
		types.TaskStatusReview,
	)
}
