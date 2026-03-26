// Package config provides configuration loading and validation for Scru-LLM.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config represents the complete application configuration
type Config struct {
	System    SystemConfig    `yaml:"system"`
	LLM       LLMConfig       `yaml:"llm"`
	Sprint    SprintConfig    `yaml:"sprint"`
	Container ContainerConfig `yaml:"container"`
	Reporting ReportingConfig `yaml:"reporting"`
	Cost      CostConfig      `yaml:"cost"`
	Security  SecurityConfig  `yaml:"security"`
	API       APIConfig       `yaml:"api"`
	Dev       DevConfig       `yaml:"development"`
}

// SystemConfig contains system-level settings
type SystemConfig struct {
	WorkspaceDir       string `yaml:"workspace_dir"`
	LogLevel           string `yaml:"log_level"`
	MaxConcurrentTasks int    `yaml:"max_concurrent_tasks"`
	LogDir             string `yaml:"log_dir"`
}

// LLMConfig contains LLM provider and role configuration
type LLMConfig struct {
	Providers []ProviderConfig      `yaml:"providers"`
	Roles     map[string]RoleConfig `yaml:"roles"`
}

// ProviderConfig defines an LLM provider
type ProviderConfig struct {
	Name         string `yaml:"name"`
	BaseURL      string `yaml:"base_url"`
	APIKey       string `yaml:"api_key"`
	DefaultModel string `yaml:"default_model"`
}

// RoleConfig defines configuration for a specific role
type RoleConfig struct {
	Provider    string  `yaml:"provider"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// SprintConfig contains sprint execution settings
type SprintConfig struct {
	SpecSprint           PhaseConfig   `yaml:"spec_sprint"`
	TestSprint           PhaseConfig   `yaml:"test_sprint"`
	ImplementationSprint PhaseConfig   `yaml:"implementation_sprint"`
	GlobalTimeout        time.Duration `yaml:"global_timeout"`
}

// PhaseConfig contains settings for a specific sprint phase
type PhaseConfig struct {
	MaxIterations       int           `yaml:"max_iterations"`
	TimeoutPerIteration time.Duration `yaml:"timeout_per_iteration"`
	MinCoveragePercent  int           `yaml:"min_coverage_percent,omitempty"`
}

// ContainerConfig contains container runtime settings
type ContainerConfig struct {
	Runtime        string            `yaml:"runtime"`
	Registry       string            `yaml:"registry,omitempty"`
	BaseImages     map[string]string `yaml:"base_images"`
	ResourceLimits ResourceLimits    `yaml:"resource_limits"`
	Volumes        []VolumeMount     `yaml:"volumes,omitempty"`
}

// ResourceLimits defines container resource constraints
type ResourceLimits struct {
	CPU     int           `yaml:"cpu"`
	Memory  string        `yaml:"memory"`
	Timeout time.Duration `yaml:"timeout"`
}

// VolumeMount defines a volume to mount in containers
type VolumeMount struct {
	HostPath      string `yaml:"host_path"`
	ContainerPath string `yaml:"container_path"`
	ReadOnly      bool   `yaml:"read_only"`
}

// ReportingConfig contains reporting and dashboard settings
type ReportingConfig struct {
	Dashboard DashboardConfig `yaml:"dashboard"`
	AuditLog  AuditLogConfig  `yaml:"audit_log"`
	AgentLogs AgentLogsConfig `yaml:"agent_logs"`
	Metrics   MetricsConfig   `yaml:"metrics"`
}

// DashboardConfig contains web dashboard settings
type DashboardConfig struct {
	Enabled         bool          `yaml:"enabled"`
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	RefreshInterval time.Duration `yaml:"refresh_interval"`
}

// AuditLogConfig contains audit logging settings
type AuditLogConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Path         string `yaml:"path"`
	RotationDays int    `yaml:"rotation_days"`
}

// AgentLogsConfig contains agent activity logging settings
type AgentLogsConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Path          string `yaml:"path"`
	RetentionDays int    `yaml:"retention_days"`
}

// MetricsConfig contains metrics export settings
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// CostConfig contains cost control settings
type CostConfig struct {
	BudgetPerTask          float64 `yaml:"budget_per_task"`
	WarningThreshold       float64 `yaml:"warning_threshold"`
	EmergencyStopThreshold float64 `yaml:"emergency_stop_threshold"`
	TrackByProvider        bool    `yaml:"track_by_provider"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	SandboxEnabled     bool     `yaml:"sandbox_enabled"`
	AllowNetworkAccess bool     `yaml:"allow_network_access"`
	AllowedCommands    []string `yaml:"allowed_commands"`
	TrustedCIDRs       []string `yaml:"trusted_cidrs"`
}

