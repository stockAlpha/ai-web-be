package db

import (
	"gorm.io/gorm"
	"time"
)

type RechargeKey struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	RechargeKey string    `gorm:"column:recharge_key" json:"recharge_key"`
	Status      uint8     `gorm:"column:status" json:"status"`
	UseAccount  string    `gorm:"column:use_account" json:"use_account"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
}

func (r *RechargeKey) TableName() string {
	return "recharge_key"
}

func (r *RechargeKey) InsertRechargeKey(db *gorm.DB) error {
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
