package aiapi

import "time"

type ImageRequest struct {
	Model  string `json:"model" default:"dall-e2"` // dall-e2/stable-diffusion
	Prompt string `json:"prompt"`
	N      int    `json:"n" default:"1"`
	Size   string `json:"size" default:"512x512"` // 256x256/512x512/1024x1024
}

type ImageResponseDataInner struct {
	URL    string `json:"url"`
	TaskId string `json:"taskId"`
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

type MjProxySubmit struct {
	Action string `json:"action"`
	Index  int    `json:"index"`
	Prompt string `json:"prompt"`
	TaskId string `json:"taskId"`
}

type MjProxySubmitRes struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

type MjProxyGetRes struct {
	Action      string `json:"action"`
	Id          string `json:"id"`
	Prompt      string `json:"prompt"`
	PromptEn    string `json:"promptEn"`
	Description string `json:"description"`
	State       string `json:"state"`
	SubmitTime  int64  `json:"submitTime"`
	StartTime   int64  `json:"startTime"`
	FinishTime  int64  `json:"finishTime"`
	ImageUrl    string `json:"imageUrl"`
	Status      string `json:"status"`
	FailReason  string `json:"failReason"`
}

type MjProxyOperate struct {
	Action string `json:"action"` // IMAGINE:出图；UPSCALE:选中放大；VARIATION：选中其中的一张图，生成四张相似的,可用值:IMAGINE,UPSCALE,VARIATION,RESET,DESCRIBE
	Index  int    `json:"index"`  // 序号: action 为 UPSCALE,VARIATION 必传，1-4
	TaskId string `json:"taskId"` // 返回的任务id
}
