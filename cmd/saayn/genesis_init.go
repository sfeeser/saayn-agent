package saayn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var genesisInitCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Generate instructional templates for a new Greenfield project",
	Long: `Scaffolds vision.md and specbook.yaml files required by the SAAYN Greenfield Protocol.

The templates contain embedded AI instructions designed to be copy-pasted into 
an AI assistant (Grok, Claude, Gemini, etc.) for collaborative design.`,

	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		targetDir, err := cmd.Flags().GetString("target")
		if err != nil {
			return fmt.Errorf("failed to read target flag: %w", err)
		}

		// Create target directory
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
		}

		// Write vision.md
		visionPath := filepath.Join(targetDir, "vision.md")
		if err := os.WriteFile(visionPath, []byte(visionTemplate(projectName)), 0644); err != nil {
			return fmt.Errorf("failed to write vision.md: %w", err)
		}

		// Write specbook.yaml
		specPath := filepath.Join(targetDir, "specbook.yaml")
		if err := os.WriteFile(specPath, []byte(specTemplate(projectName)), 0644); err != nil {
			return fmt.Errorf("failed to write specbook.yaml: %w", err)
		}

		fmt.Printf("✅ Greenfield templates successfully generated for project: %s\n", projectName)
		fmt.Printf("   Directory: %s\n", targetDir)
		fmt.Printf("   Created:   %s\n", visionPath)
		fmt.Printf("   Created:   %s\n\n", specPath)

		fmt.Println("🚀 NEXT STEPS:")
		fmt.Println("  1. Open vision.md and specbook.yaml")
		fmt.Println("  2. Copy the contents of vision.md into your AI assistant")
		fmt.Println("  3. Collaboratively refine both files with the AI")
		fmt.Println("  4. Run 'saayn genesis vet' to check quality")
		fmt.Println("  5. Run 'saayn genesis' to materialize the project")

		return nil
	},
}

func init() {
	genesisInitCmd.Flags().StringP("target", "t", ".", "Directory where templates will be generated")
	genesisCmd.AddCommand(genesisInitCmd)
}

// =============================================================================
// Templates
// =============================================================================

func visionTemplate(projectName string) string {
	return `# Project Vision: ` + projectName + `

> **AI INSTRUCTIONS (System Prompt):**
> You are helping me build a Go application using the SAAYN Greenfield Protocol.
> This document is the "Soul" of the application. It will be read directly by an autonomous synthesis engine.
>
> **Your Role:** Help me complete this document. Do NOT generate any code yet.
> **Strict Rules:**
> 1. Ask one clarifying question at a time when something is vague.
> 2. Push for concrete details on storage, error handling, and edge cases.
> 3. Enforce strict Go idioms and prefer the standard library.
> 4. Never leave placeholders like [...], TODO, or vague language in the final version.

## 1. Core Objective
(What does this software do in 2–4 clear sentences?)

## 2. Technical Physics
- **Go Version:** 1.22+
- **Architecture:** [CLI tool, HTTP API, background daemon, library, etc.]
- **State Management:** [in-memory, JSON files, SQLite, PostgreSQL, etc.]
- **Allowed Dependencies:** [standard library only, or list specific ones]
- **Concurrency:** (if relevant)

## 3. Core Entities (Nouns)
List the main data structures with their important fields and purpose.

## 4. Primary Operations (Verbs)
List the key actions the system must perform.

## 5. Edge Cases & Failure Modes
Describe important boundary conditions and expected behavior.

## 6. Success Criteria
How will we know the project is successful?
`
}

func specTemplate(projectName string) string {
	return `# SAAYN Greenfield Specbook
# This is the "Skeleton" that controls the exact build order and node logic.

project: "` + projectName + `"
version: "v0.1.0"

# =============================================================================
# GENOME NODES
# Rules:
#   name:        Exact Go identifier (must be unique)
#   path:        Relative file path inside the project
#   type:        "struct", "func", "interface", or "method"
#   receiver:    Only needed if type is "method"
#   depends_on:  List of node names that must be created first
#   logic:       Precise instruction used by the LLM when generating this node
# =============================================================================

nodes:
  - name: "Config"
    path: "internal/config/config.go"
    type: "struct"
    depends_on: []
    logic: "Configuration struct with fields for application settings. Support environment variables and validation."

  - name: "main"
    path: "cmd/` + projectName + `/main.go"
    type: "func"
    depends_on: ["Config"]
    logic: "Application entry point. Load configuration and start the main logic."
`
}
