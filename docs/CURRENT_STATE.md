# Current State

This document is the source of truth for what exists in Scru-LLM today.

## Snapshot: MVP IMPLEMENTED ✅

As of March 26, 2026, this repository contains a **working MVP implementation**.

### What's Implemented

**Core System:**
- ✅ CLI interface (`cmd/scru-llm`)
- ✅ Configuration system with YAML support
- ✅ Task, Sprint, and Event types (contracts)
- ✅ File-backed persistence (JSON for tasks/sprints, JSONL for events)
- ✅ Workspace management with artifact directories

**LLM Infrastructure:**
- ✅ HTTP client for OpenAI-compatible APIs
- ✅ Model router with role-based selection
- ✅ Cost tracking
- ✅ Multi-provider support (OpenRouter)

**Sprint Orchestration:**
- ✅ SprintManager with 5-phase workflow
- ✅ Phase transitions and iteration tracking
- ✅ Event logging
- ✅ Synchronous and asynchronous execution modes

**Worker Agents:**
- ✅ ProductOwner - creates initial specifications
- ✅ SpecEngineer - tightens specifications
- ✅ CodeEngineer - generates tests and implements code
- ✅ Reviewer - validates outputs (cross-family model)
- ✅ Verifier - runs tests locally

**CLI Commands:**
- ✅ `scru-llm task` - Create tasks
- ✅ `scru-llm run` - Execute sprints
- ✅ `scru-llm status` - View status
- ✅ `scru-llm logs` - View event logs
- ✅ `scru-llm config` - Configuration management

### What's Not Yet Implemented

- ❌ Web dashboard (planned for post-MVP)
- ❌ Container runtime / sandboxing
- ❌ API server with WebSocket
- ❌ Multi-task concurrent execution (currently sequential)
- ❌ Advanced reviewer logic (currently basic auto-approve)
- ❌ Production-quality code generation (generates stubs)

## How To Use

### Quick Start

```bash
# Build
make build

# Set API key
export OPENROUTER_API_KEY=sk-or-v1-...

# Create and run a task
./bin/scru-llm task "Create a Python calculator CLI" --title "Calculator"
./bin/scru-llm run <task-id> --sync
```

### Model Configuration

Current model mapping (all ≤35B, fits 24GB VRAM):

| Role | Model | Family |
|------|-------|--------|
| scrum_master | qwen/qwen3.5-35b-a3b | Qwen |
| product_owner | openai/gpt-oss-20b | OpenAI |
| spec_engineer | openai/gpt-oss-20b | OpenAI |
| test_engineer | openai/gpt-oss-20b | OpenAI |
| code_engineer | qwen/qwen3-coder-next | Qwen |
| reviewer | mistralai/mistral-small-24b-instruct | Mistral |

### Workspace Structure

```
~/.scru-llm/workspaces/
└── <task-id>/
    ├── spec/           # spec_v001.json, spec_v002.json
    ├── tests/          # test files
    ├── src/            # source code
    ├── artifacts/      # build outputs
    └── logs/           # execution logs
```

## Implementation Notes

### Key Design Decisions

1. **Sequential Execution**: MVP runs one sprint at a time to control costs and complexity
2. **File-Based Storage**: Simple JSON/JSONL files for MVP, database can be added later
3. **Cross-Family Reviewer**: Mistral model catches bugs that Qwen models miss
4. **Lenient Parsing**: Robust JSON extraction from markdown code blocks
5. **Relaxed Verification**: MVP accepts basic test presence, strict verification in v2

### Known Limitations

1. **Code Quality**: Generated code is basic stubs, not production-ready
2. **Test Coverage**: Tests are templates, need manual implementation
3. **Error Recovery**: Limited retry logic, failures stop the sprint
4. **Context Window**: Long specs may exceed model limits
5. **Model Availability**: Relies on OpenRouter providers being available

## Development Status

### Completed ✅
- [x] Core types and contracts
- [x] Configuration system
- [x] LLM client abstraction
- [x] Model routing
- [x] Sprint orchestration
- [x] File persistence
- [x] CLI interface
- [x] Basic workers
- [x] End-to-end workflow

### In Progress 🚧
- [ ] Enhanced code generation quality
- [ ] Better error handling and recovery
- [ ] Production-ready test generation

### Planned 📋
- [ ] Web dashboard
- [ ] Container sandboxing
- [ ] Concurrent task execution
- [ ] Advanced reviewer with multi-agent debate
- [ ] Integration tests

## Reading Order

For understanding the implementation:

1. [MVP.md](MVP.md) - What was built
2. [CONTRACTS.md](CONTRACTS.md) - Core types
3. [REFERENCE_MAP.md](REFERENCE_MAP.md) - Patterns from kokoclaw
4. Code in `internal/` - Implementation details

## Git Workflow

The implementation was done in a single session with the following structure:
- Core types and contracts
- Configuration and storage
- LLM infrastructure
- Workers and sprint orchestration
- CLI interface
- Bug fixes and polish

Future work should follow small, focused commits.
