package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FetchEmbedding calls the Gemini embedContent endpoint.
// It accepts an optional *http.Client; if nil, a default client with a 30s timeout is used.
func FetchEmbedding(ctx context.Context, client *http.Client, apiKey, modelName, baseURL, text string) ([]float32, error) {
	// 1. Pre-flight Validation
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("missing API key")
	}
	if strings.TrimSpace(modelName) == "" {
		return nil, fmt.Errorf("missing model name")
	}
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("missing base URL")
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("embedding input text is empty")
	}

	// 2. URL Construction (More robust path joining)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/models/" + modelName + ":embedContent"

	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	// 3. Request Payload
	reqBody := map[string]any{
		"model": "models/" + modelName,
		"content": map[string]any{
			"parts": []map[string]any{
				{"text": text},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	// 4. Execution
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("build embedding request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Fallback to a default client if none provided
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send embedding request: %w", err)
	}
	defer resp.Body.Close()

	// 5. Response Handling
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read embedding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := strings.TrimSpace(string(body))
		if len(msg) > 500 {
			msg = msg[:500] + "... [truncated]"
		}
		return nil, fmt.Errorf("embedding API error: HTTP %d: %s", resp.StatusCode, msg)
	}

	var res struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("decode embedding response: %w", err)
	}

	// 6. Final Integrity Check
	if len(res.Embedding.Values) == 0 {
		return nil, fmt.Errorf("embedding response contained no values")
	}

	return res.Embedding.Values, nil
}

func GenerateCompletion(prompt string, modelType string) (string, error) {
	// ... your logic to call Gemini or another LLM ...
	return "generated code", nil
}
