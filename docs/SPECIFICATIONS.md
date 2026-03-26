# Scru-LLM System Specifications

This document describes the target-state system design. It is not evidence that the current repository already implements the described behavior.

## 1. Overview

### 1.1 Purpose
Scru-LLM is an autonomous software engineering system that applies Scrum methodology to LLM-driven development. It compresses traditional Scrum cycles from days/weeks to task-completion-driven events while maintaining rigorous quality gates.

### 1.2 Goals
- **Zero-to-working-code** for well-specified tasks with minimal human intervention
- **Spec-first development** with iterative tightening before implementation
- **Test-driven guarantees** through comprehensive unit test generation
- **Reproducible builds** via containerized execution environments
- **Complete observability** through detailed audit trails and progress reporting

### 1.3 Non-Goals
- General-purpose chat interface
- Real-time collaborative editing
- IDE integration (initially)
- Multi-host distributed execution (Phase 2)

---

## 2. System Architecture

### 2.1 High-Level Components

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Interface Layer                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  CLI Tool    │  │  Web Dashboard│  │  API Server  │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                     Scrum Orchestration Layer                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │Backlog Manager│  │Sprint Manager│  │ Event Bus    │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                      LLM Agent Roles Layer                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐ │
│  │Prod Owner   │ │Scrum Master │ │Spec Engineer│ │Test Eng   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │Code Engineer│ │Review Agent │ │Container Mgr│               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                      Execution Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Container    │  │ LLM Client   │  │ File System  │          │
│  │ Runtime      │  │ Router       │  │ Sandbox      │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Data Flow

```
User Input
    ↓
Product Owner analyzes → creates Task Specification (v0)
    ↓
Scrum Master creates Sprint Plan → assigns to Backlog
    ↓
Spec Sprint: Spec Engineer ↔ Review Agent (iterative)
    ↓
Test Sprint: Test Engineer ↔ Code Engineer (red phase)
    ↓
Implementation Sprint: Code Engineer ↔ Review Agent (green phase)
    ↓
Container Manager validates in isolated environment
    ↓
Scrum Master reports completion → User receives deliverable
```

---

## 3. Component Specifications

### 3.1 User Interface Layer

#### 3.1.1 CLI Tool
**Purpose**: Primary interface for task submission and system control

**Commands**:
- `scru-llm task "<description>"` - Submit new task
- `scru-llm status` - View current sprint status
- `scru-llm dashboard` - Launch web dashboard
- `scru-llm logs [task-id]` - View detailed logs
- `scru-llm config` - Manage configuration

**Requirements**:
- Rich terminal output with progress indicators
- Support for piping specs from files
- Interactive mode for complex tasks

#### 3.1.2 Web Dashboard
**Purpose**: Visual progress tracking via Kanban-style board

**Features**:
- Task backlog view (To Do, In Progress, Done)
- Sprint progress visualization
- Real-time updates via WebSocket
- Simple, minimal UI (markdown tables as fallback)

**Requirements**:
- Single-page application
- Responsive design
- No external dependencies (pure HTML/CSS/JS)

#### 3.1.3 API Server
**Purpose**: Programmatic access for integrations

**Endpoints**:
- `POST /api/v1/tasks` - Create task
- `GET /api/v1/tasks/{id}` - Get task status
- `GET /api/v1/sprints/current` - Current sprint info
- `WS /api/v1/events` - Real-time event stream

### 3.2 Scrum Orchestration Layer

#### 3.2.1 Backlog Manager
**Purpose**: Manage task queue and prioritization

**Responsibilities**:
- Queue incoming tasks
- Track task state transitions
- Persist backlog to disk (JSON/SQLite)
- Support task dependencies

**States**:
- `pending` - Awaiting sprint assignment
- `spec_sprint` - In spec iteration
- `test_sprint` - Generating tests
- `implementation_sprint` - Coding
- `review` - Final review
- `completed` - Done
- `failed` - Terminal failure

#### 3.2.2 Sprint Manager
**Purpose**: Coordinate sprint execution and iteration loops

