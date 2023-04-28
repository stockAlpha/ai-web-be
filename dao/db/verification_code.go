package db

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type VerificationCode struct {
	ID              uint64    `gorm:"primary_key" json:"id"`
	SendSubjectName string    `gorm:"column:send_subject_name" json:"send_subject_name"`
	SendSubjectType int       `gorm:"column:send_subject_type" json:"send_subject_type"`
	Code            string    `gorm:"column:code" json:"code"`
	ExpireTime      time.Time `gorm:"column:expire_time" json:"expire_time"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"update_time"`
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

func (code *VerificationCode) GetSendCodeByCodeType(email string, realCode string, codeType int) ([]*VerificationCode, error) {
	db := DbIns.Table(code.TableName())
	var list []*VerificationCode
	err := db.Where("send_subject_name = ?", email).
		Where("code = ?", realCode).
		Where("send_subject_type = ?", codeType).
		Where("expire_time > ?", time.Now()).
		Scan(&list).Error

	return list, err
}

func (code *VerificationCode) UpdateByCode(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(code.TableName())
	}

	updateMap := map[string]interface{}{}
	if !code.ExpireTime.IsZero() {
		updateMap["expire_time"] = code.ExpireTime
	}
	updateMap["update_time"] = time.Now()
	return db.Table(code.TableName()).Where("send_subject_name = ?", code.SendSubjectName).
		Where("send_subject_type = ?", code.SendSubjectType).
		Where("code = ?", code.Code).
		Updates(updateMap).
		Error
}
