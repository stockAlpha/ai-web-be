package db

import (
	"errors"
	"time"
)

type InviteRelation struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	FromUserId uint64    `gorm:"column:from_user_id" json:"from_user_id"`
	ToUserId   uint64    `gorm:"column:to_user_id" json:"to_user_id"`
	InviteCode string    `gorm:"column:invite_code" json:"invite_code"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (i *InviteRelation) TableName() string {
	return "invite_relation"
}

func (i *InviteRelation) InsertRelation() error {
	db := DbIns.Table(i.TableName())

	err := db.Create(i).Error
	if err != nil {
		return err
	}
	if i.ID == 0 {
		return errors.New("failed to insert into users table")
	}
	return nil
}
