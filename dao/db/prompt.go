package db

import (
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Prompt struct {
	ID         uint64          `gorm:"primary_key" json:"id"`
	Name       string          `gorm:"column:name" json:"name"` // 提示名称
	UserId     uint64          `gorm:"column:user_id" json:"user_id"`
	PromptType int             `gorm:"column:prompt_type" json:"prompt_type"` // 提示类型，role/tool
	Amount     decimal.Decimal `gorm:"column:amount" json:"amount"`           // 订单金额
	Status     int             `gorm:"column:status" json:"status"`           // 提示状态，0已上线，1审核中，2已下线
	CreateTime time.Time       `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time       `gorm:"column:update_time" json:"update_time"`
}

func (order *Prompt) TableName() string {
	return "order"
}

func (order *Prompt) InsertOrder(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(order.TableName())
	}

	err := db.Create(order).Error
	if err != nil {
		return err
	}
	if order.ID == 0 {
		return errors.New("failed to insert into users table")
	}
	return nil
}

func (order *Prompt) UpdateOrderStatus(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(order.TableName())
	}
	updateMap := map[string]interface{}{}
	updateMap["status"] = order.Status
	updateMap["update_time"] = time.Now()
	return db.Table(order.TableName()).Where("id = ?", order.ID).
		Updates(updateMap).Error
}

func (order *Prompt) GetOrderById(id uint64) error {
	db := DbIns.Table(order.TableName())

	err := db.Table(order.TableName()).
		Where("id = ?", id).
		Find(order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
