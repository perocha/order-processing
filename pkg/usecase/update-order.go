package usecase

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain"
)

type UpdateOrder struct {
	// Define your dependencies such as repositories
}

func (uo *UpdateOrder) Execute(ctx context.Context, order domain.Order) error {
	// Implement your business logic for updating an order
	return nil
}
