package saayn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/sfeeser/saayn-agent/internal/genome/surgery"
)

var (
	planInputFile  string
	planOutputFile string
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate and auto-refine a strict multi-file code patch set from a surgery plan",
	RunE:  runPlan,
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().StringVarP(&planInputFile, "file", "f", "surgery.yaml", "Path to the input surgery plan")
	planCmd.Flags().StringVarP(&planOutputFile, "output", "o", "patch.yaml", "Path for the output patch set file")
}

func runPlan(cmd *cobra.Command, args []string) error {
	fmt.Printf("🧠 Initializing V5 Multi-File Planner for %s...\n", planInputFile)

	cfg, err := loadPlanConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to load planner config: %w", err)
	}

	yamlBytes, err := os.ReadFile(planInputFile)
	if err != nil {
		return fmt.Errorf("failed to read plan: %w", err)
	}

	var plan surgery.SurgeryPlan
	if err := yaml.Unmarshal(yamlBytes, &plan); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	if plan.Version == "" || plan.Target.PublicID == "" {
		return fmt.Errorf("invalid surgery plan: missing version or target")
	}

	maxRetries := 3
	var review surgery.PatchSetReview
	review.Status = "failed_validation"

	var finalSet surgery.SurgeryPatchSet
	feedback := ""

	// --- THE AUTO-REFINE LOOP ---
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("🚀 Generating multi-file patch set (Attempt %d/%d)...\n", attempt, maxRetries)
		review.Attempts = attempt

		prompt, err := buildPrompt(plan, feedback)
		if err != nil {
			return fmt.Errorf("failed to build prompt: %w", err)
		}

		jsonResponse, err := callGeminiPlanner(cmd.Context(), cfg, prompt)
		if err != nil {
			return fmt.Errorf("LLM planning failed: %w", err)
		}

		fmt.Println("👀 Performing multi-file code review...")

		var patchSet surgery.SurgeryPatchSet
		if err := json.Unmarshal([]byte(jsonResponse), &patchSet); err != nil {
			fmt.Printf("   ❌ Rejected: AI returned malformed JSON. Bouncing back for repair...\n")
			recordAttempt(&review, attempt, "json_parse", err)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: json_parse\nError: %v\n\nPlease output strictly valid JSON matching the schema.", err)
			continue
		}
		patchSet.Version = surgery.SurgeryPatchVersion

		// THE IRON GATE: Strict Identity and Coverage Validation
		if err := validatePatchSet(plan, patchSet); err != nil {
			fmt.Printf("   ❌ Rejected: Identity/Coverage drift (%v). Bouncing back for repair...\n", err)
			recordAttempt(&review, attempt, "schema_validation", err)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: schema_validation\nError: %v\n\nEvery node in the plan must have exactly one record with matching identity fields.", err)
			continue
		}

		// --- V5 AST Splice Analysis (Multi-File) ---
		var validationFailures []string
		fset := token.NewFileSet()

		for _, rec := range patchSet.Records {
			if rec.Action == surgery.ActionNone {
				continue
			}

			// Read original source for JIT validation
			originalBytes, err := os.ReadFile(rec.TargetNode.FilePath)
			if err != nil {
				validationFailures = append(validationFailures, fmt.Sprintf("[%s]: Failed to read file: %v", rec.TargetNode.PublicID, err))
				continue
			}

			// Validate Splice
			splicedSrc, err := spliceAST(rec.TargetNode.FilePath, rec.TargetNode.PublicID, originalBytes, rec.NewCode)
			if err != nil {
				shortErr := strings.Split(err.Error(), "\n")[0]
				validationFailures = append(validationFailures, fmt.Sprintf("[%s]: Splice failed: %s", rec.TargetNode.PublicID, shortErr))
				continue
			}

			// Validate Syntax
			if _, err := parser.ParseFile(fset, rec.TargetNode.FilePath, splicedSrc, parser.AllErrors); err != nil {
				shortErr := strings.Split(err.Error(), "\n")[0]
				validationFailures = append(validationFailures, fmt.Sprintf("[%s]: Syntax error: %s", rec.TargetNode.PublicID, shortErr))
				continue
			}
			fmt.Printf("   ↳ Validated syntax for: %s\n", rec.TargetNode.PublicID)
		}

		if len(validationFailures) > 0 {
			fmt.Printf("   ❌ Rejected: %d syntax/splice errors detected. Bouncing back for repair...\n", len(validationFailures))
			errCombo := errors.New(strings.Join(validationFailures, "\n"))
			recordAttempt(&review, attempt, "splice_and_parse", errCombo)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: splice_and_parse\nThe following nodes failed to compile or splice cleanly:\n%s", errCombo.Error())
			continue
		}

		fmt.Println("✅ Validation passed! All files in the set are sound.")
		recordAttempt(&review, attempt, "validation_passed", nil)
		review.Status = surgery.StatusApproved
		finalSet = patchSet
		break
	}

	if review.Status != surgery.StatusApproved {
		return fmt.Errorf("🛑 Failed to generate a valid multi-file patch set after %d attempts", maxRetries)
	}

	finalSet.Review = &review

	patchYamlBytes, err := yaml.Marshal(&finalSet)
	if err != nil {
		return fmt.Errorf("failed to marshal patch set to YAML: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(planOutputFile), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(planOutputFile, patchYamlBytes, 0644); err != nil {
		return fmt.Errorf("failed to save patch set: %w", err)
	}

	fmt.Printf("💾 Approved V5 surgery patch set saved at: %s\n", planOutputFile)
	return nil
}

