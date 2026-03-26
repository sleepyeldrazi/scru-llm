# Architecture Overview

This document provides high-level architectural context for Scru-LLM.

## Design Philosophy

Scru-LLM is built on several key principles derived from agent systems research and practical software engineering:

### 1. Bounded Autonomy
Rather than giving agents open-ended freedom, we constrain them with:
- **Iteration budgets**: Maximum attempts per phase
- **Phase gates**: Must pass quality checks to proceed
- **Timeouts**: Hard limits on execution time
- **Cost budgets**: Maximum spend per task

This prevents runaway agents and ensures predictable resource usage.

### 2. Verification-First
Every deliverable must be verified:
- **Specs**: Reviewed for completeness and testability
- **Tests**: Validated to actually test something (red phase)
- **Code**: Must pass all tests in isolated environment

Verification happens through independent Review Agent, not self-assessment.

### 3. Spec-Driven Development
Before any code is written:
1. Spec is tightened until unambiguous
2. Tests are generated from spec
3. Code is written to pass tests

This prevents "requirements drift" where implementation diverges from original intent.

### 4. Container Isolation
All work happens in containers for:
- **Reproducibility**: Same environment every time
- **Safety**: Code can't harm host system
- **Cleanliness**: Fresh start for each task
- **Testing**: Tests run in production-like environment

### 5. Complete Observability
Every action is logged:
- User sees Kanban progress
- Developers see detailed traces
- Auditors see complete history

Nothing happens in a black box.

## System Layers

### Layer 1: User Interface
**Purpose**: Human interaction with the system

**Components**:
- CLI tool for task submission
- Web dashboard for progress tracking
- REST API for programmatic access

**Design Decisions**:
- CLI is primary interface (developer-friendly)
- Dashboard is simple by design (no complex frameworks)
- API follows REST conventions

### Layer 2: Scrum Orchestration
**Purpose**: Coordinate the Scrum process

**Components**:
- Backlog Manager: Task queue and state machine
- Sprint Manager: Lifecycle coordination
- Event Bus: Decoupled communication

**Design Decisions**:
- Event-driven architecture for loose coupling
- State machine for task lifecycle (clear transitions)
- Persistent queue (survives crashes)

### Layer 3: LLM Agents
**Purpose**: Specialized AI workers

**Components**:
- Product Owner: Requirements analysis
- Scrum Master: Coordination and reporting
- Spec Engineer: Specification tightening
- Test Engineer: Test generation
- Code Engineer: Implementation
- Review Agent: Quality validation

**Design Decisions**:
- Each role has distinct responsibilities
- Roles map to different models (cost optimization)
- Review Agent is separate from workers (independent validation)

### Layer 4: Execution
**Purpose**: Run code and manage resources

**Components**:
- Container Runtime: Isolated execution
- LLM Router: Provider management
- Workspace Manager: File operations

**Design Decisions**:
- Pluggable container runtime (Docker, Podman, etc.)
- Multi-provider LLM support (failover, cost optimization)
- Git-backed workspaces (versioning and audit)

## Key Architectural Decisions

### Why Go?

**Pros**:
- Strong concurrency support (goroutines, channels)
- Static typing catches errors early
- Fast compilation and execution
- Excellent standard library
- Single binary deployment
- Good for systems programming

**Cons**:
- Verbose error handling
- Less ML ecosystem than Python
- No generics until recently (Go 1.18+)

**Alternatives Considered**:
- Python: Better ML ecosystem, but slower and less type-safe
- Rust: Safer, but steeper learning curve and slower development
- TypeScript: Good for web, less ideal for systems code

**Decision**: Go for core system, containers handle language-specific needs.

### Why Containers?

**Pros**:
- Reproducible environments
- Isolation and security
- Industry standard
- Easy to reset state

**Cons**:
- Overhead (but acceptable for our use case)
- Requires container runtime
- Build time

