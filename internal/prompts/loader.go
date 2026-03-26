// Package prompts provides prompt loading and management for LLM agents.
package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Loader handles loading and caching of prompt files
type Loader struct {
	mu         sync.RWMutex
	promptsDir string
	cache      map[string]string
}

// NewLoader creates a new prompt loader
func NewLoader(promptsDir string) *Loader {
	return &Loader{
		promptsDir: promptsDir,
		cache:      make(map[string]string),
	}
}

// LoadSystemPrompt loads a system prompt for a role
func (l *Loader) LoadSystemPrompt(role string) (string, error) {
	// Map role names to file names
	fileName := role + "-system.md"

	// Check cache first
	l.mu.RLock()
	if cached, ok := l.cache[fileName]; ok {
		l.mu.RUnlock()
		return cached, nil
	}
	l.mu.RUnlock()

	// Load from file
	path := filepath.Join(l.promptsDir, fileName)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return a default system prompt
			defaultPrompt := l.defaultSystemPrompt(role)
			return defaultPrompt, nil
		}
		return "", fmt.Errorf("failed to read prompt file %s: %w", path, err)
	}

	prompt := string(content)

	// Cache it
	l.mu.Lock()
	l.cache[fileName] = prompt
	l.mu.Unlock()

	return prompt, nil
}

// LoadPrompt loads any prompt file by name
func (l *Loader) LoadPrompt(name string) (string, error) {
	// Check cache
	l.mu.RLock()
	if cached, ok := l.cache[name]; ok {
		l.mu.RUnlock()
		return cached, nil
	}
	l.mu.RUnlock()

	// Load from file
	path := filepath.Join(l.promptsDir, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", path, err)
	}

	prompt := string(content)

	// Cache it
	l.mu.Lock()
	l.cache[name] = prompt
	l.mu.Unlock()

	return prompt, nil
}

// ClearCache clears the prompt cache
func (l *Loader) ClearCache() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[string]string)
}

// ListAvailable lists all available prompt files
func (l *Loader) ListAvailable() ([]string, error) {
	entries, err := os.ReadDir(l.promptsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	var prompts []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".txt") {
			prompts = append(prompts, name)
		}
	}

	return prompts, nil
}

// defaultSystemPrompt returns a default system prompt for a role
func (l *Loader) defaultSystemPrompt(role string) string {
	switch role {
	case "product_owner":
		return `You are the Product Owner in an LLM-driven Scrum system.
Your responsibility is to analyze user tasks and create clear, actionable specifications.

Transform user tasks into detailed specifications with:
- Title and description
- Specific, testable requirements
- Clear acceptance criteria
- Constraints and dependencies
- Complexity and risk assessment

Output valid JSON with the specification structure.`

	case "spec_engineer":
		return `You are the Spec Engineer in an LLM-driven Scrum system.
Your responsibility is to tighten specifications by removing ambiguity.

Make specifications precise and testable:
- Convert vague terms to measurable values
- Define all terms clearly
- Document edge cases
- Ensure acceptance criteria are testable
- Provide confidence score

Output valid JSON with the tightened specification.`

	case "code_engineer":
		return `You are the Code Engineer in an LLM-driven Scrum system.
Your responsibility is to implement code that passes tests.

Write clean, working code:
- Follow language idioms and best practices
- Pass all tests
- Include appropriate error handling
- Add comments for complex logic
- Follow existing code style`

	case "reviewer":
		return `You are the Reviewer in an LLM-driven Scrum system.
Your responsibility is to validate outputs against specifications.

Review thoroughly and provide:
- Pass/Fail decision
- Specific issues found
- Recommended next action
- Clear feedback

Output valid JSON with review decision.`

	default:
		return fmt.Sprintf("You are the %s in an LLM-driven Scrum system. Follow best practices for your role.", role)
	}
}
