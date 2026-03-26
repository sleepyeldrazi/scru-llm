# Current State

This document is the source of truth for what exists in Scru-LLM today.

## Snapshot

As of March 26, 2026, this repository is a design skeleton, not an implementation.

Present:

- top-level product documentation
- architecture docs
- roadmap and coding guidelines
- example configuration
- initial prompt files for `product_owner` and `spec_engineer`
- package-level `README.md` placeholders under `internal/`
- empty tracked scaffolding directories for future implementation

Absent:

- no `cmd/scru-llm` implementation yet
- no concrete Go packages under `internal/`
- no working CLI, dashboard, API server, or container runtime
- no real event bus, state store, or orchestration loop
- no tests beyond future intent described in docs

## How To Read The Repo

Use the docs in this order:

1. [MVP.md](MVP.md)
2. [CONTRACTS.md](CONTRACTS.md)
3. [SPECIFICATIONS.md](SPECIFICATIONS.md)
4. [architecture/OVERVIEW.md](architecture/OVERVIEW.md)
5. [architecture/SPRINT_LIFECYCLE.md](architecture/SPRINT_LIFECYCLE.md)
6. [REFERENCE_MAP.md](REFERENCE_MAP.md)
7. [AGENT_IMPLEMENTATION_WORKFLOW.md](AGENT_IMPLEMENTATION_WORKFLOW.md)

Interpretation rules:

- `README.md` explains product intent, not current runtime behavior.
- `docs/SPECIFICATIONS.md` and `docs/architecture/*` describe target-state behavior unless this file says otherwise.
- `TODO.md` is a roadmap, not evidence of completed implementation.

## What This Skeleton Is Supposed To Do

The repo should be sufficient for a strong implementation model to:

- understand the product boundary
- understand the intended Scrum-shaped workflow
- know the first implementation slice
- know the core data contracts
- know which prior repos to mine for reusable patterns
- know how to work in the repo without thrashing or rewriting history

## What An Implementation Agent Must Not Assume

- Do not assume commands in older target-state docs already work.
- Do not assume multi-agent execution is the first milestone.
- Do not assume containers, web UI, and API all belong in the first runnable slice.
- Do not infer missing contracts from prose if a tighter contract can be added first.

## Immediate Repository Priorities

1. Keep the docs coherent and synchronized.
2. Implement the narrow MVP before expanding scope.
3. Reuse proven runtime pieces from `delta-code` and `kokoclaw` instead of rebuilding everything from scratch.
4. Preserve a clean git history with small, reviewable commits.
