package order

import "encoding/json"

// Order definition
type Order struct {
	Id              string `json:"id"`
	ProductCategory string `json:"ProductCategory"`
	ProductID       string `json:"productId"`
	CustomerID      string `json:"customerId"`
	Status          string `json:"status"`
}

// Convert Order struct into a map[string]string
func (e *Order) ToMap() map[string]string {
	return map[string]string{
		"Id":              e.Id,
		"ProductCategory": e.ProductCategory,
		"ProductID":       e.ProductID,
		"CustomerID":      e.CustomerID,
		"Status":          e.Status,
	}
}

// Convert Order struct into a JSON
func (e *Order) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Return the Order ID
func (e *Order) GetOrderID() string {
	return e.Id
}

// Return order id as a key/value pair
func (e *Order) GetOrderIDMap() map[string]string {
	return map[string]string{
		"OrderId": e.Id,
	}
}
