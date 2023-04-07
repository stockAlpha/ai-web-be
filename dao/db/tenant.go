package db

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (tenant *Tenant) TableName() string {
	return "tenant"
}

func (tenant *Tenant) insertUser(db *gorm.DB) error {
	err := db.Create(tenant).Error
	if err != nil {
		return err
	}
	if tenant.ID == 0 {
		return errors.New("failed to insert into tenant table")
	}
	return nil
}

func (tenant *Tenant) getById(db *gorm.DB, id uint64) error {
	err := db.Table(tenant.TableName()).
		Where("id = ?", id).
		Find(tenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
