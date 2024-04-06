package application

// Define interfaces and types related to application services

// OrderService represents the interface for order-related operations.
type OrderService interface {
	ProcessOrder(orderID string) error
	// Add other order-related methods as needed
}
