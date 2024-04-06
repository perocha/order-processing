package usecase

import "github.com/perocha/order-processing/pkg/domain"

type OrderInteractor struct {
	// Define your dependencies such as repositories
}

func (oi *OrderInteractor) ProcessOrder(order domain.Order) {
	// Define your business logic for processing an order
}
