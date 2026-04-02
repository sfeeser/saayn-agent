package saayn

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/saayn-agent/internal/genome/surgery"
)

var (
	planInputFile  string
	planOutputFile string
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate and auto-refine a strict code patch from a surgery plan",
	RunE:  runPlan,
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().StringVarP(&planInputFile, "file", "f", "surgery.yaml", "Path to the input surgery plan")
	planCmd.Flags().StringVarP(&planOutputFile, "output", "o", "patch.yaml", "Path for the output patch file")
}

func runPlan(cmd *cobra.Command, args []string) error {
	fmt.Printf("🧠 Initializing AI Planner for %s...\n", planInputFile)

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

	if plan.Version == "" || plan.Target.PublicID == "" || len(plan.Context) == 0 {
		return fmt.Errorf("invalid surgery plan: missing version, target, or context")
	}

	// Read original source for in-memory splice validation
	originalBytes, err := os.ReadFile(plan.Target.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read original file for validation context: %w", err)
	}

	maxRetries := 3
	var metadata surgery.PatchMetadata
	metadata.Status = "failed_validation"

	var finalPatch surgery.SurgeryPatch
	feedback := ""

	// --- THE AUTO-REFINE LOOP ---
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("🚀 Generating patch (Attempt %d/%d)...\n", attempt, maxRetries)
		metadata.Attempts = attempt

		prompt, err := buildPrompt(plan, feedback)
		if err != nil {
			return fmt.Errorf("failed to build prompt: %w", err)
		}

		jsonResponse, err := callGeminiPlanner(cmd.Context(), cfg, prompt)
		if err != nil {
			return fmt.Errorf("LLM planning failed: %w", err)
		}

		fmt.Println("👀 Performing Code Review + AST splice analysisx...")
		var patch surgery.SurgeryPatch
		// Using strict JSON unmarshal now that schema has json tags
		if err := json.Unmarshal([]byte(jsonResponse), &patch); err != nil {
			recordAttempt(&metadata, attempt, "json_parse", err)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: json_parse\nError: %v\n\nPlease output strictly valid JSON matching the schema.", err)
			continue
		}
		patch.Version = surgery.SurgeryPatchVersion

		// 1. Strict Schema & Target Drift Validation
		if err := validateDrift(plan, patch); err != nil {
			recordAttempt(&metadata, attempt, "schema_validation", err)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: schema_validation\nError: %v\n\nEnsure target_node exactly matches the plan and action is 'replace_node'. Do NOT change the UUID, PublicID, FilePath, or NodeType.", err)
			continue
		}

		// 2. Splice & Format Validation (Uses apply.go's strict logic!)
		// spliceAST already formats via format.Source and resolves imports via imports.Process.
		splicedSrc, err := spliceAST(plan.Target.FilePath, plan.Target.PublicID, originalBytes, patch.NewCode)
		if err != nil {
			recordAttempt(&metadata, attempt, "splice_file", err)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: splice_file\nError: %v\n\nPlease fix the syntax so it can be cleanly formatted and spliced into the AST.", err)
			continue
		}

		// 3. Complete File Parse Validation
		fset := token.NewFileSet()
		if _, err := parser.ParseFile(fset, plan.Target.FilePath, splicedSrc, parser.AllErrors); err != nil {
			recordAttempt(&metadata, attempt, "parse_file", err)
			feedback = fmt.Sprintf("VALIDATION FAILURE:\nStage: parse_file\nError: %v\n\nThe resulting file fails to parse. Check for unbalanced braces or invalid declarations.", err)
			continue
		}

		// SUCCESS! The patch is pristine.
		fmt.Println("✅ Validation passed! Syntax and context are sound.")
		recordAttempt(&metadata, attempt, "validation_passed", nil)
		metadata.Status = "approved"
		finalPatch = patch
		break
	}

	if metadata.Status != "approved" {
		return fmt.Errorf("🛑 Failed to generate a valid patch after %d attempts. The LLM could not resolve the compiler errors", maxRetries)
	}

	finalPatch.Metadata = &metadata

	patchYamlBytes, err := yaml.Marshal(&finalPatch)
	if err != nil {
		return fmt.Errorf("failed to marshal patch to YAML: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(planOutputFile), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(planOutputFile, patchYamlBytes, 0644); err != nil {
		return fmt.Errorf("failed to save patch: %w", err)
	}

	fmt.Printf("💾 Approved surgery patch saved at: %s\n", planOutputFile)
	return nil
}