**Responsibilities**:
- Initialize sprints with appropriate agents
- Manage iteration budgets (max N attempts per phase)
- Detect and handle failures/stalls
- Coordinate agent handoffs
- Track sprint metrics

**Sprint Configuration**:
```yaml
sprint:
  max_spec_iterations: 3
  max_test_iterations: 3
  max_implementation_iterations: 5
  iteration_timeout: 10m
  total_timeout: 1h
```

#### 3.2.3 Event Bus
**Purpose**: Decoupled communication between components

**Requirements**:
- Async pub/sub pattern
- Persistent event log
- Support for event replay
- Structured JSON events

**Event Types**:
- `task.created`
- `sprint.started`
- `agent.assigned`
- `agent.completed`
- `iteration.completed`
- `sprint.completed`
- `sprint.failed`

### 3.3 LLM Agent Roles Layer

#### 3.3.1 Product Owner Agent
**Purpose**: Translate user tasks into initial specifications

**Model Requirements**:
- Strong reasoning capabilities
- Good at requirements analysis
- 14B+ parameter model recommended

**Inputs**:
- User task description
- Context (previous tasks, codebase if applicable)

**Outputs**:
- Task Specification v0 (JSON structure)
- Risk assessment
- Suggested approach

**Specification Schema**:
```json
{
  "task_id": "uuid",
  "title": "string",
  "description": "string",
  "requirements": ["string"],
  "acceptance_criteria": ["string"],
  "constraints": ["string"],
  "dependencies": ["task_id"],
  "estimated_complexity": "low|medium|high",
  "deliverables": ["string"],
  "context": {
    "existing_codebase": "boolean",
    "language": "string",
    "framework": "string"
  }
}
```

#### 3.3.2 Scrum Master Agent
**Purpose**: Orchestrate sprint execution and report progress

**Model Requirements**:
- Good at coordination and planning
- Can parse structured outputs
- 7B+ parameter model sufficient

**Responsibilities**:
- Create sprint plans from specs
- Assign tasks to appropriate agents
- Monitor sprint health
- Generate progress reports
- Escalate blockers

**Daily Standup (Compressed)**:
- Happens after each iteration completion
- Reviews: What was done, what's next, any blockers
- Updates Kanban board automatically

**Progress Report Format**:
```json
{
  "sprint_id": "uuid",
  "status": "active|completed|failed",
  "current_phase": "spec|test|implementation|review",
  "iteration": 2,
  "max_iterations": 5,
  "agents_active": ["agent_id"],
  "tasks_completed": 3,
  "tasks_total": 5,
  "blockers": [],
  "eta_remaining": "15m"
}
```

#### 3.3.3 Spec Engineer Agent
**Purpose**: Tighten specifications through iterative refinement

**Model Requirements**:
- Excellent technical writing
- Good at identifying ambiguities
- 7B-14B model range

**Inputs**:
- Specification vN
- Review feedback

**Outputs**:
- Specification v(N+1)
- Changelog explaining changes
- Confidence score (0-1)

**Tightening Criteria**:
- All requirements are measurable
- Acceptance criteria are testable
- No undefined terms
- Edge cases documented
- Interfaces/contracts specified

#### 3.3.4 Test Engineer Agent
**Purpose**: Generate comprehensive unit tests from specifications

**Model Requirements**:
- Strong understanding of testing patterns
- Good at edge case identification
- 14B+ parameter model recommended

**Inputs**:
- Finalized specification
- Code structure (if partial implementation exists)

**Outputs**:
- Test suite (multiple files)
- Test coverage report
- Edge case documentation

**Test Requirements**:
- 80%+ line coverage minimum
- Unit tests for all public functions
- Edge case tests
- Error condition tests
- Integration tests (if applicable)

**Output Structure**:
```
tests/
├── unit/
│   ├── test_feature_a.py
│   └── test_feature_b.py
├── integration/
│   └── test_integration.py
├── conftest.py
└── coverage_report.json
```

#### 3.3.5 Code Engineer Agent
**Purpose**: Implement code to pass tests

**Model Requirements**:
- Strong coding capabilities
- Good at debugging
- 14B+ parameter model recommended
- Multiple model sizes for cost optimization

