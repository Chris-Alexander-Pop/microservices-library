// Package anthropic provides an Anthropic Chat adapter (Claude).
package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/ai/genai/llm"
	pkgerrors "github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) Chat(ctx context.Context, messages []llm.Message, opts llm.Options) (*llm.Response, error) {
	if opts.Model == "" {
		opts.Model = "claude-3-opus-20240229"
	}

	// Convert system messages to separate field (Anthropic style)
	var system string
	var anthropicMessages []map[string]string

	for _, m := range messages {
		if m.Role == llm.RoleSystem {
			system = m.Content
		} else {
			anthropicMessages = append(anthropicMessages, map[string]string{
				"role":    string(m.Role),
				"content": m.Content,
			})
		}
	}

	body := map[string]interface{}{
		"model":      opts.Model,
		"messages":   anthropicMessages,
		"max_tokens": opts.MaxTokens,
	}
	if system != "" {
		body["system"] = system
	}
	if opts.MaxTokens == 0 {
		body["max_tokens"] = 1024
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, pkgerrors.Internal("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, pkgerrors.Internal("API request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, pkgerrors.Internal("API error", nil)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pkgerrors.Internal("failed to parse response", err)
	}

	content := ""
	if len(result.Content) > 0 {
		content = result.Content[0].Text
	}

	return &llm.Response{
		Message: llm.Message{
			Role:    llm.RoleAssistant,
			Content: content,
		},
		Usage: llm.Usage{
			PromptTokens:     result.Usage.InputTokens,
			CompletionTokens: result.Usage.OutputTokens,
			TotalTokens:      result.Usage.InputTokens + result.Usage.OutputTokens,
		},
	}, nil
}

var _ llm.Client = (*Client)(nil)
