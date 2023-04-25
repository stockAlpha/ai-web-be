package openaiapi

import "github.com/sashabaranov/go-openai"

type ChatCompletionRequest struct {
	Model            string                         `json:"model"`
	Messages         []openai.ChatCompletionMessage `json:"messages"`
	MaxTokens        int                            `json:"max_tokens,omitempty"`
	Temperature      float32                        `json:"temperature,omitempty"`
	Stream           bool                           `json:"stream,omitempty"`
	FrequencyPenalty float32                        `json:"frequency_penalty,omitempty"`
	Role             string                         `json:"role,omitempty"` // 角色
}
