package integral

type RechargeRequest struct {
	Key string `json:"key" binding:"required"`
}

type RecordRequest struct {
	Type  string `json:"type" binding:"required"`  // 计费类型，chat/image/audio
	Model string `json:"model" binding:"required"` // 使用模型
	Size  int    `json:"size" binding:"required"`  // 大小，chat为字数，image为尺寸，音频为时长
}

type BatchGenerateKeyRequest struct {
	Count int   `json:"count" default:"10"`
	Type  uint8 `json:"type" default:"1"` // 1代表100积分，2代表500积分，3代表1000积分
}
