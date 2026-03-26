# Scru-LLM Coding Guidelines

Coding standards and best practices for the Scru-LLM project.

## 1. Core Principles

### 1.1 Simplicity First
- **Write simple code**: If you can do it in 10 lines instead of 50, do it
- **Avoid over-engineering**: Start concrete, abstract only when you have 2+ implementations
- **No enterprise patterns**: No factories of factories, no complex inheritance hierarchies
- **Clarity over cleverness**: Code should be readable first, compact second

**Example - Bad:**
```go
// Factory for a single implementation
type ContainerFactory interface {
    Create(runtime string) ContainerRuntime
}

type DockerFactory struct{}
func (f *DockerFactory) Create(runtime string) ContainerRuntime {
    return &DockerRuntime{}
}
```

**Example - Good:**
```go
// Direct creation
func NewDockerRuntime() *DockerRuntime {
    return &DockerRuntime{}
}
```

### 1.2 Modularity Without Premature Abstraction
- Define interfaces at the **usage site**, not implementation site
- Start with concrete types, extract interfaces when needed
- Keep interfaces small (ideally 1-3 methods)
- Use composition over inheritance

**Example:**
```go
// Define interface where it's used
type TaskStore interface {
    Get(id string) (*Task, error)
    Save(task *Task) error
}

// Multiple implementations
type SQLiteTaskStore struct { db *sql.DB }
type MemoryTaskStore struct { tasks map[string]*Task }
```

### 1.3 Configuration Over Hardcoding
Never hardcode:
- API endpoints (except localhost defaults)
- Port numbers, timeouts, magic numbers
- File paths
- Model names or versions
- Error messages

**Example:**
```go
// Bad
const DefaultPort = 8080

// Good
type Config struct {
    Server struct {
        Port    int           `yaml:"port" default:"8080"`
        Timeout time.Duration `yaml:"timeout" default:"30s"`
    }
}
```

### 1.4 DRY - Don't Repeat Yourself
**Before writing code:**
1. Search for existing functionality: `grep -r "func.*DoSomething" internal/`
2. Check if similar functions exist
3. Abstract only when you have 3+ uses or complex logic

**When NOT to abstract:**
- Used only once or twice
- Abstraction makes code harder to understand
- Abstracting different things that happen to look similar

## 2. Project Structure

```
scru-llm/
├── cmd/scru-llm/              # Main CLI entrypoint
│   ├── main.go                 # Single file, minimal
│   └── commands/               # CLI command implementations
│       ├── task.go
│       ├── dashboard.go
│       └── config.go
│
├── internal/                   # Private application code
│   ├── scrum/                  # Scrum orchestration
│   │   ├── backlog.go          # Task queue management
│   │   ├── sprint.go           # Sprint lifecycle
│   │   ├── events.go           # Event bus
│   │   └── README.md
│   │
│   ├── workers/                # LLM agent implementations
│   │   ├── product_owner.go
│   │   ├── scrum_master.go
│   │   ├── spec_engineer.go
│   │   ├── test_engineer.go
│   │   ├── code_engineer.go
│   │   ├── reviewer.go
│   │   ├── base.go             # Shared worker functionality
│   │   └── README.md
│   │
│   ├── containers/             # Container management
│   │   ├── runtime.go          # Container runtime interface
│   │   ├── docker.go           # Docker implementation
│   │   ├── builder.go          # Image building
│   │   └── README.md
│   │
│   ├── reporting/              # Progress reporting
│   │   ├── dashboard.go        # Web dashboard
│   │   ├── kanban.go           # Kanban board generation
│   │   ├── audit.go            # Audit logging
│   │   └── README.md
│   │
│   ├── config/                 # Configuration
│   │   ├── loader.go           # Config file loading
│   │   ├── validation.go       # Config validation
│   │   └── README.md
│   │
│   ├── llm/                    # LLM client
│   │   ├── client.go           # HTTP client
│   │   ├── router.go           # Provider routing
│   │   ├── models.go           # Model definitions
│   │   └── README.md
│   │
│   └── filesystem/             # Workspace management
│       ├── sandbox.go          # Isolated workspaces
│       ├── git.go              # Git operations
│       └── README.md
│
├── config/                     # Configuration templates
│   ├── scru-llm.example.yaml
│   └── prompts/                # LLM prompts
│       ├── product-owner-system.md
│       ├── spec-engineer-system.md
│       └── ...
│
├── docs/                       # Documentation
│   ├── SPECIFICATIONS.md
│   ├── CODING_GUIDELINES.md
│   └── architecture/           # Design docs
│
├── scripts/                    # Utility scripts
│   ├── setup.sh
│   └── benchmark.sh
│
├── tests/                      # Integration tests
│   └── fixtures/               # Test data
│
├── go.mod
├── go.sum
├── README.md
└── TODO.md
```