// --- Validation Helpers ---

func validatePatchSet(plan surgery.SurgeryPlan, set surgery.SurgeryPatchSet) error {
	if len(set.Records) == 0 {
		return fmt.Errorf("patch set contains no records")
	}

	expected := make(map[string]surgery.TargetAnchor)
	expected[plan.Target.PublicID] = plan.Target
	for _, ctx := range plan.Context {
		expected[ctx.PublicID] = ctx.TargetAnchor
	}

	seen := make(map[string]bool)
	hasPrimary := false

	for _, rec := range set.Records {
		pid := rec.TargetNode.PublicID
		if pid == "" {
			return fmt.Errorf("record found with empty public_id")
		}

		exp, ok := expected[pid]
		if !ok {
			return fmt.Errorf("hallucination detected: %s was not part of the plan", pid)
		}
		if seen[pid] {
			return fmt.Errorf("duplicate record detected for %s", pid)
		}
		seen[pid] = true

		if rec.TargetNode.UUID != exp.UUID {
			return fmt.Errorf("[%s]: uuid mismatch", pid)
		}
		if rec.TargetNode.FilePath != exp.FilePath {
			return fmt.Errorf("[%s]: file_path mismatch", pid)
		}
		if rec.TargetNode.NodeType != exp.NodeType {
			return fmt.Errorf("[%s]: node_type mismatch", pid)
		}

		switch rec.Role {
		case surgery.RolePrimaryTarget:
			if pid != plan.Target.PublicID {
				return fmt.Errorf("[%s]: only the actual surgery target can have the Primary role", pid)
			}
			hasPrimary = true
		case surgery.RoleBlastRadius:
			// okay
		default:
			return fmt.Errorf("[%s]: invalid role: %s", pid, rec.Role)
		}

		if err := validateActionConsistency(rec); err != nil {
			return fmt.Errorf("[%s]: %w", pid, err)
		}
	}

	if !hasPrimary {
		return fmt.Errorf("patch set is missing the primary_target record")
	}
	if len(seen) != len(expected) {
		return fmt.Errorf("coverage failure: patch contains %d nodes, but plan requires %d", len(seen), len(expected))
	}

	return nil
}

