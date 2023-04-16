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

type FeedbackRequest struct {
	FeedbackType int    `json:"feedbackType" binding:"required"` // 反馈类型: 1-问题反馈 2-功能建议 3-咨询 4-其他
	Content      string `json:"content" binding:"required"`      // 反馈内容
}

type SendPasswordVerificationCodeRequest struct {
	SubjectType int    `json:"subjectType" default:101` // 可选字段，默认为userapi.ChangePasswordMailCode
	SubjectName string `json:"subjectName" binding:"required"`
}

type ChangePasswordRequest struct {
	SubjectType      int    `json:"subjectType" default:101` // 可选字段，默认为userapi.ChangePasswordMailCode
	SubjectName      string `json:"subjectName" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
}
