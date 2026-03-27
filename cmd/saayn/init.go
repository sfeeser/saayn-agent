package saayn

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/saayn-agent/internal/scanner"
	"github.com/saayn-agent/internal/genome"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Code Genome for the project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🧬 Initializing Code Genome at: %s\n", projectRoot)
		
		// 1. Perform the AST Scan
		nodes, err := scanner.FullScan(projectRoot)
		if err != nil {
			log.Fatalf("Failed to scan project: %v", err)
		}

		// 2. Build the Registry
		reg := genome.NewRegistry(nodes)
		
		// 3. Persist to genome.json
		if err := reg.Save(genomeFile); err != nil {
			log.Fatalf("Failed to save genome: %v", err)
		}

		fmt.Printf("✅ Success! Indexed %d nodes into %s\n", len(nodes), genomeFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
