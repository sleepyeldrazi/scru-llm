// Package cli provides CLI command implementations for Scru-LLM.
package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sleepyeldrazi/scru-llm/internal/app"
	"github.com/sleepyeldrazi/scru-llm/internal/config"
	"github.com/sleepyeldrazi/scru-llm/internal/store"
	"github.com/sleepyeldrazi/scru-llm/internal/types"
	"github.com/sleepyeldrazi/scru-llm/internal/workspace"
	"github.com/spf13/cobra"
)

// getConfig loads the configuration
func getConfig() (*config.Config, error) {
	cfg, err := config.LoadOrDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Ensure directories exist
	if err := cfg.EnsureDirs(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return cfg, nil
}

// getStores returns the task and sprint stores
func getStores(cfg *config.Config) (store.TaskStore, store.SprintStore, store.EventLog, error) {
	taskStore, err := store.NewFileTaskStore(filepath.Join(cfg.System.WorkspaceDir, "tasks"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create task store: %w", err)
	}

	sprintStore, err := store.NewFileSprintStore(cfg.System.WorkspaceDir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create sprint store: %w", err)
	}

	eventLog, err := store.NewJSONLFileEventLog(cfg.Reporting.AuditLog.Path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create event log: %w", err)
	}

	return taskStore, sprintStore, eventLog, nil
}

// generateTaskID generates a unique task ID
func generateTaskID() string {
	return fmt.Sprintf("task-%d", os.Getpid())
}

// NewTaskCommand creates the task command
func NewTaskCommand() *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "task [description]",
		Short: "Submit a new task",
		Long:  `Submit a new task to be processed by Scru-LLM.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := args[0]

			// Load config
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			// Get stores
			taskStore, _, eventLog, err := getStores(cfg)
			if err != nil {
				return err
			}
			defer eventLog.Close()

			// Generate task ID
			taskID := generateTaskID()

			// Create workspace (workspace_dir already points to workspaces directory)
			wsPath := filepath.Join(cfg.System.WorkspaceDir, taskID)
			if err := workspace.Create(wsPath); err != nil {
				return fmt.Errorf("failed to create workspace: %w", err)
			}

			// Create task
			taskTitle := title
			if taskTitle == "" {
				taskTitle = description[:min(len(description), 50)]
				if len(description) > 50 {
					taskTitle += "..."
				}
			}

			task := types.NewTask(taskID, taskTitle, description)
			task.WorkspaceDir = wsPath

			// Save task
			if err := taskStore.Save(task); err != nil {
				return fmt.Errorf("failed to save task: %w", err)
			}

			// Log event
			event := types.NewEvent(types.EventTypeTaskCreated, "cli", "Task created").
				WithTask(taskID)
			if err := eventLog.Append(event); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to log event: %v\n", err)
			}

			fmt.Printf("Task created: %s\n", taskID)
			fmt.Printf("Title: %s\n", task.Title)
			fmt.Printf("Status: %s\n", task.Status)
			fmt.Printf("Workspace: %s\n", wsPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Task title (optional)")

	return cmd
}

// NewStatusCommand creates the status command
func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [task-id]",
		Short: "Show task or system status",
		Long:  `Show the current status of a task or the overall system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			// Get stores
			taskStore, sprintStore, eventLog, err := getStores(cfg)
			if err != nil {
				return err
			}
			defer eventLog.Close()

			if len(args) > 0 {
				// Show specific task status
				taskID := args[0]
				task, err := taskStore.Get(taskID)
				if err != nil {
					return fmt.Errorf("failed to get task: %w", err)
				}

				fmt.Printf("Task: %s\n", task.ID)
				fmt.Printf("Title: %s\n", task.Title)
				fmt.Printf("Status: %s\n", task.Status)
				fmt.Printf("Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))

				if task.CurrentSprintID != "" {
					sprint, err := sprintStore.Get(task.CurrentSprintID)
					if err == nil {
						fmt.Printf("\nCurrent Sprint:\n")
						fmt.Printf("  ID: %s\n", sprint.ID)
						fmt.Printf("  Phase: %s\n", sprint.Phase)
						fmt.Printf("  Iteration: %d/%d\n", sprint.Iteration, sprint.MaxIterations[string(sprint.Phase)])
						fmt.Printf("  Started: %s\n", sprint.StartedAt.Format("2006-01-02 15:04:05"))
					}
				}
			} else {
				// Show system status
				pending, _ := taskStore.ListPending()
				active, _ := taskStore.ListActive()
				completed, _ := taskStore.List(types.TaskStatusCompleted)
				failed, _ := taskStore.List(types.TaskStatusFailed)

				fmt.Println("System Status")
				fmt.Println("=============")
				fmt.Printf("Pending:   %d\n", len(pending))
				fmt.Printf("Active:    %d\n", len(active))
				fmt.Printf("Completed: %d\n", len(completed))
				fmt.Printf("Failed:    %d\n", len(failed))

				if len(active) > 0 {
					fmt.Println("\nActive Tasks:")
					for _, task := range active {
						fmt.Printf("  - %s [%s] %s\n", task.ID, task.Status, task.Title)
					}
				}
			}

			return nil
		},
	}

	return cmd
}

