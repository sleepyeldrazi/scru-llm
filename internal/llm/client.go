// Package llm provides LLM client abstractions and routing.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest represents a request to the LLM
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// CompletionResponse represents a response from the LLM
type CompletionResponse struct {
	Content      string `json:"content"`
	TokensUsed   int    `json:"tokens_used"`
	TokensInput  int    `json:"tokens_input"`
	TokensOutput int    `json:"tokens_output"`
	Model        string `json:"model"`
}

// Client defines the interface for LLM clients
type Client interface {
	// Complete sends a completion request and returns the response
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	// CompleteWithSystem sends a completion with a system message
	CompleteWithSystem(ctx context.Context, systemPrompt, userPrompt string, opts RequestOptions) (*CompletionResponse, error)
}

// RequestOptions contains optional parameters for requests
type RequestOptions struct {
	Temperature float64
	MaxTokens   int
}

// HTTPClient implements Client using HTTP requests
type HTTPClient struct {
	baseURL      string
	apiKey       string
	defaultModel string
	httpClient   *http.Client
}

// NewHTTPClient creates a new HTTP-based LLM client
func NewHTTPClient(baseURL, apiKey, defaultModel string) *HTTPClient {
	return &HTTPClient{
		baseURL:      baseURL,
		apiKey:       apiKey,
		defaultModel: defaultModel,
		httpClient:   &http.Client{Timeout: 300 * time.Second},
	}
}

// SetTimeout sets the HTTP client timeout
func (c *HTTPClient) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// openAIRequest represents the OpenAI API request format
type openAIRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents the OpenAI API response format
type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Complete sends a completion request
func (c *HTTPClient) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	model := req.Model
	if model == "" {
		model = c.defaultModel
	}

	// Calculate input tokens (approximate)
	inputChars := 0
	for _, m := range req.Messages {
		inputChars += len(m.Content)
	}
	inputTokens := inputChars / 4 // Rough approximation

	fmt.Printf("[LLM CALL] Model: %s | Input: ~%d tokens | Temperature: %.2f\n",
		model, inputTokens, req.Temperature)

	// Build OpenAI-compatible request
	messages := make([]message, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = message{Role: m.Role, Content: m.Content}
	}

	apiReq := openAIRequest{
		Model:       model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	// Marshal request
	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	start := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	elapsed := time.Since(start)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract content
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// Log response
	outputTokens := apiResp.Usage.CompletionTokens
	totalTokens := apiResp.Usage.TotalTokens
	fmt.Printf("[LLM RESP] Model: %s | Time: %s | Input: %d | Output: %d | Total: %d\n",
		apiResp.Model, elapsed, apiResp.Usage.PromptTokens, outputTokens, totalTokens)

	return &CompletionResponse{
		Content:      apiResp.Choices[0].Message.Content,
		TokensUsed:   totalTokens,
		TokensInput:  apiResp.Usage.PromptTokens,
		TokensOutput: outputTokens,
		Model:        apiResp.Model,
	}, nil
}

// CompleteWithSystem sends a completion with system and user prompts
func (c *HTTPClient) CompleteWithSystem(ctx context.Context, systemPrompt, userPrompt string, opts RequestOptions) (*CompletionResponse, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	req := CompletionRequest{
		Messages:    messages,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
	}

	return c.Complete(ctx, req)
}

// MockClient is a mock implementation for testing
type MockClient struct {
	Responses map[string]string
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		Responses: make(map[string]string),
	}
}

// Complete returns a mock response
func (m *MockClient) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Use user message content as key
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages in request")
	}

	userContent := req.Messages[len(req.Messages)-1].Content
	if resp, ok := m.Responses[userContent]; ok {
		return &CompletionResponse{
			Content:      resp,
			TokensUsed:   100,
			TokensInput:  50,
			TokensOutput: 50,
			Model:        "mock-model",
		}, nil
	}

	return &CompletionResponse{
		Content:      `{"title": "Mock Response", "description": "This is a mock response"}`,
		TokensUsed:   50,
		TokensInput:  25,
		TokensOutput: 25,
		Model:        "mock-model",
	}, nil
}

// CompleteWithSystem returns a mock response
func (m *MockClient) CompleteWithSystem(ctx context.Context, systemPrompt, userPrompt string, opts RequestOptions) (*CompletionResponse, error) {
	return m.Complete(ctx, CompletionRequest{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	})
}

// SetResponse sets a mock response for a specific prompt
func (m *MockClient) SetResponse(prompt, response string) {
	m.Responses[prompt] = response
}
