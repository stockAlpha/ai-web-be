package user

type SendVerificationCodeRequest struct {
	Type  string `json:"type" default:"email"` // 可选字段，默认为email
	Email string `json:"email" binding:"required"`
}

type RegisterRequest struct {
	Type       string `json:"type" default:"email"`        // 可选字段，默认为email
	Email      string `json:"email" binding:"required"`    // 邮箱
	Password   string `json:"password" binding:"required"` // 密码
	Code       string `json:"code" binding:"required"`     // 验证码
	InviteCode string `json:"inviteCode"`                  // 邀请码
}

type LoginRequest struct {
	Type     string `json:"type" default:"email"` // 可选字段，默认为email
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ProfileResponse struct {
	Email      string `json:"email"`      // 邮箱
	NickName   string `json:"nickName"`   // 昵称
	Avatar     string `json:"avatar"`     // 头像
	InviteCode string `json:"inviteCode"` // 邀请码
	Integral   int    `json:"integral"`   // 用户当前积分
}

type ProfileRequest struct {
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
}
