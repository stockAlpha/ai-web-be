package integral

type RechargeRequest struct {
	Key string `json:"key" binding:"required"`
}
