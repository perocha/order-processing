package database

import (
	"github.com/perocha/order-processing/pkg/domain/order"
)

// Database represents the interface for interacting with the database.
type Database interface {
    OrderRepository() OrderRepository
}

// OrderRepository represents the interface for interacting with order data in the database.
type OrderRepository interface {
    CreateOrder(order order.Order) error
    GetOrder(orderID string) (order.Order, error)
    UpdateOrder(order order.Order) error
    DeleteOrder(orderID string) error
}
