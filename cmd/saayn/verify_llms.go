package saayn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var verifyLlmCmd = &cobra.Command{
	Use:   "verify-llm-targets",
	Short: "Verifies end-to-end LLM connectivity for both FAST and DEEP cognitive tiers",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(); err != nil {
			fmt.Println("⚠️  No .env file found. Falling back to system environment variables.")
		}

		fmt.Println("📡 Verifying SAAYN Cognitive Engine...")
		fmt.Println("==================================================")

		// --- Test FAST Tier ---
		fmt.Println("⚡ TIER 1: FAST Cognition (Enrichment/Reading)")
		errFast := verifyTarget(
			cmd.Context(),
			"SAAYN_FAST",
			strings.TrimSpace(os.Getenv("SAAYN_FAST_API_KEY")),
			strings.TrimSpace(os.Getenv("SAAYN_FAST_MODEL")),
			strings.TrimSpace(os.Getenv("SAAYN_FAST_BASE_URL")),
		)
		if errFast != nil {
			fmt.Printf("❌ FAST Tier failed: %v\n", errFast)
		}

		fmt.Println("\n--------------------------------------------------")

		// --- Test DEEP Tier ---
		fmt.Println("🧠 TIER 2: DEEP Cognition (Surgery/Reasoning)")
		errDeep := verifyTarget(
			cmd.Context(),
			"SAAYN_DEEP",
			strings.TrimSpace(os.Getenv("SAAYN_DEEP_API_KEY")),
			strings.TrimSpace(os.Getenv("SAAYN_DEEP_MODEL")),
			strings.TrimSpace(os.Getenv("SAAYN_DEEP_BASE_URL")),
		)
		if errDeep != nil {
			fmt.Printf("❌ DEEP Tier failed: %v\n", errDeep)
		}

		fmt.Println("==================================================")
		if errFast == nil && errDeep == nil {
			fmt.Println("🎉 BOTH hemispheres are verified and fully operational!")
		} else {
			fmt.Println("⚠️  Cognitive Engine degraded. Check errors above.")
			os.Exit(1)
		}
	},
}

func verifyTarget(ctx context.Context, tierName, apiKey, model, baseURL string) error {
	// 1. Strict Configuration Validation
	if apiKey == "" {
		return fmt.Errorf("Config Error: %s_API_KEY is missing from environment", tierName)
	}
	if model == "" {
		return fmt.Errorf("Config Error: %s_MODEL is missing from environment", tierName)
	}
	if baseURL == "" {
		return fmt.Errorf("Config Error: %s_BASE_URL is missing from environment", tierName)
	}

	// 2. Safe URL Construction
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("Config Error: invalid %s_BASE_URL format: %w", tierName, err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/models/" + url.PathEscape(model) + ":generateContent"

	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	// 3. Build the Payload
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": "System check. Reply with exactly one word: PONG"},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("Internal Error: failed to marshal JSON request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("Internal Error: failed to construct HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("  🔍 Sending micro-prompt to %s...\n", model)

	// 4. Execute the Request
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		printFailure(u, tierName, "Transport Layer", fmt.Sprintf("Network request failed: %v", err))
		return fmt.Errorf("transport verification failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Transport Error: failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		reason := strings.TrimSpace(string(bodyBytes))
		printFailure(u, tierName, fmt.Sprintf("HTTP %d", resp.StatusCode), reason)
		return fmt.Errorf("http verification failed: status=%d, reason=%s", resp.StatusCode, reason)
	}

	// 5. Parse the Response
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
		return fmt.Errorf("Parsing Error: failed to unmarshal Gemini JSON: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("Semantic Error: API returned a 200 OK, but the response body was empty")
	}

	// 6. Strict Semantic Assertion
	aiResponse := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
	if strings.ToUpper(aiResponse) != "PONG" {
		return fmt.Errorf("Semantic Error: expected 'PONG', but AI responded with: %q", aiResponse)
	}

	fmt.Printf("  ✅ Received assertion: %q\n", aiResponse)
	return nil
}

func printFailure(u *url.URL, tierName, stage, errMsg string) {
	fmt.Printf("  ❌ %s Verification Failed!\n", tierName)
	fmt.Printf("     Stage:  %s\n", stage)
	fmt.Printf("     Reason: %s\n\n", errMsg)
	fmt.Println("🛠️  DEBUGGING: Try running this raw curl command to isolate the issue:")
	fmt.Println("--------------------------------------------------------------------------------")

	safeURL := *u
	q := safeURL.Query()
	key := q.Get("key")
	if len(key) > 8 {
		q.Set("key", key[:4]+"..."+key[len(key)-4:])
	} else if key != "" {
		q.Set("key", "****")
	}
	safeURL.RawQuery = q.Encode()

	fmt.Printf(`curl -v -X POST "%s" \
-H 'Content-Type: application/json' \
-d '{"contents":[{"parts":[{"text":"System check. Reply with exactly one word: PONG"}]}]}'`+"\n", safeURL.String())
	fmt.Println("--------------------------------------------------------------------------------")
}
func init() {
	rootCmd.AddCommand(verifyLlmCmd)
}
