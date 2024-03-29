package db

import (
	"gorm.io/gorm"
	"time"
)

type RolePermission struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (rolePermission *RolePermission) TableName() string {
	return "role_permissions"
}

func (rolePermission *RolePermission) insertRolePermission(db *gorm.DB) error {
	err := db.Create(rolePermission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
