package surgeon

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sfeeser/saayn-agent/internal/genome/index"
	"github.com/sfeeser/saayn-agent/pkg/model"
)

// CognitiveResult represents the structured output expected from the Cognitive Audit LLM.
type CognitiveResult struct {
	Pass     bool   `json:"pass"`
	Findings string `json:"findings"`
}

// draftNode generates or remediates a single node using the DEEP LLM tier.
func draftNode(spec model.SpecNode, vision string, physicsErr error, cognitiveFindings string) (string, error) {
	prompt := buildDraftPrompt(spec, vision, physicsErr, cognitiveFindings)

	// Route to the DEEP reasoning tier (e.g., Gemini Pro)
	code, err := callLLM(prompt, "deep")
	if err != nil {
		return "", fmt.Errorf("draftNode: LLM call failed: %w", err)
	}

	return cleanGeneratedCode(code), nil
}

func buildDraftPrompt(spec model.SpecNode, vision string, physicsErr error, cognitiveFindings string) string {
	var sb strings.Builder

	// Breaking the backticks to ensure the Go compiler never panics again
	sb.WriteString(`You are SAAYN, an autonomous Go code synthesis engine.
You must output ONLY valid, raw Go source code. 
NEVER include explanations, markdown, code fences (` + "```" + `), or any other text.

PROJECT VISION:
` + strings.TrimSpace(vision) + `

TARGET NODE:
` + fmt.Sprintf("%+v", spec) + `

CRITICAL REQUIREMENTS:
- The code MUST start with the correct package declaration: package ` + inferPackageName(spec.Path) + `
- Use strict Go idioms. Prefer standard library when possible.
- Make the code clean, idiomatic, and production-ready.
`)

	if physicsErr != nil {
		sb.WriteString("\nPREVIOUS PHYSICS (AST) FAILURE:\n" + physicsErr.Error() + "\nFix this exact error.\n")
	}
	if cognitiveFindings != "" {
		sb.WriteString("\nPREVIOUS COGNITIVE FINDINGS:\n" + cognitiveFindings + "\nAddress these issues strictly.\n")
	}

	sb.WriteString("\nRespond with raw Go code only.\n")
	return sb.String()
}

// cognitiveAudit verifies that the generated code faithfully matches the vision and spec.
func cognitiveAudit(spec model.SpecNode, vision string, code string) (string, bool, error) {
	prompt := fmt.Sprintf(`You are a strict code auditor for the SAAYN Greenfield Protocol.

Review the code below against the project vision and node specification.
Return **only** valid JSON matching this exact schema. No extra text.

{
  "pass": boolean,
  "findings": "brief explanation of any drift or issues, or empty string if perfect"
}

VISION:
%s

SPEC:
%+v

CODE:
%s`, vision, spec, code)

	response, err := callLLM(prompt, "fast")
	if err != nil {
		return "", false, fmt.Errorf("cognitiveAudit: LLM call failed: %w", err)
	}

	cleaned := cleanGeneratedCode(response)
	var result CognitiveResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return "", false, fmt.Errorf("cognitiveAudit: failed to parse JSON: %w (raw: %s)", err, cleaned)
	}

	return result.Findings, result.Pass, nil
}

// MaterializeNode creates the "birth certificate" metadata after successful materialization.
func MaterializeNode(spec model.SpecNode, code string) (*model.Node, error) {
	logicHash, err := computeLogicHash(code)
	if err != nil {
		return nil, fmt.Errorf("failed to compute logic hash: %w", err)
	}

	return &model.Node{
		UUID:            uuid.New().String(),
		PublicID:        fmt.Sprintf("%s[%s]", spec.Name, filepath.Base(spec.Path)),
		Fingerprint:     "", // 🧬 TODO: Compute from AST in next iteration
		LogicHash:       logicHash,
		NodeType:        spec.Type,
		FilePath:        spec.Path,
		BusinessPurpose: spec.Logic,
		LastModified:    time.Now(),
		Version:         1,
		TestHealth:      model.TestHealth{Status: "untested"},
		Dependencies: model.DependencyMap{
			Calls:  []string{},
			Reads:  []string{},
			Writes: []string{},
		},
		// AST and Fset are nil by default, which is fine for the "Birth" phase
	}, nil
}

// -----------------------------------------------------------------------------
// Utilities
// -----------------------------------------------------------------------------

func cleanGeneratedCode(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```go")
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}

func inferPackageName(filePath string) string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(dir)
	if dir == "." || dir == "/" || base == "." {
		return "main"
	}
	return base
}

func callLLM(prompt string, tier string) (string, error) {
	modelType := "pro" // DEEP tier
	if tier == "fast" {
		modelType = "flash" // FAST tier
	}

	response, err := index.GenerateCompletion(prompt, modelType)
	if err != nil {
		return "", fmt.Errorf("%s tier LLM failure: %w", tier, err)
	}
	return response, nil
}

func computeLogicHash(code string) (string, error) {
	if len(code) == 0 {
		return "", fmt.Errorf("cannot hash empty code block")
	}
	normalized := cleanGeneratedCode(code)
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:]), nil
}
