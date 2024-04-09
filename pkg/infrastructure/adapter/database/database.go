package database

import (
	"github.com/perocha/order-processing/pkg/domain/order"
)

type Database interface {
	Init() error
	CreateOrder(order order.Order) error
	DeleteOrder(orderID string) error
	UpdateOrder(order order.Order) error
}
