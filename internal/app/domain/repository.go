package domain

// OrderRepository represents the interface for order data operations.
type OrderRepository interface {
	GetOrderByID(orderID string) (*Order, error)
	SaveOrder(order *Order) error
	// Add other order-related methods as needed
}
