package db

import (
	"errors"
	"time"
)

type VerificationCode struct {
	ID              uint64    `gorm:"primary_key" json:"id"`
	SendSubjectName string    `gorm:"column:send_subject_name" json:"send_subject_name"`
	SendSubjectType int       `gorm:"column:send_subject_type" json:"send_subject_type"`
	Code            string    `gorm:"column:code" json:"code"`
	ExpireTime      time.Time `gorm:"column:expire_time" json:"expire_time"`
	CreatedTime     time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime     time.Time `gorm:"column:updated_time" json:"updated_time"`
}

func (code *VerificationCode) TableName() string {
	return "verification_code"
}

func (code *VerificationCode) InsertCode() error {
	db := DbIns.Table(code.TableName())

	err := db.Create(code).Error
	if err != nil {
		return err
	}
	if code.ID == 0 {
		return errors.New("failed to insert into users table")
	}
	return nil
}

func (code *VerificationCode) GetSendCodeByEmailAndCode(email string, realCode string) ([]*VerificationCode, error) {
	db := DbIns.Table(code.TableName())
	var list []*VerificationCode
	err := db.Where("send_subject_name = ?", email).
		Where("code = ?", realCode).
		Where("expire_time > ?", time.Now()).
		Scan(&list).Error

	return list, err
}
