# Config Package

This package handles configuration loading and validation for Scru-LLM.

## What This Package Does

- Loads configuration from YAML files
- Supports environment variable substitution
- Validates configuration
- Provides defaults
- Hot reloading (optional)

## Files

- `loader.go` - Configuration file loading
- `validation.go` - Configuration validation
- `defaults.go` - Default values

## Configuration Sources

Priority (highest to lowest):
1. Environment variables
2. `~/.scru-llm/scru-llm.yaml`
3. `./scru-llm.yaml`
4. `config/scru-llm.example.yaml` (defaults)

## Usage Example

```go
// Load configuration
cfg, err := config.Load()
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Access values
fmt.Println("Workspace:", cfg.System.WorkspaceDir)
fmt.Println("Max concurrent tasks:", cfg.System.MaxConcurrentTasks)

// Get model for role
roleConfig := cfg.LLM.Roles["code_engineer"]
fmt.Println("Model:", roleConfig.Model)
```

## Configuration Structure

See `config/scru-llm.example.yaml` for complete structure.

Key sections:
- `system` - Workspace, logging
- `llm` - Providers, role mappings
- `sprint` - Iteration budgets, timeouts
- `container` - Runtime, base images
- `reporting` - Dashboard, audit logs
- `cost` - Budgets, thresholds
- `security` - Sandboxing, network access

## Validation

Configuration is validated at load time:
- Required fields present
- Valid paths and permissions
- Network connectivity (optional check)
- Sensible values (timeouts > 0, etc.)

## Future Improvements

- [ ] Configuration hot reloading
- [ ] Remote configuration (S3, etcd)
- [ ] Configuration versioning
- [ ] Secret management integration (Vault, etc.)
