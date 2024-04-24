package order

import (
	"encoding/json"
	"errors"
)

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

func NewOrder() *Order {
	id := ""
	productCategory := ""
	productID := ""
	customerID := ""
	status := ""

	return &Order{
		Id:              id,
		ProductCategory: productCategory,
		ProductID:       productID,
		CustomerID:      customerID,
		Status:          status,
	}
}

func (e *Order) Deserialize(data []byte) error {
	err := json.Unmarshal(data, e)
	if err != nil {
		return errors.New("order::deserialize::error deserializing Order")
	}
	return nil
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
