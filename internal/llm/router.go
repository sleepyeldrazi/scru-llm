package llm

import (
	"context"
	"fmt"
	"sync"

	"github.com/sleepyeldrazi/scru-llm/internal/config"
)

// Router routes LLM requests to the appropriate provider based on role
type Router struct {
	mu       sync.RWMutex
	clients  map[string]Client
	cfg      *config.Config
	mockMode bool
}

// NewRouter creates a new LLM router
func NewRouter(cfg *config.Config, mockMode bool) *Router {
	return &Router{
		clients:  make(map[string]Client),
		cfg:      cfg,
		mockMode: mockMode,
	}
}

// GetClient returns the appropriate client for a role
func (r *Router) GetClient(role string) (Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if we already have a cached client
	if client, ok := r.clients[role]; ok {
		return client, nil
	}

	// In mock mode, return a mock client
	if r.mockMode {
		return NewMockClient(), nil
	}

	// Get role configuration
	roleCfg, err := r.cfg.GetRoleConfig(role)
	if err != nil {
		return nil, fmt.Errorf("failed to get role config for %s: %w", role, err)
	}

	// Get provider configuration
	providerCfg, err := r.cfg.GetProvider(roleCfg.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider config for %s: %w", roleCfg.Provider, err)
	}

	// Create client
	client := NewHTTPClient(providerCfg.BaseURL, providerCfg.APIKey, roleCfg.Model)

	// Cache client
	r.clients[role] = client

	return client, nil
}

// CompleteForRole sends a completion request for a specific role
func (r *Router) CompleteForRole(ctx context.Context, role string, req CompletionRequest) (*CompletionResponse, error) {
	client, err := r.GetClient(role)
	if err != nil {
		return nil, err
	}

	// Get role config for default parameters
	if !r.mockMode {
		if roleCfg, err := r.cfg.GetRoleConfig(role); err == nil {
			if req.Temperature == 0 && roleCfg.Temperature > 0 {
				req.Temperature = roleCfg.Temperature
			}
			if req.MaxTokens == 0 && roleCfg.MaxTokens > 0 {
				req.MaxTokens = roleCfg.MaxTokens
			}
			if req.Model == "" {
				req.Model = roleCfg.Model
			}
		}
	}

	return client.Complete(ctx, req)
}

// CompleteWithSystemForRole sends a completion with system prompt for a role
func (r *Router) CompleteWithSystemForRole(ctx context.Context, role, systemPrompt, userPrompt string) (*CompletionResponse, error) {
	client, err := r.GetClient(role)
	if err != nil {
		return nil, err
	}

	// Get role config for parameters
	opts := RequestOptions{}
	if !r.mockMode {
		if roleCfg, err := r.cfg.GetRoleConfig(role); err == nil {
			opts.Temperature = roleCfg.Temperature
			opts.MaxTokens = roleCfg.MaxTokens
		}
	}

	return client.CompleteWithSystem(ctx, systemPrompt, userPrompt, opts)
}

// SetMockClient sets a mock client for a specific role (for testing)
func (r *Router) SetMockClient(role string, client Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[role] = client
}

// CostTracker tracks LLM usage costs
type CostTracker struct {
	mu          sync.RWMutex
	totalCost   float64
	costByRole  map[string]float64
	budgetLimit float64
}

// NewCostTracker creates a new cost tracker
func NewCostTracker(budgetLimit float64) *CostTracker {
	return &CostTracker{
		costByRole:  make(map[string]float64),
		budgetLimit: budgetLimit,
	}
}

// Track records a cost for a role
func (ct *CostTracker) Track(role string, tokensInput, tokensOutput int, model string) float64 {
	// Rough cost estimation (varies by provider)
	cost := estimateCost(tokensInput, tokensOutput, model)

	ct.mu.Lock()
	defer ct.mu.Unlock()

	ct.totalCost += cost
	ct.costByRole[role] += cost

	return ct.totalCost
}

// GetTotal returns the total cost
func (ct *CostTracker) GetTotal() float64 {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	return ct.totalCost
}

// GetByRole returns cost by role
func (ct *CostTracker) GetByRole(role string) float64 {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	return ct.costByRole[role]
}

// IsOverBudget checks if budget is exceeded
func (ct *CostTracker) IsOverBudget() bool {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	if ct.budgetLimit <= 0 {
		return false
	}

	return ct.totalCost >= ct.budgetLimit
}

// IsNearWarning checks if approaching budget
func (ct *CostTracker) IsNearWarning(threshold float64) bool {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	if ct.budgetLimit <= 0 {
		return false
	}

	return ct.totalCost >= ct.budgetLimit*threshold
}

// estimateCost estimates the cost of a request (rough approximation)
func estimateCost(tokensInput, tokensOutput int, model string) float64 {
	// Very rough pricing (per 1K tokens)
	var inputPrice, outputPrice float64

	switch model {
	case "claude-3.5-sonnet", "anthropic/claude-3.5-sonnet":
		inputPrice = 0.003
		outputPrice = 0.015
	case "gpt-4", "openai/gpt-4":
		inputPrice = 0.03
		outputPrice = 0.06
	case "llama-3.1-70b", "meta-llama/llama-3.1-70b":
		inputPrice = 0.0009
		outputPrice = 0.0009
	default:
		// Default to low pricing
		inputPrice = 0.001
		outputPrice = 0.002
	}

	inputCost := float64(tokensInput) * inputPrice / 1000
	outputCost := float64(tokensOutput) * outputPrice / 1000

	return inputCost + outputCost
}
