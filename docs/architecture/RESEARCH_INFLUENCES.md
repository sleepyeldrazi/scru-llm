# Research Influences

This document catalogs the research and prior work that informs Scru-LLM's design.

## Primary Research Sources

### 1. Agent Systems Research (`docs/RESEARCH.md`, imported from `~/Research.md`)

**Key Principles Applied**:
- **Simplicity First**: Start with the simplest scaffold that can pass evals
- **Verification Over Confidence**: External feedback (tests, compilation) beats self-assessment
- **Bounded Loops**: Prefer bounded retries over open-ended autonomy
- **Separate Mechanical from Reasoning**: Use cheaper models for simple tasks, expensive ones for synthesis
- **Grounded Critique**: Use CRITIC-style tool-interactive verification

**Specific Applications**:
- Spec sprint has maximum iteration budget (bounded)
- Review Agent uses deterministic tests as external verification
- Role-based routing uses different models for different tasks
- Code Engineer receives test failures as grounded feedback

### 2. Scrum Methodology

**Traditional Scrum**:
- Roles: Product Owner, Scrum Master, Development Team
- Ceremonies: Sprint Planning, Daily Standup, Sprint Review, Retrospective
- Artifacts: Product Backlog, Sprint Backlog, Increment

**Scru-LLM Adaptations**:
- **Compressed Timeframes**: Days → task-completion events
- **Automated Ceremonies**: Daily standup is automatic status update
- **LLM Agents as Team**: Each role is an agent
- **Continuous Review**: Review happens per iteration, not just at sprint end
- **Test as Definition of Done**: Passing tests required for completion

### 3. SWE-agent (Princeton NLP, 2024)

**Key Insight**: The interface between model and environment is part of the model's performance.

**Applications**:
- Careful design of agent-computer interface (ACI)
- Short read/search/edit/test loops
- File localization before editing
- Agent-friendly output formats

### 4. Agentless (Xia et al., 2024)

**Key Insight**: Simpler pipelines can outperform complex agents at lower cost.

**Applications**:
- Resist urge to add complexity without measurement
- Benchmark against simple baselines
- If full agent loop doesn't help, cut it

### 5. ReAct (Yao et al., 2022)

**Key Insight**: Interleaving reasoning with external actions outperforms monolithic reasoning.

**Applications**:
- Agents can use tools (file read, search, execution)
- Reasoning informed by environment feedback
- Test failures drive debugging reasoning

### 6. CRITIC (Gou et al., 2024)

**Key Insight**: Self-correction is stronger when critique is grounded in external tools.

**Applications**:
- Review Agent uses test results, not just code inspection
- Code Engineer debugs using compiler output
- Spec validation against testability criteria

### 7. Building Effective AI Agents (Anthropic, 2024)

**Key Insights**:
- Choose between single-agent, workflow, and multi-agent designs intentionally
- Use simple patterns: sequential, parallel, evaluator-optimizer
- Match system complexity to business value

**Applications**:
- Multi-agent design (different roles)
- Evaluator-optimizer pattern (Review Agent + Code Engineer)
- Simple workflow orchestration (Sprint Manager)

### 8. Tree of Thoughts / LATS

**Key Insight**: Search over reasoning paths helps hard problems.

**Applications**:
- Multiple spec candidates generated
- Test suite variations
- Fallback model selection

**Caveat**: Expensive. Used selectively, not for routine tasks.

## Design Patterns from Prior Projects

### From kokoclaw

**What to Copy**:
- Shared operator runtime used by multiple interfaces
- Durable runtime state (sessions, approvals, audit)
- Event pipeline for observability
- Bounded autonomy loops
- Approval gates for privileged actions

**What to Adapt**:
- Remove Discord-first orientation
- Focus on software engineering tasks
- Add Scrum methodology
- Emphasize spec-first approach

### From delta-code

**What to Copy**:
- Modular prompt system
- Deterministic verification loops
- Benchmark-first development
- Cost/quality/latency tracking
- Role-based routing

**What to Adapt**:
- Add Scrum methodology layer
- Tighten spec before coding
- Container-first execution
- Compressed sprint cycles

## Anti-Patterns Avoided

Based on research, we explicitly avoid:

1. **Too Many Homogeneous Agents**
   - Not 10 code agents, but specialized roles
   - Diversity in role, not just count

2. **Persona-Heavy Prompts**
   - Task-focused, not character-focused
   - Constraints over acting

3. **Keyword-Triggered Tools**
   - Model decides when to search/execute
   - Measured behavior, not hardcoded rules

4. **Unbounded Self-Reflection**
   - Max iterations enforced
   - Progress measured, not assumed

5. **Auto-Commits Without Validation**
   - Tests must pass
   - Review must approve
   - Container must build

