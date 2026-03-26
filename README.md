# Scru-LLM

Scru-LLM is a specification-first skeleton for an LLM orchestration system that follows a Scrum-shaped software delivery workflow.

## Status

This repository is intentionally a planning and handoff skeleton.

- It is not a runnable product yet.
- It exists to give a frontier model enough steering to implement the system later without inventing the product shape from scratch.
- The target system is documented here; the current implementation state is documented separately.

Read these first:

- [docs/CURRENT_STATE.md](docs/CURRENT_STATE.md)
- [docs/MVP.md](docs/MVP.md)
- [docs/CONTRACTS.md](docs/CONTRACTS.md)
- [docs/AGENT_IMPLEMENTATION_WORKFLOW.md](docs/AGENT_IMPLEMENTATION_WORKFLOW.md)

## Product Intent

Scru-LLM aims to turn a software task into a bounded, auditable workflow:

1. Intake the request and create an initial spec.
2. Tighten the spec until it is measurable and testable.
3. Generate tests from the finalized spec.
4. Implement code against those tests.
5. Review and verify the result with deterministic checks.
6. Report progress using Scrum-like phase transitions and artifacts.

The emphasis is on:

- spec-first development
- bounded autonomy
- grounded verification
- reproducible execution
- observable state and audit trails

## Architecture Snapshot

Target-state workflow:

```text
User Task
  -> Product Owner
  -> Scrum Master / Sprint Manager
  -> Spec Sprint
  -> Test Sprint
  -> Implementation Sprint
  -> Review + Verification
  -> Deliverable
```

This repo keeps the target design intentionally explicit, but the first implementation slice is narrower than the full vision. See [docs/MVP.md](docs/MVP.md).

## Key Docs

- [docs/CURRENT_STATE.md](docs/CURRENT_STATE.md): what exists right now
- [docs/MVP.md](docs/MVP.md): the first runnable slice to build
- [docs/CONTRACTS.md](docs/CONTRACTS.md): canonical core types and state contracts
- [docs/REFERENCE_MAP.md](docs/REFERENCE_MAP.md): what to reuse from `delta-code` and `kokoclaw`
- [docs/RESEARCH.md](docs/RESEARCH.md): imported research notes used as design guidance
- [docs/SPECIFICATIONS.md](docs/SPECIFICATIONS.md): fuller target-state specification
- [docs/architecture/OVERVIEW.md](docs/architecture/OVERVIEW.md): high-level design philosophy
- [docs/architecture/SPRINT_LIFECYCLE.md](docs/architecture/SPRINT_LIFECYCLE.md): target sprint phases
- [docs/architecture/RESEARCH_INFLUENCES.md](docs/architecture/RESEARCH_INFLUENCES.md): how research and prior projects inform the design

## Naming

Project name:

- display name: `Scru-LLM`
- binary name target: `scru-llm`
- config example: `config/scru-llm.example.yaml`
- local config path target: `config/scru-llm.yaml`

## Research Foundation

This skeleton is grounded in:

- your imported [docs/RESEARCH.md](docs/RESEARCH.md)
- prior implementation ideas from `../delta-code`
- prior runtime patterns from `../kokoclaw`

The guiding constraint from the research is simple: start with the smallest workflow that can be evaluated, then expand only when measured gains justify the added complexity.
