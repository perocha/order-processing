package main

import (
	"log"

	"github.com/perocha/order-processing/internal/app/infrastructure/config"
)

func main() {
	log.Println("Order Processor is starting...")

	cfg := config.NewConfig()

	log.Println("Event Hub Connection String:", cfg.EventHubConnectionString)
	log.Println("Cosmos DB Connection String:", cfg.CosmosDBConnectionString)
}
