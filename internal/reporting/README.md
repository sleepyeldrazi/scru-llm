# Reporting Package

This package implements progress reporting and observability for Scru-LLM.

## What This Package Does

Provides visibility into system operation:
- Web dashboard with Kanban board
- CLI progress reports
- Audit logging
- Agent activity tracking
- Metrics collection

## Files

- `dashboard.go` - Web server and dashboard UI
- `kanban.go` - Kanban board generation (HTML and Markdown)
- `audit.go` - Structured audit logging (JSONL)
- `agent_logger.go` - Per-agent activity logging
- `metrics.go` - Prometheus-style metrics

## Architecture

### Dashboard
- Echo web framework
- Static HTML/CSS/JS (no heavy frameworks)
- Real-time updates via WebSocket
- Kanban board with task cards
- Responsive design

### Kanban Board
Columns:
- Backlog
- Spec
- Testing
- Coding
- Review
- Done

Two modes:
- **Web UI**: Interactive HTML
- **CLI**: Markdown table fallback

### Audit Logging
Structured JSONL format:
```json
{
  "timestamp": "2024-03-26T14:23:45Z",
  "event_type": "agent.completion",
  "task_id": "...",
  "agent_id": "...",
  "action": "spec_tightening",
  "inputs": {...},
  "outputs": {...},
  "metadata": {...}
}
```

Features:
- Append-only (tamper-evident)
- Log rotation
- Compression after 7 days

### Agent Activity Logging
Per-agent logs tracking:
- Files read/written
- LLM calls made
- Commands executed
- Duration of operations

### Metrics
Prometheus-compatible metrics:
- Sprint success rate
- Average time per phase
- Iteration counts
- Cost per task
- Agent utilization

## Usage Example

```go
// Create dashboard
dashboard := reporting.NewDashboard(config.Dashboard)

// Start server
go dashboard.Start()

// Log event
auditLog := reporting.NewAuditLogger(config.AuditLog)
auditLog.Log(reporting.Event{
    Type: "sprint.completed",
    TaskID: taskID,
    Data: result,
})

// Update kanban
dashboard.UpdateKanban(sprint)
```

## WebSocket Events

Real-time updates sent to dashboard:
- `task.created`
- `sprint.started`
- `phase.changed`
- `iteration.completed`
- `sprint.completed`
- `sprint.failed`

## Future Improvements

- [ ] Natural language queries ("What's blocking task X?")
- [ ] Sprint post-mortem reports
- [ ] Cost analysis dashboard
- [ ] Mobile-optimized UI
- [ ] Dark mode