// APIConfig contains API server settings
type APIConfig struct {
	Enabled          bool   `yaml:"enabled"`
	Port             int    `yaml:"port"`
	Host             string `yaml:"host"`
	WebSocketEnabled bool   `yaml:"websocket_enabled"`
	RequireAuth      bool   `yaml:"require_auth"`
	AuthMethod       string `yaml:"auth_method"`
}

// DevConfig contains development settings
type DevConfig struct {
	Debug          bool `yaml:"debug"`
	SaveLLMTraces  bool `yaml:"save_llm_traces"`
	KeepWorkspaces bool `yaml:"keep_workspaces"`
	MockMode       bool `yaml:"mock_mode"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		System: SystemConfig{
			WorkspaceDir:       "~/.scru-llm/workspaces",
			LogLevel:           "info",
			MaxConcurrentTasks: 3,
			LogDir:             "~/.scru-llm/logs",
		},
		LLM: LLMConfig{
			Providers: []ProviderConfig{
				{
					Name:         "openrouter",
					BaseURL:      "https://openrouter.ai/api/v1",
					APIKey:       "${OPENROUTER_API_KEY}",
					DefaultModel: "anthropic/claude-3.5-sonnet",
				},
			},
			Roles: map[string]RoleConfig{
				"product_owner": {
					Provider:    "openrouter",
					Model:       "qwen/qwen3.5-35b-a3b",
					Temperature: 0.15,
					MaxTokens:   7000,
				},
				"scrum_master": {
					Provider:    "openrouter",
					Model:       "qwen/qwen3.5-35b-a3b",
					Temperature: 0.15,
					MaxTokens:   8192,
				},
				"spec_engineer": {
					Provider:    "openrouter",
					Model:       "qwen/qwen3.5-35b-a3b",
					Temperature: 0.15,
					MaxTokens:   7000,
				},
				"test_engineer": {
					Provider:    "openrouter",
					Model:       "qwen/qwen3.5-coder-32b",
					Temperature: 0.1,
					MaxTokens:   3000,
				},
				"code_engineer": {
					Provider:    "openrouter",
					Model:       "qwen/qwen3.5-coder-32b",
					Temperature: 0.1,
					MaxTokens:   12000,
				},
				"reviewer": {
					Provider:    "openrouter",
					Model:       "mistralai/mistral-small-24b-instruct-2501",
					Temperature: 0.05,
					MaxTokens:   4000,
				},
			},
		},
		Sprint: SprintConfig{
			SpecSprint: PhaseConfig{
				MaxIterations:       3,
				TimeoutPerIteration: 5 * time.Minute,
			},
			TestSprint: PhaseConfig{
				MaxIterations:       3,
				TimeoutPerIteration: 10 * time.Minute,
				MinCoveragePercent:  80,
			},
			ImplementationSprint: PhaseConfig{
				MaxIterations:       5,
				TimeoutPerIteration: 15 * time.Minute,
			},
			GlobalTimeout: 1 * time.Hour,
		},
		Container: ContainerConfig{
			Runtime:  "docker",
			Registry: "",
			BaseImages: map[string]string{
				"python": "python:3.11-slim",
				"node":   "node:18-alpine",
				"go":     "golang:1.21-alpine",
				"rust":   "rust:1.75-slim",
				"java":   "openjdk:21-slim",
			},
			ResourceLimits: ResourceLimits{
				CPU:     2,
				Memory:  "4g",
				Timeout: 10 * time.Minute,
			},
		},
		Reporting: ReportingConfig{
			Dashboard: DashboardConfig{
				Enabled:         true,
				Port:            8080,
				Host:            "127.0.0.1",
				RefreshInterval: 5 * time.Second,
			},
			AuditLog: AuditLogConfig{
				Enabled:      true,
				Path:         "~/.scru-llm/logs/audit.jsonl",
				RotationDays: 30,
			},
			AgentLogs: AgentLogsConfig{
				Enabled:       true,
				Path:          "~/.scru-llm/logs/agents/",
				RetentionDays: 14,
			},
			Metrics: MetricsConfig{
				Enabled: false,
				Port:    9090,
				Path:    "/metrics",
			},
		},
		Cost: CostConfig{
			BudgetPerTask:          5.00,
			WarningThreshold:       0.8,
			EmergencyStopThreshold: 1.0,
			TrackByProvider:        true,
		},
		Security: SecurityConfig{
			SandboxEnabled:     true,
			AllowNetworkAccess: false,
			AllowedCommands:    []string{},
			TrustedCIDRs: []string{
				"127.0.0.1/32",
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
			},
		},
		API: APIConfig{
			Enabled:          true,
			Port:             8081,
			Host:             "127.0.0.1",
			WebSocketEnabled: true,
			RequireAuth:      true,
			AuthMethod:       "token",
		},
		Dev: DevConfig{
			Debug:          false,
			SaveLLMTraces:  false,
			KeepWorkspaces: false,
			MockMode:       false,
		},
	}
}

// ExpandPath expands ~ to home directory
func ExpandPath(path string) (string, error) {
	if path == "" {
		return path, nil
	}

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	}

	return path, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate system config
	if c.System.MaxConcurrentTasks < 1 {
		return fmt.Errorf("max_concurrent_tasks must be at least 1")
	}

	// Validate sprint config
	if c.Sprint.SpecSprint.MaxIterations < 1 {
		return fmt.Errorf("spec_sprint.max_iterations must be at least 1")
	}
	if c.Sprint.TestSprint.MaxIterations < 1 {
		return fmt.Errorf("test_sprint.max_iterations must be at least 1")
	}
	if c.Sprint.ImplementationSprint.MaxIterations < 1 {
		return fmt.Errorf("implementation_sprint.max_iterations must be at least 1")
	}

	// Validate LLM config
	if len(c.LLM.Providers) == 0 {
		return fmt.Errorf("at least one LLM provider must be configured")
	}

	// Check that all referenced providers exist
	providerNames := make(map[string]bool)
	for _, p := range c.LLM.Providers {
		providerNames[p.Name] = true
	}

	for role, roleConfig := range c.LLM.Roles {
		if !providerNames[roleConfig.Provider] {
			return fmt.Errorf("role %q references unknown provider %q", role, roleConfig.Provider)
		}
	}

	// Validate cost config
	if c.Cost.BudgetPerTask <= 0 {
		return fmt.Errorf("budget_per_task must be positive")
	}

	return nil
}

// ExpandPaths expands all path fields in the configuration
func (c *Config) ExpandPaths() error {
	var err error

	c.System.WorkspaceDir, err = ExpandPath(c.System.WorkspaceDir)
	if err != nil {
		return fmt.Errorf("failed to expand workspace_dir: %w", err)
	}

	c.System.LogDir, err = ExpandPath(c.System.LogDir)
	if err != nil {
		return fmt.Errorf("failed to expand log_dir: %w", err)
	}

	c.Reporting.AuditLog.Path, err = ExpandPath(c.Reporting.AuditLog.Path)
	if err != nil {
		return fmt.Errorf("failed to expand audit_log.path: %w", err)
	}

	c.Reporting.AgentLogs.Path, err = ExpandPath(c.Reporting.AgentLogs.Path)
	if err != nil {
		return fmt.Errorf("failed to expand agent_logs.path: %w", err)
	}

	return nil
}
