# Contracts

This document defines the canonical core contracts for early implementation.

If other docs disagree, this file wins until code exists.

## Task Status

```text
pending
spec_sprint
test_sprint
implementation_sprint
review
completed
failed
cancelled
```

## Sprint Phase

```text
intake
spec
test
code
review
done
failed
```

## Core Types

### Task

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "pending|spec_sprint|test_sprint|implementation_sprint|review|completed|failed|cancelled",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp",
  "workspace_dir": "string",
  "current_sprint_id": "uuid or empty",
  "final_outcome": "completed|failed|cancelled or empty"
}
```

### Sprint

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "phase": "intake|spec|test|code|review|done|failed",
  "iteration": 0,
  "max_iterations": {
    "spec": 3,
    "test": 3,
    "code": 5,
    "review": 2
  },
  "started_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp",
  "completed_at": "RFC3339 timestamp or empty",
  "failure_reason": "string or empty"
}
```

### Specification

```json
{
  "task_id": "uuid",
  "version": 1,
  "title": "string",
  "description": "string",
  "requirements": ["string"],
  "acceptance_criteria": ["string"],
  "constraints": ["string"],
  "dependencies": ["string"],
  "edge_cases": [
    {
      "scenario": "string",
      "expected_behavior": "string"
    }
  ],
  "context": {
    "existing_codebase": true,
    "language": "string",
    "framework": "string"
  },
  "changes": ["string"],
  "testability_notes": "string",
  "confidence": 0.0
}
```

### Review Decision

```json
{
  "approved": true,
  "phase": "spec|test|code|review",
  "summary": "string",
  "issues": ["string"],
  "recommended_next_action": "retry|advance|fail"
}
```

### Event

```json
{
  "timestamp": "RFC3339 timestamp",
  "task_id": "uuid or empty",
  "sprint_id": "uuid or empty",
  "type": "string",
  "source": "string",
  "message": "string",
  "meta": {
    "key": "value"
  }
}
```

## Event Types

Start with this minimum set:

- `task.created`
- `task.updated`
- `sprint.started`
- `phase.changed`
- `iteration.started`
- `iteration.completed`
- `artifact.written`
- `verification.started`
- `verification.completed`
- `review.completed`
- `sprint.completed`
- `sprint.failed`

## Workspace Layout

```text
workspaces/{task_id}/
  spec/
  tests/
  src/
  artifacts/
  logs/
```

## Artifact Naming

Use deterministic names first:

- `spec/spec_v001.json`
- `tests/tests_v001.md` or language-specific test files
- `src/` for generated or modified source
- `artifacts/review_v001.json`
- `logs/events.jsonl`
- `logs/audit.jsonl`

## Implementation Rule

Before expanding the contracts, prefer adding one more field to an existing type over creating a new parallel type with overlapping meaning.
