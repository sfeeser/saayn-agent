package saayn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/printer"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/saayn-agent/internal/genome"
	"github.com/saayn-agent/internal/genome/index"
	"github.com/saayn-agent/internal/scanner"
	"github.com/saayn-agent/pkg/model"
	"github.com/spf13/cobra"
)

type EnrichConfig struct {
	Path          string
	APIKey        string
	Model         string
	EmbedModel    string
	BaseURL       string
	DelayDuration time.Duration
}

type EnrichmentStats struct {
	Updated        int
	SkippedAlready int
	SkippedMissing int
	Failed         int
}

var enrichCmd = &cobra.Command{
	Use:   "enrich",
	Short: "Uses FAST Cognition to automatically document the business purpose of code",
	RunE:  runEnrich,
}

func runEnrich(cmd *cobra.Command, args []string) error {
	cfg, err := loadEnrichConfig(cmd)
	if err != nil {
		return err
	}

	fmt.Println("🧠 Starting Semantic Enrichment Process...")

	// 1. Prepare Inputs
	regManager, liveNodes, liveASTMap, err := prepareEnrichmentInputs(cfg.Path)
	if err != nil {
		return err
	}

	// 2. Prepare Index (V2 Change)
	indexPath := cfg.Path + "/genome.index.json"
	idxStore, err := index.LoadIndex(indexPath)
	if err != nil {
		// If index is missing or corrupt, start a fresh one for this model
		fmt.Println("📭 No existing semantic index found. Creating new index...")
		idxStore = index.NewIndexStore(cfg.EmbedModel, 0)
	}

	// 3. Adopt New Nodes
	syncCount := adoptNewNodes(regManager, liveNodes)
	if syncCount > 0 {
		fmt.Printf("🌱 Adopted %d new functions into the genome...\n", syncCount)
	}

	// 4. Run Semantic Enrichment (AI Purpose Generation)
	stats := enrichRegistry(cmd.Context(), regManager, liveNodes, liveASTMap, cfg)

	// 5. Run Semantic Indexing (V2 SyncIndex Orchestration)
	fmt.Println("\n📡 Synchronizing Semantic Index...")
	syncStats, err := index.SyncIndex(
		cmd.Context(),
		idxStore,
		regManager.Registry.Nodes,
		cfg.APIKey,
		cfg.EmbedModel,
		cfg.BaseURL,
	)
	if err != nil {
		// We don't want to fail the whole command if indexing fails,
		// but we should warn the user.
		fmt.Printf("⚠️  Index Sync Warning: %v\n", err)
	}

	// 6. Report Summaries
	printEnrichmentSummary(stats)
	printIndexSummary(syncStats)

	// 7. Save Results
	if err := saveEnrichmentResults(regManager, stats, syncCount); err != nil {
		return err
	}

	// Save Index if changes occurred
	if syncStats.Created > 0 || syncStats.Updated > 0 || syncStats.Deleted > 0 {
		if err := idxStore.Save(indexPath); err != nil {
			return fmt.Errorf("failed to save semantic index: %w", err)
		}
		fmt.Println("💾 Semantic index successfully updated.")
	}

	return nil
}

func loadEnrichConfig(cmd *cobra.Command) (EnrichConfig, error) {
	if err := godotenv.Load(); err != nil {
		// Silently fallback to process environment vars, which is standard for CI/CD
	}

	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return EnrichConfig{}, fmt.Errorf("failed to parse --path flag: %w", err)
	}

	apiKey := strings.TrimSpace(os.Getenv("SAAYN_FAST_API_KEY"))
	modelName := strings.TrimSpace(os.Getenv("SAAYN_FAST_MODEL"))
	baseURL := strings.TrimSpace(os.Getenv("SAAYN_FAST_BASE_URL"))
	embedModel := strings.TrimSpace(os.Getenv("SAAYN_EMBED_MODEL"))
	if embedModel == "" {
		embedModel = "text-embedding-004" // Safe default for Google
	}

	delaySec := 4
	if val, err := strconv.Atoi(os.Getenv("SAAYN_API_DELAY_SECONDS")); err == nil && val > 0 {
		delaySec = val
	}

	if apiKey == "" || modelName == "" || baseURL == "" {
		return EnrichConfig{}, fmt.Errorf("FAST Cognition variables missing from environment. Run 'verify-llm-targets' to diagnose")
	}

	return EnrichConfig{
		Path:          path,
		APIKey:        apiKey,
		Model:         modelName,
		EmbedModel:    embedModel,
		BaseURL:       baseURL,
		DelayDuration: time.Duration(delaySec) * time.Second,
	}, nil
}

func prepareEnrichmentInputs(path string) (*genome.RegistryManager, []*model.Node, map[string]string, error) {
	regManager, err := genome.Load(path + "/genome.json")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error loading genome.json. Did you run 'init'?: %w", err)
	}

	liveNodes, err := scanner.FullScan(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error scanning directory: %w", err)
	}

	liveASTMap := buildLiveASTMap(liveNodes)

	return regManager, liveNodes, liveASTMap, nil
}

