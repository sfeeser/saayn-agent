package saayn

import (
	"fmt"

	"github.com/sfeeser/saayn-agent/internal/surgeon"
	"github.com/spf13/cobra"
)

var genesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Materialize a Greenfield project from a Vision and Specbook",
	Long: `Genesis initializes a new SAAYN-managed codebase using the Greenfield Protocol.

It reads a Markdown vision document (the Soul) and a YAML specbook (the Skeleton),
calculates the dependency graph, and executes the Surgical Inner Loop to
synthesize the complete project.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		vision, err := cmd.Flags().GetString("vision")
		if err != nil {
			return fmt.Errorf("failed to read vision flag: %w", err)
		}
		spec, err := cmd.Flags().GetString("spec")
		if err != nil {
			return fmt.Errorf("failed to read spec flag: %w", err)
		}
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return fmt.Errorf("failed to read target flag: %w", err)
		}

		fmt.Printf("🧬 GREENFIELD PROTOCOL GENESIS INITIATED\n")
		fmt.Printf("   Vision: %s\n", vision)
		fmt.Printf("   Spec:   %s\n", spec)
		fmt.Printf("   Target: %s\n", target)
		fmt.Println("--------------------------------------------------")

		// Hand off to the core Genesis Engine
		if err := surgeon.ExecuteGenesis(vision, spec, target); err != nil {
			return fmt.Errorf("genesis failed: %w", err)
		}

		fmt.Printf("\n✅ GENESIS COMPLETE — The genome has been successfully materialized.\n")
		return nil
	},
}

func init() {
	genesisCmd.Flags().StringP("vision", "v", "vision.md", "Path to the Markdown vision document (the Soul)")
	genesisCmd.Flags().StringP("spec", "s", "specbook.yaml", "Path to the YAML specbook (the Skeleton)")
	genesisCmd.Flags().StringP("target", "t", "./generated", "Target directory for the materialized project")

	// Attach this command to the root command
	rootCmd.AddCommand(genesisCmd)
}
