package db

import (
	"gorm.io/gorm"
	"time"
)

type RechargeKey struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	RechargeKey string    `gorm:"column:recharge_key" json:"recharge_key"`
	Type        uint8     `gorm:"column:type" json:"type"`     // 1代表100积分，2代表500积分，3代表1000积分
	Status      uint8     `gorm:"column:status" json:"status"` // 0代表未使用，1代表已使用，2代表已失效
	UseAccount  uint64    `gorm:"column:use_account" json:"use_account"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (r *RechargeKey) TableName() string {
	return "recharge_key"
}

func (r *RechargeKey) InsertRechargeKey() error {
	db := DbIns.Table(r.TableName())
	err := db.Create(r).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (r *RechargeKey) UpdateRechargeKey() error {
	db := DbIns.Table(r.TableName()).Where("id = ?", r.ID)
	err := db.Updates(r).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (r *RechargeKey) GetRechargeKey(key string) error {
	db := DbIns.Table(r.TableName())

	err := db.Table(r.TableName()).
		Where("recharge_key = ?", key).
		Find(r).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
