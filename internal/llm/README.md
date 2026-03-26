# LLM Package

This package provides LLM client abstraction and routing for multiple providers.

## What This Package Does

- Abstracts different LLM providers (OpenAI, OpenRouter, local)
- Routes requests based on role configuration
- Handles retries and fallbacks
- Tracks costs
- Manages context windows

## Files

- `client.go` - LLM client interface
- `openai.go` - OpenAI-compatible client implementation
- `router.go` - Role-based routing
- `models.go` - Model definitions and capabilities

## Architecture

### Client Interface
```go
type Client interface {
    Complete(ctx context.Context, messages []Message, opts Options) (*Response, error)
    Stream(ctx context.Context, messages []Message, opts Options) (chan Chunk, error)
}
```

### Provider Support
- OpenAI API
- OpenRouter
- Local models (llama.cpp, etc.)
- Any OpenAI-compatible endpoint

### Router
Maps roles to specific models:
```yaml
roles:
  code_engineer:
    provider: openrouter
    model: anthropic/claude-3.5-sonnet
    temperature: 0.7
```

Features:
- Provider fallback on failure
- Cost tracking per call
- Retry with exponential backoff
- Request timeout handling

## Usage Example

```go
// Create client from config
client := llm.NewClient(cfg.LLM)

// Route by role
messages := []llm.Message{
    {Role: "system", Content: "You are a helpful assistant"},
    {Role: "user", Content: "Write a function..."},
}

response, err := client.CompleteForRole(ctx, "code_engineer", messages)

// Or direct completion
response, err := client.Complete(ctx, messages, llm.Options{
    Model: "gpt-4",
    Temperature: 0.7,
})
```

## Message Format

```go
type Message struct {
    Role    string // "system", "user", "assistant"
    Content string
}
```

## Response Format

```go
type Response struct {
    Content      string
    TokensUsed   int
    Cost         float64
    Model        string
    FinishReason string
}
```

## Cost Tracking

Per-request tracking:
- Input tokens × input price
- Output tokens × output price
- Provider-specific pricing

Aggregated by:
- Task
- Sprint
- Agent/Role
- Time period

## Future Improvements

- [ ] Caching layer for common prompts
- [ ] Request batching
- [ ] Streaming response aggregation
- [ ] Model capability detection
- [ ] Dynamic model selection
