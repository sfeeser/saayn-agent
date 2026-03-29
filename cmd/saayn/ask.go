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
	"github.com/saayn-agent/internal/genome" // Add this line!
	"github.com/saayn-agent/internal/surgeon"
	"github.com/saayn-agent/pkg/model"
	"github.com/spf13/cobra"
)

// AgentResponse is the strict JSON format the LLM must return
type AgentResponse struct {
	Rationale  string `json:"rationale"`
	TargetUUID string `json:"target_uuid"`
	NewCode    string `json:"new_code"`
}

var askCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Ask the SAAYN AI to modify exactly one existing function",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := godotenv.Load(); err != nil {
			// Silently fallback to process env, standard for CI
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("❌ Failed to parse --path flag: %w", err)
		}

		userPrompt := strings.Join(args, " ")
		apiKey := strings.TrimSpace(os.Getenv("SAAYN_DEEP_API_KEY"))
		aiModel := strings.TrimSpace(os.Getenv("SAAYN_DEEP_MODEL")) // <-- Change to aiModel
		baseURL := strings.TrimSpace(os.Getenv("SAAYN_DEEP_BASE_URL"))

		if apiKey == "" || aiModel == "" || baseURL == "" { // <-- Update here
			return fmt.Errorf("❌ DEEP Cognition variables missing from .env. Run 'verify-llm-targets'")
		}

		fmt.Println("🧠 Waking up SAAYN Agent (DEEP Tier)...")

		// 1. Load the Semantic Memory
		regManager, err := genome.Load(path + "/genome.json")
		if err != nil {
			return fmt.Errorf("❌ Error loading genome.json. Cannot think without memory: %w", err)
		}

		// 2. Prepare Context & Build UUID Index for local validation
		var contextBuilder strings.Builder
		contextBuilder.WriteString("CODEBASE GENOME (Available Functions):\n")
		contextBuilder.WriteString("=========================================\n")

		// Create an index to quickly verify the AI's choice later
		uuidIndex := make(map[string]string) // UUID -> PublicID mapping

		for _, node := range regManager.Registry.Nodes {
			purpose := node.BusinessPurpose
			if strings.TrimSpace(purpose) == "" {
				purpose = "(Not yet enriched)"
			}
			uuidIndex[node.UUID] = node.PublicID

			contextBuilder.WriteString(fmt.Sprintf("UUID: %s\nFunction: %s\nPurpose: %s\n---\n",
				node.UUID, node.PublicID, purpose))
		}

		fmt.Printf("📦 Loaded %d functions into AI context.\n", len(regManager.Registry.Nodes))
		fmt.Printf("🗣️  Prompt: \"%s\"\n", userPrompt)
		fmt.Println("⏳ Thinking... (This may take a moment)")

		// 3. Request Mutation Decision
		aiResponseJSON, err := requestMutationDecision(cmd.Context(), apiKey, aiModel, baseURL, contextBuilder.String(), userPrompt)
		if err != nil {
			return fmt.Errorf("❌ DEEP Cognition Failure: %w", err)
		}

		// 4. Parse the Decision
		var decision AgentResponse
		if err := json.Unmarshal([]byte(aiResponseJSON), &decision); err != nil {
			raw := strings.TrimSpace(aiResponseJSON)
			if len(raw) > 500 {
				raw = raw[:500] + "... [truncated]"
			}
			return fmt.Errorf("❌ Failed to parse AI decision. Raw output: %s\nError: %w", raw, err)
		}

		fmt.Println("\n💡 AI Rationale:")
		fmt.Printf("   %s\n\n", decision.Rationale)

		// 5. LOCAL CGS VALIDATION (The Trust Boundary)
		decision.TargetUUID = strings.TrimSpace(decision.TargetUUID)
		decision.NewCode = strings.TrimSpace(decision.NewCode)

		// Handle the AI's Refusal Path
		if decision.TargetUUID == "" {
			fmt.Println("⚠️  AI determined it cannot fulfill this request by modifying exactly one function.")
			return nil
		}

		// Verify the UUID actually exists in our Genome
		targetPublicID, exists := uuidIndex[decision.TargetUUID]
		if !exists {
			return fmt.Errorf("❌ Security Exception: AI attempted to mutate unknown UUID: %q", decision.TargetUUID)
		}

		// Verify the code payload isn't completely broken
		if decision.NewCode == "" {
			return fmt.Errorf("❌ Validation Error: AI returned empty replacement code")
		}
		if strings.HasPrefix(decision.NewCode, "func ") {
			return fmt.Errorf("❌ Validation Error: AI returned a full 'func' declaration instead of just the inner body")
		}
		if strings.Contains(decision.NewCode, "```") {
			return fmt.Errorf("❌ Validation Error: AI included markdown blocks in the raw code payload")
		}

		// 6. Pre-Surgery Confirmation
		fmt.Printf("🎯 Target UUID: %s\n", decision.TargetUUID)
		fmt.Printf("🔧 Target Function: %s\n", targetPublicID)
		fmt.Println("🏥 Handing code to the Surgeon for operation...")

		// 7. Execute Surgery
		surgeryReq := model.SurgeryRequest{
			TargetUUID:  decision.TargetUUID,
			NewLogic:    decision.NewCode,
			Registry:    regManager.Registry,
			ProjectRoot: path,
		}

		err = surgeon.PerformSurgery(surgeryReq)
		if err != nil {
			return fmt.Errorf("❌ Surgery Failed: %w", err)
		}

		// --- NEW MEDICAL HISTORY LOGIC ---
		regManager.Registry.LastMutatedUUID = decision.TargetUUID
		if err := regManager.Save(); err != nil {
			return fmt.Errorf("⚠️ Surgery succeeded, but failed to save Medical History: %w", err)
		}
		// ---------------------------------

		fmt.Println("✅ Operation Complete. Medical History updated. Run './saayn verify' to check for drift!")
		return nil
	},
}

func requestMutationDecision(ctx context.Context, apiKey, model, baseURL, codeContext, userPrompt string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/models/" + url.PathEscape(model) + ":generateContent"
	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	systemPrompt := `You are SAAYN, an autonomous Code Genome Agent.
Your job is to read the user's request, look at the CODEBASE GENOME, and write the exact Go code needed to fulfill the request.

RULES:
1. You must select EXACTLY ONE Target UUID from the genome to modify. Do not invent or alter UUIDs.
2. You must write the complete, raw inner body of the function (no 'func' declaration, just the inner code).
3. Do not change the function signature, name, parameters, receiver, or return types.
4. If the request CANNOT be satisfied by modifying exactly one existing function, return an empty string for target_uuid and new_code, and explain why in the rationale.
5. Return your answer as a STRICT JSON object matching this schema:
{
  "rationale": "Explain your technical decision-making process",
  "target_uuid": "The exact UUID string (or empty if refusal)",
  "new_code": "The raw Go code string (or empty if refusal)"
}
6. Do not return placeholder code, pseudocode, or markdown fences. Return ONLY the raw JSON object.`

	fullPrompt := fmt.Sprintf("%s\n\n%s\n\nUSER REQUEST: %s", systemPrompt, codeContext, userPrompt)

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": fullPrompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"responseMimeType": "application/json",
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

	client := &http.Client{Timeout: 60 * time.Second}
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

	rawText := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")

	return strings.TrimSpace(rawText), nil
}

func init() {
	rootCmd.AddCommand(askCmd)
	askCmd.Flags().StringP("path", "p", ".", "Path to the project root")
}
