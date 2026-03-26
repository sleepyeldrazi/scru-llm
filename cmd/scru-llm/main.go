package main

import (
	"fmt"
	"os"

	"github.com/sleepyeldrazi/scru-llm/internal/cli"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scru-llm",
		Short: "Scru-LLM - A specification-first LLM orchestration system",
		Long: `Scru-LLM is an autonomous software engineering system that applies 
Scrum methodology to LLM-driven development.

It takes user tasks through a structured workflow:
  1. Intake - Create initial specification
  2. Spec Sprint - Tighten specification iteratively
  3. Test Sprint - Generate comprehensive tests
  4. Implementation Sprint - Write code to pass tests
  5. Review - Final validation and delivery`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add commands
	rootCmd.AddCommand(cli.NewTaskCommand())
	rootCmd.AddCommand(cli.NewRunCommand())
	rootCmd.AddCommand(cli.NewDeleteCommand())
	rootCmd.AddCommand(cli.NewStatusCommand())
	rootCmd.AddCommand(cli.NewLogsCommand())
	rootCmd.AddCommand(cli.NewConfigCommand())
	rootCmd.AddCommand(cli.NewDashboardCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
