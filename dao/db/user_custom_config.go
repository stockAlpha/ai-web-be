package db

import (
	"errors"
	"gorm.io/gorm"
)

type UserCustomConfig struct {
	ID          uint64 `gorm:"primary_key" json:"id"`
	UserId      uint64 `gorm:"column:user_id" json:"user_id"`
	ChatConfig  string `gorm:"column:chat_config" json:"chat_config"`
	ImageConfig string `gorm:"column:image_config" json:"image_config"`
}

func (u *UserCustomConfig) TableName() string {
	return "user_custom_config"
}

func (u *UserCustomConfig) InsertUserCustomConfig(db *gorm.DB) error {
	// todo 完善json转化
	if db == nil {
		db = DbIns.Table(u.TableName())
	}
	err := db.Create(u).Error
	if err != nil {
		return err
	}
	if u.ID == 0 {
		return errors.New("failed to insert into user_integral table")
	}
	return nil
}
