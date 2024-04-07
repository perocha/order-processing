package usecase

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/domain"
)

type OrderProcessImpl struct {
	orderRepo domain.OrderRepository
}

// Initialize the order process
func NewOrderProcess(orderRepo domain.OrderRepository) *OrderProcessImpl {
	return &OrderProcessImpl{
		orderRepo: orderRepo,
	}
}

// Process the "create_order" event
func (orderImpl *OrderProcessImpl) ProcessCreateOrder(ctx context.Context, order domain.Order) error {
	log.Printf("CreateOrder::Execute::OrderID=%v", order.OrderID)

	// Create the order in the order repository
	err := orderImpl.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	// TODO - Send a notification to the message broker
	return nil
}

func (orderImpl *OrderProcessImpl) ProcessDeleteOrder(ctx context.Context, orderID string) error {
	log.Printf("DeleteOrder::Execute::order=%v", orderID)
	return nil
}

func (orderImpl *OrderProcessImpl) ProcessUpdateOrder(ctx context.Context, order domain.Order) error {
	log.Printf("UpdateOrder::Execute::OrderID=%v", order.OrderID)
	return nil
}
