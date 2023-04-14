package db

import (
	"time"

	"gorm.io/gorm"
)

type RechargeKey struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	RechargeKey string    `gorm:"column:recharge_key" json:"recharge_key"`
	Type        uint8     `gorm:"column:type" json:"type"`     // 1代表30积分，2代表100积分，3代表500积分
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

func (r *RechargeKey) UpdateRechargeKey(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(r.TableName())
	} else {
		db = db.Table(r.TableName())
	}
	err := db.Where("id = ?", r.ID).Updates(r).Error
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