// NewLogsCommand creates the logs command
func NewLogsCommand() *cobra.Command {
	var taskID string
	var limit int

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View system logs",
		Long:  `View audit logs and event history.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			// Get event log
			_, _, eventLog, err := getStores(cfg)
			if err != nil {
				return err
			}
			defer eventLog.Close()

			// Read events
			opts := store.ReadOptions{
				TaskID: taskID,
				Limit:  limit,
			}

			events, err := eventLog.Read(opts)
			if err != nil {
				return fmt.Errorf("failed to read events: %w", err)
			}

			if len(events) == 0 {
				fmt.Println("No events found.")
				return nil
			}

			fmt.Printf("Events (%d total):\n", len(events))
			fmt.Println("==================")

			for _, event := range events {
				fmt.Printf("[%s] %s: %s\n",
					event.Timestamp.Format("2006-01-02 15:04:05"),
					event.Type,
					event.Message)
				if event.TaskID != "" {
					fmt.Printf("  Task: %s\n", event.TaskID)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&taskID, "task", "t", "", "Filter by task ID")
	cmd.Flags().IntVarP(&limit, "limit", "n", 50, "Maximum number of events to show")

	return cmd
}

// NewConfigCommand creates the config command
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `View and validate configuration.`,
	}

	// Show config
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			fmt.Printf("Configuration:\n")
			fmt.Printf("  Workspace Dir: %s\n", cfg.System.WorkspaceDir)
			fmt.Printf("  Log Level: %s\n", cfg.System.LogLevel)
			fmt.Printf("  Max Concurrent Tasks: %d\n", cfg.System.MaxConcurrentTasks)
			fmt.Printf("\n  LLM Providers: %d\n", len(cfg.LLM.Providers))
			for _, p := range cfg.LLM.Providers {
				fmt.Printf("    - %s (%s)\n", p.Name, p.BaseURL)
			}
			fmt.Printf("\n  Roles: %d\n", len(cfg.LLM.Roles))
			for role, cfg := range cfg.LLM.Roles {
				fmt.Printf("    - %s: %s/%s\n", role, cfg.Provider, cfg.Model)
			}

			return nil
		},
	}

	// Validate config
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("configuration invalid: %w", err)
			}

			fmt.Println("Configuration is valid.")
			return nil
		},
	}

	cmd.AddCommand(showCmd)
	cmd.AddCommand(validateCmd)

	return cmd
}

// NewRunCommand creates the run command
func NewRunCommand() *cobra.Command {
	var mockMode bool
	var syncMode bool

	cmd := &cobra.Command{
		Use:   "run <task-id>",
		Short: "Run a task through the sprint workflow",
		Long:  `Start processing a task through the complete sprint workflow (intake, spec, test, implementation, review).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("🚀 Starting sprint for task: %s\n", taskID)

			// Initialize application context
			ctx, err := app.New(mockMode)
			if err != nil {
				return fmt.Errorf("failed to initialize app: %w", err)
			}
			defer ctx.Close()

			// Verify task exists
			task, err := ctx.TaskStore.Get(taskID)
			if err != nil {
				return fmt.Errorf("task not found: %w", err)
			}

			fmt.Printf("📋 Task: %s\n", task.Title)
			fmt.Printf("📊 Status: %s\n", task.Status)
			fmt.Println()

			if task.Status != types.TaskStatusPending {
				return fmt.Errorf("task is not in pending state (current: %s)", task.Status)
			}

			if syncMode {
				// Run synchronously for debugging
				fmt.Println("🐛 Running in SYNC mode (for debugging)...")
				fmt.Println()

				sprint, err := ctx.SprintMgr.StartSprintSync(context.Background(), taskID)
				if err != nil {
					return fmt.Errorf("sprint failed: %w", err)
				}

				fmt.Printf("✅ Sprint completed: %s\n", sprint.ID)
				fmt.Printf("Final status: %s\n", sprint.Phase)
			} else {
				// Start the sprint asynchronously
				fmt.Println("Starting sprint workflow...")
				fmt.Println("(This will run asynchronously in the background)")
				fmt.Println()

				sprint, err := ctx.SprintMgr.StartSprint(context.Background(), taskID)
				if err != nil {
					return fmt.Errorf("failed to start sprint: %w", err)
				}

				fmt.Printf("✅ Sprint started: %s\n", sprint.ID)
				fmt.Println()
				fmt.Println("The sprint is now running through phases:")
				fmt.Println("  1. Intake - Creating initial specification")
				fmt.Println("  2. Spec Sprint - Tightening specification iteratively")
				fmt.Println("  3. Test Sprint - Generating comprehensive tests")
				fmt.Println("  4. Implementation Sprint - Writing code to pass tests")
				fmt.Println("  5. Review - Final validation and delivery")
				fmt.Println()
				fmt.Printf("Use 'scru-llm status %s' to check progress\n", taskID)
				fmt.Printf("Use 'scru-llm logs --task %s' to view detailed logs\n", taskID)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&mockMode, "mock", false, "Run in mock mode (no actual LLM calls)")
	cmd.Flags().BoolVar(&syncMode, "sync", false, "Run synchronously (for debugging)")

	return cmd
}