**Modes**:
1. **Red Phase**: Write minimal failing code to validate tests
2. **Green Phase**: Implement to pass all tests
3. **Refactor Phase**: Clean up while maintaining passing tests

**Inputs**:
- Specification
- Test suite
- Previous implementation (if iterating)
- Error output (if test failures)

**Outputs**:
- Source code
- Implementation notes
- Confidence score

**Constraints**:
- Must pass all tests
- Must follow coding guidelines
- Must include docstrings/comments

#### 3.3.6 Review Agent
**Purpose**: Validate outputs against specifications and standards

**Model Requirements**:
- Good at code review
- Understands quality standards
- 7B-14B model range

**Inputs**:
- Specification
- Output to review (spec, tests, or code)
- Previous review feedback (if iterating)

**Outputs**:
- Pass/Fail decision
- Detailed feedback
- Specific issues with line references
- Suggested fixes

**Review Dimensions**:
- **Spec Review**: Completeness, clarity, testability
- **Test Review**: Coverage, correctness, edge cases
- **Code Review**: Correctness, style, efficiency, security

#### 3.3.7 Container Manager Agent
**Purpose**: Manage reproducible build and test environments

**Not an LLM agent** - deterministic automation component

**Responsibilities**:
- Generate Dockerfiles/containerspecs
- Build container images
- Run tests in isolated environment
- Capture build/test outputs
- Manage container lifecycle

**Container Requirements**:
- Reproducible from spec
- Isolated from host system
- Pre-configured with common toolchains
- Support for multiple languages
- Ephemeral (created per task, destroyed after)

**Supported Runtimes**:
- Docker (default)
- Podman
- Containerd

### 3.4 Execution Layer

#### 3.4.1 Container Runtime
**Purpose**: Execute code in isolated, reproducible environments

**Features**:
- Dockerfile generation from project requirements
- Image building and caching
- Container execution with resource limits
- Volume mounting for source code
- Log capture

**Configuration**:
```yaml
container:
  runtime: docker
  base_images:
    python: python:3.11-slim
    node: node:18-alpine
    go: golang:1.21-alpine
    rust: rust:1.75-slim
  resource_limits:
    cpu: 2
    memory: 4g
    timeout: 10m
```

#### 3.4.2 LLM Client Router
**Purpose**: Route requests to appropriate LLM providers

**Features**:
- Multi-provider support (OpenAI, OpenRouter, local)
- Role-to-model mapping
- Cost tracking and budgeting
- Fallback chains
- Request retry with backoff

**Configuration**:
```yaml
llm:
  providers:
    - name: openrouter
      base_url: https://openrouter.ai/api/v1
      api_key: ${OPENROUTER_API_KEY}
      models:
        large: anthropic/claude-3.5-sonnet
        medium: meta-llama/llama-3.1-70b
        small: google/gemma-2-9b-it
  
  roles:
    product_owner: large
    scrum_master: medium
    spec_engineer: medium
    test_engineer: large
    code_engineer: large
    reviewer: medium
```

#### 3.4.3 File System Sandbox
**Purpose**: Isolated workspace for each task

**Features**:
- Task-scoped directories
- Atomic operations
- Git integration for versioning
- Cleanup on completion

**Structure**:
```
workspaces/
├── {task_id}/
│   ├── spec/           # Specification versions
│   ├── tests/          # Test files
│   ├── src/            # Source code
│   ├── containers/     # Container definitions
│   ├── logs/           # Execution logs
│   └── artifacts/      # Build artifacts
```

---

## 4. Sprint Lifecycle Specification

### 4.1 Phase 1: Task Intake

**Duration**: Single pass
**Goal**: Create initial task specification

**Process**:
1. User submits task via CLI/API
2. System assigns task_id and creates workspace
3. Product Owner Agent analyzes task
4. Product Owner creates Specification v0
5. Scrum Master reviews spec completeness
6. If acceptable, move to Spec Sprint
7. If unclear, request clarification from user

**Exit Criteria**:
- Specification v0 exists
- Basic requirements documented
- Context identified

### 4.2 Phase 2: Spec Sprint

**Duration**: Up to N iterations (configurable, default 3)
**Goal**: Tighten specification until it's "ready"

