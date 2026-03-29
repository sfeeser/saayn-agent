package saayn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/printer"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/saayn-agent/internal/genome"
	"github.com/saayn-agent/internal/scanner"
	"github.com/spf13/cobra"
)

var enrichCmd = &cobra.Command{
	Use:   "enrich",
	Short: "Uses FAST Cognition to automatically document the business purpose of code",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := godotenv.Load(); err != nil {
			// Silently fallback to process environment vars, which is standard for CI/CD
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("failed to parse --path flag: %w", err)
		}

		apiKey := strings.TrimSpace(os.Getenv("SAAYN_FAST_API_KEY"))
		model := strings.TrimSpace(os.Getenv("SAAYN_FAST_MODEL"))
		baseURL := strings.TrimSpace(os.Getenv("SAAYN_FAST_BASE_URL"))

		if apiKey == "" || model == "" || baseURL == "" {
			return fmt.Errorf("FAST Cognition variables missing from environment. Run 'verify-llm-targets' to diagnose")
		}

		fmt.Println("🧠 Starting Semantic Enrichment Process...")

		// 1. Load the Memory
		regManager, err := genome.Load(path + "/genome.json")
		if err != nil {
			return fmt.Errorf("error loading genome.json. Did you run 'init'?: %w", err)
		}

		// 2. Scan live code so we have the raw ASTs available
		liveNodes, err := scanner.FullScan(path)
		if err != nil {
			return fmt.Errorf("error scanning directory: %w", err)
		}

		// Map live ASTs by PublicID. Using standard formatting instead of RawFormat for better LLM readability.
		liveASTMap := make(map[string]string)
		conf := &printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
		for _, ln := range liveNodes {
			var buf strings.Builder
			_ = conf.Fprint(&buf, ln.Fset, ln.AST)
			liveASTMap[ln.PublicID] = buf.String()
		}

		// --- NEW SYNC LOGIC: Safely Adopt & Mint Undiscovered Functions ---
		// First, map what the brain already knows by PublicID
		knownFunctions := make(map[string]bool)
		for _, node := range regManager.Registry.Nodes {
			knownFunctions[node.PublicID] = true
		}

		// Now, compare live code against known functions
		syncCount := 0
		for _, liveNode := range liveNodes {
			if !knownFunctions[liveNode.PublicID] {
				// It's brand new! Mint a proper UUID for it.
				newUUID := uuid.New().String()
				liveNode.UUID = newUUID

				// Save it to the brain under the new UUID key
				regManager.Registry.Nodes[newUUID] = *liveNode
				syncCount++
			}
		}
		if syncCount > 0 {
			fmt.Printf("🌱 Minted UUIDs and adopted %d new functions into the genome...\n", syncCount)
		}
		// ------------------------------------------------------------------

		// 3. Telemetry Counters
		updatedCount := 0
		skippedAlready := 0
		skippedMissing := 0
		failedCount := 0

		// 4. The Enrichment Loop
		for id, node := range regManager.Registry.Nodes {
			if node.BusinessPurpose != "" {
				skippedAlready++
				continue
			}

			rawCode, exists := liveASTMap[node.PublicID]
			if !exists {
				fmt.Printf("  ⚠️  Skipping %s: no live code match found (node may have drifted)\n", node.PublicID)
				skippedMissing++
				continue
			}

			fmt.Printf("  🔍 Analyzing %s...\n", node.PublicID)

			purpose, err := generateBusinessPurpose(cmd.Context(), apiKey, model, baseURL, node.PublicID, rawCode)
			if err != nil {
				fmt.Printf("  ❌ Failed to analyze %s: %v\n", node.PublicID, err)
				failedCount++
				continue
			}

			// Save the insight back to the registry memory
			node.BusinessPurpose = purpose
			regManager.Registry.Nodes[id] = node
			updatedCount++

			fmt.Printf("  ✅ Insight: %s\n", purpose)
			time.Sleep(4 * time.Second) // Respect rate limits
		}

		// 5. Save and Summarize
		fmt.Println("\n📊 Enrichment Summary:")
		fmt.Printf("  - Updated: %d\n", updatedCount)
		fmt.Printf("  - Skipped (Already Enriched): %d\n", skippedAlready)
		fmt.Printf("  - Skipped (Missing/Drifted):  %d\n", skippedMissing)
		fmt.Printf("  - Failed:  %d\n", failedCount)

		if updatedCount > 0 {
			if err := regManager.Save(); err != nil {
				return fmt.Errorf("failed to save genome.json: %w", err)
			}
			fmt.Println("\n💾 Genome memory successfully updated.")
		} else if failedCount == 0 {
			fmt.Println("\n✅ Genome is fully enriched. No updates required.")
		}

		return nil
	},
}

func generateBusinessPurpose(ctx context.Context, apiKey, model, baseURL, funcName, rawCode string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/models/" + url.PathEscape(model) + ":generateContent"
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
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

	// 6. Strict Output Sanitization
	purpose := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
	purpose = strings.Trim(purpose, `"`) // Strip surrounding quotes if the AI adds them

	if purpose == "" {
		return "", fmt.Errorf("model returned empty string after normalization")
	}
	if strings.Contains(purpose, "```") {
		return "", fmt.Errorf("model returned markdown blocks despite instructions")
	}

	return purpose, nil
}

func init() {
	rootCmd.AddCommand(enrichCmd)
	enrichCmd.Flags().StringP("path", "p", ".", "Path to the project root")
}
