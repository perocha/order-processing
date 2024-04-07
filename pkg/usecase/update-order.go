package usecase

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/domain"
)

type UpdateOrder struct {
	// Define your dependencies such as repositories
}

func (uo *UpdateOrder) Execute(ctx context.Context, order domain.Order) error {
	log.Printf("UpdateOrder::Execute::OrderID=%v", order.OrderID)
	return nil
}
