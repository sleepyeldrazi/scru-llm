# Scru-LLM Development Roadmap

This document tracks the development plan for Scru-LLM, organized by priority and phase.

Snapshot:
- this repo is currently a design and handoff skeleton
- the roadmap includes both already-created scaffolding and future implementation work
- consult `docs/CURRENT_STATE.md` before treating any roadmap item as implemented runtime behavior

---

## Legend

- `[ ]` - Not started
- `[-]` - In progress
- `[x]` - Completed
- `[~]` - Deferred/Blocked
- `[?]` - Needs clarification

---

## Phase 0: Foundation (Prerequisites)

**Goal**: Project structure, tooling, and core abstractions in place.

**Duration**: 2-3 weeks

### P0.1 - Project Setup
- [x] Initialize Go module (`go mod init scru-llm`)
- [x] Set up basic directory structure
- [x] Create `.gitignore` with Go and project-specific patterns
- [x] Initialize git repository
- [ ] Set up CI/CD pipeline (GitHub Actions)
  - [ ] Linting (`golangci-lint`)
  - [ ] Unit tests
  - [ ] Build verification
- [x] Create Makefile with common tasks
  - [x] `make build`
  - [x] `make test`
  - [x] `make lint`
  - [x] `make clean`

### P0.2 - Core Abstractions
- [ ] Define core domain types
  - [ ] `Task` struct with all fields
  - [ ] `Sprint` struct
  - [ ] `Specification` struct
  - [ ] `Agent` interface
  - [ ] `Iteration` struct
- [ ] Define error types and sentinel errors
  - [ ] `ErrTaskNotFound`
  - [ ] `ErrSprintFailed`
  - [ ] `ErrMaxIterations`
  - [ ] `ErrContainerFailed`
- [ ] Define event types for event bus
  - [ ] Base `Event` interface
  - [ ] Specific event types (TaskCreated, SprintStarted, etc.)

### P0.3 - Configuration System
- [ ] Create config package
  - [ ] Config struct definitions
  - [ ] YAML loading with `viper`
  - [ ] Environment variable substitution
  - [ ] Validation logic
  - [ ] Default values
- [x] Create `config/scru-llm.example.yaml`
- [ ] Write config tests

### P0.4 - Logging & Observability
- [ ] Set up structured logging
  - [ ] Use `log/slog` (Go 1.21+)
  - [ ] JSON formatter for production
  - [ ] Text formatter for development
- [ ] Create audit log infrastructure
  - [ ] `AuditLogger` interface
  - [ ] JSONL file implementation
  - [ ] Log rotation
- [ ] Create agent activity logger
  - [ ] Track file reads/writes
  - [ ] Track LLM calls
  - [ ] Track command execution

---

## Phase 1: Core Infrastructure

**Goal**: Working infrastructure for containers, LLM routing, and workspaces.

**Duration**: 3-4 weeks

### P1.1 - Container Runtime
- [ ] Define `ContainerRuntime` interface
  - [ ] `Build(ctx, dockerfile, tag) error`
  - [ ] `Run(ctx, image, cmd, opts) (*Result, error)`
  - [ ] `Remove(ctx, containerID) error`
- [ ] Implement Docker runtime
  - [ ] Use Docker SDK
  - [ ] Build from Dockerfile
  - [ ] Run with resource limits
  - [ ] Capture stdout/stderr
  - [ ] Handle timeouts
- [ ] Implement Podman runtime (optional, similar interface)
- [ ] Create base image templates
  - [ ] Python 3.11
  - [ ] Node.js 18
  - [ ] Go 1.21
  - [ ] Rust 1.75
- [ ] Write container tests
  - [ ] Unit tests with mocks
  - [ ] Integration tests with real Docker

### P1.2 - LLM Client & Router
- [ ] Define `LLMClient` interface
  - [ ] `Complete(ctx, messages, opts) (*Response, error)`
  - [ ] `Stream(ctx, messages, opts) (chan Chunk, error)`
- [ ] Implement OpenAI-compatible client
  - [ ] Support OpenAI API
  - [ ] Support OpenRouter
  - [ ] Support local models (llama.cpp, etc.)
- [ ] Create router for role-to-model mapping
  - [ ] Load configuration
  - [ ] Route based on role
  - [ ] Track costs per call
  - [ ] Implement fallback chains
- [ ] Add retry logic with exponential backoff
- [ ] Write LLM client tests
  - [ ] Mock client for unit tests
  - [ ] Integration tests (optional, with small models)

