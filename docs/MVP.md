# MVP

This document defines the first runnable slice of Scru-LLM.

## Goal

Build the smallest implementation that proves the workflow shape is useful for software tasks.

The MVP is not "full Scrum with many agents." The MVP is a bounded local workflow that can take one task from intake to verified output with clear artifacts.

## MVP Scope

Included:

- CLI-first interface only
- one task at a time on a single machine
- local workspace creation per task
- persisted task and sprint state on disk
- event log / audit log
- sequential workflow with explicit phases:
  - intake
  - spec
  - test
  - code
  - review
- bounded iteration counts per phase
- deterministic verification via tests or compile/check commands
- prompt loading from `config/prompts/`
- role-to-model configuration from `config/scru-llm.example.yaml`

Allowed simplifications:

- one orchestrator process can call multiple role prompts without needing true concurrent agents
- container support can start as a pluggable interface with a local shell runner fallback
- dashboard and REST API can be deferred
- only a minimal set of roles must exist at first

Deferred:

- multi-task scheduling
- browser dashboard
- REST API
- distributed execution
- parallel worker swarms
- long-term memory beyond compact task artifacts
- advanced delegation and branch-per-worker automation

## Required MVP Roles

The minimum useful role set is:

1. `product_owner`
2. `spec_engineer`
3. `code_engineer`
4. `reviewer`

The `scrum_master` may exist as orchestration logic in code before it exists as a separate prompt-driven agent.

`test_engineer` can be implemented in MVP if it is straightforward, but the system should not block on a separate specialized agent if a simpler workflow reaches the same result.

## Success Criteria

The MVP is successful when it can:

1. accept a task from the CLI
2. persist a task record and workspace
3. produce a spec artifact
4. refine the spec through a bounded review loop
5. produce test artifacts and implementation artifacts
6. run deterministic verification
7. surface phase-by-phase status and audit output
8. terminate cleanly as `completed` or `failed`

## Recommended First Build Order

1. Core contracts and state types
2. Config loader
3. File-backed task store and event log
4. Prompt loader
5. LLM client abstraction and mock mode
6. CLI task submission and task status
7. Sequential sprint manager
8. Verification runner
9. Reviewer loop

## Explicit Non-Goal For MVP

Do not start by building a "general autonomous multi-agent platform." Build a narrow workflow engine for one software task at a time.
