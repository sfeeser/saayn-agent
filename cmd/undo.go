package cmd

// SAAYN:CHUNK_START:undo-imports-v1-u1n2d3o4
// BUSINESS_PURPOSE: Imports for executing system git commands and managing the internal operation tags.
// SPEC_LINK: SpecBook v1.7 Chapter 8 & 9
import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
)
// SAAYN:CHUNK_END:undo-imports-v1-u1n2d3o4

// SAAYN:CHUNK_START:undo-command-definition-v1-t5i6m7e8
// BUSINESS_PURPOSE: Defines the 'undo' command. Allows targeting a specific operation ID or defaulting to the most recent SAAYN edit.
// SPEC_LINK: SpecBook v1.7 Chapter 9
var targetOp string

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Revert the last SAAYN operation (code and registry)",
	Run: func(cmd *cobra.Command, args []string) {
		runUndo()
	},
}

func init() {
	undoCmd.Flags().StringVarP(&targetOp, "op", "o", "", "Specific Operation ID to revert (e.g., op-12345)")
	rootCmd.AddCommand(undoCmd)
}
// SAAYN:CHUNK_END:undo-command-definition-v1-t5i6m7e8

// SAAYN:CHUNK_START:undo-logic-v1-r1e2v3e4
// BUSINESS_PURPOSE: Uses git tags to perform a clean revert of both the source files and the chunk-registry.json.
// SPEC_LINK: SpecBook v1.7 Chapter 8 & 10 (Law 6)
func runUndo() {
	// 1. Identify the tag to revert
	tagToRevert := targetOp
	if tagToRevert == "" {
		// Fetch the latest SAAYN_OP tag from git
		out, err := exec.Command("git", "describe", "--tags", "--match", "SAAYN_OP:*", "--abbrev=0").Output()
		if err != nil {
			fmt.Println("❌ No previous SAAYN operations found in git history.")
			return
		}
		tagToRevert = strings.TrimSpace(string(out))
	} else if !strings.HasPrefix(tagToRevert, "SAAYN_OP:") {
		tagToRevert = "SAAYN_OP:" + tagToRevert
	}

	fmt.Printf("⏳ Reverting operation: %s...\n", tagToRevert)

	// 2. Perform the git revert
	// We use 'git revert -n' (no-commit) so the user can inspect before finalizing,
	// or just 'git revert' to keep the history clean.
	cmd := exec.Command("git", "revert", tagToRevert, "--no-edit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Git revert failed: %v\n", err)
		fmt.Println("⚠️  Ensure your working directory is clean before running undo.")
		return
	}

	// 3. Post-Undo Verification
	fmt.Printf("✅ %s reverted. Running integrity check...\n", tagToRevert)
	runVerify() // Call the existing verify logic to ensure the registry is back in sync
}
// SAAYN:CHUNK_END:undo-logic-v1-r1e2v3e4