6. **Massive Context Files**
   - Keep specs focused
   - Operational constraints, not lore

7. **Platform Before Workflow**
   - Build narrow workflow first
   - Validate with real tasks
   - Expand only when needed

## Test-Driven Development Principles

### Red-Green-Refactor in Scru-LLM

**Traditional TDD**:
1. Write failing test (Red)
2. Write code to pass test (Green)
3. Refactor while keeping tests passing

**Scru-LLM TDD**:
1. **Spec Sprint**: Ensure requirements are clear and testable
2. **Test Sprint**: Generate comprehensive tests (Red phase)
3. **Implementation Sprint**: Write code to pass tests (Green phase)
4. **Review**: Quality check (Refactor phase)

### Test Quality Criteria

From research on software testing:

1. **Test the Contract**: Verify behavior, not implementation
2. **Edge Cases**: Boundary conditions, null inputs, error paths
3. **Independence**: Tests don't depend on each other
4. **Determinism**: Same input → same output
5. **Speed**: Fast enough to run frequently
6. **Readability**: Tests as documentation

## Container Best Practices

### From Industry Standards

1. **Immutable Infrastructure**: Containers are built, not modified
2. **Ephemeral**: Created for task, destroyed after
3. **Resource Constraints**: CPU, memory, disk limits
4. **Minimal Base Images**: Reduce attack surface
5. **Layer Caching**: Speed up rebuilds
6. **Multi-stage Builds**: Separate build and runtime

### Specific to Scru-LLM

- One container per task
- Pre-built base images for common languages
- Shared package caches (mounted volumes)
- Build artifacts extracted, not left in container
- Full cleanup after task completion

## Cost Optimization Research

### SOLVE-Med / MATA (2025-2026)

**Key Insight**: Small specialized models, when orchestrated well, can outperform larger standalone systems.

**Applications**:
- Route grep/read/run to small models
- Reserve large models for coding/reasoning
- 4B orchestrator, 14B coder pattern

### Cost-Aware Routing

From multi-agent research:

- High importance tasks: Best model
- Medium importance: Balance quality/cost
- Low importance: Cheapest adequate model

**Scru-LLM Application**:
- Product Owner: Large model (requirements critical)
- Scrum Master: Medium model (coordination)
- Code Engineer: Large model (quality critical)
- Review Agent: Medium model (sufficient for review)

## Evaluation and Reliability

### From Anthropic's "Demystifying evals for AI agents"

**Principles Applied**:
- Start evals early (20-50 tasks enough)
- Unambiguous tasks with reference solutions
- Evaluate both "should X" and "should not X"
- Isolate trials (no shared state)
- Grade outputs, not traces
- Read transcripts constantly

**Scru-LLM Implementation**:
- Built-in benchmark suite
- Success/failure tracking per sprint
- Audit logs for inspection
- Isolated workspaces per task
- Metrics collection from day one

## Future Research Directions

### Areas to Explore

1. **Automatic Spec Generation**
   - Generate specs from existing code
   - Reverse-engineer requirements

2. **Mutation Testing**
   - Verify tests catch bugs
   - Improve test quality

3. **Property-Based Testing**
   - Generate random inputs
   - Find edge cases automatically

4. **Formal Verification**
   - Prove correctness for critical parts
   - Integration with coding workflow

5. **Learned Optimization**
   - Predict iteration count needed
   - Optimize model selection
   - Learn from past tasks

6. **Natural Language Queries**
   - "Why did task X fail?"
   - "What's the slowest part of the pipeline?"

## References

### Papers

1. Kojima et al. (2022) - "Large Language Models are Zero-Shot Reasoners"
2. Wang et al. (2022) - "Self-Consistency Improves Chain of Thought Reasoning"
3. Yao et al. (2022) - "ReAct: Synergizing Reasoning and Acting"
4. Yao et al. (2023) - "Tree of Thoughts"
5. Schick et al. (2023) - "Toolformer"
6. Gou et al. (2024) - "CRITIC"
7. Madaan et al. (2023) - "Self-Refine"
8. Shinn et al. (2023) - "Reflexion"
9. Yang et al. (2024) - "SWE-agent"
10. Xia et al. (2024) - "Agentless"

### Articles

1. BAIR (2024) - "The Shift from Models to Compound AI Systems"
2. Anthropic (2024) - "Building Effective AI Agents"
3. Anthropic (2026) - "Demystifying evals for AI agents"

### Projects

1. karpathy/autoresearch - Narrow loop, fixed optimization target
2. SWE-agent - Repo-focused action surface
3. OpenHands - Workspace/runtime architecture
4. aider - Architect/editor split

---

*This document should be updated as new research influences the project.*
