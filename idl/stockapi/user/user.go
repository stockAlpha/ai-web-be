package user

type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	TenantId uint64 `json:"tenant_id" binding:"required"`
	Code     string `json:"code" binding:"required"`
}
