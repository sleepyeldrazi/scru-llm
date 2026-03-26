# Agent Implementation Workflow

This file is addressed to the future implementation agent building Scru-LLM.

## Read This First

Before making code changes, read:

1. [CURRENT_STATE.md](CURRENT_STATE.md)
2. [MVP.md](MVP.md)
3. [CONTRACTS.md](CONTRACTS.md)
4. [REFERENCE_MAP.md](REFERENCE_MAP.md)
5. [RESEARCH.md](RESEARCH.md)

## Git Workflow

Use git actively while implementing.

- Create a branch for each coherent task or milestone.
- Commit in small checkpoints, not one giant dump at the end.
- Write commit messages that describe the change in terms of behavior or architecture.
- Push branches if a remote exists.
- If no remote exists, still maintain local commits as recovery points.

Recommended branch names:

- `feat/mvp-task-store`
- `feat/mvp-sprint-manager`
- `feat/mvp-llm-router`
- `fix/spec-contract-drift`
- `docs/handoff-cleanup`

Recommended commit style:

- `docs: define MVP and core contracts`
- `feat: add file-backed task store`
- `feat: add sequential sprint manager`
- `test: cover phase transition rules`
- `refactor: split prompt loading from worker setup`

## Working Rules

- Keep the mutable surface small.
- Implement the MVP before expanding to the full target-state design.
- Reuse patterns from `delta-code` and `kokoclaw` where they clearly reduce risk.
- Do not silently drift away from the contracts in [CONTRACTS.md](CONTRACTS.md).
- Update docs when the code changes the meaning of a contract or workflow.

## Branch And Push Guidance

If a remote is configured:

- push after meaningful checkpoints
- keep `main` stable
- merge only after checks pass or an explicit exception is documented

If a remote is not configured:

- still branch locally
- still commit locally
- record in the final task summary what should be pushed once a remote exists

## Safety Rules

- Do not rewrite history unless explicitly asked.
- Do not force-push shared branches unless explicitly asked.
- Do not mark roadmap items complete unless the code or docs actually satisfy them.
- Do not broaden scope just because the architecture docs mention a future feature.
