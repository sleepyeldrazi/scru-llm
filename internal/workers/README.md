# Workers Package

This package implements the LLM agent roles (workers) that perform work in the Scrum system.

## What This Package Does

Provides specialized AI agents for different Scrum roles:
- Product Owner: Task analysis and specification
- Scrum Master: Coordination and reporting
- Spec Engineer: Specification tightening
- Test Engineer: Test generation
- Code Engineer: Implementation
- Review Agent: Quality validation

## Files

- `base.go` - Base worker functionality shared by all agents
- `product_owner.go` - Requirements analysis
- `scrum_master.go` - Sprint coordination
- `spec_engineer.go` - Spec iteration
- `test_engineer.go` - Test generation
- `code_engineer.go` - Code implementation
- `reviewer.go` - Quality review

## Architecture

### Base Worker
All agents embed `BaseWorker` which provides:
- LLM client access
- Logging and auditing
- Timeout handling
- Error handling patterns

### Agent Interface
```go
type Agent interface {
    Name() string
    Role() string
    Execute(ctx context.Context, input interface{}) (interface{}, error)
}
```

### Role Specialization
Each agent implements the interface with role-specific logic:
- Different prompts
- Different models (configurable)
- Different validation
- Different output formats

### Model Routing
Configured in `config/scru-llm.yaml`:
- Product Owner → Large model (reasoning)
- Code Engineer → Large model (coding)
- Scrum Master → Medium model (coordination)
- Reviewer → Medium model (sufficient)

## Usage Example

```go
// Create LLM client
client := llm.NewClient(cfg)

// Create agents
productOwner := workers.NewProductOwner(client)
specEngineer := workers.NewSpecEngineer(client)
reviewer := workers.NewReviewer(client)

// Execute product owner
spec, err := productOwner.Execute(ctx, taskDescription)

// Tighten spec
spec, err = specEngineer.Execute(ctx, spec)

// Review
review, err := reviewer.Execute(ctx, spec)
```

## Prompts

Each agent loads its system prompt from `config/prompts/`:
- `product-owner-system.md`
- `spec-engineer-system.md`
- etc.

Prompts are versioned with the code.

## Future Improvements

- [ ] Dynamic agent creation from config
- [ ] Agent learning from past tasks
- [ ] Parallel agent execution
- [ ] Agent specialization (frontend, backend, etc.)
