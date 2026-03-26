package scrum

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sleepyeldrazi/scru-llm/internal/types"
	"github.com/sleepyeldrazi/scru-llm/internal/workers"
)

// WorkerCoordinator coordinates all the worker agents
type WorkerCoordinator struct {
	workspaceDir string
	ProductOwner *workers.ProductOwner
	SpecEngineer *workers.SpecEngineer
	CodeEngineer *workers.CodeEngineer
	Reviewer     *workers.Reviewer
	Verifier     *Verifier
}

// NewWorkerCoordinator creates a new worker coordinator
func NewWorkerCoordinator(workspaceDir string, productOwner, specEngineer, codeEngineer, reviewer interface{}, verifier *Verifier) *WorkerCoordinator {
	return &WorkerCoordinator{
		workspaceDir: workspaceDir,
		ProductOwner: productOwner.(*workers.ProductOwner),
		SpecEngineer: specEngineer.(*workers.SpecEngineer),
		CodeEngineer: codeEngineer.(*workers.CodeEngineer),
		Reviewer:     reviewer.(*workers.Reviewer),
		Verifier:     verifier,
	}
}

// SaveSpec saves a specification to the workspace
func (wc *WorkerCoordinator) SaveSpec(taskID string, spec *types.Specification) error {
	path := filepath.Join(wc.workspaceDir, taskID, "spec", fmt.Sprintf("spec_v%03d.json", spec.Version))
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal spec: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadSpec loads the latest specification from the workspace
func (wc *WorkerCoordinator) LoadSpec(taskID string) (*types.Specification, error) {
	specDir := filepath.Join(wc.workspaceDir, taskID, "spec")

	// Find the latest spec version
	entries, err := os.ReadDir(specDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec directory: %w", err)
	}

	var latestPath string
	var latestVersion int

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(entry.Name(), "spec_v%d.json", &version); err == nil {
			if version > latestVersion {
				latestVersion = version
				latestPath = filepath.Join(specDir, entry.Name())
			}
		}
	}

	if latestPath == "" {
		return nil, fmt.Errorf("no spec found for task %s", taskID)
	}

	data, err := os.ReadFile(latestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	var spec types.Specification
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}

	return &spec, nil
}

// SaveTests saves test files to the workspace
func (wc *WorkerCoordinator) SaveTests(taskID string, tests map[string]string) error {
	testDir := filepath.Join(wc.workspaceDir, taskID, "tests")

	for filename, content := range tests {
		path := filepath.Join(testDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write test file %s: %w", filename, err)
		}
	}

	return nil
}

// LoadTests loads test files from the workspace
func (wc *WorkerCoordinator) LoadTests(taskID string) (map[string]string, error) {
	testDir := filepath.Join(wc.workspaceDir, taskID, "tests")

	entries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tests directory: %w", err)
	}

	tests := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(testDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		tests[entry.Name()] = string(data)
	}

	return tests, nil
}

// SaveCode saves code files to the workspace
func (wc *WorkerCoordinator) SaveCode(taskID string, code map[string]string) error {
	srcDir := filepath.Join(wc.workspaceDir, taskID, "src")

	for filename, content := range code {
		path := filepath.Join(srcDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write code file %s: %w", filename, err)
		}
	}

	return nil
}

// LoadCode loads code files from the workspace
func (wc *WorkerCoordinator) LoadCode(taskID string) (map[string]string, error) {
	srcDir := filepath.Join(wc.workspaceDir, taskID, "src")

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read src directory: %w", err)
	}

	code := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(srcDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		code[entry.Name()] = string(data)
	}

	return code, nil
}
