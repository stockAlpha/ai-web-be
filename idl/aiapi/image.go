package aiapi

import "time"

type ImageRequest struct {
	Model     string `json:"model" default:"dall-e2"` // dall-e2/stable-diffusion
	Prompt    string `json:"prompt"`
	N         int    `json:"n" default:"1"`
	Size      string `json:"size" default:"512x512"` // 256x256/512x512/1024x1024
	UserID    uint   `json:"user_id"`
	UUID      int    `json:"uuid"`
	MessageID string `json:"message_id"`
}

type ImageResponseDataInner struct {
	URL string `json:"url,omitempty"`
}

type ReplicateStableDiffusion struct {
	Version string         `json:"version"`
	Input   ReplicateInput `json:"input"`
}

type ReplicateInput struct {
	Prompt          string `json:"prompt"`
	ImageDimensions string `json:"image_dimensions"`
	NumOutputs      int    `json:"num_outputs"`
}

type ReplicateResponse struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Urls    struct {
		Get    string `json:"get"`
		Cancel string `json:"cancel"`
	} `json:"urls"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	Status      string    `json:"status"`
	Input       any       `json:"input"`
	Output      []string  `json:"output"`
	Error       string    `json:"error"`
	Logs        string    `json:"logs"`
	Metrics     struct {
		PredictTime float64 `json:"predict_time"`
	} `json:"metrics"`
}