func buildLiveASTMap(liveNodes []*model.Node) map[string]string {
	liveASTMap := make(map[string]string)
	conf := &printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 4,
	}

	for _, ln := range liveNodes {
		var buf strings.Builder
		_ = conf.Fprint(&buf, ln.Fset, ln.AST)
		liveASTMap[ln.PublicID] = buf.String()
	}

	return liveASTMap
}

func adoptNewNodes(regManager *genome.RegistryManager, liveNodes []*model.Node) int {
	knownFunctions := make(map[string]bool)
	for _, node := range regManager.Registry.Nodes {
		knownFunctions[node.PublicID] = true
	}

	syncCount := 0
	for _, liveNode := range liveNodes {
		if knownFunctions[liveNode.PublicID] {
			continue
		}

		newUUID := uuid.New().String()
		liveNode.UUID = newUUID
		regManager.Registry.Nodes[newUUID] = *liveNode
		syncCount++
	}

	return syncCount
}

func enrichRegistry(
	ctx context.Context,
	regManager *genome.RegistryManager,
	liveNodes []*model.Node,
	liveASTMap map[string]string,
	cfg EnrichConfig,
) EnrichmentStats {
	stats := EnrichmentStats{}

	for id := range regManager.Registry.Nodes {
		processRegistryNode(ctx, regManager, id, liveNodes, liveASTMap, cfg, &stats)
	}

	return stats
}

func processRegistryNode(
	ctx context.Context,
	regManager *genome.RegistryManager,
	id string,
	liveNodes []*model.Node,
	liveASTMap map[string]string,
	cfg EnrichConfig,
	stats *EnrichmentStats,
) {
	node := regManager.Registry.Nodes[id]

	rawCode, exists := liveASTMap[node.PublicID]
	if !exists {
		fmt.Printf("   🗑️  %-45s [DELETED FROM DISK]\n", node.PublicID)
		delete(regManager.Registry.Nodes, id)
		stats.Updated++
		stats.SkippedMissing++
		return
	}

	logicChanged := syncNodeState(regManager, &node, liveNodes)

	if !logicChanged && node.BusinessPurpose != "" {
		stats.SkippedAlready++
		regManager.Registry.Nodes[id] = node
		return
	}

	fmt.Print("\r\033[2K")
	printActiveNodeHeader(node)

	if logicChanged {
		fmt.Println("   ├─ 🧬 logic changed")
		fmt.Println("   ├─ 🔄 purpose reset")
	}

	fmt.Println("   ├─ 🔍 analyzing")

	purpose, err := generateBusinessPurpose(
		ctx,
		cfg.APIKey,
		cfg.Model,
		cfg.BaseURL,
		node.PublicID,
		rawCode,
	)
	if err != nil {
		fmt.Printf("   └─ ❌ analysis failed: %v\n", err)
		stats.Failed++
		return
	}

	node.BusinessPurpose = purpose
	regManager.Registry.Nodes[id] = node
	stats.Updated++

	fmt.Println("   └─ ✅ purpose updated")
	wrapped := wrapText(purpose, 80, "      ")
	fmt.Printf("%s\n\n", wrapped)

	time.Sleep(cfg.DelayDuration)
}

func syncNodeState(
	regManager *genome.RegistryManager,
	node *model.Node,
	liveNodes []*model.Node,
) bool {
	logicChanged := false

	for _, ln := range liveNodes {
		if ln.PublicID != node.PublicID {
			continue
		}

		if node.FilePath != ln.FilePath {
			node.FilePath = ln.FilePath
		}

		newHash := regManager.NormalizeAndHash(ln)
		if node.LogicHash != newHash {
			logicChanged = true
			node.LogicHash = newHash
			node.BusinessPurpose = ""
		}

		break
	}

	return logicChanged
}

func printActiveNodeHeader(node model.Node) {
	fmt.Println()

	icon := " ⚙️ "
	if node.NodeType == "struct" {
		icon = " 📦 "
	}

	displayHash := "n/a"
	if len(node.LogicHash) >= 8 {
		displayHash = node.LogicHash[:8]
	}

	fmt.Printf(" %s %-45s %-25s [%s]\n", icon, node.PublicID, node.FilePath, displayHash)
}

func printEnrichmentSummary(stats EnrichmentStats) {
	fmt.Println("\n📊 Enrichment Summary:")
	fmt.Printf("  - Updated: %d\n", stats.Updated)
	fmt.Printf("  - Skipped (Already Enriched): %d\n", stats.SkippedAlready)
	fmt.Printf("  - Skipped (Missing/Drifted):  %d\n", stats.SkippedMissing)
	fmt.Printf("  - Failed:  %d\n", stats.Failed)
}

