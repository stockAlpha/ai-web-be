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
	NickName     string       `json:"nickName"`               // 昵称
	Avatar       string       `json:"avatar"`                 // 头像
	CustomConfig CustomConfig `json:"customConfig,omitempty"` // 自定义配置
}

type CustomConfig struct {
	ChatConfig  ChatConfig  `json:"chatConfig"`  // 聊天配置
	ImageConfig ImageConfig `json:"imageConfig"` // 图片配置
}

type ChatConfig struct {
	Model            string  `json:"model,omitempty"`                        // 模型
	Temperature      float32 `json:"temperature,omitempty" default:"1"`      // 随机性0-2,默认为1
	FrequencyPenalty float32 `json:"frequencyPenalty,omitempty" default:"0"` // 话题新鲜度,-2.0-2.0,默认为0
}

type ImageConfig struct {
	N    int    `json:"n,omitempty" default:"1"` // 返回几张图，默认1张
	Size string `json:"size,omitempty"`          // 图片大小,256x256/512x512/1024x1024
}

type FeedbackRequest struct {
	FeedbackType int    `json:"feedbackType" binding:"required"` // 反馈类型: 1-问题反馈 2-功能建议 3-咨询 4-其他
	Content      string `json:"content" binding:"required"`      // 反馈内容
}

type SendPasswordVerificationCodeRequest struct {
	SubjectType int    `json:"subjectType" default:"101"` // 可选字段，默认为userapi.ChangePasswordMailCode
	SubjectName string `json:"subjectName" binding:"required"`
}

type ChangePasswordRequest struct {
	SubjectType      int    `json:"subjectType" default:"101""` // 可选字段，默认为userapi.ChangePasswordMailCode
	SubjectName      string `json:"subjectName" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
}
