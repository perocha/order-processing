package service

import (
	"github.com/perocha/order-processing/pkg/service/eventprocessor"
	"github.com/perocha/order-processing/pkg/service/orderservice"
)

// Service aggregates all service interfaces.
type Service interface {
	// Add any necessary methods for the service
}

// ServiceImpl is a struct implementing the Service interface.
type ServiceImpl struct {
	// Embed the individual service interfaces
	eventprocessor.EventProcessor
	orderservice.OrderService
}

// NewService creates a new instance of ServiceImpl.
func NewService(
	eventProcessor eventprocessor.EventProcessor,
	orderService orderservice.OrderService,
) *ServiceImpl {
	// Add any necessary initialization logic
	return &ServiceImpl{
		EventProcessor: eventProcessor,
		OrderService:   orderService,
	}
}
