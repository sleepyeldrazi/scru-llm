// Package types defines the core domain types for Scru-LLM.
// These types represent the canonical contracts as defined in docs/CONTRACTS.md.
package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	TaskStatusPending              TaskStatus = "pending"
	TaskStatusSpecSprint           TaskStatus = "spec_sprint"
	TaskStatusTestSprint           TaskStatus = "test_sprint"
	TaskStatusImplementationSprint TaskStatus = "implementation_sprint"
	TaskStatusReview               TaskStatus = "review"
	TaskStatusCompleted            TaskStatus = "completed"
	TaskStatusFailed               TaskStatus = "failed"
	TaskStatusCancelled            TaskStatus = "cancelled"
)

// SprintPhase represents the phase of a sprint
type SprintPhase string

const (
	SprintPhaseIntake SprintPhase = "intake"
	SprintPhaseSpec   SprintPhase = "spec"
	SprintPhaseTest   SprintPhase = "test"
	SprintPhaseCode   SprintPhase = "code"
	SprintPhaseReview SprintPhase = "review"
	SprintPhaseDone   SprintPhase = "done"
	SprintPhaseFailed SprintPhase = "failed"
)

// Task represents a user-submitted task to be processed
type Task struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Status          TaskStatus `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	WorkspaceDir    string     `json:"workspace_dir"`
	CurrentSprintID string     `json:"current_sprint_id,omitempty"`
	FinalOutcome    string     `json:"final_outcome,omitempty"`
}

// NewTask creates a new task with the given description
func NewTask(id, title, description string) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      TaskStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Sprint represents an execution cycle for a task
type Sprint struct {
	ID            string         `json:"id"`
	TaskID        string         `json:"task_id"`
	Phase         SprintPhase    `json:"phase"`
	Iteration     int            `json:"iteration"`
	MaxIterations map[string]int `json:"max_iterations"`
	StartedAt     time.Time      `json:"started_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CompletedAt   *time.Time     `json:"completed_at,omitempty"`
	FailureReason string         `json:"failure_reason,omitempty"`
}

// DefaultMaxIterations returns the default iteration limits
func DefaultMaxIterations() map[string]int {
	return map[string]int{
		"spec":   3,
		"test":   3,
		"code":   5,
		"review": 2,
	}
}

// NewSprint creates a new sprint for a task
func NewSprint(id, taskID string) *Sprint {
	now := time.Now().UTC()
	return &Sprint{
		ID:            id,
		TaskID:        taskID,
		Phase:         SprintPhaseIntake,
		Iteration:     0,
		MaxIterations: DefaultMaxIterations(),
		StartedAt:     now,
		UpdatedAt:     now,
	}
}

// Specification represents a detailed task specification
type Specification struct {
	TaskID             string                 `json:"task_id"`
	Version            int                    `json:"version"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	Requirements       []string               `json:"requirements"`
	AcceptanceCriteria []string               `json:"acceptance_criteria"`
	Constraints        []string               `json:"constraints"`
	Dependencies       []string               `json:"dependencies"`
	EdgeCases          []EdgeCase             `json:"edge_cases"`
	Context            map[string]interface{} `json:"context"`
	Changes            []string               `json:"changes"`
	TestabilityNotes   string                 `json:"testability_notes"`
	Confidence         float64                `json:"confidence"`
}

// EdgeCase represents a specific edge case scenario
type EdgeCase struct {
	Scenario         string `json:"scenario"`
	ExpectedBehavior string `json:"expected_behavior"`
}

// ReviewDecision represents the outcome of a review
type ReviewDecision struct {
	Approved              bool     `json:"approved"`
	Phase                 string   `json:"phase"`
	Summary               string   `json:"summary"`
	Issues                []string `json:"issues"`
	RecommendedNextAction string   `json:"recommended_next_action"`
}

// Event represents a system event for audit and tracking
type Event struct {
	Timestamp time.Time       `json:"timestamp"`
	TaskID    string          `json:"task_id,omitempty"`
	SprintID  string          `json:"sprint_id,omitempty"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Message   string          `json:"message"`
	Meta      json.RawMessage `json:"meta,omitempty"`
}

// NewEvent creates a new event
func NewEvent(eventType, source, message string) *Event {
	return &Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		Source:    source,
		Message:   message,
	}
}

// WithTask associates a task with the event
func (e *Event) WithTask(taskID string) *Event {
	e.TaskID = taskID
	return e
}

// WithSprint associates a sprint with the event
func (e *Event) WithSprint(sprintID string) *Event {
	e.SprintID = sprintID
	return e
}

// WithMeta adds metadata to the event
func (e *Event) WithMeta(meta interface{}) (*Event, error) {
	data, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal meta: %w", err)
	}
	e.Meta = data
	return e, nil
}

// ToJSON serializes the event to JSON
func (e *Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Event types as defined in CONTRACTS.md
const (
	EventTypeTaskCreated           = "task.created"
	EventTypeTaskUpdated           = "task.updated"
	EventTypeSprintStarted         = "sprint.started"
	EventTypePhaseChanged          = "phase.changed"
	EventTypeIterationStarted      = "iteration.started"
	EventTypeIterationCompleted    = "iteration.completed"
	EventTypeArtifactWritten       = "artifact.written"
	EventTypeVerificationStarted   = "verification.started"
	EventTypeVerificationCompleted = "verification.completed"
	EventTypeReviewCompleted       = "review.completed"
	EventTypeSprintCompleted       = "sprint.completed"
	EventTypeSprintFailed          = "sprint.failed"
)
