package event

import (
	"testing"
	"time"

	"github.com/perocha/order-processing/pkg/domain/order"
)

func TestToMap(t *testing.T) {
	e := &Event{
		Type:      "TestType",
		EventID:   "TestID",
		Timestamp: time.Now(),
		OrderPayload: order.Order{
			Id:              "TestID",
			ProductCategory: "TestProductCategory",
			ProductID:       "TestProductID",
			CustomerID:      "TestCustomerID",
			Status:          "TestStatus",
		},
	}

	m := e.ToMap()

	if m["Type"] != e.Type {
		t.Errorf("Expected Type to be %s, but got %s", e.Type, m["Type"])
	}

	if m["EventID"] != e.EventID {
		t.Errorf("Expected EventID to be %s, but got %s", e.EventID, m["EventID"])
	}

	if m["Timestamp"] != e.Timestamp.Format(time.RFC3339) {
		t.Errorf("Expected Timestamp to be %s, but got %s", e.Timestamp.Format(time.RFC3339), m["Timestamp"])
	}

	if m["OrderPayload"] != e.OrderPayload.Id {
		t.Errorf("Expected OrderPayload to be %s, but got %s", e.OrderPayload.Id, m["OrderPayload"])
	}

	if m["OrderPayload.ProductCategory"] != e.OrderPayload.ProductCategory {
		t.Errorf("Expected OrderPayload.ProductCategory to be %s, but got %s", e.OrderPayload.ProductCategory, m["OrderPayload.ProductCategory"])
	}

	if m["OrderPayload.ProductID"] != e.OrderPayload.ProductID {
		t.Errorf("Expected OrderPayload.ProductID to be %s, but got %s", e.OrderPayload.ProductID, m["OrderPayload.ProductID"])
	}

	if m["OrderPayload.CustomerID"] != e.OrderPayload.CustomerID {
		t.Errorf("Expected OrderPayload.CustomerID to be %s, but got %s", e.OrderPayload.CustomerID, m["OrderPayload.CustomerID"])
	}

	if m["OrderPayload.Status"] != e.OrderPayload.Status {
		t.Errorf("Expected OrderPayload.Status to be %s, but got %s", e.OrderPayload.Status, m["OrderPayload.Status"])
	}
}
