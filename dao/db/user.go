package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	Email      string    `gorm:"column:email" json:"email"`
	Password   string    `gorm:"column:password" json:"password"`
	NickName   string    `gorm:"column:nick_name" json:"nick_name"`
	Avatar     string    `gorm:"column:avatar" json:"avatar"`
	InviteCode string    `gorm:"column:invite_code" json:"invite_code"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) InsertUser(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(user.TableName())
	}

	err := db.Create(user).Error
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("failed to insert into users table")
	}
	return nil
}

func (user *User) UpdateUser() {
	db := DbIns.Table(user.TableName()).Where("id = ?", user.ID)
	updateMap := map[string]interface{}{}
	if user.NickName != "" {
		updateMap["nick_name"] = user.NickName
	}
	if user.Avatar != "" {
		updateMap["avatar"] = user.Avatar
	}
	db.Updates(updateMap)
}

func (user *User) UpdateUserPassword(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(user.TableName())
	}

	return db.Table(user.TableName()).Where("id = ?", user.ID).
		Update("password", user.Password).Error
}

func (user *User) GetUserByEmail(email string) error {
	db := DbIns.Table(user.TableName())

	err := db.Table(user.TableName()).
		Where("email = ?", email).
		Find(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (user *User) GetUserById(id uint64) error {
	db := DbIns.Table(user.TableName())

	err := db.Table(user.TableName()).
		Where("id = ?", id).
		Find(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (user *User) GetUserByInviteCode(inviteCode string, db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(user.TableName())
	}

	err := db.Table(user.TableName()).
		Where("invite_code = ?", inviteCode).
		Find(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
