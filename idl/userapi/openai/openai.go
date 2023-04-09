package openai

type CompletionsRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	TopP        int       `json:"top_p"`
	N           int       `json:"n"`
	Stream      bool      `json:"stream"`
}

type CompletionsResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"choices"`
		FinishReason string  `json:"finish_reason"`
	}
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
