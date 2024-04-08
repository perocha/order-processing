package order

import (
	"context"

	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

type OrderProcessImpl struct {
	orderRepo OrderRepository
}

// Initialize the order process
func NewOrderProcess(orderRepo OrderRepository) *OrderProcessImpl {
	return &OrderProcessImpl{
		orderRepo: orderRepo,
	}
}

// Process the "create_order" event
func (orderImpl *OrderProcessImpl) ProcessCreateOrder(ctx context.Context, order Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	properties := order.ToMap()
	telemetryClient.TrackTrace(ctx, "ProcessCreateOrder", telemetry.Information, properties, true)

	// Create the order in the order repository
	err := orderImpl.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	// TODO - Send a notification to the message broker
	return nil
}

func (orderImpl *OrderProcessImpl) ProcessDeleteOrder(ctx context.Context, orderID string) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	properties := map[string]string{
		"OrderID": orderID,
	}
	telemetryClient.TrackTrace(ctx, "ProcessDeleteOrder", telemetry.Information, properties, true)
	return nil
}

func (orderImpl *OrderProcessImpl) ProcessUpdateOrder(ctx context.Context, order Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	properties := order.ToMap()
	telemetryClient.TrackTrace(ctx, "ProcessUpdateOrder", telemetry.Information, properties, true)
	return nil
}
