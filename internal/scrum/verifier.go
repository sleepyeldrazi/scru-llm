package scrum

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Verifier runs tests and verifies code
type Verifier struct {
	workspaceDir string
}

// NewVerifier creates a new verifier
func NewVerifier(workspaceDir string) *Verifier {
	return &Verifier{
		workspaceDir: workspaceDir,
	}
}

// VerificationResult contains the results of running tests
type VerificationResult struct {
	TestsRun    int
	TestsPassed int
	TestsFailed int
	Output      string
	Success     bool
}

// RunTests runs tests for a task
func (v *Verifier) RunTests(taskID string) (*VerificationResult, error) {
	taskDir := filepath.Join(v.workspaceDir, taskID)
	srcDir := filepath.Join(taskDir, "src")
	testDir := filepath.Join(taskDir, "tests")

	// Detect language
	language := v.detectLanguage(srcDir, testDir)

	switch language {
	case "python":
		return v.runPythonTests(taskDir, srcDir, testDir)
	case "go":
		return v.runGoTests(taskDir, srcDir, testDir)
	default:
		return &VerificationResult{
			TestsRun:    0,
			TestsPassed: 0,
			TestsFailed: 0,
			Output:      fmt.Sprintf("Unknown language: %s", language),
			Success:     false,
		}, nil
	}
}

// detectLanguage detects the programming language from files
func (v *Verifier) detectLanguage(srcDir, testDir string) string {
	// Check for Python files
	if _, err := os.Stat(filepath.Join(srcDir, "main.py")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(testDir, "test_main.py")); err == nil {
		return "python"
	}

	// Check for Go files
	if _, err := os.Stat(filepath.Join(srcDir, "main.go")); err == nil {
		return "go"
	}
	if _, err := os.Stat(filepath.Join(testDir, "main_test.go")); err == nil {
		return "go"
	}

	return "unknown"
}

// runPythonTests runs Python tests
func (v *Verifier) runPythonTests(taskDir, srcDir, testDir string) (*VerificationResult, error) {
	// Find test files
	entries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read test directory: %w", err)
	}

	var testFiles []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "test_") && strings.HasSuffix(entry.Name(), ".py") {
			testFiles = append(testFiles, filepath.Join(testDir, entry.Name()))
		}
	}

	if len(testFiles) == 0 {
		return &VerificationResult{
			TestsRun:    0,
			TestsPassed: 0,
			TestsFailed: 0,
			Output:      "No test files found",
			Success:     true,
		}, nil
	}

	// Run tests with pytest or unittest
	cmd := exec.Command("python", append([]string{"-m", "pytest", "-v"}, testFiles...)...)
	cmd.Dir = taskDir
	cmd.Env = append(os.Environ(), "PYTHONPATH="+srcDir)

	output, err := cmd.CombinedOutput()

	result := &VerificationResult{
		Output: string(output),
	}

	// Parse output for test counts
	outputStr := string(output)
	if strings.Contains(outputStr, "passed") {
		// Try to parse pytest output
		var passed, failed int
		fmt.Sscanf(outputStr, "%d passed", &passed)
		fmt.Sscanf(outputStr, "%d failed", &failed)
		result.TestsPassed = passed
		result.TestsFailed = failed
		result.TestsRun = passed + failed
		result.Success = failed == 0
	} else {
		// Check for success/failure
		result.Success = err == nil
	}

	return result, nil
}

// runGoTests runs Go tests
func (v *Verifier) runGoTests(taskDir, srcDir, testDir string) (*VerificationResult, error) {
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = taskDir

	output, err := cmd.CombinedOutput()

	result := &VerificationResult{
		Output: string(output),
	}

	// Parse output
	outputStr := string(output)
	if strings.Contains(outputStr, "PASS") || strings.Contains(outputStr, "FAIL") {
		// Count test results
		passCount := strings.Count(outputStr, "PASS")
		failCount := strings.Count(outputStr, "FAIL")

		result.TestsPassed = passCount
		result.TestsFailed = failCount
		result.TestsRun = passCount + failCount
		result.Success = failCount == 0 && err == nil
	} else {
		result.Success = err == nil
	}

	return result, nil
}

// VerifySyntax verifies code syntax without running tests
func (v *Verifier) VerifySyntax(taskID string) error {
	taskDir := filepath.Join(v.workspaceDir, taskID)
	srcDir := filepath.Join(taskDir, "src")

	language := v.detectLanguage(srcDir, "")

	switch language {
	case "python":
		return v.verifyPythonSyntax(srcDir)
	case "go":
		return v.verifyGoSyntax(taskDir)
	default:
		return fmt.Errorf("cannot verify syntax for unknown language")
	}
}

// verifyPythonSyntax checks Python syntax
func (v *Verifier) verifyPythonSyntax(srcDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".py") {
			continue
		}

		path := filepath.Join(srcDir, entry.Name())
		cmd := exec.Command("python", "-m", "py_compile", path)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("syntax error in %s: %s", entry.Name(), output)
		}
	}

	return nil
}

// verifyGoSyntax checks Go syntax
func (v *Verifier) verifyGoSyntax(taskDir string) error {
	cmd := exec.Command("go", "build", "-o", "/dev/null", ".")
	cmd.Dir = taskDir

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build error: %s", output)
	}

	return nil
}
