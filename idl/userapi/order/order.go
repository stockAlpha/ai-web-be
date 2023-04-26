package order

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"stock-web-be/dao/db"
	"time"
)

func AddOrder(userId uint64, amount decimal.Decimal, productInfo string, transaction *gorm.DB) (uint64, error) {
	order := &db.Order{
		UserId:      userId,
		OrderType:   1,
		Status:      1,
		Amount:      amount,
		ProductInfo: productInfo,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	err := order.InsertOrder(transaction)
	if err != nil {
		return 0, err
	}
	return order.ID, nil
}

func GetOrderById(orderId string) (*db.Order, error) {
	order := &db.Order{}
	err := order.GetOrderById(orderId)
	if err != nil {
		return nil, err
	}
	if order.ID == 0 {
		return nil, nil
	}
	return order, nil
}

func GetOrderByUserId(userId uint64) (*db.Order, error) {
	order := &db.Order{}
	err := order.GetOrderByUserId(userId)
	if err != nil {
		return nil, err
	}
	if order.ID == 0 {
		return nil, nil
	}
	return order, nil
}

func UpdateOrderStatus(orderId string, status int, transaction *gorm.DB) error {
	order := &db.Order{
		OrderId: orderId,
		Status:  status,
	}
	return order.UpdateOrderStatus(transaction)
}
