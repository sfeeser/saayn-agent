package saayn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/saayn-agent/internal/genome"
	"github.com/saayn-agent/internal/genome/index"
	"github.com/saayn-agent/internal/genome/surgery"
)

var draftOutputFile string

var draftCmd = &cobra.Command{
	Use:   "draft [intent]",
	Short: "Initialize a code surgery plan based on a natural language intent",
	Args:  cobra.ExactArgs(1),
	RunE:  runDraft,
}

func init() {
	rootCmd.AddCommand(draftCmd)
	draftCmd.Flags().StringVarP(&draftOutputFile, "output", "o", "surgery.yaml", "Output path for the surgery plan")
}

func runDraft(cmd *cobra.Command, args []string) error {
	// 1. Validate Intent
	intent := strings.TrimSpace(args[0])
	if intent == "" {
		return fmt.Errorf("draft intent cannot be empty")
	}

	fmt.Printf("🩺 Initializing surgery draft for intent: %q\n", intent)

	// 2. Load Configuration
	cfg, err := loadDraftConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to load search config: %w", err)
	}

	// 3. Load the Memory Systems
	store, err := index.LoadIndex(filepath.Join(cfg.Path, "genome.index.json"))
	if err != nil {
		return fmt.Errorf("failed to load semantic index. Have you run 'saayn enrich'? %w", err)
	}

	if store.Metadata.EmbeddingModel != "" && store.Metadata.EmbeddingModel != cfg.EmbedModel {
		return fmt.Errorf("index model mismatch: store uses %s, but config specifies %s", store.Metadata.EmbeddingModel, cfg.EmbedModel)
	}

	registryPath := filepath.Join(cfg.Path, "genome.json")
	registryManager, err := genome.Load(registryPath)
	if err != nil {
		return fmt.Errorf("failed to load AST registry: %w", err)
	}

	// 4. Find the Target (Semantic Search)
	fmt.Println("🔎 Scanning genome for the primary target...")
	vec, err := index.FetchEmbedding(cmd.Context(), nil, cfg.APIKey, cfg.EmbedModel, cfg.BaseURL, intent)
	if err != nil {
		return fmt.Errorf("failed to embed intent: %w", err)
	}

	queryVec := index.Normalize(vec)
	matches := store.SearchIntent(queryVec, 1)
	if len(matches) == 0 {
		return fmt.Errorf("no relevant code found in the genome for this intent")
	}

	topMatch := matches[0]
	fmt.Printf("🎯 Primary target identified: %s (Confidence: %.2f)\n", topMatch.PublicID, topMatch.Score)

	// 5. Retrieve Full Node Metadata
	// Note: Adjust .Nodes or .Registry.Nodes depending on your exact struct definition in genome.Load
	node, ok := registryManager.Registry.Nodes[topMatch.UUID]
	if !ok {
		return fmt.Errorf("target node %s found in index but missing from genome.json", topMatch.UUID)
	}

	// 6. Construct the Schema Anchor
	anchor := surgery.TargetAnchor{
		UUID:      node.UUID,
		PublicID:  node.PublicID,
		FilePath:  node.FilePath,
		NodeType:  node.NodeType,
		LogicHash: node.LogicHash,
	}

	// 7. Hydrate the Source Code
	fmt.Println("💧 Hydrating source code context...")
	contextNode, err := surgery.HydrateNode(cfg.Path, anchor, surgery.ReasonPrimaryTarget)
	if err != nil {
		return fmt.Errorf("failed to hydrate source code: %w", err)
	}

	// 8. Assemble the Blueprint
	plan := surgery.SurgeryPlan{
		Version:      surgery.SurgeryPlanVersion,
		Intent:       intent,
		PlannerModel: "gemini-2.5-pro", // Temporary hardcode for V3
		Target:       anchor,
		Context:      []surgery.ContextNode{contextNode},
	}

	// 9. Write to Disk
	yamlBytes, err := yaml.Marshal(&plan)
	if err != nil {
		return fmt.Errorf("failed to marshal surgery plan to YAML: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(draftOutputFile), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(draftOutputFile, yamlBytes, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", draftOutputFile, err)
	}

	fmt.Printf("✅ Draft successfully staged at: %s\n", draftOutputFile)
	return nil
}

// --- Configuration Loader ---

type DraftConfig struct {
	Path       string
	APIKey     string
	EmbedModel string
	BaseURL    string
}

func loadDraftConfig(cmd *cobra.Command) (DraftConfig, error) {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return DraftConfig{}, fmt.Errorf("failed to get path flag: %w", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return DraftConfig{}, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	embedModel := os.Getenv("SAAYN_EMBED_MODEL")
	if embedModel == "" {
		embedModel = "gemini-embedding-001" // Matches your V2 index
	}

	// ---------------------------------------------------------
	// THE FIX: Looking for the exact environment variable you use
	// and appending the required /v1beta route as the fallback.
	// ---------------------------------------------------------
	baseURL := os.Getenv("SAAYN_FAST_BASE_URL")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	return DraftConfig{
		Path:       path,
		APIKey:     apiKey,
		EmbedModel: embedModel,
		BaseURL:    baseURL,
	}, nil
}
