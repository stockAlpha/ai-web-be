package db

import (
	"gorm.io/gorm"
	"time"
)

type UserRole struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	UserId      uint64    `gorm:"column:user_id" json:"user_id"`
	RoleId      uint64    `gorm:"column:role_id" json:"role_id"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (userRole *UserRole) TableName() string {
	return "user_role"
}

func (userRole *UserRole) insertUserRole(db *gorm.DB) error {
	err := db.Create(userRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
