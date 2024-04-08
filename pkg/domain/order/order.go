package order

import "context"

// Order definition
type Order struct {
	OrderID    string
	ProductID  string
	CustomerID string
	Status     string
}

// Order methods
type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) error
	UpdateOrder(ctx context.Context, order Order) error
	DeleteOrder(ctx context.Context, orderID string) error
}

// Convert Order struct into a map[string]string
func (e *Order) ToMap() map[string]string {
	return map[string]string{
		"OrderID":    e.OrderID,
		"ProductID":  e.ProductID,
		"CustomerID": e.CustomerID,
		"Status":     e.Status,
	}
}
