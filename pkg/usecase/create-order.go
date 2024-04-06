package usecase

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain"
)

type CreateOrder struct {
	// Define your dependencies such as repositories
}

func (co *CreateOrder) Execute(ctx context.Context, order domain.Order) error {
	// Implement your business logic for creating an order
	return nil
}
