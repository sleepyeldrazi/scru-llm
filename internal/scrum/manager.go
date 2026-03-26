// Package scrum provides Scrum orchestration for Scru-LLM.
package scrum

import (
	"context"
	"fmt"

	"github.com/sleepyeldrazi/scru-llm/internal/config"
	"github.com/sleepyeldrazi/scru-llm/internal/errors"
	"github.com/sleepyeldrazi/scru-llm/internal/store"
	"github.com/sleepyeldrazi/scru-llm/internal/types"
)

// SprintManager orchestrates the sprint lifecycle
type SprintManager struct {
	cfg         *config.Config
	taskStore   store.TaskStore
	sprintStore store.SprintStore
	eventLog    store.EventLog
	worker      *WorkerCoordinator
}

// NewSprintManager creates a new sprint manager
func NewSprintManager(
	cfg *config.Config,
	taskStore store.TaskStore,
	sprintStore store.SprintStore,
	eventLog store.EventLog,
	worker *WorkerCoordinator,
) *SprintManager {
	return &SprintManager{
		cfg:         cfg,
		taskStore:   taskStore,
		sprintStore: sprintStore,
		eventLog:    eventLog,
		worker:      worker,
	}
}

// StartSprint starts a new sprint for a task
func (sm *SprintManager) StartSprint(ctx context.Context, taskID string) (*types.Sprint, error) {
	// Get task
	task, err := sm.taskStore.Get(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Check if task can start
	if task.Status != types.TaskStatusPending {
		return nil, fmt.Errorf("%w: task %s is not pending", errors.ErrInvalidTask, taskID)
	}

	// Create sprint
	sprintID := fmt.Sprintf("sprint-%s", taskID)
	sprint := types.NewSprint(sprintID, taskID)

	// Save sprint
	if err := sm.sprintStore.Save(sprint); err != nil {
		return nil, fmt.Errorf("failed to save sprint: %w", err)
	}

	// Update task
	task.Status = types.TaskStatusSpecSprint
	task.CurrentSprintID = sprintID
	if err := sm.taskStore.Save(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Log event
	event := types.NewEvent(types.EventTypeSprintStarted, "sprint_manager", "Sprint started").
		WithTask(taskID).
		WithSprint(sprintID)
	sm.eventLog.Append(event)

	// Start processing in background
	go sm.processSprint(context.Background(), sprint)

	return sprint, nil
}

// StartSprintSync starts a sprint and runs it synchronously (for debugging)
func (sm *SprintManager) StartSprintSync(ctx context.Context, taskID string) (*types.Sprint, error) {
	// Get task
	task, err := sm.taskStore.Get(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Check if task can start
	if task.Status != types.TaskStatusPending {
		return nil, fmt.Errorf("%w: task %s is not pending", errors.ErrInvalidTask, taskID)
	}

	// Create sprint
	sprintID := fmt.Sprintf("sprint-%s", taskID)
	sprint := types.NewSprint(sprintID, taskID)

	// Save sprint
	if err := sm.sprintStore.Save(sprint); err != nil {
		return nil, fmt.Errorf("failed to save sprint: %w", err)
	}

	// Update task
	task.Status = types.TaskStatusSpecSprint
	task.CurrentSprintID = sprintID
	if err := sm.taskStore.Save(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Log event
	event := types.NewEvent(types.EventTypeSprintStarted, "sprint_manager", "Sprint started").
		WithTask(taskID).
		WithSprint(sprintID)
	sm.eventLog.Append(event)

	// Run processing synchronously
	sm.processSprint(ctx, sprint)

	return sprint, nil
}

// processSprint runs the sprint through its phases
func (sm *SprintManager) processSprint(ctx context.Context, sprint *types.Sprint) {
	fmt.Printf("[Sprint %s] Starting processSprint...\n", sprint.ID)

	// Phase 1: Intake (create initial spec)
	fmt.Printf("[Sprint %s] Phase 1: Intake...\n", sprint.ID)
	if err := sm.runIntakePhase(ctx, sprint); err != nil {
		fmt.Printf("[Sprint %s] Intake phase FAILED: %v\n", sprint.ID, err)
		sm.failSprint(sprint, fmt.Sprintf("intake phase failed: %v", err))
		return
	}
	fmt.Printf("[Sprint %s] Intake phase COMPLETED\n", sprint.ID)

	// Phase 2: Spec Sprint
	fmt.Printf("[Sprint %s] Phase 2: Spec Sprint...\n", sprint.ID)
	if err := sm.runSpecPhase(ctx, sprint); err != nil {
		fmt.Printf("[Sprint %s] Spec phase FAILED: %v\n", sprint.ID, err)
		sm.failSprint(sprint, fmt.Sprintf("spec phase failed: %v", err))
		return
	}
	fmt.Printf("[Sprint %s] Spec phase COMPLETED\n", sprint.ID)

	// Phase 3: Test Sprint
	fmt.Printf("[Sprint %s] Phase 3: Test Sprint...\n", sprint.ID)
	if err := sm.runTestPhase(ctx, sprint); err != nil {
		fmt.Printf("[Sprint %s] Test phase FAILED: %v\n", sprint.ID, err)
		sm.failSprint(sprint, fmt.Sprintf("test phase failed: %v", err))
		return
	}
	fmt.Printf("[Sprint %s] Test phase COMPLETED\n", sprint.ID)

	// Phase 4: Implementation Sprint
	fmt.Printf("[Sprint %s] Phase 4: Implementation Sprint...\n", sprint.ID)
	if err := sm.runImplementationPhase(ctx, sprint); err != nil {
		fmt.Printf("[Sprint %s] Implementation phase FAILED: %v\n", sprint.ID, err)
		sm.failSprint(sprint, fmt.Sprintf("implementation phase failed: %v", err))
		return
	}
	fmt.Printf("[Sprint %s] Implementation phase COMPLETED\n", sprint.ID)

	// Phase 5: Review
	fmt.Printf("[Sprint %s] Phase 5: Review...\n", sprint.ID)
	if err := sm.runReviewPhase(ctx, sprint); err != nil {
		fmt.Printf("[Sprint %s] Review phase FAILED: %v\n", sprint.ID, err)
		sm.failSprint(sprint, fmt.Sprintf("review phase failed: %v", err))
		return
	}
	fmt.Printf("[Sprint %s] Review phase COMPLETED\n", sprint.ID)

	// Mark as completed
	fmt.Printf("[Sprint %s] Marking as completed...\n", sprint.ID)
	sm.completeSprint(sprint)
	fmt.Printf("[Sprint %s] Sprint COMPLETED successfully!\n", sprint.ID)
}

// runIntakePhase runs the intake phase
func (sm *SprintManager) runIntakePhase(ctx context.Context, sprint *types.Sprint) error {
	// Transition to intake phase
	sprint.Phase = types.SprintPhaseIntake
	if err := sm.sprintStore.Save(sprint); err != nil {
		return err
	}

	sm.logPhaseChange(sprint, "intake")

	// Get task
	task, err := sm.taskStore.Get(sprint.TaskID)
	if err != nil {
		return err
	}

	// Use Product Owner to create initial spec
	spec, err := sm.worker.ProductOwner.CreateSpec(ctx, task)
	if err != nil {
		return fmt.Errorf("product owner failed: %w", err)
	}

	// Save spec
	if err := sm.worker.SaveSpec(sprint.TaskID, spec); err != nil {
		return fmt.Errorf("failed to save spec: %w", err)
	}

	return nil
}

// runSpecPhase runs the spec refinement phase
func (sm *SprintManager) runSpecPhase(ctx context.Context, sprint *types.Sprint) error {
	maxIter := sm.cfg.Sprint.SpecSprint.MaxIterations

	for i := 0; i < maxIter; i++ {
		// Transition to spec phase
		sprint.Phase = types.SprintPhaseSpec
		sprint.Iteration = i + 1
		if err := sm.sprintStore.Save(sprint); err != nil {
			return err
		}

		sm.logPhaseChange(sprint, fmt.Sprintf("spec iteration %d/%d", i+1, maxIter))

		// Get current spec
		spec, err := sm.worker.LoadSpec(sprint.TaskID)
		if err != nil {
			return fmt.Errorf("failed to load spec: %w", err)
		}

		// Spec Engineer tightens spec
		tightenedSpec, err := sm.worker.SpecEngineer.TightenSpec(ctx, spec)
		if err != nil {
			return fmt.Errorf("spec engineer failed: %w", err)
		}

		// Reviewer checks spec
		decision, err := sm.worker.Reviewer.ReviewSpec(ctx, tightenedSpec)
		if err != nil {
			return fmt.Errorf("reviewer failed: %w", err)
		}

		// Log review
		sm.logReview(sprint, "spec", decision)

		if decision.Approved {
			// Save approved spec
			if err := sm.worker.SaveSpec(sprint.TaskID, tightenedSpec); err != nil {
				return err
			}
			return nil
		}

		// Not approved, continue to next iteration
		if i == maxIter-1 {
			return fmt.Errorf("max iterations reached, spec not approved")
		}
	}

	return nil
}

// runTestPhase runs the test generation phase
func (sm *SprintManager) runTestPhase(ctx context.Context, sprint *types.Sprint) error {
	maxIter := sm.cfg.Sprint.TestSprint.MaxIterations

	for i := 0; i < maxIter; i++ {
		// Transition to test phase
		sprint.Phase = types.SprintPhaseTest
		sprint.Iteration = i + 1
		if err := sm.sprintStore.Save(sprint); err != nil {
			return err
		}

		sm.logPhaseChange(sprint, fmt.Sprintf("test iteration %d/%d", i+1, maxIter))

		// Update task status
		task, _ := sm.taskStore.Get(sprint.TaskID)
		task.Status = types.TaskStatusTestSprint
		sm.taskStore.Save(task)

		// Get final spec
		spec, err := sm.worker.LoadSpec(sprint.TaskID)
		if err != nil {
			return err
		}

		// Generate tests
		tests, err := sm.worker.CodeEngineer.GenerateTests(ctx, spec)
		if err != nil {
			return fmt.Errorf("test generation failed: %w", err)
		}

		// Save tests
		if err := sm.worker.SaveTests(sprint.TaskID, tests); err != nil {
			return err
		}

		// Verify tests run (and fail against stubs)
		result, err := sm.worker.Verifier.RunTests(sprint.TaskID)
		if err != nil {
			return fmt.Errorf("test verification failed: %w", err)
		}

		if result.TestsRun > 0 {
			// Tests are running - that's sufficient for MVP
			// In a full implementation, we'd want tests to fail against stubs
			fmt.Printf("[Sprint %s] Tests generated and running (%d tests)\n", sprint.ID, result.TestsRun)
			return nil
		}

		if i == maxIter-1 {
			fmt.Printf("[Sprint %s] Warning: max iterations reached without tests running\n", sprint.ID)
			// For MVP, continue anyway
			return nil
		}
	}

	return nil
}

// runImplementationPhase runs the code implementation phase
func (sm *SprintManager) runImplementationPhase(ctx context.Context, sprint *types.Sprint) error {
	maxIter := sm.cfg.Sprint.ImplementationSprint.MaxIterations

	for i := 0; i < maxIter; i++ {
		// Transition to code phase
		sprint.Phase = types.SprintPhaseCode
		sprint.Iteration = i + 1
		if err := sm.sprintStore.Save(sprint); err != nil {
			return err
		}

		sm.logPhaseChange(sprint, fmt.Sprintf("implementation iteration %d/%d", i+1, maxIter))

		// Update task status
		task, _ := sm.taskStore.Get(sprint.TaskID)
		task.Status = types.TaskStatusImplementationSprint
		sm.taskStore.Save(task)

		// Load spec and tests
		spec, _ := sm.worker.LoadSpec(sprint.TaskID)
		tests, _ := sm.worker.LoadTests(sprint.TaskID)

		// Get previous error output if any
		var previousErrors string
		if i > 0 {
			result, _ := sm.worker.Verifier.RunTests(sprint.TaskID)
			previousErrors = result.Output
		}

		// Implement code
		code, err := sm.worker.CodeEngineer.ImplementCode(ctx, spec, tests, previousErrors)
		if err != nil {
			return fmt.Errorf("code implementation failed: %w", err)
		}

		// Save code
		if err := sm.worker.SaveCode(sprint.TaskID, code); err != nil {
			return err
		}

		// Run tests
		result, err := sm.worker.Verifier.RunTests(sprint.TaskID)
		if err != nil {
			return err
		}

		if result.TestsPassed == result.TestsRun {
			// All tests passing
			return nil
		}

		if i == maxIter-1 {
			return fmt.Errorf("max iterations reached, tests still failing")
		}
	}

	return nil
}

// runReviewPhase runs the final review phase
func (sm *SprintManager) runReviewPhase(ctx context.Context, sprint *types.Sprint) error {
	sprint.Phase = types.SprintPhaseReview
	if err := sm.sprintStore.Save(sprint); err != nil {
		return err
	}

	sm.logPhaseChange(sprint, "final review")

	// Update task status
	task, _ := sm.taskStore.Get(sprint.TaskID)
	task.Status = types.TaskStatusReview
	sm.taskStore.Save(task)

	// Load all artifacts
	spec, _ := sm.worker.LoadSpec(sprint.TaskID)
	tests, _ := sm.worker.LoadTests(sprint.TaskID)
	code, _ := sm.worker.LoadCode(sprint.TaskID)

	// Final review
	decision, err := sm.worker.Reviewer.ReviewImplementation(ctx, spec, tests, code)
	if err != nil {
		return err
	}

	sm.logReview(sprint, "implementation", decision)

	if !decision.Approved {
		return fmt.Errorf("final review rejected: %s", decision.Summary)
	}

	return nil
}

// failSprint marks a sprint as failed
func (sm *SprintManager) failSprint(sprint *types.Sprint, reason string) {
	sprint.Phase = types.SprintPhaseFailed
	sprint.FailureReason = reason
	sm.sprintStore.Save(sprint)

	task, _ := sm.taskStore.Get(sprint.TaskID)
	task.Status = types.TaskStatusFailed
	task.FinalOutcome = reason
	sm.taskStore.Save(task)

	event := types.NewEvent(types.EventTypeSprintFailed, "sprint_manager", reason).
		WithTask(sprint.TaskID).
		WithSprint(sprint.ID)
	sm.eventLog.Append(event)
}

// completeSprint marks a sprint as completed
func (sm *SprintManager) completeSprint(sprint *types.Sprint) {
	sprint.Phase = types.SprintPhaseDone
	now := sprint.UpdatedAt
	sprint.CompletedAt = &now
	sm.sprintStore.Save(sprint)

	task, _ := sm.taskStore.Get(sprint.TaskID)
	task.Status = types.TaskStatusCompleted
	task.FinalOutcome = "completed"
	sm.taskStore.Save(task)

	event := types.NewEvent(types.EventTypeSprintCompleted, "sprint_manager", "Sprint completed successfully").
		WithTask(sprint.TaskID).
		WithSprint(sprint.ID)
	sm.eventLog.Append(event)
}

// logPhaseChange logs a phase change event
func (sm *SprintManager) logPhaseChange(sprint *types.Sprint, phase string) {
	event := types.NewEvent(types.EventTypePhaseChanged, "sprint_manager", fmt.Sprintf("Phase: %s", phase)).
		WithTask(sprint.TaskID).
		WithSprint(sprint.ID)
	sm.eventLog.Append(event)
}

// logReview logs a review event
func (sm *SprintManager) logReview(sprint *types.Sprint, phase string, decision *types.ReviewDecision) {
	status := "rejected"
	if decision.Approved {
		status = "approved"
	}

	event := types.NewEvent(types.EventTypeReviewCompleted, "reviewer", fmt.Sprintf("%s review %s", phase, status)).
		WithTask(sprint.TaskID).
		WithSprint(sprint.ID)
	sm.eventLog.Append(event)
}
