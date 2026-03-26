package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // It's okay if .env doesn't exist
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			// Only set if not already set in environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
	return scanner.Err()
}

// Loader handles configuration loading from files and environment
type Loader struct {
	viper *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	v := viper.New()

	// Set config name and type
	v.SetConfigName("scru-llm")
	v.SetConfigType("yaml")

	// Add config paths
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.scru-llm")
	v.AddConfigPath("/etc/scru-llm")

	// Enable environment variables
	v.SetEnvPrefix("SCRULLM")
	v.AutomaticEnv()

	return &Loader{viper: v}
}

// Load loads configuration from files and environment
func (l *Loader) Load() (*Config, error) {
	// Load .env file first (so it can be overridden by actual env vars)
	// Try multiple locations for .env
	envPaths := []string{
		".env",
		"./.env",
		filepath.Join(os.Getenv("HOME"), ".scru-llm", ".env"),
	}
	for _, envPath := range envPaths {
		if err := loadEnvFile(envPath); err != nil {
			// Log but don't fail - .env is optional
			fmt.Fprintf(os.Stderr, "Warning: could not load .env from %s: %v\n", envPath, err)
		}
	}

	// Start with defaults
	cfg := DefaultConfig()

	// Try to read config file
	if err := l.viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal into our struct
	if err := l.viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand paths
	if err := cfg.ExpandPaths(); err != nil {
		return nil, fmt.Errorf("failed to expand paths: %w", err)
	}

	// Substitute environment variables in provider API keys
	for i := range cfg.LLM.Providers {
		cfg.LLM.Providers[i].APIKey = os.ExpandEnv(cfg.LLM.Providers[i].APIKey)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// LoadFromFile loads configuration from a specific file path
func LoadFromFile(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand paths
	if err := cfg.ExpandPaths(); err != nil {
		return nil, fmt.Errorf("failed to expand paths: %w", err)
	}

	// Substitute environment variables
	for i := range cfg.LLM.Providers {
		cfg.LLM.Providers[i].APIKey = os.ExpandEnv(cfg.LLM.Providers[i].APIKey)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// LoadOrDefault loads configuration or returns defaults if no config exists
func LoadOrDefault() (*Config, error) {
	loader := NewLoader()
	cfg, err := loader.Load()
	if err != nil {
		// If it's a validation error, return it
		return nil, err
	}
	return cfg, nil
}

// EnsureDirs creates necessary directories for the configuration
func (c *Config) EnsureDirs() error {
	dirs := []string{
		c.System.WorkspaceDir,
		c.System.LogDir,
		c.Reporting.AgentLogs.Path,
		filepath.Dir(c.Reporting.AuditLog.Path),
	}

	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetProvider returns a provider configuration by name
func (c *Config) GetProvider(name string) (*ProviderConfig, error) {
	for _, p := range c.LLM.Providers {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("provider %q not found", name)
}

// GetRoleConfig returns the configuration for a specific role
func (c *Config) GetRoleConfig(role string) (*RoleConfig, error) {
	if cfg, ok := c.LLM.Roles[role]; ok {
		return &cfg, nil
	}
	return nil, fmt.Errorf("role %q not found", role)
}

// GetBaseImage returns the base image for a language
func (c *Config) GetBaseImage(language string) (string, error) {
	if image, ok := c.Container.BaseImages[language]; ok {
		return image, nil
	}
	return "", fmt.Errorf("no base image configured for language %q", language)
}
