package order

// Order definition
type Order struct {
	Id              string
	ProductCategory string
	ProductID       string
	CustomerID      string
	Status          string
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
