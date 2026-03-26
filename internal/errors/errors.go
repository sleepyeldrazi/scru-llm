// Package errors defines sentinel errors used throughout Scru-LLM.
package errors

import "errors"

// Sentinel errors for common failure modes
var (
	// Task errors
	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskExists     = errors.New("task already exists")
	ErrInvalidTask    = errors.New("invalid task")
	ErrTaskNotPending = errors.New("task is not in pending state")

	// Sprint errors
	ErrSprintNotFound  = errors.New("sprint not found")
	ErrSprintFailed    = errors.New("sprint failed")
	ErrSprintCompleted = errors.New("sprint already completed")
	ErrMaxIterations   = errors.New("max iterations reached")
	ErrInvalidPhase    = errors.New("invalid sprint phase")

	// Container errors
	ErrContainerFailed  = errors.New("container execution failed")
	ErrContainerTimeout = errors.New("container execution timed out")

	// Config errors
	ErrInvalidConfig  = errors.New("invalid configuration")
	ErrConfigNotFound = errors.New("configuration file not found")

	// LLM errors
	ErrLLMRequestFailed = errors.New("LLM request failed")
	ErrLLMTimeout       = errors.New("LLM request timed out")
	ErrProviderNotFound = errors.New("LLM provider not found")

	// Workspace errors
	ErrWorkspaceExists   = errors.New("workspace already exists")
	ErrWorkspaceNotFound = errors.New("workspace not found")

	// Verification errors
	ErrVerificationFailed = errors.New("verification failed")
	ErrTestsFailed        = errors.New("tests failed")
)