### P1.3 - Workspace Management
- [ ] Create filesystem sandbox
  - [ ] `Sandbox` struct
  - [ ] Create isolated workspace directory
  - [ ] Enforce path restrictions
  - [ ] Cleanup on completion
- [ ] Implement workspace structure
  ```
  workspaces/{task_id}/
  ├── spec/           # Specification versions
  ├── tests/          # Test files
  ├── src/            # Source code
  ├── containers/     # Dockerfile and build files
  ├── logs/           # Execution logs
  └── artifacts/      # Build outputs
  ```
- [ ] Add Git integration
  - [ ] Initialize repo per workspace
  - [ ] Commit after each iteration
  - [ ] Tag versions (spec-v1, test-v1, impl-v1)
- [ ] Write workspace tests

### P1.4 - Event Bus
- [ ] Implement `EventBus` interface
  - [ ] `Publish(event)` error
  - [ ] `Subscribe(eventType, handler)` Subscription
  - [ ] `Close()` error
- [ ] Create in-memory implementation
  - [ ] Async processing with goroutines
  - [ ] Buffer management
  - [ ] Graceful shutdown
- [ ] Add persistent event log
  - [ ] JSONL file appender
  - [ ] Event replay capability
- [ ] Write event bus tests

---

## Phase 2: Scrum Orchestration

**Goal**: Backlog management, sprint lifecycle, and coordination.

**Duration**: 3-4 weeks

### P2.1 - Backlog Manager
- [ ] Create `BacklogManager`
  - [ ] In-memory queue with persistence
  - [ ] SQLite backend for durability
  - [ ] Task state machine
- [ ] Implement task lifecycle
  - [ ] Create task from user input
  - [ ] Update task state
  - [ ] Query tasks by state
  - [ ] Delete/completed task archival
- [ ] Add task dependencies
  - [ ] DAG validation
  - [ ] Topological ordering
  - [ ] Dependency blocking
- [ ] Write backlog manager tests

### P2.2 - Sprint Manager
- [ ] Create `SprintManager`
  - [ ] Sprint initialization
  - [ ] Phase transitions
  - [ ] Iteration tracking
  - [ ] Timeout handling
- [ ] Implement sprint lifecycle
  - [ ] Start sprint
  - [ ] Run spec phase
  - [ ] Run test phase
  - [ ] Run implementation phase
  - [ ] Complete/fail sprint
- [ ] Add iteration budgets
  - [ ] Configurable max iterations per phase
  - [ ] Iteration counter
  - [ ] Budget exhaustion handling
- [ ] Write sprint manager tests

### P2.3 - Scrum Master Agent
- [ ] Create Scrum Master agent implementation
  - [ ] `Master` struct implementing `Agent` interface
  - [ ] Sprint planning logic
  - [ ] Progress tracking
  - [ ] Report generation
- [ ] Implement "daily standup" (compressed)
  - [ ] Review completed work
  - [ ] Identify next steps
  - [ ] Detect blockers
- [ ] Create progress reports
  - [ ] JSON format for machine consumption
  - [ ] Human-readable format for CLI
- [ ] Write Scrum Master tests

---

## Phase 3: Worker Agents

**Goal**: All LLM agent roles implemented and functional.

**Duration**: 4-5 weeks

### P3.1 - Base Worker Framework
- [ ] Create base agent implementation
  - [ ] `BaseWorker` struct
  - [ ] Common functionality (LLM calls, logging, etc.)
  - [ ] Error handling patterns
  - [ ] Timeout handling
- [ ] Define agent interface
  ```go
  type Agent interface {
      Name() string
      Execute(ctx context.Context, input interface{}) (interface{}, error)
  }
  ```
- [ ] Create prompt management
  - [ ] Load prompts from files
  - [ ] Template rendering
  - [ ] Prompt versioning
- [ ] Write base worker tests

### P3.2 - Product Owner Agent
- [ ] Implement Product Owner
  - [ ] Task analysis logic
  - [ ] Initial spec generation
  - [ ] Risk assessment
  - [ ] Approach recommendation
- [ ] Create prompts
  - [ ] System prompt for Product Owner
  - [ ] Task analysis prompt template
- [ ] Write Product Owner tests

### P3.3 - Spec Engineer Agent
- [ ] Implement Spec Engineer
  - [ ] Spec tightening logic
  - [ ] Ambiguity detection
  - [ ] Edge case identification
  - [ ] Confidence scoring
- [ ] Create prompts
  - [ ] Spec tightening system prompt
  - [ ] Review feedback prompt
- [ ] Implement spec diff/merge
  - [ ] Track changes between versions
  - [ ] Generate changelogs
- [ ] Write Spec Engineer tests

