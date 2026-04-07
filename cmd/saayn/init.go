package saayn

import (
	"fmt"

	"github.com/sfeeser/saayn-agent/internal/genome"
	"github.com/sfeeser/saayn-agent/internal/scanner"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Code Genome for the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🧬 Initializing Code Genome at: %s\n", projectRoot)

		// 1. Perform the AST Scan (The Identity Factory)
		nodes, err := scanner.FullScan(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to scan project: %w", err)
		}

		if len(nodes) == 0 {
			return fmt.Errorf("no Go nodes found in %s. Is this a Go project?", projectRoot)
		}

		// 2. Build the Registry
		reg := genome.NewRegistry(nodes, genomeFile)

		// 3. Persist to genome.json
		if err := reg.Save(); err != nil {
			return fmt.Errorf("failed to save genome: %w", err)
		}

		fmt.Printf("✅ Success! Indexed %d nodes into %s\n", len(nodes), genomeFile)
		fmt.Println("💡 Next step: Run './saayn enrich' to generate semantic summaries.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
