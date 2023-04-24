package db

import (
	"time"

	"gorm.io/gorm"
)

//	type RecordResp struct {
//		Active  int `json:"active"`
//		History []struct {
//			Title string `json:"title"`
//			Uuid  int    `json:"uuid"`
//		} `json:"history"`
//		Chat []struct {
//			Uuid int `json:"uuid"`
//			Data []struct {
//				DateTime            string `json:"dateTime"`
//				Text                string `json:"text"`
//				Inversion           bool   `json:"inversion"`
//				IsImage             bool   `json:"isImage,omitempty"`
//				Error               bool   `json:"error"`
//				ConversationOptions *struct {
//					ParentMessageId string `json:"parentMessageId"`
//				} `json:"conversationOptions"`
//				RequestOptions struct {
//					Prompt  string `json:"prompt"`
//					Options *struct {
//						IsImage         bool   `json:"isImage"`
//						ParentMessageId string `json:"parentMessageId,omitempty"`
//					} `json:"options"`
//				} `json:"requestOptions"`
//				Loading bool          `json:"loading,omitempty"`
//				Images  []interface{} `json:"images,omitempty"`
//			} `json:"data"`
//		} `json:"chat"`
//	}
type ChatRecord struct {
	ID            uint64         `json:"id" form:"id" gorm:"primaryKey"`
	UserID        uint64         `json:"user_id" form:"user_id" gorm:"column:user_id"`
	Uuid          int            `json:"uuid" form:"uuid" gorm:"column:uuid"`
	Data          string         `json:"data" form:"data" gorm:"column:data;type:text;"`
	DecodeVersion int            `json:"decode_version" form:"decode_version" gorm:"column:decode_version;"` //0代表没有加密
	DeletedAt     gorm.DeletedAt `json:"deleted_at" form:"deleted_at"`
	CreatedAt     time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" form:"created_at"`
}
type ChatHistory struct {
	ID            uint64         `json:"id" form:"id" gorm:"primaryKey"`
	Uuid          int            `json:"uuid" form:"uuid" gorm:"column:uuid"`
	Title         string         `json:"title" form:"title" gorm:"column:title"`
	DecodeVersion int            `json:"decode_version" form:"decode_version" gorm:"column:decode_version;"` //0代表没有加密
	DeletedAt     gorm.DeletedAt `json:"deleted_at" form:"deleted_at"`
	CreatedAt     time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" form:"created_at"`
}

func (c *ChatRecord) TableName() string {
	return "chat_record"
}
func (c *ChatHistory) TableName() string {
	return "chat_history"
}
func InsertRecordAndHistory(records []ChatRecord, historys []ChatHistory) (err error) {
	//batchLen 分批插入个数
	batchLen := 100
	if len(historys) != 0 {
		err = DbIns.CreateInBatches(historys, batchLen).Error
		if err != nil {
			// todo 打印一下
		}
	}
	if len(records) != 0 {
		err = DbIns.CreateInBatches(records, batchLen).Error
		if err != nil {
			// todo 打印一下
		}
	}
	return
}
