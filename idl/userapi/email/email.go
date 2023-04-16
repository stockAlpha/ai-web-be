package email

type ActivityRequest struct {
	TemplateId int `json:"template_id" binding:"required"` // 模版id
	Amount     int `json:"amount" binding:"required"`      // 对积分低于多少的人群发送
}
