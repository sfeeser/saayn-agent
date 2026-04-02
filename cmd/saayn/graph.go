package saayn

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/saayn-agent/internal/genome/surgery"
)

var graphInputFile string

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Analyze the blast radius of a surgery plan (V3 Stub)",
	Long: `In V3, this command mocks the AST graph analysis. 
In V4, it will deterministically find all callers of the target node 
and hydrate them into the plan's context.`,
	RunE: runGraph,
}

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.Flags().StringVarP(&graphInputFile, "file", "f", "surgery.yaml", "Path to the surgery plan YAML")
}

func runGraph(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔍 Analyzing blast radius for %s...\n", graphInputFile)

	yamlBytes, err := os.ReadFile(graphInputFile)
	if err != nil {
		return fmt.Errorf("failed to read plan: %w", err)
	}

	var plan surgery.SurgeryPlan
	if err := yaml.Unmarshal(yamlBytes, &plan); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	if plan.Target.PublicID == "" {
		return fmt.Errorf("invalid surgery plan: missing target.public_id")
	}

	mockAffectedCount := 2
	mockCallers := []string{
		"index.SyncIndex[sync.go]",
		"saayn.runEnrich[enrich.go]",
	}

	if mockAffectedCount > 0 {
		fmt.Printf("\n⚠️  BLAST RADIUS ALERT: Changing %s impacts %d other nodes.\n",
			plan.Target.PublicID, mockAffectedCount)
		fmt.Println("Impacted nodes discovered:")
		for _, c := range mockCallers {
			fmt.Printf("  - %s\n", c)
		}

		fmt.Print("\nWould you like to hydrate these callers into the context? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user confirmation: %w", err)
		}
		input = strings.ToLower(strings.TrimSpace(input))

		if input != "y" {
			fmt.Println("🛑 Analysis aborted by user. Plan remains unchanged.")
			return nil
		}
	}

	plan.ImpactAnalysis = &surgery.ImpactAnalysis{
		RiskLevel:     surgery.RiskHigh,
		RiskReasoning: "STUB (V3): Target is a core utility used by synchronization loops and CLI handlers.",
		RiskFactors: surgery.RiskFactors{
			Visibility:        "public_api",
			BoundaryCrossings: 1,
			TransitiveDepth:   2,
			TotalBlastRadius:  mockAffectedCount,
		},
	}

	mockCallerNodes := []surgery.ContextNode{
		{
			TargetAnchor: surgery.TargetAnchor{
				UUID:      "mock-uuid-999",
				PublicID:  "index.SyncIndex[sync.go]",
				FilePath:  "internal/genome/index/sync.go",
				NodeType:  "function",
				LogicHash: "mock-hash-abc",
			},
			Reason:       surgery.ReasonKnownCaller,
			ReasonDetail: "Invokes FetchEmbedding inside a batch processing loop.",
			SourceCode: `func SyncIndex(path string) error {
	// This is a stubbed caller context
	fmt.Println("Syncing...")
	_, err := FetchEmbedding(ctx, client, key, model, url, "sample")
	return err
}`,
		},
		{
			TargetAnchor: surgery.TargetAnchor{
				UUID:      "mock-uuid-888",
				PublicID:  "saayn.runEnrich[enrich.go]",
				FilePath:  "cmd/saayn/enrich.go",
				NodeType:  "function",
				LogicHash: "mock-hash-def",
			},
			Reason:       surgery.ReasonKnownCaller,
			ReasonDetail: "Uses FetchEmbedding to index files during enrich command.",
			SourceCode: `func runEnrich(cmd *cobra.Command, args []string) error {
	// This is another stubbed caller context
	fmt.Println("Enriching...")
	vec, err := index.FetchEmbedding(ctx, client, key, model, url, content)
	return err
}`,
		},
	}

	hydratedCount := 0
	for _, mockNode := range mockCallerNodes {
		if !hasContextNode(plan, mockNode.PublicID) {
			plan.Context = append(plan.Context, mockNode)
			hydratedCount++
		}
	}

	if hydratedCount > 0 {
		fmt.Printf("💧 Hydrated %d additional caller(s) into the plan.\n", hydratedCount)
	}

	newYaml, err := yaml.Marshal(&plan)
	if err != nil {
		return fmt.Errorf("failed to marshal updated plan: %w", err)
	}

	if err := os.WriteFile(graphInputFile, newYaml, 0644); err != nil {
		return fmt.Errorf("failed to save updated plan: %w", err)
	}

	fmt.Println("✅ Surgery plan successfully updated with impact analysis.")
	return nil
}

func hasContextNode(plan surgery.SurgeryPlan, publicID string) bool {
	for _, node := range plan.Context {
		if node.PublicID == publicID {
			return true
		}
	}
	return false
}
