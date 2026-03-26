# Scrum Package

This package implements the Scrum orchestration layer for Scru-LLM.

## What This Package Does

Manages the Scrum process:
- Task backlog and queue management
- Sprint lifecycle coordination
- Event-driven communication between components
- State persistence and recovery

## Files

- `backlog.go` - Task queue management with SQLite persistence
- `sprint.go` - Sprint lifecycle, phase transitions, iteration tracking
- `events.go` - Event bus implementation for decoupled communication
- `state.go` - State machine for sprint states

## Architecture

### Backlog Manager
- Maintains ordered queue of pending tasks
- Persists to SQLite for durability
- Handles task dependencies (DAG)
- State transitions: pending → spec → test → impl → review → done

### Sprint Manager
- Creates and manages sprint instances
- Coordinates phase transitions
- Enforces iteration budgets and timeouts
- Tracks sprint metrics

### Event Bus
- Async pub/sub for loose coupling
- JSONL persistence for audit trail
- Event replay capability
- Graceful shutdown

### Data Flow

```
User Task
    ↓
Backlog Manager (queue)
    ↓
Sprint Manager (create sprint)
    ↓
Event Bus (notify listeners)
    ↓
Agent Workers (pick up tasks)
```

## Usage Example

```go
// Create backlog manager
backlog := scrum.NewBacklogManager(db)

// Add task
task := &scrum.Task{
    Title: "Create CLI tool",
    Description: "...",
}
err := backlog.AddTask(task)

// Start sprint
sprintMgr := scrum.NewSprintManager(backlog, eventBus)
sprint, err := sprintMgr.StartSprint(task)

// Listen for events
eventBus.Subscribe("sprint.completed", func(e Event) {
    fmt.Println("Sprint completed!")
})
```

## Future Improvements

- [ ] Redis backend for distributed operation
- [ ] Priority queue support
- [ ] Sprint scheduling (start at specific time)
- [ ] Batch task processing
