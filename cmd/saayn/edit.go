package saayn

import (
	"fmt"
	"os"

	"github.com/saayn-agent/internal/genome"
	"github.com/saayn-agent/internal/surgeon"
	"github.com/saayn-agent/pkg/model"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Injects new logic into a specific function by UUID",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		targetUUID, _ := cmd.Flags().GetString("uuid")
		newCodeFile, _ := cmd.Flags().GetString("code-file")

		if targetUUID == "" || newCodeFile == "" {
			fmt.Println("❌ Error: --uuid and --code-file are required.")
			return
		}

		fmt.Printf("🏥 Preparing for surgery on UUID: %s...\n", targetUUID)

		// 1. Load the Genome
		regManager, err := genome.Load(path + "/genome.json")
		if err != nil {
			fmt.Println("❌ Error loading genome.json. Cannot locate target.")
			return
		}

		// 2. Read the new code block
		newLogic, err := os.ReadFile(newCodeFile)
		if err != nil {
			fmt.Printf("❌ Error reading new code file: %v\n", err)
			return
		}

		// 3. Prepare the Surgery Request
		req := model.SurgeryRequest{
			TargetUUID:  targetUUID,
			NewLogic:    string(newLogic),
			Registry:    regManager.Registry,
			ProjectRoot: path,
		}

		// 4. Perform the Surgery
		err = surgeon.PerformSurgery(req)
		if err != nil {
			fmt.Printf("❌ Surgery failed: %v\n", err)
			return
		}

		fmt.Println("✅ Surgery successful! The function has been updated.")
		fmt.Println("⚠️  Run './saayn init' to update the genome baseline.")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringP("path", "p", ".", "Path to the project root")
	editCmd.Flags().StringP("uuid", "u", "", "The UUID of the node to edit")
	editCmd.Flags().StringP("code-file", "c", "", "Path to a text file containing the new function body")
}
