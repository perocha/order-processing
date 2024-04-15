package database

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain/order"
)

// Database represents the interface for interacting with the database.
type Database interface {
	OrderRepository() OrderRepository
}

// OrderRepository represents the interface for interacting with order data in the database.
type OrderRepository interface {
	CreateOrder(ctx context.Context, order order.Order) error
	UpdateOrder(ctx context.Context, order order.Order) error
	DeleteOrder(ctx context.Context, id string, partitionKey string) error
}
