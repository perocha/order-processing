package usecase

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/domain"
)

type CreateOrder struct {
	orderRepo domain.OrderRepository
}

func (co *CreateOrder) Execute(ctx context.Context, order domain.Order) error {
	log.Printf("CreateOrder::Execute::OrderID=%v", order.OrderID)

	// Create the order in the order repository
	err := co.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	// TODO - Send a notification to the message broker
	return nil
}
