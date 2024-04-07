package domain

import (
	"context"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) error
	UpdateOrder(ctx context.Context, order Order) error
	DeleteOrder(ctx context.Context, orderID string) error
}