// NewDashboardCommand creates the dashboard command
func NewDashboardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Launch web dashboard",
		Long:  `Launch the web dashboard for viewing tasks and progress.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Dashboard not yet implemented in MVP.")
			fmt.Println("Use 'scru-llm status' for CLI-based status viewing.")
			return nil
		},
	}

	return cmd
}

// NewDeleteCommand creates the delete command
func NewDeleteCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <task-id>",
		Short: "Delete a task and its workspace",
		Long:  `Remove a task from the system and optionally delete its workspace.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			// Load config
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			// Get stores
			taskStore, sprintStore, eventLog, err := getStores(cfg)
			if err != nil {
				return err
			}
			defer eventLog.Close()

			// Get task to check if it exists and get workspace path
			task, err := taskStore.Get(taskID)
			if err != nil {
				return fmt.Errorf("task not found: %w", err)
			}

			// Confirm deletion unless --force
			if !force {
				fmt.Printf("Are you sure you want to delete task '%s' (%s)? [y/N]: ", task.Title, taskID)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Deletion cancelled.")
					return nil
				}
			}

			// Delete workspace if it exists
			if task.WorkspaceDir != "" {
				if err := os.RemoveAll(task.WorkspaceDir); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to delete workspace: %v\n", err)
				} else {
					fmt.Printf("🗑️  Deleted workspace: %s\n", task.WorkspaceDir)
				}
			}

			// Delete associated sprints
			sprints, _ := sprintStore.ListByTask(taskID)
			for _, sprint := range sprints {
				if err := sprintStore.Delete(sprint.ID); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to delete sprint %s: %v\n", sprint.ID, err)
				}
			}

			// Delete task
			if err := taskStore.Delete(taskID); err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			// Log event
			event := types.NewEvent("task.deleted", "cli", "Task deleted").
				WithTask(taskID)
			eventLog.Append(event)

			fmt.Printf("✅ Deleted task: %s\n", taskID)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

// min returns the minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