## 3. File Organization Rules

### 3.1 File Size
- **Maximum**: 500 lines per file (ideally <300)
- If larger, split by functionality
- Exception: Generated code, test data

### 3.2 Package Size
- **Maximum**: 10-15 files per package
- Then consider sub-packages
- Exception: Large test suites

### 3.3 Every Module Gets a README
Each `internal/<module>/` directory MUST have a `README.md`:

```markdown
# Module Name

Brief description of what this module does.

## Files

- `file1.go` - Purpose
- `file2.go` - Purpose

## Architecture

How this fits into the bigger picture.

## Usage Example

```go
// Brief example
```

## Future Improvements

- Known limitations
- Planned enhancements
```

## 4. Naming Conventions

### 4.1 Packages
- Short, lowercase, no underscores: `scrum`, `workers`, `containers`
- Describe functionality, not implementation: `llm` not `openrouter`

### 4.2 Types
- PascalCase: `TaskManager`, `SprintConfig`
- Avoid stutter: `scrum.ScrumMaster` → `scrum.Master`
- Interface names: `-er` suffix for single method: `Reader`, `Handler`

### 4.3 Functions
- PascalCase for exported, camelCase for unexported
- Use verbs: `New`, `Create`, `Get`, `Update`, `Delete`, `Handle`
- Keep short but descriptive: `GetTask` not `RetrieveTaskInformationFromStorage`

### 4.4 Variables
- Short names in small scopes: `i`, `n`, `err`
- Longer names in larger scopes: `taskManager`, `configLoader`
- Acronyms uppercase: `HTTP`, `URL`, `ID`, `API`

### 4.5 Constants
- PascalCase for exported, camelCase for unexported
- Avoid ALL_CAPS (not idiomatic Go)

**Example:**
```go
const (
    DefaultTimeout = 30 * time.Second
    maxRetries     = 3
)
```

## 5. Error Handling

### 5.1 Principles
- Return errors, don't panic (except in `main` or tests)
- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Never ignore errors: Always check `if err != nil`
- Create sentinel errors for specific cases

**Example:**
```go
var (
    ErrTaskNotFound    = errors.New("task not found")
    ErrSprintFailed    = errors.New("sprint failed")
    ErrMaxIterations   = errors.New("max iterations reached")
)

func GetTask(id string) (*Task, error) {
    task, err := db.Find(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("%w: %s", ErrTaskNotFound, id)
        }
        return nil, fmt.Errorf("database error fetching task %s: %w", id, err)
    }
    return task, nil
}
```

### 5.2 Error Messages
- Start with lowercase (often chained)
- Be descriptive but concise
- Include relevant context (IDs, names)

## 6. Testing

### 6.1 Testing Requirements
- Every module MUST have unit tests
- Test files alongside implementation: `server.go` + `server_test.go`
- Use table-driven tests
- Mock at interface boundaries
- Aim for >70% coverage on critical paths

**Example:**
```go
func TestSprintManager_StartSprint(t *testing.T) {
    tests := []struct {
        name      string
        task      *Task
        wantErr   bool
        errType   error
    }{
        {
            name:    "valid task",
            task:    &Task{ID: "task-1", Title: "Test"},
            wantErr: false,
        },
        {
            name:    "nil task",
            task:    nil,
            wantErr: true,
            errType: ErrInvalidTask,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sm := NewSprintManager()
            sprint, err := sm.StartSprint(tt.task)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.ErrorIs(t, err, tt.errType)
                }
                return
            }
            
            assert.NoError(t, err)
            assert.NotNil(t, sprint)
            assert.Equal(t, tt.task.ID, sprint.TaskID)
        })
    }
}
```

### 6.2 Testing LLM Agents
Since agents use LLM calls:

**Strategy:**
- Mock LLM client for unit tests
- Use fixtures for LLM responses
- Test logic, not the LLM itself
- Integration tests use real LLM with small models

