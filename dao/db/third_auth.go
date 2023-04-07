package db

import (
	"gorm.io/gorm"
	"time"
)

type ThirdAuth struct {
	ID            uint64    `gorm:"primary_key" json:"id"`
	UserId        uint64    `gorm:"column:user_id" json:"user_id"`
	ThirdAuthType int       `gorm:"column:third_auth_type" json:"third_auth_type"`
	ThirdAuthId   string    `gorm:"column:third_auth_id" json:"third_auth_id"`
	Token         string    `gorm:"column:token" json:"token"`
	ExpireTime    time.Time `gorm:"column:expire_time" json:"expire_time"`
	CreatedTime   time.Time `gorm:"column:created_time" json:"created_time"`
	UpdateTime    time.Time `gorm:"column:update_time" json:"update_time"`
}

func (thirdAuth *ThirdAuth) TableName() string {
	return "third_auth"
}

func (thirdAuth *ThirdAuth) insertThirdAuth(db *gorm.DB) error {
	err := db.Create(thirdAuth).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