**Alternatives Considered**:
- VMs: Too heavy
- chroot: Not isolated enough
- Language-specific virtualenvs: Not language-agnostic

**Decision**: Docker/Podman with ephemeral containers per task.

### Why Separate Review Agent?

**Pros**:
- Independent validation (not self-assessed)
- Can catch issues author missed
- Clear accountability
- Different model can be used (cheaper)

**Cons**:
- Additional LLM calls (cost)
- Additional latency
- May disagree with worker (iterations)

**Alternatives Considered**:
- Self-review: Faster, but less reliable
- Deterministic checks: Fast, but can't assess quality
- Human review: Too slow for autonomous system

**Decision**: Separate Review Agent with deterministic checks where possible.

### Why Event-Driven?

**Pros**:
- Loose coupling between components
- Easy to add new listeners
- Audit trail built-in
- Supports async processing

**Cons**:
- Eventual consistency
- Harder to trace flow
- Risk of lost events (mitigated by persistence)

**Alternatives Considered**:
- Direct calls: Simpler, but tight coupling
- Message queue (RabbitMQ, etc.): Overkill for single-node

**Decision**: In-memory event bus with JSONL persistence.

## Data Flow Deep Dive

### Task Submission Flow

```
1. User → CLI/API
   "Create a Python CLI tool..."

2. CLI → Backlog Manager
   POST /api/v1/tasks
   
3. Backlog Manager → Event Bus
   TaskCreated event

4. Sprint Manager (listener) → picks up task
   Creates Sprint
   
5. Sprint Manager → Product Owner Agent
   "Analyze this task"

6. Product Owner → generates Specification v0

7. Sprint Manager → Event Bus
   SprintStarted, PhaseChanged(spec)

8. Dashboard (listener) → updates Kanban
   Task moves to "Spec" column
```

### Spec Sprint Flow

```
1. Sprint Manager → Spec Engineer
   "Tighten this spec"

2. Spec Engineer → LLM Router
   Generate tightened spec v1

3. Spec Engineer → Workspace
   Write spec_v1.json

4. Sprint Manager → Review Agent
   "Review this spec"

5. Review Agent → evaluates spec
   Pass/Fail decision

6. If Fail:
   - Review feedback → Spec Engineer
   - Loop to step 2
   
7. If Pass or Max Iterations:
   - Sprint Manager → Event Bus
   - PhaseChanged(test)
```

### Implementation Sprint Flow

```
1. Sprint Manager → Code Engineer
   "Implement to pass these tests"

2. Code Engineer → Workspace
   Write source code

3. Container Manager → builds container
   Dockerfile → Image

4. Container Manager → runs tests
   Execute in isolated container

5. If Tests Fail:
   - Test output → Code Engineer
   - "Debug and fix"
   - Loop to step 2
   
6. If Tests Pass:
   - Review Agent → code review
   - If Pass → Complete
   - If Fail → Loop with feedback
```

## Scalability Considerations

### Current Scope (Single Node)
- Concurrent tasks: 3-5
- Concurrent agents: 10-15
- Workspace storage: ~10GB
- Memory: ~4GB

### Scaling to Multiple Nodes (Phase 2)

**Challenges**:
- Shared state (backlog, sprint state)
- Event bus distribution
- Workspace storage
- Container runtime distribution

**Potential Solutions**:
1. **Shared State**: Redis or PostgreSQL
2. **Event Bus**: Redis pub/sub or NATS
3. **Workspaces**: Shared NFS or S3
4. **Containers**: Kubernetes or Nomad

**Decision**: Not needed for MVP. Design with interfaces to allow future swapping.

## Failure Modes and Recovery

### Agent Failure
**Scenario**: LLM call times out or returns garbage

**Recovery**:
1. Retry with exponential backoff
2. Escalate to larger model
3. If persistent, mark iteration failed
4. Continue to next iteration or fail sprint

### Container Failure
**Scenario**: Container crashes or hangs