**Example:**
```go
// Mock LLM client for testing
type MockLLMClient struct {
    responses map[string]string
}

func (m *MockLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
    if resp, ok := m.responses[prompt]; ok {
        return resp, nil
    }
    return "", fmt.Errorf("unexpected prompt: %s", prompt)
}

func TestSpecEngineer_TightenSpec(t *testing.T) {
    mockLLM := &MockLLMClient{
        responses: map[string]string{
            "tighten spec": `{"requirements": ["clear requirement"]}`,
        },
    }
    
    engineer := NewSpecEngineer(mockLLM)
    spec, err := engineer.TightenSpec(context.Background(), initialSpec)
    
    assert.NoError(t, err)
    assert.NotNil(t, spec)
    assert.Len(t, spec.Requirements, 1)
}
```

### 6.3 Integration Tests
- Located in `tests/` directory
- Use real dependencies (test database, container runtime)
- Clean up after tests
- Use `testify/suite` for complex scenarios

## 7. Documentation

### 7.1 Code Comments
Every exported type/function needs a comment:
- Start with the name of the thing
- Explain WHY, not WHAT (code shows what)

**Good:**
```go
// SprintManager coordinates the lifecycle of sprints from start to completion.
// It manages iteration budgets, tracks progress, and handles agent assignments.
// The manager is safe for concurrent use by multiple goroutines.
type SprintManager struct {
    // ...
}

// StartSprint begins a new sprint for the given task. It initializes the sprint
// state, assigns the appropriate agents, and begins the spec phase. Returns
// ErrInvalidTask if the task is nil or missing required fields.
func (sm *SprintManager) StartSprint(task *Task) (*Sprint, error) {
    // ...
}
```

**Bad:**
```go
// Process processes the data
func Process(data []byte) error { ... }
```

### 7.2 Package Documentation
```go
// Package scrum provides the Scrum orchestration layer for Scru-LLM.
// It manages sprints, backlogs, and the coordination between different
// agent roles following Scrum methodology.
package scrum
```

## 8. Concurrency

### 8.1 Goroutine Safety
- Protect shared state with mutexes
- Prefer channels for coordination
- Always have exit conditions for goroutines

**Example:**
```go
type EventBus struct {
    mu       sync.RWMutex
    handlers map[string][]Handler
    closed   bool
}

func (eb *EventBus) Subscribe(eventType string, handler Handler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    if eb.closed {
        return
    }
    
    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *EventBus) Close() {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.closed = true
    // Clean up resources
}
```

### 8.2 Context Usage
- Always accept `context.Context` as first parameter
- Respect cancellation
- Pass context through call chain

**Example:**
```go
func (w *Worker) Execute(ctx context.Context, task *Task) error {
    ctx, cancel := context.WithTimeout(ctx, w.timeout)
    defer cancel()
    
    select {
    case <-ctx.Done():
        return fmt.Errorf("task execution timeout: %w", ctx.Err())
    case result := <-w.runTask(ctx, task):
        return result
    }
}
```

## 9. LLM Prompts

### 9.1 Prompt Organization
- Store prompts in `config/prompts/`
- One file per role/system prompt
- Use Go templates for dynamic content
- Version control prompts with code

### 9.2 Prompt Structure
```markdown
# Role: Spec Engineer

## Objective
Tighten specifications by removing ambiguity and making requirements measurable.

## Constraints
- Only tighten, don't change intent
- Maintain all original requirements
- Add edge cases where appropriate

## Input Format
Task specification in JSON format...

## Output Format
Return JSON with...

## Examples
[If needed]
```

### 9.3 Prompt Guidelines
- Keep prompts focused and constraint-based
- Avoid persona-heavy language
- Prefer rules over examples
- Include escape hatches: "say unknown if..."

## 10. Git Practices

### 10.1 Commit Messages
Format: `<type>: <subject>`

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Code style (formatting)
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Build/dependencies

Rules:
- Use imperative mood ("Add" not "Added")
- No period at end
- Max 50 characters for subject
- Body explains what and why (not how)

**Example:**
```
feat: add container runtime abstraction

Introduce ContainerRuntime interface to support Docker, Podman,
and Containerd backends. This allows users to choose their preferred
container runtime without code changes.

Includes:
- Docker runtime implementation
- Runtime auto-detection
- Configuration options
```

