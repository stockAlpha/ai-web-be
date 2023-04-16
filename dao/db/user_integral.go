package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserIntegral struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	UserId      uint64    `gorm:"column:user_id" json:"user_id"`
	Amount      int       `gorm:"column:amount" json:"amount"`
	TotalAmount int       `gorm:"column:total_amount" json:"total_amount"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (u *UserIntegral) TableName() string {
	return "user_integral"
}

func (u *UserIntegral) InsertUserIntegral(db *gorm.DB) error {
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

func (u *UserIntegral) GetUserIntegralByUserId(userId uint64, db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(u.TableName())
	}
	err := db.Table(u.TableName()).
		Where("user_id = ?", userId).
		Find(u).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (u *UserIntegral) AddAmount(amount int, db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(u.TableName())
	}
	updateMap := map[string]interface{}{}
	updateMap["amount"] = gorm.Expr("amount + ?", amount)
	updateMap["total_amount"] = gorm.Expr("total_amount + ?", amount)
	if err := db.Model(u).Updates(updateMap).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserIntegral) SubAmount(amount int) error {
	db := DbIns.Table(u.TableName())
	return db.Transaction(func(tx *gorm.DB) error {
		if u.Amount >= amount {
			if err := tx.Model(u).UpdateColumn("amount", gorm.Expr("amount - ?", amount)).Error; err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("余额不足")
		}
	})
}
