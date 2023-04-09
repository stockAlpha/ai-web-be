package db

import (
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	TenantId    uint64    `gorm:"column:tenant_id" json:"tenant_id"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (permission *Permission) TableName() string {
	return "permissions"
}

func (permission *Permission) insertPermission(db *gorm.DB) error {
	err := db.Create(permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
