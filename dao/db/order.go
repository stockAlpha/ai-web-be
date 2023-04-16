package db

import (
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID          uint64          `gorm:"primary_key" json:"id"`
	FromUserId  uint64          `gorm:"column:from_user_id" json:"from_user_id"`
	OrderType   int             `gorm:"column:order_type" json:"order_type"`     // 充值类型，1为积分
	Amount      decimal.Decimal `gorm:"column:amount" json:"amount"`             // 订单金额
	Status      int             `gorm:"column:status" json:"status"`             // 订单状态，1为待支付，2为已支付，3为已取消
	ProductInfo string          `gorm:"column:product_info" json:"product_info"` // 商品信息
	CreateTime  time.Time       `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time       `gorm:"column:update_time" json:"update_time"`
}

func (order *Order) TableName() string {
	return "order"
}

func (order *Order) InsertOrder(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(order.TableName())
	}

	err := db.Create(order).Error
	if err != nil {
		return err
	}
	if order.ID == 0 {
		return errors.New("failed to insert into users table")
	}
	return nil
}

func (order *Order) UpdateOrderStatus(db *gorm.DB) error {
	if db == nil {
		db = DbIns.Table(order.TableName())
	}
	updateMap := map[string]interface{}{}
	updateMap["status"] = order.Status
	updateMap["update_time"] = time.Now()
	return db.Table(order.TableName()).Where("id = ?", order.ID).
		Updates(updateMap).Error
}

func (order *Order) GetOrderById(id uint64) error {
	db := DbIns.Table(order.TableName())

	err := db.Table(order.TableName()).
		Where("id = ?", id).
		Find(order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