func validateActionConsistency(rec surgery.PatchRecord) error {
	if rec.Action == surgery.ActionNone {
		if strings.TrimSpace(rec.NewCode) != "" {
			return fmt.Errorf("action is 'none' but new_code is provided")
		}
		if rec.Status != surgery.StatusNoChange {
			return fmt.Errorf("action is 'none' but status is not 'no_change'")
		}
	} else if rec.Action == surgery.ActionReplaceNode {
		if strings.TrimSpace(rec.NewCode) == "" {
			return fmt.Errorf("action is 'replace_node' but new_code is empty")
		}
		if rec.Status != surgery.StatusApproved {
			return fmt.Errorf("action is 'replace_node' but status is not 'approved'")
		}
	} else {
		return fmt.Errorf("invalid action: %s", rec.Action)
	}
	return nil
}

func recordAttempt(review *surgery.PatchSetReview, attempt int, stage string, err error) {
	rec := surgery.ReviewAttempt{
		Attempt: attempt,
		Stage:   stage,
	}
	if err != nil {
		rec.Errors = []string{err.Error()}
	}
	review.History = append(review.History, rec)
}

// --- Prompt Engineering ---

func buildPrompt(plan surgery.SurgeryPlan, feedback string) (string, error) {
	planYaml, err := yaml.Marshal(plan)
	if err != nil {
		return "", fmt.Errorf("failed to marshal plan: %w", err)
	}

	basePrompt := fmt.Sprintf(`You are an expert Go developer and AST surgeon executing multi-file refactors.
Your task is to generate a strict, compiler-ready JSON patch set based on the provided surgery plan.

SURGERY PLAN:
%s

INSTRUCTIONS:
1. Analyze the "intent" and how it affects the "target" and its "context" callers.
2. For EVERY node listed in the plan (Target + all Context nodes), you must output exactly one record.
3. If a node requires modification, set action to "replace_node", status to "approved", and provide the "new_code".
4. If a node (like a caller) does not need changes, set action to "none" and status to "no_change".
5. Target nodes MUST exactly match the identities (uuid, public_id, file_path, node_type) provided in the plan. Do not invent identities.

SCHEMA EXPECTATION (Return RAW JSON only):
{
  "version": "v1",
  "records": [
    {
      "target_node": { "uuid": "...", "public_id": "...", "file_path": "...", "node_type": "..." },
      "role": "primary_target", 
      "status": "approved",
      "action": "replace_node",
      "new_code": "...",
      "rationale": "..."
    }
  ]
}`, string(planYaml))

	if feedback != "" {
		basePrompt += fmt.Sprintf("\n\n--- COMPILER FEEDBACK ---\nYOUR PREVIOUS ATTEMPT FAILED VALIDATION:\n%s\n\nPlease correct the code while maintaining the exact node identities.", feedback)
	}

	return basePrompt, nil
}

// --- LLM API Client (Stays same) ---

type PlanConfig struct {
	APIKey       string
	PlannerModel string
	BaseURL      string
}

func loadPlanConfig(cmd *cobra.Command) (PlanConfig, error) {
	_ = godotenv.Load()
	apiKey := os.Getenv("SAAYN_DEEP_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("SAAYN_FAST_API_KEY")
		if apiKey == "" {
			return PlanConfig{}, fmt.Errorf("SAAYN_DEEP_API_KEY or SAAYN_FAST_API_KEY environment variable is required")
		}
	}
	model := os.Getenv("SAAYN_DEEP_MODEL")
	if model == "" {
		model = "gemini-2.5-pro"
	}
	baseURL := os.Getenv("SAAYN_DEEP_BASE_URL")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	return PlanConfig{APIKey: apiKey, PlannerModel: model, BaseURL: baseURL}, nil
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func callGeminiPlanner(ctx context.Context, cfg PlanConfig, prompt string) (string, error) {
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", cfg.BaseURL, cfg.PlannerModel, cfg.APIKey)
	payload := map[string]any{
		"contents":         []map[string]any{{"parts": []map[string]any{{"text": prompt}}}},
		"generationConfig": map[string]any{"responseMimeType": "application/json"},
	}
	payloadBytes, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 60 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var result geminiGenerateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}
