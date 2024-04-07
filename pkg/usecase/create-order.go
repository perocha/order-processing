package usecase

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/domain"
)

type CreateOrder struct {
	// Define your dependencies such as repositories
}

func (co *CreateOrder) Execute(ctx context.Context, order domain.Order) error {
	log.Printf("CreateOrder::Execute::OrderID=%v", order.OrderID)
	return nil
}
