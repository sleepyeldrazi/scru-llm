# Sprint Lifecycle

Detailed specification of the sprint phases and transitions.

## Overview

A sprint in Scru-LLM compresses traditional Scrum into task-completion-driven events. While traditional Scrum has 2-week sprints with daily standups, Scru-LLM has phases that complete as soon as quality gates are met.

## Sprint States

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│  PENDING │ → │ SPEC    │ → │  TEST   │ → │   CODE  │ → │  DONE   │
└─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘
     ↓              ↓              ↓              ↓
  FAILED        FAILED         FAILED         FAILED
```

## Phase 1: Task Intake

**Purpose**: Receive and analyze user task

**Duration**: Single pass (~1-5 minutes)

### Entry Criteria
- User submitted task via CLI or API
- Task has description

### Process

1. **Task Creation**
   ```go
   task := &Task{
       ID: generateUUID(),
       Description: userInput,
       Status: "pending",
       CreatedAt: time.Now(),
   }
   ```

2. **Product Owner Analysis**
   - Product Owner Agent reads task description
   - Analyzes requirements
   - Identifies deliverables
   - Assesses complexity
   - Checks for ambiguities

3. **Specification Generation**
   Output: Specification v0
   ```json
   {
     "task_id": "uuid",
     "title": "Human-readable title",
     "description": "Expanded description",
     "requirements": ["list", "of", "requirements"],
     "acceptance_criteria": ["testable", "criteria"],
     "constraints": ["limitations", "assumptions"],
     "context": {
       "language": "python",
       "existing_codebase": false
     }
   }
   ```

4. **Scrum Master Review**
   - Checks spec completeness
   - If acceptable → proceed to Spec Sprint
   - If unclear → request user clarification

### Exit Criteria
- [ ] Specification v0 created
- [ ] Basic requirements documented
- [ ] Context identified (language, framework)

### Failure Modes
- **Unclear Task**: Cannot determine what to build
  - Resolution: Ask user for clarification
- **Unsupported Language**: No base image available
  - Resolution: Reject with explanation

### Events Emitted
- `task.created`
- `specification.generated` (v0)
- `phase.changed` (pending → spec)

---

## Phase 2: Spec Sprint

**Purpose**: Tighten specification until unambiguous and testable

**Duration**: Up to N iterations (default: 3 iterations, ~5 minutes each)

### Entry Criteria
- Specification v0 exists
- Scrum Master approved proceeding

### Process per Iteration

1. **Spec Engineer Review**
   - Read current specification (vN)
   - Identify ambiguities
   - Find undefined terms
   - Check for missing edge cases
   - Verify testability

2. **Spec Tightening**
   - Rewrite requirements to be measurable
   - Define undefined terms
   - Add acceptance criteria
   - Document edge cases
   - Generate changelog

3. **Output Generation**
   ```json
   {
     "version": 2,
     "specification": { /* tightened spec */ },
     "changes": [
       "Clarified input format",
       "Added error handling requirements",
       "Defined edge cases"
     ],
     "confidence": 0.85
   }
   ```

4. **Review Agent Validation**
   Check criteria:
   - [ ] All requirements measurable?
   - [ ] All acceptance criteria testable?
   - [ ] No undefined terms?
   - [ ] Edge cases documented?
   - [ ] Confidence > 0.8?

5. **Decision**
   - If all checks pass → exit to Test Sprint
   - If any check fails and iterations remaining → continue
   - If max iterations reached → fail sprint

### Exit Criteria
- [ ] All requirements measurable
- [ ] All acceptance criteria testable
- [ ] No undefined terms
- [ ] Edge cases documented
- [ ] Confidence score > 0.8
- [ ] Review Agent approval

### Failure Modes

**Iteration Exhaustion**:
- Max iterations reached without passing review
- Resolution: Fail sprint with partial spec

**Persistent Ambiguity**:
- Specification cannot be tightened
- Resolution: Request user intervention

**Confidence Threshold**:
- Never reaches confidence > 0.8
- Resolution: Fail or proceed with warning

### Events Emitted
- `sprint.started`
- `iteration.started`
- `spec.tightened` (per iteration)
- `spec.reviewed` (per iteration)
- `phase.changed` (spec → test) or `sprint.failed`

### Quality Gates

**Gate 1: Measurability**
```
Bad:  "The system should be fast"
Good: "The system should respond to requests in <100ms"
```

**Gate 2: Testability**
```
Bad:  "Should handle errors gracefully"
Good: "Should return HTTP 400 with error message for invalid JSON input"
```

**Gate 3: Completeness**
- All functional requirements covered
- Non-functional requirements specified
- Constraints documented

---

## Phase 3: Test Sprint

**Purpose**: Generate comprehensive test suite from specification

**Duration**: Up to N iterations (default: 3 iterations, ~10 minutes each)

### Entry Criteria
- Finalized specification (ready state)
- All spec quality gates passed

### Process per Iteration

1. **Test Engineer Analysis**
   - Read finalized specification
   - Identify testable requirements
   - Plan test coverage strategy

2. **Test Generation**
   - Generate unit tests
   - Generate integration tests (if applicable)
   - Add edge case tests
   - Add error condition tests

3. **Code Engineer - Red Phase**
   - Write minimal "stub" implementation
   - Just enough to compile/import
   - All methods return zero values or panic

4. **Container Test Run**
   - Build container with stubs
   - Run test suite
   - Expect failures (confirm tests are real)
   - Capture test output

5. **Review Agent Validation**
   Check criteria:
   - [ ] 80%+ line coverage?
   - [ ] All public functions tested?
   - [ ] Edge cases covered?
   - [ ] Tests actually fail against stubs?
   - [ ] Test quality adequate?

6. **Decision**
   - If all checks pass → exit to Implementation Sprint
   - If any check fails and iterations remaining → continue
   - If max iterations reached → fail sprint

### Exit Criteria
- [ ] 80%+ line coverage
- [ ] All public functions have tests
- [ ] Edge cases covered
- [ ] Error conditions tested
- [ ] Tests fail against stubs (validated)
- [ ] Review Agent approval

### Test Structure
```
workspaces/{task_id}/tests/
├── unit/
│   ├── test_module_a.py
│   └── test_module_b.py
├── integration/
│   └── test_integration.py
├── conftest.py
└── coverage_report.json
```

### Failure Modes

**Low Coverage**:
- Cannot achieve 80% coverage
- Resolution: Add more tests or reduce scope

**Fake Tests**:
- Tests pass without real implementation
- Resolution: Review and fix tests

**Edge Case Gaps**:
- Important edge cases missing
- Resolution: Add tests, re-run

### Events Emitted
- `phase.changed` (spec → test)
- `iteration.started`
- `tests.generated`
- `tests.executed`
- `tests.reviewed`
- `phase.changed` (test → implementation) or `sprint.failed`

---

## Phase 4: Implementation Sprint

**Purpose**: Write code to pass all tests

**Duration**: Up to N iterations (default: 5 iterations, ~15 minutes each)

### Entry Criteria
- Comprehensive test suite ready
- All test quality gates passed

### Process per Iteration

1. **Code Engineer Analysis**
   - Read specification
   - Read test suite
   - Plan implementation strategy

2. **Code Generation**
   - Implement functionality
   - Follow coding guidelines
   - Add documentation/comments

3. **Container Test Run**
   - Build container with implementation
   - Run full test suite
   - Capture results

4. **Result Analysis**

   **If Tests Pass**:
   - Review Agent performs code review
   - If review passes → exit to Final Review
   - If review fails → continue with feedback

   **If Tests Fail**:
   - Code Engineer receives test output
   - Debug and identify issues
   - Fix code
   - Continue to next iteration

5. **Decision**
   - If tests pass and review passes → proceed
   - If tests fail and iterations remaining → continue
   - If max iterations reached → fail sprint

### Exit Criteria
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Code review passed
- [ ] Coding guidelines followed
- [ ] Documentation complete

### Code Quality Gates

**Gate 1: Correctness**
- All tests pass
- No regressions

**Gate 2: Style**
- Follows language conventions
- Passes linting
- Readable and maintainable

**Gate 3: Documentation**
- Functions documented
- Complex logic explained
- Examples provided (if public API)

### Failure Modes

**Test Failures**:
- Cannot make tests pass
- Resolution: Escalate to larger model or fail

**Code Quality**:
- Passes tests but review fails
- Resolution: Fix issues, re-review

**Complexity**:
- Implementation too complex
- Resolution: Refactor or simplify requirements

### Events Emitted
- `phase.changed` (test → implementation)
- `iteration.started`
- `code.implemented`
- `tests.executed`
- `code.reviewed`
- `phase.changed` (implementation → review) or `sprint.failed`

---

## Phase 5: Final Review and Delivery

**Purpose**: Validate complete deliverable

**Duration**: Single pass (~5-10 minutes)

### Entry Criteria
- Implementation passes tests
- Code review passed

### Process

1. **Clean Environment Test**
   - Build fresh container
   - Run full test suite
   - Verify reproducibility

2. **Final Code Review**
   - Review Agent checks quality
   - Verify no debug code left
   - Check for security issues
   - Validate documentation

3. **Artifact Generation**
   - Package source code
   - Generate documentation
   - Create deliverable archive

4. **Completion Report**
   ```json
   {
     "sprint_id": "uuid",
     "task_id": "uuid",
     "status": "completed",
     "phases": {
       "spec": { "iterations": 2, "duration": "8m" },
       "test": { "iterations": 1, "duration": "12m" },
       "implementation": { "iterations": 3, "duration": "35m" }
     },
     "total_duration": "55m",
     "test_coverage": 87,
     "artifacts": {
       "spec": "spec/spec_final.json",
       "tests": "tests/",
       "code": "src/",
       "docs": "README.md"
     }
   }
   ```

### Exit Criteria
- [ ] All tests pass in clean environment
- [ ] Final code review passed
- [ ] Documentation complete
- [ ] Artifacts packaged

### Failure Modes

**Clean Environment Failure**:
- Tests pass in dev but fail in clean container
- Resolution: Fix environment issues

**Final Review Failure**:
- Quality issues found
- Resolution: Fix and re-run final review

### Events Emitted
- `phase.changed` (implementation → review)
- `tests.executed` (final)
- `code.reviewed` (final)
- `sprint.completed`
- `artifacts.generated`

---

## State Transitions

### State Machine

```go
type SprintState string