**Process per Iteration**:
1. Spec Engineer reviews current specification
2. Spec Engineer identifies ambiguities/gaps
3. Spec Engineer produces tightened specification v(N+1)
4. Review Agent validates specification
5. If Review Agent approves → exit to Test Sprint
6. If Review Agent rejects → continue to next iteration
7. If max iterations reached → fail with feedback

**Ready Criteria** (all must be true):
- All requirements are measurable
- Acceptance criteria are testable
- No undefined terms
- Edge cases documented
- Confidence score > 0.8

**Iteration Budget**:
```yaml
spec_sprint:
  max_iterations: 3
  timeout_per_iteration: 5m
  temperature: 0.3  # Lower for consistency
```

### 4.3 Phase 3: Test Sprint

**Duration**: Up to N iterations (configurable, default 3)
**Goal**: Generate comprehensive, passing test suite

**Process per Iteration**:
1. Test Engineer reads finalized specification
2. Test Engineer generates test suite
3. Code Engineer writes minimal "red" code (stubs)
4. Container Manager runs tests (expecting failures)
5. Review Agent reviews test quality
6. If tests are comprehensive → exit to Implementation Sprint
7. If tests need work → continue to next iteration
8. If max iterations reached → fail with feedback

**Comprehensive Test Criteria**:
- 80%+ line coverage
- All public functions have tests
- Edge cases covered
- Error conditions tested
- Tests actually fail against stubs (confirm they're real tests)

**Iteration Budget**:
```yaml
test_sprint:
  max_iterations: 3
  timeout_per_iteration: 10m
  temperature: 0.5
```

### 4.4 Phase 4: Implementation Sprint

**Duration**: Up to N iterations (configurable, default 5)
**Goal**: Implement code that passes all tests

**Process per Iteration**:
1. Code Engineer reads specification and test suite
2. Code Engineer implements code
3. Container Manager builds and runs tests
4. If all tests pass → exit to Review
5. If tests fail:
   - Code Engineer receives test output
   - Code Engineer debugs and fixes
   - Continue to next iteration
6. If max iterations reached → fail with partial implementation

**Pass Criteria**:
- All unit tests pass
- No regressions in existing tests
- Code follows style guidelines
- Review Agent approves quality

**Iteration Budget**:
```yaml
implementation_sprint:
  max_iterations: 5
  timeout_per_iteration: 15m
  temperature: 0.7  # Higher for creativity
```

### 4.5 Phase 5: Final Review and Delivery

**Duration**: Single pass
**Goal**: Validate deliverable quality

**Process**:
1. Review Agent performs final code review
2. Container Manager runs full test suite in clean environment
3. Documentation generated
4. Artifacts packaged
5. Scrum Master generates completion report
6. Deliverable presented to user

**Exit Criteria**:
- All tests pass in clean container
- Code review passed
- Documentation complete
- Artifacts available

---

## 5. Reporting and Observability

### 5.1 User-Facing Reports

#### 5.1.1 Kanban Dashboard
**Format**: Web UI with markdown fallback

**Columns**:
- **Backlog**: Tasks awaiting assignment
- **Spec**: Tasks in specification phase
- **Testing**: Tasks generating tests
- **Coding**: Tasks in implementation
- **Review**: Tasks awaiting final review
- **Done**: Completed tasks

**Features**:
- Drag-and-drop (web UI)
- Click for task details
- Real-time updates
- Simple color coding for status

#### 5.1.2 Sprint Progress Reports
**Format**: CLI output and web UI

**Content**:
- Current phase
- Iteration count (X of Y)
- Time elapsed/remaining
- Active agents
- Recent events

**CLI Example**:
```
Sprint: task-123e4567-e89b-12d3-a456-426614174000
Phase: Implementation (Iteration 2/5)
Elapsed: 23m 15s | ETA: ~30m

Active Agents:
  • Code Engineer (claude-3.5-sonnet)

Recent Events:
  [2m ago] Code Engineer submitted implementation
  [2m ago] Tests failed: 3 errors in test_parser.py
  [1m ago] Code Engineer started debugging

Blockers: None
```

### 5.2 Developer-Facing Reports

#### 5.2.1 Audit Log
**Format**: JSONL (one JSON object per line)

**Fields**:
- `timestamp` - ISO 8601
- `event_type` - Category of event
- `task_id` - Associated task
- `agent_id` - Which agent (if applicable)
- `action` - What happened
- `inputs` - Input context (truncated)
- `outputs` - Output produced (truncated)
- `metadata` - Additional context

**Example**:
```json
{
  "timestamp": "2024-03-26T14:23:45Z",
  "event_type": "agent.completion",
  "task_id": "task-123e4567",
  "agent_id": "spec-engineer-001",
  "action": "spec_tightening",
  "inputs": {"spec_v": 1, "feedback": "ambiguous requirements"},
  "outputs": {"spec_v": 2, "changes": ["clarified input format"]},
  "metadata": {"tokens_in": 2048, "tokens_out": 1536, "cost": 0.004}
}
```

#### 5.2.2 Agent Activity Log
**Format**: Structured logs per agent

**Content**:
- Files read by agent
- Files written by agent
- Commands executed
- LLM calls made
- Duration of operations
- Success/failure status

#### 5.2.3 Execution Traces
**Format**: Detailed step-by-step logs

**Content**:
- Every tool invocation
- Full LLM prompts (optional)
- LLM responses
- Error details
- Stack traces

**Storage**:
- Local: `~/.scru-llm/logs/`
- Rotation: Keep last 30 days
- Compression: gzip after 7 days

### 5.3 Metrics Collection

**Tracked Metrics**:
- Sprint success rate
- Average time per phase
- Iteration counts per phase
- Cost per task
- Test coverage achieved
- Code quality scores
- Agent utilization

**Export**:
- Prometheus metrics endpoint
- CSV export for analysis
- Real-time dashboard updates

---

## 6. Configuration Specification

### 6.1 Configuration File
**Path**: `~/.scru-llm/scru-llm.yaml`
**Format**: YAML

### 6.2 Configuration Schema

```yaml
# System settings
system:
  workspace_dir: "~/.scru-llm/workspaces"
  log_level: "info"  # debug, info, warn, error
  max_concurrent_tasks: 3
  
# LLM provider configuration
llm:
  providers:
    - name: openrouter
      base_url: "https://openrouter.ai/api/v1"
      api_key: "${OPENROUTER_API_KEY}"
      default_model: "anthropic/claude-3.5-sonnet"
      
    - name: local
      base_url: "http://localhost:1234/v1"
      api_key: "sk-local"
      default_model: "local-model"
  
  # Role-to-model mapping
  roles:
    product_owner: 
      provider: openrouter
      model: anthropic/claude-3.5-sonnet
      temperature: 0.5
      max_tokens: 4096
      
    scrum_master:
      provider: openrouter
      model: meta-llama/llama-3.1-70b
      temperature: 0.4
      max_tokens: 2048
      
    spec_engineer:
      provider: openrouter
      model: meta-llama/llama-3.1-70b
      temperature: 0.3
      max_tokens: 4096
      
    test_engineer:
      provider: openrouter
      model: anthropic/claude-3.5-sonnet
      temperature: 0.5
      max_tokens: 8192
      
    code_engineer:
      provider: openrouter
      model: anthropic/claude-3.5-sonnet
      temperature: 0.7
      max_tokens: 8192
      
    reviewer:
      provider: openrouter
      model: meta-llama/llama-3.1-70b
      temperature: 0.3
      max_tokens: 4096

# Sprint configuration
sprint:
  spec_sprint:
    max_iterations: 3
    timeout_per_iteration: 5m
    
  test_sprint:
    max_iterations: 3
    timeout_per_iteration: 10m
    min_coverage_percent: 80
    
  implementation_sprint:
    max_iterations: 5
    timeout_per_iteration: 15m
    
  global_timeout: 1h

# Container configuration
container:
  runtime: docker  # docker, podman, containerd
  registry: ""  # Leave empty for local builds
  
  base_images:
    python: "python:3.11-slim"
    node: "node:18-alpine"
    go: "golang:1.21-alpine"
    rust: "rust:1.75-slim"
    java: "openjdk:21-slim"
    
  resource_limits:
    cpu: 2
    memory: "4g"
    timeout: 10m

# Reporting configuration
reporting:
  dashboard:
    enabled: true
    port: 8080
    refresh_interval: 5s
    
  audit_log:
    enabled: true
    path: "~/.scru-llm/logs/audit.jsonl"
    rotation_days: 30
    
  agent_logs:
    enabled: true
    path: "~/.scru-llm/logs/agents/"
    retention_days: 14

# Cost controls
cost:
  budget_per_task: 5.00  # USD
  warning_threshold: 0.8  # Alert at 80% of budget
  emergency_stop_threshold: 1.0  # Hard stop at 100%
  
# Security
security:
  sandbox_enabled: true
  allow_network_access: false  # Containers isolated by default
  allowed_commands: []  # Whitelist (empty = deny all except build/test)
```

---

## 7. API Specification

### 7.1 REST API

#### Task Management

**POST /api/v1/tasks**
Create a new task

Request:
```json
{
  "title": "Convert JSON to YAML CLI tool",
  "description": "Create a Python CLI tool that reads JSON from stdin and outputs YAML",
  "context": {
    "language": "python",
    "existing_codebase": false
  },
  "constraints": {
    "max_cost": 2.00,
    "timeout": "30m"
  }
}
```

Response:
```json
{
  "task_id": "task-uuid",
  "status": "pending",
  "created_at": "2024-03-26T14:23:45Z",
  "estimated_duration": "25m"
}
```

**GET /api/v1/tasks/{task_id}**
Get task status

Response:
```json
{
  "task_id": "task-uuid",
  "status": "implementation_sprint",
  "current_phase": "implementation",
  "iteration": 2,
  "max_iterations": 5,
  "progress": {
    "spec_sprint": "completed",
    "test_sprint": "completed",
    "implementation_sprint": "active"
  },
  "agents_active": ["code-engineer-001"],
  "created_at": "2024-03-26T14:23:45Z",
  "started_at": "2024-03-26T14:24:00Z",
  "elapsed": "12m 30s",
  "artifacts": {
    "spec": "/workspaces/task-uuid/spec/spec_v2.json",
    "tests": "/workspaces/task-uuid/tests/",
    "code": "/workspaces/task-uuid/src/"
  }
}
```

**GET /api/v1/tasks/{task_id}/logs**
Get task logs

Query params:
- `level`: debug, info, warn, error
- `agent`: Filter by agent ID
- `since`: ISO timestamp

Response:
```json
{
  "task_id": "task-uuid",
  "logs": [
    {
      "timestamp": "2024-03-26T14:25:00Z",
      "level": "info",
      "agent": "spec-engineer-001",
      "message": "Spec tightened: clarified input format"
    }
  ]
}
```

#### Sprint Management

**GET /api/v1/sprints/current**
Get current sprint status

**GET /api/v1/sprints/{sprint_id}**
Get specific sprint details

**POST /api/v1/sprints/{sprint_id}/cancel**
Cancel a running sprint

### 7.2 WebSocket API

**WS /api/v1/events**
Real-time event stream

Events:
```json
{
  "type": "sprint.updated",
  "timestamp": "2024-03-26T14:30:00Z",
  "data": {
    "sprint_id": "sprint-uuid",
    "status": "active",
    "iteration": 3,
    "agent_status": {
      "code-engineer-001": "working"
    }
  }
}
```

---

## 8. Quality Gates

### 8.1 Specification Quality Gate
- All requirements must be measurable
- Acceptance criteria must be testable
- No undefined terms
- Confidence score > 0.8

### 8.2 Test Quality Gate
- Minimum 80% line coverage
- All public functions tested
- Edge cases documented
- Tests fail against stubs (validated)

### 8.3 Implementation Quality Gate
- All tests pass
- Code review passed
- No linting errors
- Documentation complete
- Container build succeeds

### 8.4 Final Quality Gate
- All previous gates passed
- Full test suite passes in clean container
- No security vulnerabilities (basic scan)
- Performance within acceptable bounds (if specified)

---

## 9. Error Handling

### 9.1 Failure Modes

**Spec Sprint Failure**:
- Max iterations reached
- Confidence score never reached threshold
- Review Agent consistently rejects

**Test Sprint Failure**:
- Max iterations reached
- Cannot achieve minimum coverage
- Tests don't actually test anything (fake tests)

**Implementation Sprint Failure**:
- Max iterations reached
- Tests consistently fail
- Code quality gates not met

**System Failure**:
- LLM provider unavailable
- Container runtime failure
- Disk space exhausted
- Memory exhaustion

### 9.2 Recovery Strategies

**Iteration Exhaustion**:
- Escalate to larger model
- Reduce scope (if possible)
- Request user intervention with specific feedback

**Provider Failure**:
- Failover to backup provider
- Queue for retry
- Notify user

**Container Failure**:
- Retry with fresh container
- Log detailed diagnostics
- Escalate to human if persistent

---

## 10. Security Considerations

### 10.1 Container Isolation
- No host network access by default
- Read-only filesystem except workspace
- Resource limits (CPU, memory, disk)
- No privileged containers

### 10.2 Code Execution Safety
- All code runs in containers
- Command whitelist (build, test, package management)
- No arbitrary shell execution
- Timeout enforced

### 10.3 Secret Management
- API keys via environment variables only
- No secrets in logs
- Secret rotation support
- Audit access to secrets

### 10.4 Audit Trail
- All actions logged
- Immutable logs
- Tamper-evident (signed)
- Retention policy

---

## 11. Performance Requirements

### 11.1 Latency Targets
- Task submission to sprint start: <5s
- Spec iteration: <5m
- Test generation: <10m
- Implementation iteration: <15m
- Dashboard refresh: <2s

### 11.2 Throughput Targets
- Concurrent tasks: 3+ (configurable)
- Events per second: 100+
- Log ingestion: 1000 lines/second

### 11.3 Resource Usage
- Memory: <2GB base, +1GB per concurrent task
- Disk: <10GB for logs/workspaces (configurable retention)
- CPU: Scales with concurrent tasks

---

## 12. Future Considerations (Phase 2+)

### 12.1 Advanced Features
- Multi-host distributed execution
- Persistent knowledge base across tasks
- Skill learning and reuse
- Natural language progress queries
- GitHub/GitLab integration
- IDE plugins

### 12.2 Research Areas
- Automatic spec generation from codebase
- Test oracle generation
- Mutation testing
- Property-based testing
- Formal verification integration

### 12.3 Scaling Considerations
- Horizontal scaling of agents
- Distributed container orchestration
- Shared workspace caching
- Model fine-tuning on task history

---

## Appendix A: Glossary

**Sprint**: A time-boxed iteration focused on a specific deliverable
**Specification**: Detailed description of requirements and acceptance criteria
**Tightening**: Process of making specifications more precise and testable
**Red Phase**: Writing tests that fail (confirm tests are real)
**Green Phase**: Implementing code to pass tests
**Container Manager**: Component managing reproducible build environments

## Appendix B: File Formats

### B.1 Task Specification JSON Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["task_id", "title", "description", "requirements", "acceptance_criteria"],
  "properties": {
    "task_id": {"type": "string", "format": "uuid"},
    "title": {"type": "string", "maxLength": 200},
    "description": {"type": "string"},
    "requirements": {"type": "array", "items": {"type": "string"}},
    "acceptance_criteria": {"type": "array", "items": {"type": "string"}},
    "constraints": {"type": "array", "items": {"type": "string"}},
    "dependencies": {"type": "array", "items": {"type": "string"}},
    "estimated_complexity": {"enum": ["low", "medium", "high"]},
    "deliverables": {"type": "array", "items": {"type": "string"}},
    "context": {
      "type": "object",
      "properties": {
        "existing_codebase": {"type": "boolean"},
        "language": {"type": "string"},
        "framework": {"type": "string"}
      }
    }
  }
}
```

### B.2 Audit Log Entry Schema
See Section 5.2.1 for field definitions.

### B.3 Configuration Schema
See Section 6.2 for complete configuration options.
