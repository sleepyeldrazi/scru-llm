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

// SprintStore defines the interface for sprint storage
type SprintStore interface {
	// Get retrieves a sprint by ID
	Get(id string) (*types.Sprint, error)

	// Save persists a sprint
	Save(sprint *types.Sprint) error

	// Delete removes a sprint
	Delete(id string) error

	// GetByTask retrieves the current sprint for a task
	GetByTask(taskID string) (*types.Sprint, error)

	// List returns all sprints for a task
	ListByTask(taskID string) ([]*types.Sprint, error)
}

// FileSprintStore implements SprintStore using the filesystem
type FileSprintStore struct {
	mu       sync.RWMutex
	basePath string
}

// NewFileSprintStore creates a new file-based sprint store
func NewFileSprintStore(basePath string) (*FileSprintStore, error) {
	sprintPath := filepath.Join(basePath, "sprints")
	if err := os.MkdirAll(sprintPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sprint store directory: %w", err)
	}

	return &FileSprintStore{
		basePath: sprintPath,
	}, nil
}

// sprintPath returns the filesystem path for a sprint file
func (s *FileSprintStore) sprintPath(id string) string {
	return filepath.Join(s.basePath, fmt.Sprintf("%s.json", id))
}

// Get retrieves a sprint by ID
func (s *FileSprintStore) Get(id string) (*types.Sprint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.sprintPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", errors.ErrSprintNotFound, id)
		}
		return nil, fmt.Errorf("failed to read sprint file: %w", err)
	}

	var sprint types.Sprint
	if err := json.Unmarshal(data, &sprint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sprint: %w", err)
	}

	return &sprint, nil
}

// Save persists a sprint
func (s *FileSprintStore) Save(sprint *types.Sprint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update timestamp
	sprint.UpdatedAt = time.Now().UTC()

	// Marshal to JSON
	data, err := json.MarshalIndent(sprint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sprint: %w", err)
	}

	// Write atomically
	path := s.sprintPath(sprint.ID)
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write sprint file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename sprint file: %w", err)
	}

	return nil
}

// Delete removes a sprint
func (s *FileSprintStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.sprintPath(id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", errors.ErrSprintNotFound, id)
		}
		return fmt.Errorf("failed to delete sprint file: %w", err)
	}

	return nil
}

// GetByTask retrieves the most recent active sprint for a task
func (s *FileSprintStore) GetByTask(taskID string) (*types.Sprint, error) {
	sprints, err := s.ListByTask(taskID)
	if err != nil {
		return nil, err
	}

	if len(sprints) == 0 {
		return nil, fmt.Errorf("%w: no sprints found for task %s", errors.ErrSprintNotFound, taskID)
	}

	// Return the most recent sprint
	return sprints[0], nil
}

// ListByTask returns all sprints for a task, most recent first
func (s *FileSprintStore) ListByTask(taskID string) ([]*types.Sprint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sprint directory: %w", err)
	}

	var sprints []*types.Sprint

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(s.basePath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var sprint types.Sprint
		if err := json.Unmarshal(data, &sprint); err != nil {
			continue
		}

		if sprint.TaskID == taskID {
			sprints = append(sprints, &sprint)
		}
	}

	// Sort by started_at descending
	for i := 0; i < len(sprints)-1; i++ {
		for j := i + 1; j < len(sprints); j++ {
			if sprints[j].StartedAt.After(sprints[i].StartedAt) {
				sprints[i], sprints[j] = sprints[j], sprints[i]
			}
		}
	}

	return sprints, nil
}
