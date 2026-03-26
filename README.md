# Scru-LLM

Scru-LLM is a specification-first LLM orchestration system that follows a Scrum-shaped software delivery workflow. It takes user tasks through a structured pipeline: Intake → Spec Sprint → Test Sprint → Implementation Sprint → Review.

## Status: **MVP IMPLEMENTED & OPERATIONAL** ✅

This repository now contains a **working MVP implementation** that can process software tasks end-to-end.

### What's Working

- ✅ CLI interface for task management
- ✅ Full sprint workflow (5 phases)
- ✅ LLM integration with OpenRouter
- ✅ Multi-model orchestration (Qwen, Mistral, OpenAI families)
- ✅ File-backed persistence (tasks, sprints, events)
- ✅ Workspace management with artifact generation
- ✅ Synchronous and asynchronous execution modes

### Architecture

```
User Task
  -> Product Owner (GPT-OSS-20B)
  -> Scrum Master / Sprint Manager (Qwen3.5-35B-A3B)
  -> Spec Sprint (Spec Engineer: GPT-OSS-20B)
  -> Test Sprint (Test Engineer: GPT-OSS-20B)
  -> Implementation Sprint (Code Engineer: Qwen3-Coder-Next)
  -> Review + Verification (Reviewer: Mistral Small 24B)
  -> Deliverable
```

## Quick Start

```bash
# Build the binary
make build

# Configure your OpenRouter API key in .env file
echo "OPENROUTER_API_KEY=sk-or-v1-..." > .env

# Or manually edit .env in the project root or ~/.scru-llm/

# Create a task
./bin/scru-llm task "Create a Python CLI calculator" --title "Calculator"

# Run the sprint (asynchronous)
./bin/scru-llm run <task-id>

# Or run synchronously to see progress
./bin/scru-llm run <task-id> --sync

# Check status
./bin/scru-llm status <task-id>
./bin/scru-llm logs --task <task-id>
```

## Model Configuration

The system uses a cross-family model setup optimized for 24GB VRAM:

| Role | Model | Family | VRAM |
|------|-------|--------|------|
| scrum_master | qwen/qwen3.5-35b-a3b | Alibaba (Qwen) | ~17-22GB |
| product_owner | qwen/qwen3.5-35b-a3b | Alibaba (Qwen) | ~17-22GB |
| spec_engineer | qwen/qwen3.5-35b-a3b | Alibaba (Qwen) | ~17-22GB |
| test_engineer | qwen/qwen3.5-coder-32b | Alibaba (Qwen) | ~17GB |
| code_engineer | qwen/qwen3.5-coder-32b | Alibaba (Qwen) | ~17GB |
| reviewer | microsoft/phi-4-reasoning | Microsoft | ~10-12GB |

**Note**: The reviewer uses Microsoft Phi-4 (cross-family) to catch bugs that Qwen models might miss with its explicit reasoning capabilities.

## Key Docs

- [docs/CURRENT_STATE.md](docs/CURRENT_STATE.md): Current implementation status
- [docs/MVP.md](docs/MVP.md): MVP specification
- [docs/CONTRACTS.md](docs/CONTRACTS.md): Core types and state contracts
- [docs/REFERENCE_MAP.md](docs/REFERENCE_MAP.md): Reusable patterns from kokoclaw
- [docs/AGENT_IMPLEMENTATION_WORKFLOW.md](docs/AGENT_IMPLEMENTATION_WORKFLOW.md): Implementation guidelines

## Commands

### Task Management
```bash
# Create a new task
scru-llm task "description" --title "Title"

# View task status
scru-llm status [task-id]

# View system status
scru-llm status
```

### Sprint Execution
```bash
# Run sprint asynchronously
scru-llm run <task-id>

# Run sprint synchronously (see real-time output)
scru-llm run <task-id> --sync

# Run in mock mode (no LLM calls)
scru-llm run <task-id> --mock
```

### Monitoring
```bash
# View event logs
scru-llm logs
scru-llm logs --task <task-id>
scru-llm logs --limit 100
```

### Configuration
```bash
# Show current config
scru-llm config show

# Validate config
scru-llm config validate
```

## Workspace Structure

Each task gets a workspace at `~/.scru-llm/workspaces/<task-id>/`:

```
<task-id>/
├── spec/           # Specifications (spec_v001.json, spec_v002.json, ...)
├── tests/          # Test files (test_main.py, etc.)
├── src/            # Source code (main.py, etc.)
├── artifacts/      # Build outputs
└── logs/           # Execution logs
```

## Configuration

Copy the example config and customize:

```bash
cp config/scru-llm.example.yaml ~/.scru-llm/scru-llm.yaml
```

Set your OpenRouter API key in `.env`:

```bash
# In project root
echo "OPENROUTER_API_KEY=sk-or-v1-..." > .env

# Or in ~/.scru-llm/
echo "OPENROUTER_API_KEY=sk-or-v1-..." > ~/.scru-llm/.env
```

The system automatically loads `.env` from the project directory or `~/.scru-llm/.env`.

## Building

```bash
# Build binary
make build

# Run tests
make test

# Clean build artifacts
make clean
```

## Project Structure

```
.
├── bin/                    # Compiled binary
├── cmd/
│   └── scru-llm/          # Main entry point
├── config/
│   ├── prompts/           # System prompts for agents
│   └── scru-llm.example.yaml
├── internal/
│   ├── app/               # Application context
│   ├── cli/               # CLI commands
│   ├── config/            # Configuration loading
│   ├── errors/            # Sentinel errors
│   ├── llm/               # LLM client and router
│   ├── prompts/           # Prompt loader
│   ├── scrum/             # Sprint orchestration
│   ├── store/             # Persistence layer
│   ├── types/             # Core domain types
│   ├── workers/           # LLM agent workers
│   └── workspace/         # Workspace management
└── docs/                  # Documentation
```

## License

MIT License - see LICENSE file for details.

## Research Foundation

This implementation is grounded in:

- Research on specification-first LLM orchestration
- Prior patterns from `../kokoclaw` (model routing, role-based architecture)
- Modern LLM benchmarks and best practices

The guiding constraint: start with the smallest workflow that can be evaluated, then expand when measured gains justify complexity.
