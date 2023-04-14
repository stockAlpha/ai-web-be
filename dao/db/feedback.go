package db

import (
	"errors"
	"time"
)

type Feedback struct {
	ID           uint64    `gorm:"primary_key" json:"id"`
	FromUserId   uint64    `gorm:"column:from_user_id" json:"from_user_id"`
	FeedbackType int       `gorm:"column:feedback_type" json:"feedback_type"`
	Content      string    `gorm:"column:content" json:"content"`
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`
}

func (f *Feedback) TableName() string {
	return "feedback"
}

func (f *Feedback) InsertFeedback() error {
	db := DbIns.Table(f.TableName())

	err := db.Create(f).Error
	if err != nil {
		return err
	}
	if f.ID == 0 {
		return errors.New("failed to insert into users table")
	}
	return nil
}
