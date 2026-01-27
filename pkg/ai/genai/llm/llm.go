// Package llm provides an LLM client interface.
//
// Usage:
//
//	import "github.com/chris-alexander-pop/system-design-library/pkg/ai/genai/llm"
//
//	client := openai.New("key")
//	resp, err := client.Chat(ctx, []llm.Message{{Role: "user", Content: "Hello"}})
package llm

import "context"

// Role defines the message sender.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleFunction  Role = "function"
)

// Message represents a chat message.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// Response represents the model's reply.
type Response struct {
	Message Message `json:"message"`
	Usage   Usage   `json:"usage"`
}

// Usage tracks token usage.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Options configures the generation.
type Options struct {
	Model       string   `json:"model"`
	Temperature float64  `json:"temperature"`
	MaxTokens   int      `json:"max_tokens"`
	Stop        []string `json:"stop,omitempty"`
}

// Client is the LLM interface.
type Client interface {
	Chat(ctx context.Context, messages []Message, opts Options) (*Response, error)
}
