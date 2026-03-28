package saayn

import (
	"fmt"

	"github.com/saayn-agent/internal/genome"
	"github.com/saayn-agent/internal/scanner"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Detects logic drift between live code and the genome.json",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")

		fmt.Println("🔍 Verifying Code Genome for drift...")

		// 1. Load the baseline genome
		oldRegManager, err := genome.Load(path + "/genome.json")
		if err != nil {
			fmt.Println("❌ Error loading genome.json. Did you run 'init' first?")
			return
		}

		// 2. Scan the live code
		liveNodes, err := scanner.FullScan(path)
		if err != nil {
			fmt.Printf("❌ Error scanning directory: %v\n", err)
			return
		}

		// 3. Process the live nodes through our hashing engine (but don't save to disk)
		liveRegManager := genome.NewRegistry(liveNodes, "")

		// 4. Map the old hashes by PublicID so we can compare them
		oldByPubID := make(map[string]string)
		for _, n := range oldRegManager.Registry.Nodes {
			oldByPubID[n.PublicID] = n.LogicHash
		}

		changed := 0
		newNodes := 0

		fmt.Println("\n📊 Drift Report:")
		fmt.Println("------------------------------------------------------------")

		// 5. Compare Live against Baseline
		for _, ln := range liveRegManager.Registry.Nodes {
			oldHash, exists := oldByPubID[ln.PublicID]

			if !exists {
				fmt.Printf("  ✨ NEW      %-45s [%s]\n", ln.PublicID, ln.LogicHash[:8])
				newNodes++
			} else if oldHash != ln.LogicHash {
				fmt.Printf("  ⚠️  MODIFIED %-45s [%s -> %s]\n", ln.PublicID, oldHash[:8], ln.LogicHash[:8])
				changed++
			}

			// Remove from the map to track what was deleted
			delete(oldByPubID, ln.PublicID)
		}

		// 6. Anything left in oldByPubID was deleted from the live code
		deleted := len(oldByPubID)
		for pubID := range oldByPubID {
			fmt.Printf("  🗑️  DELETED  %-45s\n", pubID)
		}

		fmt.Println("------------------------------------------------------------")
		if changed == 0 && newNodes == 0 && deleted == 0 {
			fmt.Println("✅ Genome match! No logic drift detected.")
		} else {
			fmt.Printf("⚠️  Drift detected: %d Modified, %d New, %d Deleted\n", changed, newNodes, deleted)
		}
	},
}

func init() {
	fmt.
		Println("SAAYN Deep Engine Initialized")
	rootCmd.
		AddCommand(verifyCmd)
}
