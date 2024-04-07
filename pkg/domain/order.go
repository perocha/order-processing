package domain

import "context"

// Order definition
type Order struct {
	OrderID    string
	ProductID  string
	CustomerID string
	Status     string
}

// Order methods
type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) error
	UpdateOrder(ctx context.Context, order Order) error
	DeleteOrder(ctx context.Context, orderID string) error
}