### P3.4 - Test Engineer Agent
- [ ] Implement Test Engineer
  - [ ] Test generation from spec
  - [ ] Edge case coverage
  - [ ] Test file organization
  - [ ] Coverage reporting
- [ ] Create prompts
  - [ ] Test generation system prompt
  - [ ] Coverage improvement prompt
- [ ] Add test validation
  - [ ] Run tests against stubs
  - [ ] Verify tests actually test something
- [ ] Write Test Engineer tests

### P3.5 - Code Engineer Agent
- [ ] Implement Code Engineer
  - [ ] Red phase implementation (stubs)
  - [ ] Green phase implementation
  - [ ] Refactor phase
  - [ ] Debug from test failures
- [ ] Create prompts
  - [ ] Code generation system prompt
  - [ ] Debug/fix system prompt
  - [ ] Language-specific prompts
- [ ] Add code validation
  - [ ] Syntax checking
  - [ ] Linting integration
- [ ] Write Code Engineer tests

### P3.6 - Review Agent
- [ ] Implement Review Agent
  - [ ] Spec review logic
  - [ ] Test review logic
  - [ ] Code review logic
  - [ ] Pass/fail decision
- [ ] Create prompts
  - [ ] Spec review system prompt
  - [ ] Test review system prompt
  - [ ] Code review system prompt
- [ ] Implement review criteria
  - [ ] Configurable checklists
  - [ ] Quality scoring
- [ ] Write Review Agent tests

---

## Phase 4: Reporting & UI

**Goal**: Progress visualization and user interaction.

**Duration**: 2-3 weeks

### P4.1 - CLI Commands
- [ ] Implement `task` command
  - [ ] Accept task description
  - [ ] Submit to backlog
  - [ ] Show task ID
- [ ] Implement `status` command
  - [ ] Show current sprint status
  - [ ] Show task progress
  - [ ] Show agent activity
- [ ] Implement `logs` command
  - [ ] View audit logs
  - [ ] Filter by task/agent
  - [ ] Follow mode (tail)
- [ ] Implement `dashboard` command
  - [ ] Launch web dashboard
  - [ ] Port configuration
- [ ] Implement `config` command
  - [ ] View current config
  - [ ] Validate config
  - [ ] Generate example config
- [ ] Write CLI tests

### P4.2 - Kanban Dashboard
- [ ] Create web server (Echo framework)
  - [ ] Static file serving
  - [ ] API routes
  - [ ] WebSocket endpoint
- [ ] Implement Kanban board UI
  - [ ] HTML/CSS/JS (no frameworks)
  - [ ] Columns: Backlog, Spec, Testing, Coding, Review, Done
  - [ ] Task cards
  - [ ] Drag and drop (optional)
- [ ] Add real-time updates
  - [ ] WebSocket connection
  - [ ] Event streaming
  - [ ] Auto-refresh
- [ ] Create markdown fallback
  - [ ] Generate kanban as markdown table
  - [ ] CLI display
- [ ] Write dashboard tests

### P4.3 - Audit & Developer Logs
- [ ] Implement audit log viewer
  - [ ] CLI command to view
  - [ ] Filter options
  - [ ] Export to JSON/CSV
- [ ] Create agent activity reports
  - [ ] Files read/written
  - [ ] LLM calls made
  - [ ] Commands executed
  - [ ] Cost tracking
- [ ] Add log analysis tools
  - [ ] Sprint post-mortem
  - [ ] Performance metrics
  - [ ] Cost analysis
- [ ] Write log viewer tests

---

## Phase 5: Integration & Polish

**Goal**: End-to-end functionality, testing, and refinement.

**Duration**: 3-4 weeks

### P5.1 - End-to-End Testing
- [ ] Create integration test suite
  - [ ] Full task lifecycle tests
  - [ ] Multiple language support tests
  - [ ] Error handling tests
- [ ] Set up test fixtures
  - [ ] Sample tasks
  - [ ] Expected outputs
  - [ ] Mock LLM responses
- [ ] Create benchmark suite
  - [ ] Performance benchmarks
  - [ ] Cost benchmarks
  - [ ] Success rate tracking
- [ ] Write integration tests

### P5.2 - Documentation
- [ ] Write API documentation
  - [ ] REST API reference
  - [ ] WebSocket protocol
  - [ ] Example requests/responses
- [ ] Create user guide
  - [ ] Installation instructions
  - [ ] Configuration guide
  - [ ] Usage examples
  - [ ] Troubleshooting
- [ ] Write architecture documentation
  - [ ] System design
  - [ ] Data flow diagrams
  - [ ] Decision records
