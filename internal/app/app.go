// Package app provides the main application context and initialization.
package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sleepyeldrazi/scru-llm/internal/config"
	"github.com/sleepyeldrazi/scru-llm/internal/llm"
	"github.com/sleepyeldrazi/scru-llm/internal/prompts"
	"github.com/sleepyeldrazi/scru-llm/internal/scrum"
	"github.com/sleepyeldrazi/scru-llm/internal/store"
	"github.com/sleepyeldrazi/scru-llm/internal/workers"
)

// Context holds all application dependencies
type Context struct {
	Config       *config.Config
	TaskStore    store.TaskStore
	SprintStore  store.SprintStore
	EventLog     store.EventLog
	LLMRouter    *llm.Router
	PromptLoader *prompts.Loader
	SprintMgr    *scrum.SprintManager
	Logger       *slog.Logger
}

// New creates a new application context
func New(mockMode bool) (*Context, error) {
	// Load configuration
	cfg, err := config.LoadOrDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Setup logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: getLogLevel(cfg.System.LogLevel),
	}))

	// Ensure directories exist
	if err := cfg.EnsureDirs(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	// Initialize stores
	taskStore, err := store.NewFileTaskStore(filepath.Join(cfg.System.WorkspaceDir, "tasks"))
	if err != nil {
		return nil, fmt.Errorf("failed to create task store: %w", err)
	}

	sprintStore, err := store.NewFileSprintStore(cfg.System.WorkspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create sprint store: %w", err)
	}

	eventLog, err := store.NewJSONLFileEventLog(cfg.Reporting.AuditLog.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create event log: %w", err)
	}

	// Initialize LLM router
	router := llm.NewRouter(cfg, mockMode)

	// Initialize prompt loader
	promptLoader := prompts.NewLoader("config/prompts")

	// Initialize workers
	productOwner := workers.NewProductOwner(router, promptLoader)
	specEngineer := workers.NewSpecEngineer(router, promptLoader)
	codeEngineer := workers.NewCodeEngineer(router, promptLoader)
	reviewer := workers.NewReviewer(router, promptLoader)
	verifier := scrum.NewVerifier(cfg.System.WorkspaceDir)

	// Initialize worker coordinator
	workerCoord := scrum.NewWorkerCoordinator(
		cfg.System.WorkspaceDir,
		productOwner,
		specEngineer,
		codeEngineer,
		reviewer,
		verifier,
	)

	// Initialize sprint manager
	sprintMgr := scrum.NewSprintManager(
		cfg,
		taskStore,
		sprintStore,
		eventLog,
		workerCoord,
	)

	return &Context{
		Config:       cfg,
		TaskStore:    taskStore,
		SprintStore:  sprintStore,
		EventLog:     eventLog,
		LLMRouter:    router,
		PromptLoader: promptLoader,
		SprintMgr:    sprintMgr,
		Logger:       logger,
	}, nil
}

// Close cleans up resources
func (ctx *Context) Close() error {
	if ctx.EventLog != nil {
		return ctx.EventLog.Close()
	}
	return nil
}

// getLogLevel converts string to slog.Level
func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
