package saayn

import (
	"fmt"
	"sort"
	"strings"

	"github.com/saayn-agent/internal/astutil"
	"github.com/spf13/cobra"
)

var traceDir string

var traceCmd = &cobra.Command{
	Use:   "trace [PublicID]",
	Short: "Structurally grep the codebase to find all callers of a specific target",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrace,
}

func init() {
	rootCmd.AddCommand(traceCmd)
	traceCmd.Flags().StringVarP(&traceDir, "dir", "d", ".", "Root directory to scan")
}

func runTrace(cmd *cobra.Command, args []string) error {
	targetID := strings.TrimSpace(args[0])
	if targetID == "" {
		return fmt.Errorf("target PublicID cannot be empty")
	}

	fmt.Printf("🎯 Tracing target: %s\n", targetID)

	callers, err := astutil.TraceCallers(traceDir, targetID)
	if err != nil {
		return fmt.Errorf("trace failed: %w\n💡 Hint: Ensure you are using the canonical SAAYN format: pkg.Function[file.go]", err)
	}

	if len(callers) == 0 {
		fmt.Println("   ↳ No callers found. This node is isolated or only called dynamically/via interfaces.")
		return nil
	}

	fmt.Printf("💥 BLAST RADIUS: %d call sites found\n\n", len(callers))

	groupedByFile := make(map[string][]astutil.CallerInfo)
	for _, caller := range callers {
		groupedByFile[caller.FilePath] = append(groupedByFile[caller.FilePath], caller)
	}

	files := make([]string, 0, len(groupedByFile))
	for file := range groupedByFile {
		files = append(files, file)
	}
	sort.Strings(files)

	for _, file := range files {
		calls := groupedByFile[file]

		sort.Slice(calls, func(i, j int) bool {
			if calls[i].LineNumber != calls[j].LineNumber {
				return calls[i].LineNumber < calls[j].LineNumber
			}
			return calls[i].CallingFunction < calls[j].CallingFunction
		})

		fmt.Printf("📄 %s\n", file)
		for _, c := range calls {
			fmt.Printf("   ↳ %s() at line %d\n", c.CallingFunction, c.LineNumber)
		}
		fmt.Println()
	}

	fmt.Println("💡 To stage a surgery including these files, run:")
	fmt.Println("   ./saayn draft \"intent\"")
	fmt.Println("   ./saayn graph")

	return nil
}