// --- Validation Helpers ---

func validateDrift(plan surgery.SurgeryPlan, patch surgery.SurgeryPatch) error {
	if patch.TargetNode.UUID != plan.Target.UUID {
		return fmt.Errorf("target UUID mismatch")
	}
	if patch.TargetNode.PublicID != plan.Target.PublicID {
		return fmt.Errorf("target PublicID mismatch")
	}
	if patch.TargetNode.FilePath != plan.Target.FilePath {
		return fmt.Errorf("target file_path mismatch")
	}
	if patch.TargetNode.NodeType != plan.Target.NodeType {
		return fmt.Errorf("target node_type mismatch")
	}
	if patch.Action != surgery.ActionReplaceNode {
		return fmt.Errorf("unsupported action: %s", patch.Action)
	}
	if strings.TrimSpace(patch.NewCode) == "" {
		return fmt.Errorf("new_code is empty")
	}
	return nil
}

func recordAttempt(metadata *surgery.PatchMetadata, attempt int, stage string, err error) {
	rec := surgery.ReviewAttempt{
		Attempt: attempt,
		Stage:   stage,
	}
	if err != nil {
		rec.Errors = []string{err.Error()}
	}
	metadata.ReviewHistory = append(metadata.ReviewHistory, rec)
}

// --- Prompt Engineering ---

func buildPrompt(plan surgery.SurgeryPlan, feedback string) (string, error) {
	planYaml, err := yaml.Marshal(plan)
	if err != nil {
		return "", fmt.Errorf("failed to marshal plan for prompt: %w", err)
	}

	basePrompt := fmt.Sprintf(`You are an expert Go developer and AST surgeon. 
Your task is to generate a strict code patch based on the provided surgery plan.

SURGERY PLAN:
%s

INSTRUCTIONS:
1. Analyze the "intent".
2. Review the "context" nodes (which include the primary target and any callers in the blast radius).
3. Generate the modified "new_code" for the primary target node ONLY. Do not write code for the callers at this time.
4. Preserve the target node's public identity and function signature unless the surgery intent explicitly requires changing them.
5. Do not wrap the response in markdown fences. Return only raw JSON.
6. The target_node fields must exactly match the plan target. Do not invent or modify UUID, public_id, file_path, node_type, or logic_hash.
7. Output your response strictly as a JSON object matching this schema:
{
  "version": "v1",
  "target_node": {
    "uuid": "exact uuid from the target in the plan",
    "public_id": "exact public_id from the target in the plan",
    "file_path": "exact file_path from the target in the plan",
    "node_type": "exact node_type from the target in the plan",
    "logic_hash": "exact logic_hash from the target in the plan"
  },
  "action": "replace_node",
  "new_code": "the complete, newly modified raw Go code for the target node",
  "rationale": "A brief explanation of the change"
}`, string(planYaml))

	if feedback != "" {
		basePrompt += fmt.Sprintf("\n\n--- COMPILER FEEDBACK ---\nYOUR PREVIOUS ATTEMPT FAILED VALIDATION:\n%s\n\nPlease correct the error while preserving the same target identity and generate a new JSON patch.", feedback)
	}

	return basePrompt, nil
}

// --- LLM API Client ---

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
		baseURL = os.Getenv("SAAYN_FAST_BASE_URL")
		if baseURL == "" {
			baseURL = "https://generativelanguage.googleapis.com/v1beta"
		}
	}

	return PlanConfig{
		APIKey:       apiKey,
		PlannerModel: model,
		BaseURL:      baseURL,
	}, nil
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
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"responseMimeType": "application/json",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		msg := strings.TrimSpace(string(body))
		if len(msg) > 500 {
			msg = msg[:500] + "... [truncated]"
		}
		return "", fmt.Errorf("API error: HTTP %d: %s", resp.StatusCode, msg)
	}

	var result geminiGenerateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no candidates or parts returned from LLM")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
