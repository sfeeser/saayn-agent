package saayn

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/saayn-agent/internal/genome/index"
	"github.com/spf13/cobra"
)

type SearchConfig struct {
	Path       string
	APIKey     string
	Model      string
	EmbedModel string
	BaseURL    string
	TopK       int
}

var searchIntentCmd = &cobra.Command{
	Use:   "search-intent <intent>",
	Short: "Search the semantic genome index using a natural language intent",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchIntent,
}

func init() {
	rootCmd.AddCommand(searchIntentCmd)
	searchIntentCmd.Flags().StringP("path", "p", ".", "Path to the project root")
	searchIntentCmd.Flags().IntP("top", "k", 5, "Number of top matches to return")
}

// runSearchIntent is the thin Cobra wrapper orchestrating the flow.
func runSearchIntent(cmd *cobra.Command, args []string) error {
	intent := strings.TrimSpace(args[0])
	if intent == "" {
		return fmt.Errorf("search intent cannot be empty")
	}

	cfg, err := loadSearchConfig(cmd)
	if err != nil {
		return err
	}

	matches, err := executeSearchIntent(cmd.Context(), intent, cfg)
	if err != nil {
		return err
	}

	printSearchResults(intent, matches)
	return nil
}

// loadSearchConfig extracts environment variables and CLI flags.
func loadSearchConfig(cmd *cobra.Command) (SearchConfig, error) {
	_ = godotenv.Load()

	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return SearchConfig{}, fmt.Errorf("failed to parse --path flag: %w", err)
	}

	topK, err := cmd.Flags().GetInt("top")
	if err != nil {
		return SearchConfig{}, fmt.Errorf("failed to parse --top flag: %w", err)
	}
	if topK <= 0 {
		return SearchConfig{}, fmt.Errorf("--top must be greater than 0")
	}

	apiKey := strings.TrimSpace(os.Getenv("SAAYN_FAST_API_KEY"))
	modelName := strings.TrimSpace(os.Getenv("SAAYN_FAST_MODEL"))
	baseURL := strings.TrimSpace(os.Getenv("SAAYN_FAST_BASE_URL"))
	embedModel := strings.TrimSpace(os.Getenv("SAAYN_EMBED_MODEL"))
	if embedModel == "" {
		embedModel = "gemini-embedding-001"
	}

	if apiKey == "" || modelName == "" || baseURL == "" {
		return SearchConfig{}, fmt.Errorf("FAST Cognition variables missing from environment. Run 'verify-llm-targets' to diagnose")
	}

	return SearchConfig{
		Path:       path,
		APIKey:     apiKey,
		EmbedModel: embedModel,
		BaseURL:    baseURL,
		TopK:       topK,
	}, nil
}

// executeSearchIntent performs semantic retrieval against the local index.
func executeSearchIntent(ctx context.Context, intent string, cfg SearchConfig) ([]index.Match, error) {
	indexPath := filepath.Join(cfg.Path, "genome.index.json")

	store, err := index.LoadIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("semantic index not found at %s; run enrich first: %w", indexPath, err)
	}
	if len(store.Records) == 0 {
		return nil, fmt.Errorf("semantic index is empty; run enrich first")
	}

	if store.Metadata.EmbeddingModel == "" {
		return nil, fmt.Errorf("semantic index metadata is missing embedding model")
	}
	if store.Metadata.EmbeddingModel != cfg.EmbedModel {
		return nil, fmt.Errorf(
			"index model mismatch: store uses %s, but config specifies %s",
			store.Metadata.EmbeddingModel,
			cfg.EmbedModel,
		)
	}

	vec, err := index.FetchEmbedding(ctx, nil, cfg.APIKey, cfg.EmbedModel, cfg.BaseURL, intent)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch query embedding: %w", err)
	}
	if len(vec) == 0 {
		return nil, fmt.Errorf("received empty vector for query")
	}

	if store.Metadata.VectorDimension != 0 && len(vec) != store.Metadata.VectorDimension {
		return nil, fmt.Errorf(
			"query vector dimension mismatch: expected %d, got %d",
			store.Metadata.VectorDimension,
			len(vec),
		)
	}

	queryVec := index.Normalize(vec)
	matches := store.SearchIntent(queryVec, cfg.TopK)

	return matches, nil
}

// printSearchResults renders the terminal UI for the matches.
func printSearchResults(intent string, matches []index.Match) {
	fmt.Printf("\n🔎 Search Intent: %s\n\n", intent)

	if len(matches) == 0 {
		fmt.Println("No matches found.")
		return
	}

	fmt.Println("Top Matches")
	for i, match := range matches {
		displayID := match.PublicID
		if len(displayID) > 50 {
			displayID = displayID[:47] + "..."
		}

		fmt.Printf("%d. %-50s %-20s %-10s %.2f\n",
			i+1,
			displayID,
			match.FilePath,
			match.NodeType,
			match.Score,
		)
	}
	fmt.Println()
}
