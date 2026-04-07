package saayn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sfeeser/saayn-agent/internal/genome"
	"github.com/sfeeser/saayn-agent/internal/scanner"
	"github.com/spf13/cobra"
)

// AgentTestResponse defines the expected JSON from the LLM
type AgentTestResponse struct {
	TestCode    string `json:"test_code"`
	Explanation string `json:"explanation"`
}

// startSpinner creates a simple, non-blocking CLI loading animation
func startSpinner(msg string) chan bool {
	done := make(chan bool)
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Printf("\r\033[K") // Clear the line when done
				return
			default:
				fmt.Printf("\r%s %s", chars[i], msg)
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return done
}

var genTestCmd = &cobra.Command{
	Use:     "gen-test",
	Aliases: []string{"test"},
	Short:   "Generates and verifies a Go test for the last mutated function",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = godotenv.Load()

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("failed to read --path flag: %w", err)
		}

		apiKey := strings.TrimSpace(os.Getenv("SAAYN_FAST_API_KEY"))
		model := strings.TrimSpace(os.Getenv("SAAYN_FAST_MODEL"))
		baseURL := strings.TrimSpace(os.Getenv("SAAYN_FAST_BASE_URL"))

		if apiKey == "" || model == "" || baseURL == "" {
			return fmt.Errorf("LLM variables missing. Run 'verify-llm-targets'")
		}

		fmt.Println("🏥 Reading Medical History...")

		regManager, err := genome.Load(path + "/genome.json")
		if err != nil {
			return fmt.Errorf("error loading genome.json: %w", err)
		}

		lastUUID := regManager.Registry.LastMutatedUUID
		if lastUUID == "" {
			return fmt.Errorf("no recent surgeries found. Run a mutation first")
		}

		targetNode, exists := regManager.Registry.Nodes[lastUUID]
		if !exists {
			return fmt.Errorf("node %s found in history, but missing from registry", lastUUID)
		}

		fmt.Printf("🔍 Target Acquired: %s\n", targetNode.PublicID)

		liveNodes, err := scanner.FullScan(path)
		if err != nil {
			return fmt.Errorf("error scanning directory: %w", err)
		}

		var rawCode string
		var targetFilePath string // NEW: We will store the exact file path here

		conf := &printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
		for _, ln := range liveNodes {
			if ln.PublicID == targetNode.PublicID {
				var buf strings.Builder
				if err := conf.Fprint(&buf, ln.Fset, ln.AST); err != nil {
					return fmt.Errorf("failed to render live AST for %s: %w", ln.PublicID, err)
				}
				rawCode = buf.String()

				// AST MAGIC: Extract the exact file path from the Go Token parser!
				targetFilePath = ln.Fset.Position(ln.AST.Pos()).Filename
				break
			}
		}

		if rawCode == "" || targetFilePath == "" {
			return fmt.Errorf("could not find live source code for %s. Did the file move?", targetNode.PublicID)
		}

		pkgName := "main"
		if dotIdx := strings.Index(targetNode.PublicID, "."); dotIdx > 0 {
			pkgName = targetNode.PublicID[:dotIdx]
		}

		// Start the visual spinner!
		spinnerDone := startSpinner("🧪 Handing code to the QA Agent (This takes 10-20 seconds)...")

		testResp, err := generateUnitTest(cmd.Context(), apiKey, model, baseURL, targetNode.PublicID, pkgName, targetNode.BusinessPurpose, rawCode)

		// Stop the spinner
		spinnerDone <- true

		if err != nil {
			return fmt.Errorf("\nQA Agent failed: %w", err)
		}

		fset := token.NewFileSet()
		if _, err := parser.ParseFile(fset, "", testResp.TestCode, 0); err != nil {
			return fmt.Errorf("\nAI generated invalid or uncompilable Go code: %w\nExplanation: %s", err, testResp.Explanation)
		}

		if !strings.Contains(testResp.TestCode, "func Test") {
			return fmt.Errorf("\nAI generated valid Go code, but it does not contain a test function")
		}

		// Use the exact path we found in the AST, replacing .go with _test.go
		testFileName := strings.TrimSuffix(targetFilePath, ".go") + "_test.go"

		if err := os.WriteFile(testFileName, []byte(testResp.TestCode), 0644); err != nil {
			return fmt.Errorf("\nfailed to write test file: %w", err)
		}

		fmt.Printf("✅ Valid Go Test written to %s\n", testFileName)
		fmt.Printf("💡 QA Rationale: %s\n", testResp.Explanation)

		fmt.Println("\n🔬 Running test suite to verify immunity...")
		testCmdExec := exec.CommandContext(cmd.Context(), "go", "test", "./...")
		testCmdExec.Dir = path

		output, testErr := testCmdExec.CombinedOutput()
		passed := testErr == nil

		if passed {
			fmt.Println("🟢 ALL TESTS PASSED! The mutation is stable.")
		} else {
			fmt.Println("🔴 TEST FAILED! See output below:")
			fmt.Println(string(output))
		}

		targetNode.TestHealth.HasTest = true
		targetNode.TestHealth.LastPassed = passed
		targetNode.TestHealth.TestedAt = time.Now().Format(time.RFC3339)

		regManager.Registry.Nodes[lastUUID] = targetNode
		if err := regManager.Save(); err != nil {
			return fmt.Errorf("failed to update genome test health: %w", err)
		}

		fmt.Println("💾 Genome Health Record updated.")
		return nil
	},
}

func generateUnitTest(ctx context.Context, apiKey, model, baseURL, publicID, pkgName, purpose, rawCode string) (*AgentTestResponse, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/models/" + url.PathEscape(model) + ":generateContent"
	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	prompt := fmt.Sprintf(`System Context: You are the autonomous QA Engineer for SAAYN.

Task:
Write a complete, passing Go unit test for the provided function. 

Context:
- Target: %s
- Package Name: %s
- Business Purpose: %s

Raw Source Code:
%s

Strict Constraints:
1. Output ONLY raw, valid JSON. Do not include markdown fences.
2. The 'test_code' field must contain a complete, compiling Go file (including 'package %s', standard library imports, and the Test function).
3. Use ONLY the Go standard library (e.g., 'testing', 'reflect'). Do NOT import third-party libraries.
4. Do NOT invent types, structs, or constructors that are not clearly visible in the provided source code.
5. If a full behavioral test is impossible due to missing external dependencies, generate a minimal compile-safe scaffold (e.g., testing that the function exists and doesn't panic on nil inputs) and explain the limitation in the 'explanation' field.

Schema Requirement:
{
  "test_code": "package %s\n\nimport...",
  "explanation": "Strictly ONE short, punchy sentence explaining what the test does. Do not write paragraphs or defensive excuses."
}`, publicID, pkgName, purpose, rawCode, pkgName, pkgName)

	reqBody := map[string]any{
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

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		reason := strings.TrimSpace(string(bodyBytes))
		if len(reason) > 500 {
			reason = reason[:500] + "... [truncated]"
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, reason)
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
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from AI")
	}

	rawText := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)

	var agentResp AgentTestResponse
	if err := json.Unmarshal([]byte(rawText), &agentResp); err != nil {
		if len(rawText) > 1000 {
			rawText = rawText[:1000] + "... [truncated]"
		}
		return nil, fmt.Errorf("failed to parse AI JSON: %w\nRaw Output: %s", err, rawText)
	}

	return &agentResp, nil
}

func init() {
	rootCmd.AddCommand(genTestCmd)
	genTestCmd.Flags().StringP("path", "p", ".", "Path to the project root")
}
