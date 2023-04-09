package integral

type RechargeRequest struct {
	Key string `json:"key" binding:"required"`
}

type BatchGenerateKeyRequest struct {
	Count int   `json:"count" default:"10"`
	Type  uint8 `json:"type" default:"1"` // 1代表100积分，2代表500积分，3代表1000积分
}