**Recovery**:
1. Kill container after timeout
2. Log diagnostics
3. Retry with fresh container
4. If persistent, escalate to human

### System Failure
**Scenario**: Process crashes mid-sprint

**Recovery**:
1. Persistent state in SQLite
2. On restart, load active sprints
3. Resume from last known state
4. Event log enables replay/reconstruction

### Provider Failure
**Scenario**: LLM provider down

**Recovery**:
1. Automatic failover to backup provider
2. Queue tasks for retry
3. Notify user of delay

## Security Model

### Threat Model

**Threats**:
1. Malicious code execution
2. Resource exhaustion
3. Secret leakage
4. Privilege escalation

**Mitigations**:

1. **Container Isolation**
   - No host network access
   - Read-only root filesystem
   - Resource limits (CPU, memory, disk)
   - No privileged mode
   - User namespace remapping

2. **Sandbox Enforcement**
   - Workspace directory restrictions
   - Path traversal prevention
   - No file access outside workspace

3. **Command Whitelist**
   - Only build, test, package commands
   - No arbitrary shell execution
   - Timeout enforced

4. **Secret Management**
   - API keys via environment only
   - No secrets in code or logs
   - Audit access

5. **Audit Trail**
   - All actions logged
   - Immutable logs (append-only)
   - Tamper-evident (signed)

## Performance Targets

### Latency
- Task submission to sprint start: <5s
- Spec iteration: <5m
- Test generation: <10m
- Implementation iteration: <15m
- Dashboard refresh: <2s

### Throughput
- Concurrent sprints: 3+
- Events/second: 100+
- Log ingestion: 1000 lines/sec

### Resource Usage
- Base memory: <2GB
- Per-task overhead: ~1GB
- Disk: <10GB (configurable retention)

## Monitoring and Observability

### Metrics
- Sprint success rate
- Average time per phase
- Iteration counts
- Cost per task
- Test coverage
- Agent utilization
- Error rates

### Logging
- Structured JSON logging
- Multiple levels (DEBUG, INFO, WARN, ERROR)
- Separate audit log
- Agent activity logs
- Container execution logs

### Alerting
- Sprint failures
- Cost threshold exceeded
- Provider failures
- System errors

## Research Influences

This architecture is informed by:

1. **`docs/RESEARCH.md`**
   - Start simple, add complexity only when needed
   - Verification beats self-confidence
   - Bounded loops over open-ended autonomy
   - External feedback (tests, compilation) is critical

2. **Scrum Methodology**
   - Roles with clear responsibilities
   - Time-boxed iterations
   - Definition of done
   - Transparency and inspection

3. **Container Best Practices**
   - Immutable infrastructure
   - Ephemeral environments
   - Resource constraints

4. **Event-Driven Architecture**
   - Loose coupling
   - Audit trail
   - Scalability

## Future Evolution

### Short Term (MVP)
- Single-node operation
- Docker containers
- Major LLM providers (OpenAI, OpenRouter)
- Python, Node.js, Go support

### Medium Term (Phase 2)
- Multi-node support
- Kubernetes integration
- More languages (Rust, Java, etc.)
- GitHub/GitLab integration
- Skill system

### Long Term (Phase 3+)
- Automatic spec generation from code
- Learned optimization (model selection, iteration budgets)
- Natural language queries
- IDE integration
- Distributed container builds

## Questions and Decisions

### Open Questions
1. Should we support non-containerized execution for trusted environments?
2. How to handle tasks requiring external APIs (integration tests)?
3. What's the right granularity for cost tracking?
4. Should agents have access to internet for research?

### Decisions Made
1. **Go for core**: Speed, reliability, deployment simplicity
2. **Containers for isolation**: Industry standard, well-understood
3. **SQLite for state**: Simple, embedded, sufficient for single-node
4. **YAML for config**: Human-readable, standard
5. **JSON for specs**: Structured, versionable
6. **Markdown for docs**: Portable, version-controlled

---

*This document is a living document. Update it as architecture evolves.*
