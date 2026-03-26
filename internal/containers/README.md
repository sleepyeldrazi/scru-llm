# Containers Package

This package manages containerized execution environments for tasks.

## What This Package Does

Provides isolated, reproducible environments for:
- Building code
- Running tests
- Executing commands

Supports multiple container runtimes: Docker, Podman, Containerd.

## Files

- `runtime.go` - Container runtime interface
- `docker.go` - Docker implementation
- `builder.go` - Dockerfile generation and image building
- `sandbox.go` - Container execution with resource limits

## Architecture

### ContainerRuntime Interface
```go
type ContainerRuntime interface {
    Build(ctx context.Context, dockerfile, tag string) error
    Run(ctx context.Context, image string, cmd []string, opts RunOptions) (*Result, error)
    Remove(ctx context.Context, containerID string) error
    ImageExists(ctx context.Context, image string) (bool, error)
}
```

### Docker Implementation
- Uses Docker SDK for Go
- Supports Dockerfile builds
- Resource limits (CPU, memory, timeout)
- Volume mounting for workspaces
- Log capture

### Builder
- Generates Dockerfiles from templates
- Uses base images per language
- Layer caching optimization
- Multi-stage builds

### Sandbox
- Creates isolated container per task
- Mounts workspace as volume
- Enforces resource limits
- Captures all output

## Usage Example

```go
// Create runtime
runtime := containers.NewDockerRuntime()

// Build image
err := runtime.Build(ctx, dockerfile, "task-123:latest")

// Run tests
result, err := runtime.Run(ctx, "task-123:latest", 
    []string{"python", "-m", "pytest"},
    containers.RunOptions{
        WorkDir: "/workspace",
        Timeout: 10 * time.Minute,
    })

// Check results
if result.ExitCode != 0 {
    fmt.Println("Tests failed:", string(result.Stderr))
}
```

## Security

- No privileged containers
- Read-only root filesystem (optional)
- Network isolation (by default)
- Resource quotas
- User namespace remapping (Docker)

## Future Improvements

- [ ] Podman support
- [ ] Containerd support
- [ ] BuildKit integration
- [ ] Distributed builds
- [ ] Image layer caching optimization
