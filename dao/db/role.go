package db

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (role *Role) TableName() string {
	return "roles"
}

func (role *Role) insertRole(db *gorm.DB) error {
	err := db.Create(role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