func printIndexSummary(stats index.SyncStats) {
	if stats.Created == 0 && stats.Updated == 0 && stats.Deleted == 0 {
		fmt.Println("✅ Semantic index is already up to date.")
		return
	}

	fmt.Println("\n🧠 Semantic Index Summary:")
	fmt.Printf("  - Created: %d\n", stats.Created)
	fmt.Printf("  - Updated: %d\n", stats.Updated)
	fmt.Printf("  - Deleted: %d\n", stats.Deleted)
	fmt.Printf("  - Skipped: %d\n", stats.Skipped)
}

func saveEnrichmentResults(
	regManager *genome.RegistryManager,
	stats EnrichmentStats,
	syncCount int,
) error {
	if stats.Updated > 0 || syncCount > 0 {
		if err := regManager.Save(); err != nil {
			return fmt.Errorf("failed to save genome.json: %w", err)
		}
		fmt.Println("\n💾 Genome memory successfully updated.")
		return nil
	}

	if stats.Failed == 0 {
		fmt.Println("\n✅ Genome is fully enriched. No updates required.")
	}

	return nil
}

func generateBusinessPurpose(ctx context.Context, apiKey, modelName, baseURL, funcName, rawCode string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/models/" + url.PathEscape(modelName) + ":generateContent"
	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	prompt := fmt.Sprintf(`System Context: You are the analytical engine for SAAYN, an autonomous Code Genome Agent.

Task:
Read the provided Go function and explain its business purpose in no more than two sentences.

Rules:
- Explain why this function exists.
- Explain what value it provides to the system.
- Focus on business or system purpose, not line-by-line implementation.
- Do not output markdown fences (e.g. no triple backticks).
- Do not output prefixes, bullets, or pleasantries.
- Return only the raw summary text.

Function Identity: %s
Raw Code:
%s`, funcName, rawCode)

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	maxRetries := 3
	var lastErr error
	var explicitWait time.Duration

	baseDelaySec := 15
	if val, err := strconv.Atoi(os.Getenv("SAAYN_API_DELAY_SECONDS")); err == nil && val > 0 {
		baseDelaySec = val
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			var backoff time.Duration

			if explicitWait > 0 {
				backoff = explicitWait + time.Second
				fmt.Printf("      ⏳ API demanded hold. Sleeping %v...\n", backoff)
			} else {
				backoff = time.Duration(float64(baseDelaySec)*math.Pow(2, float64(attempt-1))) * time.Second
				fmt.Printf("      ⏳ Rate limited. Retrying in %v...\n", backoff)
			}

			time.Sleep(backoff)
			explicitWait = 0
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(jsonData))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)

		if err != nil {
			lastErr = err
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			bodyStr := string(bodyBytes)
			cleanMsg := "Rate Limit Exceeded"

			if match := regexp.MustCompile(`"retryDelay":\s*"([^"]+)"`).FindStringSubmatch(bodyStr); len(match) > 1 {
				cleanMsg = fmt.Sprintf("API requests a pause of %s", match[1])
				if parsed, err := time.ParseDuration(match[1]); err == nil {
					explicitWait = parsed
				}
			} else if match := regexp.MustCompile(`Please retry in ([a-zA-Z0-9.]+)`).FindStringSubmatch(bodyStr); len(match) > 1 {
				cleanMsg = fmt.Sprintf("API requests a pause of %s", match[1])
				if parsed, err := time.ParseDuration(match[1]); err == nil {
					explicitWait = parsed
				}
			}

			lastErr = fmt.Errorf("HTTP 429: %s", cleanMsg)
			fmt.Printf("      ⚠️  %s\n", lastErr.Error())
			continue
		}

		if resp.StatusCode != http.StatusOK {
			reason := strings.TrimSpace(string(bodyBytes))
			if len(reason) > 500 {
				reason = reason[:500] + "... [truncated]"
			}
			return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, reason)
		}

		var geminiResp struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}

		if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
			return "", err
		}

		if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
			return "", fmt.Errorf("empty response from AI")
		}

		purpose := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
		purpose = strings.Trim(purpose, `"`)

		if purpose == "" {
			return "", fmt.Errorf("model returned empty string after normalization")
		}
		if strings.Contains(purpose, "```") {
			return "", fmt.Errorf("model returned markdown blocks despite instructions")
		}

		return purpose, nil
	}

	return "", fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}

func init() {
	rootCmd.AddCommand(enrichCmd)
	enrichCmd.Flags().StringP("path", "p", ".", "Path to the project root")
}

func wrapText(text string, width int, indent string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var wrapped strings.Builder
	line := indent

	for _, word := range words {
		if len(line)+len(word) > width {
			wrapped.WriteString(line + "\n")
			line = indent + word
		} else {
			if line != indent {
				line += " "
			}
			line += word
		}
	}

	wrapped.WriteString(line)
	return wrapped.String()
}
