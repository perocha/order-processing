package order

import (
	"testing"
)

func TestToMap(t *testing.T) {
	o := &Order{
		Id:              "TestID",
		ProductCategory: "TestCategory",
		ProductID:       "TestProductID",
		CustomerID:      "TestCustomerID",
		Status:          "TestStatus",
	}

	m := o.ToMap()

	if m["Id"] != o.Id {
		t.Errorf("Expected Id to be %s, but got %s", o.Id, m["Id"])
	}

	if m["ProductCategory"] != o.ProductCategory {
		t.Errorf("Expected ProductCategory to be %s, but got %s", o.ProductCategory, m["ProductCategory"])
	}

	if m["ProductID"] != o.ProductID {
		t.Errorf("Expected ProductID to be %s, but got %s", o.ProductID, m["ProductID"])
	}

	if m["CustomerID"] != o.CustomerID {
		t.Errorf("Expected CustomerID to be %s, but got %s", o.CustomerID, m["CustomerID"])
	}

	if m["Status"] != o.Status {
		t.Errorf("Expected Status to be %s, but got %s", o.Status, m["Status"])
	}
}
