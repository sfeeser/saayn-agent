package cmd

// SAAYN:CHUNK_START:reconcile-imports-v1-a1b2c3d4
// BUSINESS_PURPOSE: Imports for terminal I/O, registry management, and hashing.
// SPEC_LINK: SpecBook v1.7 Chapter 9
import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"saayn/internal/adapter"
	"saayn/internal/registry"
	"github.com/spf13/cobra"
)
// SAAYN:CHUNK_END:reconcile-imports-v1-a1b2c3d4

// SAAYN:CHUNK_START:reconcile-command-definition-v1-e5f6g7h8
// BUSINESS_PURPOSE: Defines the 'reconcile' command which provides a UI for resolving cryptographic drift.
// SPEC_LINK: SpecBook v1.7 Chapter 9
var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Interactively update the registry to match manual code changes",
	Run: func(cmd *cobra.Command, args []string) {
		runReconcile()
	},
}

func init() {
	rootCmd.AddCommand(reconcileCmd)
}
// SAAYN:CHUNK_END:reconcile-command-definition-v1-e5f6g7h8

// SAAYN:CHUNK_START:reconcile-logic-v1-i9j0k1l2
// BUSINESS_PURPOSE: Implements the reconciliation loop. Scans for drift and prompts the user to accept new hashes for manual edits.
// SPEC_LINK: SpecBook v1.7 Chapter 5 (Manual Approval) & 7
func runReconcile() {
	reg := loadRegistry()
	reader := bufio.NewReader(os.Stdin)
	updated := false

	fmt.Println("🔄 Scanning for manual drift to reconcile...")

	for i, chunk := range reg.Chunks {
		content, err := os.ReadFile(chunk.FilePath)
		if err != nil {
			continue // Handled by verify, skip here
		}

		adp, _ := adapter.Get(chunk.LanguageHint)
		extracted, startLine, endLine, err := extractChunk(string(content), chunk.UUID, adp)
		if err != nil {
			continue // Marker corruption requires 'heal', skip reconcile
		}

		newContentHash := registry.ComputeContentHash(extracted)
		newMarkerHash := registry.ComputeMarkerHash(startLine, endLine)

		// Check for drift
		if newContentHash != chunk.ContentHash || newMarkerHash != chunk.MarkerHash {
			fmt.Printf("\n⚠️  Drift detected in chunk: %s (%s)\n", chunk.UUID, chunk.FilePath)
			fmt.Print("   Do you want to update the registry to match the current file state? (y/N): ")
			
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) == "y" {
				// Update registry entry
				reg.Chunks[i].ContentHash = newContentHash
				reg.Chunks[i].MarkerHash = newMarkerHash
				reg.Chunks[i].Version++
				reg.Chunks[i].LastModified = time.Now()
				updated = true
				fmt.Printf("   ✅ Registry updated to v%d for %s\n", reg.Chunks[i].Version, chunk.UUID)
			}
		}
	}

	if updated {
		saveRegistry(reg)
		fmt.Println("\n💾 chunk-registry.json has been synchronized.")
	} else {
		fmt.Println("\n✨ No manual drift requires reconciliation.")
	}
}
// SAAYN:CHUNK_END:reconcile-logic-v1-i9j0k1l2
