package aiapi

type ModelItem struct {
	Type  string   `json:"type"`  // chat/image
	Model []string `json:"model"` // 模型名称
}
