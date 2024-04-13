package order

// Order definition
type Order struct {
	OrderID    string
	ProductID  string
	CustomerID string
	Status     string
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
