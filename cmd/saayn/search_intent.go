package saayn

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sfeeser/saayn-agent/internal/genome/index"
	"github.com/spf13/cobra"
)

type SearchConfig struct {
	Path       string
	APIKey     string
	EmbedModel string
	BaseURL    string
	TopK       int
	Verbose    bool
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
	searchIntentCmd.Flags().BoolP("verbose", "v", false, "Show the business purpose of the retrieved nodes")
}

func runSearchIntent(cmd *cobra.Command, args []string) error {
	intent := strings.TrimSpace(args[0])
	if intent == "" {
		return fmt.Errorf("search intent cannot be empty")
	}

	cfg, err := loadSearchConfig(cmd)
	if err != nil {
		return err
	}

	matches, store, err := executeSearchIntent(cmd.Context(), intent, cfg)
	if err != nil {
		return err
	}

	printSearchResults(intent, matches, store, cfg.Verbose)
	return nil
}

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

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return SearchConfig{}, fmt.Errorf("failed to parse --verbose flag: %w", err)
	}

	apiKey := strings.TrimSpace(os.Getenv("SAAYN_FAST_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("SAAYN_FAST_BASE_URL"))

	embedModel := strings.TrimSpace(os.Getenv("SAAYN_EMBED_MODEL"))
	if embedModel == "" {
		embedModel = "gemini-embedding-001"
	}

	if apiKey == "" || baseURL == "" {
		return SearchConfig{}, fmt.Errorf("FAST Cognition variables missing from environment. Run 'verify-llm-targets' to diagnose")
	}

	return SearchConfig{
		Path:       path,
		APIKey:     apiKey,
		EmbedModel: embedModel,
		BaseURL:    baseURL,
		TopK:       topK,
		Verbose:    verbose,
	}, nil
}

func executeSearchIntent(ctx context.Context, intent string, cfg SearchConfig) ([]index.Match, *index.IndexStore, error) {
	indexPath := filepath.Join(cfg.Path, "genome.index.json")

	store, err := index.LoadIndex(indexPath)
	if err != nil {
		return nil, nil, fmt.Errorf("semantic index not found at %s; run enrich first: %w", indexPath, err)
	}
	if len(store.Records) == 0 {
		return nil, nil, fmt.Errorf("semantic index is empty; run enrich first")
	}

	if store.Metadata.EmbeddingModel == "" {
		return nil, nil, fmt.Errorf("semantic index metadata is missing embedding model")
	}
	if store.Metadata.EmbeddingModel != cfg.EmbedModel {
		return nil, nil, fmt.Errorf(
			"index model mismatch: store uses %s, but config specifies %s",
			store.Metadata.EmbeddingModel,
			cfg.EmbedModel,
		)
	}

	// FetchEmbedding signature strictly matches the review (no nil argument)
	vec, err := index.FetchEmbedding(ctx, nil, cfg.APIKey, cfg.EmbedModel, cfg.BaseURL, intent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch query embedding: %w", err)
	}
	if len(vec) == 0 {
		return nil, nil, fmt.Errorf("received empty vector for query")
	}

	if store.Metadata.VectorDimension != 0 && len(vec) != store.Metadata.VectorDimension {
		return nil, nil, fmt.Errorf(
			"query vector dimension mismatch: expected %d, got %d",
			store.Metadata.VectorDimension,
			len(vec),
		)
	}

	queryVec := index.Normalize(vec)
	matches := store.SearchIntent(queryVec, cfg.TopK)

	return matches, store, nil
}

func printSearchResults(intent string, matches []index.Match, store *index.IndexStore, verbose bool) {
	fmt.Printf("\n🔎 Search Intent: %s\n\n", intent)

	if len(matches) == 0 {
		fmt.Println("No matches found.")
		return
	}

	fmt.Println("Top Matches")
	for i, match := range matches {
		record, ok := store.Get(match.UUID)

		filePath := "-"
		purposeText := ""

		if ok {
			filePath = effectiveFilePath(record)
			purposeText = extractFieldFromRetrievalDoc(record.RetrievalDocument, "BusinessPurpose:")
		}

		displayID := displaySymbol(match.PublicID, filePath)

		// Truncate overly long symbols safely
		if len(displayID) > 40 {
			displayID = displayID[:37] + "..."
		}

		// Truncate overly long paths from the left, keeping the filename visible
		if len(filePath) > 30 && filePath != "-" {
			filePath = "..." + filePath[len(filePath)-27:]
		}

		fmt.Printf("%d. %-40s %-30s %-10s %.2f\n",
			i+1,
			displayID,
			filePath,
			match.NodeType,
			match.Score,
		)

		if verbose && purposeText != "" {
			printWrappedIndented(purposeText, 80, "   ")
			fmt.Println()
		}
	}

	fmt.Println()
}

// --- Helper Functions ---

func effectiveFilePath(record index.EmbeddingRecord) string {
	// Primary Source of Truth: The structured data
	if strings.TrimSpace(record.FilePath) != "" {
		return record.FilePath
	}

	// Fallback only if FilePath was not persisted correctly upstream
	if fp := extractFieldFromRetrievalDoc(record.RetrievalDocument, "FilePath:"); fp != "" {
		return fp
	}

	return "-"
}

func displaySymbol(publicID, filePath string) string {
	// Only strip [file.go] when we have a real file path column to show it in.
	if filePath == "" || filePath == "-" {
		return publicID
	}

	if idx := strings.Index(publicID, "["); idx != -1 {
		return publicID[:idx]
	}

	return publicID
}

func extractFieldFromRetrievalDoc(doc, prefix string) string {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

func printWrappedIndented(text string, width int, indent string) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}

	fmt.Print(indent)
	lineLen := len(indent)

	for _, word := range words {
		if lineLen+len(word)+1 > width {
			fmt.Print("\n" + indent)
			lineLen = len(indent)
		}
		fmt.Print(word + " ")
		lineLen += len(word) + 1
	}
	fmt.Print("\n")
}
