package openai

type ChatRecordRequest struct {
	Limit  int    `json:"limit" form:"limit" `
	Offset int    `json:"offset" form:"offset" `
	UUID   int    `json:"uuid" form:"uuid" `
	UserID uint64 `json:"user_id" form:"user_id" `
}
type ChatRecordResponse struct {
	Active int              `json:"active" `
	Chat   []ChatRecordChat `json:"chat" `
}
type ChatRecordChat struct {
	UUID int                  `json:"uuid" `
	Data []ChatRecordChatData `json:"data"`
}
type ChatRecordChatData struct {
	DateTime       string `json:"dateTime" `
	Inversion      bool   `json:"inversion" `
	IsImage        bool   `json:"isImage" `
	RequestOptions struct {
		Options struct {
			IsImage bool `json:"isImage" `
		} `json:"options" `
		Prompt string `json:"prompt" `
	} `json:"requestOptions" `
	Text string `json:"text" `
}
