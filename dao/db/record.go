package db

import (
	"errors"
	"time"

	"stock-web-be/idl/openai"

	"gorm.io/gorm"
)

type ChatRecord struct {
	ID            uint64         `json:"id" form:"id" gorm:"primaryKey"`
	UserID        uint64         `json:"user_id" form:"user_id" gorm:"column:user_id"`
	UUID          int            `json:"uuid" form:"uuid" gorm:"column:uuid"`
	MessageID     string         `json:"message_id" form:"message_id" gorm:"column:message_id"`
	DataType      int            `json:"data_type" form:"data_type" gorm:"column:data_type"` //0代表prompt，1代表文本，2代表图片
	Prompt        string         `json:"prompt" form:"prompt" gorm:"column:prompt;type:text;"`
	Data          string         `json:"data" form:"data" gorm:"column:data;type:text;"`
	DecodeVersion int            `json:"decode_version" form:"decode_version" gorm:"column:decode_version;"` //0代表没有加密
	DeletedAt     gorm.DeletedAt `json:"deleted_at" form:"deleted_at"`
	CreatedAt     time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" form:"created_at"`
}

func (c *ChatRecord) DbToOpenAIData() (data openai.ChatRecordChatData) {
	if c.Prompt == "" {
		data.RequestOptions.Prompt = c.Data
		data.Inversion = true
	} else {
		data.RequestOptions.Prompt = c.Prompt
		data.Inversion = false
	}
	data.DateTime = c.CreatedAt.Format("2006/01/02 15:04:05")
	if c.DataType == 2 {
		data.IsImage = true
		data.RequestOptions.Options.IsImage = true
	} else {
		data.IsImage = false
		data.RequestOptions.Options.IsImage = false
	}
	data.Text = c.Data
	return
}
func (c *ChatRecord) TableName() string {
	return "chat_record"
}
func InsertRecord(records []ChatRecord) (err error) {
	//batchLen 分批插入个数
	batchLen := 100
	if len(records) != 0 {
		err = DbIns.CreateInBatches(records, batchLen).Error
		if err != nil {
			// todo 打印一下
		}
	}
	return
}
func FindRecord(req openai.ChatRecordRequest) (records []ChatRecord, err error) {
	//todo 暂时不做分页
	if req.UserID == 0 {
		return nil, errors.New("need userid")
	}
	findRecord := ChatRecord{UserID: req.UserID, UUID: req.UUID}
	err = DbIns.Table(findRecord.TableName()).Order("created_at desc").Find(&records, findRecord).Error
	return
}
