# Reference Map

This document explains what Scru-LLM should borrow from nearby projects.

## Reference Inputs

- `docs/RESEARCH.md` imported from `~/Research.md`
- `../delta-code`
- `../kokoclaw`

## From `delta-code`

Copy or adapt:

- deterministic planner mindset
- bounded autonomy loops
- role-to-model routing
- prompt modularity
- benchmark and eval discipline
- compact telemetry bus patterns
- explicit work-order style thinking

Best candidate areas:

- planner and task classification ideas
- telemetry/event bus shape
- structured prompt loading
- verification-first orchestration loops

Do not copy blindly:

- experimental surfaces that expand scope before MVP
- benchmark-specific scaffolding that is not needed for the first runnable slice
- complexity added only for many-model experiments

## From `kokoclaw`

Copy or adapt:

- operator/runtime separation
- durable runtime state
- approval and audit surfaces where relevant
- bounded autonomy loop patterns
- event publication and observability

Best candidate areas:

- runtime coordination patterns
- session/task state persistence
- autonomy loop guardrails
- telemetry and logging structure

Do not copy blindly:

- Discord-specific code
- chat-product assumptions
- broad tool surface unrelated to the MVP

## From `Research.md`

Apply directly:

- default to the simplest workflow that can pass evals
- separate cheap mechanical work from expensive reasoning
- prefer bounded evaluator-optimizer loops
- use grounded verification instead of self-confidence
- avoid multi-agent expansion without measured gain

## Recommendation For The Implementation Agent

Use `delta-code` and `kokoclaw` as pattern libraries, not as upstreams to clone wholesale.

Preferred strategy:

1. lift proven primitives
2. adapt names and contracts to Scru-LLM
3. keep the first runtime smaller than either reference project
4. only expand after the MVP is working and measurable