const (
    StatePending         SprintState = "pending"
    StateSpecSprint      SprintState = "spec_sprint"
    StateTestSprint      SprintState = "test_sprint"
    StateImplementation  SprintState = "implementation"
    StateReview          SprintState = "review"
    StateCompleted       SprintState = "completed"
    StateFailed          SprintState = "failed"
)

var validTransitions = map[SprintState][]SprintState{
    StatePending:        {StateSpecSprint, StateFailed},
    StateSpecSprint:     {StateTestSprint, StateFailed},
    StateTestSprint:     {StateImplementation, StateFailed},
    StateImplementation: {StateReview, StateFailed},
    StateReview:         {StateCompleted, StateFailed},
}

func (s *Sprint) Transition(to SprintState) error {
    valid := validTransitions[s.State]
    for _, v := range valid {
        if v == to {
            s.State = to
            return nil
        }
    }
    return fmt.Errorf("invalid transition: %s → %s", s.State, to)
}
```

### Transition Triggers

| From | To | Trigger |
|------|-----|---------|
| pending | spec_sprint | Product Owner spec approved |
| pending | failed | Unclear requirements |
| spec_sprint | test_sprint | Spec review passed |
| spec_sprint | failed | Max iterations reached |
| test_sprint | implementation | Test review passed |
| test_sprint | failed | Max iterations reached |
| implementation | review | Tests pass, code review passed |
| implementation | failed | Max iterations reached |
| review | completed | Final review passed |
| review | failed | Critical issues found |

---

## Iteration Budgets

### Configuration

```yaml
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
```

### Budget Exhaustion Strategy

When iteration budget exhausted:

1. **Log diagnostics**
   - What was attempted
   - Why it failed
   - Current state

2. **Escalation Options**
   - Try with larger model
   - Reduce scope
   - Request human intervention

3. **Failure Documentation**
   - Partial deliverables saved
   - Logs preserved
   - Reason for failure recorded

---

## Parallel Execution

### What Can Run in Parallel

Within a sprint, some activities can be parallel:

1. **Spec Sprint**
   - Spec Engineer working
   - Review Agent can be warming up

2. **Test Sprint**
   - Test generation
   - Container building (if deterministic)

3. **Implementation Sprint**
   - Code implementation
   - Review preparation

### What Must Be Sequential

1. **Phase Gates**
   - Must pass review before next phase
   
2. **Container Tests**
   - Must build before testing
   - Must test before review

### Resource Constraints

- Max concurrent agents per sprint: 2-3
- Container builds: 1 at a time (to avoid conflicts)
- LLM calls: Limited by provider rate limits

---

## Monitoring and Metrics

### Per-Sprint Metrics

```go
type SprintMetrics struct {
    SprintID          string
    TaskID            string
    StartTime         time.Time
    EndTime           time.Time
    TotalDuration     time.Duration
    
    SpecPhase         PhaseMetrics
    TestPhase         PhaseMetrics
    ImplementationPhase PhaseMetrics
    
    TestCoverage      float64
    LinesOfCode       int
    TotalCost         float64
}

type PhaseMetrics struct {
    Iterations        int
    Duration          time.Duration
    Success           bool
    Cost              float64
}
```

### Health Indicators

**Healthy Sprint**:
- Completes within global timeout
- Iteration count < max for all phases
- Test coverage ≥ 80%
- Cost within budget

**Unhealthy Patterns**:
- Always hitting max iterations
- Low test coverage
- High cost per task
- Frequent failures

---

## Recovery and Resume

### Crash Recovery

If system crashes mid-sprint:

1. On restart, load active sprints from database
2. Determine current phase and iteration
3. Resume from last known state
4. Log recovery event

### Manual Intervention Points

User can intervene at:
- Any phase transition
- On failure
- Via `scru-llm cancel` command

### Partial Deliverables

Even on failure, save:
- Current specification
- Generated tests
- Partial implementation
- All logs and audit trails

---

*This document defines the core sprint lifecycle. Implementations must follow these specifications for consistency.*