- [ ] Create example projects
  - [ ] Python CLI tool
  - [ ] Node.js API
  - [ ] Go library

### P5.3 - Performance Optimization
- [ ] Profile and optimize
  - [ ] CPU profiling
  - [ ] Memory profiling
  - [ ] Identify bottlenecks
- [ ] Optimize container builds
  - [ ] Layer caching
  - [ ] Parallel builds
  - [ ] Image optimization
- [ ] Optimize LLM calls
  - [ ] Request batching
  - [ ] Response caching
  - [ ] Cost optimization
- [ ] Add performance tests

### P5.4 - Security Hardening
- [ ] Implement security measures
  - [ ] Container sandboxing
  - [ ] Command whitelisting
  - [ ] Network isolation
  - [ ] Secret management
- [ ] Security audit
  - [ ] Dependency scanning
  - [ ] Container scanning
  - [ ] Code review
- [ ] Write security documentation

---

## Phase 6: Advanced Features (Post-MVP)

**Goal**: Enhanced capabilities for production use.

**Duration**: Ongoing / 4-6 weeks

### P6.1 - Multi-Provider Support
- [ ] Add more LLM providers
  - [ ] Anthropic direct
  - [ ] Google (Gemini)
  - [ ] Cohere
  - [ ] Local model servers
- [ ] Implement provider fallback
  - [ ] Automatic failover
  - [ ] Load balancing
  - [ ] Cost optimization
- [ ] Add provider-specific optimizations

### P6.2 - Advanced Container Features
- [ ] Add caching layer
  - [ ] Shared package cache
  - [ ] Incremental builds
  - [ ] Layer optimization
- [ ] Support for complex builds
  - [ ] Multi-stage builds
  - [ ] Build arguments
  - [ ] Environment variables
- [ ] Add container health checks

### P6.3 - Skill System
- [ ] Create skill registry
  - [ ] Skill definition format
  - [ ] Skill loading
  - [ ] Skill versioning
- [ ] Implement skill execution
  - [ ] Skill contexts
  - [ ] Skill composition
  - [ ] Skill testing
- [ ] Add built-in skills
  - [ ] Common code patterns
  - [ ] Framework-specific helpers
  - [ ] Testing utilities

### P6.4 - Git Integration
- [ ] Add GitHub integration
  - [ ] PR creation
  - [ ] Issue tracking
  - [ ] CI trigger
- [ ] Add GitLab integration
- [ ] Implement branch management
  - [ ] Feature branches
  - [ ] Commit on completion
  - [ ] PR workflow

---

## Backlog (Future Ideas)

### Features
- [~] Natural language queries about progress ("What's blocking task X?")
- [~] Persistent knowledge base across tasks
- [~] Automatic spec generation from existing codebase
- [~] Property-based testing integration
- [~] Mutation testing
- [~] Formal verification hooks
- [~] IDE plugins (VSCode, JetBrains)
- [~] Multi-host distributed execution
- [~] Model fine-tuning on task history

### Research
- [~] Automatic agent role discovery
- [~] Dynamic iteration budgeting
- [~] Success prediction
- [~] Optimal model selection per task

---

## Milestones

### Milestone 1: Foundation Complete
**Criteria**:
- [ ] Project builds successfully
- [ ] All P0 tasks complete
- [ ] Basic abstractions in place
- [ ] CI/CD passing

**Target**: Week 3

### Milestone 2: Infrastructure Ready
**Criteria**:
- [ ] Containers working
- [ ] LLM routing working
- [ ] Workspaces functional
- [ ] Event bus operational

**Target**: Week 7

### Milestone 3: Core Scrum Working
**Criteria**:
- [ ] Backlog manager working
- [ ] Sprint manager working
- [ ] Basic task lifecycle functional
- [ ] Can run simple end-to-end

**Target**: Week 11

### Milestone 4: MVP Complete
**Criteria**:
- [ ] All agents implemented
- [ ] CLI functional
- [ ] Dashboard working
- [ ] Can handle real tasks
- [ ] Documentation complete
- [ ] Tests passing

**Target**: Week 16

### Milestone 5: Production Ready
**Criteria**:
- [ ] Performance optimized
- [ ] Security hardened
- [ ] Monitoring in place
- [ ] Stable for daily use

**Target**: Week 20

---

## Daily Standup Template (for human team)

**What did I complete?**
- Task: X
- Progress: Y%

**What am I working on?**
- Current: X
- Blockers: None / [description]

**What do I need?**
- Help with: X
- Decisions needed: Y

---

*Last updated: 2024-03-26*
*Next review: Weekly on Mondays*