### 10.2 Branch Strategy
- `main` - Production-ready code
- `feature/<description>` - New features
- `fix/<description>` - Bug fixes
- `refactor/<description>` - Refactoring

### 10.3 Pre-Commit Checklist
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] `gofmt -w .` run
- [ ] `go vet ./...` clean
- [ ] No secrets in code
- [ ] Commit message follows format

## 11. Dependencies

### 11.1 Principles
- Minimize external dependencies
- Prefer standard library
- If external, prefer mature, widely-used libraries
- Pin versions in go.mod
- Document WHY for each dependency

### 11.2 Approved Dependencies
Current approved dependencies:
- `github.com/labstack/echo/v4` - HTTP framework
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `github.com/stretchr/testify` - Testing utilities

### 11.3 Adding New Dependencies
1. Can stdlib do it reasonably well?
2. Is there already a similar dep?
3. Is it actively maintained?
4. Is it widely used (check pkg.go.dev)?
5. Document the justification in commit message

## 12. Configuration

### 12.1 Configuration Schema
Use structs with tags for YAML loading:

```go
type Config struct {
    System struct {
        WorkspaceDir       string `yaml:"workspace_dir" default:"~/.scru-llm/workspaces"`
        LogLevel          string `yaml:"log_level" default:"info"`
        MaxConcurrentTasks int    `yaml:"max_concurrent_tasks" default:"3"`
    } `yaml:"system"`
    
    LLM struct {
        Providers []ProviderConfig `yaml:"providers"`
        Roles     map[string]RoleConfig `yaml:"roles"`
    } `yaml:"llm"`
}
```

### 12.2 Validation
Validate configuration at startup:
- Required fields present
- Valid paths and permissions
- Network connectivity to providers
- Sensible values (timeouts > 0, etc.)

### 12.3 Defaults
- Provide sensible defaults
- Document defaults in example config
- Use struct tags for defaults where possible

## 13. Logging

### 13.1 Log Levels
- `ERROR` - Unrecoverable failures
- `WARN` - Recoverable issues, unusual conditions
- `INFO` - Normal operations, milestones
- `DEBUG` - Detailed traces, development only

### 13.2 Structured Logging
Use structured logging with key-value pairs:

```go
log.Info("sprint started",
    "sprint_id", sprint.ID,
    "task_id", task.ID,
    "phase", sprint.Phase,
)

log.Error("task execution failed",
    "task_id", task.ID,
    "error", err,
    "duration", time.Since(start),
)
```

### 13.3 What to Log
**Always log:**
- Sprint state transitions
- Agent assignments
- Iteration completions
- Errors and failures
- Container lifecycle events

**Never log:**
- API keys or secrets
- Full LLM prompts (in production)
- Large binary data
- User passwords

## 14. Agent-Friendly Code

Since agents will be reading and potentially modifying this code:

### 14.1 Make Intent Obvious
```go
// Instead of complex inline logic
if x > 0 && x < 100 && x%2 == 0 { ... }

// Use named constants and functions
const (
    MinValidID = 1
    MaxValidID = 99
)

func isValidID(id int) bool {
    return id >= MinValidID && id <= MaxValidID && id%2 == 0
}
```

### 14.2 Explicit Types
```go
// Use type aliases for domain concepts
type TaskID string
type SprintID string
type AgentID string

// Instead of primitives everywhere
func AssignTask(agentID AgentID, taskID TaskID) error
```

### 14.3 Document Complex Logic
```go
// Binary search for insertion point. Returns the index where the target
// should be inserted to maintain sorted order. O(log n) vs O(n) for linear
// search, important when processing thousands of items.
func findInsertionPoint(items []Item, target Item) int { ... }
```

---

## Code Review Checklist

Before submitting code:

- [ ] No hardcoded values (use config)
- [ ] No duplicate code (check existing modules)
- [ ] File has <500 lines
- [ ] Module has README.md
- [ ] All exported items have comments
- [ ] Errors are handled and wrapped
- [ ] Unit tests exist for new functionality
- [ ] `go vet` and `gofmt` pass
- [ ] No unnecessary dependencies
- [ ] Goroutines have proper lifecycle management
- [ ] Context is passed and respected
- [ ] No secrets in code
- [ ] Commit message follows format

---

*Keep it simple, make it modular, document everything, test thoroughly.*
