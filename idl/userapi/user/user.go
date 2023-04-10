package user

type SendVerificationCodeRequest struct {
	Type  string `json:"type" default:"email"` // 可选字段，默认为email
	Email string `json:"email" binding:"required"`
}

type RegisterRequest struct {
	Type     string `json:"type" default:"email"` // 可选字段，默认为email
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	TenantId uint64 `json:"tenantId" default:"1"`    // 租户id，默认为1
	Code     string `json:"code" binding:"required"` // 验证码
}

type LoginRequest struct {
	Type     string `json:"type" default:"email"` // 可选字段，默认为email
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ProfileResponse struct {
	Email    string `json:"email"`
	NickName string `json:"nickName"`
}
