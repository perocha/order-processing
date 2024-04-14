package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Order struct {
	Id              string `json:"id"`
	ProductCategory string `json:"productcategory"`
	ProductID       string `json:"productid"`
	CustomerID      string `json:"customerid"`
	Status          string `json:"status"`
}

type Event struct {
	Type         string `json:"Type"`
	OrderPayload Order  `json:"OrderPayload"`
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <event_type> <filename>", os.Args[0])
	}

	eventType := os.Args[1]
	filename := os.Args[2]

	order := Order{}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	err = json.Unmarshal(data, &order)
	if err != nil {
		log.Fatalf("Error unmarshalling order data: %s", err)
	}

	event := Event{
		Type:         eventType,
		OrderPayload: order,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Error occurred while encoding order data: %s", err)
	}

	// Log the event
	log.Printf("Publishing event: %s", jsonData)

	resp, err := http.Post("http://localhost:8080/publish", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error occurred while sending post request: %s", err)
	}

	defer resp.Body.Close()
}
