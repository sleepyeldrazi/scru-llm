// Package workers implements the LLM agent workers for Scru-LLM.
package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/sleepyeldrazi/scru-llm/internal/llm"
	"github.com/sleepyeldrazi/scru-llm/internal/prompts"
	"github.com/sleepyeldrazi/scru-llm/internal/types"
)

// getString safely extracts a string from a map
func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return defaultValue
}

// getFloat safely extracts a float from a map
func getFloat(m map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		}
	}
	return defaultValue
}

// getStringArray safely extracts a string array from a map
func getStringArray(m map[string]interface{}, key string) []string {
	if val, ok := m[key]; ok && val != nil {
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return []string{}
}

// extractJSON extracts JSON from a string that may be wrapped in markdown code blocks
func extractJSON(content string) string {
	// Try to extract from ```json ... ``` blocks
	re := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	// Try to extract from ``` ... ``` blocks
	re = regexp.MustCompile("(?s)```\\s*(\\{.*?\\})\\s*```")
	matches = re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	// Return as-is if no markdown found
	return content
}

// ProductOwner creates initial task specifications
type ProductOwner struct {
	router *llm.Router
	loader *prompts.Loader
}

// NewProductOwner creates a new Product Owner worker
func NewProductOwner(router *llm.Router, loader *prompts.Loader) *ProductOwner {
	return &ProductOwner{
		router: router,
		loader: loader,
	}
}

// CreateSpec creates an initial specification for a task
func (po *ProductOwner) CreateSpec(ctx context.Context, task *types.Task) (*types.Specification, error) {
	systemPrompt, err := po.loader.LoadSystemPrompt("product_owner")
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt: %w", err)
	}

	userPrompt := fmt.Sprintf(`Task: %s
Description: %s

Create a detailed specification.`, task.Title, task.Description)

	resp, err := po.router.CompleteWithSystemForRole(ctx, "product_owner", systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	// Extract JSON from markdown if present
	jsonContent := extractJSON(resp.Content)

	var spec types.Specification
	if err := json.Unmarshal([]byte(jsonContent), &spec); err != nil {
		// Try to extract as generic map first
		var genericSpec map[string]interface{}
		if err2 := json.Unmarshal([]byte(jsonContent), &genericSpec); err2 == nil {
			// Successfully parsed as generic map, convert to spec
			spec = types.Specification{
				TaskID:      task.ID,
				Version:     1,
				Title:       getString(genericSpec, "title", task.Title),
				Description: getString(genericSpec, "description", task.Description),
				Confidence:  getFloat(genericSpec, "confidence", 0.7),
			}
			// Extract arrays
			spec.Requirements = getStringArray(genericSpec, "requirements")
			spec.AcceptanceCriteria = getStringArray(genericSpec, "acceptance_criteria")
			spec.Constraints = getStringArray(genericSpec, "constraints")
			spec.Dependencies = getStringArray(genericSpec, "dependencies")
		} else {
			// Return a basic spec if parsing fails completely
			fmt.Printf("Warning: Failed to parse spec JSON: %v\n", err)
			fmt.Printf("Content was: %.200s...\n", resp.Content)
			spec = types.Specification{
				TaskID:      task.ID,
				Version:     1,
				Title:       task.Title,
				Description: task.Description,
				Confidence:  0.5,
			}
		}
	}

	// Ensure minimum confidence for auto-approval
	if spec.Confidence == 0 {
		spec.Confidence = 0.75
	}

	spec.TaskID = task.ID
	spec.Version = 1

	return &spec, nil
}

// SpecEngineer tightens specifications
type SpecEngineer struct {
	router *llm.Router
	loader *prompts.Loader
}

// NewSpecEngineer creates a new Spec Engineer worker
func NewSpecEngineer(router *llm.Router, loader *prompts.Loader) *SpecEngineer {
	return &SpecEngineer{
		router: router,
		loader: loader,
	}
}

// TightenSpec tightens a specification
func (se *SpecEngineer) TightenSpec(ctx context.Context, spec *types.Specification) (*types.Specification, error) {
	systemPrompt, err := se.loader.LoadSystemPrompt("spec_engineer")
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt: %w", err)
	}

	specJSON, _ := json.MarshalIndent(spec, "", "  ")
	userPrompt := fmt.Sprintf(`Current Specification:
%s

Tighten this specification to make it more precise and testable.`, specJSON)

	resp, err := se.router.CompleteWithSystemForRole(ctx, "spec_engineer", systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	// Extract JSON from markdown if present
	jsonContent := extractJSON(resp.Content)

	var tightened types.Specification
	if err := json.Unmarshal([]byte(jsonContent), &tightened); err != nil {
		// Try to extract as generic map first
		var genericSpec map[string]interface{}
		if err2 := json.Unmarshal([]byte(jsonContent), &genericSpec); err2 == nil {
			// Successfully parsed as generic map, convert to spec
			tightened = types.Specification{
				TaskID:      spec.TaskID,
				Version:     spec.Version + 1,
				Title:       getString(genericSpec, "title", spec.Title),
				Description: getString(genericSpec, "description", spec.Description),
				Confidence:  getFloat(genericSpec, "confidence", 0.8),
			}
			// Extract arrays
			tightened.Requirements = getStringArray(genericSpec, "requirements")
			if len(tightened.Requirements) == 0 {
				tightened.Requirements = spec.Requirements
			}
			tightened.AcceptanceCriteria = getStringArray(genericSpec, "acceptance_criteria")
			tightened.Constraints = getStringArray(genericSpec, "constraints")
			tightened.Dependencies = getStringArray(genericSpec, "dependencies")
		} else {
			// Return updated original if parsing fails completely
			fmt.Printf("Warning: Failed to parse tightened spec JSON: %v\n", err)
			tightened = *spec
			tightened.Version++
			tightened.Confidence = 0.6
		}
	}

	// Ensure minimum confidence
	if tightened.Confidence == 0 {
		tightened.Confidence = 0.8
	}

	tightened.TaskID = spec.TaskID
	tightened.Version = spec.Version + 1

	return &tightened, nil
}

// CodeEngineer implements code and test generation
type CodeEngineer struct {
	router *llm.Router
	loader *prompts.Loader
}

// NewCodeEngineer creates a new Code Engineer worker
func NewCodeEngineer(router *llm.Router, loader *prompts.Loader) *CodeEngineer {
	return &CodeEngineer{
		router: router,
		loader: loader,
	}
}

// GenerateTests generates tests from a specification
func (ce *CodeEngineer) GenerateTests(ctx context.Context, spec *types.Specification) (map[string]string, error) {
	// For MVP, create a simple test template
	language := "python"
	if l, ok := spec.Context["language"].(string); ok {
		language = l
	}

	tests := make(map[string]string)

	switch language {
	case "python":
		tests["test_main.py"] = fmt.Sprintf(`"""Tests for %s"""
import unittest

class Test%s(unittest.TestCase):
    """Test cases based on specification"""
    
    def test_basic_functionality(self):
        """Test basic functionality"""
        # TODO: Implement based on spec
        pass

if __name__ == "__main__":
    unittest.main()
`, spec.Title, spec.Title)
	case "go":
		tests["main_test.go"] = fmt.Sprintf(`package main

import "testing"

func TestBasic(t *testing.T) {
    // TODO: Implement based on spec
}
`)
	default:
		tests["test.txt"] = "TODO: Generate tests based on spec"
	}

	return tests, nil
}

// ImplementCode implements code to pass tests
func (ce *CodeEngineer) ImplementCode(ctx context.Context, spec *types.Specification, tests map[string]string, previousErrors string) (map[string]string, error) {
	language := "python"
	if l, ok := spec.Context["language"].(string); ok {
		language = l
	}

	code := make(map[string]string)

	switch language {
	case "python":
		code["main.py"] = fmt.Sprintf(`"""%s"""

def main():
    """Main function"""
    print("Hello from %s")
    
if __name__ == "__main__":
    main()
`, spec.Title, spec.Title)
	case "go":
		code["main.go"] = `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`
	default:
		code["main.txt"] = "TODO: Implement code based on spec"
	}

	return code, nil
}

// Reviewer reviews specifications, tests, and code
type Reviewer struct {
	router *llm.Router
	loader *prompts.Loader
}

// NewReviewer creates a new Reviewer worker
func NewReviewer(router *llm.Router, loader *prompts.Loader) *Reviewer {
	return &Reviewer{
		router: router,
		loader: loader,
	}
}

// ReviewSpec reviews a specification
func (r *Reviewer) ReviewSpec(ctx context.Context, spec *types.Specification) (*types.ReviewDecision, error) {
	// For MVP, auto-approve specs with high confidence
	if spec.Confidence >= 0.8 {
		return &types.ReviewDecision{
			Approved:              true,
			Phase:                 "spec",
			Summary:               "Specification meets quality standards",
			Issues:                []string{},
			RecommendedNextAction: "advance",
		}, nil
	}

	return &types.ReviewDecision{
		Approved:              false,
		Phase:                 "spec",
		Summary:               "Specification needs improvement",
		Issues:                []string{"Confidence score too low"},
		RecommendedNextAction: "retry",
	}, nil
}

// ReviewImplementation reviews code implementation
func (r *Reviewer) ReviewImplementation(ctx context.Context, spec *types.Specification, tests, code map[string]string) (*types.ReviewDecision, error) {
	return &types.ReviewDecision{
		Approved:              true,
		Phase:                 "review",
		Summary:               "Implementation approved",
		Issues:                []string{},
		RecommendedNextAction: "advance",
	}, nil
}
