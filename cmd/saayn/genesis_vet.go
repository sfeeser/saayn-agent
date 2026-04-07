package saayn

import (
	"fmt"
	"os"
	"strings"

	"github.com/sfeeser/saayn-agent/internal/surgeon"
	"github.com/sfeeser/saayn-agent/pkg/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var genesisVetCmd = &cobra.Command{
	Use:   "vet",
	Short: "Validate the Vision and Specbook before materialization",
	Long: `Runs a pre-flight purity audit on your Greenfield inputs.
Checks for unresolved placeholders in vision.md, validates specbook.yaml syntax,
and ensures the dependency graph has no circular dependencies.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		visionPath, err := cmd.Flags().GetString("vision")
		if err != nil {
			return fmt.Errorf("failed to read vision flag: %w", err)
		}
		specPath, err := cmd.Flags().GetString("spec")
		if err != nil {
			return fmt.Errorf("failed to read spec flag: %w", err)
		}

		fmt.Println("🔍 RUNNING GREENFIELD PURITY AUDIT...")
		fmt.Println("--------------------------------------------------")

		// 1. Ghost Check - Vision Document
		fmt.Printf("[1/3] Auditing Soul (vision.md)... ")
		if err := auditVision(visionPath); err != nil {
			fmt.Println("❌ FAILED")
			return err
		}
		fmt.Println("✅ PASSED")

		// 2. YAML Physics Check
		fmt.Printf("[2/3] Auditing Skeleton Syntax (specbook.yaml)... ")
		spec, err := auditSpecbook(specPath)
		if err != nil {
			fmt.Println("❌ FAILED")
			return err
		}
		fmt.Println("✅ PASSED")

		// 3. Topological Graph Audit
		fmt.Printf("[3/3] Auditing Dependency Graph... ")
		if _, err := surgeon.CalculateBuildOrder(spec); err != nil {
			fmt.Println("❌ FAILED")
			return fmt.Errorf("dependency graph error: %w", err)
		}
		fmt.Println("✅ PASSED")

		fmt.Println("--------------------------------------------------")
		fmt.Println("🟢 AUDIT COMPLETE: The Genome is clean and viable.")
		fmt.Println("   You may now run: saayn genesis")

		return nil
	},
}

func init() {
	genesisVetCmd.Flags().StringP("vision", "v", "vision.md", "Path to the Markdown vision document")
	genesisVetCmd.Flags().StringP("spec", "s", "specbook.yaml", "Path to the YAML specbook")
	genesisCmd.AddCommand(genesisVetCmd)
}

// -----------------------------------------------------------------------------
// Audit Helpers
// -----------------------------------------------------------------------------

func auditVision(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read vision document %s: %w", path, err)
	}

	text := strings.TrimSpace(string(data))
	if len(text) == 0 {
		return fmt.Errorf("vision document is empty")
	}

	// More comprehensive placeholder detection
	placeholders := []string{"[...]", "[TODO]", "TODO:", "fill in", "[What does", "[Describe"}
	for _, p := range placeholders {
		if strings.Contains(strings.ToLower(text), strings.ToLower(p)) {
			return fmt.Errorf("vision document contains unresolved placeholder: %s", p)
		}
	}

	return nil
}

func auditSpecbook(path string) (*model.Specbook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read specbook %s: %w", path, err)
	}

	var spec model.Specbook
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("invalid YAML in specbook: %w", err)
	}

	if spec.Project == "" {
		return nil, fmt.Errorf("specbook is missing required field: project")
	}
	if len(spec.Nodes) == 0 {
		return nil, fmt.Errorf("specbook contains no nodes")
	}

	// Optional: Additional validation (duplicate names, empty logic, etc.)
	seen := make(map[string]bool)
	for _, node := range spec.Nodes {
		if node.Name == "" {
			return nil, fmt.Errorf("node missing required field: name")
		}
		if seen[node.Name] {
			return nil, fmt.Errorf("duplicate node name: %s", node.Name)
		}
		seen[node.Name] = true

		if node.Logic == "" {
			return nil, fmt.Errorf("node %s has empty logic field", node.Name)
		}
	}

	return &spec, nil
}
