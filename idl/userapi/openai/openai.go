package openai

type CompletionsRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type CompletionsResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}
