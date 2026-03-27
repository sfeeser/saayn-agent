package saayn

import (
	"github.com/spf13/cobra"
)

var (
	genomeFile  string
	projectRoot string
)

var rootCmd = &cobra.Command{
	Use:   "saayn",
	Short: "SAAYN Code Genome System (CGS)",
	Long: `A deterministic, AST-based mutation and repair engine 
that treats code as a living genome.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&genomeFile, "genome", "g", "genome.json", "path to the CGS registry file")
	rootCmd.PersistentFlags().StringVarP(&projectRoot, "path", "p", ".", "root directory of the project to scan")
}
